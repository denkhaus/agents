package config

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/shared"
	"github.com/denkhaus/agents/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/chainagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/cycleagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/agent/parallelagent"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/planner/react"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// UnifiedAgentFactory provides a unified implementation for creating agents
type UnifiedAgentFactory struct {
	configProvider ConfigProvider
	toolFactory    ToolFactory
}

// NewUnifiedAgentFactory creates a new unified agent factory
func NewUnifiedAgentFactory(configPath string) AgentFactory {
	return &UnifiedAgentFactory{
		configProvider: NewCUEConfigProvider(configPath),
		toolFactory:    NewCUEToolFactory(),
	}
}

// CreateAgent creates an agent using configuration
func (f *UnifiedAgentFactory) CreateAgent(ctx context.Context, environment, agentName string) (shared.TheAgent, error) {
	// Load agent configuration
	agentConfig, err := f.configProvider.LoadAgentComposition(environment, agentName)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent config: %w", err)
	}

	// Create tools based on configuration
	tools, toolsets, err := f.toolFactory.CreateTools(agentConfig.Tools)
	if err != nil {
		return nil, fmt.Errorf("failed to create tools: %w", err)
	}

	// Get all tools (including those from toolsets)
	allTools := f.getAllTools(tools, toolsets)

	// Create the appropriate agent type based on configuration
	var ag agent.Agent
	switch agentConfig.Type {
	case shared.AgentTypeDefault:
		ag, err = f.createLLMAgent(ctx, agentConfig, allTools)
	case shared.AgentTypeChain:
		ag, err = f.createChainAgent(ctx, agentConfig)
	case shared.AgentTypeCycle:
		ag, err = f.createCycleAgent(ctx, agentConfig)
	case shared.AgentTypeParallel:
		ag, err = f.createParallelAgent(ctx, agentConfig)
	default:
		// Default to LLM agent if type is not specified or unknown
		ag, err = f.createLLMAgent(ctx, agentConfig, allTools)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return shared.NewAgent(
		ag,
		agentConfig.AgentID,
		agentConfig.Settings.Agent.StreamingEnabled,
	), nil
}

// CreateAgentByID creates an agent by its UUID using default environment
func (f *UnifiedAgentFactory) CreateAgentByID(ctx context.Context, agentID uuid.UUID) (shared.TheAgent, error) {
	agentName := getAgentNameFromID(agentID)
	if agentName == "unknown" {
		return nil, fmt.Errorf("unknown agent ID: %s", agentID)
	}

	return f.CreateAgent(ctx, "production", agentName)
}

// ValidateConfiguration validates all configurations
func (f *UnifiedAgentFactory) ValidateConfiguration() error {
	return f.configProvider.ValidateConfiguration()
}

// GetAgentConfig returns the raw configuration for an agent
func (f *UnifiedAgentFactory) GetAgentConfig(environment, agentName string) (*AgentConfig, error) {
	return f.configProvider.LoadAgentComposition(environment, agentName)
}

// getAgentNameFromID maps agent UUIDs to their names
func getAgentNameFromID(agentID uuid.UUID) string {
	agentMap := map[uuid.UUID]string{
		shared.AgentIDCoder:          "coder",
		shared.AgentIDProjectManager: "project-manager",
		shared.AgentIDResearcher:     "researcher",
	}

	if name, exists := agentMap[agentID]; exists {
		return name
	}

	return "unknown"
}

// getAllTools combines tools and tools from toolsets into a single slice
func (f *UnifiedAgentFactory) getAllTools(tools []tool.Tool, toolsets []tool.ToolSet) []tool.Tool {
	var allTools []tool.Tool
	
	// Add direct tools
	allTools = append(allTools, tools...)
	
	// Add tools from toolsets
	for _, _ = range toolsets {
		// In a real implementation, you'd get the tools from the toolset
		// This is a simplified version
	}
	
	return allTools
}

