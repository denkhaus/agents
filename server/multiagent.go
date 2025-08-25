// Package server provides multi-agent routing extracted from examples.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/server"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

// MultiAgentManager manages multiple agents with dynamic routing.
type MultiAgentManager interface {
	// RegisterAgent registers an agent with the manager
	RegisterAgent(name string, card server.AgentCard, processor MessageProcessor) error
	
	// UnregisterAgent removes an agent from the manager
	UnregisterAgent(name string) error
	
	// GetAgent retrieves an agent by name
	GetAgent(name string) (server.AgentCard, MessageProcessor, error)
	
	// ListAgents returns all registered agents
	ListAgents() map[string]server.AgentCard
	
	// CreateCardHandler creates a dynamic agent card handler
	CreateCardHandler(host string) http.Handler
	
	// CreateMessageProcessor creates a routing message processor
	CreateMessageProcessor() MessageProcessor
}

// multiAgentManager implements MultiAgentManager.
// Extracted from examples/multi_endpoint/server/main.go
type multiAgentManager struct {
	agents map[string]agentInfo
}

// agentInfo contains agent information.
type agentInfo struct {
	card      server.AgentCard
	processor MessageProcessor
}

// NewMultiAgentManager creates a new multi-agent manager.
func NewMultiAgentManager() MultiAgentManager {
	return &multiAgentManager{
		agents: make(map[string]agentInfo),
	}
}

// RegisterAgent registers an agent with the manager.
func (mam *multiAgentManager) RegisterAgent(name string, card server.AgentCard, processor MessageProcessor) error {
	if name == "" {
		return fmt.Errorf("agent name cannot be empty")
	}
	
	if processor == nil {
		return fmt.Errorf("processor cannot be nil")
	}
	
	mam.agents[name] = agentInfo{
		card:      card,
		processor: processor,
	}
	
	return nil
}

// UnregisterAgent removes an agent from the manager.
func (mam *multiAgentManager) UnregisterAgent(name string) error {
	if _, exists := mam.agents[name]; !exists {
		return fmt.Errorf("agent %s not found", name)
	}
	
	delete(mam.agents, name)
	return nil
}

// GetAgent retrieves an agent by name.
func (mam *multiAgentManager) GetAgent(name string) (server.AgentCard, MessageProcessor, error) {
	info, exists := mam.agents[name]
	if !exists {
		return server.AgentCard{}, nil, fmt.Errorf("agent %s not found", name)
	}
	
	return info.card, info.processor, nil
}

// ListAgents returns all registered agents.
func (mam *multiAgentManager) ListAgents() map[string]server.AgentCard {
	result := make(map[string]server.AgentCard)
	for name, info := range mam.agents {
		result[name] = info.card
	}
	return result
}

// CreateCardHandler creates a dynamic agent card handler.
func (mam *multiAgentManager) CreateCardHandler(host string) http.Handler {
	return &multiAgentCardHandler{
		manager: mam,
		host:    host,
	}
}

// CreateMessageProcessor creates a routing message processor.
func (mam *multiAgentManager) CreateMessageProcessor() MessageProcessor {
	return &multiAgentProcessor{
		manager: mam,
	}
}

// multiAgentCardHandler provides dynamic agent cards based on URL parameter.
type multiAgentCardHandler struct {
	manager MultiAgentManager
	host    string
}

// ServeHTTP handles agent card requests.
func (h *multiAgentCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	agentName := chi.URLParam(r, "agentName")
	if agentName == "" {
		http.Error(w, "Agent name is required", http.StatusBadRequest)
		return
	}
	
	// Check if agent exists
	card, _, err := h.manager.GetAgent(agentName)
	if err != nil {
		// Return default card for unknown agents
		card = h.createDefaultAgentCard(agentName)
	}
	
	// Update URL to include host
	if h.host != "" {
		card.URL = fmt.Sprintf("http://%s/api/v1/agent/%s/", h.host, agentName)
	}
	
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(card); err != nil {
		http.Error(w, "Failed to encode agent card", http.StatusInternalServerError)
	}
}

// createDefaultAgentCard creates a default agent card for unknown agents.
func (h *multiAgentCardHandler) createDefaultAgentCard(agentName string) server.AgentCard {
	return server.AgentCard{
		Name:        agentName,
		Description: fmt.Sprintf("Dynamic agent: %s", agentName),
		URL:         fmt.Sprintf("http://%s/api/v1/agent/%s/", h.host, agentName),
		Skills: []server.AgentSkill{
			{
				Name:        "default",
				Description: stringPtr(fmt.Sprintf("Default skill for %s", agentName)),
				InputModes:  []string{"text"},
				OutputModes: []string{"text"},
				Tags:        []string{"dynamic"},
				Examples:    []string{"Hello"},
			},
		},
		Capabilities: server.AgentCapabilities{
			Streaming: boolPtr(false),
		},
		DefaultInputModes:  []string{"text"},
		DefaultOutputModes: []string{"text"},
	}
}

// multiAgentProcessor routes messages based on agent name from context.
type multiAgentProcessor struct {
	manager MultiAgentManager
}

// ProcessMessage routes messages to the appropriate agent processor.
func (p *multiAgentProcessor) ProcessMessage(
	ctx context.Context,
	message protocol.Message,
	options taskmanager.ProcessOptions,
	handle taskmanager.TaskHandler,
) (*taskmanager.MessageProcessingResult, error) {
	// Extract agent name from context
	agentName := getAgentNameFromContext(ctx)
	if agentName == "" {
		return createErrorResult("agent name not found in context"), nil
	}
	
	// Get agent processor
	_, processor, err := p.manager.GetAgent(agentName)
	if err != nil {
		// Use default processing for unknown agents
		return p.processWithDefaultAgent(ctx, message, agentName)
	}
	
	// Delegate to agent processor
	return processor.ProcessMessage(ctx, message, options, handle)
}

// processWithDefaultAgent provides default processing for unknown agents.
func (p *multiAgentProcessor) processWithDefaultAgent(
	ctx context.Context,
	message protocol.Message,
	agentName string,
) (*taskmanager.MessageProcessingResult, error) {
	text := extractTextFromMessage(message)
	response := fmt.Sprintf("Hello from %s agent! You said: %s", agentName, text)
	return createTextResult(response), nil
}

// ContextKey for agent name in request context.
type ContextKey string

const CtxAgentNameKey ContextKey = "agentName"

// AgentNameMiddleware extracts agent name from URL and adds to context.
type AgentNameMiddleware struct{}

// NewAgentNameMiddleware creates a new agent name middleware.
func NewAgentNameMiddleware() *AgentNameMiddleware {
	return &AgentNameMiddleware{}
}

// Wrap wraps an HTTP handler with agent name extraction.
func (m *AgentNameMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		agentName := chi.URLParam(r, "agentName")
		if agentName != "" {
			ctx = context.WithValue(ctx, CtxAgentNameKey, agentName)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getAgentNameFromContext extracts agent name from context.
func getAgentNameFromContext(ctx context.Context) string {
	if agentName, ok := ctx.Value(CtxAgentNameKey).(string); ok {
		return agentName
	}
	return ""
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}