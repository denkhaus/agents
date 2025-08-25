package workspace

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/do"
)

type Provider interface {
	GetWorkspace(agentID uuid.UUID) (Workspace, error)
}

type workspaceProviderImpl struct {
	workspaces map[uuid.UUID]Workspace
}

func New(i *do.Injector) (Provider, error) {
	return &workspaceProviderImpl{
		workspaces: make(map[uuid.UUID]Workspace),
	}, nil
}

func (p *workspaceProviderImpl) GetWorkspace(agentID uuid.UUID) (Workspace, error) {
	if w, ok := p.workspaces[agentID]; ok {
		return w, nil
	}

	// This is hardcoded for now, make this dynamic in the future
	w, err := NewWorkspace(agentID, "/home/denkhaus/dev/gomodules/agents")
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace for agent %s", agentID)
	}

	p.workspaces[agentID] = w
	return w, nil
}
