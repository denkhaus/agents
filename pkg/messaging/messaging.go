package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/shared"
	agentinfo "github.com/denkhaus/agents/pkg/tools/agent_info"
	sendmessage "github.com/denkhaus/agents/pkg/tools/send_message"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// MessagingWrapper wraps any agent.Agent to add messaging capabilities
type messagingWrapper struct {
	shared.TheAgent
	broker          MessageBroker
	availableAgents []*shared.AgentInfo
}

// NewWrapper creates a new messaging wrapper with a predefined ID
func NewWrapper(baseAgent shared.TheAgent, broker MessageBroker, availableAgents ...*shared.AgentInfo) shared.TheAgent {
	// Create wrapper with predefined ID
	wrapper := &messagingWrapper{
		TheAgent:        baseAgent,
		availableAgents: availableAgents,
		broker:          broker,
	}

	// Register with broker using the predefined ID
	broker.RegisterAgent(baseAgent.GetID(), wrapper)

	return wrapper
}

// SendMessage sends a message to another agent by ID
func (mw *messagingWrapper) SendMessage(to uuid.UUID, content string) error {
	return mw.broker.SendMessage(mw.GetID(), to, content)
}

// GetMessageChannel returns the message channel for this agent
func (mw *messagingWrapper) GetMessageChannel() (<-chan *Message, error) {
	return mw.broker.GetMessageChannel(mw.GetID())
}

// Run implements the agent.Agent interface
func (mw *messagingWrapper) Run(ctx context.Context, invocation *agent.Invocation) (<-chan *event.Event, error) {
	// Get the base agent's event channel
	baseEventChan, err := mw.TheAgent.Run(ctx, invocation)
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
		Author:       msg.FromAgentID.String(),
		ID:           msg.ID,
		Timestamp:    msg.Timestamp,
	}
}

// Info implements the agent.Agent interface
func (mw *messagingWrapper) Info() agent.Info {
	return mw.TheAgent.Info()
}

// Tools implements the agent.Agent interface
func (mw *messagingWrapper) assembleTools() ([]tool.Tool, error) {
	// Get tools from the base agent
	baseTools := mw.TheAgent.Tools()

	var tools []tool.Tool
	tools = append(tools, baseTools...)

	// Add our messaging tool
	messagingTool, err := sendmessage.New(mw.broker, mw.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to create %q tool: %w", sendmessage.ToolName, err)
	}

	tools = append(tools, messagingTool)

	// Add our agent info tool
	agentInfoTool, err := agentinfo.New(mw.availableAgents, mw.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to create %q tool: %w", agentinfo.ToolName, err)
	}

	tools = append(tools, agentInfoTool)

	return tools, nil
}

func (mw *messagingWrapper) Tools() []tool.Tool {
	tools, err := mw.assembleTools()
	if err != nil {
		logger.Log.Error("failed to assemble tools for wrapped agend", zap.Error(err))
	}

	return tools
}
