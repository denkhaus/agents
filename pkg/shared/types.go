package shared

import "github.com/google/uuid"

type MessageSender interface {
	SendMessage(fromAgentID uuid.UUID, toAgentID uuid.UUID, content string) error
}
