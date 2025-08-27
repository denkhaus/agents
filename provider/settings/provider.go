package settings

import (
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/shared"
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

// GetActiveAgents returns a list of all available agents with their basic information.
// This method collects information about all agents that have been configured in the system.
func (p *agentSettingsProviderImpl) GetActiveAgents() ([]shared.AgentInfo, error) {
	// Get all settings from the settings manager
	allSettings := p.settingsManager.GetAllSettings()

	// Create a slice to hold the agent info
	agents := make([]shared.AgentInfo, 0, len(allSettings))

	// Convert settings to AgentInfo structs
	for agentID, settings := range allSettings {
		if !settings.Agent.Active {
			continue
		}
		agentInfo := shared.NewAgentInfo(
			agentID,
			settings.Agent.Name,
			settings.Agent.Description,
		)

		agents = append(agents, agentInfo)
	}

	return agents, nil
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

	return NewConfiguration(workspace, prompt, p, settings)
}
