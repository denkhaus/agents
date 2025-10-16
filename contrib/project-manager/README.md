# Project Manager CLI

A standalone CLI application for hierarchical project and task management, extracted from the agents framework.

## Features

- **Project Management**: Create, list, update, and delete projects
- **Task Management**: Create hierarchical tasks with parent-child relationships
- **Task Dependencies**: Manage task dependencies and blocking relationships
- **Multiple Storage Backends**: Support for in-memory and PostgreSQL storage
- **Task States**: Track task progress through different states (pending, in-progress, completed, blocked, cancelled)
- **Progress Tracking**: Monitor project completion and task distribution

## Installation

```bash
go build -o project-manager .
```

## Usage

### Project Commands

```bash
# Create a new project
./project-manager project create --title "My Project" --description "Project description"

# List all projects
./project-manager project list

# Get project details
./project-manager project get --id <project-id>
```

### Task Commands

```bash
# Create a new task
./project-manager task create --project-id <project-id> --title "Task title" --description "Task description"

# Create a subtask
./project-manager task create --project-id <project-id> --parent-id <parent-task-id> --title "Subtask title"

# List tasks for a project
./project-manager task list --project-id <project-id>
```

## Configuration

The application supports configuration through environment variables:

- `PM_DATABASE_URL`: PostgreSQL connection string (optional, defaults to in-memory storage)
- `PM_LOG_LEVEL`: Logging level (debug, info, warn, error)

## Development

This package is self-contained and has no dependencies on the main agents framework.