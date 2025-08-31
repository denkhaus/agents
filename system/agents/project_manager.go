package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/config"
	"github.com/denkhaus/agents/pkg/provider"
	"github.com/denkhaus/agents/pkg/provider/agent"
	"github.com/denkhaus/agents/pkg/tools/calculator"
	"github.com/denkhaus/agents/pkg/tools/file"
	"github.com/denkhaus/agents/pkg/tools/project"
	"github.com/denkhaus/agents/pkg/tools/time"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func CreateProjectManagerAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentID := shared.AgentIDProjectManager
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)
	configProvider := do.MustInvoke[config.Service](injector)

	fileFactory, err := do.InvokeNamed[file.FactoryFunc](injector, file.ToolSetName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file factory from di for agent [%s]: %w", agentID, err)
	}

	workspacePath, err := configProvider.GetWorkspacePath()
	if err != nil {
		return nil, err
	}

	fileToolSet, err := fileFactory(
		file.WithReadOnly(true),
		file.WithWorkspacePath(workspacePath),
	)

	if err != nil {
		return nil, err
	}

	projectFactory, err := do.InvokeNamed[project.FactoryFunc](injector, project.ToolSetName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve project factory from di for agent [%s]: %w", agentID, err)
	}

	projectManagerToolSet, err := projectFactory()
	if err != nil {
		return nil, err
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
