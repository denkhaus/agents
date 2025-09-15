/**
 * Convex Integration Helpers
 * Utilities for working with Convex real-time data
 */

import { Task, TaskState } from '../types/task.types';
import { Project, UUID } from '../types/project.types';
import { Agent } from '../types/agent.types';

// Convex document types (will match Convex schema)
export interface ConvexProject {
  _id: string;
  _creationTime: number;
  title: string;
  description: string;
  totalTasks: number;
  completedTasks: number;
  progress: number;
}

export interface ConvexTask {
  _id: string;
  _creationTime: number;
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
  completedAt?: number;
}

export interface ConvexAgent {
  _id: string;
  _creationTime: number;
  name: string;
  role: string;
  description: string;
  status: string;
  isStreaming: boolean;
  capabilities: string[];
  currentTasks: string[];
  lastActiveAt?: number;
}

/**
 * Convert Convex project to frontend Project type
 */
export function convexProjectToProject(convexProject: ConvexProject): Project {
  return {
    id: convexProject._id as UUID,
    title: convexProject.title,
    description: convexProject.description,
    createdAt: new Date(convexProject._creationTime),
    updatedAt: new Date(convexProject._creationTime), // Convex doesn't track updates separately
    totalTasks: convexProject.totalTasks,
    completedTasks: convexProject.completedTasks,
    progress: convexProject.progress
  };
}

/**
 * Convert Convex task to frontend Task type
 */
export function convexTaskToTask(convexTask: ConvexTask): Task {
  return {
    id: convexTask._id as UUID,
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
    createdAt: new Date(convexTask._creationTime),
    updatedAt: new Date(convexTask._creationTime),
    completedAt: convexTask.completedAt ? new Date(convexTask.completedAt) : undefined,
    position: convexTask.positionX !== undefined && convexTask.positionY !== undefined 
      ? { x: convexTask.positionX, y: convexTask.positionY }
      : undefined
  };
}

/**
 * Convert Convex agent to frontend Agent type
 */
export function convexAgentToAgent(convexAgent: ConvexAgent): Agent {
  return {
    id: convexAgent._id as UUID,
    name: convexAgent.name,
    role: convexAgent.role as any, // Will be properly typed with enum
    description: convexAgent.description,
    status: convexAgent.status as any, // Will be properly typed with enum
    isStreaming: convexAgent.isStreaming,
    capabilities: convexAgent.capabilities,
    currentTasks: convexAgent.currentTasks as UUID[],
    createdAt: new Date(convexAgent._creationTime),
    updatedAt: new Date(convexAgent._creationTime),
    lastActiveAt: convexAgent.lastActiveAt ? new Date(convexAgent.lastActiveAt) : undefined
  };
}

/**
 * Convert frontend Project to Convex format (for updates)
 */
export function projectToConvexProject(project: Project): Partial<ConvexProject> {
  return {
    title: project.title,
    description: project.description,
    totalTasks: project.totalTasks,
    completedTasks: project.completedTasks,
    progress: project.progress
  };
}

/**
 * Convert frontend Task to Convex format (for updates)
 */
export function taskToConvexTask(task: Task): Partial<ConvexTask> {
  return {
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
    completedAt: task.completedAt?.getTime()
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
    maxDepth: options.maxDepth ?? 10
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
      timestamp: Date.now()
    });
  }

  removeUpdate(id: string) {
    this.updates.delete(id);
  }

  applyUpdates(items: T[]): T[] {
    return items.map(item => {
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