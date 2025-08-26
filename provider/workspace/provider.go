package workspace

import (
	"fmt"

	"github.com/denkhaus/agents/provider"
	"github.com/denkhaus/agents/provider/generic"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type workspaceProviderImpl struct {
	workspaces *generic.ResourceManager[provider.Workspace]
}

func New(i *do.Injector) (provider.WorkspaceProvider, error) {
	return &workspaceProviderImpl{
		workspaces: generic.NewResourceManager[provider.Workspace](),
	}, nil
}

func (p *workspaceProviderImpl) GetWorkspace(agentID uuid.UUID) (provider.Workspace, error) {
	// Use GetOrSetWithError for thread-safe lazy initialization with proper error handling
	workspace, err := p.workspaces.GetOrSetWithError(agentID, func() (provider.Workspace, error) {
		// This is hardcoded for now, make this dynamic in the future
		w, err := NewWorkspace(agentID, "/home/denkhaus/dev/gomodules/agents")
		if err != nil {
			return nil, fmt.Errorf("failed to create workspace for agent %s: %w", agentID, err)
		}
		return w, nil
	})

	if err != nil {
		return nil, err
	}

	return workspace, nil
}
