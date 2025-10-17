package task

import (
	"context"
	"fmt"
	"strings"

	"github.com/denkhaus/knot/internal/manager"
	"github.com/denkhaus/knot/internal/shared"

	"github.com/denkhaus/knot/internal/errors"
	"github.com/denkhaus/knot/internal/types"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// DeletionCommands returns task deletion related CLI commands
func DeletionCommands(appCtx *shared.AppContext) []*cli.Command {
	return []*cli.Command{
		{
			Name:   "delete",
			Usage:  "Delete a single task (only if no children exist)",
			Action: deleteAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Usage:    "Task ID to delete",
					Required: true,
				},
				&cli.BoolFlag{
					Name:  "force",
					Usage: "Skip confirmation prompt",
					Value: false,
				},
				&cli.BoolFlag{
					Name:  "dry-run",
					Usage: "Show what would be deleted without actually deleting",
					Value: false,
				},
			},
		},
		{
			Name:   "delete-subtree",
			Usage:  "Delete a task and all its descendants recursively",
			Action: deleteSubtreeAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Usage:    "Root task ID to delete (with all children)",
					Required: true,
				},
				&cli.BoolFlag{
					Name:  "force",
					Usage: "Skip confirmation prompt",
					Value: false,
				},
				&cli.BoolFlag{
					Name:  "dry-run",
					Usage: "Show what would be deleted without actually deleting",
					Value: false,
				},
			},
		},
	}
}

// deleteAction handles single task deletion
func deleteAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		taskIDStr := c.String("id")
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return errors.InvalidUUIDError("task-id", taskIDStr)
		}

		force := c.Bool("force")
		dryRun := c.Bool("dry-run")

		appCtx.Logger.Info("Deleting single task",
			zap.String("taskID", taskID.String()),
			zap.Bool("force", force),
			zap.Bool("dryRun", dryRun))

		// Get task details for validation and confirmation
		task, err := appCtx.ProjectManager.GetTask(context.Background(), taskID)
		if err != nil {
			appCtx.Logger.Error("Failed to get task", zap.Error(err))
			return errors.TaskNotFoundError(taskID)
		}

		// Check if task has children
		children, err := appCtx.ProjectManager.GetChildTasks(context.Background(), taskID)
		if err != nil {
			appCtx.Logger.Error("Failed to check for child tasks", zap.Error(err))
			return errors.WrapWithSuggestion(err, "checking child tasks")
		}

		if len(children) > 0 {
			return &errors.EnhancedError{
				Operation:   "deleting task",
				Cause:       fmt.Errorf("task has %d child task(s)", len(children)),
				Suggestion:  "Delete child tasks first, or use 'delete-subtree' to delete the entire hierarchy",
				Example:     fmt.Sprintf("knot task delete-subtree --id %s", taskID),
				HelpCommand: "knot task children --task-id " + taskID.String(),
			}
		}

		// Show what will be deleted
		fmt.Printf("Task to delete:\n")
		fmt.Printf("  • %s (ID: %s)\n", task.Title, task.ID)
		if task.Description != "" {
			fmt.Printf("    %s\n", task.Description)
		}
		fmt.Printf("    State: %s | Complexity: %d\n", task.State, task.Complexity)

		// Check for dependencies
		dependencies, err := appCtx.ProjectManager.GetTaskDependencies(context.Background(), taskID)
		if err == nil && len(dependencies) > 0 {
			fmt.Printf("\n  This task depends on %d other task(s):\n", len(dependencies))
			for _, dep := range dependencies {
				fmt.Printf("    • %s (ID: %s)\n", dep.Title, dep.ID)
			}
		}

		dependents, err := appCtx.ProjectManager.GetDependentTasks(context.Background(), taskID)
		if err == nil && len(dependents) > 0 {
			fmt.Printf("\n  %d task(s) depend on this task:\n", len(dependents))
			for _, dep := range dependents {
				fmt.Printf("    • %s (ID: %s)\n", dep.Title, dep.ID)
			}
			fmt.Printf("    These dependencies will be removed.\n")
		}

		if dryRun {
			fmt.Printf("\n DRY RUN: Task would be deleted (no actual changes made)\n")
			return nil
		}

		// Confirmation prompt
		if !force {
			if !confirmDeletion("task", task.Title) {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		// Perform deletion
		err = appCtx.ProjectManager.DeleteTask(context.Background(), taskID)
		if err != nil {
			appCtx.Logger.Error("Failed to delete task", zap.Error(err))
			return errors.WrapWithSuggestion(err, "deleting task")
		}

		appCtx.Logger.Info("Task deleted successfully")
		fmt.Printf("✅ Task deleted successfully: %s\n", task.Title)

		return nil
	}
}

