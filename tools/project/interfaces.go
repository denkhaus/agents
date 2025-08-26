package project

import (
	"context"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// ProjectManager defines the public interface for project task management
type ProjectManager interface {
	// Project operations
	CreateProject(ctx context.Context, title, description string) (*Project, error)
	GetProject(ctx context.Context, projectID uuid.UUID) (*Project, error)
	UpdateProject(ctx context.Context, projectID uuid.UUID, title, description string) (*Project, error)
	UpdateProjectDescription(ctx context.Context, projectID uuid.UUID, description string) (*Project, error)
	DeleteProject(ctx context.Context, projectID uuid.UUID) error
	ListProjects(ctx context.Context) ([]*Project, error)

	// Task operations
	CreateTask(ctx context.Context, projectID uuid.UUID, parentID *uuid.UUID, title, description string, complexity, priority int) (*Task, error)
	GetTask(ctx context.Context, taskID uuid.UUID) (*Task, error)
	UpdateTask(ctx context.Context, taskID uuid.UUID, title, description string, complexity, priority int, state TaskState) (*Task, error)
	UpdateTaskDescription(ctx context.Context, taskID uuid.UUID, description string) (*Task, error)
	UpdateTaskState(ctx context.Context, taskID uuid.UUID, state TaskState) (*Task, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	DeleteTaskSubtree(ctx context.Context, taskID uuid.UUID) error

	// Task queries and analysis
	GetParentTask(ctx context.Context, taskID uuid.UUID) (*Task, error)
	GetChildTasks(ctx context.Context, taskID uuid.UUID) ([]*Task, error)
	GetRootTasks(ctx context.Context, projectID uuid.UUID) ([]*Task, error)
	ListTasksForProject(ctx context.Context, projectID uuid.UUID) ([]*Task, error)
	FindNextActionableTask(ctx context.Context, projectID uuid.UUID) (*Task, error)
	FindTasksNeedingBreakdown(ctx context.Context, projectID uuid.UUID) ([]*Task, error)
	GetProjectProgress(ctx context.Context, projectID uuid.UUID) (*ProjectProgress, error)
	ListTasksByState(ctx context.Context, projectID uuid.UUID, state TaskState) ([]*Task, error)
	BulkUpdateTasks(ctx context.Context, taskIDs []uuid.UUID, updates TaskUpdates) error
	DuplicateTask(ctx context.Context, taskID uuid.UUID, newProjectID uuid.UUID) (*Task, error)
	SetTaskEstimate(ctx context.Context, taskID uuid.UUID, estimate int64) (*Task, error)
	AddTaskDependency(ctx context.Context, taskID uuid.UUID, dependsOnTaskID uuid.UUID) (*Task, error)
	RemoveTaskDependency(ctx context.Context, taskID uuid.UUID, dependsOnTaskID uuid.UUID) (*Task, error)
	GetTaskDependencies(ctx context.Context, taskID uuid.UUID) ([]*Task, error)
	GetDependentTasks(ctx context.Context, taskID uuid.UUID) ([]*Task, error)

	// Configuration
	GetConfig() *Config
	UpdateConfig(config *Config)
}

// ToolSetProvider defines the interface for creating project task tool sets
type ToolSetProvider interface {
	CreateToolSet(opts ...Option) (tool.ToolSet, error)
}

// repository defines the internal interface for task and project persistence
type repository interface {
	// Project operations
	CreateProject(ctx context.Context, project *Project) error
	GetProject(ctx context.Context, id uuid.UUID) (*Project, error)
	UpdateProject(ctx context.Context, project *Project) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
	ListProjects(ctx context.Context) ([]*Project, error)

	// Task operations
	CreateTask(ctx context.Context, task *Task) error
	GetTask(ctx context.Context, id uuid.UUID) (*Task, error)
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error

	// Task queries
	ListTasks(ctx context.Context, filter TaskFilter) ([]*Task, error)
	GetTasksByProject(ctx context.Context, projectID uuid.UUID) ([]*Task, error)
	GetTasksByParent(ctx context.Context, parentID uuid.UUID) ([]*Task, error)
	GetRootTasks(ctx context.Context, projectID uuid.UUID) ([]*Task, error)
	GetParentTask(ctx context.Context, taskID uuid.UUID) (*Task, error)

	// Hierarchy operations
	DeleteTaskSubtree(ctx context.Context, taskID uuid.UUID) error

	// Metrics and analysis
	GetProjectProgress(ctx context.Context, projectID uuid.UUID) (*ProjectProgress, error)
	GetTaskCountByDepth(ctx context.Context, projectID uuid.UUID, maxDepth int) (map[int]int, error)
}

// Config holds configuration for the task management system
type Config struct {
	MaxTasksPerDepth    int // Maximum tasks allowed per depth level (applies to all depths)
	ComplexityThreshold int // Threshold for task breakdown suggestions
	MaxDepth            int // Maximum allowed depth
	DefaultPriority     int // Default priority for new tasks
	MaxDescriptionLength int // Maximum length for descriptions
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		MaxTasksPerDepth:    20, // Max 50 tasks per depth level
		ComplexityThreshold: 8,
		MaxDepth:            5,
		DefaultPriority:     5,
		MaxDescriptionLength: 2000, // Default maximum description length
	}
}
