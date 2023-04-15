package chatgpt

import (
	"context"
	"fmt"
)

func (client *Client) Ask(ctx context.Context, content string) (ans string, err error) {
	var (
		req *CompReq
		res *CompRes
	)

	req = &CompReq{Messages: []CompMsg{{Role: "user", Content: content}}}
	if res, err = client.Completions(ctx, req); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return res.Choices[0].Message.Content, nil
}

func (client *Client) Trans2En(ctx context.Context, content string) (ans string, err error) {
	var (
		req *CompReq
		res *CompRes
	)

	input, _ := TransDecorate(content, English)
	req = &CompReq{Messages: []CompMsg{{Role: "user", Content: input}}}
	if res, err = client.Completions(ctx, req); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return res.Choices[0].Message.Content, nil
}

func (client *Client) Trans2Cn(ctx context.Context, content string) (ans string, err error) {
	var (
		req *CompReq
		res *CompRes
	)

	input, _ := TransDecorate(content, Chinese)
	req = &CompReq{Messages: []CompMsg{{Role: "user", Content: input}}}
	if res, err = client.Completions(ctx, req); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return res.Choices[0].Message.Content, nil
}
