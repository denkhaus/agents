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
	"trpc.group/trpc-go/trpc-agent-go/runner"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
)

// ChatProcessor defines the interface for managing multi-agent chat interactions.
// It provides methods for sending messages between agents, retrieving agent information,
// and setting up message interception for monitoring communication.
type ChatProcessor interface {
	// SetMessageInterceptor sets a function to intercept and monitor messages between agents.
	// The interceptor receives the sender ID, receiver ID, and message content.
	SetMessageInterceptor(interceptor messaging.Interceptor)
	
	// SendMessage sends a message from one agent to another and returns a channel of events.
	// The caller is responsible for processing the events from the returned channel.
	SendMessage(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) (<-chan *event.Event, error)
	
	// SendMessageWithProcessing sends a message and automatically processes all resulting events.
	// This is a convenience method that handles event processing internally.
	SendMessageWithProcessing(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) error
	
	// GetAgentInfoByAuthor retrieves agent information by author name or UUID string.
	// Returns nil if no agent is found with the given identifier.
	GetAgentInfoByAuthor(author string) *shared.AgentInfo
	
	// GetAllAgentInfos returns information for all registered agents in the chat processor.
	GetAllAgentInfos() []shared.AgentInfo
	
	// GetAgentNameByID returns the name of an agent given its UUID.
	// Returns empty string if no agent is found with the given ID.
	GetAgentNameByID(agentID uuid.UUID) string
}

// chatProcessorImpl implements the ChatProcessor interface and manages
// the lifecycle and communication between multiple agents.
type chatProcessorImpl struct {
	Options
	agents map[uuid.UUID]*AgentRunner
	broker messaging.MessageBroker
}

// NewChatProcessor creates a new ChatProcessor instance with the given options.
// It initializes the message broker, sets up default configuration, and registers all agents.
func NewChatProcessor(opts ...ChatProcessorOption) ChatProcessor {
	processor := &chatProcessorImpl{
		agents: make(map[uuid.UUID]*AgentRunner),
		broker: messaging.NewMessageBroker(),
		Options: Options{
			applicationName: "chat-app-default",
		},
	}

	for _, opt := range opts {
		opt(&processor.Options)
	}

	processor.initAgents()
	return processor
}

// initAgents initializes all agents in the processor by creating AgentRunner instances
// and setting up message processing for each agent.
func (p *chatProcessorImpl) initAgents() {
	for _, agent := range p.Options.agents {
		if _, exists := p.agents[agent.ID()]; exists {
			logger.Log.Warn("agent already registered in chat processor",
				zap.String("app_name", p.applicationName),
				zap.Any("agent_id", agent.ID()),
			)
			continue
		}

		wrapper := messaging.NewMessagingWrapper(agent, p.broker)

		ar := &AgentRunner{
			wrapper: wrapper,
			runner: runner.NewRunner(
				p.applicationName,
				wrapper,
				runner.WithSessionService(
					inmemory.NewSessionService(),
				),
			),
		}

		p.agents[agent.ID()] = ar
		p.startMessageProcessing(ar)
	}
}

// SetMessageInterceptor sets a message interceptor on the underlying message broker.
// The interceptor function will be called for every message sent between agents.
func (p *chatProcessorImpl) SetMessageInterceptor(interceptor messaging.Interceptor) {
	p.broker.SetMessageInterceptor(interceptor)
}

// GetAllAgentInfos returns a slice containing information about all registered agents.
func (p *chatProcessorImpl) GetAllAgentInfos() []shared.AgentInfo {
	var infos []shared.AgentInfo
	for _, agent := range p.agents {
		infos = append(infos, *agent.Info())
	}

	return infos
}

// GetAgentInfoByID retrieves agent information by UUID.
// Returns nil if no agent is found with the given ID.
func (p *chatProcessorImpl) GetAgentInfoByID(agentID uuid.UUID) *shared.AgentInfo {
	if v, ok := p.agents[agentID]; ok {
		return v.Info()
	}

	return nil
}

// GetAgentInfoByAuthor retrieves agent information by author identifier.
// The author can be either a UUID string or an agent name.
// Returns nil if no agent is found with the given identifier.
func (p *chatProcessorImpl) GetAgentInfoByAuthor(author string) *shared.AgentInfo {
	// Try to parse as UUID first
	if authorID, err := uuid.Parse(author); err == nil {
		if info := p.GetAgentInfoByID(authorID); info != nil {
			return info
		}
	}

	// If not UUID or not found, check if it's already a name
	for _, agent := range p.agents {
		if agent.Name() == author {
			return agent.Info()
		}
	}

	return nil
}

// GetAgentNameByID returns the agent name for a given AgentID
func (p *chatProcessorImpl) GetAgentNameByID(agentID uuid.UUID) string {
	// Check all agents
	for _, agent := range p.agents {
		if agent.ID() == agentID {
			return agent.Name()
		}
	}

	return ""
}

// startMessageProcessing starts a goroutine to process incoming messages for the given agent.
// It listens on the agent's message channel and forwards messages to the agent's runner.
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

// SendMessage sends a message from one agent to another and returns a channel of events.
// The caller is responsible for processing the events from the returned channel.
func (p *chatProcessorImpl) SendMessage(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) (<-chan *event.Event, error) {
	agent, exists := p.agents[toAgentID]
	if !exists {
		return nil, fmt.Errorf("agent %q not found", toAgentID)
	}

	userMessage := model.NewUserMessage(message)
	return agent.Run(ctx, fromAgentID, userMessage)
}

// SendMessageWithProcessing sends a message to an agent and automatically processes all resulting events.
// This method handles event processing internally and provides progress updates through callbacks.
func (p *chatProcessorImpl) SendMessageWithProcessing(ctx context.Context, fromAgentID, toAgentID uuid.UUID, message string) error {
	agent, exists := p.agents[toAgentID]
	if !exists {
		return fmt.Errorf("agent %q not found", toAgentID)
	}

	userMessage := model.NewUserMessage(message)
	p.onProgress(SystemMessageSending, "sending message to %s...", agent)

	events, err := agent.Run(ctx, fromAgentID, userMessage)
	if err != nil {
		return fmt.Errorf("failed to send message from %s to %s: %w", fromAgentID, toAgentID, err)
	}

	p.onProgress(SystemMessageDelivered, "message delivered to %s - Processing...", agent)

	// Process events
	for event := range events {
		p.processEvent(event)
	}

	p.onProgress(SystemMessageProcessed, "%s finished processing", agent)
	return nil
}

// processEvent processes a single event from an agent's response.
// It handles errors, assistant messages, and tool calls by invoking the appropriate callbacks.
func (p *chatProcessorImpl) processEvent(event *event.Event) {
	if event.Error != nil {
		info := p.GetAgentInfoByAuthor(event.Author)
		p.onError(info, errors.New(event.Error.Message))
	}

	if event.Response != nil && len(event.Response.Choices) > 0 {
		choice := event.Response.Choices[0]

		// Show reasoning content first if present (future-proof detection)
		if choice.Message.ReasoningContent != "" {
			info := p.GetAgentInfoByAuthor(event.Author)
			p.onReasoningMessage(info, choice.Message.ReasoningContent)
		}

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
