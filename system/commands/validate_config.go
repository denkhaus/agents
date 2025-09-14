package commands

import (
	"fmt"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/urfave/cli/v2"
)

// ValidateConfigCommand validates the system configuration
func ValidateConfigCommand(c *cli.Context, configProvider config.ConfigProvider) error {
	logger.Log.Info("Validating system configuration...")

	if err := configProvider.ValidateConfiguration(); err != nil {
		return cli.Exit(fmt.Sprintf("Configuration validation failed: %v", err), 1)
	}

	// Load and validate system config
	systemConfig, err := configProvider.LoadSystemConfig()
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
