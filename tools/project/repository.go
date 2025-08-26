package project

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/denkhaus/agents/shared/resource"
	"github.com/google/uuid"
)

// memoryRepository implements repository interface with in-memory storage using resource managers
type memoryRepository struct {
	projects       *resource.Manager[*Project]
	tasks          *resource.Manager[*Task]
	tasksByProject *resource.Manager[[]uuid.UUID]
	tasksByParent  *resource.Manager[[]uuid.UUID]
}

// newMemoryRepository creates a new in-memory repository
func newMemoryRepository() repository {
	return &memoryRepository{
		projects:       resource.NewManager[*Project](),
		tasks:          resource.NewManager[*Task](),
		tasksByProject: resource.NewManager[[]uuid.UUID](),
		tasksByParent:  resource.NewManager[[]uuid.UUID](),
	}
}

// Project operations

func (r *memoryRepository) CreateProject(ctx context.Context, project *Project) error {
	if r.projects.Exists(project.ID) {
		return fmt.Errorf("project with ID %s already exists", project.ID)
	}

	projectCopy := *project
	projectCopy.CreatedAt = time.Now()
	projectCopy.UpdatedAt = projectCopy.CreatedAt
	projectCopy.TotalTasks = 0
	projectCopy.CompletedTasks = 0
	projectCopy.Progress = 0.0

	r.projects.Set(project.ID, &projectCopy)
	r.tasksByProject.Set(project.ID, make([]uuid.UUID, 0))

	return nil
}

func (r *memoryRepository) GetProject(ctx context.Context, id uuid.UUID) (*Project, error) {
	project, exists := r.projects.Get(id)
	if !exists {
		return nil, fmt.Errorf("project with ID %s not found", id)
	}

	projectCopy := *project
	return &projectCopy, nil
}

func (r *memoryRepository) UpdateProject(ctx context.Context, project *Project) error {
	return r.projects.UpdateWithError(project.ID, func(existing *Project) (*Project, error) {
		projectCopy := *project
		projectCopy.CreatedAt = existing.CreatedAt
		projectCopy.UpdatedAt = time.Now()
		return &projectCopy, nil
	})
}

func (r *memoryRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	if !r.projects.Exists(id) {
		return fmt.Errorf("project with ID %s not found", id)
	}

	if taskIDs, exists := r.tasksByProject.Get(id); exists {
		for _, taskID := range taskIDs {
			r.tasks.Delete(taskID)
			r.tasksByParent.Delete(taskID)
		}
	}

	r.projects.Delete(id)
	r.tasksByProject.Delete(id)

	return nil
}

func (r *memoryRepository) ListProjects(ctx context.Context) ([]*Project, error) {
	allProjects := r.projects.GetAll()
	projects := make([]*Project, 0, len(allProjects))

	for _, project := range allProjects {
		projectCopy := *project
		projects = append(projects, &projectCopy)
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].CreatedAt.Before(projects[j].CreatedAt)
	})

	return projects, nil
}

// Task operations

func (r *memoryRepository) CreateTask(ctx context.Context, task *Task) error {
	if r.tasks.Exists(task.ID) {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	if !r.projects.Exists(task.ProjectID) {
		return fmt.Errorf("project with ID %s not found", task.ProjectID)
	}

	if task.ParentID != nil {
		parentTask, exists := r.tasks.Get(*task.ParentID)
		if !exists {
			return fmt.Errorf("parent task with ID %s not found", *task.ParentID)
		}
		if parentTask.ProjectID != task.ProjectID {
			return fmt.Errorf("parent task must be in the same project")
		}
		task.Depth = parentTask.Depth + 1
	} else {
		task.Depth = 0
	}

	taskCopy := *task
	taskCopy.CreatedAt = time.Now()
	taskCopy.UpdatedAt = taskCopy.CreatedAt

	r.tasks.Set(task.ID, &taskCopy)
	r.tasksByParent.Set(task.ID, make([]uuid.UUID, 0))

	// Update indexes
	r.tasksByProject.Upsert(task.ProjectID, func(tasks []uuid.UUID) []uuid.UUID {
		if tasks == nil {
			tasks = make([]uuid.UUID, 0)
		}
		return append(tasks, task.ID)
	})

	if task.ParentID != nil {
		r.tasksByParent.Upsert(*task.ParentID, func(tasks []uuid.UUID) []uuid.UUID {
			if tasks == nil {
				tasks = make([]uuid.UUID, 0)
			}
			return append(tasks, task.ID)
		})
	}

	r.updateProjectMetrics(task.ProjectID)
	return nil
}

func (r *memoryRepository) GetTask(ctx context.Context, id uuid.UUID) (*Task, error) {
	task, exists := r.tasks.Get(id)
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", id)
	}

	taskCopy := *task
	return &taskCopy, nil
}

