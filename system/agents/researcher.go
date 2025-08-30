package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/agent"
	"github.com/denkhaus/agents/provider/tools"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
)

func CreateResearcherAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentID := shared.AgentIDResearcher
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)
	toolProvider := do.MustInvoke[tools.ToolProvider](injector)

	// Get tools from ToolProvider based on agent ID
	agentTools, agentToolSets, err := toolProvider.GetToolsForAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tools for researcher agent: %w", err)
	}

	agent, err := agentProvider.GetAgent(ctx, agentID,
		agent.WithLLMAgentOptions(
			llmagent.WithTools(agentTools),
			llmagent.WithToolSets(agentToolSets),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return agent, nil
}
