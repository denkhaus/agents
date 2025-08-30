package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap" // Add this import
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const (
	ToolSetName = "project_toolset"
)

// projectTaskToolSet implements the ToolSet interface for project task management
type projectTaskToolSet struct {
	manager    ProjectManager
	logger     *zap.Logger
	isReadOnly bool
	tools      []tool.CallableTool
}

// NewToolSet creates a new project task management tool set
func NewToolSet(opts ...Option) (tool.ToolSet, error) {
	toolSet := &projectTaskToolSet{
		manager: NewManager(DefaultConfig()),
		logger:  zap.NewNop(), // Use null logger by default to avoid interfering with chat output
	}

	// Apply options
	for _, opt := range opts {
		opt(toolSet)
	}

	if toolSet.manager == nil {
		return nil, fmt.Errorf("manager cannot be nil")
	}

	if toolSet.isReadOnly {
		// Initialize readonly tools
		toolSet.tools = []tool.CallableTool{
			toolSet.getProjectTool(),
			toolSet.listProjectsTool(),
			toolSet.getTaskTool(),
			toolSet.getProjectProgressTool(),
			toolSet.getChildTasksTool(),
			toolSet.getParentTaskTool(),
			toolSet.findNextActionableTaskTool(),
			toolSet.findTasksNeedingBreakdownTool(),
			toolSet.getRootTasksTool(),
			toolSet.listTasksByStateTool(),
			toolSet.listTasksForProjectTool(),
			toolSet.getTaskDependenciesTool(),
			toolSet.getDependentTasksTool(),
		}
	} else {
		// Initialize all tools
		toolSet.tools = []tool.CallableTool{
			toolSet.createProjectTool(),
			toolSet.getProjectTool(),
			toolSet.updateProjectDescriptionTool(), // Add this line
			toolSet.listProjectsTool(),
			toolSet.createTaskTool(),
			toolSet.getTaskTool(),
			toolSet.updateTaskDescriptionTool(), // Add this line
			toolSet.updateTaskStateTool(),
			toolSet.getProjectProgressTool(),
			toolSet.getChildTasksTool(),
			toolSet.getParentTaskTool(),
			toolSet.findNextActionableTaskTool(),
			toolSet.findTasksNeedingBreakdownTool(),
			toolSet.getRootTasksTool(),
			toolSet.listTasksByStateTool(),
			toolSet.deleteTaskSubtreeTool(),
			toolSet.updateTaskTool(),
			toolSet.deleteTaskTool(),
			toolSet.updateProjectTool(),
			toolSet.deleteProjectTool(),
			toolSet.listTasksForProjectTool(),
			toolSet.bulkUpdateTasksTool(),
			toolSet.duplicateTaskTool(),
			toolSet.setTaskEstimateTool(),
			toolSet.addTaskDependencyTool(),
			toolSet.removeTaskDependencyTool(),
			toolSet.getTaskDependenciesTool(),
			toolSet.getDependentTasksTool(),
		}
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

// createProject performs project creation
func (pts *projectTaskToolSet) createProject(ctx context.Context, args createProjectArgs) (createProjectResult, error) {
	pts.logger.Info("Creating project", zap.String("title", args.Title))

	project, err := pts.manager.CreateProject(ctx, args.Title, args.Details)
	if err != nil {
		pts.logger.Error("Failed to create project", zap.Error(err))
		return createProjectResult{}, err
	}

	pts.logger.Info("Created project successfully", zap.String("projectID", project.ID.String()))
	return createProjectResult{
		Project: project,
		Message: "Project created successfully",
	}, nil
}

func (pts *projectTaskToolSet) createProjectTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.createProject,
		function.WithName("create_project"),
		function.WithDescription("Create a new project for task management"),
	)
}

// getProject performs project retrieval
func (pts *projectTaskToolSet) getProject(ctx context.Context, args getProjectArgs) (*Project, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Getting project", zap.String("projectID", projectID.String()))

	project, err := pts.manager.GetProject(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to get project", zap.Error(err))
		return nil, err
	}

	return project, nil
}

