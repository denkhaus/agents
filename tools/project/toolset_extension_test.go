package project

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func TestFindNextActionableTask(t *testing.T) {
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
		"priority":    8,
	}
	task1InputJSON, _ := json.Marshal(task1Input)

	_, err = createTaskTool.Call(ctx, task1InputJSON)
	require.NoError(t, err)

	// Create second task with higher priority
	task2Input := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Task 2",
		"description": "Second task",
		"complexity":  3,
		"priority":    9,
	}
	task2InputJSON, _ := json.Marshal(task2Input)

	task2Result, err := createTaskTool.Call(ctx, task2InputJSON)
	require.NoError(t, err)
	task2 := task2Result.(*Task)

	// Find next actionable task
	findNextActionableTaskTool := findTool("find_next_actionable_task")
	findInput := map[string]interface{}{
		"project_id": projectID,
	}
	findInputJSON, _ := json.Marshal(findInput)

	findResult, err := findNextActionableTaskTool.Call(ctx, findInputJSON)
	require.NoError(t, err)

	result := findResult.(findNextActionableTaskResult)
	assert.Equal(t, task2.ID, result.Task.ID)
	assert.Equal(t, "Task 2", result.Task.Title)
}

func TestFindTasksNeedingBreakdown(t *testing.T) {
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

	// Create a complex task that needs breakdown
	createTaskTool := findTool("create_task")
	complexTaskInput := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Complex Task",
		"description": "A complex task that needs to be broken down",
		"complexity":  9, // High complexity
		"priority":    5,
	}
	complexTaskInputJSON, _ := json.Marshal(complexTaskInput)

	complexTaskResult, err := createTaskTool.Call(ctx, complexTaskInputJSON)
	require.NoError(t, err)
	complexTask := complexTaskResult.(*Task)

	// Create a simple task that doesn't need breakdown
	simpleTaskInput := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Simple Task",
		"description": "A simple task",
		"complexity":  3, // Low complexity
		"priority":    5,
	}
	simpleTaskInputJSON, _ := json.Marshal(simpleTaskInput)

	_, err = createTaskTool.Call(ctx, simpleTaskInputJSON)
	require.NoError(t, err)

	// Find tasks needing breakdown
	findTasksNeedingBreakdownTool := findTool("find_tasks_needing_breakdown")
	findInput := map[string]interface{}{
		"project_id": projectID,
	}
	findInputJSON, _ := json.Marshal(findInput)

	findResult, err := findTasksNeedingBreakdownTool.Call(ctx, findInputJSON)
	require.NoError(t, err)

	result := findResult.(findTasksNeedingBreakdownResult)
	assert.Equal(t, 1, result.Count)
	assert.Equal(t, complexTask.ID, result.Tasks[0].ID)
	assert.Equal(t, "Complex Task", result.Tasks[0].Title)
}

func TestGetRootTasks(t *testing.T) {
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

	// Get root tasks
	getRootTasksTool := findTool("get_root_tasks")
	getInput := map[string]interface{}{
		"project_id": projectID,
	}
	getInputJSON, _ := json.Marshal(getInput)

	getResult, err := getRootTasksTool.Call(ctx, getInputJSON)
	require.NoError(t, err)

	result := getResult.(getRootTasksResult)
	assert.Equal(t, 1, result.Count)
	assert.Equal(t, rootTask.ID, result.Tasks[0].ID)
	assert.Equal(t, "Root Task", result.Tasks[0].Title)
}

func TestListTasksByState(t *testing.T) {
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

	// Update task state to completed
	updateTaskStateTool := findTool("update_task_state")
	stateInput := map[string]interface{}{
		"task_id": task.ID.String(),
		"state":   TaskStateCompleted,
	}
	stateInputJSON, _ := json.Marshal(stateInput)

	_, err = updateTaskStateTool.Call(ctx, stateInputJSON)
	require.NoError(t, err)

	// List tasks by state
	listTasksByStateTool := findTool("list_tasks_by_state")
	listInput := map[string]interface{}{
		"project_id": projectID,
		"state":      TaskStateCompleted,
	}
	listInputJSON, _ := json.Marshal(listInput)

	listResult, err := listTasksByStateTool.Call(ctx, listInputJSON)
	require.NoError(t, err)

	result := listResult.(listTasksByStateResult)
	assert.Equal(t, 1, result.Count)
	assert.Equal(t, task.ID, result.Tasks[0].ID)
	assert.Equal(t, TaskStateCompleted, result.Tasks[0].State)
}