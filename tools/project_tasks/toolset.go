package projecttasks

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// projectTaskToolSet implements the ToolSet interface for project task management
type projectTaskToolSet struct {
	manager ProjectTaskManager
	tools   []tool.CallableTool
}

// Option is a functional option for configuring the project task tool set
type Option func(*projectTaskToolSet)

// WithManager sets a custom manager instance
func WithManager(manager ProjectTaskManager) Option {
	return func(pts *projectTaskToolSet) {
		pts.manager = manager
	}
}

// WithConfig sets a custom configuration
func WithConfig(config *Config) Option {
	return func(pts *projectTaskToolSet) {
		if pts.manager != nil {
			pts.manager.UpdateConfig(config)
		}
	}
}

// NewToolSet creates a new project task management tool set
func NewToolSet(opts ...Option) (tool.ToolSet, error) {
	toolSet := &projectTaskToolSet{
		manager: NewManager(DefaultConfig()),
	}

	// Apply options
	for _, opt := range opts {
		opt(toolSet)
	}

	if toolSet.manager == nil {
		return nil, fmt.Errorf("manager cannot be nil")
	}

	// Initialize tools
	toolSet.tools = []tool.CallableTool{
		toolSet.createProjectTool(),
		toolSet.getProjectTool(),
		toolSet.listProjectsTool(),
		toolSet.createTaskTool(),
		toolSet.getTaskTool(),
		toolSet.updateTaskStateTool(),
		toolSet.listTasksHierarchicalTool(),
		toolSet.getProjectProgressTool(),
	}

	return toolSet, nil
}

// Tools returns the list of available tools
func (pts *projectTaskToolSet) Tools(ctx context.Context) []tool.CallableTool {
	return pts.tools
}

// Close cleans up resources
func (pts *projectTaskToolSet) Close() error {
	return nil
}

// Project management tools

func (pts *projectTaskToolSet) createProjectTool() tool.CallableTool {
	createProject := func(ctx context.Context, input struct {
		Title       string `json:"title" description:"Project title (required, max 200 chars)"`
		Description string `json:"description" description:"Project description (optional, max 2000 chars)"`
	}) (interface{}, error) {
		log.Printf("Creating project: %s", input.Title)
		
		project, err := pts.manager.CreateProject(ctx, input.Title, input.Description)
		if err != nil {
			log.Printf("Failed to create project: %v", err)
			return nil, err
		}
		
		log.Printf("Created project %s successfully", project.ID)
		return project, nil
	}
	
	return function.NewFunctionTool(
		createProject,
		function.WithName("create_project"),
		function.WithDescription("Create a new project for task management"),
	)
}

func (pts *projectTaskToolSet) getProjectTool() tool.CallableTool {
	getProject := func(ctx context.Context, input struct {
		ProjectID string `json:"project_id" description:"Project UUID"`
	}) (interface{}, error) {
		projectID, err := uuid.Parse(input.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("invalid project ID format: %w", err)
		}

		log.Printf("Getting project: %s", projectID)
		
		project, err := pts.manager.GetProject(ctx, projectID)
		if err != nil {
			log.Printf("Failed to get project: %v", err)
			return nil, err
		}
		
		return project, nil
	}
	
	return function.NewFunctionTool(
		getProject,
		function.WithName("get_project"),
		function.WithDescription("Get project details by ID"),
	)
}

func (pts *projectTaskToolSet) listProjectsTool() tool.CallableTool {
	listProjects := func(ctx context.Context, input struct{}) (interface{}, error) {
		log.Printf("Listing all projects")
		
		projects, err := pts.manager.ListProjects(ctx)
		if err != nil {
			log.Printf("Failed to list projects: %v", err)
			return nil, err
		}
		
		return map[string]interface{}{
			"projects": projects,
			"count":    len(projects),
		}, nil
	}
	
	return function.NewFunctionTool(
		listProjects,
		function.WithName("list_projects"),
		function.WithDescription("List all projects"),
	)
}