func (r *memoryRepository) UpdateTask(ctx context.Context, task *Task) error {
	var projectID uuid.UUID
	err := r.tasks.UpdateWithError(task.ID, func(existing *Task) (*Task, error) {
		taskCopy := *task
		taskCopy.CreatedAt = existing.CreatedAt
		taskCopy.UpdatedAt = time.Now()
		taskCopy.ProjectID = existing.ProjectID
		taskCopy.ParentID = existing.ParentID
		taskCopy.Depth = existing.Depth
		projectID = existing.ProjectID

		if task.State == TaskStateCompleted && existing.State != TaskStateCompleted {
			now := time.Now()
			taskCopy.CompletedAt = &now
		} else if task.State != TaskStateCompleted {
			taskCopy.CompletedAt = nil
		}

		return &taskCopy, nil
	})

	// Update metrics outside the lock to avoid deadlock
	if err == nil {
		r.updateProjectMetrics(projectID)
	}

	return err
}

func (r *memoryRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	task, exists := r.tasks.Get(id)
	if !exists {
		return fmt.Errorf("task with ID %s not found", id)
	}

	if children, exists := r.tasksByParent.Get(id); exists && len(children) > 0 {
		return fmt.Errorf("cannot delete task with children, use DeleteTaskSubtree instead")
	}

	return r.deleteTaskInternal(id, task.ProjectID)
}

func (r *memoryRepository) deleteTaskInternal(taskID uuid.UUID, projectID uuid.UUID) error {
	task, exists := r.tasks.Get(taskID)
	if !exists {
		return fmt.Errorf("task not found")
	}

	// Remove from parent's children list
	if task.ParentID != nil {
		if children, exists := r.tasksByParent.Get(*task.ParentID); exists {
			for i, childID := range children {
				if childID == taskID {
					children = append(children[:i], children[i+1:]...)
					r.tasksByParent.Set(*task.ParentID, children)
					break
				}
			}
		}
	}

	// Remove from project's task list
	if projectTasks, exists := r.tasksByProject.Get(projectID); exists {
		for i, id := range projectTasks {
			if id == taskID {
				projectTasks = append(projectTasks[:i], projectTasks[i+1:]...)
				r.tasksByProject.Set(projectID, projectTasks)
				break
			}
		}
	}

	r.tasks.Delete(taskID)
	r.tasksByParent.Delete(taskID)

	return nil
}

func (r *memoryRepository) updateProjectMetrics(projectID uuid.UUID) {
	r.projects.Update(projectID, func(project *Project) *Project {
		taskIDs, exists := r.tasksByProject.Get(projectID)
		if !exists {
			return project
		}

		totalTasks := len(taskIDs)
		completedTasks := 0

		for _, taskID := range taskIDs {
			if task, exists := r.tasks.Get(taskID); exists && task.State == TaskStateCompleted {
				completedTasks++
			}
		}

		progress := 0.0
		if totalTasks > 0 {
			progress = float64(completedTasks) / float64(totalTasks) * 100.0
		}

		updatedProject := *project
		updatedProject.TotalTasks = totalTasks
		updatedProject.CompletedTasks = completedTasks
		updatedProject.Progress = progress
		updatedProject.UpdatedAt = time.Now()

		return &updatedProject
	})
}

// Simplified implementations for the remaining methods

