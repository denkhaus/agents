package multi

import (
	"context"
	"errors"
	"fmt"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/messaging"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	"trpc.group/trpc-go/trpc-agent-go/session"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
)

// chatProcessorImpl implements the ChatProcessor interface and manages
// the lifecycle and communication between multiple agents.
type chatProcessorImpl struct {
	Options
	agents map[uuid.UUID]*AgentRunner
	broker messaging.MessageBroker
}

// NewChatProcessor creates a new ChatProcessor instance with the given options.
// It initializes the message broker, sets up default configuration, and registers all agents.
func NewChatProcessor(sessionID uuid.UUID, opts ...ChatProcessorOption) (ChatProcessor, error) {
	processor := &chatProcessorImpl{
		agents: make(map[uuid.UUID]*AgentRunner),
		broker: messaging.NewMessageBroker(sessionID),
		Options: Options{
			sessionService:  inmemory.NewSessionService(),
			applicationName: "chat-app-default",
		},
	}

	for _, opt := range opts {
		opt(&processor.Options)
	}

	// Ensure all callbacks have default implementations to prevent nil panics.
	if processor.onProgress == nil {
		processor.onProgress = func(info *messaging.RoutingInfo, messageType SystemMessageType, format string, a ...any) {
			logger.Log.Warn("onProgress callback not initialized",
				zap.String("app_name", processor.applicationName),
				zap.Any("from_agent_id", info.FromAgentID),
				zap.Any("to_agent_id", info.ToAgentID),
				zap.Any("session_id", info.SessionID),
			)
		}
	}
	if processor.onMessage == nil {
		processor.onMessage = func(info *messaging.RoutingInfo, content string) {
			logger.Log.Warn("onMessage callback not initialized",
				zap.String("app_name", processor.applicationName),
				zap.Any("from_agent_id", info.FromAgentID),
				zap.Any("to_agent_id", info.ToAgentID),
				zap.Any("session_id", info.SessionID),
			)
		}
	}
	if processor.onReasoningMessage == nil {
		processor.onReasoningMessage = func(info *messaging.RoutingInfo, content string) {
			logger.Log.Warn("onReasoningMessage callback not initialized",
				zap.String("app_name", processor.applicationName),
				zap.Any("from_agent_id", info.FromAgentID),
				zap.Any("to_agent_id", info.ToAgentID),
				zap.Any("session_id", info.SessionID),
			)
		}
	}
	if processor.onToolCall == nil {
		processor.onToolCall = func(info *messaging.RoutingInfo, functionDef model.FunctionDefinitionParam) {
			logger.Log.Warn("onToolCall callback not initialized",
				zap.String("app_name", processor.applicationName),
				zap.Any("from_agent_id", info.FromAgentID),
				zap.Any("to_agent_id", info.ToAgentID),
				zap.Any("session_id", info.SessionID),
			)
		}
	}
	if processor.onError == nil {
		processor.onError = func(info *messaging.RoutingInfo, err error) {
			logger.Log.Warn("onError callback not initialized",
				zap.String("app_name", processor.applicationName),
				zap.Any("from_agent_id", info.FromAgentID),
				zap.Any("to_agent_id", info.ToAgentID),
				zap.Any("session_id", info.SessionID),
				zap.Error(err),
			)
		}
	}

	err := processor.initAgents()
	if err != nil {
		return nil, fmt.Errorf("failed to init agents:%w", err)
	}

	return processor, nil
}

func (p *chatProcessorImpl) GetApplicationName() string {
	return p.applicationName
}

func (p *chatProcessorImpl) getAgentInfo() []*shared.AgentInfo {
	result := make([]*shared.AgentInfo, len(p.availableAgents))
	for idx, agent := range p.availableAgents {
		result[idx] = agent.GetInfo()
	}

	return result
}

