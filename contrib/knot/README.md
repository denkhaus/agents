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
- **Audit Trail**: Track changes with actor information using `--actor` flag
- **Bulk Operations**: Update multiple tasks simultaneously, duplicate tasks
- **Health Checks**: Database connectivity and validation tools
- **State Validation**: Task state validation and transition checks
- **Comprehensive Dependency Analysis**: Dependency chains, cycles detection, validation

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

# Create a project with actor tracking
knot project create --title "My Project" --actor "john.doe"

# List all projects
knot project list

# Get project details
knot project get --id <project-id>

# Delete a project (two-step confirmation)
knot project delete --id <project-id>
```

### Task Commands

```bash
# Create a new task
knot task create --project-id <project-id> --title "Task title" --description "Task description"

# Create a task with actor tracking
knot task create --project-id <project-id> --title "Task title" --actor "jane.smith"

# Create a subtask
knot task create --project-id <project-id> --parent-id <parent-task-id> --title "Subtask title"

# List tasks for a project
knot task list --project-id <project-id>

# Update task state
knot task update-state --id <task-id> --state completed

# Update task state with actor tracking
knot task update-state --id <task-id> --state in-progress --actor "dev.bot"

# Update task title or description
knot task update-title --id <task-id> --title "New title"
knot task update-description --id <task-id> --description "New description"

# Update task with actor tracking
knot task update-title --id <task-id> --title "New title" --actor "admin"
knot task update-description --id <task-id> --description "New description" --actor "admin"

# Delete single task (only if no children)
knot task delete --id <task-id>

# Delete task and all descendants
knot task delete-subtree --id <task-id>
```

### Task Hierarchy Commands

```bash
# Show task hierarchy as tree
knot task tree --project-id <project-id>

# Show tree starting from specific task
knot task tree --project-id <project-id> --root-task-id <task-id>

# Show tree with depth limit
knot task tree --project-id <project-id> --max-depth 3

# Get direct children of a task
knot task children --task-id <task-id>

# Get children recursively
knot task children --task-id <task-id> --recursive

# Get parent task
knot task parent --task-id <task-id>

# Get root tasks (no parent)
knot task roots --project-id <project-id>

# Get root tasks with limit
knot task roots --project-id <project-id> --limit 10
```

### Task Analysis Commands

```bash
# Show tasks ready to work on (no blockers)
knot ready --project-id <project-id>

# Show tasks ready to work on with limit
knot ready --project-id <project-id> --limit 5

# Show tasks blocked by dependencies
knot blocked --project-id <project-id>

# Show blocked tasks with limit
knot blocked --project-id <project-id> --limit 5

# Find next actionable task
knot actionable --project-id <project-id>

# Find tasks needing breakdown (high complexity)
knot breakdown --project-id <project-id>

# Find tasks needing breakdown with custom threshold
knot breakdown --project-id <project-id> --threshold 7

# List tasks filtered by state
knot task list-by-state --project-id <project-id> --state "in-progress"
```

### Task Bulk Operations

```bash
# Bulk update multiple tasks
knot task bulk-update --task-ids "task-id-1,task-id-2,task-id-3" --state completed --complexity 3

# Duplicate a task to another project
knot task duplicate --task-id <task-id> --target-project-id <project-id>
```

### Dependency Commands

```bash
# Add task dependency
knot dependency add --task-id <task-id> --depends-on <dependency-task-id>

# Add task dependency with actor tracking
knot dependency add --task-id <task-id> --depends-on <dependency-task-id> --actor "project.manager"

# Remove dependency
knot dependency remove --task-id <task-id> --depends-on <dependency-task-id>

# Remove dependency with actor tracking
knot dependency remove --task-id <task-id> --depends-on <dependency-task-id> --actor "project.manager"

# List dependencies
knot dependency list --task-id <task-id>

# List tasks that depend on this task
knot dependency dependents --task-id <task-id>

# Show dependency chain for a task (both upstream and downstream)
knot dependency chain --task-id <task-id> --upstream --downstream

# Show dependency chain for a task (upstream only)
knot dependency chain --task-id <task-id> --upstream

# Show dependency chain for a task (downstream only)
knot dependency chain --task-id <task-id> --downstream

# Detect circular dependencies in project
knot dependency cycles --project-id <project-id>

# Validate all dependencies in project
knot dependency validate --project-id <project-id>
```

### Configuration Commands

```bash
# Show current configuration
knot config show

# Set configuration values
knot config set --key complexity-threshold --value 8
knot config set --key auto-reduce-complexity --value 1

# Reset to defaults
knot config reset
```

### Health and Validation Commands

```bash
# Check database health
knot health check

# Check database health with JSON output
knot health check --json

# Simple database connectivity test
knot health ping

# Comprehensive database validation
knot health validate

# Show valid task states and transitions
knot validate states

# Show complete transition matrix
knot validate states --matrix

# Show valid transitions from specific state
knot validate states --from "pending"

# Validate a state transition without applying it
knot validate transition --task-id <task-id> --to completed

# Validate all task states in a project
knot validate project --project-id <project-id>

# Validate all task states and attempt to fix invalid ones
knot validate project --project-id <project-id> --fix
```

## Configuration

### Environment Variables

- `PM_LOG_LEVEL`: Logging level (`debug`, `warn`, `error`, `off`) - defaults to `error`
- `KNOT_DEFAULT_COMPLEXITY`: Default complexity for new tasks (1-10) - defaults to `5`
- `KNOT_MAX_TASKS_PER_DEPTH`: Maximum tasks per hierarchy level - defaults to `100`
- `KNOT_TASK_LIMIT`: Maximum number of tasks to show in commands - defaults to `10`
- `KNOT_COMPLEXITY_THRESHOLD`: Threshold for breakdown suggestions (default: `8`)
- `KNOT_ACTOR`: Default actor name for audit trail (default: $USER)

### Configuration Commands

```bash
# Show current configuration
knot config show

# Set configuration values (0/1 for boolean values)
knot config set --key complexity-threshold --value 8
knot config set --key auto-reduce-complexity --value 1

# Reset to defaults
knot config reset
```

### Actor/Audit Trail Configuration

The `--actor` flag is available on major operations to track who made changes:

- `knot project create --actor "john.doe"`
- `knot task create --actor "ai.agent"`
- `knot task update-state --actor "project.manager"`
- `knot task update-title --actor "project.manager"`
- `knot task update-description --actor "project.manager"`
- `knot dependency add --actor "project.manager"`
- `knot dependency remove --actor "project.manager"`

If no actor is specified, the system will use:
- The value from the `KNOT_ACTOR` environment variable if set
- The value from the `USER` environment variable if `KNOT_ACTOR` is not set
- "unknown" if neither environment variable is set

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
- Audit trail for tracking changes

Perfect for organizing complex AI workflows and project structures.
