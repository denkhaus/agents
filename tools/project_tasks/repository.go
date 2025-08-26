package projecttasks

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// memoryRepository implements Repository interface with in-memory storage
type memoryRepository struct {
	mu       sync.RWMutex
	projects map[uuid.UUID]*Project
	tasks    map[uuid.UUID]*Task
	// Index for efficient queries
	tasksByProject map[uuid.UUID][]uuid.UUID
	tasksByParent  map[uuid.UUID][]uuid.UUID
}

// newMemoryRepository creates a new in-memory repository
func newMemoryRepository() repository {
	return &memoryRepository{
		projects:       make(map[uuid.UUID]*Project),
		tasks:          make(map[uuid.UUID]*Task),
		tasksByProject: make(map[uuid.UUID][]uuid.UUID),
		tasksByParent:  make(map[uuid.UUID][]uuid.UUID),
	}
}

// Project operations

func (r *memoryRepository) CreateProject(ctx context.Context, project *Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.projects[project.ID]; exists {
		return fmt.Errorf("project with ID %s already exists", project.ID)
	}

	// Create a copy to avoid external modifications
	projectCopy := *project
	projectCopy.CreatedAt = time.Now()
	projectCopy.UpdatedAt = projectCopy.CreatedAt
	projectCopy.TotalTasks = 0
	projectCopy.CompletedTasks = 0
	projectCopy.Progress = 0.0

	r.projects[project.ID] = &projectCopy
	r.tasksByProject[project.ID] = make([]uuid.UUID, 0)

	return nil
}

func (r *memoryRepository) GetProject(ctx context.Context, id uuid.UUID) (*Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	project, exists := r.projects[id]
	if !exists {
		return nil, fmt.Errorf("project with ID %s not found", id)
	}

	// Return a copy to prevent external modifications
	projectCopy := *project
	return &projectCopy, nil
}

func (r *memoryRepository) UpdateProject(ctx context.Context, project *Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.projects[project.ID]
	if !exists {
		return fmt.Errorf("project with ID %s not found", project.ID)
	}

	// Preserve creation time and update timestamp
	projectCopy := *project
	projectCopy.CreatedAt = existing.CreatedAt
	projectCopy.UpdatedAt = time.Now()

	r.projects[project.ID] = &projectCopy
	return nil
}

func (r *memoryRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.projects[id]; !exists {
		return fmt.Errorf("project with ID %s not found", id)
	}

	// Delete all tasks in the project
	taskIDs := r.tasksByProject[id]
	for _, taskID := range taskIDs {
		delete(r.tasks, taskID)
		// Clean up parent-child relationships
		delete(r.tasksByParent, taskID)
	}

	delete(r.projects, id)
	delete(r.tasksByProject, id)

	return nil
}

func (r *memoryRepository) ListProjects(ctx context.Context) ([]*Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projects := make([]*Project, 0, len(r.projects))
	for _, project := range r.projects {
		projectCopy := *project
		projects = append(projects, &projectCopy)
	}

	// Sort by creation time
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].CreatedAt.Before(projects[j].CreatedAt)
	})

	return projects, nil
}

// Task operations

func (r *memoryRepository) CreateTask(ctx context.Context, task *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	// Validate project exists
	if _, exists := r.projects[task.ProjectID]; !exists {
		return fmt.Errorf("project with ID %s not found", task.ProjectID)
	}

	// Validate parent task if specified
	if task.ParentID != nil {
		parentTask, exists := r.tasks[*task.ParentID]
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

	// Create a copy and set timestamps
	taskCopy := *task
	taskCopy.CreatedAt = time.Now()
	taskCopy.UpdatedAt = taskCopy.CreatedAt

	r.tasks[task.ID] = &taskCopy

	// Update indexes
	r.tasksByProject[task.ProjectID] = append(r.tasksByProject[task.ProjectID], task.ID)
	if task.ParentID != nil {
		r.tasksByParent[*task.ParentID] = append(r.tasksByParent[*task.ParentID], task.ID)
	}

	// Update project metrics
	r.updateProjectMetrics(task.ProjectID)

	return nil
}

func (r *memoryRepository) GetTask(ctx context.Context, id uuid.UUID) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", id)
	}

	taskCopy := *task
	return &taskCopy, nil
}

func (r *memoryRepository) UpdateTask(ctx context.Context, task *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.tasks[task.ID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	// Preserve immutable fields
	taskCopy := *task
	taskCopy.CreatedAt = existing.CreatedAt
	taskCopy.UpdatedAt = time.Now()
	taskCopy.ProjectID = existing.ProjectID
	taskCopy.ParentID = existing.ParentID
	taskCopy.Depth = existing.Depth

	// Set completion time if task is being completed
	if task.State == TaskStateCompleted && existing.State != TaskStateCompleted {
		now := time.Now()
		taskCopy.CompletedAt = &now
	} else if task.State != TaskStateCompleted {
		taskCopy.CompletedAt = nil
	}

	r.tasks[task.ID] = &taskCopy

	// Update project metrics
	r.updateProjectMetrics(task.ProjectID)

	return nil
}

func (r *memoryRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %s not found", id)
	}

	// Check if task has children
	if children := r.tasksByParent[id]; len(children) > 0 {
		return fmt.Errorf("cannot delete task with children, use DeleteTaskSubtree instead")
	}

	return r.deleteTaskInternal(id, task.ProjectID)
}

