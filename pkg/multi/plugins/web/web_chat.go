package web

import (
	"context"

	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/multi/plugins"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

type webMultiAgentChatImpl struct {
	plugins.Options
}

// NewWebMultiAgentChat creates a new Web-based multi-agent chat plugin.
// It sets up the chat processor with the provided options and configures message handling.
func NewWebMultiAgentChat(opts ...plugins.MultiAgentChatOption) plugins.ChatPlugin {
	chat := &webMultiAgentChatImpl{
		Options: plugins.Options{},
	}

	for _, opt := range opts {
		opt(&chat.Options)
	}

	processorOptions := []multi.ChatProcessorOption{
		multi.WithOnProgress(chat.handleOnProgress),
		multi.WithOnMessage(chat.handleOnMessage),
		multi.WithOnReasoningMessage(chat.handleOnReasoningMessage),
		multi.WithOnError(chat.handleOnError),
		multi.WithOnToolCall(chat.handleOnToolCall),
	}

	processorOptions = append(processorOptions, chat.ProcessorOptions...)
	chat.Processor = multi.NewChatProcessor(processorOptions...)
	chat.setupMessageListener()

	return chat
}

// setupMessageListener configures the message interceptor to display inter-agent communication.
func (p *webMultiAgentChatImpl) setupMessageListener() {
	// Add a message interceptor to the broker
	p.Processor.SetMessageInterceptor(func(fromID, toID uuid.UUID, content string) {
		// fromName := p.Processor.GetAgentNameByID(fromID)
		// toName := p.Processor.GetAgentNameByID(toID)

		// if fromName != "" && toName != "" {
		// 	// Format: "FromName (FromID) -> ToName (ToID)"
		// 	header := fmt.Sprintf("%s (%s) -> %s (%s)",
		// 		fromName, fromID, toName, toID,
		// 	)
		// }
	})
}

// handleOnProgress handles progress updates by printing them to stdout.
func (p *webMultiAgentChatImpl) handleOnProgress(messageType multi.SystemMessageType, format string, a ...any) {

}

// handleOnMessage handles agent messages by displaying them with a formatted border.
func (p *webMultiAgentChatImpl) handleOnMessage(info *shared.AgentInfo, content string) {

}

// handleOnError handles agent errors by displaying them with a formatted border.
func (p *webMultiAgentChatImpl) handleOnError(info *shared.AgentInfo, err error) {

}

// handleOnToolCall handles tool calls made by agents by displaying them with a formatted border.
func (p *webMultiAgentChatImpl) handleOnToolCall(info *shared.AgentInfo, functionDef model.FunctionDefinitionParam) {

}

// handleOnReasoningMessage handles reasoning messages from agents.
func (p *webMultiAgentChatImpl) handleOnReasoningMessage(info *shared.AgentInfo, reasoning string) {
}

// Start runs the interactive chat loop, handling user input and agent communication.
// It supports commands like /exit, /list, /agent-name to select agents, and direct messaging.
func (p *webMultiAgentChatImpl) Start(ctx context.Context) error {

	return nil
}
