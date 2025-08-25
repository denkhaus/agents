package settings

import (
	"fmt"

	"github.com/denkhaus/agents/provider/prompt"
	"github.com/denkhaus/agents/provider/workspace"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type Provider interface {
	GetAgentConfiguration(agentID uuid.UUID) (AgentConfiguration, error)
}

type agentSettingsProviderImpl struct {
	workspaceProvider workspace.Provider
	promptProvider    prompt.Provider
}

func New(i *do.Injector) (Provider, error) {
	workspaceProvider := do.MustInvoke[workspace.Provider](i)
	promptProvider := do.MustInvoke[prompt.Provider](i)
	return &agentSettingsProviderImpl{
		workspaceProvider: workspaceProvider,
		promptProvider:    promptProvider,
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

	return NewConfiguration(workspace, prompt)
}
