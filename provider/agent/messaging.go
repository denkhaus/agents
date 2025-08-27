package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// Message represents a message between agents
type Message struct {
	ID        string
	From      uuid.UUID
	To        uuid.UUID
	Content   string
	Timestamp time.Time
}

// MessageBroker handles routing messages between agents
type MessageBroker struct {
	mu          sync.RWMutex
	agents      map[uuid.UUID]agent.Agent
	channels    map[uuid.UUID]chan *Message
	interceptor func(fromID, toID uuid.UUID, content string)
}

// NewMessageBroker creates a new message broker
func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		agents:   make(map[uuid.UUID]agent.Agent),
		channels: make(map[uuid.UUID]chan *Message),
	}
}

// RegisterAgent registers an agent with the broker
func (mb *MessageBroker) RegisterAgent(agent agent.Agent) uuid.UUID {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	id := uuid.New()
	mb.agents[id] = agent
	mb.channels[id] = make(chan *Message, 100) // Buffered channel

	return id
}

// RegisterAgentWithID registers an agent with a predefined ID
func (mb *MessageBroker) RegisterAgentWithID(agentID uuid.UUID, agent agent.Agent) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.agents[agentID] = agent
	mb.channels[agentID] = make(chan *Message, 100) // Buffered channel
}

// UnregisterAgent removes an agent from the broker
func (mb *MessageBroker) UnregisterAgent(agentID uuid.UUID) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	delete(mb.agents, agentID)
	if ch, exists := mb.channels[agentID]; exists {
		close(ch)
		delete(mb.channels, agentID)
	}
}

// SendMessage sends a message from one agent to another
func (mb *MessageBroker) SendMessage(from, to uuid.UUID, content string) error {
	mb.mu.RLock()
	interceptor := mb.interceptor
	mb.mu.RUnlock()

	// Call interceptor if set (for displaying messages in chat)
	if interceptor != nil {
		interceptor(from, to, content)
	}

	mb.mu.RLock()
	defer mb.mu.RUnlock()

	// Check if recipient exists
	if _, exists := mb.agents[to]; !exists {
		return fmt.Errorf("agent %s not found", to)
	}

	// Check if channel exists
	ch, exists := mb.channels[to]
	if !exists {
		return fmt.Errorf("channel for agent %s not found", to)
	}

	// Create and send message
	message := &Message{
		ID:        uuid.New().String(),
		From:      from,
		To:        to,
		Content:   content,
		Timestamp: time.Now(),
	}

	// Non-blocking send with timeout
	select {
	case ch <- message:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout sending message to %s", to)
	}
}

// GetMessageChannel returns the message channel for an agent
func (mb *MessageBroker) GetMessageChannel(agentID uuid.UUID) (<-chan *Message, error) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	ch, exists := mb.channels[agentID]
	if !exists {
		return nil, fmt.Errorf("channel for agent %s not found", agentID)
	}

	return ch, nil
}

// ListAgentIDs returns a list of registered agent IDs
func (mb *MessageBroker) ListAgentIDs() []uuid.UUID {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	ids := make([]uuid.UUID, 0, len(mb.agents))
	for id := range mb.agents {
		ids = append(ids, id)
	}

	return ids
}

// MessagingWrapper wraps any agent.Agent to add messaging capabilities
type MessagingWrapper struct {
	agent.Agent
	broker *MessageBroker
	id     uuid.UUID
}

// NewMessagingWrapper creates a new messaging wrapper for an agent
func NewMessagingWrapper(baseAgent agent.Agent, broker *MessageBroker) *MessagingWrapper {
	// Create wrapper with a temporary ID
	wrapper := &MessagingWrapper{
		Agent:  baseAgent,
		broker: broker,
		id:     uuid.New(),
	}

	// Register with broker and get the actual ID
	wrapper.id = broker.RegisterAgent(wrapper)

	return wrapper
}

// NewMessagingWrapperWithID creates a new messaging wrapper with a predefined ID
func NewMessagingWrapperWithID(baseAgent agent.Agent, broker *MessageBroker, agentID uuid.UUID) *MessagingWrapper {
	// Create wrapper with predefined ID
	wrapper := &MessagingWrapper{
		Agent:  baseAgent,
		broker: broker,
		id:     agentID,
	}

	// Register with broker using the predefined ID
	broker.RegisterAgentWithID(agentID, wrapper)

	return wrapper
}

// ID returns the unique ID of this agent
func (mw *MessagingWrapper) ID() uuid.UUID {
	return mw.id
}

// SendMessage sends a message to another agent by ID
func (mw *MessagingWrapper) SendMessage(to uuid.UUID, content string) error {
	return mw.broker.SendMessage(mw.id, to, content)
}

// GetMessageChannel returns the message channel for this agent
func (mw *MessagingWrapper) GetMessageChannel() (<-chan *Message, error) {
	return mw.broker.GetMessageChannel(mw.id)
}

// Run implements the agent.Agent interface
func (mw *MessagingWrapper) Run(ctx context.Context, invocation *agent.Invocation) (<-chan *event.Event, error) {
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
func (mw *MessagingWrapper) messageToEvent(msg *Message) *event.Event {
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
func (mw *MessagingWrapper) Info() agent.Info {
	info := mw.Agent.Info()
	info.Description = fmt.Sprintf("%s (with messaging capabilities, ID: %s)", info.Description, mw.id)
	return info
}

// Tools implements the agent.Agent interface
func (mw *MessagingWrapper) Tools() []tool.Tool {
	// Get tools from the base agent
	baseTools := mw.Agent.Tools()
	
	// Add our messaging tool
	messagingTool := NewMessagingTool(mw.broker, mw.id)
	
	// Convert to the expected tool type
	tools := make([]tool.Tool, 0, len(baseTools)+1)
	for _, t := range baseTools {
		tools = append(tools, t)
	}
	tools = append(tools, messagingTool)
	
	return tools
}

// MessagingTool is a tool that allows agents to send messages
type MessagingTool struct {
	broker  *MessageBroker
	agentID uuid.UUID
}

// NewMessagingTool creates a new messaging tool
func NewMessagingTool(broker *MessageBroker, agentID uuid.UUID) *MessagingTool {
	return &MessagingTool{
		broker:  broker,
		agentID: agentID,
	}
}

// Declaration returns the tool declaration
func (mt *MessagingTool) Declaration() *tool.Declaration {
	return &tool.Declaration{
		Name:        "send_message",
		Description: "Send a message to another agent by ID",
		InputSchema: &tool.Schema{
			Type: "object",
			Properties: map[string]*tool.Schema{
				"to": {
					Type:        "string",
					Description: "The UUID of the recipient agent",
				},
				"content": {
					Type:        "string",
					Description: "The message content",
				},
			},
			Required: []string{"to", "content"},
		},
	}
}

// Call executes the tool
func (mt *MessagingTool) Call(ctx context.Context, jsonArgs []byte) (any, error) {
	// Parse the arguments
	var args map[string]interface{}
	if err := json.Unmarshal(jsonArgs, &args); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	toStr, ok := args["to"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'to' parameter")
	}

	to, err := uuid.Parse(toStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'to' parameter: %w", err)
	}

	content, ok := args["content"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'content' parameter")
	}

	err = mt.broker.SendMessage(mt.agentID, to, content)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return map[string]interface{}{
		"status":  "sent",
		"to":      to.String(),
		"content": content,
	}, nil
}

// SetMessageInterceptor sets a function to intercept all messages
func (mb *MessageBroker) SetMessageInterceptor(interceptor func(fromID, toID uuid.UUID, content string)) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.interceptor = interceptor
}