// initAgents initializes all agents in the processor by creating AgentRunner instances
// and setting up message processing for each agent.
func (p *chatProcessorImpl) initAgents() error {
	agentInfos := p.getAgentInfo()
	for _, agent := range p.availableAgents {
		agentID := agent.GetID()
		if _, exists := p.agents[agentID]; exists {
			logger.Log.Warn("agent already registered in chat processor",
				zap.String("app_name", p.applicationName),
				zap.Any("agent_id", agentID),
			)
			continue
		}

		wrapper, err := messaging.NewAgentWrapper(agent, p.broker, agentInfos...)
		if err != nil {
			return fmt.Errorf("failed to create the wrapper for the %q agent: %w",
				agent.GetInfo().Name, err,
			)
		}

		ar := &AgentRunner{
			wrapper: wrapper,
			runner: runner.NewRunner(
				p.applicationName,
				wrapper,
				runner.WithSessionService(
					p.sessionService,
				),
			),
		}

		p.agents[agentID] = ar
		p.startMessageProcessing(ar)
	}

	return nil
}

func (p *chatProcessorImpl) DeleteSession(ctx context.Context, key session.Key, options ...session.Option) error {
	return p.sessionService.DeleteSession(ctx, key, options...)
}

func (p *chatProcessorImpl) ListSessions(ctx context.Context, userKey session.UserKey, options ...session.Option) ([]*session.Session, error) {
	return p.sessionService.ListSessions(ctx, userKey, options...)
}

func (p *chatProcessorImpl) GetSession(ctx context.Context, key session.Key, options ...session.Option) (*session.Session, error) {
	return p.sessionService.GetSession(ctx, key, options...)
}

func (p *chatProcessorImpl) CreateSession(ctx context.Context, key session.Key, state session.StateMap, options ...session.Option) (*session.Session, error) {
	return p.sessionService.CreateSession(ctx, key, state, options...)
}

// SetOnMessageCallback sets the message callback function for the ChatProcessor.
func (p *chatProcessorImpl) SetOnMessageCallback(onMessage OnMessage) {
	p.onMessage = onMessage
}

// SetOnToolCallCallback sets the tool call callback function for the ChatProcessor.
func (p *chatProcessorImpl) SetOnToolCallCallback(onToolCall OnToolCall) {
	p.onToolCall = onToolCall
}

// SetOnRawEventCallback sets the raw event callback function for the ChatProcessor.
func (p *chatProcessorImpl) SetOnRawEventCallback(onRawEvent OnRawEvent) {
	p.onRawEvent = onRawEvent
}

// SetMessageInterceptor sets a message interceptor on the underlying message broker.
// The interceptor function will be called for every message sent between agents.
func (p *chatProcessorImpl) SetMessageInterceptor(interceptor messaging.Interceptor) {
	p.broker.SetMessageInterceptor(interceptor)
}

