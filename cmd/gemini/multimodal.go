package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	imagePath, err := filepath.Abs("assets/images/gvisor.png")
	if err != nil {
		log.Fatal(err)
	}
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	parts := []*genai.Part{
		genai.NewPartFromText("What is in this image?"),
		&genai.Part{
			InlineData: &genai.Blob{
				Data:     imageData,
				MIMEType: "image/png",
			},
		},
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	var thinkingBudget int32 = 0 // Disable thinking
	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudget, // Disable thinking
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		contents,
		config,
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result.Text())
}
