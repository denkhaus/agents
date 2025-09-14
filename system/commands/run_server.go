package commands

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/multi"
	"github.com/denkhaus/agents/pkg/multi/plugins/web"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/shared"
	sys_shared "github.com/denkhaus/agents/system/shared"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func RunServerCommand(
	ctx *cli.Context,
	configProvider config.ConfigProvider,
	agentFactory config.AgentFactory,
) error {

	return runServerCommand(
		ctx.Context,
		ctx.String("environment"),
		ctx.String("server-addr"),
		configProvider,
		ctx.App.Version,
		agentFactory,
		ctx.App.Name,
	)
}

// runServerCommand starts the multi-agent system
func runServerCommand(
	ctx context.Context,
	environmentName, serverAddr string,
	configProvider config.ConfigProvider,
	appVersion string,
	agentFactory config.AgentFactory,
	appName string,
) error {

	envName := config.EnvironmentName(environmentName)

	logger.Log.Info("Starting agents system",
		zap.String("version", appVersion),
		zap.String("environment", string(envName)),
	)

	// Validate environment exists
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
	agents, err := agentFactory.CreateAllAgentsInEnvironment(ctx, envName)
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
		return startServer(ctx, envConfig, appName, agents, serverAddr)
	})

	return g.Wait()
}

// startServer starts the debug HTTP server
func startServer(
	ctx context.Context,
	envConfig *config.EnvironmentConfig,
	appName string,
	agents []shared.TheAgent,
	addr string,
) error {

	condenserService, err := sys_shared.CreateCondenser(ctx, envConfig)
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
