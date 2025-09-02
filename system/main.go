package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/multi"
	"github.com/denkhaus/agents/multi/plugins"
	"github.com/denkhaus/agents/multi/plugins/cli"
	"github.com/denkhaus/agents/pkg/provider/config"
	"github.com/denkhaus/agents/pkg/shared"

	"github.com/google/uuid"
	"github.com/samber/do"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/server/debug"
)

const debugServerDefaultListenAddr = ":6999"

func startup(ctx context.Context, selectedEnvironment string) error {

	injector := di.NewContainer()

	agentFactory := do.MustInvoke[config.AgentFactory](injector)

	researcherAgent, err := agentFactory.CreateAgent(ctx, selectedEnvironment, shared.AgentRoleResearcher)
	if err != nil {
		return fmt.Errorf("failed to create researcher agent: %w", err)
	}

	projectManagerAgent, err := agentFactory.CreateAgent(ctx, selectedEnvironment, shared.AgentRoleProjectManager)
	if err != nil {
		return fmt.Errorf("failed to create project manager agent: %w", err)
	}

	coderAgent, err := agentFactory.CreateAgent(ctx, selectedEnvironment, shared.AgentRoleCoder)
	if err != nil {
		return fmt.Errorf("failed to create coder agent: %w", err)
	}

	// // Create session service.
	// // Use the in-memory session service as the base.
	// baseSessionService := inmemory.NewSessionService()

	// // Create the condenser service, wrapping the base session service.
	// condenserSessionService, err := condenser.NewWithOptions(
	// 	baseSessionService,
	// 	modelInstance,                        // Pass the LLM model instance for summarization
	// 	logger.Log,                           // Provide the global logger for visibility
	// 	condenser.WithMaxContextTokens(4000), // Example: max context size of 4000 tokens
	// 	condenser.WithTriggerThreshold(0.75), // Example: trigger condensation at 75% capacity
	// 	condenser.WithSummaryPrompt("Summarize the conversation history for an LLM to continue the conversation without losing context. Keep it concise and focused on key facts."),
	// 	condenser.WithTokenCountingMethod(condenser.TokenCountingAuto), // Auto-select best token counting method
	// 	condenser.WithRecentEventsToKeep(2),                            // Keep 2 most recent events
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to create condenser service: %w", err)
	// }

	// Enhanced Bubble Tea Chat with real LLM calls and spinners
	chat := cli.NewCLIMultiAgentChat(
		plugins.WithProcessorOptions(
			//	multi.WithSessionService(condenserSessionService),
			multi.WithSessionID(uuid.New()),
			multi.WithApplicationName("denkhaus-multi-agent"),
			multi.WithAgents(
				shared.NewHumanAgent(shared.AgentInfoHuman),
				researcherAgent,
				projectManagerAgent,
				coderAgent,
			),
		),
	)

	go func() {
		agents := map[string]agent.Agent{}
		agents[researcherAgent.Info().Name] = researcherAgent
		agents[projectManagerAgent.Info().Name] = projectManagerAgent
		agents[coderAgent.Info().Name] = coderAgent

		server := debug.New(agents)
		if err := http.ListenAndServe(debugServerDefaultListenAddr, server.Handler()); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	return chat.Start(ctx)
}

func main() {
	// TODO: Allow environment selection via command-line argument or environment variable
	selectedEnvironment := "production" // Default environment

	if err := startup(context.Background(), selectedEnvironment); err != nil {
		logger.Log.Fatal("application error", zap.Error(err))
	}
}
