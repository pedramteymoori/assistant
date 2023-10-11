package main

import "github.com/housinganywhere/assistant/internal/assistant"

func main() {
	asi, err := assistant.NewAssistant()
	if err != nil {
		panic(err)
	}
	asi.Start()
}
