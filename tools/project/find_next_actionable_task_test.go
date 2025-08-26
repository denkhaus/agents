package project

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

func TestFindNextActionableTaskWithDependencies(t *testing.T) {
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

	// Scenario 1: Independent tasks (no dependencies)
	t.Run("IndependentTasks", func(t *testing.T) {
		// Create first task
		task1Input := map[string]interface{}{
			"project_id":  projectID,
			"title":       "Independent Task 1",
			"description": "First independent task",
			"complexity":  5,
		}
		task1InputJSON, _ := json.Marshal(task1Input)

		task1Result, err := createTaskTool.Call(ctx, task1InputJSON)
		require.NoError(t, err)
		task1 := task1Result.(*Task)

		// Create second task
		task2Input := map[string]interface{}{
			"project_id":  projectID,
			"title":       "Independent Task 2",
			"description": "Second independent task",
			"complexity":  3,
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
		// Either task could be returned as next actionable since both are independent
		taskIDs := []string{task1.ID.String(), task2.ID.String()}
		found := false
		for _, id := range taskIDs {
			if result.Task.ID.String() == id {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected task ID to be one of the created tasks")
	})

	// Scenario 2: Dependent tasks (one task depends on another)
	t.Run("DependentTasks", func(t *testing.T) {
		// Create a new project for this scenario
		projectInput := map[string]interface{}{
			"title":   "Dependent Tasks Project",
			"details": "A project for testing dependent tasks",
		}
		projectInputJSON, _ := json.Marshal(projectInput)

		projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
		require.NoError(t, err)

		project := projectResult.(createProjectResult).Project
		projectID := project.ID.String()

		// Create first task (dependency)
		dependencyTaskInput := map[string]interface{}{
			"project_id":  projectID,
			"title":       "Dependency Task",
			"description": "This task must be completed first",
			"complexity":  5,
		}
		dependencyTaskInputJSON, _ := json.Marshal(dependencyTaskInput)

		dependencyTaskResult, err := createTaskTool.Call(ctx, dependencyTaskInputJSON)
		require.NoError(t, err)
		dependencyTask := dependencyTaskResult.(*Task)

		// Create second task (dependent on the first)
		dependentTaskInput := map[string]interface{}{
			"project_id":  projectID,
			"title":       "Dependent Task",
			"description": "This task depends on the first task",
			"complexity":  3,
		}
		dependentTaskInputJSON, _ := json.Marshal(dependentTaskInput)

		dependentTaskResult, err := createTaskTool.Call(ctx, dependentTaskInputJSON)
		require.NoError(t, err)
		dependentTask := dependentTaskResult.(*Task)

		// Add dependency
		addDependencyTool := findTool("add_task_dependency")
		addDependencyInput := map[string]interface{}{
			"task_id":            dependentTask.ID.String(),
			"depends_on_task_id": dependencyTask.ID.String(),
		}
		addDependencyInputJSON, _ := json.Marshal(addDependencyInput)

		_, err = addDependencyTool.Call(ctx, addDependencyInputJSON)
		require.NoError(t, err)

		// Find next actionable task - should be the dependency task since it has no dependencies
		findNextActionableTaskTool := findTool("find_next_actionable_task")
		findInput := map[string]interface{}{
			"project_id": projectID,
		}
		findInputJSON, _ := json.Marshal(findInput)

		findResult, err := findNextActionableTaskTool.Call(ctx, findInputJSON)
		require.NoError(t, err)

		result := findResult.(findNextActionableTaskResult)
		assert.Equal(t, dependencyTask.ID, result.Task.ID)
		assert.Equal(t, "Dependency Task", result.Task.Title)

		// Mark the dependency task as completed
		updateTaskStateTool := findTool("update_task_state")
		updateStateInput := map[string]interface{}{
			"task_id": dependencyTask.ID.String(),
			"state":   "completed",
		}
		updateStateInputJSON, _ := json.Marshal(updateStateInput)

		_, err = updateTaskStateTool.Call(ctx, updateStateInputJSON)
		require.NoError(t, err)

		// Find next actionable task - should now be the dependent task
		findResult, err = findNextActionableTaskTool.Call(ctx, findInputJSON)
		require.NoError(t, err)

		result = findResult.(findNextActionableTaskResult)
		assert.Equal(t, dependentTask.ID, result.Task.ID)
		assert.Equal(t, "Dependent Task", result.Task.Title)
	})

	// Scenario 3: Multiple dependencies (task depends on multiple tasks)
	t.Run("MultipleDependencies", func(t *testing.T) {
		// Create a new project for this scenario
		projectInput := map[string]interface{}{
			"title":   "Multiple Dependencies Project",
			"details": "A project for testing multiple dependencies",
		}
		projectInputJSON, _ := json.Marshal(projectInput)

		projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
		require.NoError(t, err)

		project := projectResult.(createProjectResult).Project
		projectID := project.ID.String()

		// Create two dependency tasks
		dependencyTask1Input := map[string]interface{}{
			"project_id":  projectID,
			"title":       "Dependency Task 1",
			"description": "First dependency task",
			"complexity":  5,
		}
		dependencyTask1InputJSON, _ := json.Marshal(dependencyTask1Input)

		dependencyTask1Result, err := createTaskTool.Call(ctx, dependencyTask1InputJSON)
		require.NoError(t, err)
		dependencyTask1 := dependencyTask1Result.(*Task)

		dependencyTask2Input := map[string]interface{}{
			"project_id":  projectID,
			"title":       "Dependency Task 2",
			"description": "Second dependency task",
			"complexity":  3,
		}
		dependencyTask2InputJSON, _ := json.Marshal(dependencyTask2Input)

		dependencyTask2Result, err := createTaskTool.Call(ctx, dependencyTask2InputJSON)
		require.NoError(t, err)
		dependencyTask2 := dependencyTask2Result.(*Task)

		// Create dependent task that depends on both
		dependentTaskInput := map[string]interface{}{
			"project_id":  projectID,
			"title":       "Dependent Task",
			"description": "This task depends on both dependency tasks",
			"complexity":  4,
		}
		dependentTaskInputJSON, _ := json.Marshal(dependentTaskInput)

		dependentTaskResult, err := createTaskTool.Call(ctx, dependentTaskInputJSON)
		require.NoError(t, err)
		dependentTask := dependentTaskResult.(*Task)

		// Add dependencies
		addDependencyTool := findTool("add_task_dependency")

		// Add first dependency
		addDependencyInput1 := map[string]interface{}{
			"task_id":            dependentTask.ID.String(),
			"depends_on_task_id": dependencyTask1.ID.String(),
		}
		addDependencyInput1JSON, _ := json.Marshal(addDependencyInput1)

		_, err = addDependencyTool.Call(ctx, addDependencyInput1JSON)
		require.NoError(t, err)

		// Add second dependency
		addDependencyInput2 := map[string]interface{}{
			"task_id":            dependentTask.ID.String(),
			"depends_on_task_id": dependencyTask2.ID.String(),
		}
		addDependencyInput2JSON, _ := json.Marshal(addDependencyInput2)

		_, err = addDependencyTool.Call(ctx, addDependencyInput2JSON)
		require.NoError(t, err)

		// Find next actionable task - should be one of the dependency tasks
		findNextActionableTaskTool := findTool("find_next_actionable_task")
		findInput := map[string]interface{}{
			"project_id": projectID,
		}
		findInputJSON, _ := json.Marshal(findInput)

		findResult, err := findNextActionableTaskTool.Call(ctx, findInputJSON)
		require.NoError(t, err)

		result := findResult.(findNextActionableTaskResult)
		dependencyTaskIDs := []string{dependencyTask1.ID.String(), dependencyTask2.ID.String()}
		found := false
		for _, id := range dependencyTaskIDs {
			if result.Task.ID.String() == id {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected task ID to be one of the dependency tasks")

		// Mark first dependency task as completed
		updateTaskStateTool := findTool("update_task_state")
		updateStateInput := map[string]interface{}{
			"task_id": dependencyTask1.ID.String(),
			"state":   "completed",
		}
		updateStateInputJSON, _ := json.Marshal(updateStateInput)

		_, err = updateTaskStateTool.Call(ctx, updateStateInputJSON)
		require.NoError(t, err)

		// Find next actionable task - should be the second dependency task (since the first is done)
		findResult, err = findNextActionableTaskTool.Call(ctx, findInputJSON)
		require.NoError(t, err)

		result = findResult.(findNextActionableTaskResult)
		assert.Equal(t, dependencyTask2.ID, result.Task.ID)
		assert.Equal(t, "Dependency Task 2", result.Task.Title)

		// Mark second dependency task as completed
		updateStateInput = map[string]interface{}{
			"task_id": dependencyTask2.ID.String(),
			"state":   "completed",
		}
		updateStateInputJSON, _ = json.Marshal(updateStateInput)

		_, err = updateTaskStateTool.Call(ctx, updateStateInputJSON)
		require.NoError(t, err)

		// Find next actionable task - should now be the dependent task
		findResult, err = findNextActionableTaskTool.Call(ctx, findInputJSON)
		require.NoError(t, err)

		result = findResult.(findNextActionableTaskResult)
		assert.Equal(t, dependentTask.ID, result.Task.ID)
		assert.Equal(t, "Dependent Task", result.Task.Title)
	})

	// Scenario 4: Complex dependency chain
	t.Run("ComplexDependencyChain", func(t *testing.T) {
		// Create a new project for this scenario
		projectInput := map[string]interface{}{
			"title":   "Complex Dependency Chain Project",
			"details": "A project for testing complex dependency chains",
		}
		projectInputJSON, _ := json.Marshal(projectInput)

		projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
		require.NoError(t, err)

		project := projectResult.(createProjectResult).Project
		projectID := project.ID.String()

		// Create tasks in a chain: Task 1 -> Task 2 -> Task 3 -> Task 4
		// Where each task depends on the previous one
		taskInputs := []map[string]interface{}{
			{
				"project_id":  projectID,
				"title":       "Task 1",
				"description": "First task in chain",
				"complexity":  5,
			},
			{
				"project_id":  projectID,
				"title":       "Task 2",
				"description": "Second task in chain",
				"complexity":  4,
			},
			{
				"project_id":  projectID,
				"title":       "Task 3",
				"description": "Third task in chain",
				"complexity":  3,
			},
			{
				"project_id":  projectID,
				"title":       "Task 4",
				"description": "Fourth task in chain",
				"complexity":  2,
			},
		}

		tasks := make([]*Task, len(taskInputs))
		for i, taskInput := range taskInputs {
			taskInputJSON, _ := json.Marshal(taskInput)
			taskResult, err := createTaskTool.Call(ctx, taskInputJSON)
			require.NoError(t, err)
			tasks[i] = taskResult.(*Task)
		}

		// Add dependencies to create the chain
		addDependencyTool := findTool("add_task_dependency")
		for i := 1; i < len(tasks); i++ {
			addDependencyInput := map[string]interface{}{
				"task_id":            tasks[i].ID.String(),
				"depends_on_task_id": tasks[i-1].ID.String(),
			}
			addDependencyInputJSON, _ := json.Marshal(addDependencyInput)
			_, err = addDependencyTool.Call(ctx, addDependencyInputJSON)
			require.NoError(t, err)
		}

		// Find next actionable task - should be Task 1 (since it has no dependencies)
		findNextActionableTaskTool := findTool("find_next_actionable_task")
		findInput := map[string]interface{}{
			"project_id": projectID,
		}
		findInputJSON, _ := json.Marshal(findInput)

		findResult, err := findNextActionableTaskTool.Call(ctx, findInputJSON)
		require.NoError(t, err)

		result := findResult.(findNextActionableTaskResult)
		assert.Equal(t, tasks[0].ID, result.Task.ID)
		assert.Equal(t, "Task 1", result.Task.Title)

		// Mark tasks as completed one by one and verify the next actionable task
		updateTaskStateTool := findTool("update_task_state")
		for i := 0; i < len(tasks)-1; i++ {
			// Mark current task as completed
			updateStateInput := map[string]interface{}{
				"task_id": tasks[i].ID.String(),
				"state":   "completed",
			}
			updateStateInputJSON, _ := json.Marshal(updateStateInput)
			_, err = updateTaskStateTool.Call(ctx, updateStateInputJSON)
			require.NoError(t, err)

			// Find next actionable task - should be the next task in the chain
			findResult, err = findNextActionableTaskTool.Call(ctx, findInputJSON)
			require.NoError(t, err)

			result = findResult.(findNextActionableTaskResult)
			assert.Equal(t, tasks[i+1].ID, result.Task.ID)
			assert.Equal(t, taskInputs[i+1]["title"], result.Task.Title)
		}
	})

	// Scenario 5: In-progress tasks with dependencies
	t.Run("InProgressTasksWithDependencies", func(t *testing.T) {
		// Create a new project for this scenario
		projectInput := map[string]interface{}{
			"title":   "In Progress Tasks Project",
			"details": "A project for testing in-progress tasks with dependencies",
		}
		projectInputJSON, _ := json.Marshal(projectInput)

		projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
		require.NoError(t, err)

		project := projectResult.(createProjectResult).Project
		projectID := project.ID.String()

		// Create dependency task
		dependencyTaskInput := map[string]interface{}{
			"project_id":  projectID,
			"title":       "Dependency Task",
			"description": "This task must be completed first",
			"complexity":  5,
		}
		dependencyTaskInputJSON, _ := json.Marshal(dependencyTaskInput)

		dependencyTaskResult, err := createTaskTool.Call(ctx, dependencyTaskInputJSON)
		require.NoError(t, err)
		dependencyTask := dependencyTaskResult.(*Task)

		// Create dependent task
		dependentTaskInput := map[string]interface{}{
			"project_id":  projectID,
			"title":       "Dependent Task",
			"description": "This task depends on the first task",
			"complexity":  3,
		}
		dependentTaskInputJSON, _ := json.Marshal(dependentTaskInput)

		dependentTaskResult, err := createTaskTool.Call(ctx, dependentTaskInputJSON)
		require.NoError(t, err)
		dependentTask := dependentTaskResult.(*Task)

		// Add dependency
		addDependencyTool := findTool("add_task_dependency")
		addDependencyInput := map[string]interface{}{
			"task_id":            dependentTask.ID.String(),
			"depends_on_task_id": dependencyTask.ID.String(),
		}
		addDependencyInputJSON, _ := json.Marshal(addDependencyInput)

		_, err = addDependencyTool.Call(ctx, addDependencyInputJSON)
		require.NoError(t, err)

		// Mark dependency task as completed
		updateTaskStateTool := findTool("update_task_state")
		updateStateInput := map[string]interface{}{
			"task_id": dependencyTask.ID.String(),
			"state":   "completed",
		}
		updateStateInputJSON, _ := json.Marshal(updateStateInput)

		_, err = updateTaskStateTool.Call(ctx, updateStateInputJSON)
		require.NoError(t, err)

		// Mark dependent task as in-progress
		updateStateInput = map[string]interface{}{
			"task_id": dependentTask.ID.String(),
			"state":   "in-progress",
		}
		updateStateInputJSON, _ = json.Marshal(updateStateInput)

		_, err = updateTaskStateTool.Call(ctx, updateStateInputJSON)
		require.NoError(t, err)

		// Find next actionable task - should be the in-progress task since it has its dependencies met
		findNextActionableTaskTool := findTool("find_next_actionable_task")
		findInput := map[string]interface{}{
			"project_id": projectID,
		}
		findInputJSON, _ := json.Marshal(findInput)

		findResult, err := findNextActionableTaskTool.Call(ctx, findInputJSON)
		require.NoError(t, err)

		result := findResult.(findNextActionableTaskResult)
		assert.Equal(t, dependentTask.ID, result.Task.ID)
		assert.Equal(t, "Dependent Task", result.Task.Title)
	})

	// Scenario 6: Deadlock scenario (circular dependencies are prevented by the system)
	// This is more of a validation test to ensure the system prevents circular dependencies
	t.Run("CircularDependencyPrevention", func(t *testing.T) {
		// Create a new project for this scenario
		projectInput := map[string]interface{}{
			"title":   "Circular Dependency Prevention Project",
			"details": "A project for testing circular dependency prevention",
		}
		projectInputJSON, _ := json.Marshal(projectInput)

		projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
		require.NoError(t, err)

		project := projectResult.(createProjectResult).Project
		projectID := project.ID.String()

		// Create two tasks
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

		// Try to create a circular dependency: Task 1 depends on Task 2
		addDependencyTool := findTool("add_task_dependency")
		addDependencyInput1 := map[string]interface{}{
			"task_id":            task1.ID.String(),
			"depends_on_task_id": task2.ID.String(),
		}
		addDependencyInput1JSON, _ := json.Marshal(addDependencyInput1)

		_, err = addDependencyTool.Call(ctx, addDependencyInput1JSON)
		require.NoError(t, err)

		// Try to create the reverse dependency: Task 2 depends on Task 1
		// This should fail
		addDependencyInput2 := map[string]interface{}{
			"task_id":            task2.ID.String(),
			"depends_on_task_id": task1.ID.String(),
		}
		addDependencyInput2JSON, _ := json.Marshal(addDependencyInput2)

		_, err = addDependencyTool.Call(ctx, addDependencyInput2JSON)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circular dependency")
	})

	// Scenario 7: Tasks with dependencies
	t.Run("TasksWithDependencies", func(t *testing.T) {
		// Create a new project for this scenario
		projectInput := map[string]interface{}{
			"title":   "Tasks With Dependencies Project",
			"details": "A project for testing tasks with dependencies",
		}
		projectInputJSON, _ := json.Marshal(projectInput)

		projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
		require.NoError(t, err)

		project := projectResult.(createProjectResult).Project
		projectID := project.ID.String()

		// Create two tasks
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

		// Add dependencies: Task 2 depends on Task 1
		addDependencyTool := findTool("add_task_dependency")
		addDependencyInput := map[string]interface{}{
			"task_id":           task2.ID.String(),
			"depends_on_task_id": task1.ID.String(),
		}
		addDependencyInputJSON, _ := json.Marshal(addDependencyInput)

		_, err = addDependencyTool.Call(ctx, addDependencyInputJSON)
		require.NoError(t, err)

		// Both tasks are in pending state, but Task 2 depends on Task 1
		// Find next actionable task - should return Task 1 since it has no dependencies
		findNextActionableTaskTool := findTool("find_next_actionable_task")
		findInput := map[string]interface{}{
			"project_id": projectID,
		}
		findInputJSON, _ := json.Marshal(findInput)

		findResult, err := findNextActionableTaskTool.Call(ctx, findInputJSON)
		require.NoError(t, err)

		result := findResult.(findNextActionableTaskResult)
		assert.Equal(t, task1.ID, result.Task.ID)
		assert.Equal(t, "Task 1", result.Task.Title)

		// Mark Task 1 as completed
		updateTaskStateTool := findTool("update_task_state")
		updateStateInput := map[string]interface{}{
			"task_id": task1.ID.String(),
			"state":   "completed",
		}
		updateStateInputJSON, _ := json.Marshal(updateStateInput)

		_, err = updateTaskStateTool.Call(ctx, updateStateInputJSON)
		require.NoError(t, err)

		// Now Task 2 should be actionable since its dependency is completed
		findResult, err = findNextActionableTaskTool.Call(ctx, findInputJSON)
		require.NoError(t, err)

		result = findResult.(findNextActionableTaskResult)
		assert.Equal(t, task2.ID, result.Task.ID)
		assert.Equal(t, "Task 2", result.Task.Title)
	})
}
