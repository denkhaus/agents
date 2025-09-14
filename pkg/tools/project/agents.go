package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// Agent Management Tools

// assignTaskToAgent assigns a task to a specific agent
func (pts *projectTaskToolSet) assignTaskToAgent(ctx context.Context, args assignTaskToAgentArgs) (assignTaskToAgentResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return assignTaskToAgentResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	agentID, err := uuid.Parse(args.AgentID)
	if err != nil {
		return assignTaskToAgentResult{}, fmt.Errorf("invalid agent ID format: %w", err)
	}

	pts.logger.Info("Assigning task to agent", zap.String("taskID", taskID.String()), zap.String("agentID", agentID.String()))

	// Validate that the agent exists in available agents
	agentExists := false
	var agentName string
	for _, agent := range pts.availableAgents {
		if agent.ID == agentID {
			agentExists = true
			agentName = agent.Name
			break
		}
	}

	if !agentExists {
		return assignTaskToAgentResult{}, fmt.Errorf("agent with ID %s is not available in the system", agentID)
	}

	task, err := pts.manager.AssignTaskToAgent(ctx, taskID, agentID)
	if err != nil {
		pts.logger.Error("Failed to assign task to agent", zap.Error(err))
		return assignTaskToAgentResult{}, err
	}

	pts.logger.Info("Successfully assigned task to agent", zap.String("taskID", taskID.String()), zap.String("agentName", agentName))
	return assignTaskToAgentResult{
		Task:    task,
		Message: fmt.Sprintf("Successfully assigned task '%s' to agent '%s'", task.Title, agentName),
	}, nil
}

func (pts *projectTaskToolSet) assignTaskToAgentTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.assignTaskToAgent,
		function.WithName("assign_task_to_agent"),
		function.WithDescription("Assign a specific task to an agent. Use this when you want to delegate a task to a particular agent based on their capabilities. First use 'list_available_agents' to see which agents are available and their descriptions, then choose the most suitable agent for the task based on their role and capabilities."),
	)
}

// unassignTaskFromAgent removes agent assignment from a task
func (pts *projectTaskToolSet) unassignTaskFromAgent(ctx context.Context, args unassignTaskFromAgentArgs) (unassignTaskFromAgentResult, error) {
	taskID, err := uuid.Parse(args.TaskID)
	if err != nil {
		return unassignTaskFromAgentResult{}, fmt.Errorf("invalid task ID format: %w", err)
	}

	pts.logger.Info("Unassigning task from agent", zap.String("taskID", taskID.String()))

	task, err := pts.manager.UnassignTaskFromAgent(ctx, taskID)
	if err != nil {
		pts.logger.Error("Failed to unassign task from agent", zap.Error(err))
		return unassignTaskFromAgentResult{}, err
	}

	pts.logger.Info("Successfully unassigned task from agent", zap.String("taskID", taskID.String()))
	return unassignTaskFromAgentResult{
		Task:    task,
		Message: fmt.Sprintf("Successfully unassigned task '%s' from agent", task.Title),
	}, nil
}

func (pts *projectTaskToolSet) unassignTaskFromAgentTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.unassignTaskFromAgent,
		function.WithName("unassign_task_from_agent"),
		function.WithDescription("Remove agent assignment from a task. Use this when you want to make a task available for reassignment or when the current agent assignment is no longer appropriate."),
	)
}

// listTasksByAgent lists all tasks assigned to a specific agent
func (pts *projectTaskToolSet) listTasksByAgent(ctx context.Context, args listTasksByAgentArgs) (listTasksByAgentResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return listTasksByAgentResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	agentID, err := uuid.Parse(args.AgentID)
	if err != nil {
		return listTasksByAgentResult{}, fmt.Errorf("invalid agent ID format: %w", err)
	}

	pts.logger.Info("Listing tasks assigned to agent", zap.String("projectID", projectID.String()), zap.String("agentID", agentID.String()))

	// Find agent name for better logging
	var agentName string
	for _, agent := range pts.availableAgents {
		if agent.ID == agentID {
			agentName = agent.Name
			break
		}
	}

	tasks, err := pts.manager.ListTasksByAgent(ctx, projectID, agentID)
	if err != nil {
		pts.logger.Error("Failed to list tasks by agent", zap.Error(err))
		return listTasksByAgentResult{}, err
	}

	pts.logger.Info("Found tasks assigned to agent", zap.Int("count", len(tasks)), zap.String("agentName", agentName))
	return listTasksByAgentResult{
		Tasks:   tasks,
		Count:   len(tasks),
		Message: fmt.Sprintf("Found %d tasks assigned to agent '%s'", len(tasks), agentName),
	}, nil
}

func (pts *projectTaskToolSet) listTasksByAgentTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.listTasksByAgent,
		function.WithName("list_tasks_by_agent"),
		function.WithDescription("List all tasks assigned to a specific agent in a project. Use this to see the workload of a particular agent or to review what tasks are currently assigned to them."),
	)
}

// listUnassignedTasks lists all tasks that have no agent assigned
func (pts *projectTaskToolSet) listUnassignedTasks(ctx context.Context, args listUnassignedTasksArgs) (listUnassignedTasksResult, error) {
	projectID, err := uuid.Parse(args.ProjectID)
	if err != nil {
		return listUnassignedTasksResult{}, fmt.Errorf("invalid project ID format: %w", err)
	}

	pts.logger.Info("Listing unassigned tasks", zap.String("projectID", projectID.String()))

	tasks, err := pts.manager.ListUnassignedTasks(ctx, projectID)
	if err != nil {
		pts.logger.Error("Failed to list unassigned tasks", zap.Error(err))
		return listUnassignedTasksResult{}, err
	}

	pts.logger.Info("Found unassigned tasks", zap.Int("count", len(tasks)))
	return listUnassignedTasksResult{
		Tasks:   tasks,
		Count:   len(tasks),
		Message: fmt.Sprintf("Found %d unassigned tasks that need agent assignment", len(tasks)),
	}, nil
}

func (pts *projectTaskToolSet) listUnassignedTasksTool() tool.CallableTool {
	return function.NewFunctionTool(
		pts.listUnassignedTasks,
		function.WithName("list_unassigned_tasks"),
		function.WithDescription("List all tasks in a project that have no agent assigned. Use this to identify tasks that need to be assigned to agents. These tasks are available for assignment and should be distributed among available agents based on their capabilities and current workload."),
	)
}
