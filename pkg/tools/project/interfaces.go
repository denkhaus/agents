package project

import (
	"context"

	"github.com/denkhaus/agents/pkg/tools/project/shared"
	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// ProjectManager defines the public interface for project task management
type ProjectManager interface {
	// Project operations
	CreateProject(ctx context.Context, title, description string) (*shared.Project, error)
	GetProject(ctx context.Context, projectID uuid.UUID) (*shared.Project, error)
	UpdateProject(ctx context.Context, projectID uuid.UUID, title, description string) (*shared.Project, error)
	UpdateProjectDescription(ctx context.Context, projectID uuid.UUID, description string) (*shared.Project, error)
	DeleteProject(ctx context.Context, projectID uuid.UUID) error
	ListProjects(ctx context.Context) ([]*shared.Project, error)

	// Task operations
	CreateTask(ctx context.Context, projectID uuid.UUID, parentID *uuid.UUID, title, description string, complexity int) (*shared.Task, error)
	GetTask(ctx context.Context, taskID uuid.UUID) (*shared.Task, error)
	UpdateTask(ctx context.Context, taskID uuid.UUID, title, description string, complexity int, state shared.TaskState) (*shared.Task, error)
	UpdateTaskDescription(ctx context.Context, taskID uuid.UUID, description string) (*shared.Task, error)
	UpdateTaskState(ctx context.Context, taskID uuid.UUID, state shared.TaskState) (*shared.Task, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	DeleteTaskSubtree(ctx context.Context, taskID uuid.UUID) error

	// Task queries and analysis
	GetParentTask(ctx context.Context, taskID uuid.UUID) (*shared.Task, error)
	GetChildTasks(ctx context.Context, taskID uuid.UUID) ([]*shared.Task, error)
	GetRootTasks(ctx context.Context, projectID uuid.UUID) ([]*shared.Task, error)
	ListTasksForProject(ctx context.Context, projectID uuid.UUID) ([]*shared.Task, error)
	FindNextActionableTask(ctx context.Context, projectID uuid.UUID) (*shared.Task, error)
	FindTasksNeedingBreakdown(ctx context.Context, projectID uuid.UUID) ([]*shared.Task, error)
	GetProjectProgress(ctx context.Context, projectID uuid.UUID) (*shared.ProjectProgress, error)
	ListTasksByState(ctx context.Context, projectID uuid.UUID, state shared.TaskState) ([]*shared.Task, error)
	BulkUpdateTasks(ctx context.Context, taskIDs []uuid.UUID, updates shared.TaskUpdates) error
	DuplicateTask(ctx context.Context, taskID uuid.UUID, newProjectID uuid.UUID) (*shared.Task, error)
	SetTaskEstimate(ctx context.Context, taskID uuid.UUID, estimate int64) (*shared.Task, error)

	// Agent assignment management
	AssignTaskToAgent(ctx context.Context, taskID uuid.UUID, agentID uuid.UUID) (*shared.Task, error)
	UnassignTaskFromAgent(ctx context.Context, taskID uuid.UUID) (*shared.Task, error)
	ListTasksByAgent(ctx context.Context, projectID uuid.UUID, agentID uuid.UUID) ([]*shared.Task, error)
	ListUnassignedTasks(ctx context.Context, projectID uuid.UUID) ([]*shared.Task, error)

	// Dependency management
	AddTaskDependency(ctx context.Context, taskID uuid.UUID, dependsOnTaskID uuid.UUID) (*shared.Task, error)
	RemoveTaskDependency(ctx context.Context, taskID uuid.UUID, dependsOnTaskID uuid.UUID) (*shared.Task, error)
	GetTaskDependencies(ctx context.Context, taskID uuid.UUID) ([]*shared.Task, error)
	GetDependentTasks(ctx context.Context, taskID uuid.UUID) ([]*shared.Task, error)

	// Configuration
	GetConfig() *Config
	UpdateConfig(config *Config)
}

// ToolSetProvider defines the interface for creating project task tool sets
type ToolSetProvider interface {
	CreateToolSet(opts ...Option) (tool.ToolSet, error)
}

// Config holds configuration for the task management system
type Config struct {
	MaxTasksPerDepth     int // Maximum tasks allowed per depth level (applies to all depths)
	ComplexityThreshold  int // Threshold for task breakdown suggestions
	MaxDepth             int // Maximum allowed depth
	MaxDescriptionLength int // Maximum length for descriptions
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		MaxTasksPerDepth:     20, // Max 50 tasks per depth level
		ComplexityThreshold:  8,
		MaxDepth:             5,
		MaxDescriptionLength: 2000, // Default maximum description length
	}
}
