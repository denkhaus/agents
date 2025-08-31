package main

import (
	"context"
	"fmt"
	"log"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/shared"
)

func main() {
	// Create a new agent factory
	// Assuming your CUE configurations are in the ./config directory

	container := di.NewContainer()

	factory, err := config.NewUnifiedAgentFactory(container)
	if err != nil {
		log.Fatalf("Failed to create agent factory: %v", err)
	}

	// Validate the configuration
	if err := factory.ValidateConfiguration(); err != nil {
		log.Printf("Configuration validation warning: %v", err)
		// Not failing here as we want to demonstrate usage even with warnings
	}

	// Create an agent by name
	ctx := context.Background()
	agent, err := factory.CreateAgent(ctx, "development", "coder")
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	fmt.Printf("Created agent: %s (ID: %s)\n", agent.Info().Name, agent.ID())

	// Create an agent by ID
	agentByID, err := factory.CreateAgentByID(ctx, shared.AgentIDProjectManager)
	if err != nil {
		log.Fatalf("Failed to create agent by ID: %v", err)
	}

	fmt.Printf("Created agent by ID: %s (ID: %s)\n", agentByID.Info().Name, agentByID.ID())

	// Get raw configuration
	config, err := factory.GetAgentConfig("development", "researcher")
	if err != nil {
		log.Fatalf("Failed to get agent config: %v", err)
	}

	fmt.Printf("Researcher agent config name: %s\n", config.Name)
	fmt.Printf("Researcher agent config type: %s\n", config.Type)
}
