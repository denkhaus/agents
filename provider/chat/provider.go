package chat

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type chatProviderImpl struct {
	agentProvider provider.AgentProvider
}

func New(i *do.Injector) (provider.ChatProvider, error) {
	agentProvider := do.MustInvoke[provider.AgentProvider](i)

	return &chatProviderImpl{
		agentProvider: agentProvider,
	}, nil
}

func (p *chatProviderImpl) GetChat(
	ctx context.Context,
	agentID uuid.UUID,
	opts ...provider.ChatProviderOption,
) (provider.Chat, error) {

	var options provider.ChatProviderOptions = provider.ChatProviderOptions{
		AppName: "generic-chat",
	}

	for _, opt := range opts {
		opt(&options)
	}

	agent, err := p.agentProvider.GetAgent(ctx, agentID, options.AgentProviderOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent with id %s", agentID)
	}

	chat := NewChat(agent, options)
	return chat, nil
}
