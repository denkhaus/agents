package main

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/multi"
	"github.com/denkhaus/agents/multi/plugins"
	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/agent"
	"github.com/denkhaus/agents/shared"
	"github.com/denkhaus/agents/tools/project"
	"github.com/google/uuid"
	"github.com/samber/do"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/file"
)

func startup(ctx context.Context) error {

	injector := di.NewContainer()
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)
	workspaceProvider := do.MustInvoke[provider.WorkspaceProvider](injector)

	projectManagerWorkspace, err := workspaceProvider.GetWorkspace(shared.AgentIDProjectManager)
	if err != nil {
		return fmt.Errorf("failed to get workspace for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	path, err := projectManagerWorkspace.GetWorkspacePath()
	if err != nil {
		if err != nil {
			return fmt.Errorf("failed to get workspacePath for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
		}
	}
	// Create file operation tools.
	fileToolSet, err := file.NewToolSet(
		file.WithBaseDir(path),
	)
	if err != nil {
		return nil, fmt.Errorf("create file tool set: %w", err)
	}

	projectManagerToolSet, err := project.NewToolSet()
	if err != nil {
		return fmt.Errorf("failed to create project manager toolset: %w", err)
	}

	projectManagerAgent, err := agentProvider.GetAgent(ctx, shared.AgentIDProjectManager,
		agent.WithLLMAgentOptions(
			llmagent.WithToolSets([]tool.ToolSet{projectManagerToolSet}),
		),
	)

	if err != nil {
		return fmt.Errorf("failed to create project manager agent: %w", err)
	}

	chat := plugins.NewCLIMultiAgentChat(
		plugins.WithProcessorOptions(
			multi.WithSessionID(uuid.New()),
			multi.WithHumanAgent(shared.NewHumanAgent(shared.AgentInfoHuman)),
			multi.WithApplicationName("denkhaus-multi-chat"),
			multi.WithAgents(projectManagerAgent),
		),
	)

	return chat.Start(ctx)
}

func main() {
	if err := startup(context.Background()); err != nil {
		logger.Log.Fatal("application error", zap.Error(err))
	}
}

// func startup(ctx context.Context) error {

// 	injector := di.NewContainer()
// 	multiAgentSystem := do.MustInvoke[multi.MultiAgentSystem](injector)

// 	chat, err := multiAgentSystem.CreateProjectManagerChat(ctx, injector,
// 		chat.WithAppName("denkhaus-system-chat"),
// 		chat.WithSessionID(uuid.New()),
// 		chat.WithUserID(uuid.New()),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	return multiAgentSystem.EnterChat(ctx, chat)
// }

// func main() {
// 	if err := startup(context.Background()); err != nil {
// 		log.Fatal(err)
// 	}
// }
