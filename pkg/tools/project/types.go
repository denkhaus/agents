package project

import (
	"fmt"

	"github.com/denkhaus/agents/pkg/tools/project/shared"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// Tool Input/Output Types

// createProjectArgs holds the input for creating a project
type createProjectArgs struct {
	Title   string `json:"title" description:"Project title (required, max 200 chars)"`
	Details string `json:"details" description:"Project details (optional, max 2000 chars)"`
}

// createProjectResult holds the output for creating a project
type createProjectResult struct {
	Project *shared.Project `json:"project"`
	Message string          `json:"message"`
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
	Projects []*shared.Project `json:"projects"`
	Count    int               `json:"count"`
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
	TaskID string           `json:"task_id" description:"Task UUID"`
	State  shared.TaskState `json:"state" description:"New task state (pending, in-progress, completed, blocked, cancelled)"`
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
	Tasks []*shared.Task `json:"tasks"`
	Count int            `json:"count"`
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
	TaskID      string           `json:"task_id" description:"The ID of the task to update"`
	Title       string           `json:"title" description:"The new title for the task"`
	Description string           `json:"description" description:"The new description for the task"`
	Complexity  int              `json:"complexity" description:"The new complexity for the task (1-10)"`
	State       shared.TaskState `json:"state" description:"The new state for the task"`
}

// updateTaskResult defines the result of updating a task
type updateTaskResult struct {
	Task    *shared.Task `json:"task,omitempty" description:"The updated task"`
	Message string       `json:"message" description:"A message describing the result"`
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
	ProjectID string           `json:"project_id" description:"The ID of the project to list tasks from"`
	State     shared.TaskState `json:"state" description:"The state of tasks to list"`
}

// listTasksByStateResult defines the result of listing tasks by state
type listTasksByStateResult struct {
	Tasks   []*shared.Task `json:"tasks,omitempty" description:"The tasks with the specified state, if any"`
	Count   int            `json:"count" description:"The number of tasks with the specified state"`
	Message string         `json:"message" description:"A message describing the result"`
}

// getRootTasksArgs defines the arguments for getting root tasks
type getRootTasksArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to get root tasks from"`
}

// getRootTasksResult defines the result of getting root tasks
type getRootTasksResult struct {
	Tasks   []*shared.Task `json:"tasks,omitempty" description:"The root tasks, if any"`
	Count   int            `json:"count" description:"The number of root tasks"`
	Message string         `json:"message" description:"A message describing the result"`
}

// findTasksNeedingBreakdownArgs defines the arguments for finding tasks needing breakdown
type findTasksNeedingBreakdownArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to find tasks needing breakdown in"`
}

// findTasksNeedingBreakdownResult defines the result of finding tasks needing breakdown
type findTasksNeedingBreakdownResult struct {
	Tasks   []*shared.Task `json:"tasks,omitempty" description:"The tasks needing breakdown, if any"`
	Count   int            `json:"count" description:"The number of tasks needing breakdown"`
	Message string         `json:"message" description:"A message describing the result"`
}

// findNextActionableTaskArgs defines the arguments for finding the next actionable task
type findNextActionableTaskArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to find the next actionable task in"`
}

// findNextActionableTaskResult defines the result of finding the next actionable task
type findNextActionableTaskResult struct {
	Task    *shared.Task `json:"task,omitempty" description:"The next actionable task, if found"`
	Message string       `json:"message" description:"A message describing the result"`
}

// updateProjectArgs defines the arguments for updating a project
type updateProjectArgs struct {
	ProjectID   string `json:"project_id" description:"The ID of the project to update"`
	Title       string `json:"title" description:"The new title for the project"`
	Description string `json:"description" description:"The new description for the project"`
}

// updateProjectResult defines the result of updating a project
type updateProjectResult struct {
	Project *shared.Project `json:"project,omitempty" description:"The updated project"`
	Message string          `json:"message" description:"A message describing the result"`
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
	Tasks   []*shared.Task `json:"tasks,omitempty" description:"All tasks in the project"`
	Count   int            `json:"count" description:"The number of tasks in the project"`
	Message string         `json:"message" description:"A message describing the result"`
}

// bulkUpdateTasksArgs defines the arguments for bulk updating tasks
type bulkUpdateTasksArgs struct {
	TaskIDs    []string          `json:"task_ids" description:"The IDs of the tasks to update"`
	State      *shared.TaskState `json:"state,omitempty" description:"The new state for the tasks"`
	Complexity *int              `json:"complexity,omitempty" description:"The new complexity for the tasks (1-10)"`
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
	Task    *shared.Task `json:"task,omitempty" description:"The duplicated task"`
	Message string       `json:"message" description:"A message describing the result"`
}

