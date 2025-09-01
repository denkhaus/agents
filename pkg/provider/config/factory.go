package config

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/pkg/shared"
	"github.com/denkhaus/agents/pkg/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
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
func NewUnifiedAgentFactory(injector *do.Injector) (AgentFactory, error) {
	toolFactory := do.MustInvoke[ToolFactory](injector)
	configProvider := do.MustInvoke[ConfigProvider](injector)

	f := &UnifiedAgentFactory{
		configProvider: configProvider,
		toolFactory:    toolFactory,
	}

	return f, nil
}

// CreateAgent creates an agent using configuration
func (f *UnifiedAgentFactory) CreateAgent(ctx context.Context, environment string, agentRole shared.AgentRole) (shared.TheAgent, error) {
	// Load agent configuration
	agentConfig, err := f.configProvider.LoadAgentComposition(environment, agentRole)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent config: %w", err)
	}

	// Create tools based on configuration
	tools, toolsets, err := f.toolFactory.CreateTools(agentConfig.Tool)
	if err != nil {
		return nil, fmt.Errorf("failed to create tools: %w", err)
	}

	agentInfo, err := f.configProvider.GetAgentsInEnvironment(environment, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent info from environment %q: %w", environment, err)
	}

	// Get all tools (including those from toolsets)
	allTools := f.getAllTools(ctx, tools, toolsets)

	// Create the appropriate agent type based on configuration
	var ag agent.Agent
	switch agentConfig.Type {
	case shared.AgentTypeDefault:
		ag, err = f.createLLMAgent(ctx, environment, agentConfig, allTools, agentInfo)
	case shared.AgentTypeChain:
		ag, err = f.createChainAgent(ctx, environment, agentConfig, allTools)
	case shared.AgentTypeCycle:
		ag, err = f.createCycleAgent(ctx, environment, agentConfig, allTools)
	case shared.AgentTypeParallel:
		ag, err = f.createParallelAgent(ctx, environment, agentConfig, allTools)
	default:
		// Default to LLM agent if type is not specified or unknown
		ag, err = f.createLLMAgent(ctx, environment, agentConfig, allTools, agentInfo)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return shared.NewAgent(
		ag,
		agentConfig.AgentID,
		agentConfig.Setting.Agent.StreamingEnabled,
	), nil
}

// CreateAgentByID creates an agent by its UUID using default environment
func (f *UnifiedAgentFactory) CreateAgentByID(ctx context.Context, environment string, agentID uuid.UUID) (shared.TheAgent, error) {
	agentInfo, err := f.configProvider.GetAgentInfoByID(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent info for ID %s: %w", agentID, err)
	}

	return f.CreateAgent(ctx, environment, agentInfo.Role())
}

// ValidateConfiguration validates all configurations
func (f *UnifiedAgentFactory) ValidateConfiguration() error {
	return f.configProvider.ValidateConfiguration()
}

// GetAgentConfig returns the raw configuration for an agent
func (f *UnifiedAgentFactory) GetAgentConfig(environment string, agentRole shared.AgentRole) (*AgentConfig, error) {
	return f.configProvider.LoadAgentComposition(environment, agentRole)
}

// GetAgentNameFromID maps agent UUIDs to their names
// This function now uses the ConfigProvider to get the agent's name.
func (f *UnifiedAgentFactory) GetAgentNameFromID(agentID uuid.UUID) string {
	agentInfo, err := f.configProvider.GetAgentInfoByID(agentID)
	if err != nil {
		logger.Log.Warn("Failed to get agent info by ID", zap.Any("agent_id", agentID), zap.Error(err))
		return "unknown"
	}
	return agentInfo.Name
}

// getAllTools combines tools and tools from toolsets into a single slice
func (f *UnifiedAgentFactory) getAllTools(ctx context.Context, tools []tool.Tool, toolsets []tool.ToolSet) []tool.Tool {
	var allTools []tool.Tool

	// Add direct tools
	allTools = append(allTools, tools...)

	// Add tools from toolsets
	for _, toolset := range toolsets {
		callableTools := toolset.Tools(ctx)
		for _, t := range callableTools {
			allTools = append(allTools, t)
		}
	}

	return allTools
}

