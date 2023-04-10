package stable_diffusion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Img2ImgReq struct {
	Prompt         string   `json:"prompt,omitempty"`
	NegativePrompt string   `json:"negative_prompt,omitempty"`
	InitImages     []string `json:"init_images,omitempty"`
	// ResizeMode     uint     `json:"resize_mode,omitempty"` // default 0, Just resize
	SamplerName       string  `json:"sampler_name,omitempty"`       // default "Euler a", sampling method
	Steps             uint    `json:"steps,omitempty"`              // default 20, sampling steps
	RestoreFaces      bool    `json:"restore_faces,omitempty"`      // default: false, Restore faces
	Tiling            bool    `json:"tiling,omitempty"`             // default: false, Tiling
	Width             uint32  `json:"width,omitempty"`              // default 512
	Height            uint32  `json:"height,omitempty"`             // default 512
	CfgScale          float64 `json:"cfg_scale,omitempty"`          // default 7.0
	DenoisingStrength float64 `json:"denoising_strength,omitempty"` // default 0.75
	BatchSize         uint32  `json:"batch_size,omitempty"`         // default 1
	// NIter          uint32 `json:"n_iter,omitempty"`     // default 1
	AlwaysonScripts map[string]any `json:"alwayson_scripts,omitempty"` // extensions
}

type Img2ImgRes struct {
	Images []string `json:"images,omitempty"`
	// Parameters map[string]any   `json:"parameters,omitempty"`
	Info string `json:"info,omitempty"`
}

func (req *Img2ImgReq) Validate() (err error) {
	// if len(req.InitImages) == 0 {
	// 	return fmt.Errorf("empty init_images")
	// }

	if req.Width == 0 {
		req.Width = 512
	}
	if req.Height == 0 {
		req.Height = 512
	}
	if req.BatchSize == 0 {
		req.BatchSize = 1
	}

	return nil
}

func (client *Client) img2img_url() string {
	return fmt.Sprintf("%s%s", client.config.Url, _Data.GetString("apis.img2img"))
}

func (client *Client) Img2Img(ctx context.Context, req *Img2ImgReq, exts ...Extension) (
	res *Img2ImgRes, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	if err = req.Validate(); err != nil {
		return nil, err
	}

	if len(exts) > 0 && req.AlwaysonScripts == nil {
		req.AlwaysonScripts = make(map[string]any, len(exts))
	}
	for i := range exts {
		req.AlwaysonScripts[exts[i].Key()] = exts[i].Value()
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	_ = encoder.Encode(req)

	request, _ = http.NewRequest("POST", client.img2img_url(), buf)
	request.WithContext(ctx)
	if response, err = client.cli.Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case 422:
		bts, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("%s: %s", response.Status, bts)
	default:
		return nil, fmt.Errorf(response.Status)
	}

	res = new(Img2ImgRes)
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}
