package stable_diffusion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type InterrogateReq struct {
	Image string `json:"image"`
	Model string `json:"model"`
}

type InterrogateRes struct {
	Caption string `json:"caption"`
}

func (client *Client) Interrogate_url() string {
	return fmt.Sprintf("%s%s", client.config.Url, _Data.GetString("apis.Interrogate"))
}

func (req *InterrogateReq) Validate() (err error) {
	if req.Image == "" {
		return fmt.Errorf("empty image")
	}

	if req.Model == "" {
		return fmt.Errorf("empty model")
	}

	return nil
}

// Interrogate
func (client *Client) Interrogate(ctx context.Context, req *InterrogateReq) (
	res *InterrogateRes, err error) {
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

	request, _ = http.NewRequest("POST", client.Interrogate_url(), buf)
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
	// fmt.Println("~~~", string(bts))
	res = new(InterrogateRes)
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (client *Client) InterrogateCLIP(fp string) (res *InterrogateRes, err error) {
	var (
		b64 string
		req *InterrogateReq
	)

	if b64, err = Base64File(fp); err != nil {
		return nil, err
	}

	req = &InterrogateReq{Image: b64, Model: "clip"}

	return client.Interrogate(context.Background(), req)
}

func (client *Client) InterrogateDeepDanBooru(fp string) (res *InterrogateRes, err error) {
	var (
		b64 string
		req *InterrogateReq
	)

	if b64, err = Base64File(fp); err != nil {
		return nil, err
	}

	req = &InterrogateReq{Image: b64, Model: "deepdanbooru"}

	return client.Interrogate(context.Background(), req)
}
