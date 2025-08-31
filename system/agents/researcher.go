package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/provider"
	"github.com/denkhaus/agents/pkg/provider/agent"
	"github.com/denkhaus/agents/pkg/tools"
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

	timeFactory := do.MustInvokeNamed[tools.ToolFactoryFunc](injector, time.ToolName)
	timeTool, err := timeFactory(tools.ConfigPayload{})
	if err != nil {
		return nil, fmt.Errorf("failed to create time tool: %w", err)
	}

	calculatorFactory := do.MustInvokeNamed[tools.ToolFactoryFunc](injector, calculator.ToolName)
	calculatorTool, err := calculatorFactory(tools.ConfigPayload{})
	if err != nil {
		return nil, fmt.Errorf("failed to create calculator tool: %w", err)
	}

	fetchFactory := do.MustInvokeNamed[tools.ToolFactoryFunc](injector, fetch.ToolName)
	fetchTool, err := fetchFactory(tools.ConfigPayload{})
	if err != nil {
		return nil, fmt.Errorf("failed to create fetch tool: %w", err)
	}

	tavilyFactory := do.MustInvokeNamed[tools.ToolSetFactoryFunc](injector, tavily.ToolSetName)
	tavilyToolSet, err := tavilyFactory(tools.ConfigPayload{})
	if err != nil {
		return nil, fmt.Errorf("failed to create tavily toolset: %w", err)
	}

	agent, err := agentProvider.GetAgent(ctx, agentID,
		agent.WithLLMAgentOptions(
			llmagent.WithTools([]tool.Tool{timeTool, calculatorTool, fetchTool}),
			llmagent.WithToolSets([]tool.ToolSet{tavilyToolSet}),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return agent, nil
}
