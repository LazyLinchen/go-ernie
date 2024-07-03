package go_ernie

import (
	"context"
	"fmt"
	"testing"
)

var (
	ak = "ak"
	sk = "sk"
)

func TestClient_CreateChatCompletionStream(t *testing.T) {
	client, err := NewClient(ak, sk)
	if err != nil {
		t.Error(err)
		return
	}
	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Model: "ERNIE-4.0-Turbo-8K",
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "你好!",
			},
		},
		Stream: true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("%+v\n", resp)
	}
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(ak, sk)
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := client.CreateChatCompletion(context.Background(), ChatCompletionRequest{
		Model: "ERNIE-4.0-Turbo-8K",
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

func TestClient_CreateCompletion(t *testing.T) {
	client, err := NewClient(ak, sk)
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := client.CreateCompletion(context.Background(), CompletionRequest{
		Model:       "SQLCoder-7B",
		Prompt:      "你是一个通古晓今的历史老师，我需要你给我讲解一下秦始皇的故事。",
		Temperature: 0.5,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v\n", resp)
}
