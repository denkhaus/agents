/**
 * Validation utilities
 * Type-safe validation for forms and data
 */

import { UUID } from '../types/project.types';
import { Task, TaskState } from '../types/task.types';
import { Project } from '../types/project.types';
import { isValidUUID } from './uuid';

// Validation result type
export interface ValidationResult {
  isValid: boolean;
  errors: string[];
}

// Project validation
export function validateProject(project: Partial<Project>): ValidationResult {
  const errors: string[] = [];
  
  if (!project.title || project.title.trim().length === 0) {
    errors.push('Project title is required');
  }
  
  if (project.title && project.title.length > 200) {
    errors.push('Project title must be 200 characters or less');
  }
  
  if (project.description && project.description.length > 2000) {
    errors.push('Project description must be 2000 characters or less');
  }
  
  if (project.progress !== undefined) {
    if (project.progress < 0 || project.progress > 100) {
      errors.push('Project progress must be between 0 and 100');
    }
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
}

// Task validation
export function validateTask(task: Partial<Task>): ValidationResult {
  const errors: string[] = [];
  
  if (!task.title || task.title.trim().length === 0) {
    errors.push('Task title is required');
  }
  
  if (task.title && task.title.length > 200) {
    errors.push('Task title must be 200 characters or less');
  }
  
  if (task.description && task.description.length > 2000) {
    errors.push('Task description must be 2000 characters or less');
  }
  
  if (task.complexity !== undefined) {
    if (!Number.isInteger(task.complexity) || task.complexity < 1 || task.complexity > 10) {
      errors.push('Task complexity must be an integer between 1 and 10');
    }
  }
  
  if (task.depth !== undefined) {
    if (!Number.isInteger(task.depth) || task.depth < 0) {
      errors.push('Task depth must be a non-negative integer');
    }
  }
  
  if (task.estimate !== undefined) {
    if (!Number.isInteger(task.estimate) || task.estimate < 0) {
      errors.push('Task estimate must be a non-negative integer (minutes)');
    }
  }
  
  if (task.state && !Object.values(TaskState).includes(task.state)) {
    errors.push('Invalid task state');
  }
  
  if (task.projectId && !isValidUUID(task.projectId)) {
    errors.push('Invalid project ID format');
  }
  
  if (task.parentId && !isValidUUID(task.parentId)) {
    errors.push('Invalid parent task ID format');
  }
  
  if (task.assignedAgent && !isValidUUID(task.assignedAgent)) {
    errors.push('Invalid assigned agent ID format');
  }
  
  // Validate dependencies array
  if (task.dependencies) {
    task.dependencies.forEach((depId, index) => {
      if (!isValidUUID(depId)) {
        errors.push(`Invalid dependency ID at index ${index}`);
      }
    });
  }
  
  // Validate dependents array
  if (task.dependents) {
    task.dependents.forEach((depId, index) => {
      if (!isValidUUID(depId)) {
        errors.push(`Invalid dependent ID at index ${index}`);
      }
    });
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
}

// Circular dependency validation
export function validateNoCycles(
  taskId: UUID,
  dependencies: UUID[],
  allTasks: Task[]
): ValidationResult {
  const errors: string[] = [];
  const visited = new Set<UUID>();
  const visiting = new Set<UUID>();
  
  // Create a temporary task with the new dependencies to check for cycles
  const tempTask = {
    id: taskId,
    dependencies: dependencies
  };
  
  // Create a map of tasks for quick lookup
  const taskMap = new Map<UUID, Task | { id: UUID; dependencies: UUID[] }>();
  for (const task of allTasks) {
    taskMap.set(task.id, task);
  }
  taskMap.set(taskId, tempTask);
  
  function hasCycle(currentId: UUID): boolean {
    if (visiting.has(currentId)) {
      return true; // Cycle detected
    }
    
    if (visited.has(currentId)) {
      return false; // Already processed
    }
    
    visiting.add(currentId);
    
    const task = taskMap.get(currentId);
    if (task) {
      for (const depId of task.dependencies) {
        if (hasCycle(depId)) {
          return true;
        }
      }
    }
    
    visiting.delete(currentId);
    visited.add(currentId);
    return false;
  }
  
  // Check if adding these dependencies would create a cycle
  if (hasCycle(taskId)) {
    errors.push('Adding these dependencies would create a circular dependency');
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
}

// Form field validation helpers
export function validateRequired(value: any, fieldName: string): string | null {
  if (value === null || value === undefined || value === '') {
    return `${fieldName} is required`;
  }
  return null;
}

export function validateMaxLength(value: string, maxLength: number, fieldName: string): string | null {
  if (value && value.length > maxLength) {
    return `${fieldName} must be ${maxLength} characters or less`;
  }
  return null;
}

export function validateMinLength(value: string, minLength: number, fieldName: string): string | null {
  if (value && value.length < minLength) {
    return `${fieldName} must be at least ${minLength} characters`;
  }
  return null;
}

export function validateRange(value: number, min: number, max: number, fieldName: string): string | null {
  if (value < min || value > max) {
    return `${fieldName} must be between ${min} and ${max}`;
  }
  return null;
}