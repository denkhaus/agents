package config

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/shared"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/chainagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/cycleagent"

	"trpc.group/trpc-go/trpc-agent-go/agent/parallelagent"
)

// createChainAgent creates a chain agent with the provided configuration
func (f *UnifiedAgentFactory) createChainAgent(ctx context.Context, environment EnvironmentName, agentConfig *AgentConfig) (agent.Agent, error) {
	options := []chainagent.Option{}

	// Add sub-agents
	if len(agentConfig.Setting.Agent.SubAgents) > 0 {
		subAgents, err := f.getSubAgents(ctx, environment, agentConfig.Setting.Agent.SubAgents)
		if err != nil {
			return nil, fmt.Errorf("failed to get sub agents: %w", err)
		}
		if len(subAgents) > 0 {
			options = append(options, chainagent.WithSubAgents(subAgents))
		}
	}

	// Add channel buffer size if specified
	if agentConfig.Setting.Agent.ChannelBufferSize > 0 {
		options = append(options, chainagent.WithChannelBufferSize(agentConfig.Setting.Agent.ChannelBufferSize))
	}

	// Create and return the chain agent
	return chainagent.New(agentConfig.Name, options...), nil
}

// createCycleAgent creates a cycle agent with the provided configuration
func (f *UnifiedAgentFactory) createCycleAgent(ctx context.Context, environment EnvironmentName, agentConfig *AgentConfig) (agent.Agent, error) {
	options := []cycleagent.Option{}

	// Add sub-agents
	if len(agentConfig.Setting.Agent.SubAgents) > 0 {
		subAgents, err := f.getSubAgents(ctx, environment, agentConfig.Setting.Agent.SubAgents)
		if err != nil {
			return nil, fmt.Errorf("failed to get sub agents: %w", err)
		}
		if len(subAgents) > 0 {
			options = append(options, cycleagent.WithSubAgents(subAgents))
		}
	}

	// Add max iterations if specified
	if agentConfig.Setting.Agent.MaxIterations > 0 {
		options = append(options, cycleagent.WithMaxIterations(agentConfig.Setting.Agent.MaxIterations))
	}

	// Add channel buffer size if specified
	if agentConfig.Setting.Agent.ChannelBufferSize > 0 {
		options = append(options, cycleagent.WithChannelBufferSize(agentConfig.Setting.Agent.ChannelBufferSize))
	}

	// Create and return the cycle agent
	return cycleagent.New(agentConfig.Name, options...), nil
}

// createParallelAgent creates a parallel agent with the provided configuration
func (f *UnifiedAgentFactory) createParallelAgent(ctx context.Context, environment EnvironmentName, agentConfig *AgentConfig) (agent.Agent, error) {
	options := []parallelagent.Option{}

	// Add sub-agents
	if len(agentConfig.Setting.Agent.SubAgents) > 0 {
		subAgents, err := f.getSubAgents(ctx, environment, agentConfig.Setting.Agent.SubAgents)
		if err != nil {
			return nil, fmt.Errorf("failed to get sub agents: %w", err)
		}
		if len(subAgents) > 0 {
			options = append(options, parallelagent.WithSubAgents(subAgents))
		}
	}

	// Add channel buffer size if specified
	if agentConfig.Setting.Agent.ChannelBufferSize > 0 {
		options = append(options, parallelagent.WithChannelBufferSize(agentConfig.Setting.Agent.ChannelBufferSize))
	}

	// Create and return the parallel agent
	return parallelagent.New(agentConfig.Name, options...), nil
}

// getSubAgents creates sub-agent instances based on their roles
func (f *UnifiedAgentFactory) getSubAgents(ctx context.Context, environment EnvironmentName, subAgentRoles []shared.AgentRole) ([]agent.Agent, error) {
	var subAgents []agent.Agent
	for _, role := range subAgentRoles {

		// Recursively create the sub-agent
		subAgent, err := f.CreateAgent(ctx, environment, role)
		if err != nil {
			return nil, fmt.Errorf("failed to create sub-agent %s: %w", role, err)
		}
		subAgents = append(subAgents, subAgent)
	}

	return subAgents, nil
}
