package project

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddTaskDependency(t *testing.T) {
	ctx := context.Background()
	manager := NewManager(DefaultConfig())

	// Create a project
	project, err := manager.CreateProject(ctx, "Test Project", "A test project")
	require.NoError(t, err)

	// Create two tasks
	task1, err := manager.CreateTask(ctx, project.ID, nil, "Task 1", "First task", 5)
	require.NoError(t, err)

	task2, err := manager.CreateTask(ctx, project.ID, nil, "Task 2", "Second task", 5)
	require.NoError(t, err)

	// Add a dependency
	updatedTask, err := manager.AddTaskDependency(ctx, task1.ID, task2.ID)
	require.NoError(t, err)

	// Verify the dependency was added
	assert.Len(t, updatedTask.Dependencies, 1)
	assert.Equal(t, task2.ID, updatedTask.Dependencies[0])

	// Verify the dependent was added
	dependentTask, err := manager.GetTask(ctx, task2.ID)
	require.NoError(t, err)
	assert.Len(t, dependentTask.Dependents, 1)
	assert.Equal(t, task1.ID, dependentTask.Dependents[0])
}

func TestRemoveTaskDependency(t *testing.T) {
	ctx := context.Background()
	manager := NewManager(DefaultConfig())

	// Create a project
	project, err := manager.CreateProject(ctx, "Test Project", "A test project")
	require.NoError(t, err)

	// Create two tasks
	task1, err := manager.CreateTask(ctx, project.ID, nil, "Task 1", "First task", 5)
	require.NoError(t, err)

	task2, err := manager.CreateTask(ctx, project.ID, nil, "Task 2", "Second task", 5)
	require.NoError(t, err)

	// Add a dependency
	_, err = manager.AddTaskDependency(ctx, task1.ID, task2.ID)
	require.NoError(t, err)

	// Remove the dependency
	updatedTask, err := manager.RemoveTaskDependency(ctx, task1.ID, task2.ID)
	require.NoError(t, err)

	// Verify the dependency was removed
	assert.Len(t, updatedTask.Dependencies, 0)

	// Verify the dependent was removed
	dependentTask, err := manager.GetTask(ctx, task2.ID)
	require.NoError(t, err)
	assert.Len(t, dependentTask.Dependents, 0)
}

func TestGetTaskDependencies(t *testing.T) {
	ctx := context.Background()
	manager := NewManager(DefaultConfig())

	// Create a project
	project, err := manager.CreateProject(ctx, "Test Project", "A test project")
	require.NoError(t, err)

	// Create three tasks
	task1, err := manager.CreateTask(ctx, project.ID, nil, "Task 1", "First task", 5)
	require.NoError(t, err)

	task2, err := manager.CreateTask(ctx, project.ID, nil, "Task 2", "Second task", 5)
	require.NoError(t, err)

	task3, err := manager.CreateTask(ctx, project.ID, nil, "Task 3", "Third task", 5)
	require.NoError(t, err)

	// Add dependencies
	_, err = manager.AddTaskDependency(ctx, task1.ID, task2.ID)
	require.NoError(t, err)

	_, err = manager.AddTaskDependency(ctx, task1.ID, task3.ID)
	require.NoError(t, err)

	// Get dependencies
	dependencies, err := manager.GetTaskDependencies(ctx, task1.ID)
	require.NoError(t, err)

	// Verify we got the right dependencies
	assert.Len(t, dependencies, 2)
	dependencyIDs := make(map[uuid.UUID]bool)
	for _, dep := range dependencies {
		dependencyIDs[dep.ID] = true
	}
	assert.True(t, dependencyIDs[task2.ID])
	assert.True(t, dependencyIDs[task3.ID])
}

func TestGetDependentTasks(t *testing.T) {
	ctx := context.Background()
	manager := NewManager(DefaultConfig())

	// Create a project
	project, err := manager.CreateProject(ctx, "Test Project", "A test project")
	require.NoError(t, err)

	// Create three tasks
	task1, err := manager.CreateTask(ctx, project.ID, nil, "Task 1", "First task", 5)
	require.NoError(t, err)

	task2, err := manager.CreateTask(ctx, project.ID, nil, "Task 2", "Second task", 5)
	require.NoError(t, err)

	task3, err := manager.CreateTask(ctx, project.ID, nil, "Task 3", "Third task", 5)
	require.NoError(t, err)

	// Add dependencies
	_, err = manager.AddTaskDependency(ctx, task2.ID, task1.ID)
	require.NoError(t, err)

	_, err = manager.AddTaskDependency(ctx, task3.ID, task1.ID)
	require.NoError(t, err)

	// Get dependents
	dependents, err := manager.GetDependentTasks(ctx, task1.ID)
	require.NoError(t, err)

	// Verify we got the right dependents
	assert.Len(t, dependents, 2)
	dependentIDs := make(map[uuid.UUID]bool)
	for _, dep := range dependents {
		dependentIDs[dep.ID] = true
	}
	assert.True(t, dependentIDs[task2.ID])
	assert.True(t, dependentIDs[task3.ID])
}

func TestCircularDependencyDetection(t *testing.T) {
	ctx := context.Background()
	manager := NewManager(DefaultConfig())

	// Create a project
	project, err := manager.CreateProject(ctx, "Test Project", "A test project")
	require.NoError(t, err)

	// Create three tasks
	task1, err := manager.CreateTask(ctx, project.ID, nil, "Task 1", "First task", 5)
	require.NoError(t, err)

	task2, err := manager.CreateTask(ctx, project.ID, nil, "Task 2", "Second task", 5)
	require.NoError(t, err)

	task3, err := manager.CreateTask(ctx, project.ID, nil, "Task 3", "Third task", 5)
	require.NoError(t, err)

	// Add a normal dependency
	_, err = manager.AddTaskDependency(ctx, task2.ID, task1.ID)
	require.NoError(t, err)

	// Add another normal dependency
	_, err = manager.AddTaskDependency(ctx, task3.ID, task2.ID)
	require.NoError(t, err)

	// Try to create a circular dependency - this should fail
	_, err = manager.AddTaskDependency(ctx, task1.ID, task3.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular dependency")
}

func TestCrossProjectDependencyRejection(t *testing.T) {
	ctx := context.Background()
	manager := NewManager(DefaultConfig())

	// Create two projects
	project1, err := manager.CreateProject(ctx, "Test Project 1", "First test project")
	require.NoError(t, err)

	project2, err := manager.CreateProject(ctx, "Test Project 2", "Second test project")
	require.NoError(t, err)

	// Create tasks in different projects
	task1, err := manager.CreateTask(ctx, project1.ID, nil, "Task 1", "First task", 5)
	require.NoError(t, err)

	task2, err := manager.CreateTask(ctx, project2.ID, nil, "Task 2", "Second task", 5)
	require.NoError(t, err)

	// Try to create a dependency between tasks in different projects - this should fail
	_, err = manager.AddTaskDependency(ctx, task1.ID, task2.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "same project")
}