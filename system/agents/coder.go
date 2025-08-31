package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/config"
	"github.com/denkhaus/agents/pkg/provider"
	"github.com/denkhaus/agents/pkg/provider/agent"
	"github.com/denkhaus/agents/pkg/tools"
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

	fileFactory, err := do.InvokeNamed[tools.ToolSetFactoryFunc](injector, file.ToolSetName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file factory from di for agent [%s]: %w", agentID, err)
	}

	workspacePath, err := configProvider.GetWorkspacePath()
	if err != nil {
		return nil, err
	}

	fileToolSet, err := fileFactory(tools.ConfigPayload{
		"WorkspacePath": workspacePath,
		"ReadOnly":      false,
	})

	if err != nil {
		return nil, err
	}

	shellFactory, err := do.InvokeNamed[tools.ToolSetFactoryFunc](injector, shell.ToolSetName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shell factory from di for agent [%s]: %w", agentID, err)
	}

	shellToolSet, err := shellFactory(tools.ConfigPayload{
		"baseDir":               workspacePath,
		"executeCommandEnabled": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create shell toolset: %w", err)
	}

	projectFactory, err := do.InvokeNamed[tools.ToolSetFactoryFunc](injector, project.ToolSetName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve project factory from di for agent [%s]: %w", agentID, err)
	}

	readOnlyProjectManagerToolSet, err := projectFactory(tools.ConfigPayload{
		"isReadOnly": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create project manager toolset: %w", err)
	}

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
