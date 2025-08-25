package agent

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider/settings"
	"github.com/google/uuid"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
)

type Provider interface {
	GetAgent(ctx context.Context, agentID uuid.UUID) (agent.Agent, error)
}

type agentProviderImpl struct {
	settingsProvider settings.Provider
}

func New(i *do.Injector) (Provider, error) {
	settingsProvider := do.MustInvoke[settings.Provider](i)
	return &agentProviderImpl{
		settingsProvider: settingsProvider,
	}, nil
}

func (p *agentProviderImpl) GetAgent(ctx context.Context, agentID uuid.UUID) (agent.Agent, error) {

	agentConfig, err := p.settingsProvider.GetAgentConfiguration(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent settings for agent %s", agentID)
	}

	options, err := agentConfig.GetOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get options for agent %s", agentID)
	}

	newAgent := llmagent.New(
		agentConfig.GetAgentName(), options...,
	)

	return newAgent, nil

}
