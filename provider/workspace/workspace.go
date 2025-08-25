package workspace

import (
	"github.com/denkhaus/agents/provider"
	"github.com/google/uuid"
)

type workspaceImpl struct {
	agentID uuid.UUID
	path    string
}

func NewWorkspace(agentID uuid.UUID, path string) (provider.Workspace, error) {
	return &workspaceImpl{
		agentID: agentID,
		path:    path,
	}, nil
}

func (p *workspaceImpl) GetWorkspacePath() (string, error) {
	return p.path, nil
}
