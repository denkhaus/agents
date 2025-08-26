package projecttasks

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func TestProjectTaskToolSet(t *testing.T) {
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

	// Test project creation
	createProjectTool := findTool("create_project")
	projectInput := map[string]interface{}{
		"title":       "Test Project",
		"description": "A test project for unit testing",
	}
	projectInputJSON, _ := json.Marshal(projectInput)

	projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
	require.NoError(t, err)

	project := projectResult.(*Project)
	assert.Equal(t, "Test Project", project.Title)
	assert.Equal(t, "A test project for unit testing", project.Description)
	assert.NotEqual(t, uuid.Nil, project.ID)

	projectID := project.ID.String()

	// Test task creation
	createTaskTool := findTool("create_task")
	taskInput := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Root Task",
		"description": "A root level task",
		"complexity":  5,
		"priority":    8,
	}
	taskInputJSON, _ := json.Marshal(taskInput)

	taskResult, err := createTaskTool.Call(ctx, taskInputJSON)
	require.NoError(t, err)

	rootTask := taskResult.(*Task)
	assert.Equal(t, "Root Task", rootTask.Title)
	assert.Equal(t, 5, rootTask.Complexity)
	assert.Equal(t, 8, rootTask.Priority)
	assert.Equal(t, 0, rootTask.Depth)
	assert.Nil(t, rootTask.ParentID)

	// Test subtask creation
	subtaskInput := map[string]interface{}{
		"project_id":  projectID,
		"parent_id":   rootTask.ID.String(),
		"title":       "Subtask 1",
		"description": "A subtask",
		"complexity":  3,
		"priority":    6,
	}
	subtaskInputJSON, _ := json.Marshal(subtaskInput)

	subtaskResult, err := createTaskTool.Call(ctx, subtaskInputJSON)
	require.NoError(t, err)

	subtask := subtaskResult.(*Task)
	assert.Equal(t, "Subtask 1", subtask.Title)
	assert.Equal(t, 1, subtask.Depth)
	assert.Equal(t, rootTask.ID, *subtask.ParentID)

	// Test hierarchical listing
	listHierarchicalTool := findTool("list_tasks_hierarchical")
	listInput := map[string]interface{}{
		"project_id": projectID,
	}
	listInputJSON, _ := json.Marshal(listInput)

	listResult, err := listHierarchicalTool.Call(ctx, listInputJSON)
	require.NoError(t, err)

	listData := listResult.(map[string]interface{})
	hierarchy := listData["hierarchy"].([]*TaskHierarchy)
	assert.Len(t, hierarchy, 1)
	assert.Equal(t, "Root Task", hierarchy[0].Task.Title)
	assert.Len(t, hierarchy[0].Children, 1)
	assert.Equal(t, "Subtask 1", hierarchy[0].Children[0].Task.Title)

	// Test task state update
	updateStateTool := findTool("update_task_state")
	stateInput := map[string]interface{}{
		"task_id": subtask.ID.String(),
		"state":   TaskStateCompleted,
	}
	stateInputJSON, _ := json.Marshal(stateInput)

	stateResult, err := updateStateTool.Call(ctx, stateInputJSON)
	require.NoError(t, err)

	updatedTask := stateResult.(*Task)
	assert.Equal(t, TaskStateCompleted, updatedTask.State)
	assert.NotNil(t, updatedTask.CompletedAt)

	// Test project progress
	progressTool := findTool("get_project_progress")
	progressInput := map[string]interface{}{
		"project_id": projectID,
	}
	progressInputJSON, _ := json.Marshal(progressInput)

	progressResult, err := progressTool.Call(ctx, progressInputJSON)
	require.NoError(t, err)

	progress := progressResult.(*ProjectProgress)
	assert.Equal(t, 2, progress.TotalTasks)
	assert.Equal(t, 1, progress.CompletedTasks)
	assert.Equal(t, 50.0, progress.OverallProgress)

	// Test finding next actionable task
	nextTaskTool := findTool("find_next_actionable_task")
	nextTaskResult, err := nextTaskTool.Call(ctx, progressInputJSON)
	require.NoError(t, err)

	nextTask := nextTaskResult.(*Task)
	assert.Equal(t, rootTask.ID, nextTask.ID) // Root task should be next (pending, higher priority)

	// Test tasks needing breakdown
	breakdownTool := findTool("find_tasks_needing_breakdown")
	breakdownResult, err := breakdownTool.Call(ctx, progressInputJSON)
	require.NoError(t, err)

	breakdownData := breakdownResult.(map[string]interface{})
	breakdownTasks := breakdownData["tasks"].([]*Task)
	assert.Empty(t, breakdownTasks) // No tasks with complexity >= 8

	// Create high complexity task
	highComplexityInput := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Complex Task",
		"description": "A very complex task",
		"complexity":  9,
		"priority":    5,
	}
	highComplexityInputJSON, _ := json.Marshal(highComplexityInput)

	_, err = createTaskTool.Call(ctx, highComplexityInputJSON)
	require.NoError(t, err)

	// Test breakdown again
	breakdownResult, err = breakdownTool.Call(ctx, progressInputJSON)
	require.NoError(t, err)

	breakdownData = breakdownResult.(map[string]interface{})
	breakdownTasks = breakdownData["tasks"].([]*Task)
	assert.Len(t, breakdownTasks, 1)
	assert.Equal(t, "Complex Task", breakdownTasks[0].Title)

	// Test task deletion
	deleteTaskTool := findTool("delete_task")
	deleteInput := map[string]interface{}{
		"task_id": subtask.ID.String(),
	}
	deleteInputJSON, _ := json.Marshal(deleteInput)

	deleteResult, err := deleteTaskTool.Call(ctx, deleteInputJSON)
	require.NoError(t, err)

	deleteData := deleteResult.(map[string]interface{})
	assert.True(t, deleteData["success"].(bool))

	// Test subtree deletion
	deleteSubtreeTool := findTool("delete_task_subtree")
	deleteSubtreeInput := map[string]interface{}{
		"task_id": rootTask.ID.String(),
	}
	deleteSubtreeInputJSON, _ := json.Marshal(deleteSubtreeInput)

	_, err = deleteSubtreeTool.Call(ctx, deleteSubtreeInputJSON)
	require.NoError(t, err)

	// Verify hierarchy is now empty except for the complex task
	listResult, err = listHierarchicalTool.Call(ctx, listInputJSON)
	require.NoError(t, err)

	listData = listResult.(map[string]interface{})
	hierarchy = listData["hierarchy"].([]*TaskHierarchy)
	assert.Len(t, hierarchy, 1)
	assert.Equal(t, "Complex Task", hierarchy[0].Task.Title)
}