func (pts *projectTaskToolSet) getProjectTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getProject,
		function.WithName("get_project"),
		function.WithDescription("Get project details by ID"),
	)
}

// updateProjectDescription performs project description update
func (pts *projectTaskToolSet) updateProjectDescription(ctx context.Context, args updateProjectDescriptionArgs) (*Project, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Updating project description", zap.String("projectID", projectID.String()))

	project, err := pts.manager.UpdateProjectDescription(ctx, projectID, args.Description)
	if err != nil {
		pts.logger.Error("Failed to update project description", zap.Error(err))
		return nil, err
	}

	return project, nil
}

func (pts *projectTaskToolSet) updateProjectDescriptionTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.updateProjectDescription,
		function.WithName("update_project_description"),
		function.WithDescription("Update only the project description"),
	)
}

// listProjects performs project listing
func (pts *projectTaskToolSet) listProjects(ctx context.Context, args listProjectsArgs) (listProjectsResult, error) {
	pts.logger.Info("Listing all projects")

	projects, err := pts.manager.ListProjects(ctx)
	if err != nil {
		pts.logger.Error("Failed to list projects", zap.Error(err))
		return listProjectsResult{}, err
	}

	return listProjectsResult{
		Projects: projects,
		Count:    len(projects),
	}, nil
}

func (pts *projectTaskToolSet) listProjectsTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.listProjects,
		function.WithName("list_projects"),
		function.WithDescription("List all projects"),
	)
}

// createTask performs task creation
func (pts *projectTaskToolSet) createTask(ctx context.Context, args createTaskArgs) (*Task, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID format: %w", err)
	}

	var parentID *uuid.UUID
	if args.ParentID != nil {
		pid, err := uuid.Parse(*args.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent ID format: %w", err)
		}
		parentID = &pid
	}

	pts.logger.Info("Creating task in project", zap.String("projectID", projectID.String()), zap.String("title", args.Title))

	task, err := pts.manager.CreateTask(ctx, projectID, parentID, args.Title, args.Description, args.Complexity)
	if err != nil {
		pts.logger.Error("Failed to create task", zap.Error(err))
		return nil, err
	}

	pts.logger.Info("Created task successfully", zap.String("taskID", task.ID.String()))
	return task, nil
}

func (pts *projectTaskToolSet) createTaskTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.createTask,
		function.WithName("create_task"),
		function.WithDescription("Create a new task in a project"),
	)
}

// getTask performs task retrieval
func (pts *projectTaskToolSet) getTask(ctx context.Context, args getTaskArgs) (*Task, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Getting task", zap.String("taskID", taskID.String()))

	task, err := pts.manager.GetTask(ctx, taskID)
	if err != nil {
		pts.logger.Error("Failed to get task", zap.Error(err))
		return nil, err
	}

	return task, nil
}

func (pts *projectTaskToolSet) getTaskTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getTask,
		function.WithName("get_task"),
		function.WithDescription("Get task details by ID"),
	)
}

// updateTaskDescription performs task description update
func (pts *projectTaskToolSet) updateTaskDescription(ctx context.Context, args updateTaskDescriptionArgs) (*Task, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Updating task description", zap.String("taskID", taskID.String()))

	task, err := pts.manager.UpdateTaskDescription(ctx, taskID, args.Description)
	if err != nil {
		pts.logger.Error("Failed to update task description", zap.Error(err))
		return nil, err
	}

	return task, nil
}

func (pts *projectTaskToolSet) updateTaskDescriptionTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.updateTaskDescription,
		function.WithName("update_task_description"),
		function.WithDescription("Update only the task description"),
	)
}

// updateTaskState performs task state update
func (pts *projectTaskToolSet) updateTaskState(ctx context.Context, args updateTaskStateArgs) (*Task, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Updating task state", zap.String("taskID", taskID.String()), zap.String("state", string(args.State)))

	task, err := pts.manager.UpdateTaskState(ctx, taskID, args.State)
	if err != nil {
		pts.logger.Error("Failed to update task state", zap.Error(err))
		return nil, err
	}

	return task, nil
}

