package chatgpt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Model struct {
	Id         string       `json:"id,omitempty"`
	Object     string       `json:"object,omitempty"`
	Created    int64        `json:"created,omitempty"`
	OwnedBy    string       `json:"owned_by,omitempty"`
	Root       string       `json:"root,omitempty"`
	Permission []Permission `json:"permission,omitempty"`
	// "parent": null
}

type Permission struct {
	Id                  string `json:"id,omitempty"`
	Object              string `json:"object,omitempty"`
	Created             int64  `json:"created,omitempty"`
	Organization        string `json:"organization,omitempty"`
	AllowCreateEngine   bool   `json:"allow_create_engine,omitempty"`
	AllowSampling       bool   `json:"allow_sampling,omitempty"`
	AllowLogprobs       bool   `json:"allow_logprobs,omitempty"`
	AllowSearch_indices bool   `json:"allow_search_indices,omitempty"`
	AllowView           bool   `json:"allow_view,omitempty"`
	AllowFineTuning     bool   `json:"allow_fine_tuning,omitempty"`
	IsBlocking          bool   `json:"is_blocking,omitempty"`
	// group
}

type ModelsRes struct {
	Oobject string  `json:"object,omitempty"`
	Data    []Model `json:"data,omitempty"`
}

func (client *Client) Models(ctx context.Context) (res *ModelsRes, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	request, _ = http.NewRequest("GET", client.modelsUrl(), nil)
	request.WithContext(ctx)
	client.setAuth(request, true)
	// fmt.Println("~~~", request)

	if response, err = client.cli.Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(response.Status)
	}

	res = new(ModelsRes)
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}
