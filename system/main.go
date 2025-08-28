package main

import (
	"context"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/multi"
	"github.com/denkhaus/agents/multi/plugins"
	"github.com/denkhaus/agents/shared"
	"github.com/denkhaus/agents/system/agents"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func startup(ctx context.Context) error {

	injector := di.NewContainer()

	projectManager, err := agents.CreateProjectManagerAgent(ctx, injector)
	if err != nil {
		return err
	}

	chat := plugins.NewCLIMultiAgentChat(
		plugins.WithDisplayWidth(120),
		plugins.WithProcessorOptions(
			multi.WithSessionID(uuid.New()),
			multi.WithApplicationName("denkhaus-multi-chat"),
			multi.WithAgents(
				shared.NewHumanAgent(shared.AgentInfoHuman),
				projectManager,
			),
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
