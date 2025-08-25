package agent

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/chainagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/cycleagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
)

type agentProviderImpl struct {
	settingsProvider provider.SettingsProvider
}

func New(i *do.Injector) (provider.AgentProvider, error) {
	settingsProvider := do.MustInvoke[provider.SettingsProvider](i)
	return &agentProviderImpl{
		settingsProvider: settingsProvider,
	}, nil
}

func (p *agentProviderImpl) getDefaultAgent(
	ctx context.Context,
	agentConfig provider.AgentConfiguration,
) (agent.Agent, error) {
	options, err := agentConfig.GetDefaultOptions(ctx, p)
	if err != nil {
		return nil, err
	}

	return llmagent.New(
		agentConfig.GetName(), options...,
	), nil
}

func (p *agentProviderImpl) getChainAgent(
	ctx context.Context,
	agentConfig provider.AgentConfiguration,
) (agent.Agent, error) {
	options, err := agentConfig.GetChainOptions(ctx, p)
	if err != nil {
		return nil, err
	}

	return chainagent.New(
		agentConfig.GetName(), options...,
	), nil
}

func (p *agentProviderImpl) getCycleAgent(
	ctx context.Context,
	agentConfig provider.AgentConfiguration,
) (agent.Agent, error) {
	options, err := agentConfig.GetCycleOptions(ctx, p)
	if err != nil {
		return nil, err
	}

	return cycleagent.New(
		agentConfig.GetName(), options...,
	), nil
}

func (p *agentProviderImpl) GetAgent(
	ctx context.Context,
	agentID uuid.UUID,
) (agent agent.Agent, isStreamingEnabled bool, err error) {

	agentConfig, err := p.settingsProvider.GetAgentConfiguration(agentID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get agent settings for agent %s", agentID)
	}

	switch agentConfig.GetType() {
	case shared.AgentTypeDefault:
		agent, err = p.getDefaultAgent(ctx, agentConfig)

	case shared.AgentTypeChain:
		agent, err = p.getChainAgent(ctx, agentConfig)

	case shared.AgentTypeCycle:
		agent, err = p.getCycleAgent(ctx, agentConfig)

	default:
		agent, err = p.getDefaultAgent(ctx, agentConfig)
	}

	return agent, agentConfig.IsStreamingEnabled(), err
}