// createLLMAgent creates an LLM agent with the provided configuration
func (f *UnifiedAgentFactory) createLLMAgent(ctx context.Context, agentConfig *AgentConfig, tools []tool.Tool) (agent.Agent, error) {
	options := []llmagent.Option{}
	
	// Add generation config
	generationConfig := model.GenerationConfig{
		MaxTokens:   utils.IntPtr(agentConfig.Settings.Agent.LLM.MaxTokens),
		Temperature: utils.FloatPtr(agentConfig.Settings.Agent.LLM.Temperature),
		Stream:      agentConfig.Settings.Agent.StreamingEnabled,
	}
	
	options = append(options, llmagent.WithGenerationConfig(generationConfig))
	
	// Add model
	modelInstance, err := f.getModel(agentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}
	options = append(options, llmagent.WithModel(modelInstance))
	
	// Add instruction
	promptContent := agentConfig.Prompt.Content
	if agentConfig.Prompt.GlobalInstruction != "" {
		promptContent = agentConfig.Prompt.GlobalInstruction + "\n\n" + promptContent
	}
	options = append(options, llmagent.WithInstruction(promptContent))
	
	// Add global instruction
	options = append(options, llmagent.WithGlobalInstruction(agentConfig.Prompt.GlobalInstruction))
	
	// Add planner if enabled
	if agentConfig.Settings.Agent.PlanningEnabled {
		reactPlanner := react.New()
		options = append(options, llmagent.WithPlanner(reactPlanner))
	}
	
	// Add sub-agents if any
	if len(agentConfig.Settings.Agent.SubAgents) > 0 {
		subAgents, err := f.getSubAgents(ctx, agentConfig.Settings.Agent.SubAgents)
		if err != nil {
			return nil, fmt.Errorf("failed to get sub agents: %w", err)
		}
		if len(subAgents) > 0 {
			options = append(options, llmagent.WithSubAgents(subAgents))
		}
	}
	
	// Add schemas and other configurations
	if agentConfig.Settings.Agent.InputSchema != nil {
		options = append(options, llmagent.WithInputSchema(agentConfig.Settings.Agent.InputSchema))
	}
	
	if agentConfig.Settings.Agent.OutputSchema != nil {
		options = append(options, llmagent.WithOutputSchema(agentConfig.Settings.Agent.OutputSchema))
	}
	
	if agentConfig.Settings.Agent.OutputKey != "" {
		options = append(options, llmagent.WithOutputKey(agentConfig.Settings.Agent.OutputKey))
	}
	
	if agentConfig.Settings.Agent.ChannelBufferSize > 0 {
		options = append(options, llmagent.WithChannelBufferSize(agentConfig.Settings.Agent.ChannelBufferSize))
	}
	
	// Add tools
	if len(tools) > 0 {
		options = append(options, llmagent.WithTools(tools))
	}
	
	// Create and return the LLM agent
	return llmagent.New(agentConfig.Name, options...), nil
}

// createChainAgent creates a chain agent with the provided configuration
func (f *UnifiedAgentFactory) createChainAgent(ctx context.Context, agentConfig *AgentConfig) (agent.Agent, error) {
	options := []chainagent.Option{}
	
	// Add sub-agents
	if len(agentConfig.Settings.Agent.SubAgents) > 0 {
		subAgents, err := f.getSubAgents(ctx, agentConfig.Settings.Agent.SubAgents)
		if err != nil {
			return nil, fmt.Errorf("failed to get sub agents: %w", err)
		}
		if len(subAgents) > 0 {
			options = append(options, chainagent.WithSubAgents(subAgents))
		}
	}
	
	// Add channel buffer size if specified
	if agentConfig.Settings.Agent.ChannelBufferSize > 0 {
		options = append(options, chainagent.WithChannelBufferSize(agentConfig.Settings.Agent.ChannelBufferSize))
	}
	
	// Create and return the chain agent
	return chainagent.New(agentConfig.Name, options...), nil
}

