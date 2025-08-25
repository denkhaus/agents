package prompt

import (
	"github.com/google/uuid"
)

// promptManagerImpl is an unexported implementation of PromptManager.
type promptManagerImpl struct {
	prompts map[uuid.UUID]*promptEntry
}

// NewManager creates a new instance of promptManagerImpl.
func NewManager(prompts map[uuid.UUID]*promptEntry) PromptManager {
	return &promptManagerImpl{prompts: prompts}
}

func (pm *promptManagerImpl) GetPrompt(agentID uuid.UUID) (Prompt, error) {
	entry, ok := pm.prompts[agentID]
	if !ok {
		return nil, &PromptError{
			Message: "prompt template not found",
			AgentID: agentID,
		}
	}

	return entry, nil
}
