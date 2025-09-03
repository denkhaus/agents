package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	s.ssePool.BroadcastToAgent(info.Name, event)
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
	s.ssePool.BroadcastToAgent(info.Name, event)
}

// MultiChatRequest represents a request to send a message in multi-agent chat
type MultiChatRequest struct {
	FromAgent string `json:"fromAgent"` // Agent name or "user"
	ToAgent   string `json:"toAgent"`   // Target agent name
	Message   string `json:"message"`   // Message content
	SessionID string `json:"sessionId"` // Session ID
	UserID    string `json:"userId"`    // User ID
}

// InterAgentEvent represents an inter-agent communication event
type InterAgentEvent struct {
	Type      string    `json:"type"`      // Always "inter_agent"
	FromAgent string    `json:"fromAgent"` // Source agent name
	ToAgent   string    `json:"toAgent"`   // Target agent name
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

	s.chatProcessor.SetMessageInterceptor(func(fromID, toID uuid.UUID, content string) {
		fromName := s.chatProcessor.GetAgentNameByID(fromID)
		toName := s.chatProcessor.GetAgentNameByID(toID)

		if fromName != "" && toName != "" {
			// Create inter-agent event in ADK format
			interAgentEvent := s.createInterAgentEvent(fromName, toName, content)

			// Broadcast to all active SSE connections
			s.broadcastInterAgentEvent(interAgentEvent)
		}
	})
}

// createInterAgentEvent creates an ADK-compatible event for inter-agent communication
func (s *Server) createInterAgentEvent(fromAgent, toAgent, message string) map[string]interface{} {
	eventID := uuid.New().String()
	timestamp := time.Now()

	return map[string]interface{}{
		"id":           eventID,
		"invocationId": eventID,
		"author":       fromAgent,
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
			"fromAgent": fromAgent,
			"toAgent":   toAgent,
			"type":      "communication",
		},
	}
}

// createReceivedInterAgentEvent creates an event for the receiving agent
func (s *Server) createReceivedInterAgentEvent(fromAgent, toAgent string, originalEvent map[string]interface{}) map[string]interface{} {
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
		"author":       fromAgent,
		"timestamp":    timestamp.Unix(),
		"object":       "inter_agent",
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
			"fromAgent": fromAgent,
			"toAgent":   toAgent,
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

	fromAgent, ok := interAgentData["fromAgent"].(string)
	if !ok {
		s.logger.Error("Invalid inter-agent event format: missing fromAgent")
		return
	}

	toAgent, ok := interAgentData["toAgent"].(string)
	if !ok {
		s.logger.Error("Invalid inter-agent event format: missing toAgent")
		return
	}

	s.logger.Info("Broadcasting inter-agent event",
		zap.String("fromAgent", fromAgent),
		zap.String("toAgent", toAgent),
	)

	senderCount := s.ssePool.BroadcastToAgent(fromAgent, event)
	receivedEvent := s.createReceivedInterAgentEvent(fromAgent, toAgent, event)
	receiverCount := s.ssePool.BroadcastToAgent(toAgent, receivedEvent)

	totalSent := senderCount + receiverCount
	if totalSent == 0 {
		s.logger.Info("No active SSE connections for inter-agent event",
			zap.String("fromAgent", fromAgent),
			zap.String("toAgent", toAgent),
		)
	} else {
		s.logger.Info("Sent inter-agent event",
			zap.Int("totalSent", totalSent),
			zap.Int("senderCount", senderCount),
			zap.Int("receiverCount", receiverCount),
		)
	}
}

// registerMultiChatRoutes adds multi-agent chat endpoints to the router
func (s *Server) registerMultiChatRoutes() {
	s.router.HandleFunc("/multi-chat/send", s.handleMultiChatSend).Methods(http.MethodPost)
	s.router.HandleFunc("/multi-chat/start_sse", s.handleMultiChatSSE).Methods(http.MethodGet)

	preflight := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	s.router.HandleFunc("/multi-chat/send", preflight).Methods(http.MethodOptions)
	s.router.HandleFunc("/multi-chat/start_sse", preflight).Methods(http.MethodOptions)
}