func TestProjectTaskToolSetValidation(t *testing.T) {
	ctx := context.Background()

	toolSet, err := NewToolSet()
	require.NoError(t, err)
	defer toolSet.Close()

	tools := toolSet.Tools(ctx)
	
	// Find the create_project tool
	var createProjectTool tool.CallableTool
	for _, t := range tools {
		if t.Declaration().Name == "create_project" {
			createProjectTool = t
			break
		}
	}

	// Test empty title validation
	invalidInput := map[string]interface{}{
		"title":       "",
		"description": "Test",
	}
	invalidInputJSON, _ := json.Marshal(invalidInput)

	_, err = createProjectTool.Call(ctx, invalidInputJSON)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title cannot be empty")

	// Test title too long
	longTitle := make([]byte, 201)
	for i := range longTitle {
		longTitle[i] = 'a'
	}

	invalidInput = map[string]interface{}{
		"title":       string(longTitle),
		"description": "Test",
	}
	invalidInputJSON, _ = json.Marshal(invalidInput)

	_, err = createProjectTool.Call(ctx, invalidInputJSON)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title cannot exceed 200 characters")
}

func TestProjectTaskToolSetConcurrency(t *testing.T) {
	ctx := context.Background()

	toolSet, err := NewToolSet()
	require.NoError(t, err)
	defer toolSet.Close()

	tools := toolSet.Tools(ctx)
	
	// Find tools by name
	var createProjectTool, createTaskTool tool.CallableTool
	for _, t := range tools {
		switch t.Declaration().Name {
		case "create_project":
			createProjectTool = t
		case "create_task":
			createTaskTool = t
		}
	}

	// Create project
	projectInput := map[string]interface{}{
		"title":       "Concurrent Test Project",
		"description": "Testing concurrent operations",
	}
	projectInputJSON, _ := json.Marshal(projectInput)

	projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
	require.NoError(t, err)

	project := projectResult.(*Project)
	projectID := project.ID.String()

	// Create multiple tasks concurrently
	const numTasks = 10
	results := make(chan error, numTasks)

	for i := 0; i < numTasks; i++ {
		go func(taskNum int) {
			taskInput := map[string]interface{}{
				"project_id":  projectID,
				"title":       fmt.Sprintf("Concurrent Task %d", taskNum),
				"description": fmt.Sprintf("Task created concurrently %d", taskNum),
				"complexity":  5,
				"priority":    5,
			}
			taskInputJSON, _ := json.Marshal(taskInput)

			_, err := createTaskTool.Call(ctx, taskInputJSON)
			results <- err
		}(i)
	}

	// Wait for all tasks to complete
	for i := 0; i < numTasks; i++ {
		err := <-results
		assert.NoError(t, err)
	}

	// Verify all tasks were created
	listTool := tools[11] // Assuming list_tasks_hierarchical
	listInput := map[string]interface{}{
		"project_id": projectID,
	}
	listInputJSON, _ := json.Marshal(listInput)

	listResult, err := listTool.Call(ctx, listInputJSON)
	require.NoError(t, err)

	listData := listResult.(map[string]interface{})
	hierarchy := listData["hierarchy"].([]*TaskHierarchy)
	assert.Len(t, hierarchy, numTasks)
}