// GetAllAgentInfos returns a slice containing information about all registered agents.
func (p *chatProcessorImpl) GetAllAgentInfos() []*shared.AgentInfo {
	var infos []*shared.AgentInfo
	for _, agent := range p.agents {
		infos = append(infos, agent.Info())
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

// GetAgentByName returns the actual agent instance by name.
// Returns nil if no agent is found with the given name.
func (p *chatProcessorImpl) GetAgentByName(name string) shared.TheAgent {
	for _, agent := range p.agents {
		if agent.Name() == name {
			// Return the messaging wrapper (which implements shared.TheAgent and has send_message tool)
			return agent.wrapper
		}
	}

	return nil
}

// startMessageProcessing starts a goroutine to process incoming inter-agent messages for the given agent.
// It listens on the agent's message channel and forwards messages to the agent's runner.
func (p *chatProcessorImpl) startMessageProcessing(agent *AgentRunner) {
	go func() {
		// Get the message channel for this agent
		msgChan, err := p.broker.GetMessageChannel(agent.ID())
		if err != nil {
			logger.Log.Error("failed to get message channel for agent", zap.String("agent", agent.Name()), zap.Error(err))
			return
		}

		routingInfo := &messaging.RoutingInfo{
			Streaming: utils.BoolPtr(agent.IsStreaming()),
			ToAgentID: agent.ID(),
		}

		// Process incoming messages
		for msg := range msgChan {
			// Create a context for message processing
			ctx := context.Background()

			// Format the message content
			messageContent := fmt.Sprintf("Message from %s-[%s]:\n\n%s",
				p.GetAgentNameByID(msg.FromAgentID), msg.FromAgentID, msg.Content,
			)

			routingInfo.SessionID = msg.SessionID
			routingInfo.FromAgentID = msg.FromAgentID

			// Send to the agent's runner
			events, err := agent.Run(ctx, routingInfo, model.NewUserMessage(messageContent))
			if err != nil {
				logger.Log.Error("failed to process message for agent", zap.String("agent", agent.Name()), zap.Error(err))
				continue
			}

			// Process events from the agent's response
			go func(info *messaging.RoutingInfo) {
				for event := range events {
					p.processEvent(info.SwapFromTo(), event)
				}
			}(routingInfo)
		}
	}()
}

// SendMessage sends a message from one agent to another and returns a channel of events.
// The caller is responsible for processing the events from the returned channel.
func (p *chatProcessorImpl) SendMessage(
	ctx context.Context,
	routingInfo *messaging.RoutingInfo,
	message model.Message,
) (<-chan *event.Event, error) {
	ag, exists := p.agents[routingInfo.ToAgentID]
	if !exists {
		return nil, fmt.Errorf("agent %q not found", routingInfo.ToAgentID)
	}

	return ag.Run(ctx, routingInfo, message)
}

// SendMessageWithProcessing sends a message to an agent and automatically processes all resulting events.
// This method handles event processing internally and provides progress updates through callbacks.
func (p *chatProcessorImpl) SendMessageWithProcessing(
	ctx context.Context,
	routingInfo *messaging.RoutingInfo,
	message model.Message,
) error {
	agent, exists := p.agents[routingInfo.ToAgentID]
	if !exists {
		return fmt.Errorf("agent %q not found", routingInfo.ToAgentID)
	}

	p.onProgress(routingInfo, SystemMessageSending, "sending message to %s...", agent)

	events, err := agent.Run(ctx, routingInfo, message)
	if err != nil {
		return fmt.Errorf("failed to send message from %s to %s: %w", routingInfo.FromAgentID, routingInfo.ToAgentID, err)
	}

	p.onProgress(routingInfo, SystemMessageDelivered, "message delivered to %s - Processing...", agent)

	// Process events
	for event := range events {
		p.processEvent(routingInfo.SwapFromTo(), event)
	}

	p.onProgress(routingInfo, SystemMessageProcessed, "%s finished processing", agent)
	return nil
}

// processEvent processes a single event from an agent's response.
// It handles errors, assistant messages, and tool calls by invoking the appropriate callbacks.
func (p *chatProcessorImpl) processEvent(info *messaging.RoutingInfo, event *event.Event) {
	if event.Error != nil {
		p.onError(info, errors.New(event.Error.Message))
	}

	p.onRawEvent(info, event)

	if event.Response != nil && len(event.Response.Choices) > 0 {
		choice := event.Response.Choices[0]

		// Show reasoning content first if present (future-proof detection)
		if choice.Message.ReasoningContent != "" {
			p.onReasoningMessage(info, choice.Message.ReasoningContent)
		}

		// Show assistant messages
		if choice.Message.Role == model.RoleAssistant && choice.Message.Content != "" {
			p.onMessage(info, choice.Message.Content)
		}

		// Show tool calls (but suppress the generic "sending message" for cleaner output)
		if len(choice.Message.ToolCalls) > 0 {
			for _, toolCall := range choice.Message.ToolCalls {
				p.onToolCall(info, toolCall.Function)
			}
		}
	}
}
