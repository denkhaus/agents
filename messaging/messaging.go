package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// MessagingWrapper wraps any agent.Agent to add messaging capabilities
type messagingWrapper struct {
	agent.Agent
	broker MessageBroker
	id     uuid.UUID
}

// NewMessagingWrapper creates a new messaging wrapper with a predefined ID
func NewMessagingWrapper(baseAgent agent.Agent, broker MessageBroker, agentID uuid.UUID) shared.TheAgent {
	// Create wrapper with predefined ID
	wrapper := &messagingWrapper{
		Agent:  baseAgent,
		broker: broker,
		id:     agentID,
	}

	// Register with broker using the predefined ID
	broker.RegisterAgent(agentID, wrapper)

	return wrapper
}

// ID returns the unique ID of this agent
func (mw *messagingWrapper) ID() uuid.UUID {
	return mw.id
}

// SendMessage sends a message to another agent by ID
func (mw *messagingWrapper) SendMessage(to uuid.UUID, content string) error {
	return mw.broker.SendMessage(mw.id, to, content)
}

// GetMessageChannel returns the message channel for this agent
func (mw *messagingWrapper) GetMessageChannel() (<-chan *Message, error) {
	return mw.broker.GetMessageChannel(mw.id)
}

// Run implements the agent.Agent interface
func (mw *messagingWrapper) Run(ctx context.Context, invocation *agent.Invocation) (<-chan *event.Event, error) {
	// Get the base agent's event channel
	baseEventChan, err := mw.Agent.Run(ctx, invocation)
	if err != nil {
		return nil, err
	}

	// Create a new event channel that merges base events with message events
	eventChan := make(chan *event.Event, 256)

	// Get message channel
	msgChan, err := mw.GetMessageChannel()
	if err != nil {
		return nil, err
	}

	// Create a context that we can cancel to stop the goroutine
	mergeCtx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(eventChan)
		defer cancel() // Cancel the context when we're done

		// Merge base events and message events
		for {
			select {
			case <-mergeCtx.Done():
				return
			case baseEvent, ok := <-baseEventChan:
				if !ok {
					// Base event channel closed
					return
				}
				select {
				case eventChan <- baseEvent:
				case <-mergeCtx.Done():
					return
				}
			case msg, ok := <-msgChan:
				if !ok {
					// Message channel closed
					return
				}
				// Convert message to event
				msgEvent := mw.messageToEvent(msg)
				select {
				case eventChan <- msgEvent:
				case <-mergeCtx.Done():
					return
				}
			}
		}
	}()

	return eventChan, nil
}

// messageToEvent converts a message to an event
func (mw *messagingWrapper) messageToEvent(msg *Message) *event.Event {
	// Create a message in the content
	message := model.NewAssistantMessage(msg.Content)

	response := &model.Response{
		Object:    model.ObjectTypeChatCompletion,
		Done:      true,
		Created:   time.Now().Unix(),
		Choices:   []model.Choice{{Message: message}},
		Timestamp: msg.Timestamp,
	}

	return &event.Event{
		Response:     response,
		InvocationID: uuid.New().String(),
		Author:       msg.From.String(),
		ID:           msg.ID,
		Timestamp:    msg.Timestamp,
	}
}

// Info implements the agent.Agent interface
func (mw *messagingWrapper) Info() agent.Info {
	info := mw.Agent.Info()
	info.Description = fmt.Sprintf("%s (with messaging capabilities, ID: %s)", info.Description, mw.id)
	return info
}

// Tools implements the agent.Agent interface
func (mw *messagingWrapper) Tools() []tool.Tool {
	// Get tools from the base agent
	baseTools := mw.Agent.Tools()

	// Add our messaging tool
	messagingTool := NewMessagingTool(mw.broker, mw.id)

	// Convert to the expected tool type
	tools := make([]tool.Tool, 0, len(baseTools)+1)
	tools = append(tools, baseTools...)
	tools = append(tools, messagingTool)

	return tools
}
