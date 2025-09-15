/**
 * Project-related TypeScript types
 * Based on Go struct from pkg/tools/project/shared/shared.go
 */

export type UUID = string; // UUID v4 format

export interface Project {
  id: UUID;
  title: string;
  description: string;
  createdAt: Date;
  updatedAt: Date;
  // Progress metrics
  totalTasks: number;
  completedTasks: number;
  progress: number; // Percentage (0-100)
}

export interface ProjectProgress {
  projectId: UUID;
  totalTasks: number;
  completedTasks: number;
  inProgressTasks: number;
  pendingTasks: number;
  blockedTasks: number;
  cancelledTasks: number;
  overallProgress: number;
  tasksByDepth: Record<number, number>;
}

export interface ProjectFilter {
  searchTerm?: string;
  sortBy?: 'title' | 'createdAt' | 'updatedAt' | 'progress';
  sortOrder?: 'asc' | 'desc';
}

export interface ProjectCreateInput {
  title: string;
  description: string;
}

export interface ProjectUpdateInput {
  title?: string;
  description?: string;
}