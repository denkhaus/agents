package main

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/examples/messaging/multi"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/agent"
	"github.com/denkhaus/agents/shared"
	"github.com/denkhaus/agents/tools/project"
	"github.com/samber/do"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func startup(ctx context.Context) error {

	injector := di.NewContainer()
	agentProvider := do.MustInvoke[provider.AgentProvider](injector)

	projectManagerToolSet, err := project.NewToolSet()
	if err != nil {
		return fmt.Errorf("failed to create project manager toolset: %w", err)
	}

	agentProvider.GetAgent(ctx, shared.AgentIDProjectManager, agent.WithLLMAgentOptions(
		llmagent.WithToolSets([]tool.ToolSet{projectManagerToolSet}),
	))

	chat := multi.NewMultiAgentChat("mult")
	return nil
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
