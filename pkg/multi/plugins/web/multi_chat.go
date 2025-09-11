package web

import (
	"time"

	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/event"
)

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

		s.chatProcessor = processor
		s.chatProcessor.SetOnRawEventCallback(func(routing *messaging.RoutingInfo, event *event.Event) {

			llmEvent, err := messaging.NewLLMEvent(routing, event)
			if err != nil {
				s.logger.Error("failed to create LLM event", zap.Error(err))
				return
			}

			if llmEvent != nil {
				s.ssePool.BroadcastToAgent(llmEvent)
			}
		})

		// Now, set the fully wired-up processor on the server.
		s.chatProcessor.SetMessageInterceptor(func(routing *messaging.RoutingInfo, content string) {
			if routing.FromAgentID != uuid.Nil && routing.ToAgentID != uuid.Nil {
				interAgentEvent := s.createInterAgentEvent(routing, content)
				// Broadcast to all active SSE connections
				s.broadcastInterAgentEvent(interAgentEvent)
			}
		})
	}
}

// createInterAgentEvent creates a modern LLMEvent for inter-agent communication
func (s *Server) createInterAgentEvent(routing *messaging.RoutingInfo, message string) *messaging.LLMEvent {
	return messaging.NewInterAgentEvent(
		routing,
		message,
		messaging.InterAgentCommunication,
	)
}

// createReceivedInterAgentEvent creates an event for the receiving agent
func (s *Server) createReceivedInterAgentEvent(routing *messaging.RoutingInfo, originalEvent *messaging.LLMEvent) *messaging.LLMEvent {
	// Extract message from original event
	if len(originalEvent.Parts) > 0 {
		if textPart, ok := originalEvent.Parts[0].(*messaging.TextPart); ok {
			return messaging.NewInterAgentEvent(
				routing,
				textPart.Content,
				messaging.InterAgentReceived,
			)
		}
	}

	return nil
}

// broadcastInterAgentEvent broadcasts an inter-agent event to both sender and receiver agents
func (s *Server) broadcastInterAgentEvent(event *messaging.LLMEvent) {
	if event.InterAgent == nil {
		s.logger.Error("Invalid inter-agent event: missing InterAgent data")
		return
	}

	fromAgentID := event.Routing.FromAgentID
	toAgentID := event.Routing.ToAgentID

	s.logger.Info("Broadcasting inter-agent event",
		zap.Any("fromAgent", fromAgentID),
		zap.Any("toAgent", toAgentID),
	)

	senderCount := s.ssePool.BroadcastToAgent(fromAgentID, event)
	receivedEvent := s.createReceivedInterAgentEvent(event.Routing, event)
	receiverCount := s.ssePool.BroadcastToAgent(toAgentID, receivedEvent)

	totalSent := senderCount + receiverCount
	if totalSent == 0 {
		s.logger.Info("No active SSE connections for inter-agent event",
			zap.Any("fromAgent", fromAgentID),
			zap.Any("toAgent", toAgentID),
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
