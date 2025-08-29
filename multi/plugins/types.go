package plugins

import "context"

// MessageType represents different types of messages for styling
type MessageType int

const (
	MessageTypeNormal MessageType = iota
	MessageTypeReasoningMessage
	MessageTypeToolCall
	MessageTypeIntercept
	MessageTypeError
	MessageTypeSystem
	MessageTypeAgentError
)

// ChatPlugin defines the interface for chat plugins that can be started.
type ChatPlugin interface {
	// Start begins the chat plugin operation with the given context.
	Start(ctx context.Context) error
}
