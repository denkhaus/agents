package project

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func TestDuplicateTask(t *testing.T) {
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

	// Create source project
	createProjectTool := findTool("create_project")
	sourceProjectInput := map[string]interface{}{
		"title":   "Source Project",
		"details": "A source project for duplication",
	}
	sourceProjectInputJSON, _ := json.Marshal(sourceProjectInput)

	sourceProjectResult, err := createProjectTool.Call(ctx, sourceProjectInputJSON)
	require.NoError(t, err)

	sourceProject := sourceProjectResult.(createProjectResult).Project
	sourceProjectID := sourceProject.ID.String()

	// Create target project
	targetProjectInput := map[string]interface{}{
		"title":   "Target Project",
		"details": "A target project for duplication",
	}
	targetProjectInputJSON, _ := json.Marshal(targetProjectInput)

	targetProjectResult, err := createProjectTool.Call(ctx, targetProjectInputJSON)
	require.NoError(t, err)

	targetProject := targetProjectResult.(createProjectResult).Project
	targetProjectID := targetProject.ID.String()

	// Create task in source project
	createTaskTool := findTool("create_task")
	taskInput := map[string]interface{}{
		"project_id":  sourceProjectID,
		"title":       "Original Task",
		"description": "This is the original task",
		"complexity":  7,
		"priority":    8,
	}
	taskInputJSON, _ := json.Marshal(taskInput)

	taskResult, err := createTaskTool.Call(ctx, taskInputJSON)
	require.NoError(t, err)

	originalTask := taskResult.(*Task)

	// Duplicate the task to the target project
	duplicateTaskTool := findTool("duplicate_task")
	duplicateInput := map[string]interface{}{
		"task_id":       originalTask.ID.String(),
		"new_project_id": targetProjectID,
	}
	duplicateInputJSON, _ := json.Marshal(duplicateInput)

	duplicateResult, err := duplicateTaskTool.Call(ctx, duplicateInputJSON)
	require.NoError(t, err)

	result := duplicateResult.(duplicateTaskResult)
	duplicatedTask := result.Task

	// Verify the duplicated task
	assert.NotEqual(t, originalTask.ID, duplicatedTask.ID)
	assert.Equal(t, targetProjectID, duplicatedTask.ProjectID.String())
	assert.Equal(t, originalTask.Title, duplicatedTask.Title)
	assert.Equal(t, originalTask.Description, duplicatedTask.Description)
	assert.Equal(t, originalTask.Complexity, duplicatedTask.Complexity)
	assert.Equal(t, originalTask.Priority, duplicatedTask.Priority)
	assert.Equal(t, TaskStatePending, duplicatedTask.State)
	assert.Equal(t, 0, duplicatedTask.Depth)
	assert.Nil(t, duplicatedTask.ParentID)
	assert.Nil(t, duplicatedTask.CompletedAt)
	assert.WithinDuration(t, originalTask.CreatedAt, duplicatedTask.CreatedAt, 1000000000) // Within 1 second

	// Verify the original task is unchanged
	getTaskTool := findTool("get_task")
	getTaskInput := map[string]interface{}{
		"task_id": originalTask.ID.String(),
	}
	getTaskInputJSON, _ := json.Marshal(getTaskInput)

	getTaskResult, err := getTaskTool.Call(ctx, getTaskInputJSON)
	require.NoError(t, err)

	stillOriginalTask := getTaskResult.(*Task)
	assert.Equal(t, originalTask.ID, stillOriginalTask.ID)
	assert.Equal(t, originalTask.ProjectID, stillOriginalTask.ProjectID)
	assert.Equal(t, originalTask.State, stillOriginalTask.State)
}