func TestProjectTaskToolSetDepthLimits(t *testing.T) {
	ctx := context.Background()

	// Create toolset with custom config
	config := &Config{
		MaxTasksPerDepth: map[int]int{
			0: 2, // Only 2 root tasks allowed
			1: 3, // Only 3 tasks at depth 1
		},
		ComplexityThreshold: 8,
		MaxDepth:           2,
		DefaultPriority:    5,
	}

	toolSet, err := NewToolSet(WithConfig(config))
	require.NoError(t, err)
	defer toolSet.Close()

	tools := toolSet.Tools(ctx)
	createProjectTool := tools[0]
	createTaskTool := tools[5]

	// Create project
	projectInput := map[string]interface{}{
		"title":       "Depth Limit Test",
		"description": "Testing depth limits",
	}
	projectInputJSON, _ := json.Marshal(projectInput)

	projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
	require.NoError(t, err)

	project := projectResult.(*Project)
	projectID := project.ID.String()

	// Create 2 root tasks (should succeed)
	for i := 0; i < 2; i++ {
		taskInput := map[string]interface{}{
			"project_id":  projectID,
			"title":       fmt.Sprintf("Root Task %d", i+1),
			"description": "Root task",
			"complexity":  5,
			"priority":    5,
		}
		taskInputJSON, _ := json.Marshal(taskInput)

		_, err := createTaskTool.Call(ctx, taskInputJSON)
		require.NoError(t, err)
	}

	// Try to create 3rd root task (should fail)
	taskInput := map[string]interface{}{
		"project_id":  projectID,
		"title":       "Root Task 3",
		"description": "This should fail",
		"complexity":  5,
		"priority":    5,
	}
	taskInputJSON, _ := json.Marshal(taskInput)

	_, err = createTaskTool.Call(ctx, taskInputJSON)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum tasks per depth")
}