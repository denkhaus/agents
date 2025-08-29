package time

import (
	"context"
	"time"

	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const (
	ToolName = "current_time"
)

// timeArgs holds the input for the time tool.
type timeArgs struct {
	Timezone string `json:"timezone" description:"Timezone (UTC, EST, PST, CST) or leave empty for local"`
}

// timeResult holds the output for the time tool.
type timeResult struct {
	Timezone string `json:"timezone"`
	Time     string `json:"time"`
	Date     string `json:"date"`
	Weekday  string `json:"weekday"`
}

// Time tool implementation.
// getCurrentTime returns the current time for the specified timezone.
// If the timezone is invalid or empty, it defaults to local time.
func getCurrentTime(ctx context.Context, args timeArgs) (timeResult, error) {
	loc := time.Local
	zone := args.Timezone
	// Attempt to load the specified timezone.
	if zone != "" {
		var err error
		loc, err = time.LoadLocation(zone)
		if err != nil {
			loc = time.Local
		}
	}
	now := time.Now().In(loc)
	return timeResult{
		Timezone: loc.String(),
		Time:     now.Format("15:04:05"),
		Date:     now.Format("2006-01-02"),
		Weekday:  now.Weekday().String(),
	}, nil
}

func NewTool(i *do.Injector) (tool.Tool, error) {
	// Create time tool for timezone queries.
	timeTool := function.NewFunctionTool(
		getCurrentTime,
		function.WithName(ToolName),
		function.WithDescription(
			"Get the current time and date for a specific timezone",
		),
	)

	return timeTool, nil
}
