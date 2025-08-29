package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/agent"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/duckduckgo"
)

func CreateResearcherAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentID := shared.AgentIDResearcher
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)

	searchTool := duckduckgo.NewTool()
	agent, err := agentProvider.GetAgent(ctx, agentID,
		agent.WithLLMAgentOptions(
			llmagent.WithTools([]tool.Tool{searchTool}),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return agent, nil
}
