package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/multi/plugins/web/schema"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

func Test_Processor(t *testing.T) {

	os.Setenv("AGENTS_WORKSPACE_PATH", "/home/denkhaus/dev/gomodules/agents/test_workspace")
	os.Setenv("AGENTS_CONFIG_PATH", "/home/denkhaus/dev/gomodules/agents/config")
	os.Setenv("AGENTS_DATABASE_URL", "postgres://agents:agents@localhost:6888/agents?sslmode=disable")

	app, err := NewApp()
	assert.NoError(t, err)

	ctx := context.Background()
	environmentName := "production"
	sessionID := uuid.New()

	envName := config.EnvironmentName(environmentName)

	logger.Log.Info("Starting agents system",
		zap.String("version", appVersion),
		zap.String("environment", string(envName)),
	)

	// Validate environment exists
	envConfig, err := app.configProvider.LoadEnvironmentConfig(envName)
	if err != nil {
		assert.NoError(t, err)
	}

	logger.Log.Info("Environment loaded successfully",
		zap.String("name", envConfig.Name),
		zap.String("description", envConfig.Description),
		zap.Int("agents", len(envConfig.Agents)),
		zap.Int("roles", len(envConfig.Roles)),
		zap.Bool("condenser_enabled", envConfig.Condenser.LoggingEnabled),
	)

	// Create all agents automatically
	agents, err := app.agentFactory.CreateAllAgentsInEnvironment(ctx, envName)
	if err != nil {
		assert.NoError(t, err)
	}

	if len(agents) == 0 {
		return
	}

	// Log created agents
	for _, ag := range agents {
		logger.Log.Info("Agent ready",
			zap.String("name", ag.Info().Name),
			zap.String("role", string(ag.GetRole())),
			zap.String("id", ag.ID().String()),
		)
	}

	condenserService, err := app.createCondenser(ctx, envConfig)
	if err != nil {
		assert.NoError(t, err)
	}

	processor := multi.NewChatProcessor(
		multi.WithSessionService(condenserService),
		multi.WithApplicationName(fmt.Sprintf("%s-%s", appName, envConfig.Name)),
		multi.WithAgents(agents...),
	)

	evts, err := processor.SendMessage(
		ctx,
		shared.AgentIDHuman,
		shared.AgentIDCoder,
		sessionID,
		model.NewUserMessage("hi, what's your name?"),
	)

	if err != nil {
		assert.NoError(t, err)
	}

	for ev := range evts {
		llmEvent, err := schema.NewLLMEvent(ev, false)
		if err != nil {
			logger.Log.Error("failed to create llm event", zap.Error(err))
		}

		spew.Dump(llmEvent)
	}
}
