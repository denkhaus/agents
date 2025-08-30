package main

import (
	"context"
	"fmt"
	"log"

	"github.com/denkhaus/agents/provider/config"
	"github.com/denkhaus/agents/shared"
)

func main() {
	// Create a new agent factory
	// Assuming your CUE configurations are in the ./config directory
	factory := config.NewUnifiedAgentFactory("./config")

	// Validate the configuration
	if err := factory.ValidateConfiguration(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Create an agent by name
	ctx := context.Background()
	agent, err := factory.CreateAgent(ctx, "production", "coder")
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
	config, err := factory.GetAgentConfig("production", "researcher")
	if err != nil {
		log.Fatalf("Failed to get agent config: %v", err)
	}

	fmt.Printf("Researcher agent config: %+v\n", config)
}