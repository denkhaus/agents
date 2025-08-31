package project

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AddTaskDependency adds a dependency relationship between two tasks
func (s *service) AddTaskDependency(ctx context.Context, taskID uuid.UUID, dependsOnTaskID uuid.UUID) (*Task, error) {
	// Validate that both tasks exist
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	dependsOnTask, err := s.repo.GetTask(ctx, dependsOnTaskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get depends on task: %w", err)
	}

	// Validate that both tasks are in the same project
	if task.ProjectID != dependsOnTask.ProjectID {
		return nil, ValidationError{Field: "depends_on_task_id", Message: "tasks must be in the same project"}
	}

	// Check for circular dependencies
	if err := s.checkCircularDependency(ctx, taskID, dependsOnTaskID); err != nil {
		return nil, err
	}

	// Check if dependency already exists
	for _, depID := range task.Dependencies {
		if depID == dependsOnTaskID {
			// Dependency already exists, return the task as-is
			return task, nil
		}
	}

	// Add the dependency to the task
	task.Dependencies = append(task.Dependencies, dependsOnTaskID)

	// Add this task to the dependents list of the dependsOnTask
	dependsOnTask.Dependents = append(dependsOnTask.Dependents, taskID)

	// Update both tasks
	task.UpdatedAt = time.Now()
	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task with dependency: %w", err)
	}

	dependsOnTask.UpdatedAt = time.Now()
	if err := s.repo.UpdateTask(ctx, dependsOnTask); err != nil {
		return nil, fmt.Errorf("failed to update depends on task with dependent: %w", err)
	}

	return s.repo.GetTask(ctx, taskID)
}

// RemoveTaskDependency removes a dependency relationship between two tasks
func (s *service) RemoveTaskDependency(ctx context.Context, taskID uuid.UUID, dependsOnTaskID uuid.UUID) (*Task, error) {
	// Get both tasks
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	dependsOnTask, err := s.repo.GetTask(ctx, dependsOnTaskID)
	if err != nil {
		// If the dependsOnTask doesn't exist, we still want to remove it from the dependencies list
		// but we don't need to update the dependents list
		dependsOnTask = nil
	}

	// Remove the dependency from the task
	newDependencies := make([]uuid.UUID, 0, len(task.Dependencies))
	for _, depID := range task.Dependencies {
		if depID != dependsOnTaskID {
			newDependencies = append(newDependencies, depID)
		}
	}
	task.Dependencies = newDependencies

	// Remove this task from the dependents list of the dependsOnTask
	if dependsOnTask != nil {
		newDependents := make([]uuid.UUID, 0, len(dependsOnTask.Dependents))
		for _, depID := range dependsOnTask.Dependents {
			if depID != taskID {
				newDependents = append(newDependents, depID)
			}
		}
		dependsOnTask.Dependents = newDependents

		// Update the dependsOnTask
		dependsOnTask.UpdatedAt = time.Now()
		if err := s.repo.UpdateTask(ctx, dependsOnTask); err != nil {
			return nil, fmt.Errorf("failed to update depends on task: %w", err)
		}
	}

	// Update the task
	task.UpdatedAt = time.Now()
	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return s.repo.GetTask(ctx, taskID)
}

// GetTaskDependencies returns all tasks that the given task depends on
func (s *service) GetTaskDependencies(ctx context.Context, taskID uuid.UUID) ([]*Task, error) {
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	var dependencies []*Task
	for _, depID := range task.Dependencies {
		depTask, err := s.repo.GetTask(ctx, depID)
		if err != nil {
			// Skip dependencies that can't be found
			continue
		}
		dependencies = append(dependencies, depTask)
	}

	return dependencies, nil
}

// GetDependentTasks returns all tasks that depend on the given task
func (s *service) GetDependentTasks(ctx context.Context, taskID uuid.UUID) ([]*Task, error) {
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	var dependents []*Task
	for _, depID := range task.Dependents {
		depTask, err := s.repo.GetTask(ctx, depID)
		if err != nil {
			// Skip dependents that can't be found
			continue
		}
		dependents = append(dependents, depTask)
	}

	return dependents, nil
}

// checkCircularDependency checks if adding a dependency would create a circular dependency
func (s *service) checkCircularDependency(ctx context.Context, taskID uuid.UUID, dependsOnTaskID uuid.UUID) error {
	// If task depends on itself, that's a circular dependency
	if taskID == dependsOnTaskID {
		return ValidationError{Field: "depends_on_task_id", Message: "task cannot depend on itself"}
	}

	// Get all tasks that dependsOnTaskID depends on
	visited := make(map[uuid.UUID]bool)
	return s.checkCircularDependencyHelper(ctx, dependsOnTaskID, taskID, visited)
}

// checkCircularDependencyHelper is a recursive helper for detecting circular dependencies
func (s *service) checkCircularDependencyHelper(ctx context.Context, currentTaskID uuid.UUID, targetTaskID uuid.UUID, visited map[uuid.UUID]bool) error {
	// If we've already visited this task, skip it to avoid infinite loops
	if visited[currentTaskID] {
		return nil
	}

	// Mark this task as visited
	visited[currentTaskID] = true

	// Get the current task
	task, err := s.repo.GetTask(ctx, currentTaskID)
	if err != nil {
		// If we can't find the task, we can't check its dependencies
		return nil
	}

	// Check if any of the current task's dependencies is the target task
	for _, depID := range task.Dependencies {
		if depID == targetTaskID {
			return ValidationError{Field: "depends_on_task_id", Message: "adding this dependency would create a circular dependency"}
		}
		// Recursively check the dependencies of this dependency
		if err := s.checkCircularDependencyHelper(ctx, depID, targetTaskID, visited); err != nil {
			return err
		}
	}

	return nil
}