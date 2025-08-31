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
	ID           uuid.UUID   `json:"id"`
	ProjectID    uuid.UUID   `json:"project_id"`
	ParentID     *uuid.UUID  `json:"parent_id,omitempty"` // nil for root tasks
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	State        TaskState   `json:"state"`
	Complexity   int         `json:"complexity"`             // Used for breakdown decisions
	Depth        int         `json:"depth"`                  // 0 for root tasks
	Estimate     *int64      `json:"estimate,omitempty"`     // Time estimate in minutes
	Dependencies []uuid.UUID `json:"dependencies,omitempty"` // Tasks this task depends on
	Dependents   []uuid.UUID `json:"dependents,omitempty"`   // Tasks that depend on this task
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	CompletedAt  *time.Time  `json:"completed_at,omitempty"`
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

// deleteTaskArgs defines the arguments for deleting a task
type deleteTaskArgs struct {
	TaskID string `json:"task_id" description:"The ID of the task to delete"`
}

// deleteTaskResult defines the result of deleting a task
type deleteTaskResult struct {
	Message string `json:"message" description:"A message describing the result"`
}

// updateTaskArgs defines the arguments for updating a task
type updateTaskArgs struct {
	TaskID      string    `json:"task_id" description:"The ID of the task to update"`
	Title       string    `json:"title" description:"The new title for the task"`
	Description string    `json:"description" description:"The new description for the task"`
	Complexity  int       `json:"complexity" description:"The new complexity for the task (1-10)"`
	State       TaskState `json:"state" description:"The new state for the task"`
}

// updateTaskResult defines the result of updating a task
type updateTaskResult struct {
	Task    *Task  `json:"task,omitempty" description:"The updated task"`
	Message string `json:"message" description:"A message describing the result"`
}

// deleteTaskSubtreeArgs defines the arguments for deleting a task subtree
type deleteTaskSubtreeArgs struct {
	TaskID string `json:"task_id" description:"The ID of the task whose subtree to delete"`
}

// deleteTaskSubtreeResult defines the result of deleting a task subtree
type deleteTaskSubtreeResult struct {
	Message string `json:"message" description:"A message describing the result"`
}

// listTasksByStateArgs defines the arguments for listing tasks by state
type listTasksByStateArgs struct {
	ProjectID string    `json:"project_id" description:"The ID of the project to list tasks from"`
	State     TaskState `json:"state" description:"The state of tasks to list"`
}

// listTasksByStateResult defines the result of listing tasks by state
type listTasksByStateResult struct {
	Tasks   []*Task `json:"tasks,omitempty" description:"The tasks with the specified state, if any"`
	Count   int     `json:"count" description:"The number of tasks with the specified state"`
	Message string  `json:"message" description:"A message describing the result"`
}

// getRootTasksArgs defines the arguments for getting root tasks
type getRootTasksArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to get root tasks from"`
}

// getRootTasksResult defines the result of getting root tasks
type getRootTasksResult struct {
	Tasks   []*Task `json:"tasks,omitempty" description:"The root tasks, if any"`
	Count   int     `json:"count" description:"The number of root tasks"`
	Message string  `json:"message" description:"A message describing the result"`
}

// findTasksNeedingBreakdownArgs defines the arguments for finding tasks needing breakdown
type findTasksNeedingBreakdownArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to find tasks needing breakdown in"`
}

// findTasksNeedingBreakdownResult defines the result of finding tasks needing breakdown
type findTasksNeedingBreakdownResult struct {
	Tasks   []*Task `json:"tasks,omitempty" description:"The tasks needing breakdown, if any"`
	Count   int     `json:"count" description:"The number of tasks needing breakdown"`
	Message string  `json:"message" description:"A message describing the result"`
}

// findNextActionableTaskArgs defines the arguments for finding the next actionable task
type findNextActionableTaskArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to find the next actionable task in"`
}

// findNextActionableTaskResult defines the result of finding the next actionable task
type findNextActionableTaskResult struct {
	Task    *Task  `json:"task,omitempty" description:"The next actionable task, if found"`
	Message string `json:"message" description:"A message describing the result"`
}

// updateProjectArgs defines the arguments for updating a project
type updateProjectArgs struct {
	ProjectID   string `json:"project_id" description:"The ID of the project to update"`
	Title       string `json:"title" description:"The new title for the project"`
	Description string `json:"description" description:"The new description for the project"`
}

// updateProjectResult defines the result of updating a project
type updateProjectResult struct {
	Project *Project `json:"project,omitempty" description:"The updated project"`
	Message string   `json:"message" description:"A message describing the result"`
}

// deleteProjectArgs defines the arguments for deleting a project
type deleteProjectArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to delete"`
}

// deleteProjectResult defines the result of deleting a project
type deleteProjectResult struct {
	Message string `json:"message" description:"A message describing the result"`
}

// listTasksForProjectArgs defines the arguments for listing all tasks in a project
type listTasksForProjectArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to list tasks from"`
}

// listTasksForProjectResult defines the result of listing all tasks in a project
type listTasksForProjectResult struct {
	Tasks   []*Task `json:"tasks,omitempty" description:"All tasks in the project"`
	Count   int     `json:"count" description:"The number of tasks in the project"`
	Message string  `json:"message" description:"A message describing the result"`
}

// bulkUpdateTasksArgs defines the arguments for bulk updating tasks
type bulkUpdateTasksArgs struct {
	TaskIDs    []string   `json:"task_ids" description:"The IDs of the tasks to update"`
	State      *TaskState `json:"state,omitempty" description:"The new state for the tasks"`
	Complexity *int       `json:"complexity,omitempty" description:"The new complexity for the tasks (1-10)"`
}

// bulkUpdateTasksResult defines the result of bulk updating tasks
type bulkUpdateTasksResult struct {
	Message string `json:"message" description:"A message describing the result"`
	Count   int    `json:"count" description:"The number of tasks that were updated"`
}

// duplicateTaskArgs defines the arguments for duplicating a task
type duplicateTaskArgs struct {
	TaskID       string `json:"task_id" description:"The ID of the task to duplicate"`
	NewProjectID string `json:"new_project_id" description:"The ID of the project to duplicate the task to"`
}

// duplicateTaskResult defines the result of duplicating a task
type duplicateTaskResult struct {
	Task    *Task  `json:"task,omitempty" description:"The duplicated task"`
	Message string `json:"message" description:"A message describing the result"`
}

// setTaskEstimateArgs defines the arguments for setting a task estimate
type setTaskEstimateArgs struct {
	TaskID   string `json:"task_id" description:"The ID of the task to set estimate for"`
	Estimate int64  `json:"estimate" description:"The time estimate in minutes"`
}

// setTaskEstimateResult defines the result of setting a task estimate
type setTaskEstimateResult struct {
	Task    *Task  `json:"task,omitempty" description:"The updated task"`
	Message string `json:"message" description:"A message describing the result"`
}
