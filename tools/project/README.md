# Project Task Management Tool

This package provides a comprehensive project and task management system designed to enable an LLM model to administer projects with hierarchical tasks.

## Features

### Project Management
- Create, retrieve, update, and delete projects
- List all projects
- Update project descriptions independently

### Task Management
- Create hierarchical tasks with parent-child relationships
- Retrieve task details
- Update task descriptions independently
- Update task state (pending, in-progress, completed, blocked, cancelled)
- Delete tasks and subtrees

### Task Organization
- Get parent task of a given task
- Get child tasks of a given task
- Get root tasks of a project

### Analysis and Progress Tracking
- Find next actionable task in a project
- Find tasks needing breakdown based on complexity
- Get detailed project progress metrics
- List tasks by state

### Configuration
- Configurable limits for tasks per depth level
- Configurable complexity threshold for task breakdown suggestions
- Configurable maximum depth for task hierarchy
- Configurable default priority for new tasks
- Configurable maximum description length

## Usage

The package provides a `ToolSetProvider` that can be used to create a toolset for LLM agents. The toolset includes the following functions:

- `create_project` - Create a new project
- `get_project` - Get project details by ID
- `update_project_description` - Update only the project description
- `list_projects` - List all projects
- `create_task` - Create a new task in a project
- `get_task` - Get task details by ID
- `update_task_description` - Update only the task description
- `update_task_state` - Update only the task state
- `get_project_progress` - Get detailed progress metrics for a project
- `get_child_tasks` - Get the child tasks of a given task
- `get_parent_task` - Get the parent task of a given task

## Configuration

The system can be configured with the following options:

- `MaxTasksPerDepth` - Maximum tasks allowed per depth level (default: 20)
- `ComplexityThreshold` - Threshold for task breakdown suggestions (default: 8)
- `MaxDepth` - Maximum allowed depth for task hierarchy (default: 5)
- `DefaultPriority` - Default priority for new tasks (default: 5)
- `MaxDescriptionLength` - Maximum length for descriptions (default: 2000)

## Data Models

### Project
- `ID` - Unique identifier
- `Title` - Project title
- `Description` - Project description
- `CreatedAt` - Creation timestamp
- `UpdatedAt` - Last update timestamp
- `TotalTasks` - Total number of tasks in the project
- `CompletedTasks` - Number of completed tasks
- `Progress` - Overall progress percentage

### Task
- `ID` - Unique identifier
- `ProjectID` - Associated project ID
- `ParentID` - Parent task ID (nil for root tasks)
- `Title` - Task title
- `Description` - Task description
- `State` - Task state (pending, in-progress, completed, blocked, cancelled)
- `Complexity` - Task complexity (1-10)
- `Priority` - Task priority (1-10)
- `Depth` - Task depth in hierarchy (0 for root tasks)
- `CreatedAt` - Creation timestamp
- `UpdatedAt` - Last update timestamp
- `CompletedAt` - Completion timestamp (when applicable)

### ProjectProgress
- `ProjectID` - Associated project ID
- `TotalTasks` - Total number of tasks
- `CompletedTasks` - Number of completed tasks
- `InProgressTasks` - Number of in-progress tasks
- `PendingTasks` - Number of pending tasks
- `BlockedTasks` - Number of blocked tasks
- `CancelledTasks` - Number of cancelled tasks
- `OverallProgress` - Overall progress percentage
- `TasksByDepth` - Task count by depth level