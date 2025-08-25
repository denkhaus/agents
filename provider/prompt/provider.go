package prompt

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/do"
)

type promptProviderImpl struct {
	manager PromptManager
}

func New(i *do.Injector) (Provider, error) {
	manager, err := NewPromptManager(promptFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize prompt manager: %w", err)
	}

	return &promptProviderImpl{
		manager: manager,
	}, nil
}

func (p *promptProviderImpl) GetPrompt(agentID uuid.UUID, data interface{}) (Prompt, error) {
	return p.manager.GetPrompt(agentID)
}
