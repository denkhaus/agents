# Knot

A standalone CLI tool for hierarchical project and task management with dependencies. Perfect for LLM agents to organize complex workflows.

## Features

- **Project Management**: Create, list, update, and delete projects
- **Task Management**: Create hierarchical tasks with parent-child relationships
- **Task Dependencies**: Manage task dependencies and blocking relationships
- **Smart Complexity Management**: Auto-reduce parent task complexity when subtasks are added
- **Workflow Analysis**: Ready/blocked task views, actionable task recommendations, breakdown suggestions
- **Local SQLite Storage**: Automatic .knot directory with persistent SQLite database
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
knot task update-state --id <task-id> --state completed
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

### Workflow Commands

```bash
# Show tasks ready to work on (no blockers)
knot ready --project-id <project-id>

# Show tasks blocked by dependencies
knot blocked --project-id <project-id>

# Find next actionable task
knot actionable --project-id <project-id>

# Find tasks needing breakdown (high complexity)
knot breakdown --project-id <project-id>
```

### Hierarchy Commands

```bash
# Show task hierarchy as tree
knot task tree --project-id <project-id>

# Get direct children of a task
knot task children --task-id <task-id>

# Get parent task
knot task parent --task-id <task-id>

# Get root tasks (no parent)
knot task roots --project-id <project-id>
```

### Advanced Task Commands

```bash
# Delete single task (only if no children)
knot task delete --id <task-id>

# Delete task and all descendants
knot task delete-subtree --id <task-id>
```

## Configuration

### Environment Variables

- `PM_LOG_LEVEL`: Logging level (`debug`, `warn`, `error`, `off`) - defaults to `error`
- `KNOT_DEFAULT_COMPLEXITY`: Default complexity for new tasks (1-10) - defaults to `5`
- `KNOT_MAX_TASKS_PER_DEPTH`: Maximum tasks per hierarchy level - defaults to `100`

### Configuration Commands

```bash
# Show current configuration
knot config show

# Set configuration values
knot config set --key complexity-threshold --value 8
knot config set --key auto-reduce-complexity --value true

# Reset to defaults
knot config reset
```

### Smart Complexity Management

Knot features intelligent complexity management that automatically reduces parent task complexity when subtasks are added:

- **Problem Solved**: Tasks with high complexity disappear from `knot breakdown` even though they still have high complexity numbers
- **Logic**: When a task is broken into subtasks, the parent becomes a coordination task with lower complexity
- **Automatic Reduction**: Based on number of subtasks (1 subtask: -2 complexity, 2-3 subtasks: complexity 4, 4-5 subtasks: complexity 3, 6+ subtasks: complexity 2)
- **Configurable**: Can be disabled via `auto-reduce-complexity` setting
- **User-Friendly**: Shows feedback about complexity changes during task creation

### Storage

Knot automatically creates a `.knot` directory in your current working directory to store:
- SQLite database with all projects and tasks
- Configuration file (`config.json`)

## For LLM Agents

Knot is designed for LLM agents with:

- Clean, parsable output
- Short, memorable commands
- Hierarchical task organization
- Dependency management
- Local persistence

Perfect for organizing complex AI workflows and project structures.