// createCycleAgent creates a cycle agent with the provided configuration
func (f *UnifiedAgentFactory) createCycleAgent(ctx context.Context, agentConfig *AgentConfig) (agent.Agent, error) {
	options := []cycleagent.Option{}
	
	// Add sub-agents
	if len(agentConfig.Settings.Agent.SubAgents) > 0 {
		subAgents, err := f.getSubAgents(ctx, agentConfig.Settings.Agent.SubAgents)
		if err != nil {
			return nil, fmt.Errorf("failed to get sub agents: %w", err)
		}
		if len(subAgents) > 0 {
			options = append(options, cycleagent.WithSubAgents(subAgents))
		}
	}
	
	// Add max iterations if specified
	if agentConfig.Settings.Agent.MaxIterations > 0 {
		options = append(options, cycleagent.WithMaxIterations(agentConfig.Settings.Agent.MaxIterations))
	}
	
	// Add channel buffer size if specified
	if agentConfig.Settings.Agent.ChannelBufferSize > 0 {
		options = append(options, cycleagent.WithChannelBufferSize(agentConfig.Settings.Agent.ChannelBufferSize))
	}
	
	// Create and return the cycle agent
	return cycleagent.New(agentConfig.Name, options...), nil
}

// createParallelAgent creates a parallel agent with the provided configuration
func (f *UnifiedAgentFactory) createParallelAgent(ctx context.Context, agentConfig *AgentConfig) (agent.Agent, error) {
	options := []parallelagent.Option{}
	
	// Add sub-agents
	if len(agentConfig.Settings.Agent.SubAgents) > 0 {
		subAgents, err := f.getSubAgents(ctx, agentConfig.Settings.Agent.SubAgents)
		if err != nil {
			return nil, fmt.Errorf("failed to get sub agents: %w", err)
		}
		if len(subAgents) > 0 {
			options = append(options, parallelagent.WithSubAgents(subAgents))
		}
	}
	
	// Add channel buffer size if specified
	if agentConfig.Settings.Agent.ChannelBufferSize > 0 {
		options = append(options, parallelagent.WithChannelBufferSize(agentConfig.Settings.Agent.ChannelBufferSize))
	}
	
	// Create and return the parallel agent
	return parallelagent.New(agentConfig.Name, options...), nil
}

// getModel creates a model instance based on the configuration
func (f *UnifiedAgentFactory) getModel(agentConfig *AgentConfig) (model.Model, error) {
	switch agentConfig.Settings.Agent.LLM.Provider {
	case shared.ModelProviderOpenAI:
		modelOptions := []openai.Option{}

		if len(agentConfig.Settings.Agent.LLM.BaseURL) > 0 {
			modelOptions = append(modelOptions,
				openai.WithBaseURL(
					agentConfig.Settings.Agent.LLM.BaseURL,
				),
			)
		}

		if len(agentConfig.Settings.Agent.LLM.APIKey) > 0 {
			modelOptions = append(modelOptions,
				openai.WithAPIKey(
					agentConfig.Settings.Agent.LLM.APIKey,
				),
			)
		}

		if agentConfig.Settings.Agent.LLM.ChannelBufferSize > 0 {
			modelOptions = append(modelOptions,
				openai.WithChannelBufferSize(
					agentConfig.Settings.Agent.LLM.ChannelBufferSize,
				),
			)
		}

		modelInstance := openai.New(agentConfig.Settings.Agent.LLM.Model, modelOptions...)
		return modelInstance, nil
	}

	return nil, fmt.Errorf("model provider %s is unknown", agentConfig.Settings.Agent.LLM.Provider)
}

// getSubAgents creates sub-agent instances based on their UUIDs
// Note: In a complete implementation, this would recursively create sub-agents
func (f *UnifiedAgentFactory) getSubAgents(ctx context.Context, subAgentIDs []uuid.UUID) ([]agent.Agent, error) {
	// This is a simplified implementation that returns empty agents
	// In a real implementation, you would create the actual sub-agents
	var subAgents []agent.Agent
	for _, id := range subAgentIDs {
		// Create a simple placeholder agent for now
		subAgent := llmagent.New(fmt.Sprintf("sub-agent-%s", id.String()))
		subAgents = append(subAgents, subAgent)
	}
	
	// Log that this is a simplified implementation
	zap.L().Warn("getSubAgents: using simplified implementation - sub-agents are placeholders")
	
	return subAgents, nil
}
