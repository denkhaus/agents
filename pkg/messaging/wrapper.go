package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/denkhaus/agents/pkg/shared"
	agentinfo "github.com/denkhaus/agents/pkg/tools/agent_info"
	sendmessage "github.com/denkhaus/agents/pkg/tools/send_message"
	"github.com/google/uuid"
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
	tools           []tool.Tool
}

// NewAgentWrapper creates a new messaging wrapper around a base agent to enable inter-agent communication
func NewAgentWrapper(
	baseAgent shared.TheAgent,
	broker MessageBroker,
	availableAgents ...*shared.AgentInfo,
) (shared.TheAgent, error) {

	// Create wrapper with predefined ID
	wrapper := &messagingWrapper{
		TheAgent:        baseAgent,
		availableAgents: availableAgents,
		broker:          broker,
	}

	// Register with broker using the predefined ID
	broker.RegisterAgent(baseAgent.GetID(), wrapper)

	if err := wrapper.assembleTools(); err != nil {
		return nil, err
	}

	return wrapper, nil
}

// Info implements the agent.Agent interface
func (mw *messagingWrapper) Info() agent.Info {
	return mw.TheAgent.Info()
}

func (mw *messagingWrapper) Tools() []tool.Tool {
	return mw.tools
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

func (mw *messagingWrapper) appendToolSafe(tool tool.Tool) error {
	for _, t := range mw.tools {
		currentToolName := t.Declaration().Name
		if currentToolName == tool.Declaration().Name {
			return fmt.Errorf("the %q tool already exists in the tool collection of %q agent",
				currentToolName, mw.TheAgent.GetInfo().Name,
			)
		}
	}

	mw.tools = append(mw.tools, tool)
	return nil
}

// Tools implements the agent.Agent interface
func (mw *messagingWrapper) assembleTools() error {

	mw.tools = mw.TheAgent.Tools()

	messagingTool, err := sendmessage.New(mw.broker, mw.GetID(), mw.GetAllowedToCommunicateWith())
	if err != nil {
		return fmt.Errorf("failed to create %q tool: %w", sendmessage.ToolName, err)
	}

	err = mw.appendToolSafe(messagingTool)
	if err != nil {
		return err
	}

	agentInfoTool, err := agentinfo.New(mw.availableAgents, mw.GetID(), mw.GetAllowedToCommunicateWith())
	if err != nil {
		return fmt.Errorf("failed to create %q tool: %w", agentinfo.ToolName, err)
	}

	err = mw.appendToolSafe(agentInfoTool)
	if err != nil {
		return err
	}

	return nil
}
