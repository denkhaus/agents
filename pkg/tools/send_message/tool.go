//cue:generate cue get go github.com/denkhaus/agents/pkg/tools/sendmessage

package sendmessage

import (
	"context"
	"errors"
	"fmt"

	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/tools"
	"github.com/google/uuid"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// -----------------------------------------------------------------------------
// Constants & helpers ----------------------------------------------------------
// -----------------------------------------------------------------------------
const (
	ToolName = "send_message"
)

// ToolConfig holds configuration for the messaging tool
type ToolConfig struct {
	Sender                   shared.MessageSender `json:"broker" mapstructure:"broker"`
	AllowedToCommunicateWith []uuid.UUID          `json:"allowed_to_communicate_with" mapstructure:"allowed_to_communicate_with"`
	AgentID                  uuid.UUID            `json:"agent_id" mapstructure:"agent_id"`
}

// sendMessageArgs holds the input for the messaging tool.
type sendMessageArgs struct {
	To      string `json:"to" description:"The UUID of the recipient agent"`
	Content string `json:"content" description:"The message content"`
}

// sendMessageResult holds the output for the messaging tool.
type sendMessageResult struct {
	Status  string `json:"status"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type SendMessageFunc func(context.Context, sendMessageArgs) (sendMessageResult, error)

// sendMessage performs the message sending operation.
// It sends a message from the current agent to another agent by ID.
func sendMessage(sender shared.MessageSender, agentID uuid.UUID, allowedToCommunicateWith []uuid.UUID) SendMessageFunc {

	isAllowed := func(agentID uuid.UUID) bool {
		for _, id := range allowedToCommunicateWith {
			if agentID == id {
				return true
			}
		}

		return false
	}

	return func(ctx context.Context, args sendMessageArgs) (sendMessageResult, error) {
		// Parse the recipient UUID
		to, err := uuid.Parse(args.To)
		if err != nil {
			return sendMessageResult{}, fmt.Errorf("invalid 'to' parameter: %w", err)
		}

		if to == agentID {
			return sendMessageResult{}, errors.New("you cannot send a message to yourself")
		}

		if !isAllowed(to) {
			return sendMessageResult{}, fmt.Errorf("you are not allowed to send a message to the agent with id %q", to)
		}

		// Send the message through the broker
		err = sender.SendMessage(agentID, to, args.Content)
		if err != nil {
			return sendMessageResult{}, fmt.Errorf("failed to send message: %w", err)
		}

		return sendMessageResult{
			Status:  "sent",
			To:      to.String(),
			Content: args.Content,
		}, nil
	}
}

func NewWithDI(injector *do.Injector) (tools.ToolFactoryFunc, error) {
	return func(config tools.ConfigPayload, _ []*shared.AgentInfo) (tool.Tool, error) {
		var toolConfig ToolConfig
		if err := config.Bind(&toolConfig); err != nil {
			return nil, fmt.Errorf("failed to bind messaging tool config: %w", err)
		}

		return New(
			toolConfig.Sender,
			toolConfig.AgentID,
			toolConfig.AllowedToCommunicateWith,
		)

	}, nil
}

func New(sender shared.MessageSender, agentID uuid.UUID, allowedToCommunicateWith []uuid.UUID) (tool.Tool, error) {
	// Create messaging tool for inter-agent communication.
	messagingTool := function.NewFunctionTool(
		sendMessage(sender, agentID, allowedToCommunicateWith),
		function.WithName(ToolName),
		function.WithDescription(
			"Send a message to another agent by ID",
		),
	)

	return messagingTool, nil
}