// createLLMAgent creates an LLM agent with the provided configuration
func (f *UnifiedAgentFactory) createLLMAgent(
	ctx context.Context,
	environment string,
	agentConfig *AgentConfig,
	tools []tool.Tool,
	agentInfo []*shared.AgentInfo,
) (agent.Agent, error) {
	options := []llmagent.Option{}

	// Add generation config
	generationConfig := model.GenerationConfig{
		MaxTokens:   utils.IntPtr(agentConfig.Setting.Agent.LLM.MaxTokens),
		Temperature: utils.FloatPtr(agentConfig.Setting.Agent.LLM.Temperature),
		Stream:      agentConfig.Setting.Agent.StreamingEnabled,
	}

	options = append(options, llmagent.WithGenerationConfig(generationConfig))

	// Add model
	modelInstance, err := f.getModel(agentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}
	options = append(options, llmagent.WithModel(modelInstance))

	// Prepare prompt context for rendering
	promptContext := map[string]interface{}{
		shared.ContextKeyToolInfo:  utils.GetToolInfo(tools...),
		shared.ContextKeyAgentInfo: agentInfo,
	}

	// Validate prompt context against schema
	if err := utils.ValidateSchema(promptContext, agentConfig.Prompt.Schema); err != nil {
		return nil, fmt.Errorf("prompt context validation failed for agent %s: %w", agentConfig.Name, err)
	}

	// Render prompt content using text/template
	renderedContent, err := f.renderPromptContent(agentConfig.Prompt.Content, promptContext)
	if err != nil {
		return nil, fmt.Errorf("failed to render prompt content for agent %s: %w", agentConfig.Name, err)
	}

	// Add instruction with rendered content
	options = append(options, llmagent.WithInstruction(renderedContent))

	// Add global instruction
	options = append(options, llmagent.WithGlobalInstruction(agentConfig.Prompt.GlobalInstruction))

	// Add planner if enabled
	if agentConfig.Setting.Agent.PlanningEnabled {
		reactPlanner := react.New()
		options = append(options, llmagent.WithPlanner(reactPlanner))
	}

	// Add sub-agents if any
	if len(agentConfig.Setting.Agent.SubAgents) > 0 {
		subAgents, err := f.getSubAgents(ctx, environment, agentConfig.Setting.Agent.SubAgents)
		if err != nil {
			return nil, fmt.Errorf("failed to get sub agents: %w", err)
		}
		if len(subAgents) > 0 {
			options = append(options, llmagent.WithSubAgents(subAgents))
		}
	}

	// Add schemas and other configurations
	if agentConfig.Setting.Agent.InputSchema != nil {
		options = append(options, llmagent.WithInputSchema(agentConfig.Setting.Agent.InputSchema))
	}

	if agentConfig.Setting.Agent.OutputSchema != nil {
		options = append(options, llmagent.WithOutputSchema(agentConfig.Setting.Agent.OutputSchema))
	}

	if agentConfig.Setting.Agent.OutputKey != "" {
		options = append(options, llmagent.WithOutputKey(agentConfig.Setting.Agent.OutputKey))
	}

	if agentConfig.Setting.Agent.ChannelBufferSize > 0 {
		options = append(options, llmagent.WithChannelBufferSize(agentConfig.Setting.Agent.ChannelBufferSize))
	}

	// Add tools
	if len(tools) > 0 {
		options = append(options, llmagent.WithTools(tools))
	}

	// Create and return the LLM agent
	return llmagent.New(agentConfig.Name, options...), nil
}

// getModel creates a model instance based on the configuration
func (f *UnifiedAgentFactory) getModel(agentConfig *AgentConfig) (model.Model, error) {
	switch agentConfig.Setting.Agent.LLM.Provider {
	case shared.ModelProviderOpenAI:
		modelOptions := []openai.Option{}

		if len(agentConfig.Setting.Agent.LLM.BaseURL) > 0 {
			modelOptions = append(modelOptions,
				openai.WithBaseURL(
					agentConfig.Setting.Agent.LLM.BaseURL,
				),
			)
		}

		if len(agentConfig.Setting.Agent.LLM.APIKey) > 0 {
			modelOptions = append(modelOptions,
				openai.WithAPIKey(
					agentConfig.Setting.Agent.LLM.APIKey,
				),
			)
		}

		if agentConfig.Setting.Agent.LLM.ChannelBufferSize > 0 {
			modelOptions = append(modelOptions,
				openai.WithChannelBufferSize(
					agentConfig.Setting.Agent.LLM.ChannelBufferSize,
				),
			)
		}

		modelInstance := openai.New(agentConfig.Setting.Agent.LLM.Model, modelOptions...)
		return modelInstance, nil
	}

	return nil, fmt.Errorf("model provider %s is unknown", agentConfig.Setting.Agent.LLM.Provider)
}
