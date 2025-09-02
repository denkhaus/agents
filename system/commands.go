package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/server/debug"
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

// runCommand starts the multi-agent system
func (a *App) runCommand(c *cli.Context) error {
	ctx := c.Context
	envName := config.EnvironmentName(c.String("environment"))

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
			zap.String("id", ag.ID().String()),
		)
	}

	// Start debug server if enabled
	if c.Bool("debug-server") {
		go a.startDebugServer(agents, c.String("debug-addr"))
	}

	// Start chat system
	return a.startChatSystem(ctx, agents, envConfig)
}

// startDebugServer starts the debug HTTP server
func (a *App) startDebugServer(agents []shared.TheAgent, addr string) {
	debugAgents := make(map[string]agent.Agent)
	for _, ag := range agents {
		debugAgents[ag.Info().Name] = ag
	}

	server := debug.New(debugAgents)
	logger.Log.Info("Debug server starting", zap.String("address", addr))

	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		logger.Log.Fatal("Debug server failed", zap.Error(err))
	}
}
