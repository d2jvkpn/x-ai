package chatgpt

import (
	"context"
	"fmt"
)

func TransDecorate(content, language string) (string, error) {
	switch language {
	case English:
		return fmt.Sprintf("Translate the following into English:\n%q", content), nil
	case Chinese:
		return fmt.Sprintf("Translate the following into Chinese:\n%q", content), nil
	default:
		return "", fmt.Errorf("unknown languge")
	}
}

func (client *Client) Ask(ctx context.Context, content string) (ans string, err error) {
	var (
		req *ChatCompReq
		res *ChatCompRes
	)

	req = &ChatCompReq{Messages: []ChatCompMsg{{Role: "user", Content: content}}}
	if res, err = client.ChatCompletions(ctx, req); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return res.Choices[0].Message.Content, nil
}

func (client *Client) Trans2En(ctx context.Context, content string) (ans string, err error) {
	var (
		req *ChatCompReq
		res *ChatCompRes
	)

	input, _ := TransDecorate(content, English)
	req = &ChatCompReq{Messages: []ChatCompMsg{{Role: "user", Content: input}}}
	if res, err = client.ChatCompletions(ctx, req); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return res.Choices[0].Message.Content, nil
}

func (client *Client) Trans2Cn(ctx context.Context, content string) (ans string, err error) {
	var (
		req *ChatCompReq
		res *ChatCompRes
	)

	input, _ := TransDecorate(content, Chinese)
	req = &ChatCompReq{Messages: []ChatCompMsg{{Role: "user", Content: input}}}
	if res, err = client.ChatCompletions(ctx, req); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return res.Choices[0].Message.Content, nil
}
