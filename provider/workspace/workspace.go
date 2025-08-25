package workspace

import (
	"github.com/google/uuid"
)

type Workspace interface {
	GetWorkspacePath() (string, error)
}

type workspaceImpl struct {
	agentID uuid.UUID
	path    string
}

func NewWorkspace(agentID uuid.UUID, path string) (Workspace, error) {
	return &workspaceImpl{
		agentID: agentID,
		path:    path,
	}, nil
}

func (p *workspaceImpl) GetWorkspacePath() (string, error) {
	return p.path, nil
}
