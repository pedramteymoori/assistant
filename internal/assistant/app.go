package assistant

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	Tenant   ResponseMessage `json:"tenant"`
	Landlord ResponseMessage `json:"landlord"`
}

type ResponseMessage struct {
	Text string `json:"text"`
	Slug string `json:"slug"`
	//Link string `json:"link"`
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
	conv := "We are looking out for a tone which is not rude, abusive, disrespectful and making someone feel uncomfortable. The tone should be such that it helps the people involved in the conversation trust each other and make them feel secured\nConsidering these aspects can you please help me highlight whether the below conversation has acceptable tone or not? Can you please provide an average tone rating for tenant and landlord for the below conversation between the scale of 1 to 5, where 1 is acceptable and 5 is non-acceptable in json format\n"

	for _, row := range req.Messages {
		conv = fmt.Sprintf("%s%s:%s\n", conv, row.User, row.Message)
	}

	return conv
}

func (a *Assistant) provideResponse(aiResp *ai.AIResponse) Response {

	//var tenantText string
	//
	//if aiResp.Tenant > 8 {
	//	tenantText = "good job"
	//} else {
	//
	//}

	return Response{
		Tenant: ResponseMessage{
			Text: fmt.Sprintf("%d", aiResp.Tenant),
			Slug: "hi",
		},
		Landlord: ResponseMessage{
			Text: fmt.Sprintf("%d", aiResp.Landlord),
			Slug: "hi",
		},
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
