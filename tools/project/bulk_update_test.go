package project

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func TestBulkUpdateTasks(t *testing.T) {
	ctx := context.Background()

	// Create toolset
	toolSet, err := NewToolSet()
	require.NoError(t, err)
	defer toolSet.Close()

	tools := toolSet.Tools(ctx)
	require.NotEmpty(t, tools)

	// Helper function to find tool by name
	findTool := func(name string) tool.CallableTool {
		for _, tool := range tools {
			if tool.Declaration().Name == name {
				return tool
			}
		}
		t.Fatalf("Tool %s not found", name)
		return nil
	}

	// Create project
	createProjectTool := findTool("create_project")
	projectInput := map[string]interface{}{
		"title":   "Test Project",
		"details": "A test project for unit testing",
	}
	projectInputJSON, _ := json.Marshal(projectInput)

	projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
	require.NoError(t, err)

	project := projectResult.(createProjectResult).Project
	projectID := project.ID.String()

	// Create tasks
	createTaskTool := findTool("create_task")

	// Create first task
	task1Input := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Task 1",
		"description": "First task",
		"complexity":  5,
	}
	task1InputJSON, _ := json.Marshal(task1Input)

	task1Result, err := createTaskTool.Call(ctx, task1InputJSON)
	require.NoError(t, err)
	task1 := task1Result.(*Task)

	// Create second task
	task2Input := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Task 2",
		"description": "Second task",
		"complexity":  3,
	}
	task2InputJSON, _ := json.Marshal(task2Input)

	task2Result, err := createTaskTool.Call(ctx, task2InputJSON)
	require.NoError(t, err)
	task2 := task2Result.(*Task)

	// Bulk update tasks - change state to completed
	bulkUpdateTool := findTool("bulk_update_tasks")
	newState := "completed"
	bulkInput := map[string]interface{}{
		"task_ids": []string{task1.ID.String(), task2.ID.String()},
		"state":    newState,
	}
	bulkInputJSON, _ := json.Marshal(bulkInput)

	bulkResult, err := bulkUpdateTool.Call(ctx, bulkInputJSON)
	require.NoError(t, err)

	result := bulkResult.(bulkUpdateTasksResult)
	assert.Equal(t, 2, result.Count)
	assert.Equal(t, "Successfully updated 2 tasks", result.Message)

	// Verify tasks were updated
	getTaskTool := findTool("get_task")

	// Check first task
	getTask1Input := map[string]interface{}{
		"task_id": task1.ID.String(),
	}
	getTask1InputJSON, _ := json.Marshal(getTask1Input)

	getTask1Result, err := getTaskTool.Call(ctx, getTask1InputJSON)
	require.NoError(t, err)

	updatedTask1 := getTask1Result.(*Task)
	assert.Equal(t, TaskStateCompleted, updatedTask1.State)
	assert.NotNil(t, updatedTask1.CompletedAt)

	// Check second task
	getTask2Input := map[string]interface{}{
		"task_id": task2.ID.String(),
	}
	getTask2InputJSON, _ := json.Marshal(getTask2Input)

	getTask2Result, err := getTaskTool.Call(ctx, getTask2InputJSON)
	require.NoError(t, err)

	updatedTask2 := getTask2Result.(*Task)
	assert.Equal(t, TaskStateCompleted, updatedTask2.State)
	assert.NotNil(t, updatedTask2.CompletedAt)
}

func TestBulkUpdateTasksWithMultipleFields(t *testing.T) {
	ctx := context.Background()

	// Create toolset
	toolSet, err := NewToolSet()
	require.NoError(t, err)
	defer toolSet.Close()

	tools := toolSet.Tools(ctx)
	require.NotEmpty(t, tools)

	// Helper function to find tool by name
	findTool := func(name string) tool.CallableTool {
		for _, tool := range tools {
			if tool.Declaration().Name == name {
				return tool
			}
		}
		t.Fatalf("Tool %s not found", name)
		return nil
	}

	// Create project
	createProjectTool := findTool("create_project")
	projectInput := map[string]interface{}{
		"title":   "Test Project",
		"details": "A test project for unit testing",
	}
	projectInputJSON, _ := json.Marshal(projectInput)

	projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
	require.NoError(t, err)

	project := projectResult.(createProjectResult).Project
	projectID := project.ID.String()

	// Create task
	createTaskTool := findTool("create_task")
	taskInput := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Test Task",
		"description": "A test task",
		"complexity":  5,
		"priority":    5,
	}
	taskInputJSON, _ := json.Marshal(taskInput)

	taskResult, err := createTaskTool.Call(ctx, taskInputJSON)
	require.NoError(t, err)
	task := taskResult.(*Task)

	// Bulk update task with multiple fields
	bulkUpdateTool := findTool("bulk_update_tasks")
	newComplexity := 7
	newState := "in-progress"
	bulkInput := map[string]interface{}{
		"task_ids":   []string{task.ID.String()},
		"state":      newState,
		"complexity": newComplexity,
	}
	bulkInputJSON, _ := json.Marshal(bulkInput)

	_, err = bulkUpdateTool.Call(ctx, bulkInputJSON)
	require.NoError(t, err)

	// Verify task was updated
	getTaskTool := findTool("get_task")
	getTaskInput := map[string]interface{}{
		"task_id": task.ID.String(),
	}
	getTaskInputJSON, _ := json.Marshal(getTaskInput)

	getTaskResult, err := getTaskTool.Call(ctx, getTaskInputJSON)
	require.NoError(t, err)

	updatedTask := getTaskResult.(*Task)
	assert.Equal(t, TaskStateInProgress, updatedTask.State)
	assert.Equal(t, 7, updatedTask.Complexity)
	assert.Nil(t, updatedTask.CompletedAt)
}
