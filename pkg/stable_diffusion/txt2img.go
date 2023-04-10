package stable_diffusion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Txt2ImgReq struct {
	Prompt         string `json:"prompt,omitempty"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	SamplerName    string `json:"sampler_name,omitempty"`  // default "Euler a", sampling method
	Steps          uint   `json:"steps,omitempty"`         // default 20, sampling steps
	RestoreFaces   bool   `json:"restore_faces,omitempty"` // default false, Restore faces
	Tiling         bool   `json:"tiling,omitempty"`        // default false, Tiling
	// EnableHr bool `json:"enable_hr,omitempty"` // default false, Hires. fix
	Width     uint32 `json:"width,omitempty"`      // default 512
	Height    uint32 `json:"height,omitempty"`     // default 512
	BatchSize uint32 `json:"batch_size,omitempty"` // default 1
	// NIter          uint32 `json:"n_iter,omitempty"`     // default 1
	CfgScale        float64        `json:"cfg_scale,omitempty"`        // default 7.0
	AlwaysonScripts map[string]any `json:"alwayson_scripts,omitempty"` // extensions
}

type Txt2ImgRes struct {
	Images []string `json:"images,omitempty"`
	// Parameters map[string]any   `json:"parameters,omitempty"`
	Info string `json:"info,omitempty"`
}

func (req *Txt2ImgReq) Validate() (err error) {
	if req.Prompt == "" {
		return fmt.Errorf("empty prompt")
	}

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

func (client *Client) txt2img_url() string {
	return fmt.Sprintf("%s%s", client.config.Url, _Data.GetString("apis.txt2img"))
}

func (client *Client) Txt2Img(ctx context.Context, req *Txt2ImgReq, exts ...Extension) (
	res *Txt2ImgRes, err error) {
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

	request, _ = http.NewRequest("POST", client.txt2img_url(), buf)
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

	// bts, _ := ioutil.ReadAll(response.Body)
	// fmt.Printf("~~~ response body: %s\n", bts)

	res = new(Txt2ImgRes)
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (client *Client) Txt2ImgPromot(prompt string, batchSize uint32) (res *Txt2ImgRes, err error) {
	req := &Txt2ImgReq{Prompt: prompt, BatchSize: batchSize}
	return client.Txt2Img(context.Background(), req)
}
