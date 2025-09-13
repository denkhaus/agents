package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/multi/plugins/web"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// validateConfigCommand validates the system configuration
func (a *App) validateConfigCommand(c *cli.Context) error {
	logger.Log.Info("Validating system configuration...")

	if err := a.configProvider.ValidateConfiguration(); err != nil {
		return cli.Exit(fmt.Sprintf("Configuration validation failed: %v", err), 1)
	}

	// Load and validate system config
	systemConfig, err := a.configProvider.LoadSystemConfig()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to load system config: %v", err), 1)
	}

	fmt.Printf("Configuration validation successful!\n\n")
	fmt.Printf("System configuration:\n")
	fmt.Printf("  Environments: %d\n", len(systemConfig.Environments))

	for envName, envConfig := range systemConfig.Environments {
		fmt.Printf("  * %s: %d agents, %d roles\n", envName, len(envConfig.Agents), len(envConfig.Roles))
	}

	return nil
}

// listEnvironmentsCommand shows available environments
func (a *App) listEnvironmentsCommand(c *cli.Context) error {
	systemConfig, err := a.configProvider.LoadSystemConfig()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to load system config: %v", err), 1)
	}

	fmt.Printf("Available environments (%d):\n\n", len(systemConfig.Environments))

	// Sort environments for consistent output
	envNames := make([]string, 0, len(systemConfig.Environments))
	for name := range systemConfig.Environments {
		envNames = append(envNames, string(name))
	}
	sort.Strings(envNames)

	for _, name := range envNames {
		envConfig := systemConfig.Environments[config.EnvironmentName(name)]
		fmt.Printf("* %s\n", name)
		fmt.Printf("  Description: %s\n", envConfig.Description)
		fmt.Printf("  Agents: %d\n", len(envConfig.Agents))
		fmt.Printf("  Role mappings: %d\n", len(envConfig.Roles))

		if envConfig.Condenser.LoggingEnabled {
			fmt.Printf("  Condenser: enabled (threshold: %.2f)\n", envConfig.Condenser.TriggerThreshold)
		} else {
			fmt.Printf("  Condenser: disabled\n")
		}
		fmt.Println()
	}

	return nil
}

func (a *App) RunCommand(ctx *cli.Context) error {
	return a.runCommand(
		ctx.Context,
		ctx.String("environment"),
		ctx.String("server-addr"),
	)
}

// runCommand starts the multi-agent system
func (a *App) runCommand(ctx context.Context, environmentName, serverAddr string) error {

	envName := config.EnvironmentName(environmentName)

	logger.Log.Info("Starting agents system",
		zap.String("version", appVersion),
		zap.String("environment", string(envName)),
	)

	// Validate environment exists
	envConfig, err := a.configProvider.LoadEnvironmentConfig(envName)
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
	agents, err := a.agentFactory.CreateAllAgentsInEnvironment(ctx, envName)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to create agents: %v", err), 1)
	}

	if len(agents) == 0 {
		return cli.Exit(fmt.Sprintf("No agents found in environment '%s'", envName), 1)
	}

	// Log created agents
	for _, ag := range agents {
		logger.Log.Info("Agent ready",
			zap.String("name", ag.Info().Name),
			zap.String("role", string(ag.GetRole())),
			zap.Any("id", ag.GetID()),
		)
	}

	g := new(errgroup.Group)
	g.Go(func() error {
		return a.startServer(ctx, envConfig, agents, serverAddr)
	})

	return g.Wait()
}

// startServer starts the debug HTTP server
func (a *App) startServer(
	ctx context.Context,
	envConfig *config.EnvironmentConfig,
	agents []shared.TheAgent,
	addr string,
) error {

	condenserService, err := a.createCondenser(ctx, envConfig)
	if err != nil {
		return err
	}

	processor := multi.NewChatProcessor(
		uuid.New(),
		multi.WithSessionService(condenserService),
		multi.WithApplicationName(fmt.Sprintf("%s-%s", appName, envConfig.Name)),
		multi.WithAgents(agents...),
	)

	server := web.New(
		web.WithChatProcessor(processor),
		web.WithLogger(logger.Log),
	)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: server.Handler(),
	}

	logger.Log.Info("server is starting", zap.String("address", addr))

	// Start server in goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("server failed", zap.Error(err))
		}
	}()

	// Wait for context cancellation (shutdown signal)
	<-ctx.Done()
	logger.Log.Info("shutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("server shutdown failed", zap.Error(err))
		return err
	}

	logger.Log.Info("server shutdown complete")
	return nil
}
