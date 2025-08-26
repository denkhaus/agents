package chat

import (
	"context"
	"fmt"
	"strings"

	"github.com/denkhaus/agents/provider"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/runner"
)

type State struct {
	fullContent       string
	toolCallsDetected bool
	assistantStarted  bool
	currentAgent      string
}

func (p *State) Reset() {
	p.currentAgent = ""
	p.fullContent = ""
	p.toolCallsDetected = false
	p.assistantStarted = false
}

type chatImpl struct {
	State
	provider.ChatProviderOptions
	runner    runner.Runner
	streaming bool
}

func NewChat(agent agent.Agent, streaming bool, options provider.ChatProviderOptions) provider.Chat {
	runnerOptions := []runner.Option{}
	if options.SessionService != nil {
		runnerOptions = append(runnerOptions,
			runner.WithSessionService(
				options.SessionService,
			),
		)
	}

	impl := &chatImpl{
		streaming:           streaming,
		ChatProviderOptions: options,
		runner: runner.NewRunner(
			options.AppName,
			agent,
			runnerOptions...,
		),
	}

	return impl
}

// processMessage handles a single message exchange.
func (c *chatImpl) ProcessMessage(ctx context.Context, userMessage string) error {
	c.Reset()

	message := model.NewUserMessage(userMessage)

	// Run the agent through the runner.
	eventChan, err := c.runner.Run(ctx, c.UserID, c.SessionID, message)
	if err != nil {
		return fmt.Errorf("failed to run agent: %w", err)
	}

	// Process response.
	return c.processResponse(eventChan)
}

// processResponse handles both streaming and non-streaming responses with tool call visualization.
func (c *chatImpl) processResponse(eventChan <-chan *event.Event) error {
	fmt.Print("🤖 Assistant: ")

	for event := range eventChan {
		if err := c.handleEvent(event); err != nil {
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
func (c *chatImpl) handleEvent(event *event.Event) error {
	// Handle errors.
	if event.Error != nil {
		fmt.Printf("\n❌ Error: %s\n", event.Error.Message)
		return nil
	}

	c.handleAgentTransition(event)

	// Handle tool calls.
	if c.handleToolCalls(event) {
		return nil
	}

	// Handle tool responses.
	if c.handleToolResponses(event) {
		return nil
	}

	// Handle content.
	c.handleContent(event)

	return nil
}

// handleAgentTransition manages agent switching and display.
func (c *chatImpl) handleAgentTransition(event *event.Event) {
	if event.Author != c.currentAgent {
		if c.assistantStarted {
			fmt.Printf("\n")
		}

		c.currentAgent = event.Author
		c.assistantStarted = true
		c.toolCallsDetected = false

		fmt.Printf("[%s]:\n", c.currentAgent)
	}
}

// handleToolCalls detects and displays tool calls.
func (c *chatImpl) handleToolCalls(event *event.Event) bool {
	if len(event.Choices) > 0 && len(event.Choices[0].Message.ToolCalls) > 0 {
		c.toolCallsDetected = true
		if c.assistantStarted {
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
func (c *chatImpl) handleContent(event *event.Event) {
	if len(event.Choices) > 0 {
		choice := event.Choices[0]
		content := c.extractContent(choice)

		if content != "" {
			c.displayContent(content)
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
func (c *chatImpl) displayContent(content string) {
	if !c.assistantStarted {
		if c.toolCallsDetected {
			fmt.Printf("\n🤖 Assistant ")
		}
		c.assistantStarted = true
	}

	fmt.Print(content)
	c.fullContent += content
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
