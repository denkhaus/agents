package settings

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/logger"
	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/shared"
	"go.uber.org/zap"

	"github.com/denkhaus/agents/utils"
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

type agentSettingsImpl struct {
	*Settings
	workspaceProvider provider.Workspace
	promptProvider    provider.Prompt
	settingsProvider  provider.SettingsProvider
}

func NewConfiguration(
	workspaceProvider provider.Workspace,
	promptProvider provider.Prompt,
	settingsProvider provider.SettingsProvider,
	settings *Settings,
) (provider.AgentConfiguration, error) {
	return &agentSettingsImpl{
		Settings:          settings,
		settingsProvider:  settingsProvider,
		workspaceProvider: workspaceProvider,
		promptProvider:    promptProvider,
	}, nil
}

func (p *agentSettingsImpl) GetName() string {
	return p.Agent.Name
}

func (p *agentSettingsImpl) GetType() shared.AgentType {
	return p.Agent.Type
}

func (p *agentSettingsImpl) IsStreamingEnabled() bool {
	return p.Agent.StreamingEnabled
}

func (p *agentSettingsImpl) getGenerationConfig() (model.GenerationConfig, error) {
	return model.GenerationConfig{
		MaxTokens:   utils.IntPtr(p.Agent.MaxTokens),
		Temperature: utils.FloatPtr(p.Agent.Temperature),
		Stream:      p.Agent.StreamingEnabled,
	}, nil
}

func (p *agentSettingsImpl) getModel() (model.Model, error) {
	switch p.Model.Provider {
	case shared.ModelProviderOpenAI:
		modelOptions := []openai.Option{}

		if len(p.Model.BaseURL) > 0 {
			modelOptions = append(modelOptions,
				openai.WithBaseURL(
					p.Model.BaseURL,
				),
			)
		}

		if len(p.Model.APIKey) > 0 {
			modelOptions = append(modelOptions,
				openai.WithAPIKey(
					p.Model.APIKey,
				),
			)
		}

		if p.Model.ChannelBufferSize > 0 {
			modelOptions = append(modelOptions,
				openai.WithChannelBufferSize(
					p.Model.ChannelBufferSize,
				),
			)
		}

		modelInstance := openai.New(p.Model.Name, modelOptions...)
		return modelInstance, nil
	}

	return nil, fmt.Errorf("model provider %s is unknown", p.Model.Provider)
}

func (p *agentSettingsImpl) getSubAgents(
	ctx context.Context,
	provider provider.AgentProvider,
) ([]agent.Agent, error) {

	var subAgents []agent.Agent
	for _, agentID := range p.Agent.SubAgents {
		agent, err := provider.GetAgent(ctx, agentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get agent with id %s", agentID)
		}

		subAgents = append(subAgents, agent)
	}

	return subAgents, nil
}

// getToolsFromOptions processes llmagent.Option, extracts tools, and deduplicates them.
func (p *agentSettingsImpl) getToolsFromOptions(ctx context.Context, options ...llmagent.Option) ([]tool.Tool, error) {
	var llmOptions llmagent.Options

	// Helper function to add a tool to the map, handling duplicates.
	addToolToMap := func(toolsMap map[string]tool.Tool, t tool.Tool) {
		name := t.Declaration().Name
		if _, exists := toolsMap[name]; exists {
			logger.Log.Warn("getToolsFromOptions: tool already registered, skipping", zap.String("tool_name", name))
		} else {
			toolsMap[name] = t
		}
	}

	for _, opt := range options {
		opt(&llmOptions)
	}

	// toolsMap is used to deduplicate tools by their name.
	toolsMap := make(map[string]tool.Tool)

	for _, t := range llmOptions.Tools {
		addToolToMap(toolsMap, t)
	}

	for _, toolSet := range llmOptions.ToolSets {
		for _, t := range toolSet.Tools(ctx) {
			addToolToMap(toolsMap, t)
		}
	}

	// Convert the map of unique tools back to a slice.
	// Pre-allocate slice capacity for minor performance optimization.
	allTools := make([]tool.Tool, 0, len(toolsMap))
	for _, t := range toolsMap {
		allTools = append(allTools, t)
	}

	return allTools, nil
}

