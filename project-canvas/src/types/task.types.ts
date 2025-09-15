/**
 * Task-related TypeScript types
 * Based on Go struct from pkg/tools/project/shared/shared.go
 */

import { UUID } from './project.types';

export enum TaskState {
  PENDING = "pending",
  IN_PROGRESS = "in-progress",
  COMPLETED = "completed",
  BLOCKED = "blocked",
  CANCELLED = "cancelled"
}

export interface Task {
  id: UUID;
  projectId: UUID;
  parentId?: UUID; // undefined for root tasks
  title: string;
  description: string;
  state: TaskState;
  complexity: number; // 1-10, used for breakdown decisions
  depth: number; // 0 for root tasks
  estimate?: number; // Time estimate in minutes
  assignedAgent?: UUID; // Agent assigned to this task
  dependencies: UUID[]; // Tasks this task depends on
  dependents: UUID[]; // Tasks that depend on this task
  createdAt: Date;
  updatedAt: Date;
  completedAt?: Date;
  // UI-specific properties
  position?: { x: number; y: number };
}

export interface TaskFilter {
  projectId?: UUID;
  parentId?: UUID;
  state?: TaskState;
  minDepth?: number;
  maxDepth?: number;
  minComplexity?: number;
  maxComplexity?: number;
  assignedAgent?: UUID;
  searchTerm?: string;
}

export interface TaskUpdates {
  state?: TaskState;
  complexity?: number;
  title?: string;
  description?: string;
  estimate?: number;
  assignedAgent?: UUID;
  position?: { x: number; y: number };
}

export interface TaskCreateInput {
  projectId: UUID;
  parentId?: UUID;
  title: string;
  description: string;
  complexity: number;
  estimate?: number;
  assignedAgent?: UUID;
}