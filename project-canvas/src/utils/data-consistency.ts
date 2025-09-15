/**
 * Data Consistency Utilities
 * Ensures data matches Go model structure exactly
 */

import { Task, TaskState } from '../types/task.types';
import { Project, UUID } from '../types/project.types';
import { isValidUUID } from './uuid';

/**
 * Validate task data consistency with Go model
 */
export function validateTaskConsistency(task: Task): string[] {
  const errors: string[] = [];
  
  // Required fields from Go model
  if (!task.id || !isValidUUID(task.id)) {
    errors.push('Task ID must be a valid UUID');
  }
  
  if (!task.projectId || !isValidUUID(task.projectId)) {
    errors.push('Project ID must be a valid UUID');
  }
  
  if (!task.title || task.title.trim().length === 0) {
    errors.push('Task title is required');
  }
  
  if (!task.description) {
    errors.push('Task description is required');
  }
  
  if (!Object.values(TaskState).includes(task.state)) {
    errors.push('Invalid task state');
  }
  
  if (!Number.isInteger(task.complexity) || task.complexity < 1 || task.complexity > 10) {
    errors.push('Task complexity must be integer 1-10');
  }
  
  if (!Number.isInteger(task.depth) || task.depth < 0) {
    errors.push('Task depth must be non-negative integer');
  }
  
  // Optional fields validation
  if (task.parentId && !isValidUUID(task.parentId)) {
    errors.push('Parent ID must be valid UUID if provided');
  }
  
  if (task.assignedAgent && !isValidUUID(task.assignedAgent)) {
    errors.push('Assigned agent must be valid UUID if provided');
  }
  
  if (task.estimate !== undefined && (!Number.isInteger(task.estimate) || task.estimate < 0)) {
    errors.push('Estimate must be non-negative integer (minutes) if provided');
  }
  
  // Dependencies must be UUID arrays
  if (!Array.isArray(task.dependencies)) {
    errors.push('Dependencies must be an array');
  } else {
    task.dependencies.forEach((depId, index) => {
      if (!isValidUUID(depId)) {
        errors.push(`Dependency at index ${index} must be valid UUID`);
      }
    });
  }
  
  // Dependents must be UUID arrays
  if (!Array.isArray(task.dependents)) {
    errors.push('Dependents must be an array');
  } else {
    task.dependents.forEach((depId, index) => {
      if (!isValidUUID(depId)) {
        errors.push(`Dependent at index ${index} must be valid UUID`);
      }
    });
  }
  
  // Date validation
  if (!(task.createdAt instanceof Date) || isNaN(task.createdAt.getTime())) {
    errors.push('Created date must be valid Date object');
  }
  
  if (!(task.updatedAt instanceof Date) || isNaN(task.updatedAt.getTime())) {
    errors.push('Updated date must be valid Date object');
  }
  
  if (task.completedAt && (!(task.completedAt instanceof Date) || isNaN(task.completedAt.getTime()))) {
    errors.push('Completed date must be valid Date object if provided');
  }
  
  return errors;
}

/**
 * Validate project data consistency with Go model
 */
export function validateProjectConsistency(project: Project): string[] {
  const errors: string[] = [];
  
  // Required fields from Go model
  if (!project.id || !isValidUUID(project.id)) {
    errors.push('Project ID must be a valid UUID');
  }
  
  if (!project.title || project.title.trim().length === 0) {
    errors.push('Project title is required');
  }
  
  if (!project.description) {
    errors.push('Project description is required');
  }
  
  // Progress metrics validation
  if (!Number.isInteger(project.totalTasks) || project.totalTasks < 0) {
    errors.push('Total tasks must be non-negative integer');
  }
  
  if (!Number.isInteger(project.completedTasks) || project.completedTasks < 0) {
    errors.push('Completed tasks must be non-negative integer');
  }
  
  if (project.completedTasks > project.totalTasks) {
    errors.push('Completed tasks cannot exceed total tasks');
  }
  
  if (typeof project.progress !== 'number' || project.progress < 0 || project.progress > 100) {
    errors.push('Progress must be number between 0 and 100');
  }
  
  // Date validation
  if (!(project.createdAt instanceof Date) || isNaN(project.createdAt.getTime())) {
    errors.push('Created date must be valid Date object');
  }
  
  if (!(project.updatedAt instanceof Date) || isNaN(project.updatedAt.getTime())) {
    errors.push('Updated date must be valid Date object');
  }
  
  return errors;
}

/**
 * Ensure task dependencies are consistent
 */
export function validateTaskDependencies(tasks: Task[]): string[] {
  const errors: string[] = [];
  const taskIds = new Set(tasks.map(t => t.id));
  
  tasks.forEach(task => {
    // Check if all dependencies exist
    task.dependencies.forEach(depId => {
      if (!taskIds.has(depId)) {
        errors.push(`Task ${task.id} has dependency ${depId} that doesn't exist`);
      }
    });
    
    // Check if all dependents exist
    task.dependents.forEach(depId => {
      if (!taskIds.has(depId)) {
        errors.push(`Task ${task.id} has dependent ${depId} that doesn't exist`);
      }
    });
    
    // Check for self-dependencies
    if (task.dependencies.includes(task.id)) {
      errors.push(`Task ${task.id} cannot depend on itself`);
    }
    
    if (task.dependents.includes(task.id)) {
      errors.push(`Task ${task.id} cannot be dependent on itself`);
    }
  });
  
  return errors;
}

/**
 * Convert Go time.Time string to JavaScript Date
 */
export function parseGoTime(timeString: string): Date {
  return new Date(timeString);
}

/**
 * Convert JavaScript Date to Go time.Time format
 */
export function formatGoTime(date: Date): string {
  return date.toISOString();
}

/**
 * Ensure data structure matches Go model exactly
 */
export function normalizeTaskFromAPI(apiTask: any): Task {
  return {
    id: apiTask.id,
    projectId: apiTask.project_id,
    parentId: apiTask.parent_id || undefined,
    title: apiTask.title,
    description: apiTask.description,
    state: apiTask.state as TaskState,
    complexity: apiTask.complexity,
    depth: apiTask.depth,
    estimate: apiTask.estimate || undefined,
    assignedAgent: apiTask.assigned_agent || undefined,
    dependencies: apiTask.dependencies || [],
    dependents: apiTask.dependents || [],
    createdAt: parseGoTime(apiTask.created_at),
    updatedAt: parseGoTime(apiTask.updated_at),
    completedAt: apiTask.completed_at ? parseGoTime(apiTask.completed_at) : undefined,
    position: apiTask.position // UI-only field
  };
}

/**
 * Ensure project structure matches Go model exactly
 */
export function normalizeProjectFromAPI(apiProject: any): Project {
  return {
    id: apiProject.id,
    title: apiProject.title,
    description: apiProject.description,
    createdAt: parseGoTime(apiProject.created_at),
    updatedAt: parseGoTime(apiProject.updated_at),
    totalTasks: apiProject.total_tasks,
    completedTasks: apiProject.completed_tasks,
    progress: apiProject.progress
  };
}