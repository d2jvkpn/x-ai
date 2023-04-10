package stable_diffusion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ControlnetUnit struct {
	InputImage string  // base64
	Module     string  // "depth"
	Model      string  // "control_sd15_depth [fef5e48e]"
	Guidance   float64 // 1.00
}

type CNTxt2ImgReq struct {
	Txt2ImgReq
	ControlnetUnits []ControlnetUnit
}

type CNImg2ImgReq struct {
	Txt2ImgReq
	ControlnetUnits []ControlnetUnit
}

func (client *Client) cn_txt2img_url() string {
	return fmt.Sprintf("%s%s", client.config.Url, _Data.GetString("apis.controlnet_txt2img"))
}

func (client *Client) cn_img2img_url() string {
	return fmt.Sprintf("%s%s", client.config.Url, _Data.GetString("apis.controlnet_img2img"))
}

func (client *Client) CNTxt2Img(ctx context.Context, req *CNTxt2ImgReq) (res *Txt2ImgRes, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	if err = req.Validate(); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	_ = encoder.Encode(req)

	request, _ = http.NewRequest("POST", client.cn_txt2img_url(), buf)
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

func (client *Client) CNImg2Img(ctx context.Context, req *CNImg2ImgReq) (res *Img2ImgRes, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	if err = req.Validate(); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	_ = encoder.Encode(req)

	request, _ = http.NewRequest("POST", client.cn_img2img_url(), buf)
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
