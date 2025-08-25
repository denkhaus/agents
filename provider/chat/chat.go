package chat

import (
	"context"
	"fmt"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/runner"
)

type chatImpl struct {
	Options
	runner runner.Runner
}

func NewChat(options Options) Chat {
	runnerOptions := []runner.Option{}
	if options.sessionService != nil {
		runnerOptions = append(runnerOptions,
			runner.WithSessionService(
				options.sessionService,
			),
		)
	}

	impl := &chatImpl{
		Options: options,
		runner: runner.NewRunner(
			options.appName,
			options.agent,
			runnerOptions...,
		),
	}

	return impl
}

// processMessage handles a single message exchange.
func (c *chatImpl) ProcessMessage(ctx context.Context, userMessage string) error {
	message := model.NewUserMessage(userMessage)

	// Run the agent through the runner.
	eventChan, err := c.runner.Run(ctx, c.userID, c.sessionID, message)
	if err != nil {
		return fmt.Errorf("failed to run agent: %w", err)
	}

	// Process response.
	return c.processResponse(eventChan)
}

// processResponse handles both streaming and non-streaming responses with tool call visualization.
func (c *chatImpl) processResponse(eventChan <-chan *event.Event) error {
	fmt.Print("🤖 Assistant: ")

	var (
		fullContent       string
		toolCallsDetected bool
		assistantStarted  bool
	)

	for event := range eventChan {
		if err := c.handleEvent(event, &toolCallsDetected, &assistantStarted, &fullContent); err != nil {
			return err
		}

		// Check if this is the final event.
		// Don't break on tool response events (Done=true but not final assistant response).
		if event.Done && !c.isToolEvent(event) {
			fmt.Printf("\n")
			break
		}
	}

	return nil
}

// handleEvent processes a single event from the event channel.
func (c *chatImpl) handleEvent(
	event *event.Event,
	toolCallsDetected *bool,
	assistantStarted *bool,
	fullContent *string,
) error {
	// Handle errors.
	if event.Error != nil {
		fmt.Printf("\n❌ Error: %s\n", event.Error.Message)
		return nil
	}

	// Handle tool calls.
	if c.handleToolCalls(event, toolCallsDetected, assistantStarted) {
		return nil
	}

	// Handle tool responses.
	if c.handleToolResponses(event) {
		return nil
	}

	// Handle content.
	c.handleContent(event, toolCallsDetected, assistantStarted, fullContent)

	return nil
}

// handleToolCalls detects and displays tool calls.
func (c *chatImpl) handleToolCalls(
	event *event.Event,
	toolCallsDetected *bool,
	assistantStarted *bool,
) bool {
	if len(event.Choices) > 0 && len(event.Choices[0].Message.ToolCalls) > 0 {
		*toolCallsDetected = true
		if *assistantStarted {
			fmt.Printf("\n")
		}
		fmt.Printf("🔧 CallableTool calls initiated:\n")
		for _, toolCall := range event.Choices[0].Message.ToolCalls {
			fmt.Printf("   • %s (ID: %s)\n", toolCall.Function.Name, toolCall.ID)
			if len(toolCall.Function.Arguments) > 0 {
				fmt.Printf("     Args: %s\n", string(toolCall.Function.Arguments))
			}
		}
		fmt.Printf("\n🔄 Executing tools...\n")
		return true
	}
	return false
}

// handleToolResponses detects and displays tool responses.
func (c *chatImpl) handleToolResponses(event *event.Event) bool {
	if event.Response != nil && len(event.Response.Choices) > 0 {
		hasToolResponse := false
		for _, choice := range event.Response.Choices {
			if choice.Message.Role == model.RoleTool && choice.Message.ToolID != "" {
				fmt.Printf("✅ CallableTool response (ID: %s): %s\n",
					choice.Message.ToolID,
					strings.TrimSpace(choice.Message.Content))
				hasToolResponse = true
			}
		}
		if hasToolResponse {
			return true
		}
	}
	return false
}

// handleContent processes and displays content.
func (c *chatImpl) handleContent(
	event *event.Event,
	toolCallsDetected *bool,
	assistantStarted *bool,
	fullContent *string,
) {
	if len(event.Choices) > 0 {
		choice := event.Choices[0]
		content := c.extractContent(choice)

		if content != "" {
			c.displayContent(content, toolCallsDetected, assistantStarted, fullContent)
		}
	}
}

// extractContent extracts content based on streaming mode.
func (c *chatImpl) extractContent(choice model.Choice) string {
	if c.streaming {
		// Streaming mode: use delta content.
		return choice.Delta.Content
	}
	// Non-streaming mode: use full message content.
	return choice.Message.Content
}

// displayContent prints content to console.
func (c *chatImpl) displayContent(
	content string,
	toolCallsDetected *bool,
	assistantStarted *bool,
	fullContent *string,
) {
	if !*assistantStarted {
		if *toolCallsDetected {
			fmt.Printf("\n🤖 Assistant: ")
		}
		*assistantStarted = true
	}
	fmt.Print(content)
	*fullContent += content
}

// isToolEvent checks if an event is a tool response (not a final response).
func (c *chatImpl) isToolEvent(event *event.Event) bool {
	if event.Response == nil {
		return false
	}
	if len(event.Choices) > 0 && len(event.Choices[0].Message.ToolCalls) > 0 {
		return true
	}
	if len(event.Choices) > 0 && event.Choices[0].Message.ToolID != "" {
		return true
	}

	// Check if this is a tool response by examining choices.
	for _, choice := range event.Response.Choices {
		if choice.Message.Role == model.RoleTool {
			return true
		}
	}

	return false
}