func (r *memoryRepository) ListTasks(ctx context.Context, filter TaskFilter) ([]*Task, error) {
	allTasks := r.tasks.GetAll()
	var tasks []*Task

	for _, task := range allTasks {
		if r.matchesFilter(task, filter) {
			taskCopy := *task
			tasks = append(tasks, &taskCopy)
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Priority != tasks[j].Priority {
			return tasks[i].Priority > tasks[j].Priority
		}
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *memoryRepository) matchesFilter(task *Task, filter TaskFilter) bool {
	if filter.ProjectID != nil && task.ProjectID != *filter.ProjectID {
		return false
	}
	if filter.ParentID != nil {
		if task.ParentID == nil || *task.ParentID != *filter.ParentID {
			return false
		}
	}
	if filter.State != nil && task.State != *filter.State {
		return false
	}
	if filter.MinDepth != nil && task.Depth < *filter.MinDepth {
		return false
	}
	if filter.MaxDepth != nil && task.Depth > *filter.MaxDepth {
		return false
	}
	if filter.MinComplexity != nil && task.Complexity < *filter.MinComplexity {
		return false
	}
	if filter.MaxComplexity != nil && task.Complexity > *filter.MaxComplexity {
		return false
	}
	return true
}

func (r *memoryRepository) GetTasksByProject(ctx context.Context, projectID uuid.UUID) ([]*Task, error) {
	return r.ListTasks(ctx, TaskFilter{ProjectID: &projectID})
}

func (r *memoryRepository) GetTasksByParent(ctx context.Context, parentID uuid.UUID) ([]*Task, error) {
	childIDs, exists := r.tasksByParent.Get(parentID)
	if !exists {
		return []*Task{}, nil
	}

	var childTasks []*Task
	for _, childID := range childIDs {
		if task, exists := r.tasks.Get(childID); exists {
			taskCopy := *task
			childTasks = append(childTasks, &taskCopy)
		}
	}

	// Sort by priority then creation time
	sort.Slice(childTasks, func(i, j int) bool {
		if childTasks[i].Priority != childTasks[j].Priority {
			return childTasks[i].Priority > childTasks[j].Priority
		}
		return childTasks[i].CreatedAt.Before(childTasks[j].CreatedAt)
	})

	return childTasks, nil
}

func (r *memoryRepository) GetRootTasks(ctx context.Context, projectID uuid.UUID) ([]*Task, error) {
	taskIDs, exists := r.tasksByProject.Get(projectID)
	if !exists {
		return []*Task{}, nil
	}

	var rootTasks []*Task
	for _, taskID := range taskIDs {
		if task, exists := r.tasks.Get(taskID); exists && task.ParentID == nil {
			taskCopy := *task
			rootTasks = append(rootTasks, &taskCopy)
		}
	}

	sort.Slice(rootTasks, func(i, j int) bool {
		if rootTasks[i].Priority != rootTasks[j].Priority {
			return rootTasks[i].Priority > rootTasks[j].Priority
		}
		return rootTasks[i].CreatedAt.Before(rootTasks[j].CreatedAt)
	})

	return rootTasks, nil
}

func (r *memoryRepository) GetParentTask(ctx context.Context, taskID uuid.UUID) (*Task, error) {
	task, err := r.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task.ParentID == nil {
		return nil, nil // No parent
	}
	return r.GetTask(ctx, *task.ParentID)
}

func (r *memoryRepository) DeleteTaskSubtree(ctx context.Context, taskID uuid.UUID) error {
	task, exists := r.tasks.Get(taskID)
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	// Recursively delete all children first
	if childIDs, exists := r.tasksByParent.Get(taskID); exists {
		for _, childID := range childIDs {
			if err := r.DeleteTaskSubtree(ctx, childID); err != nil {
				return err
			}
		}
	}

	// Delete the task itself
	if err := r.deleteTaskInternal(taskID, task.ProjectID); err != nil {
		return err
	}

	r.updateProjectMetrics(task.ProjectID)
	return nil
}

func (r *memoryRepository) GetProjectProgress(ctx context.Context, projectID uuid.UUID) (*ProjectProgress, error) {
	if !r.projects.Exists(projectID) {
		return nil, fmt.Errorf("project with ID %s not found", projectID)
	}

	taskIDs, exists := r.tasksByProject.Get(projectID)
	if !exists {
		taskIDs = []uuid.UUID{}
	}

	progress := &ProjectProgress{
		ProjectID:    projectID,
		TasksByDepth: make(map[int]int),
	}

	for _, taskID := range taskIDs {
		task, exists := r.tasks.Get(taskID)
		if !exists {
			continue
		}

		progress.TotalTasks++
		progress.TasksByDepth[task.Depth]++

		switch task.State {
		case TaskStateCompleted:
			progress.CompletedTasks++
		case TaskStateInProgress:
			progress.InProgressTasks++
		case TaskStatePending:
			progress.PendingTasks++
		case TaskStateBlocked:
			progress.BlockedTasks++
		case TaskStateCancelled:
			progress.CancelledTasks++
		}
	}

	if progress.TotalTasks > 0 {
		progress.OverallProgress = float64(progress.CompletedTasks) / float64(progress.TotalTasks) * 100.0
	}

	return progress, nil
}

func (r *memoryRepository) GetTaskCountByDepth(ctx context.Context, projectID uuid.UUID, maxDepth int) (map[int]int, error) {
	if !r.projects.Exists(projectID) {
		return nil, fmt.Errorf("project with ID %s not found", projectID)
	}

	counts := make(map[int]int)
	taskIDs, exists := r.tasksByProject.Get(projectID)
	if !exists {
		return counts, nil
	}

	for _, taskID := range taskIDs {
		task, exists := r.tasks.Get(taskID)
		if exists && task.Depth <= maxDepth {
			counts[task.Depth]++
		}
	}

	return counts, nil
}