func (pts *projectTaskToolSet) updateTaskStateTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.updateTaskState,
		function.WithName("update_task_state"),
		function.WithDescription("Update only the task state"),
	)
}

func (pts *projectTaskToolSet) getProjectProgress(ctx context.Context, args getProjectProgressArgs) (*ProjectProgress, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Getting project progress", zap.String("projectID", projectID.String()))

	progress, err := pts.manager.GetProjectProgress(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to get project progress", zap.Error(err))
		return nil, err
	}

	return progress, nil
}

func (pts *projectTaskToolSet) getProjectProgressTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getProjectProgress,
		function.WithName("get_project_progress"),
		function.WithDescription("Get detailed progress metrics for a project"),
	)
}

func (pts *projectTaskToolSet) getChildTasks(ctx context.Context, args getChildTasksArgs) (getChildTasksResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return getChildTasksResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Getting child tasks", zap.String("taskID", taskID.String()))

	tasks, err := pts.manager.GetChildTasks(ctx, taskID)
	if err != nil {
		pts.logger.Error("Failed to get child tasks", zap.Error(err))
		return getChildTasksResult{}, err
	}

	return getChildTasksResult{
		Tasks: tasks,
		Count: len(tasks),
	}, nil
}

func (pts *projectTaskToolSet) getChildTasksTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getChildTasks,
		function.WithName("get_child_tasks"),
		function.WithDescription("Get the child tasks of a given task"),
	)
}

func (pts *projectTaskToolSet) getParentTask(ctx context.Context, args getParentTaskArgs) (*Task, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Getting parent task", zap.String("taskID", taskID.String()))

	task, err := pts.manager.GetParentTask(ctx, taskID)
	if err != nil {
		pts.logger.Error("Failed to get parent task", zap.Error(err))
		return nil, err
	}

	return task, nil
}

func (pts *projectTaskToolSet) getParentTaskTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getParentTask,
		function.WithName("get_parent_task"),
		function.WithDescription("Get the parent task of a given task"),
	)
}

// findNextActionableTask finds the next actionable task in a project
func (pts *projectTaskToolSet) findNextActionableTask(ctx context.Context, args findNextActionableTaskArgs) (findNextActionableTaskResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return findNextActionableTaskResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Finding next actionable task", zap.String("projectID", projectID.String()))

	task, err := pts.manager.FindNextActionableTask(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to find next actionable task", zap.Error(err))
		return findNextActionableTaskResult{
			Message: fmt.Sprintf("No actionable task found: %v", err),
		}, nil
	}

	pts.logger.Info("Found next actionable task", zap.String("taskID", task.ID.String()))
	return findNextActionableTaskResult{
		Task:    task,
		Message: "Successfully found next actionable task",
	}, nil
}

// findNextActionableTaskTool creates a tool for finding the next actionable task
func (pts *projectTaskToolSet) findNextActionableTaskTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.findNextActionableTask,
		function.WithName("find_next_actionable_task"),
		function.WithDescription("Find the next actionable task in a project"),
	)
}

// findTasksNeedingBreakdown finds tasks that need to be broken down
func (pts *projectTaskToolSet) findTasksNeedingBreakdown(ctx context.Context, args findTasksNeedingBreakdownArgs) (findTasksNeedingBreakdownResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return findTasksNeedingBreakdownResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Finding tasks needing breakdown", zap.String("projectID", projectID.String()))

	tasks, err := pts.manager.FindTasksNeedingBreakdown(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to find tasks needing breakdown", zap.Error(err))
		return findTasksNeedingBreakdownResult{}, err
	}

	pts.logger.Info("Found tasks needing breakdown", zap.Int("count", len(tasks)))
	return findTasksNeedingBreakdownResult{
		Tasks:   tasks,
		Count:   len(tasks),
		Message: fmt.Sprintf("Found %d tasks needing breakdown", len(tasks)),
	}, nil
}

