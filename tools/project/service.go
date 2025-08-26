package project

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// service provides business logic for project task management
type service struct {
	repo   repository
	config *Config
}

// newService creates a new task management service
func newService(repo repository, config *Config) *service {
	if config == nil {
		config = DefaultConfig()
	}
	return &service{
		repo:   repo,
		config: config,
	}
}

// Project operations

func (s *service) CreateProject(ctx context.Context, title, description string) (*Project, error) {
	if err := s.validateProjectInput(title, description); err != nil {
		return nil, err
	}

	project := &Project{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
	}

	if err := s.repo.CreateProject(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return s.repo.GetProject(ctx, project.ID)
}

func (s *service) GetProject(ctx context.Context, projectID uuid.UUID) (*Project, error) {
	return s.repo.GetProject(ctx, projectID)
}

func (s *service) UpdateProject(ctx context.Context, projectID uuid.UUID, title, description string) (*Project, error) {
	if err := s.validateProjectInput(title, description); err != nil {
		return nil, err
	}

	project, err := s.repo.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	project.Title = title
	project.Description = description

	if err := s.repo.UpdateProject(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return s.repo.GetProject(ctx, projectID)
}

func (s *service) UpdateProjectDescription(ctx context.Context, projectID uuid.UUID, description string) (*Project, error) {
	// Validate description length
	if len(description) > s.config.MaxDescriptionLength {
		return nil, ValidationError{Field: "description", Message: fmt.Sprintf("description cannot exceed %d characters", s.config.MaxDescriptionLength)}
	}

	project, err := s.repo.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	project.Description = description

	if err := s.repo.UpdateProject(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project description: %w", err)
	}

	return s.repo.GetProject(ctx, projectID)
}

func (s *service) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	return s.repo.DeleteProject(ctx, projectID)
}

func (s *service) ListProjects(ctx context.Context) ([]*Project, error) {
	return s.repo.ListProjects(ctx)
}

// Task operations

func (s *service) CreateTask(ctx context.Context, projectID uuid.UUID, parentID *uuid.UUID, title, description string, complexity int) (*Task, error) {
	if err := s.validateTaskInput(title, description, complexity); err != nil {
		return nil, err
	}

	// Validate project exists
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Calculate depth and validate constraints
	depth := 0
	if parentID != nil {
		parentTask, err := s.repo.GetTask(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("parent task not found: %w", err)
		}
		if parentTask.ProjectID != projectID {
			return nil, fmt.Errorf("parent task must be in the same project")
		}
		depth = parentTask.Depth + 1
	}

	// Check depth constraints
	if depth > s.config.MaxDepth {
		return nil, fmt.Errorf("maximum depth of %d exceeded", s.config.MaxDepth)
	}

	// Check task count constraints for this depth
	counts, err := s.repo.GetTaskCountByDepth(ctx, projectID, depth)
	if err != nil {
		return nil, fmt.Errorf("failed to check task count constraints: %w", err)
	}
	if counts[depth] >= s.config.MaxTasksPerDepth {
		return nil, fmt.Errorf("maximum tasks per depth (%d) exceeded for depth %d", s.config.MaxTasksPerDepth, depth)
	}

	task := &Task{
		ID:          uuid.New(),
		ProjectID:   projectID,
		ParentID:    parentID,
		Title:       title,
		Description: description,
		State:       TaskStatePending,
		Complexity:  complexity,
		Depth:       depth,
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return s.repo.GetTask(ctx, task.ID)
}

func (s *service) GetTask(ctx context.Context, taskID uuid.UUID) (*Task, error) {
	return s.repo.GetTask(ctx, taskID)
}

func (s *service) UpdateTaskState(ctx context.Context, taskID uuid.UUID, state TaskState) (*Task, error) {
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task.State = state
	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task state: %w", err)
	}

	return s.repo.GetTask(ctx, taskID)
}

func (s *service) UpdateTask(ctx context.Context, taskID uuid.UUID, title, description string, complexity int, state TaskState) (*Task, error) {
	if err := s.validateTaskInput(title, description, complexity); err != nil {
		return nil, err
	}

	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task.Title = title
	task.Description = description
	task.Complexity = complexity
	task.State = state

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return s.repo.GetTask(ctx, taskID)
}

func (s *service) UpdateTaskDescription(ctx context.Context, taskID uuid.UUID, description string) (*Task, error) {
	// Validate description length
	if len(description) > s.config.MaxDescriptionLength {
		return nil, ValidationError{Field: "description", Message: fmt.Sprintf("description cannot exceed %d characters", s.config.MaxDescriptionLength)}
	}

	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task.Description = description

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task description: %w", err)
	}

	return s.repo.GetTask(ctx, taskID)
}

