package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

const (
	serverDefaultAddr = ":6999"
	appName           = "agents-system"
	appUsage          = "Multi-agent system with environment-based configuration"
	appVersion        = "1.0.0"
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
			&cli.StringFlag{
				Name:  "server-addr",
				Value: serverDefaultAddr,
				Usage: "Server listen address",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"start", "r"},
				Usage:   "Start the multi-agent system",
				Action:  app.RunCommand,
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
		Action:         app.RunCommand, // Default action when no command specified
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
	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Log.Info("Received shutdown signal, initiating graceful shutdown...")
		cancel()
	}()

	// Create application instance
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Ensure cleanup on exit
	defer func() {
		logger.Log.Info("Cleaning up application resources...")
		app.Close()
		// Give some time for cleanup
		time.Sleep(100 * time.Millisecond)
	}()

	// Run the CLI application with context
	cliApp := createCLIApp(app)

	// Create a context-aware version of os.Args
	if err := cliApp.RunContext(ctx, os.Args); err != nil {
		logger.Log.Fatal("Application error", zap.Error(err))
	}
}
