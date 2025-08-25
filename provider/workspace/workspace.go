package workspace

import (
	"errors"

	"github.com/denkhaus/agents/workspace"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type Provider interface {
	GetWorkspace(agentID uuid.UUID) (workspace.Workspace, error)
}

type workspaceProviderImpl struct {
}

func New(i *do.Injector) (Provider, error) {
	return &workspaceProviderImpl{}, nil
}

func (p *workspaceProviderImpl) GetWorkspace(agentID uuid.UUID) (workspace.Workspace, error) {
	return nil, errors.New("not yet implemented")
}
