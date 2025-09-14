package commands

import (
	"fmt"
	"sort"

	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/urfave/cli/v2"
)

// listEnvironmentsCommand shows available environments
func ListEnvironmentsCommand(c *cli.Context, configProvider config.ConfigProvider) error {
	systemConfig, err := configProvider.LoadSystemConfig()
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
