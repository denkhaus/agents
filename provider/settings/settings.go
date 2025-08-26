package settings

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/shared"

	shelltoolset "github.com/denkhaus/agents/tools/shell"
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
	"trpc.group/trpc-go/trpc-agent-go/tool/file"
)

type agentSettingsImpl struct {
	*Settings
	workspace provider.Workspace
	prompt    provider.Prompt
}

func NewConfiguration(
	workspace provider.Workspace,
	prompt provider.Prompt,
	settings *Settings,
) (provider.AgentConfiguration, error) {
	return &agentSettingsImpl{
		Settings:  settings,
		workspace: workspace,
		prompt:    prompt,
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
	modelInstance := openai.New(p.Model.Name, openai.WithChannelBufferSize(
		p.Model.ChannelBufferSize,
	))

	return modelInstance, nil
}

func (p *agentSettingsImpl) getSubAgents(
	ctx context.Context,
	provider provider.AgentProvider,
) ([]agent.Agent, error) {

	var subAgents []agent.Agent
	for _, agentID := range p.Agent.SubAgents {
		agent, _, err := provider.GetAgent(ctx, agentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get agent with id %s", agentID)
		}

		subAgents = append(subAgents, agent)
	}

	return subAgents, nil
}

func (p *agentSettingsImpl) getToolSets() ([]tool.ToolSet, error) {
	workspacePath, err := p.workspace.GetWorkspacePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspacePath for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}
	// Create file operation tools.
	fileToolSet, err := file.NewToolSet(
		file.WithBaseDir(workspacePath),
	)
	if err != nil {
		return nil, fmt.Errorf("create file tool set: %w", err)
	}

	shellToolSet, err := shelltoolset.NewToolSet(
		shelltoolset.WithBaseDir(workspacePath),
	)
	if err != nil {
		return nil, fmt.Errorf("create shell tool set: %w", err)
	}

	return []tool.ToolSet{fileToolSet, shellToolSet}, nil
}

func (p *agentSettingsImpl) GetDefaultOptions(
	ctx context.Context,
	agentProvider provider.AgentProvider,
	opt ...llmagent.Option,
) ([]llmagent.Option, error) {

	options := []llmagent.Option{}
	options = append(options, opt...)

	toolSets, err := p.getToolSets()
	if err != nil {
		return nil, fmt.Errorf("failed to get toolsets for [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	options = append(options, llmagent.WithToolSets(toolSets))

	generationConfig, err := p.getGenerationConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get generation config for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	options = append(options, llmagent.WithGenerationConfig(generationConfig))

	toolInfo := utils.GetToolInfoFromSets(ctx, toolSets)

	promptContext := map[string]interface{}{
		shared.ContextKeyToolInfo: toolInfo,
	}

	instruction, err := p.prompt.GetInstruction(promptContext)
	if err != nil {
		return nil, fmt.Errorf("failed to get instruction prompt for agent [%s]-[%s]: %w", p.Agent.Role, p.AgentID, err)
	}

	options = append(options, llmagent.WithInstruction(instruction))

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
	options = append(options, llmagent.WithGlobalInstruction(p.prompt.GetGlobalInstruction()))
	options = append(options, llmagent.WithDescription(p.prompt.GetDescription()))
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
