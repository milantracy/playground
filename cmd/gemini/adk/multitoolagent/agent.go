package main

import (
	"context"
	"os"

	"github.com/milantracy/playground/cmd/gemini/adk"
	"github.com/milantracy/playground/pkg/gemini/adk/multitoolagent"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

// https://google.github.io/adk-docs/get-started/quickstart/#build-a-multi-tool-agent

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, adk.GeminiFlashModel, &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})

	if err != nil {
		panic(err)
	}

	timeTool, err := multitoolagent.NewTimeTool()
	if err != nil {
		panic(err)
	}
	weatherTool, err := multitoolagent.NewWeatherTool()
	if err != nil {
		panic(err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "weather_time_agent",
		Model:       model,
		Description: "Agent to answer questions about the time and weather in a city.",
		Instruction: `
        You are the **Weather time agent**. Your **sole function** is to determine and state the current time and weather for a specified city or location using the provided tools.
        
        **MANDATORY RULES:**
        1. Always be concise and provide only the current time and city.
        2. If the user asks about your capabilities (e.g., "What can you do?"), you must respond only with: "I am the time weather Agent. I can tell you the current time and weather in any city you specify."
        3. Do not engage in general conversation or answer questions unrelated to time.
		`,
		Tools: []tool.Tool{
			timeTool,
			weatherTool,
		},
	})
	if err != nil {
		panic(err)
	}
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(a),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		panic(err)
	}
}
