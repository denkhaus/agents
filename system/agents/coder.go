package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/agent"
	"github.com/denkhaus/agents/shared"
	"github.com/denkhaus/agents/tools/calculator"
	shelltoolset "github.com/denkhaus/agents/tools/shell"
	"github.com/denkhaus/agents/tools/time"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func CreateCoderAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentID := shared.AgentIDCoder
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)
	wkspceProvider := do.MustInvoke[provider.WorkspaceProvider](injector)

	wkspce, err := wkspceProvider.GetWorkspace(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace for agent [%s]: %w", agentID, err)
	}

	wkspcePath, err := wkspce.GetPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspacePath for agent [%s]: %w", agentID, err)
	}

	fileToolSet, err := CreateFileToolset(wkspcePath, false)
	if err != nil {
		return nil, err
	}

	shellToolSet, err := shelltoolset.NewToolSet()
	if err != nil {
		return nil, fmt.Errorf("failed to create shell toolset: %w", err)
	}

	timeTool := do.MustInvokeNamed[tool.Tool](injector, time.ToolName)
	calculatorTool := do.MustInvokeNamed[tool.Tool](injector, calculator.ToolName)

	coderAgent, err := agentProvider.GetAgent(ctx, agentID,
		agent.WithLLMAgentOptions(
			llmagent.WithTools([]tool.Tool{timeTool, calculatorTool}),
			llmagent.WithToolSets([]tool.ToolSet{shellToolSet, fileToolSet}),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return coderAgent, nil
}
