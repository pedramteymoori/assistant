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
	Tenant    Tenant    `json:"tenant"`
	Landlord  LandLord  `json:"landlord"`
	Questions Questions `json:"questions"`
}

type Tenant struct {
	Score      float32    `json:"score"`
	Suggestion Suggestion `json:"suggestion"`
}

type LandLord struct {
	Score      float32    `json:"score"`
	Suggestion Suggestion `json:"suggestion"`
}

type Questions struct {
	PropertyRent bool `json:"property_rent"`
	Deposit      bool `json:"deposit"`
	Amenities    bool `json:"amenities"`
	Payment      bool `json:"payment"`
}

type Suggestion struct {
	Text string `json:"text"`
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
		Model: chatgpt.GPT4,
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
		Questions: Questions{
			Deposit:      response.Questions.Deposit,
			PropertyRent: response.Questions.PropertyRent,
			Payment:      response.Questions.Payment,
			Amenities:    response.Questions.Amenities,
		},
	}, nil
}
