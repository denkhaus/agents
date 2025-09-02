package postgres

import (
	"context"
	"fmt"

	"github.com/denkhaus/agents/pkg/tools/project/repository/postgres/ent"
	"github.com/denkhaus/agents/pkg/tools/project/repository/postgres/ent/project"
	"github.com/denkhaus/agents/pkg/tools/project/repository/postgres/ent/task"
	"github.com/denkhaus/agents/pkg/tools/project/repository/postgres/ent/taskdependency"
	"github.com/denkhaus/agents/pkg/tools/project/shared"
	"github.com/google/uuid"
)

// Project CRUD Operations

func (r *postgresRepository) CreateProject(ctx context.Context, project *shared.Project) error {
	_, err := projectToEntProjectCreate(project, r.client).Save(ctx)
	if err != nil {
		return r.mapError("create project", err)
	}
	return nil
}

// GetProject retrieves a project by ID using ent
func (r *postgresRepository) GetProject(ctx context.Context, id uuid.UUID) (*shared.Project, error) {
	entProject, err := r.client.Project.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, NewNotFoundError("project", id.String())
		}
		return nil, r.mapError("get project", err)
	}
	return entProjectToProject(entProject), nil
}

// UpdateProject updates an existing project using ent
func (r *postgresRepository) UpdateProject(ctx context.Context, project *shared.Project) error {
	err := r.client.Project.UpdateOneID(project.ID).
		SetTitle(project.Title).
		SetDescription(project.Description).
		SetUpdatedAt(project.UpdatedAt).
		SetTotalTasks(project.TotalTasks).
		SetCompletedTasks(project.CompletedTasks).
		SetProgress(project.Progress).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return NewNotFoundError("project", project.ID.String())
		}
		return r.mapError("update project", err)
	}
	return nil
}

// DeleteProject deletes a project and all its tasks using ent transaction
func (r *postgresRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return r.withTx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		// First, delete all task dependencies for tasks in this project
		taskIDs, err := tx.Task.Query().
			Where(task.ProjectID(id)).
			IDs(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("failed to get task IDs for project: %w", err)
		}

		if len(taskIDs) > 0 {
			// Delete all task dependencies
			_, err = tx.TaskDependency.Delete().
				Where(taskdependency.Or(
					taskdependency.TaskIDIn(taskIDs...),
					taskdependency.DependsOnTaskIDIn(taskIDs...),
				)).
				Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to delete task dependencies: %w", err)
			}

			// Delete all tasks in the project
			_, err = tx.Task.Delete().
				Where(task.ProjectID(id)).
				Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to delete tasks: %w", err)
			}
		}

		// Delete the project
		err = tx.Project.DeleteOneID(id).Exec(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return NewNotFoundError("project", id.String())
			}
			return fmt.Errorf("failed to delete project: %w", err)
		}

		return nil
	})
}

// ListProjects retrieves all projects using ent
func (r *postgresRepository) ListProjects(ctx context.Context) ([]*shared.Project, error) {
	entProjects, err := r.client.Project.Query().
		Order(ent.Asc(project.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, r.mapError("list projects", err)
	}

	return entProjectsToProjects(entProjects), nil
}

// updateProjectMetrics updates project metrics (total tasks, completed tasks, progress)
func (r *postgresRepository) updateProjectMetrics(ctx context.Context, projectID uuid.UUID) error {
	// Get task counts by state using ent aggregation
	var totalTasks int
	var completedTasks int

	// Count total tasks
	totalTasks, err := r.client.Task.Query().
		Where(task.ProjectID(projectID)).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count total tasks: %w", err)
	}

	// Count completed tasks
	completedTasks, err = r.client.Task.Query().
		Where(
			task.ProjectID(projectID),
			task.StateEQ(task.StateCompleted),
		).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count completed tasks: %w", err)
	}

	// Calculate progress
	progress := 0.0
	if totalTasks > 0 {
		progress = float64(completedTasks) / float64(totalTasks) * 100.0
	}

	// Update project metrics
	err = r.client.Project.UpdateOneID(projectID).
		SetTotalTasks(totalTasks).
		SetCompletedTasks(completedTasks).
		SetProgress(progress).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update project metrics: %w", err)
	}

	return nil
}

// updateProjectMetricsInTx updates project metrics within a transaction
func (r *postgresRepository) updateProjectMetricsInTx(ctx context.Context, tx *ent.Tx, projectID uuid.UUID) error {
	// Get task counts by state using ent aggregation within transaction
	var totalTasks int
	var completedTasks int

	// Count total tasks
	totalTasks, err := tx.Task.Query().
		Where(task.ProjectID(projectID)).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count total tasks: %w", err)
	}

	// Count completed tasks
	completedTasks, err = tx.Task.Query().
		Where(
			task.ProjectID(projectID),
			task.StateEQ(task.StateCompleted),
		).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count completed tasks: %w", err)
	}

	// Calculate progress
	progress := 0.0
	if totalTasks > 0 {
		progress = float64(completedTasks) / float64(totalTasks) * 100.0
	}

	// Update project metrics within transaction
	err = tx.Project.UpdateOneID(projectID).
		SetTotalTasks(totalTasks).
		SetCompletedTasks(completedTasks).
		SetProgress(progress).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update project metrics: %w", err)
	}

	return nil
}
