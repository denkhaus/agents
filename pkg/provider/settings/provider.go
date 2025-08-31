package settings

import (
	"fmt"

	"github.com/denkhaus/agents/pkg/config"
	"github.com/denkhaus/agents/pkg/provider"
	"github.com/denkhaus/agents/shared"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type agentSettingsProviderImpl struct {
	appConfigService config.Service
	promptProvider   provider.PromptProvider
	settingsManager  SettingsManager
}

func NewWithDI(i *do.Injector) (provider.SettingsProvider, error) {
	appConfigService := do.MustInvoke[config.Service](i)
	promptProvider := do.MustInvoke[provider.PromptProvider](i)

	settingsManager, err := NewSettingsManager(SettingsFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize settings manager: %w", err)
	}

	return &agentSettingsProviderImpl{
		appConfigService: appConfigService,
		promptProvider:   promptProvider,
		settingsManager:  settingsManager,
	}, nil
}

// GetActiveAgents returns a list of all available agents with their basic information.
// This method collects information about all agents that have been configured in the system.
func (p *agentSettingsProviderImpl) GetActiveAgents(includeHumanAgent bool) ([]shared.AgentInfo, error) {
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
			settings.Agent.Role,
			settings.Agent.StreamingEnabled,
			settings.Agent.Name,
			settings.Agent.Description,
		)

		agentInfo.InputSchema = settings.Agent.InputSchema
		agentInfo.OutputSchema = settings.Agent.OutputSchema
		agents = append(agents, agentInfo)
	}

	if includeHumanAgent {
		agents = append(agents, shared.AgentInfoHuman)
	}

	return agents, nil
}

func (p *agentSettingsProviderImpl) GetAgentConfiguration(agentID uuid.UUID) (provider.AgentConfiguration, error) {

	prompt, err := p.promptProvider.GetPrompt(agentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt for agent %s", agentID)
	}

	settings, err := p.settingsManager.GetSettings(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings for agent %s", agentID)
	}

	return NewConfiguration(prompt, p, settings)
}