func (r *memoryRepository) deleteTaskInternal(taskID uuid.UUID, projectID uuid.UUID) error {
	// Remove from parent's children list
	task := r.tasks[taskID]
	if task.ParentID != nil {
		children := r.tasksByParent[*task.ParentID]
		for i, childID := range children {
			if childID == taskID {
				r.tasksByParent[*task.ParentID] = append(children[:i], children[i+1:]...)
				break
			}
		}
	}

	// Remove from project's task list
	projectTasks := r.tasksByProject[projectID]
	for i, id := range projectTasks {
		if id == taskID {
			r.tasksByProject[projectID] = append(projectTasks[:i], projectTasks[i+1:]...)
			break
		}
	}

	// Delete the task
	delete(r.tasks, taskID)
	delete(r.tasksByParent, taskID)

	return nil
}

// Helper method to update project metrics (must be called with lock held)
func (r *memoryRepository) updateProjectMetrics(projectID uuid.UUID) {
	project := r.projects[projectID]
	if project == nil {
		return
	}

	taskIDs := r.tasksByProject[projectID]
	totalTasks := len(taskIDs)
	completedTasks := 0

	for _, taskID := range taskIDs {
		if task := r.tasks[taskID]; task != nil && task.State == TaskStateCompleted {
			completedTasks++
		}
	}

	progress := 0.0
	if totalTasks > 0 {
		progress = float64(completedTasks) / float64(totalTasks) * 100.0
	}

	project.TotalTasks = totalTasks
	project.CompletedTasks = completedTasks
	project.Progress = progress
	project.UpdatedAt = time.Now()
}

// Task query operations

func (r *memoryRepository) ListTasks(ctx context.Context, filter TaskFilter) ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tasks []*Task
	for _, task := range r.tasks {
		if r.matchesFilter(task, filter) {
			taskCopy := *task
			tasks = append(tasks, &taskCopy)
		}
	}

	// Sort by priority (descending) then by creation time
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
	return r.ListTasks(ctx, TaskFilter{ParentID: &parentID})
}

func (r *memoryRepository) GetRootTasks(ctx context.Context, projectID uuid.UUID) ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var rootTasks []*Task
	taskIDs := r.tasksByProject[projectID]
	
	for _, taskID := range taskIDs {
		if task := r.tasks[taskID]; task != nil && task.ParentID == nil {
			taskCopy := *task
			rootTasks = append(rootTasks, &taskCopy)
		}
	}

	// Sort by priority then creation time
	sort.Slice(rootTasks, func(i, j int) bool {
		if rootTasks[i].Priority != rootTasks[j].Priority {
			return rootTasks[i].Priority > rootTasks[j].Priority
		}
		return rootTasks[i].CreatedAt.Before(rootTasks[j].CreatedAt)
	})

	return rootTasks, nil
}

// Hierarchy operations

func (r *memoryRepository) GetTaskHierarchy(ctx context.Context, projectID uuid.UUID) ([]*TaskHierarchy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rootTasks, err := r.GetRootTasks(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var hierarchies []*TaskHierarchy
	for _, rootTask := range rootTasks {
		hierarchy := r.buildTaskHierarchy(rootTask)
		hierarchies = append(hierarchies, hierarchy)
	}

	return hierarchies, nil
}

func (r *memoryRepository) buildTaskHierarchy(task *Task) *TaskHierarchy {
	hierarchy := &TaskHierarchy{
		Task:     task,
		Children: make([]*TaskHierarchy, 0),
	}

	childIDs := r.tasksByParent[task.ID]
	for _, childID := range childIDs {
		if childTask := r.tasks[childID]; childTask != nil {
			childCopy := *childTask
			childHierarchy := r.buildTaskHierarchy(&childCopy)
			hierarchy.Children = append(hierarchy.Children, childHierarchy)
		}
	}

	// Sort children by priority then creation time
	sort.Slice(hierarchy.Children, func(i, j int) bool {
		if hierarchy.Children[i].Task.Priority != hierarchy.Children[j].Task.Priority {
			return hierarchy.Children[i].Task.Priority > hierarchy.Children[j].Task.Priority
		}
		return hierarchy.Children[i].Task.CreatedAt.Before(hierarchy.Children[j].Task.CreatedAt)
	})

	return hierarchy
}

func (r *memoryRepository) GetTaskSubtree(ctx context.Context, taskID uuid.UUID) (*TaskHierarchy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", taskID)
	}

	taskCopy := *task
	return r.buildTaskHierarchy(&taskCopy), nil
}

func (r *memoryRepository) DeleteTaskSubtree(ctx context.Context, taskID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	// Recursively delete all children first
	childIDs := r.tasksByParent[taskID]
	for _, childID := range childIDs {
		if err := r.DeleteTaskSubtree(ctx, childID); err != nil {
			return err
		}
	}

	// Delete the task itself
	if err := r.deleteTaskInternal(taskID, task.ProjectID); err != nil {
		return err
	}

	// Update project metrics
	r.updateProjectMetrics(task.ProjectID)

	return nil
}

// Metrics and analysis

func (r *memoryRepository) GetProjectProgress(ctx context.Context, projectID uuid.UUID) (*ProjectProgress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, exists := r.projects[projectID]; !exists {
		return nil, fmt.Errorf("project with ID %s not found", projectID)
	}

	taskIDs := r.tasksByProject[projectID]
	progress := &ProjectProgress{
		ProjectID:    projectID,
		TasksByDepth: make(map[int]int),
	}

	for _, taskID := range taskIDs {
		task := r.tasks[taskID]
		if task == nil {
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
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, exists := r.projects[projectID]; !exists {
		return nil, fmt.Errorf("project with ID %s not found", projectID)
	}

	counts := make(map[int]int)
	taskIDs := r.tasksByProject[projectID]

	for _, taskID := range taskIDs {
		task := r.tasks[taskID]
		if task != nil && task.Depth <= maxDepth {
			counts[task.Depth]++
		}
	}

	return counts, nil
}