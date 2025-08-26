package project

import (
	"time"

	"github.com/google/uuid"
)

// TaskState represents the current state of a task
type TaskState string

const (
	TaskStatePending    TaskState = "pending"
	TaskStateInProgress TaskState = "in-progress"
	TaskStateCompleted  TaskState = "completed"
	TaskStateBlocked    TaskState = "blocked"
	TaskStateCancelled  TaskState = "cancelled"
)

// Task represents a single task in the project hierarchy
type Task struct {
	ID           uuid.UUID    `json:"id"`
	ProjectID    uuid.UUID    `json:"project_id"`
	ParentID     *uuid.UUID   `json:"parent_id,omitempty"` // nil for root tasks
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	State        TaskState    `json:"state"`
	Complexity   int          `json:"complexity"` // Used for breakdown decisions
	Priority     int          `json:"priority"`   // Higher number = higher priority
	Depth        int          `json:"depth"`      // 0 for root tasks
	Estimate     *int64       `json:"estimate,omitempty"` // Time estimate in minutes
	Dependencies []uuid.UUID  `json:"dependencies,omitempty"` // Tasks this task depends on
	Dependents   []uuid.UUID  `json:"dependents,omitempty"`   // Tasks that depend on this task
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	CompletedAt  *time.Time   `json:"completed_at,omitempty"`
}

// Project represents a project containing hierarchical tasks
type Project struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Progress metrics
	TotalTasks     int     `json:"total_tasks"`
	CompletedTasks int     `json:"completed_tasks"`
	Progress       float64 `json:"progress"` // Percentage (0-100)
}

// ProjectProgress represents detailed progress information
type ProjectProgress struct {
	ProjectID       uuid.UUID   `json:"project_id"`
	TotalTasks      int         `json:"total_tasks"`
	CompletedTasks  int         `json:"completed_tasks"`
	InProgressTasks int         `json:"in_progress_tasks"`
	PendingTasks    int         `json:"pending_tasks"`
	BlockedTasks    int         `json:"blocked_tasks"`
	CancelledTasks  int         `json:"cancelled_tasks"`
	OverallProgress float64     `json:"overall_progress"`
	TasksByDepth    map[int]int `json:"tasks_by_depth"`
}

// TaskFilter represents filtering options for task queries
type TaskFilter struct {
	ProjectID     *uuid.UUID `json:"project_id,omitempty"`
	ParentID      *uuid.UUID `json:"parent_id,omitempty"`
	State         *TaskState `json:"state,omitempty"`
	MinDepth      *int       `json:"min_depth,omitempty"`
	MaxDepth      *int       `json:"max_depth,omitempty"`
	MinComplexity *int       `json:"min_complexity,omitempty"`
	MaxComplexity *int       `json:"max_complexity,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// TaskUpdates represents the fields that can be updated in bulk
type TaskUpdates struct {
	State      *TaskState `json:"state,omitempty"`
	Priority   *int       `json:"priority,omitempty"`
	Complexity *int       `json:"complexity,omitempty"`
}

// Tool Input/Output Types

// createProjectArgs holds the input for creating a project
type createProjectArgs struct {
	Title   string `json:"title" description:"Project title (required, max 200 chars)"`
	Details string `json:"details" description:"Project details (optional, max 2000 chars)"`
}

// createProjectResult holds the output for creating a project
type createProjectResult struct {
	Project *Project `json:"project"`
	Message string   `json:"message"`
}

// getProjectArgs holds the input for getting a project
type getProjectArgs struct {
	ProjectID string `json:"project_id" description:"Project UUID"`
}

// updateProjectDescriptionArgs holds the input for updating a project description
type updateProjectDescriptionArgs struct {
	ProjectID   string `json:"project_id" description:"Project UUID"`
	Description string `json:"description" description:"New project description (max 2000 chars)"`
}

// listProjectsArgs holds the input for listing projects (empty struct for no parameters)
type listProjectsArgs struct{}

// listProjectsResult holds the output for listing projects
type listProjectsResult struct {
	Projects []*Project `json:"projects"`
	Count    int        `json:"count"`
}

// createTaskArgs holds the input for creating a task
type createTaskArgs struct {
	ProjectID   string  `json:"project_id" description:"Project UUID (required)"`
	ParentID    *string `json:"parent_id,omitempty" description:"Parent task UUID (optional, for subtasks)"`
	Title       string  `json:"title" description:"Task title (required, max 200 chars)"`
	Description string  `json:"description" description:"Task description (optional, max 2000 chars)"`
	Complexity  int     `json:"complexity" description:"Task complexity (1-10, used for breakdown decisions)"`
	Priority    int     `json:"priority" description:"Task priority (1-10, higher = more important)"`
}

// updateTaskDescriptionArgs holds the input for updating a task description
type updateTaskDescriptionArgs struct {
	TaskID      string `json:"task_id" description:"Task UUID"`
	Description string `json:"description" description:"New task description (max 2000 chars)"`
}

// getTaskArgs holds the input for getting a task
type getTaskArgs struct {
	TaskID string `json:"task_id" description:"Task UUID"`
}

// updateTaskStateArgs holds the input for updating task state
type updateTaskStateArgs struct {
	TaskID string    `json:"task_id" description:"Task UUID"`
	State  TaskState `json:"state" description:"New task state (pending, in-progress, completed, blocked, cancelled)"`
}

// getProjectProgressArgs holds the input for getting project progress
type getProjectProgressArgs struct {
	ProjectID string `json:"project_id" description:"Project UUID"`
}

// getChildTasksArgs holds the input for getting child tasks
type getChildTasksArgs struct {
	TaskID string `json:"task_id" description:"Task UUID"`
}

// getChildTasksResult holds the output for getting child tasks
type getChildTasksResult struct {
	Tasks []*Task `json:"tasks"`
	Count int     `json:"count"`
}

// getParentTaskArgs holds the input for getting the parent task
type getParentTaskArgs struct {
	TaskID string `json:"task_id" description:"Task UUID"`
}
