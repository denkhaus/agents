package settings

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/provider/prompt"
	"github.com/denkhaus/agents/shared"

	shelltoolset "github.com/denkhaus/agents/tools/shell"
	"github.com/denkhaus/agents/utils"
	"github.com/denkhaus/agents/workspace"

	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/file"
)

type agentSettingsImpl struct {
	Settings
	workspace workspace.Workspace
	prompt    prompt.Prompt
}

func NewConfiguration(workspace workspace.Workspace, prompt prompt.Prompt) (AgentConfiguration, error) {
	return &agentSettingsImpl{
		workspace: workspace,
		prompt:    prompt,
	}, nil
}

func NewConfigurationWithSettings(workspace workspace.Workspace, prompt prompt.Prompt, settings *Settings) (AgentConfiguration, error) {
	return &agentSettingsImpl{
		Settings:  *settings,
		workspace: workspace,
		prompt:    prompt,
	}, nil
}

func (p *agentSettingsImpl) GetAgentName() string {
	return p.Agent.Name
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

func (p *agentSettingsImpl) GetModel() (model.Model, error) {
	modelInstance := openai.New(p.Model.Name, openai.WithChannelBufferSize(
		p.Model.ChannelBufferSize,
	))

	return modelInstance, nil
}

func (p *agentSettingsImpl) getToolSets() ([]tool.ToolSet, error) {
	workspacePath, err := p.workspace.GetWorkspacePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspacePath for agent [%s]-[%s]", p.Agent.Role, p.AgentID)
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

func (p *agentSettingsImpl) GetOptions(ctx context.Context) ([]llmagent.Option, error) {
	options := []llmagent.Option{}

	toolSets, err := p.getToolSets()
	if err != nil {
		return nil, fmt.Errorf("failed to get toolsets for [%s]-[%s]", p.Agent.Role, p.AgentID)
	}

	options = append(options, llmagent.WithToolSets(toolSets))

	generationConfig, err := p.getGenerationConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get generation config for agent [%s]-[%s]", p.Agent.Role, p.AgentID)
	}

	toolInfo := utils.GetToolInfoFromSets(ctx, toolSets)

	promptContext := map[string]interface{}{
		shared.ContextKeyToolInfo: toolInfo,
	}

	instruction, err := p.prompt.GetInstruction(promptContext)
	if err != nil {
		return nil, fmt.Errorf("failed to get instruction prompt for agent [%s]-[%s]", p.Agent.Role, p.AgentID)
	}

	options = append(options, llmagent.WithInstruction(instruction))
	options = append(options, llmagent.WithGlobalInstruction(p.prompt.GetGlobalInstruction()))
	options = append(options, llmagent.WithDescription(p.prompt.GetDescription()))

	options = append(options, llmagent.WithGenerationConfig(generationConfig))
	options = append(options, llmagent.WithChannelBufferSize(p.Agent.ChannelBufferSize))

	return options, nil
}
