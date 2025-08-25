package workspace

import (
	"github.com/google/uuid"
)

type Workspace interface {
	GetWorkspacePath() (string, error)
}

type workspaceImpl struct {
	agentID uuid.UUID
}

func New(agentID uuid.UUID) (Workspace, error) {
	return &workspaceImpl{
		agentID: agentID,
	}, nil
}

func (p *workspaceImpl) GetWorkspacePath() (string, error) {
	return "/home/denkhaus/dev/gomodules/agents", nil
}
