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

	coder, err := agents.CreateCoderAgent(ctx, injector)
	if err != nil {
		return err
	}

	// Enhanced Bubble Tea Chat with real LLM calls and spinners
	chat := plugins.NewEnhancedBubbleTeaChatPlugin(
		plugins.WithProcessorOptions(
			multi.WithSessionID(uuid.New()),
			multi.WithApplicationName("denkhaus-multi-agent"),
			multi.WithAgents(
				shared.NewHumanAgent(shared.AgentInfoHuman),
				projectManager,
				coder,
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