func (pts *projectTaskToolSet) createTaskTool() tool.CallableTool {
	createTask := func(ctx context.Context, input struct {
		ProjectID   string  `json:"project_id" description:"Project UUID (required)"`
		ParentID    *string `json:"parent_id,omitempty" description:"Parent task UUID (optional, for subtasks)"`
		Title       string  `json:"title" description:"Task title (required, max 200 chars)"`
		Description string  `json:"description" description:"Task description (optional, max 2000 chars)"`
		Complexity  int     `json:"complexity" description:"Task complexity (1-10, used for breakdown decisions)"`
		Priority    int     `json:"priority" description:"Task priority (1-10, higher = more important)"`
	}) (interface{}, error) {
		projectID, err := uuid.Parse(input.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("invalid project ID format: %w", err)
		}

		var parentID *uuid.UUID
		if input.ParentID != nil {
			pid, err := uuid.Parse(*input.ParentID)
			if err != nil {
				return nil, fmt.Errorf("invalid parent ID format: %w", err)
			}
			parentID = &pid
		}

		log.Printf("Creating task in project %s: %s", projectID, input.Title)
		
		task, err := pts.manager.CreateTask(ctx, projectID, parentID, input.Title, input.Description, input.Complexity, input.Priority)
		if err != nil {
			log.Printf("Failed to create task: %v", err)
			return nil, err
		}
		
		log.Printf("Created task %s successfully", task.ID)
		return task, nil
	}
	
	return function.NewFunctionTool(
		createTask,
		function.WithName("create_task"),
		function.WithDescription("Create a new task in a project"),
	)
}

func (pts *projectTaskToolSet) getTaskTool() tool.CallableTool {
	getTask := func(ctx context.Context, input struct {
		TaskID string `json:"task_id" description:"Task UUID"`
	}) (interface{}, error) {
		taskID, err := uuid.Parse(input.TaskID)
		if err != nil {
			return nil, fmt.Errorf("invalid task ID format: %w", err)
		}

		log.Printf("Getting task: %s", taskID)
		
		task, err := pts.manager.GetTask(ctx, taskID)
		if err != nil {
			log.Printf("Failed to get task: %v", err)
			return nil, err
		}
		
		return task, nil
	}
	
	return function.NewFunctionTool(
		getTask,
		function.WithName("get_task"),
		function.WithDescription("Get task details by ID"),
	)
}

func (pts *projectTaskToolSet) updateTaskStateTool() tool.CallableTool {
	updateTaskState := func(ctx context.Context, input struct {
		TaskID string    `json:"task_id" description:"Task UUID"`
		State  TaskState `json:"state" description:"New task state (pending, in-progress, completed, blocked, cancelled)"`
	}) (interface{}, error) {
		taskID, err := uuid.Parse(input.TaskID)
		if err != nil {
			return nil, fmt.Errorf("invalid task ID format: %w", err)
		}

		log.Printf("Updating task state: %s to %s", taskID, input.State)
		
		task, err := pts.manager.UpdateTaskState(ctx, taskID, input.State)
		if err != nil {
			log.Printf("Failed to update task state: %v", err)
			return nil, err
		}
		
		return task, nil
	}
	
	return function.NewFunctionTool(
		updateTaskState,
		function.WithName("update_task_state"),
		function.WithDescription("Update only the task state"),
	)
}

func (pts *projectTaskToolSet) listTasksHierarchicalTool() tool.CallableTool {
	listTasksHierarchical := func(ctx context.Context, input struct {
		ProjectID string `json:"project_id" description:"Project UUID"`
	}) (interface{}, error) {
		projectID, err := uuid.Parse(input.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("invalid project ID format: %w", err)
		}

		log.Printf("Listing hierarchical tasks for project: %s", projectID)
		
		hierarchy, err := pts.manager.ListTasksHierarchical(ctx, projectID)
		if err != nil {
			log.Printf("Failed to list hierarchical tasks: %v", err)
			return nil, err
		}
		
		return map[string]interface{}{
			"project_id": projectID.String(),
			"hierarchy":  hierarchy,
		}, nil
	}
	
	return function.NewFunctionTool(
		listTasksHierarchical,
		function.WithName("list_tasks_hierarchical"),
		function.WithDescription("List all tasks in a project with hierarchical structure"),
	)
}

func (pts *projectTaskToolSet) getProjectProgressTool() tool.CallableTool {
	getProjectProgress := func(ctx context.Context, input struct {
		ProjectID string `json:"project_id" description:"Project UUID"`
	}) (interface{}, error) {
		projectID, err := uuid.Parse(input.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("invalid project ID format: %w", err)
		}

		log.Printf("Getting project progress: %s", projectID)
		
		progress, err := pts.manager.GetProjectProgress(ctx, projectID)
		if err != nil {
			log.Printf("Failed to get project progress: %v", err)
			return nil, err
		}
		
		return progress, nil
	}
	
	return function.NewFunctionTool(
		getProjectProgress,
		function.WithName("get_project_progress"),
		function.WithDescription("Get detailed progress metrics for a project"),
	)
}