// findTasksNeedingBreakdownTool creates a tool for finding tasks needing breakdown
func (pts *projectTaskToolSet) findTasksNeedingBreakdownTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.findTasksNeedingBreakdown,
		function.WithName("find_tasks_needing_breakdown"),
		function.WithDescription("Find tasks that need to be broken down into smaller subtasks"),
	)
}

// getRootTasks gets the root tasks of a project
func (pts *projectTaskToolSet) getRootTasks(ctx context.Context, args getRootTasksArgs) (getRootTasksResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return getRootTasksResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Getting root tasks", zap.String("projectID", projectID.String()))

	tasks, err := pts.manager.GetRootTasks(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to get root tasks", zap.Error(err))
		return getRootTasksResult{}, err
	}

	pts.logger.Info("Found root tasks", zap.Int("count", len(tasks)))
	return getRootTasksResult{
		Tasks:   tasks,
		Count:   len(tasks),
		Message: fmt.Sprintf("Found %d root tasks", len(tasks)),
	}, nil
}

// getRootTasksTool creates a tool for getting root tasks
func (pts *projectTaskToolSet) getRootTasksTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getRootTasks,
		function.WithName("get_root_tasks"),
		function.WithDescription("Get the root tasks of a project"),
	)
}

// listTasksByState lists tasks by their state
func (pts *projectTaskToolSet) listTasksByState(ctx context.Context, args listTasksByStateArgs) (listTasksByStateResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return listTasksByStateResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Listing tasks with state", zap.String("state", string(args.State)), zap.String("projectID", projectID.String()))

	tasks, err := pts.manager.ListTasksByState(ctx, projectID, args.State)
	if err != nil {
		pts.logger.Error("Failed to list tasks by state", zap.Error(err))
		return listTasksByStateResult{}, err
	}

	pts.logger.Info("Found tasks with state", zap.Int("count", len(tasks)), zap.String("state", string(args.State)))
	return listTasksByStateResult{
		Tasks:   tasks,
		Count:   len(tasks),
		Message: fmt.Sprintf("Found %d tasks with state %s", len(tasks), args.State),
	}, nil
}

// listTasksByStateTool creates a tool for listing tasks by state
func (pts *projectTaskToolSet) listTasksByStateTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.listTasksByState,
		function.WithName("list_tasks_by_state"),
		function.WithDescription("List tasks by their state in a project"),
	)
}

// deleteTaskSubtree deletes a task and all its descendants
func (pts *projectTaskToolSet) deleteTaskSubtree(ctx context.Context, args deleteTaskSubtreeArgs) (deleteTaskSubtreeResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return deleteTaskSubtreeResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Deleting task subtree", zap.String("taskID", taskID.String()))

	err = pts.manager.DeleteTaskSubtree(ctx, taskID)
	if err != nil {
		pts.logger.Error("Failed to delete task subtree", zap.Error(err))
		return deleteTaskSubtreeResult{}, err
	}

	pts.logger.Info("Successfully deleted task subtree", zap.String("taskID", taskID.String()))
	return deleteTaskSubtreeResult{
		Message: fmt.Sprintf("Successfully deleted task subtree: %s", taskID),
	}, nil
}

// deleteTaskSubtreeTool creates a tool for deleting a task subtree
func (pts *projectTaskToolSet) deleteTaskSubtreeTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.deleteTaskSubtree,
		function.WithName("delete_task_subtree"),
		function.WithDescription("Delete a task and all its descendants"),
	)
}

// updateTask updates a task with all fields
func (pts *projectTaskToolSet) updateTask(ctx context.Context, args updateTaskArgs) (updateTaskResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return updateTaskResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Updating task", zap.String("taskID", taskID.String()))

	task, err := pts.manager.UpdateTask(ctx, taskID, args.Title, args.Description, args.Complexity, args.State)
	if err != nil {
		pts.logger.Error("Failed to update task", zap.Error(err))
		return updateTaskResult{}, err
	}

	pts.logger.Info("Successfully updated task", zap.String("taskID", task.ID.String()))
	return updateTaskResult{
		Task:    task,
		Message: fmt.Sprintf("Successfully updated task: %s", task.ID),
	}, nil
}

