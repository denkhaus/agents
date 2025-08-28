package agents

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/agent"
	"github.com/denkhaus/agents/shared"
	"github.com/denkhaus/agents/tools/project"
	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/file"
)

func CreateFileToolset(workspacePath string, readOnly bool) (toolset tool.ToolSet, err error) {
	options := []file.Option{
		file.WithBaseDir(workspacePath),
	}

	if readOnly {
		// Create readonly file operation tools.
		options = append(options,
			file.WithListFileEnabled(true),
			file.WithReadFileEnabled(true),
			file.WithReplaceContentEnabled(false),
			file.WithSaveFileEnabled(false),
			file.WithSearchFileEnabled(true),
			file.WithSearchContentEnabled(true),
		)
	}

	toolset, err = file.NewToolSet(options...)
	if err != nil {
		return nil, fmt.Errorf("create file tool set: %w", err)
	}

	return toolset, err
}

func CreateProjectManagerAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentID := shared.AgentIDProjectManager
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)
	workspaceProvider := do.MustInvoke[provider.WorkspaceProvider](injector)

	projectManagerWorkspace, err := workspaceProvider.GetWorkspace(shared.AgentIDProjectManager)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace for agent [%s]: %w", agentID, err)
	}

	workspacePath, err := projectManagerWorkspace.GetPath()
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

	projectManagerAgent, err := agentProvider.GetAgent(ctx, shared.AgentIDProjectManager,
		agent.WithLLMAgentOptions(
			llmagent.WithToolSets([]tool.ToolSet{projectManagerToolSet, fileToolSet}),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create project manager agent: %w", err)
	}

	return projectManagerAgent, nil
}
