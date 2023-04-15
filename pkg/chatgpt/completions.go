package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CompMsg struct {
	// enum: user, system, assistant
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompReq struct {
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature"`
	Messages    []CompMsg `json:"messages"`
}

type CompChoice struct {
	Message      CompMsg `json:"message,omitempty"`
	FinishReason string  `json:"finish_reason,omitempty"`
	Index        uint    `json:"index,omitempty"`
}

type CompRes struct {
	Id      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	Model   string `json:"model,omitempty"`
	Usage   struct {
		PromptTokens     uint32 `json:"prompt_tokens,omitempty"`
		CompletionTokens uint32 `json:"completion_tokens,omitempty"`
		TotalTokens      uint32 `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`
	Choices []CompChoice `json:"choices,omitempty"`
}

func (req *CompReq) Validate() (err error) {
	if req.Model == "" {
		req.Model = default_model()
	}
	if req.Temperature == 0 {
		req.Temperature = default_temperature()
	}
	if len(req.Messages) == 0 {
		return fmt.Errorf("empty messages")
	}

	return nil
}

func (client *Client) Completions(ctx context.Context, req *CompReq) (res *CompRes, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	_ = encoder.Encode(req)

	request, _ = http.NewRequest("POST", client.completions_url(), buf)
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

	res = new(CompRes)
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}
