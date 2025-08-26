package projecttasks

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
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"project_id"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"` // nil for root tasks
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       TaskState `json:"state"`
	Complexity  int       `json:"complexity"` // Used for breakdown decisions
	Priority    int       `json:"priority"`   // Higher number = higher priority
	Depth       int       `json:"depth"`      // 0 for root tasks
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
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

// TaskHierarchy represents a task with its children for display purposes
type TaskHierarchy struct {
	Task     *Task            `json:"task"`
	Children []*TaskHierarchy `json:"children,omitempty"`
}

// ProjectProgress represents detailed progress information
type ProjectProgress struct {
	ProjectID        uuid.UUID `json:"project_id"`
	TotalTasks       int       `json:"total_tasks"`
	CompletedTasks   int       `json:"completed_tasks"`
	InProgressTasks  int       `json:"in_progress_tasks"`
	PendingTasks     int       `json:"pending_tasks"`
	BlockedTasks     int       `json:"blocked_tasks"`
	CancelledTasks   int       `json:"cancelled_tasks"`
	OverallProgress  float64   `json:"overall_progress"`
	TasksByDepth     map[int]int `json:"tasks_by_depth"`
}

// TaskFilter represents filtering options for task queries
type TaskFilter struct {
	ProjectID *uuid.UUID  `json:"project_id,omitempty"`
	ParentID  *uuid.UUID  `json:"parent_id,omitempty"`
	State     *TaskState  `json:"state,omitempty"`
	MinDepth  *int        `json:"min_depth,omitempty"`
	MaxDepth  *int        `json:"max_depth,omitempty"`
	MinComplexity *int    `json:"min_complexity,omitempty"`
	MaxComplexity *int    `json:"max_complexity,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}