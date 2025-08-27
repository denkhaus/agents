package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// messagingToolImpl is a tool that allows agents to send messages
type messagingToolImpl struct {
	broker  MessageBroker
	agentID uuid.UUID
}

// NewMessagingTool creates a new messaging tool
func NewMessagingTool(broker MessageBroker, agentID uuid.UUID) tool.Tool {
	return &messagingToolImpl{
		broker:  broker,
		agentID: agentID,
	}
}

// Declaration returns the tool declaration
func (mt *messagingToolImpl) Declaration() *tool.Declaration {
	return &tool.Declaration{
		Name:        "send_message",
		Description: "Send a message to another agent by ID",
		InputSchema: &tool.Schema{
			Type: "object",
			Properties: map[string]*tool.Schema{
				"to": {
					Type:        "string",
					Description: "The UUID of the recipient agent",
				},
				"content": {
					Type:        "string",
					Description: "The message content",
				},
			},
			Required: []string{"to", "content"},
		},
	}
}

// Call executes the tool
func (mt *messagingToolImpl) Call(ctx context.Context, jsonArgs []byte) (any, error) {
	// Parse the arguments
	var args map[string]interface{}
	if err := json.Unmarshal(jsonArgs, &args); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	toStr, ok := args["to"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'to' parameter")
	}

	to, err := uuid.Parse(toStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'to' parameter: %w", err)
	}

	content, ok := args["content"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'content' parameter")
	}

	err = mt.broker.SendMessage(mt.agentID, to, content)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return map[string]interface{}{
		"status":  "sent",
		"to":      to.String(),
		"content": content,
	}, nil
}
