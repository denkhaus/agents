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
	"trpc.group/trpc-go/trpc-agent-go/log"
)

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
func WithChatProcessor(processor multi.ChatProcessor) Option {
	return func(s *Server) {
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
		"actions": map[string]interface{}{
			"stateDelta":           map[string]interface{}{},
			"artifactDelta":        map[string]interface{}{},
			"requestedAuthConfigs": map[string]interface{}{},
		},
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
		"actions": map[string]interface{}{
			"stateDelta":           map[string]interface{}{},
			"artifactDelta":        map[string]interface{}{},
			"requestedAuthConfigs": map[string]interface{}{},
		},
		"interAgent": map[string]interface{}{
			"fromAgent": fromAgent,
			"toAgent":   toAgent,
			"type":      "received",
		},
	}
}

// broadcastInterAgentEvent broadcasts an inter-agent event to both sender and receiver agents
func (s *Server) broadcastInterAgentEvent(event map[string]interface{}) {
	// Extract sender agent name from the inter-agent event
	interAgentData, ok := event["interAgent"].(map[string]interface{})
	if !ok {
		log.Errorf("Invalid inter-agent event format: missing interAgent data")
		return
	}
	
	fromAgent, ok := interAgentData["fromAgent"].(string)
	if !ok {
		log.Errorf("Invalid inter-agent event format: missing fromAgent")
		return
	}
	
	toAgent, ok := interAgentData["toAgent"].(string)
	if !ok {
		log.Errorf("Invalid inter-agent event format: missing toAgent")
		return
	}

	log.Infof("Broadcasting inter-agent event: %s -> %s", fromAgent, toAgent)

	// Broadcast to sender agent sessions
	senderCount := s.ssePool.BroadcastToAgent(fromAgent, event)
	
	// Create and broadcast received message event to receiver agent
	receivedEvent := s.createReceivedInterAgentEvent(fromAgent, toAgent, event)
	receiverCount := s.ssePool.BroadcastToAgent(toAgent, receivedEvent)

	totalSent := senderCount + receiverCount
	if totalSent == 0 {
		log.Infof("No active SSE connections found for agents: %s or %s", fromAgent, toAgent)
	} else {
		log.Infof("Sent inter-agent event to %d connections (sender: %d, receiver: %d)", totalSent, senderCount, receiverCount)
	}
}

// registerMultiChatRoutes adds multi-agent chat endpoints to the router
func (s *Server) registerMultiChatRoutes() {
	// Multi-Agent Chat APIs
	s.router.HandleFunc("/multi-chat/send", s.handleMultiChatSend).Methods(http.MethodPost)
	s.router.HandleFunc("/multi-chat/start_sse", s.handleMultiChatSSE).Methods(http.MethodGet)

	// OPTIONS handlers for CORS
	preflight := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	s.router.HandleFunc("/multi-chat/send", preflight).Methods(http.MethodOptions)
	s.router.HandleFunc("/multi-chat/start_sse", preflight).Methods(http.MethodOptions)
}

// handleMultiChatSend handles sending messages between agents
func (s *Server) handleMultiChatSend(w http.ResponseWriter, r *http.Request) {
	log.Infof("handleMultiChatSend called: path=%s", r.URL.Path)

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

	// Validate request
	if req.ToAgent == "" || req.Message == "" {
		http.Error(w, "toAgent and message are required", http.StatusBadRequest)
		return
	}

	// Get agent info
	toAgentInfo := s.chatProcessor.GetAgentInfoByAuthor(req.ToAgent)
	if toAgentInfo == nil {
		http.Error(w, fmt.Sprintf("Agent '%s' not found", req.ToAgent), http.StatusNotFound)
		return
	}

	// Determine sender ID
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

	// Send message using chat processor
	go func() {
		ctx := context.Background()
		err := s.chatProcessor.SendMessageWithProcessing(ctx, fromAgentID, toAgentInfo.ID(), req.Message)
		if err != nil {
			log.Errorf("Error sending multi-chat message: %v", err)
		}
	}()

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Message sent successfully",
		"from":    req.FromAgent,
		"to":      req.ToAgent,
	}
	s.writeJSON(w, response)
}

// handleMultiChatSSE handles SSE connections for multi-agent chat
func (s *Server) handleMultiChatSSE(w http.ResponseWriter, r *http.Request) {
	log.Infof("handleMultiChatSSE called: path=%s", r.URL.Path)

	if s.chatProcessor == nil {
		http.Error(w, "Multi-agent chat not configured", http.StatusServiceUnavailable)
		return
	}

	// Parse query parameters
	agents := r.URL.Query().Get("agents")
	sessionID := r.URL.Query().Get("sessionId")
	_ = r.URL.Query().Get("userId") // userID for future use

	if agents == "" {
		http.Error(w, "agents parameter is required", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Parse agent list
	agentNames := strings.Split(agents, ",")

	// Send initial agent list
	agentList := make([]map[string]interface{}, 0)
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

	initialEvent := map[string]interface{}{
		"type":      "agent_list",
		"agents":    agentList,
		"timestamp": time.Now().Unix(),
	}

	data, _ := json.Marshal(initialEvent)
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()

	// Register SSE connections for each agent in the connection pool
	var connections []*SSEConnection
	for _, agentName := range agentNames {
		agentName = strings.TrimSpace(agentName)
		if agentName != "" {
			conn := s.ssePool.RegisterConnection(sessionID, agentName, "user", w, r.Context())
			if conn != nil {
				connections = append(connections, conn)
				log.Infof("Registered SSE connection for agent: %s, session: %s", agentName, sessionID)
			}
		}
	}

	// Cleanup connections when done
	defer func() {
		for _, conn := range connections {
			s.ssePool.UnregisterConnection(conn.SessionID, conn.AgentName)
		}
	}()

	// Keep the connection open and send periodic heartbeats
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			log.Infof("Multi-chat SSE connection closed for session %s", sessionID)
			return
		case <-ticker.C:
			// Send heartbeat
			heartbeat := map[string]interface{}{
				"type":      "heartbeat",
				"timestamp": time.Now().Unix(),
			}
			data, _ := json.Marshal(heartbeat)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
			log.Debugf("Sent heartbeat for session %s", sessionID)
		}
	}
}

// GetMultiChatAgents returns the list of available agents for multi-agent chat
func (s *Server) GetMultiChatAgents() []map[string]interface{} {
	if s.chatProcessor == nil {
		return nil
	}

	agents := make([]map[string]interface{}, 0)
	for _, info := range s.chatProcessor.GetAllAgentInfos() {
		agents = append(agents, map[string]interface{}{
			"id":   info.ID().String(),
			"name": info.Name,
			"role": string(info.Role()),
		})
	}
	return agents
}
