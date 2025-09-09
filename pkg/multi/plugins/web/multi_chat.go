package web

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/multi/plugins/web/schema"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// BroadcastMessage sends a message from an agent to all connected clients for that agent.
func (s *Server) BroadcastMessage(info *shared.AgentInfo, content string) {
	if info == nil {
		s.logger.Error("BroadcastMessage called with nil agent info")
		return
	}

	event := &schema.LLMEvent{
		ID:           uuid.New().String(),
		InvocationID: uuid.New().String(),
		Author:       info.Name,
		Timestamp:    time.Now().Unix(),
		Type:         schema.EventTypeAssistant,
		Done:         true,
		Partial:      false,
		Role:         model.RoleAssistant,
		Parts: []schema.Part{
			&schema.TextPart{Content: content},
		},
	}

	s.logger.Info("Broadcasting agent message",
		zap.String("agent", info.Name),
		zap.String("content", content),
	)

	s.ssePool.BroadcastToAgent(info.ID(), event)
}

// BroadcastToolCall sends a tool call from an agent to all connected clients for that agent.
func (s *Server) BroadcastToolCall(info *shared.AgentInfo, functionDef model.FunctionDefinitionParam) {
	if info == nil {
		s.logger.Error("BroadcastToolCall called with nil agent info")
		return
	}

	// Parse arguments for structured storage
	var args interface{}
	if err := json.Unmarshal(functionDef.Arguments, &args); err != nil {
		args = string(functionDef.Arguments) // fallback to string
	}

	event := &schema.LLMEvent{
		ID:           uuid.New().String(),
		InvocationID: uuid.New().String(),
		Author:       info.Name,
		Timestamp:    time.Now().Unix(),
		Type:         schema.EventTypeToolCall,
		Done:         true,
		Partial:      false,
		Role:         model.RoleTool,
		Parts: []schema.Part{
			&schema.FunctionCallPart{
				Name: functionDef.Name,
				Args: args,
			},
		},
	}

	s.logger.Info("Broadcasting agent tool call",
		zap.String("agent", info.Name),
		zap.String("toolName", functionDef.Name),
	)
	s.ssePool.BroadcastToAgent(info.ID(), event)
}

// InterAgentEvent represents an inter-agent communication event
type InterAgentEvent struct {
	Type      string    `json:"type"`      // Always "inter_agent"
	FromAgent uuid.UUID `json:"fromAgent"` // Source agent name
	ToAgent   uuid.UUID `json:"toAgent"`   // Target agent name
	Message   string    `json:"message"`   // Message content
	Timestamp time.Time `json:"timestamp"` // Event timestamp
}

// WithChatProcessor sets the multi-agent chat processor for the server.
// It also wires up the necessary callbacks for the processor to communicate back to the web server.
func WithChatProcessor(processor multi.ChatProcessor) Option {
	return func(s *Server) {
		// Define the callback for regular messages.
		onMessageCallback := func(info *shared.AgentInfo, content string) {
			s.BroadcastMessage(info, content)
		}

		// Define the callback for tool calls.
		onToolCallCallback := func(info *shared.AgentInfo, functionDef model.FunctionDefinitionParam) {
			s.BroadcastToolCall(info, functionDef)
		}

		// Set the callbacks on the processor instance.
		processor.SetOnMessageCallback(onMessageCallback)
		processor.SetOnToolCallCallback(onToolCallCallback)

		// Now, set the fully wired-up processor on the server.
		s.chatProcessor = processor
		s.setupInterAgentInterceptor()
	}
}

// setupInterAgentInterceptor configures the message interceptor for inter-agent communication
func (s *Server) setupInterAgentInterceptor() {
	if s.chatProcessor == nil {
		return
	}

	s.chatProcessor.SetMessageInterceptor(func(fromAgentID, toAgentID uuid.UUID, content string) {
		if fromAgentID != uuid.Nil && toAgentID != uuid.Nil {
			interAgentEvent := s.createInterAgentEvent(fromAgentID, toAgentID, content)

			// Broadcast to all active SSE connections
			s.broadcastInterAgentEvent(interAgentEvent)
		}
	})
}

// createInterAgentEvent creates a modern LLMEvent for inter-agent communication
func (s *Server) createInterAgentEvent(fromAgent, toAgent uuid.UUID, message string) *schema.LLMEvent {
	return schema.NewInterAgentEvent(
		fromAgent.String(),
		toAgent.String(),
		message,
		schema.InterAgentCommunication,
	)
}

// createReceivedInterAgentEvent creates an event for the receiving agent
func (s *Server) createReceivedInterAgentEvent(fromAgent, toAgent uuid.UUID, originalEvent *schema.LLMEvent) *schema.LLMEvent {
	// Extract message from original event
	var message string
	if len(originalEvent.Parts) > 0 {
		if textPart, ok := originalEvent.Parts[0].(*schema.TextPart); ok {
			message = textPart.Content
		}
	}

	receivedMessage := fmt.Sprintf("Received from %s: %s", fromAgent, message)
	return schema.NewInterAgentEvent(
		fromAgent.String(),
		toAgent.String(),
		receivedMessage,
		schema.InterAgentReceived,
	)
}

// broadcastInterAgentEvent broadcasts an inter-agent event to both sender and receiver agents
func (s *Server) broadcastInterAgentEvent(event *schema.LLMEvent) {
	if event.InterAgent == nil {
		s.logger.Error("Invalid inter-agent event: missing InterAgent data")
		return
	}

	fromAgent, err := uuid.Parse(event.InterAgent.FromAgent)
	if err != nil {
		s.logger.Error("Invalid fromAgent UUID", zap.String("uuid", event.InterAgent.FromAgent))
		return
	}

	toAgent, err := uuid.Parse(event.InterAgent.ToAgent)
	if err != nil {
		s.logger.Error("Invalid toAgent UUID", zap.String("uuid", event.InterAgent.ToAgent))
		return
	}

	s.logger.Info("Broadcasting inter-agent event",
		zap.String("fromAgent", fromAgent.String()),
		zap.String("toAgent", toAgent.String()),
	)

	senderCount := s.ssePool.BroadcastToAgent(fromAgent, event)
	receivedEvent := s.createReceivedInterAgentEvent(fromAgent, toAgent, event)
	receiverCount := s.ssePool.BroadcastToAgent(toAgent, receivedEvent)

	totalSent := senderCount + receiverCount
	if totalSent == 0 {
		s.logger.Info("No active SSE connections for inter-agent event",
			zap.String("fromAgent", fromAgent.String()),
			zap.String("toAgent", toAgent.String()),
		)
	} else {
		s.logger.Info("Sent inter-agent event",
			zap.Int("totalSent", totalSent),
			zap.Int("senderCount", senderCount),
			zap.Int("receiverCount", receiverCount),
		)
	}
}

// GetMultiChatAgents returns the list of available agents for multi-agent chat
func (s *Server) GetMultiChatAgents() []*shared.AgentInfo {
	if s.chatProcessor == nil {
		return nil
	}
	return s.chatProcessor.GetAllAgentInfos()
}
