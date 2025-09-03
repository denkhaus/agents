package main

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/session/condenser"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
)

func (a *App) createCondenser(_ context.Context, envConfig *config.EnvironmentConfig) (*condenser.Service, error) {

	options := []condenser.ConfigOption{
		condenser.WithTokenCountingMethod(envConfig.Condenser.TokenCountingMethod),
	}

	if envConfig.Condenser.LoggingEnabled {
		options = append(options, condenser.WithLogger(logger.Log))
	}

	if envConfig.Condenser.MaxContextTokens > 0 {
		options = append(
			options, condenser.WithMaxContextTokens(
				envConfig.Condenser.MaxContextTokens,
			),
		)
	}

	if envConfig.Condenser.TriggerThreshold > 0 {
		options = append(
			options, condenser.WithTriggerThreshold(
				envConfig.Condenser.TriggerThreshold,
			),
		)
	}

	if envConfig.Condenser.SummaryPrompt != "" {
		options = append(
			options, condenser.WithSummaryPrompt(
				envConfig.Condenser.SummaryPrompt,
			),
		)
	}

	if envConfig.Condenser.RecentEventsToKeep > 0 {
		options = append(
			options, condenser.WithRecentEventsToKeep(
				envConfig.Condenser.RecentEventsToKeep,
			),
		)
	}

	// TODO: get Model by provider and name from a real system wide ModelProvider
	llm := openai.New(envConfig.Condenser.ModelName)

	// Create the condenser service, wrapping the base session service.
	condenserSessionService, err := condenser.NewWithOptions(
		inmemory.NewSessionService(), llm, options...,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create condenser service: %w", err)
	}

	return condenserSessionService, nil
}
