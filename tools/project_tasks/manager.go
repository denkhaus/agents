package projecttasks

import (
	"context"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// manager implements ProjectTaskManager using resource manager pattern
type manager struct {
	service *service
}

// NewManager creates a new project task manager
func NewManager(config *Config) ProjectTaskManager {
	repo := newMemoryRepository()
	svc := newService(repo, config)
	
	return &manager{
		service: svc,
	}
}

// NewToolSetProvider creates a new tool set provider
func NewToolSetProvider() ToolSetProvider {
	return &toolSetProvider{}
}

// toolSetProvider implements ToolSetProvider
type toolSetProvider struct{}

func (p *toolSetProvider) CreateToolSet(opts ...Option) (tool.ToolSet, error) {
	return NewToolSet(opts...)
}

// Project operations

func (m *manager) CreateProject(ctx context.Context, title, description string) (*Project, error) {
	return m.service.CreateProject(ctx, title, description)
}

func (m *manager) GetProject(ctx context.Context, projectID uuid.UUID) (*Project, error) {
	return m.service.GetProject(ctx, projectID)
}

func (m *manager) UpdateProject(ctx context.Context, projectID uuid.UUID, title, description string) (*Project, error) {
	return m.service.UpdateProject(ctx, projectID, title, description)
}

func (m *manager) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	return m.service.DeleteProject(ctx, projectID)
}

func (m *manager) ListProjects(ctx context.Context) ([]*Project, error) {
	return m.service.ListProjects(ctx)
}

// Task operations

func (m *manager) CreateTask(ctx context.Context, projectID uuid.UUID, parentID *uuid.UUID, title, description string, complexity, priority int) (*Task, error) {
	return m.service.CreateTask(ctx, projectID, parentID, title, description, complexity, priority)
}

func (m *manager) GetTask(ctx context.Context, taskID uuid.UUID) (*Task, error) {
	return m.service.GetTask(ctx, taskID)
}

func (m *manager) UpdateTask(ctx context.Context, taskID uuid.UUID, title, description string, complexity, priority int, state TaskState) (*Task, error) {
	return m.service.UpdateTask(ctx, taskID, title, description, complexity, priority, state)
}

func (m *manager) UpdateTaskState(ctx context.Context, taskID uuid.UUID, state TaskState) (*Task, error) {
	return m.service.UpdateTaskState(ctx, taskID, state)
}

func (m *manager) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	return m.service.DeleteTask(ctx, taskID)
}

func (m *manager) DeleteTaskSubtree(ctx context.Context, taskID uuid.UUID) error {
	return m.service.DeleteTaskSubtree(ctx, taskID)
}

// Task queries and analysis

func (m *manager) ListTasksHierarchical(ctx context.Context, projectID uuid.UUID) ([]*TaskHierarchy, error) {
	return m.service.ListTasksHierarchical(ctx, projectID)
}

func (m *manager) GetTaskSubtree(ctx context.Context, taskID uuid.UUID) (*TaskHierarchy, error) {
	return m.service.GetTaskSubtree(ctx, taskID)
}

func (m *manager) FindNextActionableTask(ctx context.Context, projectID uuid.UUID) (*Task, error) {
	return m.service.FindNextActionableTask(ctx, projectID)
}

func (m *manager) FindTasksNeedingBreakdown(ctx context.Context, projectID uuid.UUID) ([]*Task, error) {
	return m.service.FindTasksNeedingBreakdown(ctx, projectID)
}

func (m *manager) GetProjectProgress(ctx context.Context, projectID uuid.UUID) (*ProjectProgress, error) {
	return m.service.GetProjectProgress(ctx, projectID)
}

func (m *manager) ListTasksByState(ctx context.Context, projectID uuid.UUID, state TaskState) ([]*Task, error) {
	return m.service.ListTasksByState(ctx, projectID, state)
}

// Configuration

func (m *manager) GetConfig() *Config {
	return m.service.GetConfig()
}

func (m *manager) UpdateConfig(config *Config) {
	m.service.UpdateConfig(config)
}

// Ensure manager implements ProjectTaskManager
var _ ProjectTaskManager = (*manager)(nil)
var _ ToolSetProvider = (*toolSetProvider)(nil)