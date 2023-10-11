package ai

import (
	"context"
	"encoding/json"
	"os"

	chatgpt "github.com/ayush6624/go-chatgpt"
)

type AI struct {
	client *chatgpt.Client
}

type AIResponse struct {
	Tenant   int `json:"tenant"`
	Landlord int `json:"landlord"`
}

func NewAI() (*AI, error) {

	key := os.Getenv("OPENAI_KEY")
	client, err := chatgpt.NewClient(key)
	if err != nil {
		return nil, err
	}
	return &AI{
		client: client,
	}, nil
}

func (a *AI) ProvideTips(conversations string) (*AIResponse, error) {
	ctx := context.Background()

	res, err := a.client.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT35Turbo,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: conversations,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var response AIResponse
	err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &response)
	if err != nil {
		return nil, err
	}

	return &AIResponse{
		Tenant:   response.Tenant,
		Landlord: response.Landlord,
	}, nil
}