// updateTaskTool creates a tool for updating a task
func (pts *projectTaskToolSet) updateTaskTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.updateTask,
		function.WithName("update_task"),
		function.WithDescription("Update a task with all fields"),
	)
}

// deleteTask deletes a task
func (pts *projectTaskToolSet) deleteTask(ctx context.Context, args deleteTaskArgs) (deleteTaskResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return deleteTaskResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Deleting task", zap.String("taskID", taskID.String()))

	err = pts.manager.DeleteTask(ctx, taskID)
	if err != nil {
		pts.logger.Error("Failed to delete task", zap.Error(err))
		return deleteTaskResult{}, err
	}

	pts.logger.Info("Successfully deleted task", zap.String("taskID", taskID.String()))
	return deleteTaskResult{
		Message: fmt.Sprintf("Successfully deleted task: %s", taskID),
	}, nil
}

// deleteTaskTool creates a tool for deleting a task
func (pts *projectTaskToolSet) deleteTaskTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.deleteTask,
		function.WithName("delete_task"),
		function.WithDescription("Delete a task"),
	)
}

// updateProject updates a project with all fields
func (pts *projectTaskToolSet) updateProject(ctx context.Context, args updateProjectArgs) (updateProjectResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return updateProjectResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Updating project", zap.String("projectID", projectID.String()))

	project, err := pts.manager.UpdateProject(ctx, projectID, args.Title, args.Description)
	if err != nil {
		pts.logger.Error("Failed to update project", zap.Error(err))
		return updateProjectResult{}, err
	}

	pts.logger.Info("Successfully updated project", zap.String("projectID", project.ID.String()))
	return updateProjectResult{
		Project: project,
		Message: fmt.Sprintf("Successfully updated project: %s", project.ID),
	}, nil
}

// updateProjectTool creates a tool for updating a project
func (pts *projectTaskToolSet) updateProjectTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.updateProject,
		function.WithName("update_project"),
		function.WithDescription("Update a project with all fields"),
	)
}

// deleteProject deletes a project
func (pts *projectTaskToolSet) deleteProject(ctx context.Context, args deleteProjectArgs) (deleteProjectResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return deleteProjectResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Deleting project", zap.String("projectID", projectID.String()))

	err = pts.manager.DeleteProject(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to delete project", zap.Error(err))
		return deleteProjectResult{}, err
	}

	pts.logger.Info("Successfully deleted project", zap.String("projectID", projectID.String()))
	return deleteProjectResult{
		Message: fmt.Sprintf("Successfully deleted project: %s", projectID),
	}, nil
}

// deleteProjectTool creates a tool for deleting a project
func (pts *projectTaskToolSet) deleteProjectTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.deleteProject,
		function.WithName("delete_project"),
		function.WithDescription("Delete a project"),
	)
}

// listTasksForProject lists all tasks in a project regardless of hierarchy level
func (pts *projectTaskToolSet) listTasksForProject(ctx context.Context, args listTasksForProjectArgs) (listTasksForProjectResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return listTasksForProjectResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Listing all tasks for project", zap.String("projectID", projectID.String()))

	tasks, err := pts.manager.ListTasksForProject(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to list tasks for project", zap.Error(err))
		return listTasksForProjectResult{}, err
	}

	pts.logger.Info("Found tasks in project", zap.Int("count", len(tasks)))
	return listTasksForProjectResult{
		Tasks:   tasks,
		Count:   len(tasks),
		Message: fmt.Sprintf("Found %d tasks in project", len(tasks)),
	}, nil
}

// listTasksForProjectTool creates a tool for listing all tasks in a project
func (pts *projectTaskToolSet) listTasksForProjectTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.listTasksForProject,
		function.WithName("list_tasks_for_project"),
		function.WithDescription("List all tasks in a project regardless of hierarchy level"),
	)
}

