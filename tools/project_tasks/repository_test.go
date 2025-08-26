package projecttasks

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepository(t *testing.T) {
	ctx := context.Background()
	repo := newMemoryRepository()

	t.Run("Project Operations", func(t *testing.T) {
		// Test project creation
		project := &Project{
			ID:          uuid.New(),
			Title:       "Test Project",
			Description: "A test project",
		}

		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)

		// Test duplicate project creation
		err = repo.CreateProject(ctx, project)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")

		// Test project retrieval
		retrieved, err := repo.GetProject(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, project.Title, retrieved.Title)
		assert.Equal(t, project.Description, retrieved.Description)
		assert.False(t, retrieved.CreatedAt.IsZero())
		assert.False(t, retrieved.UpdatedAt.IsZero())

		// Test project update
		retrieved.Title = "Updated Project"
		retrieved.Description = "Updated description"
		err = repo.UpdateProject(ctx, retrieved)
		require.NoError(t, err)

		updated, err := repo.GetProject(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Project", updated.Title)
		assert.Equal(t, "Updated description", updated.Description)
		assert.True(t, updated.UpdatedAt.After(updated.CreatedAt))

		// Test project listing
		projects, err := repo.ListProjects(ctx)
		require.NoError(t, err)
		assert.Len(t, projects, 1)
		assert.Equal(t, updated.ID, projects[0].ID)

		// Test project deletion
		err = repo.DeleteProject(ctx, project.ID)
		require.NoError(t, err)

		_, err = repo.GetProject(ctx, project.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Task Operations", func(t *testing.T) {
		// Create a project first
		project := &Project{
			ID:          uuid.New(),
			Title:       "Task Test Project",
			Description: "For testing tasks",
		}
		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)

		// Test root task creation
		rootTask := &Task{
			ID:          uuid.New(),
			ProjectID:   project.ID,
			Title:       "Root Task",
			Description: "A root task",
			State:       TaskStatePending,
			Complexity:  5,
			Priority:    8,
		}

		err = repo.CreateTask(ctx, rootTask)
		require.NoError(t, err)

		// Test task retrieval
		retrieved, err := repo.GetTask(ctx, rootTask.ID)
		require.NoError(t, err)
		assert.Equal(t, rootTask.Title, retrieved.Title)
		assert.Equal(t, 0, retrieved.Depth)
		assert.Nil(t, retrieved.ParentID)
		assert.False(t, retrieved.CreatedAt.IsZero())

		// Test subtask creation
		subtask := &Task{
			ID:          uuid.New(),
			ProjectID:   project.ID,
			ParentID:    &rootTask.ID,
			Title:       "Subtask",
			Description: "A subtask",
			State:       TaskStatePending,
			Complexity:  3,
			Priority:    6,
		}

		err = repo.CreateTask(ctx, subtask)
		require.NoError(t, err)

		retrievedSubtask, err := repo.GetTask(ctx, subtask.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, retrievedSubtask.Depth)
		assert.Equal(t, rootTask.ID, *retrievedSubtask.ParentID)

		// Test task update
		retrievedSubtask.State = TaskStateCompleted
		err = repo.UpdateTask(ctx, retrievedSubtask)
		require.NoError(t, err)

		updated, err := repo.GetTask(ctx, subtask.ID)
		require.NoError(t, err)
		assert.Equal(t, TaskStateCompleted, updated.State)
		assert.NotNil(t, updated.CompletedAt)

		// Test task queries
		projectTasks, err := repo.GetTasksByProject(ctx, project.ID)
		require.NoError(t, err)
		assert.Len(t, projectTasks, 2)

		rootTasks, err := repo.GetRootTasks(ctx, project.ID)
		require.NoError(t, err)
		assert.Len(t, rootTasks, 1)
		assert.Equal(t, rootTask.ID, rootTasks[0].ID)

		parentTasks, err := repo.GetTasksByParent(ctx, rootTask.ID)
		require.NoError(t, err)
		assert.Len(t, parentTasks, 1)
		assert.Equal(t, subtask.ID, parentTasks[0].ID)

		// Test task filtering
		completedState := TaskStateCompleted
		completedTasks, err := repo.ListTasks(ctx, TaskFilter{
			ProjectID: &project.ID,
			State:     &completedState,
		})
		require.NoError(t, err)
		assert.Len(t, completedTasks, 1)
		assert.Equal(t, subtask.ID, completedTasks[0].ID)

		// Test task deletion
		err = repo.DeleteTask(ctx, subtask.ID)
		require.NoError(t, err)

		_, err = repo.GetTask(ctx, subtask.ID)
		assert.Error(t, err)

		// Verify parent task still exists
		_, err = repo.GetTask(ctx, rootTask.ID)
		require.NoError(t, err)
	})

	t.Run("Hierarchy Operations", func(t *testing.T) {
		// Create project
		project := &Project{
			ID:          uuid.New(),
			Title:       "Hierarchy Test",
			Description: "Testing hierarchy",
		}
		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)

		// Create task hierarchy
		//   Root1
		//   ├── Child1.1
		//   │   └── Grandchild1.1.1
		//   └── Child1.2
		//   Root2

		root1 := &Task{
			ID:        uuid.New(),
			ProjectID: project.ID,
			Title:     "Root1",
			State:     TaskStatePending,
			Complexity: 5,
			Priority:  8,
		}
		err = repo.CreateTask(ctx, root1)
		require.NoError(t, err)

		child11 := &Task{
			ID:        uuid.New(),
			ProjectID: project.ID,
			ParentID:  &root1.ID,
			Title:     "Child1.1",
			State:     TaskStatePending,
			Complexity: 3,
			Priority:  6,
		}
		err = repo.CreateTask(ctx, child11)
		require.NoError(t, err)

		grandchild111 := &Task{
			ID:        uuid.New(),
			ProjectID: project.ID,
			ParentID:  &child11.ID,
			Title:     "Grandchild1.1.1",
			State:     TaskStateCompleted,
			Complexity: 2,
			Priority:  4,
		}
		err = repo.CreateTask(ctx, grandchild111)
		require.NoError(t, err)

		child12 := &Task{
			ID:        uuid.New(),
			ProjectID: project.ID,
			ParentID:  &root1.ID,
			Title:     "Child1.2",
			State:     TaskStateInProgress,
			Complexity: 4,
			Priority:  7,
		}
		err = repo.CreateTask(ctx, child12)
		require.NoError(t, err)

		root2 := &Task{
			ID:        uuid.New(),
			ProjectID: project.ID,
			Title:     "Root2",
			State:     TaskStatePending,
			Complexity: 6,
			Priority:  5,
		}
		err = repo.CreateTask(ctx, root2)
		require.NoError(t, err)

		// Test hierarchy retrieval
		hierarchy, err := repo.GetTaskHierarchy(ctx, project.ID)
		require.NoError(t, err)
		assert.Len(t, hierarchy, 2)

		// Find Root1 in hierarchy (should be first due to higher priority)
		var root1Hierarchy *TaskHierarchy
		for _, h := range hierarchy {
			if h.Task.Title == "Root1" {
				root1Hierarchy = h
				break
			}
		}
		require.NotNil(t, root1Hierarchy)
		assert.Len(t, root1Hierarchy.Children, 2)

		// Verify child ordering (Child1.2 should be first due to higher priority)
		assert.Equal(t, "Child1.2", root1Hierarchy.Children[0].Task.Title)
		assert.Equal(t, "Child1.1", root1Hierarchy.Children[1].Task.Title)

		// Verify grandchild
		child11Hierarchy := root1Hierarchy.Children[1]
		assert.Len(t, child11Hierarchy.Children, 1)
		assert.Equal(t, "Grandchild1.1.1", child11Hierarchy.Children[0].Task.Title)

		// Test subtree retrieval
		subtree, err := repo.GetTaskSubtree(ctx, root1.ID)
		require.NoError(t, err)
		assert.Equal(t, "Root1", subtree.Task.Title)
		assert.Len(t, subtree.Children, 2)

		// Test subtree deletion
		err = repo.DeleteTaskSubtree(ctx, root1.ID)
		require.NoError(t, err)

		// Verify all tasks in subtree are deleted
		_, err = repo.GetTask(ctx, root1.ID)
		assert.Error(t, err)
		_, err = repo.GetTask(ctx, child11.ID)
		assert.Error(t, err)
		_, err = repo.GetTask(ctx, grandchild111.ID)
		assert.Error(t, err)
		_, err = repo.GetTask(ctx, child12.ID)
		assert.Error(t, err)

		// Verify Root2 still exists
		_, err = repo.GetTask(ctx, root2.ID)
		require.NoError(t, err)

		// Verify hierarchy now only has Root2
		hierarchy, err = repo.GetTaskHierarchy(ctx, project.ID)
		require.NoError(t, err)
		assert.Len(t, hierarchy, 1)
		assert.Equal(t, "Root2", hierarchy[0].Task.Title)
	})

	t.Run("Progress Metrics", func(t *testing.T) {
		// Create project
		project := &Project{
			ID:          uuid.New(),
			Title:       "Progress Test",
			Description: "Testing progress metrics",
		}
		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)

		// Create tasks with different states
		tasks := []*Task{
			{
				ID:        uuid.New(),
				ProjectID: project.ID,
				Title:     "Completed Task",
				State:     TaskStateCompleted,
				Complexity: 5,
				Priority:  5,
			},
			{
				ID:        uuid.New(),
				ProjectID: project.ID,
				Title:     "In Progress Task",
				State:     TaskStateInProgress,
				Complexity: 3,
				Priority:  6,
			},
			{
				ID:        uuid.New(),
				ProjectID: project.ID,
				Title:     "Pending Task",
				State:     TaskStatePending,
				Complexity: 4,
				Priority:  7,
			},
			{
				ID:        uuid.New(),
				ProjectID: project.ID,
				Title:     "Blocked Task",
				State:     TaskStateBlocked,
				Complexity: 2,
				Priority:  3,
			},
		}

		for _, task := range tasks {
			err = repo.CreateTask(ctx, task)
			require.NoError(t, err)
		}

		// Test progress metrics
		progress, err := repo.GetProjectProgress(ctx, project.ID)
		require.NoError(t, err)

		assert.Equal(t, project.ID, progress.ProjectID)
		assert.Equal(t, 4, progress.TotalTasks)
		assert.Equal(t, 1, progress.CompletedTasks)
		assert.Equal(t, 1, progress.InProgressTasks)
		assert.Equal(t, 1, progress.PendingTasks)
		assert.Equal(t, 1, progress.BlockedTasks)
		assert.Equal(t, 0, progress.CancelledTasks)
		assert.Equal(t, 25.0, progress.OverallProgress)
		assert.Equal(t, 4, progress.TasksByDepth[0])

		// Test task count by depth
		counts, err := repo.GetTaskCountByDepth(ctx, project.ID, 2)
		require.NoError(t, err)
		assert.Equal(t, 4, counts[0])
		assert.Equal(t, 0, counts[1])
		assert.Equal(t, 0, counts[2])
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		// Create project
		project := &Project{
			ID:          uuid.New(),
			Title:       "Concurrent Test",
			Description: "Testing concurrent operations",
		}
		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)

		// Create tasks concurrently
		const numTasks = 100
		results := make(chan error, numTasks)

		for i := 0; i < numTasks; i++ {
			go func(taskNum int) {
				task := &Task{
					ID:        uuid.New(),
					ProjectID: project.ID,
					Title:     fmt.Sprintf("Concurrent Task %d", taskNum),
					State:     TaskStatePending,
					Complexity: 5,
					Priority:  5,
				}
				results <- repo.CreateTask(ctx, task)
			}(i)
		}

		// Wait for all operations to complete
		for i := 0; i < numTasks; i++ {
			err := <-results
			assert.NoError(t, err)
		}

		// Verify all tasks were created
		tasks, err := repo.GetTasksByProject(ctx, project.ID)
		require.NoError(t, err)
		assert.Len(t, tasks, numTasks)

		// Update tasks concurrently
		updateResults := make(chan error, numTasks)
		for _, task := range tasks {
			go func(t *Task) {
				t.State = TaskStateCompleted
				updateResults <- repo.UpdateTask(ctx, t)
			}(task)
		}

		// Wait for all updates
		for i := 0; i < numTasks; i++ {
			err := <-updateResults
			assert.NoError(t, err)
		}

		// Verify progress
		progress, err := repo.GetProjectProgress(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, numTasks, progress.TotalTasks)
		assert.Equal(t, numTasks, progress.CompletedTasks)
		assert.Equal(t, 100.0, progress.OverallProgress)
	})

	t.Run("Error Cases", func(t *testing.T) {
		// Test getting non-existent project
		_, err := repo.GetProject(ctx, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		// Test getting non-existent task
		_, err = repo.GetTask(ctx, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		// Test creating task with non-existent project
		task := &Task{
			ID:        uuid.New(),
			ProjectID: uuid.New(),
			Title:     "Invalid Task",
			State:     TaskStatePending,
			Complexity: 5,
			Priority:  5,
		}
		err = repo.CreateTask(ctx, task)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		// Test creating task with non-existent parent
		project := &Project{
			ID:    uuid.New(),
			Title: "Test Project",
		}
		err = repo.CreateProject(ctx, project)
		require.NoError(t, err)

		parentID := uuid.New()
		task.ProjectID = project.ID
		task.ParentID = &parentID
		err = repo.CreateTask(ctx, task)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent task")

		// Test deleting task with children
		rootTask := &Task{
			ID:        uuid.New(),
			ProjectID: project.ID,
			Title:     "Root",
			State:     TaskStatePending,
			Complexity: 5,
			Priority:  5,
		}
		err = repo.CreateTask(ctx, rootTask)
		require.NoError(t, err)

		childTask := &Task{
			ID:        uuid.New(),
			ProjectID: project.ID,
			ParentID:  &rootTask.ID,
			Title:     "Child",
			State:     TaskStatePending,
			Complexity: 3,
			Priority:  5,
		}
		err = repo.CreateTask(ctx, childTask)
		require.NoError(t, err)

		err = repo.DeleteTask(ctx, rootTask.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete task with children")
	})
}