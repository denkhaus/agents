package web

import (
	"fmt"
	"time"

	"github.com/denkhaus/agents/pkg/multi"
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

	eventID := uuid.New().String()
	timestamp := time.Now()

	// Create an ADK-compatible event
	event := map[string]interface{}{
		"id":           eventID,
		"invocationId": eventID,
		"author":       info.Name,
		"timestamp":    timestamp.Unix(),
		"object":       "message",
		"done":         true,
		"partial":      false,
		"content": map[string]interface{}{
			"role": "assistant",
			"parts": []map[string]interface{}{
				{
					"text": content,
				},
			},
		},
		"actions": map[string]interface{}{},
	}

	s.logger.Info("Broadcasting agent message",
		zap.String("agent", info.Name),
		zap.Any("event", event),
	)

	s.ssePool.BroadcastToAgent(info.ID(), event)
}

// BroadcastToolCall sends a tool call from an agent to all connected clients for that agent.
func (s *Server) BroadcastToolCall(info *shared.AgentInfo, functionDef model.FunctionDefinitionParam) {
	if info == nil {
		s.logger.Error("BroadcastToolCall called with nil agent info")
		return
	}

	eventID := uuid.New().String()
	timestamp := time.Now()

	// Create an ADK-compatible event for tool call
	event := map[string]interface{}{
		"id":           eventID,
		"invocationId": eventID,
		"author":       info.Name,
		"timestamp":    timestamp.Unix(),
		"object":       "tool_code",
		"done":         true,
		"partial":      false,
		"content": map[string]interface{}{
			"role": "tool",
			"parts": []map[string]interface{}{
				{
					"functionCall": map[string]interface{}{
						"name":      functionDef.Name,
						"arguments": string(functionDef.Arguments),
					},
				},
			},
		},
		"actions": map[string]interface{}{},
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

// createInterAgentEvent creates an ADK-compatible event for inter-agent communication
func (s *Server) createInterAgentEvent(fromAgent, toAgent uuid.UUID, message string) map[string]interface{} {
	eventID := uuid.New().String()
	timestamp := time.Now()

	return map[string]interface{}{
		"id":           eventID,
		"invocationId": eventID,
		"author":       fromAgent.String(),
		"timestamp":    timestamp.Unix(),
		"object":       "inter_agent",
		"done":         true,
		"partial":      false,
		"content": map[string]interface{}{
			"role": "assistant",
			"parts": []map[string]interface{}{
				{
					"text": message,
				},
			},
		},
		"actions": map[string]interface{}{},
		"interAgent": map[string]interface{}{
			"fromAgent": fromAgent.String(),
			"toAgent":   toAgent.String(),
			"type":      "communication",
		},
	}
}

// createReceivedInterAgentEvent creates an event for the receiving agent
func (s *Server) createReceivedInterAgentEvent(fromAgent, toAgent uuid.UUID, originalEvent map[string]interface{}) map[string]interface{} {
	eventID := uuid.New().String()
	timestamp := time.Now()

	// Extract message from original event
	var message string
	if content, ok := originalEvent["content"].(map[string]interface{}); ok {
		if parts, ok := content["parts"].([]map[string]interface{}); ok && len(parts) > 0 {
			if text, ok := parts[0]["text"].(string); ok {
				message = text
			}
		}
	}

	return map[string]interface{}{
		"id":           eventID,
		"invocationId": eventID,
		"author":       fromAgent.String(),
		"timestamp":    timestamp.Unix(),
		"type":         "inter_agent",
		"done":         true,
		"partial":      false,
		"content": map[string]interface{}{
			"role": "user",
			"parts": []map[string]interface{}{
				{
					"text": fmt.Sprintf("Received from %s: %s", fromAgent, message),
				},
			},
		},
		"actions": map[string]interface{}{},
		"interAgent": map[string]interface{}{
			"fromAgent": fromAgent.String(),
			"toAgent":   toAgent.String(),
			"type":      "received",
		},
	}
}

// broadcastInterAgentEvent broadcasts an inter-agent event to both sender and receiver agents
func (s *Server) broadcastInterAgentEvent(event map[string]interface{}) {
	interAgentData, ok := event["interAgent"].(map[string]interface{})
	if !ok {
		s.logger.Error("Invalid inter-agent event format: missing interAgent data")
		return
	}

	fromAgentStr, ok := interAgentData["fromAgent"].(string)
	if !ok {
		s.logger.Error("Invalid inter-agent event format: missing fromAgent")
		return
	}
	fromAgent, err := uuid.Parse(fromAgentStr)
	if err != nil {
		s.logger.Error("Invalid fromAgent UUID", zap.String("uuid", fromAgentStr))
		return
	}

	toAgentStr, ok := interAgentData["toAgent"].(string)
	if !ok {
		s.logger.Error("Invalid inter-agent event format: missing toAgent")
		return
	}
	toAgent, err := uuid.Parse(toAgentStr)
	if err != nil {
		s.logger.Error("Invalid toAgent UUID", zap.String("uuid", toAgentStr))
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
