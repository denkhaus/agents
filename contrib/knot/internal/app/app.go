package app

import (
	configCommands "github.com/denkhaus/knot/internal/commands/config"
	"github.com/denkhaus/knot/internal/commands/dependency"
	"github.com/denkhaus/knot/internal/commands/health"
	"github.com/denkhaus/knot/internal/commands/project"
	"github.com/denkhaus/knot/internal/commands/task"
	validationCommands "github.com/denkhaus/knot/internal/commands/validation"
	"github.com/denkhaus/knot/internal/logger"
	"github.com/denkhaus/knot/internal/manager"
	"github.com/denkhaus/knot/internal/repository/inmemory"
	"github.com/denkhaus/knot/internal/repository/sqlite"
	"github.com/denkhaus/knot/internal/types"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// App represents the CLI application
type App struct {
	*cli.App
	context *Context
}

// New creates a new CLI application with all dependencies initialized
func New() (*App, error) {
	// Initialize logger
	appLogger := logger.GetLogger()
	
	// Initialize repository (SQLite with fallback to in-memory)
	var repo types.Repository
	var err error
	
	repo, err = sqlite.NewRepository(
		sqlite.WithLogger(appLogger),
		sqlite.WithAutoMigrate(true),
	)
	if err != nil {
		appLogger.Warn("Failed to initialize SQLite repository, falling back to in-memory", zap.Error(err))
		repo = inmemory.NewMemoryRepository()
	} else {
		appLogger.Info("SQLite repository initialized successfully")
	}
	
	// Initialize project manager
	config := manager.DefaultConfig()
	projectManager := manager.NewManagerWithRepository(repo, config)
	
	// Create application context
	appCtx := NewContext(projectManager, appLogger)
	
	// Create CLI app
	cliApp := &cli.App{
		Name:    "knot",
		Usage:   "A CLI tool for hierarchical project and task management with dependencies",
		Version: "1.0.0",
		Authors: []*cli.Author{
			{
				Name:  "denkhaus",
				Email: "denkhaus@example.com",
			},
		},
		Before: func(c *cli.Context) error {
			appLogger.Info("Knot CLI started", zap.String("version", "1.0.0"))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:        "project",
				Aliases:     []string{"p"},
				Usage:       "Project management commands",
				Subcommands: project.Commands(projectManager, appLogger),
			},
			{
				Name:        "task",
				Aliases:     []string{"t"},
				Usage:       "Task management commands",
				Subcommands: task.Commands(projectManager, appLogger),
			},
			{
				Name:        "dependency",
				Aliases:     []string{"dep"},
				Usage:       "Task dependency management",
				Subcommands: dependency.Commands(projectManager, appLogger),
			},
			{
				Name:        "config",
				Aliases:     []string{"cfg"},
				Usage:       "Configuration management",
				Subcommands: configCommands.Commands(projectManager, appLogger),
			},
			{
				Name:        "health",
				Usage:       "Database health and connectivity checks",
				Subcommands: health.Commands(projectManager, appLogger),
			},
			{
				Name:        "validate",
				Usage:       "Task state validation and transition checks",
				Subcommands: validationCommands.Commands(projectManager, appLogger),
			},
			{
				Name:    "ready",
				Usage:   "Show tasks with no blockers (ready to work on)",
				Action:  task.ReadyAction(projectManager, appLogger),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "project-id",
						Aliases:  []string{"p"},
						Usage:    "Project ID",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "Maximum number of tasks to show (default: 10)",
						Value:   10,
						EnvVars: []string{"KNOT_TASK_LIMIT"},
					},
				},
			},
			{
				Name:    "blocked",
				Usage:   "Show tasks blocked by dependencies",
				Action:  task.BlockedAction(projectManager, appLogger),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "project-id",
						Aliases:  []string{"p"},
						Usage:    "Project ID",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "Maximum number of tasks to show (default: 10)",
						Value:   10,
						EnvVars: []string{"KNOT_TASK_LIMIT"},
					},
				},
			},
			{
				Name:    "actionable",
				Usage:   "Find the next actionable task in a project",
				Action:  task.ActionableAction(projectManager, appLogger),
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
				Name:    "breakdown",
				Usage:   "Find tasks that need breakdown based on complexity",
				Action:  task.BreakdownAction(projectManager, appLogger),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "project-id",
						Aliases:  []string{"p"},
						Usage:    "Project ID",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "threshold",
						Aliases: []string{"t"},
						Usage:   "Complexity threshold for breakdown (default: 8)",
						Value:   8,
						EnvVars: []string{"KNOT_COMPLEXITY_THRESHOLD"},
					},
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "Maximum number of tasks to show (default: 10)",
						Value:   10,
						EnvVars: []string{"KNOT_TASK_LIMIT"},
					},
				},
			},
		},
	}
	
	return &App{
		App:     cliApp,
		context: appCtx,
	}, nil
}

// Run starts the CLI application
func (a *App) Run(args []string) error {
	defer logger.Sync()
	
	if err := a.App.Run(args); err != nil {
		a.context.Logger.Error("Application error", zap.Error(err))
		return err
	}
	
	return nil
}