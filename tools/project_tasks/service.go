package projecttasks

import (
	"context"
	"fmt"

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

func (s *service) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	return s.repo.DeleteProject(ctx, projectID)
}

func (s *service) ListProjects(ctx context.Context) ([]*Project, error) {
	return s.repo.ListProjects(ctx)
}

// Task operations

func (s *service) CreateTask(ctx context.Context, projectID uuid.UUID, parentID *uuid.UUID, title, description string, complexity, priority int) (*Task, error) {
	if err := s.validateTaskInput(title, description, complexity, priority); err != nil {
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
	if maxTasks, exists := s.config.MaxTasksPerDepth[depth]; exists {
		counts, err := s.repo.GetTaskCountByDepth(ctx, projectID, depth)
		if err != nil {
			return nil, fmt.Errorf("failed to check task count constraints: %w", err)
		}
		if counts[depth] >= maxTasks {
			return nil, fmt.Errorf("maximum tasks per depth (%d) exceeded for depth %d", maxTasks, depth)
		}
	}

	task := &Task{
		ID:          uuid.New(),
		ProjectID:   projectID,
		ParentID:    parentID,
		Title:       title,
		Description: description,
		State:       TaskStatePending,
		Complexity:  complexity,
		Priority:    priority,
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

func (s *service) UpdateTask(ctx context.Context, taskID uuid.UUID, title, description string, complexity, priority int, state TaskState) (*Task, error) {
	if err := s.validateTaskInput(title, description, complexity, priority); err != nil {
		return nil, err
	}

	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task.Title = title
	task.Description = description
	task.Complexity = complexity
	task.Priority = priority
	task.State = state

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
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

func (s *service) ListTasksHierarchical(ctx context.Context, projectID uuid.UUID) ([]*TaskHierarchy, error) {
	// Validate project exists
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	return s.repo.GetTaskHierarchy(ctx, projectID)
}

func (s *service) GetTaskSubtree(ctx context.Context, taskID uuid.UUID) (*TaskHierarchy, error) {
	return s.repo.GetTaskSubtree(ctx, taskID)
}

func (s *service) FindNextActionableTask(ctx context.Context, projectID uuid.UUID) (*Task, error) {
	// Validate project exists
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Get all pending and in-progress tasks, sorted by priority
	pendingState := TaskStatePending
	inProgressState := TaskStateInProgress
	
	pendingTasks, err := s.repo.ListTasks(ctx, TaskFilter{
		ProjectID: &projectID,
		State:     &pendingState,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pending tasks: %w", err)
	}

	inProgressTasks, err := s.repo.ListTasks(ctx, TaskFilter{
		ProjectID: &projectID,
		State:     &inProgressState,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get in-progress tasks: %w", err)
	}

	// Prioritize in-progress tasks first, then pending
	if len(inProgressTasks) > 0 {
		return inProgressTasks[0], nil
	}
	if len(pendingTasks) > 0 {
		return pendingTasks[0], nil
	}

	return nil, fmt.Errorf("no actionable tasks found")
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
	if len(description) > 2000 {
		return ValidationError{Field: "description", Message: "description cannot exceed 2000 characters"}
	}
	return nil
}

func (s *service) validateTaskInput(title, description string, complexity, priority int) error {
	if title == "" {
		return ValidationError{Field: "title", Message: "title cannot be empty"}
	}
	if len(title) > 200 {
		return ValidationError{Field: "title", Message: "title cannot exceed 200 characters"}
	}
	if len(description) > 2000 {
		return ValidationError{Field: "description", Message: "description cannot exceed 2000 characters"}
	}
	if complexity < 1 || complexity > 10 {
		return ValidationError{Field: "complexity", Message: "complexity must be between 1 and 10"}
	}
	if priority < 1 || priority > 10 {
		return ValidationError{Field: "priority", Message: "priority must be between 1 and 10"}
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