// handleMultiChatSend handles sending messages between agents
func (s *Server) handleMultiChatSend(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("handleMultiChatSend called", zap.String("path", r.URL.Path))

	if s.chatProcessor == nil {
		http.Error(w, "Multi-agent chat not configured", http.StatusServiceUnavailable)
		return
	}

	var req MultiChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.ToAgent == "" || req.Message == "" {
		http.Error(w, "toAgent and message are required", http.StatusBadRequest)
		return
	}

	toAgentInfo := s.chatProcessor.GetAgentInfoByAuthor(req.ToAgent)
	if toAgentInfo == nil {
		http.Error(w, fmt.Sprintf("Agent '%s' not found", req.ToAgent), http.StatusNotFound)
		return
	}

	var fromAgentID uuid.UUID
	if req.FromAgent == "" || req.FromAgent == "user" {
		fromAgentID = shared.AgentIDHuman
	} else {
		fromAgentInfo := s.chatProcessor.GetAgentInfoByAuthor(req.FromAgent)
		if fromAgentInfo == nil {
			http.Error(w, fmt.Sprintf("Source agent '%s' not found", req.FromAgent), http.StatusNotFound)
			return
		}
		fromAgentID = fromAgentInfo.ID()
	}

	go func() {
		ctx := context.Background()
		err := s.chatProcessor.SendMessageWithProcessing(ctx, fromAgentID, toAgentInfo.ID(), req.Message)
		if err != nil {
			s.logger.Error("Error sending multi-chat message", zap.Error(err))
		}
	}()

	response := map[string]interface{}{"success": true, "message": "Message sent successfully"}
	s.writeJSON(w, response)
}

// handleMultiChatSSE handles SSE connections for multi-agent chat
func (s *Server) handleMultiChatSSE(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("handleMultiChatSSE called",
		zap.String("path", r.URL.Path),
		zap.String("query", r.URL.RawQuery),
	)

	if s.chatProcessor == nil {
		http.Error(w, "Multi-agent chat not configured", http.StatusServiceUnavailable)
		return
	}

	agents := r.URL.Query().Get("agents")
	sessionID := r.URL.Query().Get("sessionId")
	if agents == "" {
		http.Error(w, "agents parameter is required", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	agentNames := strings.Split(agents, ",")

	var agentList []map[string]interface{}
	for _, agentName := range agentNames {
		agentName = strings.TrimSpace(agentName)
		if agentInfo := s.chatProcessor.GetAgentInfoByAuthor(agentName); agentInfo != nil {
			agentList = append(agentList, map[string]interface{}{
				"id":   agentInfo.ID().String(),
				"name": agentInfo.Name,
				"role": string(agentInfo.Role()),
			})
		}
	}

	initialEvent := map[string]interface{}{"type": "agent_list", "agents": agentList, "timestamp": time.Now().Unix()}
	data, _ := json.Marshal(initialEvent)
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()

	var connections []*SSEConnection
	for _, agentName := range agentNames {
		agentName = strings.TrimSpace(agentName)
		if agentName != "" {
			conn := s.ssePool.RegisterConnection(sessionID, agentName, "user", w, r.Context())
			if conn != nil {
				connections = append(connections, conn)
				s.logger.Info("Registered SSE connection for agent",
					zap.String("agentName", agentName),
					zap.String("sessionID", sessionID),
				)
			}
		}
	}

	defer func() {
		for _, conn := range connections {
			s.ssePool.UnregisterConnection(conn.SessionID, conn.AgentName)
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Multi-chat SSE connection closed", zap.String("sessionID", sessionID))
			return
		case <-ticker.C:
			heartbeat := map[string]interface{}{"type": "heartbeat", "timestamp": time.Now().Unix()}
			data, _ := json.Marshal(heartbeat)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
			s.logger.Debug("Sent heartbeat", zap.String("sessionID", sessionID))
		}
	}
}

// GetMultiChatAgents returns the list of available agents for multi-agent chat
func (s *Server) GetMultiChatAgents() []map[string]interface{} {
	if s.chatProcessor == nil {
		return nil
	}

	var agents []map[string]interface{}
	for _, info := range s.chatProcessor.GetAllAgentInfos() {
		agents = append(agents, map[string]interface{}{
			"id":   info.ID().String(),
			"name": info.Name,
			"role": string(info.Role()),
		})
	}
	return agents
}
