package prompt

import (
	"github.com/denkhaus/agents/provider"
	"github.com/google/uuid"
)

// promptManagerImpl is an unexported implementation of PromptManager.
type promptManagerImpl struct {
	prompts map[uuid.UUID]*promptEntry
}

// NewManager creates a new instance of promptManagerImpl.
func NewManager(prompts map[uuid.UUID]*promptEntry) provider.PromptManager {
	return &promptManagerImpl{prompts: prompts}
}

func (pm *promptManagerImpl) GetPrompt(agentID uuid.UUID) (provider.Prompt, error) {
	entry, ok := pm.prompts[agentID]
	if !ok {
		return nil, &PromptError{
			Message: "prompt template not found",
			AgentID: agentID,
		}
	}

	return entry, nil
}
