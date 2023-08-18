package go_ernie

import (
	"context"
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("akk", "akkk")
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := client.CreateChatCompletion(context.Background(), ChatCompletionRequest{
		Model: ErnieBot,
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "你好!",
			},
		},
	})

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Result)
}
