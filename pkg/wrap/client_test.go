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

	req := openai.ChatCompletionRequest{
		Model:    "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{{Role: "user", Content: "who are you?"}},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)

	if err != nil {
		t.Fatalf("%v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
