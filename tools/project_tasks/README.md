# Project Task Management Toolset

A hierarchical project task planning toolset that provides comprehensive task management capabilities with tree-structured tasks, parent-child relationships, and configurable complexity thresholds.

## Features

### Core Capabilities
- **Hierarchical Task Structure**: Tasks organized in tree structure with parent-child relationships
- **Project Management**: Create and manage projects containing multiple task hierarchies
- **State Management**: Track task states (pending, in-progress, completed, blocked, cancelled)
- **Complexity-Based Breakdown**: Configurable complexity thresholds for task breakdown suggestions
- **Progress Tracking**: Real-time project progress metrics and completion tracking
- **Depth Control**: Configurable maximum tasks per depth level to control branching density

### Functional Tools
- **Project Operations**: Create, read, update, delete projects
- **Task Operations**: Full CRUD operations on tasks with automatic hierarchy management
- **Hierarchy Queries**: Recursive task listing with hierarchical display
- **Smart Task Finding**: Find next actionable tasks based on state and priority
- **Breakdown Analysis**: Identify tasks that need decomposition based on complexity
- **Progress Analytics**: Calculate detailed project progress metrics
- **Safe Deletion**: Remove tasks with subtree deletion support

## Architecture

### Components
- **Repository Layer**: Pluggable persistence with in-memory implementation
- **Service Layer**: Business logic and validation
- **Tool Layer**: Function-based tools following established patterns
- **Types**: Comprehensive type definitions with validation

### Configuration
```go
type Config struct {
    MaxTasksPerDepth      map[int]int // Maximum tasks allowed per depth level
    ComplexityThreshold   int         // Threshold for task breakdown suggestions
    MaxDepth              int         // Maximum allowed depth
    DefaultPriority       int         // Default priority for new tasks
}
```

## Usage Example

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/denkhaus/agents/tools/project_tasks"
)

func main() {
    ctx := context.Background()
    
    // Create toolset with custom configuration
    config := &projecttasks.Config{
        MaxTasksPerDepth: map[int]int{
            0: 10,  // Max 10 root tasks
            1: 20,  // Max 20 tasks at depth 1
            2: 50,  // Max 50 tasks at depth 2
        },
        ComplexityThreshold: 8,
        MaxDepth:           3,
        DefaultPriority:    5,
    }
    
    toolSet, err := projecttasks.NewToolSet(
        projecttasks.WithConfig(config),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer toolSet.Close()
    
    tools := toolSet.Tools(ctx)
    
    // Find the create_project tool
    var createProjectTool tool.CallableTool
    for _, t := range tools {
        if t.Declaration().Name == "create_project" {
            createProjectTool = t
            break
        }
    }
    
    // Create a project
    projectInput := map[string]interface{}{
        "title":       "Software Development Project",
        "description": "Building a new web application",
    }
    projectInputJSON, _ := json.Marshal(projectInput)
    
    projectResult, err := createProjectTool.Call(ctx, projectInputJSON)
    if err != nil {
        log.Fatal(err)
    }
    
    project := projectResult.(*projecttasks.Project)
    fmt.Printf("Created project: %s (ID: %s)\n", project.Title, project.ID)
    
    // Create a high-level task
    var createTaskTool tool.CallableTool
    for _, t := range tools {
        if t.Declaration().Name == "create_task" {
            createTaskTool = t
            break
        }
    }
    
    taskInput := map[string]interface{}{
        "project_id":  project.ID.String(),
        "title":       "Design Database Schema",
        "description": "Design and implement the database schema for the application",
        "complexity":  9, // High complexity - will need breakdown
        "priority":    8,
    }
    taskInputJSON, _ := json.Marshal(taskInput)
    
    taskResult, err := createTaskTool.Call(ctx, taskInputJSON)
    if err != nil {
        log.Fatal(err)
    }
    
    task := taskResult.(*projecttasks.Task)
    fmt.Printf("Created task: %s (ID: %s, Complexity: %d)\n", 
        task.Title, task.ID, task.Complexity)
    
    // Find tasks that need breakdown
    var breakdownTool tool.CallableTool
    for _, t := range tools {
        if t.Declaration().Name == "find_tasks_needing_breakdown" {
            breakdownTool = t
            break
        }
    }
    
    breakdownInput := map[string]interface{}{
        "project_id": project.ID.String(),
    }
    breakdownInputJSON, _ := json.Marshal(breakdownInput)
    
    breakdownResult, err := breakdownTool.Call(ctx, breakdownInputJSON)
    if err != nil {
        log.Fatal(err)
    }
    
    breakdownData := breakdownResult.(map[string]interface{})
    tasks := breakdownData["tasks"].([]*projecttasks.Task)
    
    fmt.Printf("Tasks needing breakdown: %d\n", len(tasks))
    for _, t := range tasks {
        fmt.Printf("- %s (Complexity: %d)\n", t.Title, t.Complexity)
    }
}
```

## Available Tools

### Project Management
- `create_project` - Create a new project
- `get_project` - Get project details by ID
- `update_project` - Update project title and description
- `delete_project` - Delete a project and all its tasks
- `list_projects` - List all projects

### Task Management
- `create_task` - Create a new task in a project
- `get_task` - Get task details by ID
- `update_task` - Update task properties
- `update_task_state` - Update only the task state
- `delete_task` - Delete a single task (fails if task has children)
- `delete_task_subtree` - Delete a task and all its subtasks recursively

### Task Queries
- `list_tasks_hierarchical` - List all tasks in hierarchical structure
- `get_task_subtree` - Get a task and all its subtasks
- `find_next_actionable_task` - Find the next task to work on
- `find_tasks_needing_breakdown` - Find high-complexity tasks needing breakdown
- `get_project_progress` - Get detailed progress metrics
- `list_tasks_by_state` - List tasks filtered by state

## Task States

- `pending` - Task is ready to be worked on
- `in-progress` - Task is currently being worked on
- `completed` - Task has been finished
- `blocked` - Task is blocked by dependencies
- `cancelled` - Task has been cancelled

## Safety Features

- **Concurrency Safe**: All operations are thread-safe with proper locking
- **Validation**: Comprehensive input validation with descriptive error messages
- **Depth Limits**: Configurable maximum depth to prevent infinite hierarchies
- **Branching Control**: Configurable maximum tasks per depth level
- **Safe Deletion**: Prevents accidental deletion of tasks with children
- **Project Context**: All operations require project ID for safe context isolation

## Testing

The toolset includes comprehensive unit tests covering:
- Basic CRUD operations
- Hierarchical operations and edge cases
- Concurrent operations
- Validation and error handling
- Configuration limits and constraints

Run tests with:
```bash
go test ./tools/project_tasks/...
```