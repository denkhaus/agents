package commands

import (
	"fmt"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/multi/plugins"
	cli_chat "github.com/denkhaus/agents/pkg/multi/plugins/cli"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func RunChatCommand(
	ctx *cli.Context,
	configProvider config.ConfigProvider,
	agentFactory config.AgentFactory,
) error {

	envName := config.EnvironmentName(ctx.String("environment"))

	logger.Log.Info("Starting agents system",
		zap.String("version", ctx.App.Version),
		zap.String("environment", string(envName)),
	)

	envConfig, err := configProvider.LoadEnvironmentConfig(envName)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to load environment '%s': %v", envName, err), 1)
	}

	logger.Log.Info("Environment loaded successfully",
		zap.String("name", envConfig.Name),
		zap.String("description", envConfig.Description),
		zap.Int("agents", len(envConfig.Agents)),
		zap.Int("roles", len(envConfig.Roles)),
		zap.Bool("condenser_enabled", envConfig.Condenser.LoggingEnabled),
	)

	// Create all agents automatically
	agents, err := agentFactory.CreateAllAgentsInEnvironment(ctx.Context, envName)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to create agents: %v", err), 1)
	}

	if len(agents) == 0 {
		return cli.Exit(fmt.Sprintf("No agents found in environment '%s'", envName), 1)
	}

	processorOptions := []multi.ChatProcessorOption{
		multi.WithApplicationName(ctx.App.Name),
		multi.WithAgents(agents...),
	}

	chat := cli_chat.NewCLIMultiAgentChat(
		plugins.WithSessionID(uuid.New()),
		plugins.WithProcessorOptions(processorOptions...),
	)

	return chat.Start(ctx.Context)
}