func (p *agentSettingsImpl) GetDefaultOptions(
	ctx context.Context,
	agentProvider provider.AgentProvider,
	opt ...llmagent.Option,
) ([]llmagent.Option, error) {

	options := []llmagent.Option{}
	options = append(options, opt...)

	tools, err := p.getToolsFromOptions(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to get toolsets for [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	generationConfig, err := p.getGenerationConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get generation config for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	options = append(options, llmagent.WithGenerationConfig(generationConfig))

	//TODO: make the includeHumanAgent setting configurable
	availableAgentsVal, err := p.settingsProvider.GetActiveAgents(true) // Renamed to avoid conflict
	if err != nil {
		return nil, fmt.Errorf("failed to get active agents: %w", err)
	}

	// Convert []shared.AgentInfo to []*shared.AgentInfo
	availableAgentsPtr := make([]*shared.AgentInfo, len(availableAgentsVal))
	for i := range availableAgentsVal {
		availableAgentsPtr[i] = &availableAgentsVal[i]
	}

	// prepare context for prompt renderring
	promptContext := map[string]interface{}{
		shared.ContextKeyToolInfo:  utils.GetToolInfo(tools...),
		shared.ContextKeyAgentInfo: utils.GetAgentInfoForAgent(p.AgentID, availableAgentsPtr...), // Pass pointers
	}

	instruction, err := p.promptProvider.GetInstruction(promptContext)
	if err != nil {
		return nil, fmt.Errorf("failed to get instruction prompt for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	options = append(options, llmagent.WithInstruction(instruction))
	options = append(options, llmagent.WithGlobalInstruction(p.promptProvider.GetGlobalInstruction()))

	model, err := p.getModel()
	if err != nil {
		return nil, fmt.Errorf("failed to get model for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	options = append(options, llmagent.WithModel(model))

	if p.Agent.PlanningEnabled {
		reactPlanner := react.New()
		options = append(options, llmagent.WithPlanner(reactPlanner))
	}

	subAgents, err := p.getSubAgents(ctx, agentProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to get subagents for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	if len(subAgents) > 0 {
		options = append(options, llmagent.WithSubAgents(subAgents))
	}

	options = append(options, llmagent.WithInputSchema(p.Agent.InputSchema))
	options = append(options, llmagent.WithOutputSchema(p.Agent.OutputSchema))
	options = append(options, llmagent.WithOutputKey(p.Agent.OutputKey))
	options = append(options, llmagent.WithDescription(p.Agent.Description))
	options = append(options, llmagent.WithChannelBufferSize(p.Agent.ChannelBufferSize))

	return options, nil
}

func (p *agentSettingsImpl) GetCycleOptions(
	ctx context.Context,
	agentProvider provider.AgentProvider,
	opt ...cycleagent.Option,
) ([]cycleagent.Option, error) {

	options := []cycleagent.Option{}
	options = append(options, opt...)

	subAgents, err := p.getSubAgents(ctx, agentProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to get subagents for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	if len(subAgents) > 0 {
		options = append(options, cycleagent.WithSubAgents(subAgents))
	}

	options = append(options, cycleagent.WithMaxIterations(p.Agent.MaxIterations))
	options = append(options, cycleagent.WithChannelBufferSize(p.Agent.ChannelBufferSize))

	return options, nil
}

func (p *agentSettingsImpl) GetChainOptions(
	ctx context.Context,
	agentProvider provider.AgentProvider,
	opt ...chainagent.Option,
) ([]chainagent.Option, error) {

	options := []chainagent.Option{}
	options = append(options, opt...)

	subAgents, err := p.getSubAgents(ctx, agentProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to get subagents for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	if len(subAgents) > 0 {
		options = append(options, chainagent.WithSubAgents(subAgents))
	}

	options = append(options, chainagent.WithChannelBufferSize(p.Agent.ChannelBufferSize))

	return options, nil
}

func (p *agentSettingsImpl) GetParallelOptions(
	ctx context.Context,
	agentProvider provider.AgentProvider,
	opt ...parallelagent.Option,
) ([]parallelagent.Option, error) {

	options := []parallelagent.Option{}
	options = append(options, opt...)

	subAgents, err := p.getSubAgents(ctx, agentProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to get subagents for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	if len(subAgents) > 0 {
		options = append(options, parallelagent.WithSubAgents(subAgents))
	}

	options = append(options, parallelagent.WithChannelBufferSize(p.Agent.ChannelBufferSize))

	return options, nil
}
