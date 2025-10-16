package project

import (
	"context"
	"fmt"

	"github.com/denkhaus/knot/internal/manager"
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
