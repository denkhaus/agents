package project

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
func NewManager(config *Config) ProjectManager {
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

func (m *manager) UpdateProjectDescription(ctx context.Context, projectID uuid.UUID, description string) (*Project, error) {
	return m.service.UpdateProjectDescription(ctx, projectID, description)
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

func (m *manager) UpdateTaskDescription(ctx context.Context, taskID uuid.UUID, description string) (*Task, error) {
	return m.service.UpdateTaskDescription(ctx, taskID, description)
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

func (m *manager) GetParentTask(ctx context.Context, taskID uuid.UUID) (*Task, error) {
	return m.service.GetParentTask(ctx, taskID)
}

func (m *manager) GetChildTasks(ctx context.Context, taskID uuid.UUID) ([]*Task, error) {
	return m.service.GetChildTasks(ctx, taskID)
}

func (m *manager) GetRootTasks(ctx context.Context, projectID uuid.UUID) ([]*Task, error) {
	return m.service.GetRootTasks(ctx, projectID)
}

func (m *manager) ListTasksForProject(ctx context.Context, projectID uuid.UUID) ([]*Task, error) {
	return m.service.ListTasksForProject(ctx, projectID)
}

func (m *manager) BulkUpdateTasks(ctx context.Context, taskIDs []uuid.UUID, updates TaskUpdates) error {
	return m.service.BulkUpdateTasks(ctx, taskIDs, updates)
}

func (m *manager) DuplicateTask(ctx context.Context, taskID uuid.UUID, newProjectID uuid.UUID) (*Task, error) {
	return m.service.DuplicateTask(ctx, taskID, newProjectID)
}

func (m *manager) SetTaskEstimate(ctx context.Context, taskID uuid.UUID, estimate int64) (*Task, error) {
	return m.service.SetTaskEstimate(ctx, taskID, estimate)
}

func (m *manager) AddTaskDependency(ctx context.Context, taskID uuid.UUID, dependsOnTaskID uuid.UUID) (*Task, error) {
	return m.service.AddTaskDependency(ctx, taskID, dependsOnTaskID)
}

func (m *manager) RemoveTaskDependency(ctx context.Context, taskID uuid.UUID, dependsOnTaskID uuid.UUID) (*Task, error) {
	return m.service.RemoveTaskDependency(ctx, taskID, dependsOnTaskID)
}

func (m *manager) GetTaskDependencies(ctx context.Context, taskID uuid.UUID) ([]*Task, error) {
	return m.service.GetTaskDependencies(ctx, taskID)
}

func (m *manager) GetDependentTasks(ctx context.Context, taskID uuid.UUID) ([]*Task, error) {
	return m.service.GetDependentTasks(ctx, taskID)
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
var _ ProjectManager = (*manager)(nil)
var _ ToolSetProvider = (*toolSetProvider)(nil)
