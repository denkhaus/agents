package project

import (
	"github.com/denkhaus/agents/pkg/tools/project/repository/inmemory"
	"github.com/denkhaus/agents/pkg/tools/project/shared"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// NewManager creates a new project task manager with default in-memory repository
func NewManager(config *Config) ProjectManager {
	repo := inmemory.NewMemoryRepository()
	svc := newService(repo, config)
	return svc
}

// NewManagerWithRepository creates a new project task manager with custom repository
func NewManagerWithRepository(repo shared.Repository, config *Config) ProjectManager {
	svc := newService(repo, config)
	return svc
}

// NewToolSetProvider creates a new tool set provider
func NewToolSetProvider() ToolSetProvider {
	return &toolSetProvider{}
}

// toolSetProvider implements ToolSetProvider
type toolSetProvider struct{}

func (p *toolSetProvider) CreateToolSet(opts ...Option) (tool.ToolSet, error) {
	return New(opts...)
}
