package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/agent"
	"github.com/denkhaus/agents/shared"
	"github.com/denkhaus/agents/tools/calculator"
	"github.com/denkhaus/agents/tools/project"
	"github.com/denkhaus/agents/tools/time"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func CreateProjectManagerAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentID := shared.AgentIDProjectManager
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)
	workspaceProvider := do.MustInvoke[provider.WorkspaceProvider](injector)

	wkspce, err := workspaceProvider.GetWorkspace(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace for agent [%s]: %w", agentID, err)
	}

	workspacePath, err := wkspce.GetPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspacePath for agent [%s]: %w", agentID, err)
	}

	fileToolSet, err := CreateFileToolset(workspacePath, true)
	if err != nil {
		return nil, err
	}

	projectManagerToolSet, err := project.NewToolSet()
	if err != nil {
		return nil, fmt.Errorf("failed to create project manager toolset: %w", err)
	}

	timeTool := do.MustInvokeNamed[tool.Tool](injector, time.ToolName)
	calculatorTool := do.MustInvokeNamed[tool.Tool](injector, calculator.ToolName)

	projectManagerAgent, err := agentProvider.GetAgent(ctx, agentID,
		agent.WithLLMAgentOptions(
			llmagent.WithTools([]tool.Tool{timeTool, calculatorTool}),
			llmagent.WithToolSets([]tool.ToolSet{projectManagerToolSet, fileToolSet}),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return projectManagerAgent, nil
}
