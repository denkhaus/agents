package messaging

import (
	"fmt"
	"sync"
	"time"

	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/shared/resource"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

type RoutingInfo struct {
	FromAgentID uuid.UUID `json:"from_agent_id"`
	ToAgentID   uuid.UUID `json:"to_agent_id"`
	SessionID   uuid.UUID `json:"session_id"`
	Streaming   *bool     `json:"streaming"`
}

func (p *RoutingInfo) String() string {
	return fmt.Sprintf("[%s]->[%s]", p.FromAgentID, p.ToAgentID)
}

// Message represents a message between agents
type Message struct {
	RoutingInfo
	ID        string
	Content   string
	Timestamp time.Time
}

type Interceptor func(routing *RoutingInfo, content string)

type MessageBroker interface {
	shared.MessageSender
	RegisterAgent(agentID uuid.UUID, agent shared.TheAgent)
	UnregisterAgent(agentID uuid.UUID)
	GetMessageChannel(agentID uuid.UUID) (<-chan *Message, error)
	SetMessageInterceptor(interceptor Interceptor)
	ListAgentIDs() []uuid.UUID
}

// messageBrokerImpl handles routing messages between agents
type messageBrokerImpl struct {
	mu          sync.RWMutex
	sessionID   uuid.UUID
	agents      *resource.Manager[shared.TheAgent]
	channels    *resource.Manager[chan *Message]
	interceptor func(routing *RoutingInfo, content string)
}

type UsageMetaData struct {
	PromptTokenCount     int `json:"prompt_token_count,omitempty"`
	CandidatesTokenCount int `json:"candidates_token_count,omitempty"`
	TotalTokenCount      int `json:"total_token_count,omitempty"`
}

type Part interface {
}

type EventPartType string

const (
	EventPartTypeText             EventPartType = "text"
	EventPartTypeTextFragment     EventPartType = "text-fragment"
	EventPartTypeFunctionCall     EventPartType = "function-call"
	EventPartTypeFunctionResponse EventPartType = "function-response"
	EventPartTypeinterAgent       EventPartType = "inter-agent"
)

type TextPart struct {
	Type    EventPartType `json:"type,omitempty"`
	Content string        `json:"content,omitempty"`
}

type FunctionCallPart struct {
	Type EventPartType `json:"type,omitempty"`
	Name string        `json:"name,omitempty"`
	Args interface{}   `json:"args,omitempty"`
	ID   string        `json:"id,omitempty"`
}

type FunctionResponsePart struct {
	Type     EventPartType `json:"type,omitempty"`
	Name     string        `json:"name,omitempty"`
	Args     interface{}   `json:"args,omitempty"`
	ID       string        `json:"id,omitempty"`
	Response interface{}   `json:"response,omitempty"`
}

type InterAgentPart struct {
	Type        EventPartType `json:"type,omitempty"`
	FromAgentID string        `json:"from_agent_id,omitempty"`
	ToAgentID   string        `json:"to_agent_id,omitempty"`
	Message     string        `json:"message,omitempty"`
	Direction   string        `json:"direction,omitempty"`
}

type InterAgentEventType string

const (
	InterAgentCommunication InterAgentEventType = "communication"
	InterAgentReceived      InterAgentEventType = "received"
)

type LLMEvent struct {
	base             *event.Event         `json:"-"`
	Routing          *RoutingInfo         `json:"routing,omitempty"`
	Usage            *UsageMetaData       `json:"usage,omitempty"`
	Done             bool                 `json:"done,omitempty"`
	Partial          bool                 `json:"partial,omitempty"`
	Type             EventType            `json:"type,omitempty"`
	Created          time.Time            `json:"created,omitempty"`
	CreatedTimestamp int64                `json:"created_ts,omitempty"`
	Model            string               `json:"model,omitempty"`
	Role             model.Role           `json:"role,omitempty"`
	Parts            []Part               `json:"parts,omitempty"`
	ID               uuid.UUID            `json:"id,omitempty"`
	InvocationID     uuid.UUID            `json:"invocation_id,omitempty"`
	Author           string               `json:"author,omitempty"`
	InterAgent       *InterAgentEventType `json:"interagent_type,omitempty"`
}
