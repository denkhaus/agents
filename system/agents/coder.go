package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/config"
	"github.com/denkhaus/agents/pkg/provider"
	"github.com/denkhaus/agents/pkg/provider/agent"
	"github.com/denkhaus/agents/pkg/tools/calculator"
	"github.com/denkhaus/agents/pkg/tools/fetch"
	"github.com/denkhaus/agents/pkg/tools/file"
	"github.com/denkhaus/agents/pkg/tools/project"
	"github.com/denkhaus/agents/pkg/tools/shell"
	"github.com/denkhaus/agents/pkg/tools/time"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func CreateCoderAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentID := shared.AgentIDCoder
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
		file.WithReadOnly(false),
		file.WithWorkspacePath(workspacePath),
	)

	if err != nil {
		return nil, err
	}

	shellFactory, err := do.InvokeNamed[shell.FactoryFunc](injector, shell.ToolSetName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shell factory from di for agent [%s]: %w", agentID, err)
	}

	shellToolSet, err := shellFactory()
	if err != nil {
		return nil, fmt.Errorf("failed to create shell toolset: %w", err)
	}

	projectFactory, err := do.InvokeNamed[project.FactoryFunc](injector, project.ToolSetName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve project factory from di for agent [%s]: %w", agentID, err)
	}

	readOnlyProjectManagerToolSet, err := projectFactory(
		project.WithReadOnly(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project manager toolset: %w", err)
	}

	timeTool := do.MustInvokeNamed[tool.Tool](injector, time.ToolName)
	calculatorTool := do.MustInvokeNamed[tool.Tool](injector, calculator.ToolName)
	fetchTool := do.MustInvokeNamed[tool.Tool](injector, fetch.ToolName)

	coderAgent, err := agentProvider.GetAgent(ctx, agentID,
		agent.WithLLMAgentOptions(
			llmagent.WithTools([]tool.Tool{timeTool, calculatorTool, fetchTool}),
			llmagent.WithToolSets([]tool.ToolSet{
				shellToolSet,
				fileToolSet,
				readOnlyProjectManagerToolSet,
			}),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return coderAgent, nil
}
