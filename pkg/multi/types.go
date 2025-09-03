package multi

import (
	"context"

	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/event"
)

// ChatProcessor defines the interface for managing multi-agent chat interactions.
// It provides methods for sending messages between agents, retrieving agent information,
// and setting up message interception for monitoring communication.
type ChatProcessor interface {
	// SetMessageInterceptor sets a function to intercept and monitor messages between agents.
	// The interceptor receives the sender ID, receiver ID, and message content.
	SetMessageInterceptor(interceptor messaging.Interceptor)

	// SendMessage sends a message from one agent to another and returns a channel of events.
	// The caller is responsible for processing the events from the returned channel.
	SendMessage(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) (<-chan *event.Event, error)

	// SendMessageWithProcessing sends a message and automatically processes all resulting events.
	// This is a convenience method that handles event processing internally.
	SendMessageWithProcessing(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) error

	// GetAgentInfoByAuthor retrieves agent information by author name or UUID string.
	// Returns nil if no agent is found with the given identifier.
	GetAgentInfoByAuthor(author string) *shared.AgentInfo

	// GetAllAgentInfos returns information for all registered agents in the chat processor.
	GetAllAgentInfos() []shared.AgentInfo

	// GetAgentNameByID returns the name of an agent given its UUID.
	// Returns empty string if no agent is found with the given ID.
	GetAgentNameByID(agentID uuid.UUID) string

	// GetAgentByName returns the actual agent instance by name.
	// Returns nil if no agent is found with the given name.
	GetAgentByName(name string) shared.TheAgent

	// SetOnMessageCallback sets the message callback function for the ChatProcessor.
	SetOnMessageCallback(onMessage OnMessage)
	// SetOnToolCallCallback sets the tool call callback function for the ChatProcessor.
	SetOnToolCallCallback(onToolCall OnToolCall)
}
