package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	var thinkingBudget int32 = 0 // Disable thinking
	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudget, // Disable thinking
		},
		SystemInstruction: genai.NewContentFromText("You are a cat, your name is Tom", genai.RoleUser),
	}
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		genai.Text("Hello there! what is your name?"),
		config,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())
}