func (s *service) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	return s.repo.DeleteTask(ctx, taskID)
}

func (s *service) DeleteTaskSubtree(ctx context.Context, taskID uuid.UUID) error {
	return s.repo.DeleteTaskSubtree(ctx, taskID)
}

// Task queries and analysis

func (s *service) GetParentTask(ctx context.Context, taskID uuid.UUID) (*Task, error) {
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task.ParentID == nil {
		return nil, nil // No parent
	}
	return s.repo.GetTask(ctx, *task.ParentID)
}

func (s *service) GetChildTasks(ctx context.Context, taskID uuid.UUID) ([]*Task, error) {
	// Validate task exists
	if _, err := s.repo.GetTask(ctx, taskID); err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	return s.repo.GetTasksByParent(ctx, taskID)
}

func (s *service) GetRootTasks(ctx context.Context, projectID uuid.UUID) ([]*Task, error) {
	return s.repo.GetRootTasks(ctx, projectID)
}

// ListTasksForProject returns all tasks in a project regardless of hierarchy level
func (s *service) ListTasksForProject(ctx context.Context, projectID uuid.UUID) ([]*Task, error) {
	// Validate project exists
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	return s.repo.GetTasksByProject(ctx, projectID)
}

// BulkUpdateTasks updates multiple tasks with the same updates
func (s *service) BulkUpdateTasks(ctx context.Context, taskIDs []uuid.UUID, updates TaskUpdates) error {
	if len(taskIDs) == 0 {
		return nil // Nothing to update
	}

	// Validate updates
	if updates.State == nil && updates.Complexity == nil {
		return ValidationError{Field: "updates", Message: "at least one field must be specified for update"}
	}

	// Validate complexity if provided
	if updates.Complexity != nil && (*updates.Complexity < 1 || *updates.Complexity > 10) {
		return ValidationError{Field: "complexity", Message: "complexity must be between 1 and 10"}
	}

	// Update each task
	for _, taskID := range taskIDs {
		task, err := s.repo.GetTask(ctx, taskID)
		if err != nil {
			return fmt.Errorf("failed to get task %s: %w", taskID, err)
		}

		// Apply updates
		if updates.State != nil {
			task.State = *updates.State
			if task.State == TaskStateCompleted && task.CompletedAt == nil {
				now := time.Now()
				task.CompletedAt = &now
			} else if task.State != TaskStateCompleted && task.CompletedAt != nil {
				task.CompletedAt = nil
			}
		}
		if updates.Complexity != nil {
			task.Complexity = *updates.Complexity
		}

		task.UpdatedAt = time.Now()

		if err := s.repo.UpdateTask(ctx, task); err != nil {
			return fmt.Errorf("failed to update task %s: %w", taskID, err)
		}
	}

	return nil
}

// DuplicateTask creates a copy of a task in a new project
func (s *service) DuplicateTask(ctx context.Context, taskID uuid.UUID, newProjectID uuid.UUID) (*Task, error) {
	// Validate source task exists
	originalTask, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source task: %w", err)
	}

	// Validate target project exists
	if _, err := s.repo.GetProject(ctx, newProjectID); err != nil {
		return nil, fmt.Errorf("target project not found: %w", err)
	}

	// Create a copy of the task
	newTask := &Task{
		ID:          uuid.New(),
		ProjectID:   newProjectID,
		ParentID:    originalTask.ParentID, // This will be nil for the duplicated task
		Title:       originalTask.Title,
		Description: originalTask.Description,
		State:       TaskStatePending, // Reset state to pending
		Complexity:  originalTask.Complexity,
		Depth:       0, // Reset depth to 0 as it's now a root task
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CompletedAt: nil, // Reset completion status
	}

	// Save the new task
	if err := s.repo.CreateTask(ctx, newTask); err != nil {
		return nil, fmt.Errorf("failed to create duplicated task: %w", err)
	}

	return s.repo.GetTask(ctx, newTask.ID)
}

// SetTaskEstimate sets the time estimate for a task
func (s *service) SetTaskEstimate(ctx context.Context, taskID uuid.UUID, estimate int64) (*Task, error) {
	// Validate estimate
	if estimate < 0 {
		return nil, ValidationError{Field: "estimate", Message: "estimate must be non-negative"}
	}

	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task.Estimate = &estimate
	task.UpdatedAt = time.Now()

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task estimate: %w", err)
	}

	return s.repo.GetTask(ctx, taskID)
}

