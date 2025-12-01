package multitoolagent

import (
	"fmt"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

var cityWeather = map[string]string{
	"new york":    "Sunny, 75°F",
	"los angeles": "Cloudy, 68°F",
	"chicago":     "Rainy, 60°F",
	"london":      "Foggy, 55°F",
	"tokyo":       "Clear, 80°F",
	"paris":       "Windy, 65°F",
}

type weatherToolInput struct {
	City string `json:"city"`
}

type weatherToolOutput struct {
	Message string `json:"message"`
}

func getWeather(ctx tool.Context, input weatherToolInput) (weatherToolOutput, error) {
	weather, exists := cityWeather[strings.ToLower(input.City)]
	if !exists {
		return weatherToolOutput{}, fmt.Errorf("no weather information for the city %v", input.City)
	}
	return weatherToolOutput{Message: fmt.Sprintf("The weather in %s is %s", input.City, weather)}, nil
}

func NewWeatherTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "get_weather",
		Description: "Get weather for a specified city.",
	}, getWeather)
}
