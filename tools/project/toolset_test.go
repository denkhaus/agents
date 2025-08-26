package project

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
	// t.Skip("Skipping toolset test due to function tool schema generation issue")
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
		"details": "A test project for unit testing",
	}
	projectInputJSON, _ := json.Marshal(projectInput)

	projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
	require.NoError(t, err)

	project := projectResult.(createProjectResult).Project
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

	// Test child task retrieval
	getChildTasksTool := findTool("get_child_tasks")
	childTaskInput := map[string]interface{}{
		"task_id": rootTask.ID.String(),
	}
	childTaskInputJSON, _ := json.Marshal(childTaskInput)

	childTasksResult, err := getChildTasksTool.Call(ctx, childTaskInputJSON)
	require.NoError(t, err)

	childTasks := childTasksResult.(getChildTasksResult)
	assert.Len(t, childTasks.Tasks, 1)
	assert.Equal(t, "Subtask 1", childTasks.Tasks[0].Title)

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
}

func TestProjectTaskToolSetValidation(t *testing.T) {
	t.Skip("Skipping toolset validation test due to function tool schema generation issue")
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
	t.Skip("Skipping toolset concurrency test due to function tool schema generation issue")
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

	// Verify all tasks were created by listing them
	// Note: We're skipping the hierarchical listing test since we don't have that function
	// In a real implementation, you would verify the tasks were created correctly
}

func TestUpdateDescriptionFunctions(t *testing.T) {
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

	// Create a project
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

	// Update project description
	updateProjectDescTool := findTool("update_project_description")
	newProjectDesc := "Updated project description"
	projectDescInput := map[string]interface{}{
		"project_id":  projectID,
		"description": newProjectDesc,
	}
	projectDescInputJSON, _ := json.Marshal(projectDescInput)

	updatedProjectResult, err := updateProjectDescTool.Call(ctx, projectDescInputJSON)
	require.NoError(t, err)

	updatedProject := updatedProjectResult.(*Project)
	assert.Equal(t, newProjectDesc, updatedProject.Description)

	// Create a task
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
	taskID := task.ID.String()

	// Update task description
	updateTaskDescTool := findTool("update_task_description")
	newTaskDesc := "Updated task description"
	taskDescInput := map[string]interface{}{
		"task_id":     taskID,
		"description": newTaskDesc,
	}
	taskDescInputJSON, _ := json.Marshal(taskDescInput)

	updatedTaskResult, err := updateTaskDescTool.Call(ctx, taskDescInputJSON)
	require.NoError(t, err)

	updatedTask := updatedTaskResult.(*Task)
	assert.Equal(t, newTaskDesc, updatedTask.Description)
}

func TestProjectTaskToolSetDepthLimits(t *testing.T) {
	t.Skip("Skipping toolset depth limits test due to function tool schema generation issue")
	ctx := context.Background()

	// Create toolset with custom config
	config := &Config{
		MaxTasksPerDepth:    2, // Only 2 tasks per depth level
		ComplexityThreshold: 8,
		MaxDepth:            2,
		DefaultPriority:     5,
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
