package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/denkhaus/project-manager/internal/manager"
	"github.com/denkhaus/project-manager/internal/repository/inmemory"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

var projectManager manager.ProjectManager

func main() {
	// Initialize the project manager with in-memory repository
	repo := inmemory.NewMemoryRepository()
	config := manager.DefaultConfig()
	projectManager = manager.NewManagerWithRepository(repo, config)

	app := &cli.App{
		Name:    "project-manager",
		Usage:   "A CLI tool for hierarchical project and task management",
		Version: "1.0.0",
		Authors: []*cli.Author{
			{
				Name:  "denkhaus",
				Email: "denkhaus@example.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "project",
				Aliases: []string{"p"},
				Usage:   "Project management commands",
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new project",
						Action: createProject,
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
						Action: listProjects,
					},
					{
						Name:   "get",
						Usage:  "Get project details",
						Action: getProject,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "id",
								Usage:    "Project ID",
								Required: true,
							},
						},
					},
				},
			},
			{
				Name:    "task",
				Aliases: []string{"t"},
				Usage:   "Task management commands",
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new task",
						Action: createTask,
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
							},
						},
					},
					{
						Name:   "list",
						Usage:  "List tasks",
						Action: listTasks,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "project-id",
								Aliases:  []string{"p"},
								Usage:    "Project ID",
								Required: true,
							},
						},
					},
					{
						Name:   "update",
						Usage:  "Update task state",
						Action: updateTaskState,
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
				},
			},
			{
				Name:    "dependency",
				Aliases: []string{"dep"},
				Usage:   "Task dependency management",
				Subcommands: []*cli.Command{
					{
						Name:   "add",
						Usage:  "Add task dependency",
						Action: addDependency,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "task-id",
								Usage:    "Task ID",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "depends-on",
								Usage:    "Task ID that this task depends on",
								Required: true,
							},
						},
					},
					{
						Name:   "remove",
						Usage:  "Remove task dependency",
						Action: removeDependency,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "task-id",
								Usage:    "Task ID",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "depends-on",
								Usage:    "Task ID to remove dependency from",
								Required: true,
							},
						},
					},
					{
						Name:   "list",
						Usage:  "List task dependencies",
						Action: listDependencies,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "task-id",
								Usage:    "Task ID",
								Required: true,
							},
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func createProject(c *cli.Context) error {
	title := c.String("title")
	description := c.String("description")

	project, err := projectManager.CreateProject(context.Background(), title, description)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	fmt.Printf("Created project: %s (ID: %s)\n", project.Title, project.ID)
	if project.Description != "" {
		fmt.Printf("  Description: %s\n", project.Description)
	}
	return nil
}

func listProjects(c *cli.Context) error {
	projects, err := projectManager.ListProjects(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

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

func getProject(c *cli.Context) error {
	idStr := c.String("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid project ID: %w", err)
	}

	project, err := projectManager.GetProject(context.Background(), projectID)
	if err != nil {
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

func createTask(c *cli.Context) error {
	projectIDStr := c.String("project-id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return fmt.Errorf("invalid project ID: %w", err)
	}

	title := c.String("title")
	description := c.String("description")
	complexity := c.Int("complexity")

	var parentID *uuid.UUID
	if parentIDStr := c.String("parent-id"); parentIDStr != "" {
		parsed, err := uuid.Parse(parentIDStr)
		if err != nil {
			return fmt.Errorf("invalid parent ID: %w", err)
		}
		parentID = &parsed
	}

	task, err := projectManager.CreateTask(context.Background(), projectID, parentID, title, description, complexity)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	fmt.Printf("Created task: %s (ID: %s)\n", task.Title, task.ID)
	if task.Description != "" {
		fmt.Printf("  Description: %s\n", task.Description)
	}
	fmt.Printf("  Complexity: %d\n", task.Complexity)
	fmt.Printf("  State: %s\n", task.State)
	if parentID != nil {
		fmt.Printf("  Parent: %s\n", *parentID)
	}
	return nil
}

func listTasks(c *cli.Context) error {
	projectIDStr := c.String("project-id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return fmt.Errorf("invalid project ID: %w", err)
	}

	tasks, err := projectManager.ListTasksForProject(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found for this project.")
		return nil
	}

	fmt.Printf("Found %d task(s):\n\n", len(tasks))
	for _, task := range tasks {
		indent := ""
		for i := 0; i < task.Depth; i++ {
			indent += "  "
		}
		fmt.Printf("%s• %s (ID: %s)\n", indent, task.Title, task.ID)
		if task.Description != "" {
			fmt.Printf("%s  %s\n", indent, task.Description)
		}
		fmt.Printf("%s  State: %s | Complexity: %d\n", indent, task.State, task.Complexity)
		fmt.Println()
	}
	return nil
}