package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/provider"
	"github.com/denkhaus/agents/pkg/provider/agent"
	"github.com/denkhaus/agents/pkg/tools/calculator"
	"github.com/denkhaus/agents/pkg/tools/fetch"
	"github.com/denkhaus/agents/pkg/tools/tavily"
	"github.com/denkhaus/agents/pkg/tools/time"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func CreateResearcherAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentID := shared.AgentIDResearcher
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)

	timeTool := do.MustInvokeNamed[tool.Tool](injector, time.ToolName)
	calculatorTool := do.MustInvokeNamed[tool.Tool](injector, calculator.ToolName)
	fetchTool := do.MustInvokeNamed[tool.Tool](injector, fetch.ToolName)
	tavilyTool := do.MustInvokeNamed[tool.ToolSet](injector, tavily.ToolSetName)

	agent, err := agentProvider.GetAgent(ctx, agentID,
		agent.WithLLMAgentOptions(
			llmagent.WithTools([]tool.Tool{timeTool, calculatorTool, fetchTool}),
			llmagent.WithToolSets([]tool.ToolSet{tavilyTool}),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return agent, nil
}
