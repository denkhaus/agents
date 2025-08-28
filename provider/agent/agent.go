// Package agent provides an implementation of the AgentProvider interface
// for creating different types of agents based on configuration.
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
	"trpc.group/trpc-go/trpc-agent-go/agent/parallelagent"
)

// agentProviderImpl implements the AgentProvider interface
type agentProviderImpl struct {
	settingsProvider provider.SettingsProvider
}

// New creates a new AgentProvider instance using dependency injection
func New(i *do.Injector) (provider.AgentProvider, error) {
	settingsProvider := do.MustInvoke[provider.SettingsProvider](i)
	return &agentProviderImpl{
		settingsProvider: settingsProvider,
	}, nil
}

// getDefaultAgent creates a default LLM agent with the provided configuration
func (p *agentProviderImpl) getDefaultAgent(
	ctx context.Context,
	agentConfig provider.AgentConfiguration,
	opt ...llmagent.Option,
) (agent.Agent, error) {
	options, err := agentConfig.GetDefaultOptions(ctx, p, opt...)
	if err != nil {
		return nil, err
	}

	return llmagent.New(
		agentConfig.GetName(), options...,
	), nil
}

// getChainAgent creates a chain agent with the provided configuration
func (p *agentProviderImpl) getChainAgent(
	ctx context.Context,
	agentConfig provider.AgentConfiguration,
	opt ...chainagent.Option,
) (agent.Agent, error) {
	options, err := agentConfig.GetChainOptions(ctx, p, opt...)
	if err != nil {
		return nil, err
	}

	return chainagent.New(
		agentConfig.GetName(), options...,
	), nil
}

// getCycleAgent creates a cycle agent with the provided configuration
func (p *agentProviderImpl) getCycleAgent(
	ctx context.Context,
	agentConfig provider.AgentConfiguration,
	opt ...cycleagent.Option,
) (agent.Agent, error) {
	options, err := agentConfig.GetCycleOptions(ctx, p, opt...)
	if err != nil {
		return nil, err
	}

	return cycleagent.New(
		agentConfig.GetName(), options...,
	), nil
}

// getParallelAgent creates a parallel agent with the provided configuration
func (p *agentProviderImpl) getParallelAgent(
	ctx context.Context,
	agentConfig provider.AgentConfiguration,
	opt ...parallelagent.Option,
) (agent.Agent, error) {
	options, err := agentConfig.GetParallelOptions(ctx, p, opt...)
	if err != nil {
		return nil, err
	}

	return parallelagent.New(
		agentConfig.GetName(), options...,
	), nil
}

// GetAgent creates and returns an agent based on the provided agent ID and configuration.
// It also returns whether streaming is enabled for the agent and any error that occurred.
func (p *agentProviderImpl) GetAgent(
	ctx context.Context,
	agentID uuid.UUID,
	opt ...provider.AgentProviderOption,
) (shared.TheAgent, error) {

	var options provider.AgentProviderOptions
	for _, o := range opt {
		o(&options)
	}

	agentConfig, err := p.settingsProvider.GetAgentConfiguration(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent settings for agent %s", agentID)
	}

	var ag agent.Agent
	switch agentConfig.GetType() {
	case shared.AgentTypeDefault:
		ag, err = p.getDefaultAgent(ctx, agentConfig, options.LLMOpt...)
	case shared.AgentTypeChain:
		ag, err = p.getChainAgent(ctx, agentConfig, options.ChainOpt...)
	case shared.AgentTypeCycle:
		ag, err = p.getCycleAgent(ctx, agentConfig, options.CycleOpt...)
	case shared.AgentTypeParallel:
		ag, err = p.getParallelAgent(ctx, agentConfig, options.ParallelOpt...)
	default:
		ag, err = p.getDefaultAgent(ctx, agentConfig, options.LLMOpt...)
	}

	return shared.NewAgent(
		ag,
		agentID,
		agentConfig.IsStreamingEnabled(),
	), err
}
