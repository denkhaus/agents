package project

import (
	"context"
	"fmt"

	"github.com/denkhaus/knot/internal/errors"
	"github.com/denkhaus/knot/internal/manager"
	"github.com/denkhaus/knot/internal/types"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// Commands returns all project-related CLI commands
func Commands(projectManager manager.ProjectManager, logger *zap.Logger) []*cli.Command {
	return []*cli.Command{
		{
			Name:   "create",
			Usage:  "Create a new project",
			Action: createAction(projectManager, logger),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "title",
					Aliases:  []string{"t"},
					Usage:    "Project title",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "description",
					Aliases: []string{"d"},
					Usage:   "Project description",
				},
			},
		},
		{
			Name:   "list",
			Usage:  "List all projects",
			Action: listAction(projectManager, logger),
		},
		{
			Name:   "get",
			Usage:  "Get project details",
			Action: getAction(projectManager, logger),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Usage:    "Project ID",
					Required: true,
				},
			},
		},
		{
			Name:   "delete",
			Usage:  "Delete a project with two-step confirmation",
			Action: deleteAction(projectManager, logger),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Usage:    "Project ID to delete",
					Required: true,
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

func createAction(projectManager manager.ProjectManager, logger *zap.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		title := c.String("title")
		description := c.String("description")

		logger.Info("Creating project", zap.String("title", title), zap.String("description", description))

		project, err := projectManager.CreateProject(context.Background(), title, description)
		if err != nil {
			logger.Error("Failed to create project", zap.Error(err))
			return fmt.Errorf("failed to create project: %w", err)
		}

		logger.Info("Project created successfully", zap.String("projectID", project.ID.String()), zap.String("title", project.Title))
		fmt.Printf("Created project: %s (ID: %s)\n", project.Title, project.ID)
		if project.Description != "" {
			fmt.Printf("  Description: %s\n", project.Description)
		}
		return nil
	}
}

func listAction(projectManager manager.ProjectManager, logger *zap.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		logger.Info("Listing projects")

		projects, err := projectManager.ListProjects(context.Background())
		if err != nil {
			logger.Error("Failed to list projects", zap.Error(err))
			return fmt.Errorf("failed to list projects: %w", err)
		}

		logger.Info("Projects retrieved", zap.Int("count", len(projects)))

		if len(projects) == 0 {
			fmt.Println("No projects found.")
			return nil
		}

		fmt.Printf("Found %d project(s):\n\n", len(projects))
		for _, project := range projects {
			fmt.Printf("• %s (ID: %s)\n", project.Title, project.ID)
			if project.Description != "" {
				fmt.Printf("  %s\n", project.Description)
			}
			fmt.Printf("  Progress: %.1f%% (%d/%d tasks completed)\n",
				project.Progress, project.CompletedTasks, project.TotalTasks)
			fmt.Println()
		}
		return nil
	}
}

func getAction(projectManager manager.ProjectManager, logger *zap.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		idStr := c.String("id")
		projectID, err := uuid.Parse(idStr)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}

		logger.Info("Getting project", zap.String("projectID", projectID.String()))

		project, err := projectManager.GetProject(context.Background(), projectID)
		if err != nil {
			logger.Error("Failed to get project", zap.Error(err))
			return fmt.Errorf("failed to get project: %w", err)
		}

		fmt.Printf("Project: %s\n", project.Title)
		fmt.Printf("ID: %s\n", project.ID)
		if project.Description != "" {
			fmt.Printf("Description: %s\n", project.Description)
		}
		fmt.Printf("Progress: %.1f%% (%d/%d tasks completed)\n",
			project.Progress, project.CompletedTasks, project.TotalTasks)
		fmt.Printf("Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", project.UpdatedAt.Format("2006-01-02 15:04:05"))

		return nil
	}
}

