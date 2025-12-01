package multitoolagent

import (
	"fmt"
	"strings"
	"time"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

var cityTimeZones = map[string]string{
	"new york":    "America/New_York",
	"los angeles": "America/Los_Angeles",
	"chicago":     "America/Chicago",
	"london":      "Europe/London",
	"tokyo":       "Asia/Tokyo",
	"paris":       "Europe/Paris",
}

type timeToolInput struct {
	City string `json:"city"`
}

type timeToolOutput struct {
	Message string `json:"message"`
}

func getCurrentTime(ctx tool.Context, input timeToolInput) (timeToolOutput, error) {
	timeZone, exists := cityTimeZones[strings.ToLower(input.City)]
	if !exists {
		return timeToolOutput{}, fmt.Errorf("no timezone information for the city %v", input.City)
	}
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return timeToolOutput{}, fmt.Errorf("failed to load location for timezone %v: %v", timeZone, err)
	}
	currentTime := time.Now().In(location).Format(time.UnixDate)
	return timeToolOutput{Message: fmt.Sprintf("The current time in %s is %s", input.City, currentTime)}, nil
}

func NewTimeTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "get_current_time",
		Description: "Get the current time for a specified city.",
	}, getCurrentTime)
}
