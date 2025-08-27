package multi

import (
	"context"
	"errors"
	"fmt"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/messaging"
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

type ChatProcessor interface {
	SetMessageInterceptor(interceptor messaging.Interceptor)
	SendMessage(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) (<-chan *event.Event, error)
	SendMessageWithProcessing(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) error
	GetAgentInfoByAuthor(author string) *shared.AgentInfo
	GetAllAgentInfos() []shared.AgentInfo
	GetAgentNameByID(agentID uuid.UUID) string
}

type chatProcessorImpl struct {
	Options
	agents map[uuid.UUID]*AgentRunner
	broker messaging.MessageBroker
}

func NewChatProcessor(opts ...ChatProcessorOption) ChatProcessor {
	processor := &chatProcessorImpl{
		agents: make(map[uuid.UUID]*AgentRunner),
		broker: messaging.NewMessageBroker(),
	}

	for _, opt := range opts {
		opt(&processor.Options)
	}

	return processor
}

func (p *chatProcessorImpl) SetMessageInterceptor(interceptor messaging.Interceptor) {
	p.broker.SetMessageInterceptor(interceptor)
}

func (p *chatProcessorImpl) GetAllAgentInfos() []shared.AgentInfo {
	var infos []shared.AgentInfo
	for _, agent := range p.agents {
		infos = append(infos, *agent.Info())
	}

	return infos
}

func (p *chatProcessorImpl) GetAgentInfoByID(agentID uuid.UUID) *shared.AgentInfo {
	if v, ok := p.agents[agentID]; ok {
		return v.Info()
	}

	return nil
}

func (p *chatProcessorImpl) GetAgentInfoByAuthor(author string) *shared.AgentInfo {
	// Try to parse as UUID first
	if authorID, err := uuid.Parse(author); err == nil {
		if info := p.GetAgentInfoByID(authorID); info != nil {
			return info
		}
		if p.humanAgent.ID() == authorID {
			return shared.TheAgentToInfo(p.humanAgent)
		}
	}

	// If not UUID or not found, check if it's already a name
	for _, agent := range p.agents {
		if agent.Name() == author {
			return agent.Info()
		}
	}

	if p.humanAgent.Info().Name == author {
		return shared.TheAgentToInfo(p.humanAgent)
	}

	return nil
}

// GetAgentNameByID returns the agent name for a given AgentID
func (p *chatProcessorImpl) GetAgentNameByID(agentID uuid.UUID) string {
	// Check human
	if p.humanAgent.ID() == agentID {
		return p.humanAgent.Info().Name
	}

	// Check AI agents
	for _, agent := range p.agents {
		if agent.ID() == agentID {
			return agent.Name()
		}
	}

	return ""
}

func (p *chatProcessorImpl) startMessageProcessing(agent *AgentRunner) {
	go func() {
		// Get the message channel for this agent
		msgChan, err := p.broker.GetMessageChannel(agent.ID())
		if err != nil {
			logger.Log.Error("failed to get message channel for agent", zap.String("agent", agent.Name()), zap.Error(err))
			return
		}

		// Process incoming messages
		for msg := range msgChan {
			// Create a context for message processing
			ctx := context.Background()

			// Format the message content
			messageContent := fmt.Sprintf("Message from %s: %s", p.GetAgentNameByID(msg.From), msg.Content)

			// Send to the agent's runner
			events, err := agent.Run(ctx, msg.From, model.NewUserMessage(messageContent))
			if err != nil {
				logger.Log.Error("failed to process message for agent", zap.String("agent", agent.Name()), zap.Error(err))
				continue
			}

			// Process events from the agent's response
			go func() {
				for event := range events {
					p.processEvent(event)
				}
			}()
		}
	}()
}

func (p *chatProcessorImpl) SendMessage(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) (<-chan *event.Event, error) {
	agent, exists := p.agents[toAgentID]
	if !exists {
		return nil, fmt.Errorf("agent %q not found", toAgentID)
	}

	userMessage := model.NewUserMessage(message)
	return agent.Run(ctx, fromAgentID, userMessage)
}

func (p *chatProcessorImpl) SendMessageWithProcessing(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) error {
	agent, exists := p.agents[toAgentID]
	if !exists {
		return fmt.Errorf("agent %q not found", toAgentID)
	}

	userMessage := model.NewUserMessage(message)
	p.onProgress("system >> sending message to %s...\n", agent)

	events, err := agent.Run(ctx, fromAgentID, userMessage)
	if err != nil {
		return fmt.Errorf("failed to send message from %s to %s: %w", fromAgentID, toAgentID, err)
	}

	p.onProgress("system >> message delivered to %s - Processing...\n", agent)

	// Process events
	for event := range events {
		p.processEvent(event)
	}

	p.onProgress("system >> %s finished processing\n", agent)
	return nil
}

func (p *chatProcessorImpl) processEvent(event *event.Event) {
	if event.Error != nil {
		info := p.GetAgentInfoByAuthor(event.Author)
		p.onError(info, errors.New(event.Error.Message))
	}

	if event.Response != nil && len(event.Response.Choices) > 0 {
		choice := event.Response.Choices[0]

		// Show assistant messages
		if choice.Message.Role == model.RoleAssistant && choice.Message.Content != "" {
			info := p.GetAgentInfoByAuthor(event.Author)
			p.onMessage(info, choice.Message.Content)
		}

		// Show tool calls (but suppress the generic "sending message" for cleaner output)
		if len(choice.Message.ToolCalls) > 0 {
			for _, toolCall := range choice.Message.ToolCalls {
				info := p.GetAgentInfoByAuthor(event.Author)
				p.onToolCall(info, toolCall.Function)
			}
		}
	}
}
