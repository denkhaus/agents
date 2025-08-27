package messaging

import (
	"sync"
	"time"

	"github.com/denkhaus/agents/shared/resource"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/agent"
)

// Message represents a message between agents
type Message struct {
	ID        string
	From      uuid.UUID
	To        uuid.UUID
	Content   string
	Timestamp time.Time
}

type Interceptor func(fromID, toID uuid.UUID, content string)

type MessageBroker interface {
	RegisterAgent(agentID uuid.UUID, agent agent.Agent)
	UnregisterAgent(agentID uuid.UUID)
	GetMessageChannel(agentID uuid.UUID) (<-chan *Message, error)
	SetMessageInterceptor(interceptor Interceptor)
	SendMessage(from, to uuid.UUID, content string) error
	ListAgentIDs() []uuid.UUID
}

// messageBrokerImpl handles routing messages between agents
type messageBrokerImpl struct {
	mu          sync.RWMutex
	agents      *resource.Manager[agent.Agent]
	channels    *resource.Manager[chan *Message]
	interceptor func(fromID, toID uuid.UUID, content string)
}
