package task

import (
	"context"
	"fmt"
	"strings"

	"github.com/denkhaus/knot/internal/shared"

	"github.com/denkhaus/knot/internal/types"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// bulk.go contains bulk operations on tasks
// - bulk-update: update multiple tasks at once
// - duplicate: duplicate tasks
// - state filtering and bulk operations

// TODO: Implement bulk operations
// REFERENCE: pkg/tools/project/main.go line 133 (bulkUpdateTasksTool)
// REFERENCE: pkg/tools/project/main.go line 134 (duplicateTaskTool)

// BulkUpdateAction updates multiple tasks simultaneously
func BulkUpdateAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		taskIDsStr := c.String("task-ids")
		if taskIDsStr == "" {
			return fmt.Errorf("task-ids are required")
		}

		// Parse comma-separated task IDs
		taskIDStrings := strings.Split(taskIDsStr, ",")
		var taskIDs []uuid.UUID
		for _, idStr := range taskIDStrings {
			idStr = strings.TrimSpace(idStr)
			taskID, err := uuid.Parse(idStr)
			if err != nil {
				return fmt.Errorf("invalid task ID '%s': %w", idStr, err)
			}
			taskIDs = append(taskIDs, taskID)
		}

		if len(taskIDs) == 0 {
			return fmt.Errorf("no valid task IDs provided")
		}

		// Build updates struct
		var updates types.TaskUpdates

		if stateStr := c.String("state"); stateStr != "" {
			state := types.TaskState(stateStr)
			updates.State = &state
		}

		if complexity := c.Int("complexity"); complexity > 0 {
			updates.Complexity = &complexity
		}

		if updates.State == nil && updates.Complexity == nil {
			return fmt.Errorf("at least one field (state or complexity) must be specified")
		}

		appCtx.Logger.Info("Bulk updating tasks",
			zap.Int("taskCount", len(taskIDs)),
			zap.Any("updates", updates))

		err := appCtx.ProjectManager.BulkUpdateTasks(context.Background(), taskIDs, updates)
		if err != nil {
			appCtx.Logger.Error("Failed to bulk update tasks", zap.Error(err))
			return fmt.Errorf("failed to bulk update tasks: %w", err)
		}

		fmt.Printf("Successfully updated %d tasks\n", len(taskIDs))
		if updates.State != nil {
			fmt.Printf("  State: %s\n", *updates.State)
		}
		if updates.Complexity != nil {
			fmt.Printf("  Complexity: %d\n", *updates.Complexity)
		}

		return nil
	}
}

// DuplicateAction creates a copy of a task
func DuplicateAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		taskIDStr := c.String("task-id")
		if taskIDStr == "" {
			return fmt.Errorf("task-id is required")
		}

		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return fmt.Errorf("invalid task ID: %w", err)
		}

		projectIDStr := c.String("target-project-id")
		if projectIDStr == "" {
			return fmt.Errorf("target-project-id is required")
		}

		targetProjectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			return fmt.Errorf("invalid target project ID: %w", err)
		}

		appCtx.Logger.Info("Duplicating task",
			zap.String("taskID", taskID.String()),
			zap.String("targetProjectID", targetProjectID.String()))

		duplicatedTask, err := appCtx.ProjectManager.DuplicateTask(context.Background(), taskID, targetProjectID)
		if err != nil {
			appCtx.Logger.Error("Failed to duplicate task", zap.Error(err))
			return fmt.Errorf("failed to duplicate task: %w", err)
		}

		fmt.Printf("Task duplicated successfully:\n")
		fmt.Printf("  Original: %s\n", taskID)
		fmt.Printf("  New: %s (ID: %s)\n", duplicatedTask.Title, duplicatedTask.ID)
		fmt.Printf("  Target Project: %s\n", targetProjectID)
		fmt.Printf("  State: %s (reset to pending)\n", duplicatedTask.State)
		fmt.Printf("  Complexity: %d\n", duplicatedTask.Complexity)

		return nil
	}
}

// ListByStateAction lists tasks filtered by state
func ListByStateAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		projectIDStr := c.String("project-id")
		if projectIDStr == "" {
			return fmt.Errorf("project-id is required")
		}

		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}

		stateStr := c.String("state")
		if stateStr == "" {
			return fmt.Errorf("state is required")
		}

		state := types.TaskState(stateStr)

		appCtx.Logger.Info("Listing tasks by state",
			zap.String("projectID", projectID.String()),
			zap.String("state", stateStr))

		tasks, err := appCtx.ProjectManager.ListTasksByState(context.Background(), projectID, state)
		if err != nil {
			appCtx.Logger.Error("Failed to list tasks by state", zap.Error(err))
			return fmt.Errorf("failed to list tasks by state: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Printf("No tasks found with state '%s' in project %s\n", state, projectID)
			return nil
		}

		// Check if JSON output is requested
		if c.Bool("json") {
			return outputTasksAsJSON(tasks)
		}

		fmt.Printf("Tasks with state '%s' (%d found):\n\n", state, len(tasks))
		for i, task := range tasks {
			fmt.Printf("%d. %s (ID: %s)\n", i+1, task.Title, task.ID)
			if task.Description != "" {
				fmt.Printf("   %s\n", task.Description)
			}
			fmt.Printf("   Complexity: %d | Depth: %d\n", task.Complexity, task.Depth)
			if task.ParentID != nil {
				fmt.Printf("   Parent: %s\n", *task.ParentID)
			}
			fmt.Println()
		}

		return nil
	}
}

// BulkCommands returns bulk operation CLI commands
func BulkCommands(appCtx *shared.AppContext) []*cli.Command {
	return []*cli.Command{
		{
			Name:   "bulk-update",
			Usage:  "Update multiple tasks simultaneously",
			Action: BulkUpdateAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "task-ids",
					Usage:    "Comma-separated list of task IDs",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "state",
					Usage: "New state (pending, in-progress, completed, blocked, cancelled)",
				},
				&cli.IntFlag{
					Name:  "complexity",
					Usage: "New complexity (1-10)",
				},
			},
		},
		{
			Name:   "duplicate",
			Usage:  "Duplicate a task to another project",
			Action: DuplicateAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "task-id",
					Usage:    "Task ID to duplicate",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "target-project-id",
					Usage:    "Target project ID",
					Required: true,
				},
			},
		},
		{
			Name:   "list-by-state",
			Usage:  "List tasks filtered by state",
			Action: ListByStateAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "project-id",
					Aliases:  []string{"p"},
					Usage:    "Project ID",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "state",
					Aliases:  []string{"s"},
					Usage:    "Task state (pending, in-progress, completed, blocked, cancelled)",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "json",
					Aliases: []string{"j"},
					Usage:   "Output in JSON format",
				},
			},
		},
	}
}
