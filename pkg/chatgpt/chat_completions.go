package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatCompMsg struct {
	// enum: user, system, assistant
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompReq struct {
	Model       string        `json:"model"`
	Temperature float32       `json:"temperature"`
	Messages    []ChatCompMsg `json:"messages"`
}

type ChatCompChoice struct {
	Message      ChatCompMsg `json:"message,omitempty"`
	FinishReason string      `json:"finish_reason,omitempty"`
	Index        uint        `json:"index,omitempty"`
}

type ChatCompRes struct {
	Id      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	Model   string `json:"model,omitempty"`
	Usage   struct {
		PromptTokens     uint32 `json:"prompt_tokens,omitempty"`
		CompletionTokens uint32 `json:"completion_tokens,omitempty"`
		TotalTokens      uint32 `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`
	Choices []ChatCompChoice `json:"choices,omitempty"`
}

func (req *ChatCompReq) Validate() (err error) {
	if req.Model == "" {
		req.Model = defaultChatModel()
	}
	if req.Temperature == 0 {
		req.Temperature = defaultChatTemperature()
	}
	if len(req.Messages) == 0 {
		return fmt.Errorf("empty messages")
	}

	return nil
}

func (client *Client) ChatCompletions(ctx context.Context, req *ChatCompReq) (res *ChatCompRes, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	_ = encoder.Encode(req)

	request, _ = http.NewRequest("POST", client.chatCompletionsUrl(), buf)
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

	res = new(ChatCompRes)
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}
