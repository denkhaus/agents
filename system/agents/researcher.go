package agents

import (
	"context"
	"fmt"

	providerConfig "github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
)

func CreateResearcherAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentFactory := do.MustInvoke[providerConfig.AgentFactory](injector)

	// Create agent using the unified factory - it will handle tool loading automatically
	researcherAgent, err := agentFactory.CreateAgent(ctx, "production", "researcher")
	if err != nil {
		return nil, fmt.Errorf("failed to create researcher agent: %w", err)
	}

	return researcherAgent, nil
}