func deleteAction(projectManager manager.ProjectManager, logger *zap.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		projectIDStr := c.String("id")
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			return &errors.EnhancedError{
				Operation:   "parsing project ID",
				Cause:       err,
				Suggestion:  "Provide a valid UUID for the project ID",
				Example:     "knot project delete --id 550e8400-e29b-41d4-a716-446655440000",
				HelpCommand: "knot project delete --help",
			}
		}

		dryRun := c.Bool("dry-run")

		// Get project details
		project, err := projectManager.GetProject(context.Background(), projectID)
		if err != nil {
			return &errors.EnhancedError{
				Operation:   "retrieving project",
				Cause:       err,
				Suggestion:  "Verify the project ID exists",
				Example:     "knot project list",
				HelpCommand: "knot project get --help",
			}
		}

		// Check if project has tasks
		tasks, err := projectManager.ListTasksForProject(context.Background(), projectID)
		if err != nil {
			return &errors.EnhancedError{
				Operation:   "checking project tasks",
				Cause:       err,
				Suggestion:  "Unable to verify if project has tasks",
				HelpCommand: "knot task list --help",
			}
		}

		// Two-step deletion process
		if project.State == types.ProjectStateDeletionPending {
			// Second call - actually delete the project
			if dryRun {
				fmt.Printf("🔍 DRY RUN: Project would be permanently deleted (no actual changes made)\n")
				return nil
			}

			// Show what will be deleted
			fmt.Printf("🗑️  Final deletion of project:\n")
			fmt.Printf("  • %s (ID: %s)\n", project.Title, project.ID)
			if project.Description != "" {
				fmt.Printf("    %s\n", project.Description)
			}
			if len(tasks) > 0 {
				fmt.Printf("    ⚠️  This will also delete %d task(s)\n", len(tasks))
			}

			// Perform deletion
			err = projectManager.DeleteProject(context.Background(), projectID)
			if err != nil {
				return &errors.EnhancedError{
					Operation:   "deleting project",
					Cause:       err,
					Suggestion:  "Check if the project still exists or if there are constraint violations",
					HelpCommand: "knot project get --help",
				}
			}

			fmt.Printf("✅ Project permanently deleted: %s\n", project.Title)
			return nil
		} else {
			// First call - mark for deletion
			if dryRun {
				fmt.Printf("🔍 DRY RUN: Project would be marked for deletion (no actual changes made)\n")
				return nil
			}

			// Show what will be marked for deletion
			fmt.Printf("📋 Project to be marked for deletion:\n")
			fmt.Printf("  • %s (ID: %s)\n", project.Title, project.ID)
			if project.Description != "" {
				fmt.Printf("    %s\n", project.Description)
			}
			fmt.Printf("    Current State: %s\n", project.State)
			fmt.Printf("    Progress: %.1f%% (%d/%d tasks)\n", project.Progress, project.CompletedTasks, project.TotalTasks)

			if len(tasks) > 0 {
				fmt.Printf("\n  ⚠️  This project contains %d task(s):\n", len(tasks))
				for i, task := range tasks {
					if i < 5 { // Show first 5 tasks
						fmt.Printf("    • %s (%s)\n", task.Title, task.State)
					} else if i == 5 {
						fmt.Printf("    • ... and %d more task(s)\n", len(tasks)-5)
						break
					}
				}
				fmt.Printf("    All tasks will be deleted with the project.\n")
			}

			// Mark project for deletion
			_, err = projectManager.UpdateProjectState(context.Background(), projectID, types.ProjectStateDeletionPending)
			if err != nil {
				return &errors.EnhancedError{
					Operation:   "marking project for deletion",
					Cause:       err,
					Suggestion:  "Check if the project state transition is valid",
					HelpCommand: "knot project get --help",
				}
			}

			fmt.Printf("\n⚠️  Project marked for deletion. To confirm deletion, run the same command again:\n")
			fmt.Printf("    knot project delete --id %s\n", projectID)
			fmt.Printf("\n💡 To cancel deletion, change the project state:\n")
			fmt.Printf("    knot project update-state --id %s --state active\n", projectID)

			return nil
		}
	}
}
