package chat

import (
	"context"

	"github.com/google/uuid"
)

type Chat interface {
	ProcessMessage(ctx context.Context, userMessage string) error
}

type ChatProvider interface {
	GetChat(ctx context.Context, agentID uuid.UUID, opts ...Option) (Chat, error)
}
