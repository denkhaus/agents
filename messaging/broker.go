package messaging

import (
	"fmt"
	"time"

	"github.com/denkhaus/agents/shared/resource"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
)

// NewMessageBroker creates a new message broker
func NewMessageBroker() MessageBroker {
	return &messageBrokerImpl{
		agents:   resource.NewManager[agent.Agent](),
		channels: resource.NewManager[chan *Message](),
	}
}

// RegisterAgent registers an agent with a predefined ID
func (mb *messageBrokerImpl) RegisterAgent(agentID uuid.UUID, agent agent.Agent) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.agents.Set(agentID, agent)
	mb.channels.Set(agentID, make(chan *Message, 100)) // Buffered channel
}

// UnregisterAgent removes an agent from the broker
func (mb *messageBrokerImpl) UnregisterAgent(agentID uuid.UUID) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.agents.Delete(agentID)
	if ch, exists := mb.channels.Get(agentID); exists {
		close(ch)
		mb.channels.Delete(agentID)
	}
}

// SendMessage sends a message from one agent to another
func (mb *messageBrokerImpl) SendMessage(from, to uuid.UUID, content string) error {
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
	if !mb.agents.Exists(to) {
		return fmt.Errorf("agent %s not found", to)
	}

	// Check if channel exists
	ch, exists := mb.channels.Get(to)
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
func (mb *messageBrokerImpl) GetMessageChannel(agentID uuid.UUID) (<-chan *Message, error) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	ch, exists := mb.channels.Get(agentID)
	if !exists {
		return nil, fmt.Errorf("channel for agent %s not found", agentID)
	}

	return ch, nil
}

// ListAgentIDs returns a list of registered agent IDs
func (mb *messageBrokerImpl) ListAgentIDs() []uuid.UUID {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	agents := mb.agents.GetAll()
	ids := make([]uuid.UUID, 0, len(agents))
	for id := range agents {
		ids = append(ids, id)
	}

	return ids
}

// SetMessageInterceptor sets a function to intercept all messages
func (mb *messageBrokerImpl) SetMessageInterceptor(interceptor Interceptor) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.interceptor = interceptor
}
