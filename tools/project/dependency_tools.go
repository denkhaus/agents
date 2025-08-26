package project

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// addTaskDependencyArgs defines the arguments for adding a task dependency
type addTaskDependencyArgs struct {
	TaskID          string `json:"task_id" description:"The ID of the task that will depend on another task"`
	DependsOnTaskID string `json:"depends_on_task_id" description:"The ID of the task that the first task will depend on"`
}

// addTaskDependencyResult defines the result of adding a task dependency
type addTaskDependencyResult struct {
	Task    *Task  `json:"task,omitempty" description:"The updated task"`
	Message string `json:"message" description:"A message describing the result"`
}

// addTaskDependency adds a dependency relationship between two tasks
func (pts *projectTaskToolSet) addTaskDependency(ctx context.Context, args addTaskDependencyArgs) (addTaskDependencyResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return addTaskDependencyResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	dependsOnTaskID, err := uuid.Parse(args.DependsOnTaskID)
	if err != nil {
		return addTaskDependencyResult{}, fmt.Errorf("invalid depends on task ID format: %w", err)
	}

	log.Printf("Adding dependency: task %s depends on task %s", taskID, dependsOnTaskID)

	task, err := pts.manager.AddTaskDependency(ctx, taskID, dependsOnTaskID)
	if err != nil {
		log.Printf("Failed to add task dependency: %v", err)
		return addTaskDependencyResult{}, err
	}

	log.Printf("Successfully added dependency: task %s depends on task %s", taskID, dependsOnTaskID)
	return addTaskDependencyResult{
		Task:    task,
		Message: fmt.Sprintf("Successfully added dependency: task %s depends on task %s", taskID, dependsOnTaskID),
	}, nil
}

// addTaskDependencyTool creates a tool for adding a task dependency
func (pts *projectTaskToolSet) addTaskDependencyTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.addTaskDependency,
		function.WithName("add_task_dependency"),
		function.WithDescription("Add a dependency relationship between two tasks"),
	)
}

// removeTaskDependencyArgs defines the arguments for removing a task dependency
type removeTaskDependencyArgs struct {
	TaskID          string `json:"task_id" description:"The ID of the task that depends on another task"`
	DependsOnTaskID string `json:"depends_on_task_id" description:"The ID of the task that the first task depends on"`
}

// removeTaskDependencyResult defines the result of removing a task dependency
type removeTaskDependencyResult struct {
	Task    *Task  `json:"task,omitempty" description:"The updated task"`
	Message string `json:"message" description:"A message describing the result"`
}

// removeTaskDependency removes a dependency relationship between two tasks
func (pts *projectTaskToolSet) removeTaskDependency(ctx context.Context, args removeTaskDependencyArgs) (removeTaskDependencyResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return removeTaskDependencyResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	dependsOnTaskID, err := uuid.Parse(args.DependsOnTaskID)
	if err != nil {
		return removeTaskDependencyResult{}, fmt.Errorf("invalid depends on task ID format: %w", err)
	}

	log.Printf("Removing dependency: task %s no longer depends on task %s", taskID, dependsOnTaskID)

	task, err := pts.manager.RemoveTaskDependency(ctx, taskID, dependsOnTaskID)
	if err != nil {
		log.Printf("Failed to remove task dependency: %v", err)
		return removeTaskDependencyResult{}, err
	}

	log.Printf("Successfully removed dependency: task %s no longer depends on task %s", taskID, dependsOnTaskID)
	return removeTaskDependencyResult{
		Task:    task,
		Message: fmt.Sprintf("Successfully removed dependency: task %s no longer depends on task %s", taskID, dependsOnTaskID),
	}, nil
}

// removeTaskDependencyTool creates a tool for removing a task dependency
func (pts *projectTaskToolSet) removeTaskDependencyTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.removeTaskDependency,
		function.WithName("remove_task_dependency"),
		function.WithDescription("Remove a dependency relationship between two tasks"),
	)
}

// getTaskDependenciesArgs defines the arguments for getting task dependencies
type getTaskDependenciesArgs struct {
	TaskID string `json:"task_id" description:"The ID of the task to get dependencies for"`
}

// getTaskDependenciesResult defines the result of getting task dependencies
type getTaskDependenciesResult struct {
	Tasks   []*Task `json:"tasks,omitempty" description:"The tasks that the specified task depends on"`
	Count   int     `json:"count" description:"The number of dependencies"`
	Message string  `json:"message" description:"A message describing the result"`
}

// getTaskDependencies gets all tasks that the specified task depends on
func (pts *projectTaskToolSet) getTaskDependencies(ctx context.Context, args getTaskDependenciesArgs) (getTaskDependenciesResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return getTaskDependenciesResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	log.Printf("Getting dependencies for task: %s", taskID)

	tasks, err := pts.manager.GetTaskDependencies(ctx, taskID)
	if err != nil {
		log.Printf("Failed to get task dependencies: %v", err)
		return getTaskDependenciesResult{}, err
	}

	log.Printf("Found %d dependencies for task %s", len(tasks), taskID)
	return getTaskDependenciesResult{
		Tasks:   tasks,
		Count:   len(tasks),
		Message: fmt.Sprintf("Found %d dependencies for task %s", len(tasks), taskID),
	}, nil
}

// getTaskDependenciesTool creates a tool for getting task dependencies
func (pts *projectTaskToolSet) getTaskDependenciesTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getTaskDependencies,
		function.WithName("get_task_dependencies"),
		function.WithDescription("Get all tasks that the specified task depends on"),
	)
}

// getDependentTasksArgs defines the arguments for getting dependent tasks
type getDependentTasksArgs struct {
	TaskID string `json:"task_id" description:"The ID of the task to get dependents for"`
}

// getDependentTasksResult defines the result of getting dependent tasks
type getDependentTasksResult struct {
	Tasks   []*Task `json:"tasks,omitempty" description:"The tasks that depend on the specified task"`
	Count   int     `json:"count" description:"The number of dependent tasks"`
	Message string  `json:"message" description:"A message describing the result"`
}

// getDependentTasks gets all tasks that depend on the specified task
func (pts *projectTaskToolSet) getDependentTasks(ctx context.Context, args getDependentTasksArgs) (getDependentTasksResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return getDependentTasksResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	log.Printf("Getting dependent tasks for task: %s", taskID)

	tasks, err := pts.manager.GetDependentTasks(ctx, taskID)
	if err != nil {
		log.Printf("Failed to get dependent tasks: %v", err)
		return getDependentTasksResult{}, err
	}

	log.Printf("Found %d dependent tasks for task %s", len(tasks), taskID)
	return getDependentTasksResult{
		Tasks:   tasks,
		Count:   len(tasks),
		Message: fmt.Sprintf("Found %d dependent tasks for task %s", len(tasks), taskID),
	}, nil
}

// getDependentTasksTool creates a tool for getting dependent tasks
func (pts *projectTaskToolSet) getDependentTasksTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.getDependentTasks,
		function.WithName("get_dependent_tasks"),
		function.WithDescription("Get all tasks that depend on the specified task"),
	)
}