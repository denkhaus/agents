/**
 * Data consistency utilities
 * Ensuring data integrity between frontend and backend
 */

import { Task } from "../types/task.types";
import { Project } from "../types/project.types";
import { Agent } from "../types/agent.types";

// Ensure task data consistency
export function ensureTaskConsistency(task: Task): Task {
  // Ensure required fields have default values
  const consistentTask: Task = {
    id: task.id,
    projectId: task.projectId,
    parentId: task.parentId,
    title: task.title || "Untitled Task",
    description: task.description || "",
    state: task.state || "pending",
    complexity: task.complexity || 1,
    depth: task.depth || 0,
    estimate: task.estimate,
    assignedAgent: task.assignedAgent,
    dependencies: task.dependencies || [],
    dependents: task.dependents || [],
    createdAt: task.createdAt || new Date(),
    updatedAt: task.updatedAt || new Date(),
    completedAt: task.completedAt,
    position: task.position,
  };

  return consistentTask;
}

// Ensure project data consistency
export function ensureProjectConsistency(project: Project): Project {
  const consistentProject: Project = {
    id: project.id,
    title: project.title || "Untitled Project",
    description: project.description || "",
    createdAt: project.createdAt || new Date(),
    updatedAt: project.updatedAt || new Date(),
    totalTasks: project.totalTasks || 0,
    completedTasks: project.completedTasks || 0,
    progress: project.progress || 0,
  };

  return consistentProject;
}

// Ensure agent data consistency
export function ensureAgentConsistency(agent: Agent): Agent {
  const consistentAgent: Agent = {
    id: agent.id,
    name: agent.name || "Unnamed Agent",
    role: agent.role || "human",
    description: agent.description || "",
    status: agent.status || "offline",
    isStreaming: agent.isStreaming || false,
    capabilities: agent.capabilities || [],
    currentTasks: agent.currentTasks || [],
    createdAt: agent.createdAt || new Date(),
    updatedAt: agent.updatedAt || new Date(),
    lastActiveAt: agent.lastActiveAt,
  };

  return consistentAgent;
}

// Normalize task data from external sources
export function normalizeTaskData(rawTask: Record<string, any>): Task {
  const normalizedTask: Task = {
    id: rawTask.id || rawTask._id || "",
    projectId: rawTask.projectId || "",
    parentId: rawTask.parentId,
    title: rawTask.title || rawTask.name || "Untitled Task",
    description: rawTask.description || rawTask.desc || "",
    state: rawTask.state || rawTask.status || "pending",
    complexity: rawTask.complexity || rawTask.difficulty || 1,
    depth: rawTask.depth || rawTask.level || 0,
    estimate: rawTask.estimate || rawTask.duration,
    assignedAgent: rawTask.assignedAgent || rawTask.assignee,
    dependencies: rawTask.dependencies || rawTask.deps || [],
    dependents: rawTask.dependents || rawTask.dependentsOf || [],
    createdAt: rawTask.createdAt ? new Date(rawTask.createdAt) : new Date(),
    updatedAt: rawTask.updatedAt ? new Date(rawTask.updatedAt) : new Date(),
    completedAt: rawTask.completedAt
      ? new Date(rawTask.completedAt)
      : undefined,
    position: rawTask.position
      ? {
          x: Number(rawTask.position.x) || 0,
          y: Number(rawTask.position.y) || 0,
        }
      : undefined,
  };

  return normalizedTask;
}

// Normalize project data from external sources
export function normalizeProjectData(rawProject: Record<string, any>): Project {
  const normalizedProject: Project = {
    id: rawProject.id || rawProject._id || "",
    title: rawProject.title || rawProject.name || "Untitled Project",
    description: rawProject.description || rawProject.desc || "",
    createdAt: rawProject.createdAt
      ? new Date(rawProject.createdAt)
      : new Date(),
    updatedAt: rawProject.updatedAt
      ? new Date(rawProject.updatedAt)
      : new Date(),
    totalTasks: rawProject.totalTasks || rawProject.taskCount || 0,
    completedTasks: rawProject.completedTasks || rawProject.completedCount || 0,
    progress: rawProject.progress || rawProject.completion || 0,
  };

  return normalizedProject;
}
