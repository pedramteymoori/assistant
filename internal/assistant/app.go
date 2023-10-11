package assistant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/housinganywhere/assistant/internal/ai"
)

type Assistant struct {
	ai *ai.AI
}

type Request struct {
	Messages []RequestMessage `json:"messages"`
}

type RequestMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type Response struct {
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
	Slug string `json:"slug"`
}

type Suggestion struct {
	Text string `json:"text"`
}

func NewAssistant() (*Assistant, error) {
	aiInstance, err := ai.NewAI()
	if err != nil {
		return nil, err
	}

	return &Assistant{
		ai: aiInstance,
	}, nil
}

func (a *Assistant) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/assist", a.getSuggestions)

	err := http.ListenAndServe(":8080", mux)
	panic(err)
}

func (a *Assistant) getSuggestions(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var req Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conv := a.prepareMessage(req)

	aiResp, err := a.ai.ProvideTips(conv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := a.provideResponse(aiResp)

	marshalled, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalled)
}

func (a *Assistant) prepareMessage(req Request) string {
	re := regexp.MustCompile(`\r?\n`)

	conv := "We are looking out for a tone which is not rude, abusive, disrespectful and making someone feel uncomfortable. The tone should be such that it helps the people involved in the conversation trust each other and make them feel secured\nConsidering these aspects can you please help me highlight whether the below conversation has acceptable tone or not? \n Can you please provide an average tone rating for tenant and landlord for the below conversation between the scale of 1 to 5, where 1 is acceptable and 5 is non-acceptable\nAlso, in case the score is towards unacceptable, please provide some suggestion how landlord and tenant can improve\nDoes the last message of the Tenant contain question around property-rental, deposit, payment or amenities?  if the last message is from landlord then please return false for all of these area.\nCan you please provide json structure like:\n{\n    \"tenant\": {\n        \"score\": 4,\n        \"suggestion\": {\n            \"text\": \"XYZ\"\n        }\n    },\n    \"landlord\": {\n        \"score\": 4,\n        \"suggestion\": {\n            \"text\": \"ABC\"\n        }\n    },\n    \"questions\": {\n        \"property_rent\": false,\n        \"deposit\": false,\n        \"amenities\": false,\n        \"payment\": true   \n    }\n}\n\n\n"
	for _, row := range req.Messages {
		userMessage := re.ReplaceAllString(row.Message, " ")
		conv = fmt.Sprintf("%s%s: %s\n", conv, row.User, userMessage)
	}

	return conv
}

func (a *Assistant) provideResponse(aiResp *ai.AIResponse) Response {
	slug := ""
	if aiResp.Questions.Amenities {
		slug = "amenities"
	} else if aiResp.Questions.PropertyRent {
		slug = "property-rent"
	} else if aiResp.Questions.Deposit {
		slug = "deposit"
	} else if aiResp.Questions.Payment {
		slug = "payment"
	}

	res := Response{
		Tenant: Tenant{
			Score: aiResp.Tenant.Score,
			Suggestion: Suggestion{
				Text: aiResp.Tenant.Suggestion.Text,
			},
		},
		Landlord: LandLord{
			Score: aiResp.Landlord.Score,
			Suggestion: Suggestion{
				Text: aiResp.Landlord.Suggestion.Text,
			},
		},
		Questions: Questions{
			slug,
		},
	}

	return res
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
