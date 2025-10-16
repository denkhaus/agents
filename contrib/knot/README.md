# Knot

A standalone CLI tool for hierarchical project and task management with dependencies. Perfect for LLM agents to organize complex workflows.

## Features

- **Project Management**: Create, list, update, and delete projects
- **Task Management**: Create hierarchical tasks with parent-child relationships
- **Task Dependencies**: Manage task dependencies and blocking relationships
- **Local SQLite Storage**: Automatic .project directory with persistent SQLite database
- **Clean CLI Output**: No JSON logs in normal operation, debug mode available
- **LLM-Friendly**: Structured, parsable outputs perfect for AI agents

## Installation

```bash
go install github.com/denkhaus/knot/cmd/knot@latest
```

Or build locally:

```bash
go build -o knot ./cmd/knot
```

## Usage

### Project Commands

```bash
# Create a new project
knot project create --title "My Project" --description "Project description"

# List all projects
knot project list

# Get project details
knot project get --id <project-id>
```

### Task Commands

```bash
# Create a new task
knot task create --project-id <project-id> --title "Task title" --description "Task description"

# Create a subtask
knot task create --project-id <project-id> --parent-id <parent-task-id> --title "Subtask title"

# List tasks for a project
knot task list --project-id <project-id>

# Update task state
knot task update --id <task-id> --state completed
```

### Dependency Commands

```bash
# Add task dependency
knot dependency add --task-id <task-id> --depends-on <dependency-task-id>

# Remove dependency
knot dependency remove --task-id <task-id> --depends-on <dependency-task-id>

# List dependencies
knot dependency list --task-id <task-id>
```

## Configuration

Environment variables:

- `PM_LOG_LEVEL`: Logging level (`debug`, `warn`, `error`, `off`) - defaults to `error`

Knot automatically creates a `.project` directory in your current working directory to store the SQLite database with all your projects and tasks.

## For LLM Agents

Knot is designed for LLM agents with:

- Clean, parsable output
- Short, memorable commands
- Hierarchical task organization
- Dependency management
- Local persistence

Perfect for organizing complex AI workflows and project structures.
