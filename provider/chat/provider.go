package chat

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider/agent"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type chatProviderImpl struct {
	agentProvider agent.Provider
}

func New(i *do.Injector) (Provider, error) {
	agentProvider := do.MustInvoke[agent.Provider](i)

	return &chatProviderImpl{
		agentProvider: agentProvider,
	}, nil
}

func (p *chatProviderImpl) GetChat(ctx context.Context, agentID uuid.UUID, opts ...Option) (Chat, error) {
	agent, streaming, err := p.agentProvider.GetAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent with id %s", agentID)
	}

	var options Options = Options{
		appName:   "generic-chat",
		streaming: streaming,
		agent:     agent,
	}

	for _, opt := range opts {
		opt(&options)
	}

	chat := NewChat(options)
	return chat, nil
}
