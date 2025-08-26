package project

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func TestListTasksForProject(t *testing.T) {
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

	// Create root task
	createTaskTool := findTool("create_task")
	rootTaskInput := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Root Task",
		"description": "A root task",
		"complexity":  5,
		"priority":    5,
	}
	rootTaskInputJSON, _ := json.Marshal(rootTaskInput)

	rootTaskResult, err := createTaskTool.Call(ctx, rootTaskInputJSON)
	require.NoError(t, err)
	rootTask := rootTaskResult.(*Task)

	// Create subtask
	subTaskInput := map[string]interface{}{
		"project_id":  projectID,
		"parent_id":   rootTask.ID.String(),
		"title":       "Sub Task",
		"description": "A sub task",
		"complexity":  3,
		"priority":    3,
	}
	subTaskInputJSON, _ := json.Marshal(subTaskInput)

	_, err = createTaskTool.Call(ctx, subTaskInputJSON)
	require.NoError(t, err)

	// List all tasks for project
	listTasksTool := findTool("list_tasks_for_project")
	listInput := map[string]interface{}{
		"project_id": projectID,
	}
	listInputJSON, _ := json.Marshal(listInput)

	listResult, err := listTasksTool.Call(ctx, listInputJSON)
	require.NoError(t, err)

	result := listResult.(listTasksForProjectResult)
	assert.Equal(t, 2, result.Count)
	
	// Verify that both tasks are in the result (order may vary)
	taskTitles := make(map[string]bool)
	for _, task := range result.Tasks {
		taskTitles[task.Title] = true
	}
	assert.True(t, taskTitles["Root Task"])
	assert.True(t, taskTitles["Sub Task"])
}