// setTaskEstimateArgs defines the arguments for setting a task estimate
type setTaskEstimateArgs struct {
	TaskID   string `json:"task_id" description:"The ID of the task to set estimate for"`
	Estimate int64  `json:"estimate" description:"The time estimate in minutes"`
}

// setTaskEstimateResult defines the result of setting a task estimate
type setTaskEstimateResult struct {
	Task    *shared.Task `json:"task,omitempty" description:"The updated task"`
	Message string       `json:"message" description:"A message describing the result"`
}

// listAvailableAgentsArgs defines the arguments for listing available agents (empty struct)
type listAvailableAgentsArgs struct{}

// listAvailableAgentsResult defines the result of listing available agents
type listAvailableAgentsResult struct {
	Agents  []AgentInfo `json:"agents" description:"List of available agents in the system"`
	Count   int         `json:"count" description:"Number of available agents"`
	Message string      `json:"message" description:"A message describing the result"`
}

// AgentInfo represents information about an available agent
type AgentInfo struct {
	ID          string `json:"id" description:"Unique identifier of the agent"`
	Name        string `json:"name" description:"Display name of the agent"`
	Role        string `json:"role" description:"Role/type of the agent (e.g., researcher, coder, project-manager)"`
	Description string `json:"description" description:"Detailed description of the agent's capabilities and purpose"`
}

// assignTaskToAgentArgs defines the arguments for assigning a task to an agent
type assignTaskToAgentArgs struct {
	TaskID  string `json:"task_id" description:"The ID of the task to assign"`
	AgentID string `json:"agent_id" description:"The ID of the agent to assign the task to"`
}

// assignTaskToAgentResult defines the result of assigning a task to an agent
type assignTaskToAgentResult struct {
	Task    *shared.Task `json:"task,omitempty" description:"The updated task with agent assignment"`
	Message string       `json:"message" description:"A message describing the result"`
}

// unassignTaskFromAgentArgs defines the arguments for unassigning a task from an agent
type unassignTaskFromAgentArgs struct {
	TaskID string `json:"task_id" description:"The ID of the task to unassign"`
}

// unassignTaskFromAgentResult defines the result of unassigning a task from an agent
type unassignTaskFromAgentResult struct {
	Task    *shared.Task `json:"task,omitempty" description:"The updated task with removed agent assignment"`
	Message string       `json:"message" description:"A message describing the result"`
}

// listTasksByAgentArgs defines the arguments for listing tasks assigned to a specific agent
type listTasksByAgentArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to search in"`
	AgentID   string `json:"agent_id" description:"The ID of the agent to find tasks for"`
}

// listTasksByAgentResult defines the result of listing tasks assigned to a specific agent
type listTasksByAgentResult struct {
	Tasks   []*shared.Task `json:"tasks,omitempty" description:"Tasks assigned to the specified agent"`
	Count   int            `json:"count" description:"Number of tasks assigned to the agent"`
	Message string         `json:"message" description:"A message describing the result"`
}

// listUnassignedTasksArgs defines the arguments for listing unassigned tasks
type listUnassignedTasksArgs struct {
	ProjectID string `json:"project_id" description:"The ID of the project to search in"`
}

// listUnassignedTasksResult defines the result of listing unassigned tasks
type listUnassignedTasksResult struct {
	Tasks   []*shared.Task `json:"tasks,omitempty" description:"Tasks that have no agent assigned"`
	Count   int            `json:"count" description:"Number of unassigned tasks"`
	Message string         `json:"message" description:"A message describing the result"`
}

type ProjectRepositoryType string

const (
	ProjectRepositoryTypeInMemory ProjectRepositoryType = "inmemory"
	ProjectRepositoryTypePostgres ProjectRepositoryType = "postgres"
)

func (p ProjectRepositoryType) String() string {
	return string(p)
}

// Validate checks if the AgentRole is a valid defined role
func (p ProjectRepositoryType) Validate() error {
	switch p {
	case ProjectRepositoryTypeInMemory,
		ProjectRepositoryTypePostgres:
		return nil
	default:
		return fmt.Errorf("invalid project repository type: %s. Valid types are: %s, %s",
			p, ProjectRepositoryTypeInMemory, ProjectRepositoryTypePostgres)
	}
}

// ToolSetConfig holds configuration for the project management toolset
type ToolSetConfig struct {
	RepositoryType ProjectRepositoryType `json:"repository_type" mapstructure:"repository_type"`
	DatabaseURL    *string               `json:"database_url" mapstructure:"database_url"`
	IsReadOnly     bool                  `json:"read_only" mapstructure:"read_only"`
}
