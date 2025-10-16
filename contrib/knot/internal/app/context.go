package app

import (
	"github.com/denkhaus/knot/internal/manager"
	"go.uber.org/zap"
)

// Context holds the application dependencies
type Context struct {
	ProjectManager manager.ProjectManager
	Logger         *zap.Logger
}

// NewContext creates a new application context with all dependencies
func NewContext(projectManager manager.ProjectManager, logger *zap.Logger) *Context {
	return &Context{
		ProjectManager: projectManager,
		Logger:         logger,
	}
}