// deleteSubtreeAction handles recursive task deletion
func deleteSubtreeAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		taskIDStr := c.String("id")
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return errors.InvalidUUIDError("task-id", taskIDStr)
		}

		force := c.Bool("force")
		dryRun := c.Bool("dry-run")

		appCtx.Logger.Info("Deleting task subtree",
			zap.String("taskID", taskID.String()),
			zap.Bool("force", force),
			zap.Bool("dryRun", dryRun))

		// Get task details
		task, err := appCtx.ProjectManager.GetTask(context.Background(), taskID)
		if err != nil {
			appCtx.Logger.Error("Failed to get task", zap.Error(err))
			return errors.TaskNotFoundError(taskID)
		}

		// Get all descendants for preview
		descendants, err := getTaskDescendants(appCtx.ProjectManager, taskID)
		if err != nil {
			appCtx.Logger.Error("Failed to get descendants", zap.Error(err))
			return errors.WrapWithSuggestion(err, "getting task descendants")
		}

		// Show what will be deleted
		fmt.Printf("Task subtree to delete:\n")
		fmt.Printf("  📁 %s (ID: %s) [ROOT]\n", task.Title, task.ID)

		if len(descendants) > 0 {
			fmt.Printf("  └── %d descendant task(s):\n", len(descendants))
			for _, desc := range descendants {
				indent := strings.Repeat("  ", desc.Depth-task.Depth+1)
				fmt.Printf("  %s├─ %s (ID: %s)\n", indent, desc.Title, desc.ID)
			}
		}

		totalTasks := 1 + len(descendants)
		fmt.Printf("\nTotal tasks to delete: %d\n", totalTasks)

		if dryRun {
			fmt.Printf("\n🔍 DRY RUN: %d task(s) would be deleted (no actual changes made)\n", totalTasks)
			return nil
		}

		// Confirmation prompt
		if !force {
			if !confirmDeletion("task subtree", fmt.Sprintf("%s and %d descendants", task.Title, len(descendants))) {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		// Perform deletion
		err = appCtx.ProjectManager.DeleteTaskSubtree(context.Background(), taskID)
		if err != nil {
			appCtx.Logger.Error("Failed to delete task subtree", zap.Error(err))
			return errors.WrapWithSuggestion(err, "deleting task subtree")
		}

		appCtx.Logger.Info("Task subtree deleted successfully", zap.Int("totalDeleted", totalTasks))
		fmt.Printf("✅ Task subtree deleted successfully: %d task(s) removed\n", totalTasks)

		return nil
	}
}

// confirmDeletion prompts user for confirmation
func confirmDeletion(itemType, itemName string) bool {
	fmt.Printf("\nAre you sure you want to delete this %s?\n", itemType)
	fmt.Printf("   %s\n", itemName)
	fmt.Printf("\nThis action cannot be undone. Type 'yes' to confirm: ")

	var response string
	fmt.Scanln(&response)

	return strings.ToLower(strings.TrimSpace(response)) == "yes"
}

// getTaskDescendants recursively gets all descendants of a task (renamed to avoid conflict)
func getTaskDescendants(projectManager manager.ProjectManager, taskID uuid.UUID) ([]*types.Task, error) {
	var result []*types.Task
	visited := make(map[uuid.UUID]bool)

	var collectDescendants func(uuid.UUID) error
	collectDescendants = func(id uuid.UUID) error {
		if visited[id] {
			return nil
		}
		visited[id] = true

		children, err := projectManager.GetChildTasks(context.Background(), id)
		if err != nil {
			return err
		}

		for _, child := range children {
			result = append(result, child)
			if err := collectDescendants(child.ID); err != nil {
				return err
			}
		}

		return nil
	}

	if err := collectDescendants(taskID); err != nil {
		return nil, err
	}

	return result, nil
}
