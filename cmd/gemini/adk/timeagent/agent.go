package main

import (
	"context"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

const geminiModel = "gemini-2.5-flash"

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, geminiModel, &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		panic(err)
	}

	timeAgent, err := llmagent.New(llmagent.Config{
		Name:        "hello_time_agent",
		Model:       model,
		Description: "Tells the current time in a specified city.",
		Instruction: `
        You are the **Hello Time Agent**. Your **sole function** is to determine and state the current time for a specified city or location using the provided Google Search tool.
        
        **MANDATORY RULES:**
        1. Always be concise and provide only the current time and city.
        2. If the user asks about your capabilities (e.g., "What can you do?"), you must respond only with: "I am the Hello Time Agent. I can tell you the current time in any city you specify."
        3. Do not engage in general conversation or answer questions unrelated to time.
		`,
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		panic(err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(timeAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		panic(err)
	}
}
