# Agent Management Tools for Project Management

## Overview

The project management toolset now includes comprehensive agent management capabilities that allow LLMs to discover available agents in the system and assign tasks to them based on their capabilities and descriptions.

## New Features

### 1. Agent Discovery
- **Tool**: `list_available_agents`
- **Purpose**: Discover all available agents in the system
- **Returns**: Agent ID, name, role, and detailed capability descriptions

### 2. Task Assignment
- **Tool**: `assign_task_to_agent`
- **Purpose**: Assign specific tasks to agents based on their capabilities
- **Validation**: Ensures agent exists and is available

### 3. Task Management
- **Tools**: `unassign_task_from_agent`, `list_tasks_by_agent`, `list_unassigned_tasks`
- **Purpose**: Manage agent workloads and task distribution

## Available Tools

### `list_available_agents`
**Description**: List all available agents in the system with their capabilities and descriptions.

**Usage**: Use this first to discover which agents are available for task assignment. Each agent has a unique ID, name, role, and detailed description of their capabilities.

**Parameters**: None

**Returns**:
```json
{
  "agents": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440004",
      "name": "researcher",
      "role": "researcher",
      "description": "Research and information gathering agent"
    }
  ],
  "count": 1,
  "message": "Found 1 available agents"
}
```

### `assign_task_to_agent`
**Description**: Assign a specific task to an agent based on their capabilities.

**Usage**: First use `list_available_agents` to see available agents and their descriptions, then choose the most suitable agent for the task.

**Parameters**:
- `task_id` (string): The ID of the task to assign
- `agent_id` (string): The ID of the agent to assign the task to

**Returns**: Updated task with agent assignment

**Example Decision Process**:
1. Use `list_available_agents` to see available agents
2. Review task requirements and agent capabilities
3. Match task type with agent expertise:
   - Research tasks → researcher agent
   - Coding tasks → coder agent
   - Project planning → project-manager agent

### `unassign_task_from_agent`
**Description**: Remove agent assignment from a task.

**Usage**: Use when you want to make a task available for reassignment or when the current assignment is no longer appropriate.

**Parameters**:
- `task_id` (string): The ID of the task to unassign

### `list_tasks_by_agent`
**Description**: List all tasks assigned to a specific agent in a project.

**Usage**: Use to see the workload of a particular agent or review what tasks are currently assigned to them.

**Parameters**:
- `project_id` (string): The ID of the project to search in
- `agent_id` (string): The ID of the agent to find tasks for

### `list_unassigned_tasks`
**Description**: List all tasks in a project that have no agent assigned.

**Usage**: Use to identify tasks that need to be assigned to agents. These tasks should be distributed among available agents based on their capabilities and current workload.

**Parameters**:
- `project_id` (string): The ID of the project to search in

## LLM Decision Making Guidelines

### Agent Selection Criteria

When assigning tasks to agents, consider:

1. **Agent Role and Capabilities**:
   - **researcher**: Information gathering, web searches, data analysis
   - **coder**: Software development, code review, technical implementation
   - **project-manager**: Planning, coordination, task breakdown, progress tracking

2. **Task Complexity and Type**:
   - Research tasks → researcher agent
   - Implementation tasks → coder agent
   - Planning and coordination → project-manager agent

3. **Current Workload**:
   - Use `list_tasks_by_agent` to check current assignments
   - Distribute tasks evenly among available agents

4. **Task Dependencies**:
   - Consider task dependencies when assigning
   - Ensure prerequisite tasks are completed or assigned appropriately

### Workflow Examples

#### Example 1: New Project Setup
```
1. list_available_agents() // Discover available agents
2. create_project() // Create the project
3. create_task() // Create initial tasks
4. list_unassigned_tasks() // See what needs assignment
5. assign_task_to_agent() // Assign based on capabilities
```

#### Example 2: Workload Balancing
```
1. list_available_agents() // See all agents
2. For each agent: list_tasks_by_agent() // Check current workload
3. list_unassigned_tasks() // Find tasks needing assignment
4. assign_task_to_agent() // Distribute evenly based on capabilities
```

#### Example 3: Task Reassignment
```
1. list_tasks_by_agent() // Check current assignments
2. unassign_task_from_agent() // Remove inappropriate assignments
3. assign_task_to_agent() // Reassign to more suitable agent
```

## Integration with Existing Tools

The agent management tools integrate seamlessly with existing project management tools:

- **Task Creation**: New tasks can be immediately assigned during creation
- **Task Updates**: Assignment information is preserved during task updates
- **Project Progress**: Agent assignments are visible in task listings
- **Dependencies**: Agent assignments work with task dependencies

## Data Model Changes

### Task Structure
Tasks now include an `assigned_agent` field:
```json
{
  "id": "task-uuid",
  "title": "Task Title",
  "description": "Task Description",
  "assigned_agent": "agent-uuid", // NEW FIELD
  "state": "pending",
  "complexity": 5
}
```

### Agent Information
Agent information includes:
```json
{
  "id": "agent-uuid",
  "name": "agent-name",
  "role": "agent-role",
  "description": "Detailed capability description"
}
```

## Best Practices for LLMs

1. **Always discover agents first**: Use `list_available_agents` before making assignments
2. **Read agent descriptions**: Use the detailed descriptions to make informed decisions
3. **Check workloads**: Use `list_tasks_by_agent` to avoid overloading agents
4. **Monitor unassigned tasks**: Regularly use `list_unassigned_tasks` to ensure all work is assigned
5. **Consider task types**: Match task requirements with agent capabilities
6. **Validate assignments**: Ensure agents exist before attempting assignment

## Error Handling

The tools include comprehensive error handling:
- Invalid UUIDs are caught and reported clearly
- Non-existent agents are validated before assignment
- Missing projects/tasks are handled gracefully
- Clear error messages guide correct usage

## Performance Considerations

- Agent lists are cached within the toolset
- Task filtering is done efficiently in memory
- Database queries are optimized for common operations
- Logging provides visibility into operations without affecting performance