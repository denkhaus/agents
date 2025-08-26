package project

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func TestSetTaskEstimate(t *testing.T) {
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
	}
	taskInputJSON, _ := json.Marshal(taskInput)

	taskResult, err := createTaskTool.Call(ctx, taskInputJSON)
	require.NoError(t, err)

	task := taskResult.(*Task)

	// Set task estimate
	setEstimateTool := findTool("set_task_estimate")
	estimate := int64(120) // 2 hours
	setEstimateInput := map[string]interface{}{
		"task_id":  task.ID.String(),
		"estimate": estimate,
	}
	setEstimateInputJSON, _ := json.Marshal(setEstimateInput)

	setEstimateResult, err := setEstimateTool.Call(ctx, setEstimateInputJSON)
	require.NoError(t, err)

	result := setEstimateResult.(setTaskEstimateResult)
	updatedTask := result.Task

	// Verify the estimate was set
	assert.Equal(t, estimate, *updatedTask.Estimate)
	assert.WithinDuration(t, task.UpdatedAt, updatedTask.UpdatedAt, 1000000000) // Within 1 second

	// Verify the task can be retrieved with the estimate
	getTaskTool := findTool("get_task")
	getTaskInput := map[string]interface{}{
		"task_id": task.ID.String(),
	}
	getTaskInputJSON, _ := json.Marshal(getTaskInput)

	getTaskResult, err := getTaskTool.Call(ctx, getTaskInputJSON)
	require.NoError(t, err)

	retrievedTask := getTaskResult.(*Task)
	assert.Equal(t, estimate, *retrievedTask.Estimate)
}

func TestSetTaskEstimateToZero(t *testing.T) {
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

	// Set task estimate to zero
	setEstimateTool := findTool("set_task_estimate")
	estimate := int64(0) // 0 minutes
	setEstimateInput := map[string]interface{}{
		"task_id":  task.ID.String(),
		"estimate": estimate,
	}
	setEstimateInputJSON, _ := json.Marshal(setEstimateInput)

	setEstimateResult, err := setEstimateTool.Call(ctx, setEstimateInputJSON)
	require.NoError(t, err)

	result := setEstimateResult.(setTaskEstimateResult)
	updatedTask := result.Task

	// Verify the estimate was set to zero
	assert.Equal(t, estimate, *updatedTask.Estimate)
}