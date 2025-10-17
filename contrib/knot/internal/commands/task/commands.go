package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/denkhaus/knot/internal/shared"

	"github.com/denkhaus/knot/internal/errors"
	"github.com/denkhaus/knot/internal/types"
	"github.com/denkhaus/knot/internal/validation"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// Commands returns all task-related CLI commands
func Commands(appCtx *shared.AppContext) []*cli.Command {
	// Basic task commands
	basicCommands := []*cli.Command{
		{
			Name:   "create",
			Usage:  "Create a new task",
			Action: createAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "project-id",
					Aliases:  []string{"p"},
					Usage:    "Project ID",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "title",
					Aliases:  []string{"t"},
					Usage:    "Task title",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "description",
					Aliases: []string{"d"},
					Usage:   "Task description",
				},
				&cli.StringFlag{
					Name:  "parent-id",
					Usage: "Parent task ID (for subtasks)",
				},
				&cli.IntFlag{
					Name:    "complexity",
					Aliases: []string{"c"},
					Usage:   "Task complexity (1-10)",
					Value:   5,
					EnvVars: []string{"KNOT_DEFAULT_COMPLEXITY"},
				},
			},
		},
		{
			Name:   "list",
			Usage:  "List tasks",
			Action: listAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "project-id",
					Aliases:  []string{"p"},
					Usage:    "Project ID",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "json",
					Aliases: []string{"j"},
					Usage:   "Output in JSON format",
				},
			},
		},
		{
			Name:   "update-state",
			Usage:  "Update task state",
			Action: updateStateAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Usage:    "Task ID",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "state",
					Aliases:  []string{"s"},
					Usage:    "New state (pending, in-progress, completed, blocked, cancelled)",
					Required: true,
				},
			},
		},
		{
			Name:   "update-title",
			Usage:  "Update task title",
			Action: updateTitleAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Usage:    "Task ID",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "title",
					Aliases:  []string{"t"},
					Usage:    "New task title",
					Required: true,
				},
			},
		},
		{
			Name:   "update-description",
			Usage:  "Update task description",
			Action: updateDescriptionAction(appCtx),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Usage:    "Task ID",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "description",
					Aliases:  []string{"d"},
					Usage:    "New task description",
					Required: true,
				},
			},
		},
	}

	// Hierarchy navigation commands
	hierarchyCommands := HierarchyCommands(appCtx)

	// Task deletion commands
	deletionCommands := DeletionCommands(appCtx)

	// Bulk operation commands
	bulkCommands := BulkCommands(appCtx)

	// Combine all commands
	allCommands := make([]*cli.Command, 0, len(basicCommands)+len(hierarchyCommands)+len(deletionCommands)+len(bulkCommands))
	allCommands = append(allCommands, basicCommands...)
	allCommands = append(allCommands, hierarchyCommands...)
	allCommands = append(allCommands, deletionCommands...)
	allCommands = append(allCommands, bulkCommands...)

	return allCommands
}

func createAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		projectIDStr := c.String("project-id")
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}

		title := c.String("title")
		description := c.String("description")
		complexity := c.Int("complexity")

		// Validate complexity
		if err := errors.ValidateComplexity(complexity); err != nil {
			return err
		}

		var parentID *uuid.UUID
		if parentIDStr := c.String("parent-id"); parentIDStr != "" {
			parsed, err := uuid.Parse(parentIDStr)
			if err != nil {
				return errors.InvalidUUIDError("parent-id", parentIDStr)
			}
			parentID = &parsed
		}

		appCtx.Logger.Info("Creating task",
			zap.String("title", title),
			zap.String("projectID", projectID.String()),
			zap.Int("complexity", complexity))

		task, err := appCtx.ProjectManager.CreateTask(context.Background(), projectID, parentID, title, description, complexity)
		if err != nil {
			appCtx.Logger.Error("Failed to create task", zap.Error(err))
			return errors.WrapWithSuggestion(err, "creating task")
		}

		appCtx.Logger.Info("Task created successfully", zap.String("taskID", task.ID.String()))

		fmt.Printf("Created task: %s (ID: %s)\n", task.Title, task.ID)
		if task.Description != "" {
			fmt.Printf("  Description: %s\n", task.Description)
		}
		fmt.Printf("  Complexity: %d\n", task.Complexity)
		fmt.Printf("  State: %s\n", task.State)
		if parentID != nil {
			fmt.Printf("  Parent: %s\n", *parentID)
		}

		// Show breakdown suggestion for high complexity tasks
		if complexity >= 8 {
			fmt.Printf("\nNote: This task has high complexity (%d >= 8 threshold).\n", complexity)
			fmt.Printf("Consider breaking it down into smaller subtasks:\n")
			fmt.Printf("  knot task create --project-id %s --parent-id %s --title \"Subtask 1\"\n", projectID, task.ID)
			fmt.Printf("  knot breakdown --project-id %s  # to see all tasks needing breakdown\n", projectID)
		}

		return nil
	}
}

func listAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		projectIDStr := c.String("project-id")
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}

		appCtx.Logger.Info("Listing tasks", zap.String("projectID", projectID.String()))

		tasks, err := appCtx.ProjectManager.ListTasksForProject(context.Background(), projectID)
		if err != nil {
			appCtx.Logger.Error("Failed to list tasks", zap.Error(err))
			return errors.WrapWithSuggestion(err, "listing tasks")
		}

		appCtx.Logger.Info("Tasks retrieved", zap.Int("count", len(tasks)))

		if len(tasks) == 0 {
			return errors.EmptyResultError("list tasks", fmt.Sprintf("project %s", projectID))
		}

		// Check if JSON output is requested
		if c.Bool("json") {
			return outputTasksAsJSON(tasks)
		}

		fmt.Printf("Found %d task(s):\n\n", len(tasks))
		for _, task := range tasks {
			indent := ""
			for i := 0; i < task.Depth; i++ {
				indent += "  "
			}
			fmt.Printf("%s* %s (ID: %s)\n", indent, task.Title, task.ID)
			if task.Description != "" {
				fmt.Printf("%s  %s\n", indent, task.Description)
			}
			fmt.Printf("%s  State: %s | Complexity: %d\n", indent, task.State, task.Complexity)
			fmt.Println()
		}
		return nil
	}
}

func updateStateAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		taskIDStr := c.String("id")
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return errors.InvalidUUIDError("task-id", taskIDStr)
		}

		stateStr := c.String("state")

		// Basic state validation
		if err := errors.ValidateTaskState(stateStr); err != nil {
			return err
		}

		newState := types.TaskState(stateStr)

		appCtx.Logger.Info("Updating task state",
			zap.String("taskID", taskID.String()),
			zap.String("newState", stateStr))

		// Get current task to preserve other fields
		task, err := appCtx.ProjectManager.GetTask(context.Background(), taskID)
		if err != nil {
			appCtx.Logger.Error("Failed to get task", zap.Error(err))
			return errors.TaskNotFoundError(taskID)
		}

		// Validate state transition
		validator := validation.NewStateValidator()
		if err := validator.ValidateTransition(task.State, newState, task); err != nil {
			appCtx.Logger.Error("Invalid state transition",
				zap.String("from", string(task.State)),
				zap.String("to", string(newState)),
				zap.Error(err))
			return err
		}

		// Update task state
		updatedTask, err := appCtx.ProjectManager.UpdateTaskState(context.Background(), taskID, newState)
		if err != nil {
			appCtx.Logger.Error("Failed to update task state", zap.Error(err))
			return errors.WrapWithSuggestion(err, "updating task state")
		}

		appCtx.Logger.Info("Task state updated successfully")
		fmt.Printf("Updated task state: %s -> %s\n", task.State, updatedTask.State)
		return nil
	}
}

// outputTasksAsJSON outputs tasks in JSON format
func outputTasksAsJSON(tasks []*types.Task) error {
	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks to JSON: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

// outputSingleTaskAsJSON outputs a single task in JSON format
func outputSingleTaskAsJSON(task *types.Task) error {
	jsonData, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task to JSON: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

func updateTitleAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		taskIDStr := c.String("id")
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return errors.InvalidUUIDError("task-id", taskIDStr)
		}

		newTitle := c.String("title")
		if newTitle == "" {
			return fmt.Errorf("title cannot be empty")
		}

		appCtx.Logger.Info("Updating task title",
			zap.String("taskID", taskID.String()),
			zap.String("newTitle", newTitle))

		// Get current task to check if it exists and get old title
		task, err := appCtx.ProjectManager.GetTask(context.Background(), taskID)
		if err != nil {
			appCtx.Logger.Error("Failed to get task", zap.Error(err))
			return errors.TaskNotFoundError(taskID)
		}

		oldTitle := task.Title

		// Update task title
		updatedTask, err := appCtx.ProjectManager.UpdateTaskTitle(context.Background(), taskID, newTitle)
		if err != nil {
			appCtx.Logger.Error("Failed to update task title", zap.Error(err))
			return errors.WrapWithSuggestion(err, "updating task title")
		}

		appCtx.Logger.Info("Task title updated successfully")
		fmt.Printf("Updated task title: \"%s\" -> \"%s\"\n", oldTitle, updatedTask.Title)
		return nil
	}
}

func updateDescriptionAction(appCtx *shared.AppContext) cli.ActionFunc {
	return func(c *cli.Context) error {
		taskIDStr := c.String("id")
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return errors.InvalidUUIDError("task-id", taskIDStr)
		}

		newDescription := c.String("description")

		appCtx.Logger.Info("Updating task description",
			zap.String("taskID", taskID.String()),
			zap.String("newDescription", newDescription))

		// Get current task to check if it exists and get old description
		task, err := appCtx.ProjectManager.GetTask(context.Background(), taskID)
		if err != nil {
			appCtx.Logger.Error("Failed to get task", zap.Error(err))
			return errors.TaskNotFoundError(taskID)
		}

		oldDescription := task.Description

		// Update task description
		updatedTask, err := appCtx.ProjectManager.UpdateTaskDescription(context.Background(), taskID, newDescription)
		if err != nil {
			appCtx.Logger.Error("Failed to update task description", zap.Error(err))
			return errors.WrapWithSuggestion(err, "updating task description")
		}

		appCtx.Logger.Info("Task description updated successfully")
		if oldDescription == "" {
			fmt.Printf("Updated task description: (empty) -> \"%s\"\n", updatedTask.Description)
		} else {
			fmt.Printf("Updated task description: \"%s\" -> \"%s\"\n", oldDescription, updatedTask.Description)
		}
		return nil
	}
}
