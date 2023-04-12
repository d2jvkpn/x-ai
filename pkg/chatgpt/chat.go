package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatMsg struct {
	// enum: user, system, assistant
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatReq struct {
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature"`
	Messages    []ChatMsg `json:"messages"`
}

type ChatChoice struct {
	Message      ChatMsg `json:"message,omitempty"`
	FinishReason string  `json:"finish_reason,omitempty"`
	Index        uint    `json:"index,omitempty"`
}

type ChatRes struct {
	Id      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	Model   string `json:"model,omitempty"`
	Usage   struct {
		PromptTokens     uint32 `json:"prompt_tokens,omitempty"`
		CompletionTokens uint32 `json:"completion_tokens,omitempty"`
		TotalTokens      uint32 `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`
	Choices []ChatChoice `json:"choices,omitempty"`
}

func (req *ChatReq) Validate() (err error) {
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

func (client *Client) Chat(ctx context.Context, req *ChatReq) (res *ChatRes, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	_ = encoder.Encode(req)

	request, _ = http.NewRequest("POST", client.chat_url(), buf)
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

	res = new(ChatRes)
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (client *Client) Ask(ctx context.Context, content string) (ans string, err error) {
	var (
		req *ChatReq
		res *ChatRes
	)

	req = &ChatReq{Messages: []ChatMsg{{Role: "user", Content: content}}}
	if res, err = client.Chat(ctx, req); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return res.Choices[0].Message.Content, nil
}

func (client *Client) Trans2En(ctx context.Context, content string) (ans string, err error) {
	var (
		req *ChatReq
		res *ChatRes
	)

	input, _ := TransDecorate(content, English)
	req = &ChatReq{Messages: []ChatMsg{{Role: "user", Content: input}}}
	if res, err = client.Chat(ctx, req); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return res.Choices[0].Message.Content, nil
}

func (client *Client) Trans2Cn(ctx context.Context, content string) (ans string, err error) {
	var (
		req *ChatReq
		res *ChatRes
	)

	input, _ := TransDecorate(content, Chinese)
	req = &ChatReq{Messages: []ChatMsg{{Role: "user", Content: input}}}
	if res, err = client.Chat(ctx, req); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return res.Choices[0].Message.Content, nil
}
