package settings

import (
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type agentSettingsProviderImpl struct {
	workspaceProvider provider.WorkspaceProvider
	promptProvider    provider.PromptProvider
	settingsManager   SettingsManager
}

func New(i *do.Injector) (provider.SettingsProvider, error) {
	workspaceProvider := do.MustInvoke[provider.WorkspaceProvider](i)
	promptProvider := do.MustInvoke[provider.PromptProvider](i)

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

func (p *agentSettingsProviderImpl) GetAgentConfiguration(agentID uuid.UUID) (provider.AgentConfiguration, error) {
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

	return NewConfiguration(workspace, prompt, settings)
}