// bulkUpdateTasks bulk updates multiple tasks with the same updates
func (pts *projectTaskToolSet) bulkUpdateTasks(ctx context.Context, args bulkUpdateTasksArgs) (bulkUpdateTasksResult, error) {
	if len(args.TaskIDs) == 0 {
		return bulkUpdateTasksResult{
			Message: "No task IDs provided",
			Count:   0,
		}, nil
	}

	// Parse task IDs
	taskIDs := make([]uuid.UUID, len(args.TaskIDs))
	for i, taskIDStr := range args.TaskIDs {
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return bulkUpdateTasksResult{}, fmt.Errorf("invalid task ID format at index %d: %w", i, err)
		}
		taskIDs[i] = taskID
	}

	// Create updates object
	updates := TaskUpdates{
		State:      args.State,
		Complexity: args.Complexity,
	}

	pts.logger.Info("Bulk updating tasks", zap.Int("count", len(taskIDs)))

	err := pts.manager.BulkUpdateTasks(ctx, taskIDs, updates)
	if err != nil {
		pts.logger.Error("Failed to bulk update tasks", zap.Error(err))
		return bulkUpdateTasksResult{}, err
	}

	pts.logger.Info("Successfully bulk updated tasks", zap.Int("count", len(taskIDs)))
	return bulkUpdateTasksResult{
		Message: fmt.Sprintf("Successfully updated %d tasks", len(taskIDs)),
		Count:   len(taskIDs),
	}, nil
}

// bulkUpdateTasksTool creates a tool for bulk updating tasks
func (pts *projectTaskToolSet) bulkUpdateTasksTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.bulkUpdateTasks,
		function.WithName("bulk_update_tasks"),
		function.WithDescription("Bulk update multiple tasks with the same updates"),
	)
}

// duplicateTask duplicates a task in a new project
func (pts *projectTaskToolSet) duplicateTask(ctx context.Context, args duplicateTaskArgs) (duplicateTaskResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return duplicateTaskResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	newProjectID, err := uuid.Parse(args.NewProjectID)
	if err != nil {
		return duplicateTaskResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Duplicating task", zap.String("taskID", taskID.String()), zap.String("newProjectID", newProjectID.String()))

	task, err := pts.manager.DuplicateTask(ctx, taskID, newProjectID)
	if err != nil {
		pts.logger.Error("Failed to duplicate task", zap.Error(err))
		return duplicateTaskResult{}, err
	}

	pts.logger.Info("Successfully duplicated task", zap.String("originalTaskID", taskID.String()), zap.String("newTaskID", task.ID.String()))
	return duplicateTaskResult{
		Task:    task,
		Message: fmt.Sprintf("Successfully duplicated task as %s", task.ID),
	}, nil
}

// duplicateTaskTool creates a tool for duplicating a task
func (pts *projectTaskToolSet) duplicateTaskTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.duplicateTask,
		function.WithName("duplicate_task"),
		function.WithDescription("Duplicate a task in a new project"),
	)
}

// setTaskEstimate sets the time estimate for a task
func (pts *projectTaskToolSet) setTaskEstimate(ctx context.Context, args setTaskEstimateArgs) (setTaskEstimateResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return setTaskEstimateResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Setting estimate for task", zap.String("taskID", taskID.String()), zap.Int64("estimateMinutes", args.Estimate))

	task, err := pts.manager.SetTaskEstimate(ctx, taskID, args.Estimate)
	if err != nil {
		pts.logger.Error("Failed to set task estimate", zap.Error(err))
		return setTaskEstimateResult{}, err
	}

	pts.logger.Info("Successfully set estimate for task", zap.String("taskID", taskID.String()))
	return setTaskEstimateResult{
		Task:    task,
		Message: fmt.Sprintf("Successfully set estimate for task %s", taskID),
	}, nil
}

// setTaskEstimateTool creates a tool for setting a task estimate
func (pts *projectTaskToolSet) setTaskEstimateTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.setTaskEstimate,
		function.WithName("set_task_estimate"),
		function.WithDescription("Set the time estimate for a task in minutes"),
	)
}
