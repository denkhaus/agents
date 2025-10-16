package shared

import "github.com/google/uuid"

// AgentInfo represents basic agent information (simplified for CLI)
type AgentInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// MessageSender interface (not used in CLI but kept for compatibility)
type MessageSender interface {
	SendMessage(fromAgentID uuid.UUID, toAgentID uuid.UUID, content string) error
}
