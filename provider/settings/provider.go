package settings

import (
	"fmt"

	"github.com/denkhaus/agents/provider/prompt"
	"github.com/denkhaus/agents/provider/workspace"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type agentSettingsProviderImpl struct {
	workspaceProvider workspace.Provider
	promptProvider    prompt.Provider
	settingsManager   SettingsManager
}

func New(i *do.Injector) (Provider, error) {
	workspaceProvider := do.MustInvoke[workspace.Provider](i)
	promptProvider := do.MustInvoke[prompt.Provider](i)

	settingsManager, err := NewSettingsManager(SettingsFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize settings manager: %w", err)
	}

	return &agentSettingsProviderImpl{
		workspaceProvider: workspaceProvider,
		promptProvider:    promptProvider,
		settingsManager:   settingsManager,
	}, nil
}

func (p *agentSettingsProviderImpl) GetAgentConfiguration(agentID uuid.UUID) (AgentConfiguration, error) {
	workspace, err := p.workspaceProvider.GetWorkspace(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace for agent %s", agentID)
	}

	prompt, err := p.promptProvider.GetPrompt(agentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt for agent %s", agentID)
	}

	settings, err := p.settingsManager.GetSettings(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings for agent %s", agentID)
	}

	return NewConfigurationWithSettings(workspace, prompt, settings)
}
