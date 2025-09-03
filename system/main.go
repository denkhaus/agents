package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/multi"
	"github.com/denkhaus/agents/multi/plugins"
	multicli "github.com/denkhaus/agents/multi/plugins/cli"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/shared"

	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

const (
	debugServerDefaultAddr = ":6999"
	appName                = "agents-system"
	appUsage               = "Multi-agent system with environment-based configuration"
	appVersion             = "1.0.0"
)

// App represents the main application with dependency injection
type App struct {
	injector       *do.Injector
	configProvider config.ConfigProvider
	agentFactory   config.AgentFactory
}

// NewApp creates a new application instance with dependency injection
func NewApp() (*App, error) {
	injector := di.NewContainer()

	configProvider := do.MustInvoke[config.ConfigProvider](injector)
	agentFactory := do.MustInvoke[config.AgentFactory](injector)

	return &App{
		injector:       injector,
		configProvider: configProvider,
		agentFactory:   agentFactory,
	}, nil
}

// Close cleans up application resources
func (a *App) Close() {
	if a.injector != nil {
		a.injector.Shutdown()
	}
}

// startChatSystem initializes and starts the chat interface
func (a *App) startChatSystem(ctx context.Context, agents []shared.TheAgent, envConfig *config.EnvironmentConfig) error {
	chatAgents := []shared.TheAgent{shared.NewHumanAgent(shared.AgentInfoHuman)}
	chatAgents = append(chatAgents, agents...)

	condenserService, err := a.createCondenser(ctx, envConfig)
	if err != nil {
		return err
	}

	chat := multicli.NewCLIMultiAgentChat(
		plugins.WithProcessorOptions(
			multi.WithSessionService(condenserService),
			multi.WithSessionID(uuid.New()),
			multi.WithApplicationName(fmt.Sprintf("%s-%s", appName, envConfig.Name)),
			multi.WithAgents(chatAgents...),
		),
	)

	logger.Log.Info("Chat system ready",
		zap.Int("agents", len(chatAgents)),
	)

	return chat.Start(ctx)
}

// createCLIApp creates and configures the CLI application
func createCLIApp(app *App) *cli.App {
	return &cli.App{
		Name:    appName,
		Usage:   appUsage,
		Version: appVersion,
		Authors: []*cli.Author{
			{
				Name:  "denkhaus",
				Email: "denkhaus@example.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "environment",
				Aliases: []string{"env", "e"},
				Value:   "production",
				Usage:   "Environment to run (development, production, experiment)",
				EnvVars: []string{"AGENTS_ENVIRONMENT"},
			},
			&cli.BoolFlag{
				Name:    "debug-server",
				Aliases: []string{"d"},
				Value:   true,
				Usage:   "Enable debug HTTP server",
			},
			&cli.StringFlag{
				Name:  "debug-addr",
				Value: debugServerDefaultAddr,
				Usage: "Debug server listen address",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"start", "r"},
				Usage:   "Start the multi-agent system",
				Action:  app.runCommand,
			},
			{
				Name:    "list-environments",
				Aliases: []string{"list", "ls", "envs"},
				Usage:   "List available environments",
				Action:  app.listEnvironmentsCommand,
			},
			{
				Name:    "validate",
				Aliases: []string{"check", "v"},
				Usage:   "Validate system configuration",
				Action:  app.validateConfigCommand,
			},
		},
		DefaultCommand: "run",
		Action:         app.runCommand, // Default action when no command specified
		Before: func(c *cli.Context) error {
			// Setup logging based on environment
			env := c.String("environment")
			if env == "development" {
				logger.Log.Info("Running in development mode with enhanced logging")
			}
			return nil
		},
		After: func(c *cli.Context) error {
			// Cleanup
			app.Close()
			return nil
		},
	}
}

func main() {
	// Create application instance
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run the CLI application
	if err := createCLIApp(app).Run(os.Args); err != nil {
		logger.Log.Fatal("Application error", zap.Error(err))
	}
}
