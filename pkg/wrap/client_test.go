package wrap

import (
	"context"
	"fmt"
	"log"
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func TestClient(t *testing.T) {
	/*
		config := openai.DefaultConfig("OPENAI_KEY")
		fmt.Printf("~~~ config: %s, %s\n", config.BaseURL, config.OrgID)
		fmt.Printf("~~~ config: %s, %s, %s\n", config.APIType, config.APIVersion, config.Engine)
	*/

	var (
		err    error
		client *openai.Client
	)

	if client, err = NewOpenAiClient("../../configs/local.yaml", "chatgpt"); err != nil {
		log.Fatalln(err)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
