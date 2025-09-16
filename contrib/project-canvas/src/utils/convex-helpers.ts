/**
 * Convex Integration Helpers
 * Utilities for working with Convex real-time data
 */

import { Task, TaskState } from "../types/task.types";
import { Project, UUID } from "../types/project.types";
import { Agent, AgentRole, AgentStatus } from "../types/agent.types";

// Convex document types (will match Convex schema)
export interface ConvexProject {
  _id: string;
  _creationTime: number;
  id: string;
  title: string;
  description: string;
  totalTasks: number;
  completedTasks: number;
  progress: number;
  createdAt: number;
  updatedAt: number;
}

export interface ConvexTask {
  _id: string;
  _creationTime: number;
  id: string;
  projectId: string;
  parentId?: string;
  title: string;
  description: string;
  state: TaskState;
  complexity: number;
  depth: number;
  estimate?: number;
  assignedAgent?: string;
  dependencies: string[];
  dependents: string[];
  positionX?: number;
  positionY?: number;
  createdAt: number;
  completedAt?: number;
  updatedAt: number;
}

export interface ConvexAgent {
  _id: string;
  _creationTime: number;
  name: string;
  role?: string;
  description?: string;
  status?: string;
  isStreaming?: boolean;
  capabilities?: string[];
  currentTasks?: string[];
  lastActiveAt?: number;
  id: string;
}

/**
 * Convert Convex project to frontend Project type
 */
export function convexProjectToProject(convexProject: ConvexProject): Project {
  return {
    id: convexProject._id as UUID,
    title: convexProject.title,
    description: convexProject.description,
    createdAt: new Date(convexProject.createdAt),
    updatedAt: new Date(convexProject.updatedAt),
    totalTasks: convexProject.totalTasks,
    completedTasks: convexProject.completedTasks,
    progress: convexProject.progress,
  };
}

/**
 * Convert Convex task to frontend Task type
 */
export function convexTaskToTask(convexTask: ConvexTask): Task {
  return {
    id: convexTask.id as UUID,
    projectId: convexTask.projectId as UUID,
    parentId: convexTask.parentId as UUID | undefined,
    title: convexTask.title,
    description: convexTask.description,
    state: convexTask.state,
    complexity: convexTask.complexity,
    depth: convexTask.depth,
    estimate: convexTask.estimate,
    assignedAgent: convexTask.assignedAgent as UUID | undefined,
    dependencies: convexTask.dependencies as UUID[],
    dependents: convexTask.dependents as UUID[],
    createdAt: new Date(convexTask.createdAt),
    updatedAt: new Date(convexTask.updatedAt),
    completedAt: convexTask.completedAt
      ? new Date(convexTask.completedAt)
      : undefined,
    position:
      convexTask.positionX !== undefined && convexTask.positionY !== undefined
        ? { x: convexTask.positionX, y: convexTask.positionY }
        : undefined,
  };
}

/**
 * Convert Convex agent to frontend Agent type
 */
export function convexAgentToAgent(convexAgent: ConvexAgent): Agent {
  // Provide defaults for required fields that might be missing
  const role = convexAgent.role || "coder"; // Default role
  const description = convexAgent.description || ""; // Default empty description
  const status = convexAgent.status || "offline"; // Default status
  const isStreaming = convexAgent.isStreaming || false; // Default streaming status
  const capabilities = convexAgent.capabilities || []; // Default empty capabilities
  const currentTasks = convexAgent.currentTasks || []; // Default empty tasks

  return {
    id: convexAgent.id as UUID,
    name: convexAgent.name,
    role: role as AgentRole,
    description: description,
    status: status as AgentStatus,
    isStreaming: isStreaming,
    capabilities: capabilities,
    currentTasks: currentTasks as UUID[],
    createdAt: new Date(convexAgent._creationTime),
    updatedAt: new Date(convexAgent._creationTime),
    lastActiveAt: convexAgent.lastActiveAt
      ? new Date(convexAgent.lastActiveAt)
      : undefined,
  };
}

/**
 * Convert frontend Project to Convex format (for updates)
 */
export function projectToConvexProject(
  project: Project
): Partial<ConvexProject> {
  return {
    id: project.id,
    title: project.title,
    description: project.description,
    totalTasks: project.totalTasks,
    completedTasks: project.completedTasks,
    progress: project.progress,
    createdAt: project.createdAt.getTime(),
    updatedAt: project.updatedAt.getTime(),
  };
}

/**
 * Convert frontend Task to Convex format (for updates)
 */
export function taskToConvexTask(task: Task): Partial<ConvexTask> {
  return {
    id: task.id,
    projectId: task.projectId,
    parentId: task.parentId,
    title: task.title,
    description: task.description,
    state: task.state,
    complexity: task.complexity,
    depth: task.depth,
    estimate: task.estimate,
    assignedAgent: task.assignedAgent,
    dependencies: task.dependencies,
    dependents: task.dependents,
    positionX: task.position?.x,
    positionY: task.position?.y,
    createdAt: task.createdAt.getTime(),
    completedAt: task.completedAt?.getTime(),
    updatedAt: task.updatedAt.getTime(),
  };
}

/**
 * Real-time subscription helpers
 */
export interface SubscriptionOptions {
  projectId?: UUID;
  includeCompleted?: boolean;
  maxDepth?: number;
}

/**
 * Generate Convex query args for tasks
 */
export function getTaskQueryArgs(options: SubscriptionOptions = {}) {
  return {
    projectId: options.projectId,
    includeCompleted: options.includeCompleted ?? true,
    maxDepth: options.maxDepth ?? 10,
  };
}

/**
 * Optimistic update helpers
 */
export interface OptimisticUpdate<T> {
  id: string;
  data: Partial<T>;
  timestamp: number;
}

export class OptimisticUpdateManager<T> {
  private updates = new Map<string, OptimisticUpdate<T>>();

  addUpdate(id: string, data: Partial<T>) {
    this.updates.set(id, {
      id,
      data,
      timestamp: Date.now(),
    });
  }

  removeUpdate(id: string) {
    this.updates.delete(id);
  }

  applyUpdates(items: T[]): T[] {
    return items.map((item) => {
      const update = this.updates.get((item as any).id);
      return update ? { ...item, ...update.data } : item;
    });
  }

  clearOldUpdates(maxAge: number = 5000) {
    const now = Date.now();
    for (const [id, update] of this.updates.entries()) {
      if (now - update.timestamp > maxAge) {
        this.updates.delete(id);
      }
    }
  }
}
