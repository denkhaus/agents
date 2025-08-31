package agents

import (
	"context"
	"fmt"

	providerConfig "github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/shared"
	"github.com/samber/do"
)

func CreateCoderAgent(ctx context.Context, injector *do.Injector) (shared.TheAgent, error) {
	agentFactory := do.MustInvoke[providerConfig.AgentFactory](injector)

	// Create agent using the unified factory - it will handle tool loading automatically
	coderAgent, err := agentFactory.CreateAgent(ctx, "production", "coder")
	if err != nil {
		return nil, fmt.Errorf("failed to create coder agent: %w", err)
	}

	return coderAgent, nil
}
