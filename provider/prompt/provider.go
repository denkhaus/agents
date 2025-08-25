package prompt

import (
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type promptProviderImpl struct {
	manager provider.PromptManager
}

func New(i *do.Injector) (provider.PromptProvider, error) {
	manager, err := NewPromptManager(promptFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize prompt manager: %w", err)
	}

	return &promptProviderImpl{
		manager: manager,
	}, nil
}

func (p *promptProviderImpl) GetPrompt(agentID uuid.UUID, data interface{}) (provider.Prompt, error) {
	return p.manager.GetPrompt(agentID)
}
