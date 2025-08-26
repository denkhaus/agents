package prompt

import (
	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/generic"
	"github.com/google/uuid"
)

// promptManagerImpl is an unexported implementation of PromptManager.
type promptManagerImpl struct {
	prompts *generic.ResourceManager[*promptEntry]
}

// NewManager creates a new instance of promptManagerImpl.
func NewManager(prompts map[uuid.UUID]*promptEntry) provider.PromptManager {
	manager := generic.NewResourceManager[*promptEntry]()
	
	// Populate the generic manager with existing prompts
	for id, entry := range prompts {
		manager.Set(id, entry)
	}
	
	return &promptManagerImpl{prompts: manager}
}

func (pm *promptManagerImpl) GetPrompt(agentID uuid.UUID) (provider.Prompt, error) {
	entry, ok := pm.prompts.Get(agentID)
	if !ok {
		return nil, &PromptError{
			Message: "prompt template not found",
			AgentID: agentID,
		}
	}

	return entry, nil
}
