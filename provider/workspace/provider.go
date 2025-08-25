package workspace

import (
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type workspaceProviderImpl struct {
	workspaces map[uuid.UUID]provider.Workspace
}

func New(i *do.Injector) (provider.WorkspaceProvider, error) {
	return &workspaceProviderImpl{
		workspaces: make(map[uuid.UUID]provider.Workspace),
	}, nil
}

func (p *workspaceProviderImpl) GetWorkspace(agentID uuid.UUID) (provider.Workspace, error) {
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