func (s *service) FindNextActionableTask(ctx context.Context, projectID uuid.UUID) (*Task, error) {
	// Validate project exists
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Get all tasks in the project
	allTasks, err := s.repo.GetTasksByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project tasks: %w", err)
	}

	// Create a map of task IDs to tasks for quick lookup
	taskMap := make(map[uuid.UUID]*Task)
	for _, task := range allTasks {
		taskMap[task.ID] = task
	}

	// Separate tasks by state
	var pendingTasks, inProgressTasks []*Task
	for _, task := range allTasks {
		switch task.State {
		case TaskStatePending:
			pendingTasks = append(pendingTasks, task)
		case TaskStateInProgress:
			inProgressTasks = append(inProgressTasks, task)
		}
	}

	// Prioritize in-progress tasks first
	if len(inProgressTasks) > 0 {
		// For in-progress tasks, find one that has all its dependencies met
		for _, task := range inProgressTasks {
			if s.areDependenciesMet(task, taskMap) {
				return task, nil
			}
		}
		// If no in-progress task has its dependencies met, this indicates an inconsistency
		// Since we prevent circular dependencies and should maintain data integrity,
		// this should not happen. We'll return an error to highlight the issue.
		return nil, fmt.Errorf("in-progress tasks exist but none have all dependencies met - possible data inconsistency")
	}

	// For pending tasks, find one that has all its dependencies met
	for _, task := range pendingTasks {
		if s.areDependenciesMet(task, taskMap) {
			return task, nil
		}
	}

	// If we reach here, it means either:
	// 1. There are no pending or in-progress tasks
	// 2. All pending tasks have unmet dependencies (potential deadlock scenario)
	// Since we prevent circular dependencies, case 2 suggests a logical error in task setup
	if len(pendingTasks) > 0 {
		return nil, fmt.Errorf("pending tasks exist but none have all dependencies met - possible deadlock scenario")
	}

	// No actionable tasks found
	return nil, fmt.Errorf("no actionable tasks found")
}

// areDependenciesMet checks if all dependencies of a task are completed
func (s *service) areDependenciesMet(task *Task, taskMap map[uuid.UUID]*Task) bool {
	for _, depID := range task.Dependencies {
		depTask, exists := taskMap[depID]
		if !exists || depTask.State != TaskStateCompleted {
			// If a dependency doesn't exist or isn't completed, the dependencies aren't met
			return false
		}
	}
	return true
}

func (s *service) FindTasksNeedingBreakdown(ctx context.Context, projectID uuid.UUID) ([]*Task, error) {
	// Validate project exists
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Find tasks with complexity above threshold that have no children
	tasks, err := s.repo.GetTasksByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project tasks: %w", err)
	}

	var needsBreakdown []*Task
	for _, task := range tasks {
		if task.Complexity >= s.config.ComplexityThreshold {
			// Check if task has children
			children, err := s.repo.GetTasksByParent(ctx, task.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to check task children: %w", err)
			}
			if len(children) == 0 {
				needsBreakdown = append(needsBreakdown, task)
			}
		}
	}

	return needsBreakdown, nil
}

func (s *service) GetProjectProgress(ctx context.Context, projectID uuid.UUID) (*ProjectProgress, error) {
	// Validate project exists
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	return s.repo.GetProjectProgress(ctx, projectID)
}

func (s *service) ListTasksByState(ctx context.Context, projectID uuid.UUID, state TaskState) ([]*Task, error) {
	// Validate project exists
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	return s.repo.ListTasks(ctx, TaskFilter{
		ProjectID: &projectID,
		State:     &state,
	})
}

// Validation helpers

func (s *service) validateProjectInput(title, description string) error {
	if title == "" {
		return ValidationError{Field: "title", Message: "title cannot be empty"}
	}
	if len(title) > 200 {
		return ValidationError{Field: "title", Message: "title cannot exceed 200 characters"}
	}
	if len(description) > s.config.MaxDescriptionLength {
		return ValidationError{Field: "description", Message: fmt.Sprintf("description cannot exceed %d characters", s.config.MaxDescriptionLength)}
	}
	return nil
}

func (s *service) validateTaskInput(title, description string, complexity int) error {
	if title == "" {
		return ValidationError{Field: "title", Message: "title cannot be empty"}
	}
	if len(title) > 200 {
		return ValidationError{Field: "title", Message: "title cannot exceed 200 characters"}
	}
	if len(description) > s.config.MaxDescriptionLength {
		return ValidationError{Field: "description", Message: fmt.Sprintf("description cannot exceed %d characters", s.config.MaxDescriptionLength)}
	}
	if complexity < 1 || complexity > 10 {
		return ValidationError{Field: "complexity", Message: "complexity must be between 1 and 10"}
	}
	return nil
}

// Config management

func (s *service) GetConfig() *Config {
	return s.config
}

func (s *service) UpdateConfig(config *Config) {
	if config != nil {
		s.config = config
	}
}
