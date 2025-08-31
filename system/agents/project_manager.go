package agents

import (
	"context"
	"fmt"

	providerConfig "github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
)

func CreateProjectManagerAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentFactory := do.MustInvoke[providerConfig.AgentFactory](injector)

	// Create agent using the unified factory - it will handle tool loading automatically
	projectManagerAgent, err := agentFactory.CreateAgent(ctx, "production", "project_manager")
	if err != nil {
		return nil, fmt.Errorf("failed to create project manager agent: %w", err)
	}

	return projectManagerAgent, nil
}
