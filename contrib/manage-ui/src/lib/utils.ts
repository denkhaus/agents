import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"
import type { Task, TaskState, TaskHierarchy } from './types';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// Date formatting utilities
export function formatDate(date: string | Date): string {
  const d = new Date(date);
  return d.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

export function formatDateTime(date: string | Date): string {
  const d = new Date(date);
  return d.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function formatRelativeTime(date: string | Date): string {
  const d = new Date(date);
  const now = new Date();
  const diffMs = now.getTime() - d.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMins / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffMins < 1) return 'just now';
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;
  
  return formatDate(d);
}

// Task state utilities
export function getTaskStateColor(state: TaskState): string {
  switch (state) {
    case 'pending':
      return 'bg-gray-100 text-gray-800 border-gray-200';
    case 'in-progress':
      return 'bg-blue-100 text-blue-800 border-blue-200';
    case 'completed':
      return 'bg-green-100 text-green-800 border-green-200';
    case 'blocked':
      return 'bg-red-100 text-red-800 border-red-200';
    case 'cancelled':
      return 'bg-gray-100 text-gray-600 border-gray-200';
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200';
  }
}

export function getTaskStateIcon(state: TaskState): string {
  switch (state) {
    case 'pending':
      return 'clock';
    case 'in-progress':
      return 'play-circle';
    case 'completed':
      return 'check-circle';
    case 'blocked':
      return 'x-circle';
    case 'cancelled':
      return 'minus-circle';
    default:
      return 'circle';
  }
}

export function getComplexityColor(complexity: number): string {
  if (complexity <= 3) return 'bg-green-100 text-green-800';
  if (complexity <= 6) return 'bg-yellow-100 text-yellow-800';
  if (complexity <= 8) return 'bg-orange-100 text-orange-800';
  return 'bg-red-100 text-red-800';
}

// Time estimation utilities
export function formatEstimate(minutes: number | null | undefined): string {
  if (!minutes) return 'No estimate';
  
  if (minutes < 60) return `${minutes}m`;
  
  const hours = Math.floor(minutes / 60);
  const remainingMins = minutes % 60;
  
  if (remainingMins === 0) return `${hours}h`;
  return `${hours}h ${remainingMins}m`;
}

// Task hierarchy utilities
export function buildTaskHierarchy(tasks: Task[]): TaskHierarchy[] {
  const taskMap = new Map<string, Task>();
  const rootTasks: Task[] = [];
  
  // Create task map and identify root tasks
  tasks.forEach(task => {
    taskMap.set(task.id, task);
    if (!task.parent_id) {
      rootTasks.push(task);
    }
  });
  
  // Build hierarchy recursively
  function buildChildren(parentId: string, level: number = 0): TaskHierarchy[] {
    const children = tasks.filter(task => task.parent_id === parentId);
    
    return children.map(task => ({
      task,
      children: buildChildren(task.id, level + 1),
      level
    }));
  }
  
  return rootTasks.map(task => ({
    task,
    children: buildChildren(task.id, 1),
    level: 0
  }));
}

export function flattenTaskHierarchy(hierarchy: TaskHierarchy[]): Task[] {
  const result: Task[] = [];
  
  function flatten(items: TaskHierarchy[]) {
    items.forEach(item => {
      result.push(item.task);
      if (item.children.length > 0) {
        flatten(item.children);
      }
    });
  }
  
  flatten(hierarchy);
  return result;
}

export function getTaskDepth(task: Task, allTasks: Task[]): number {
  if (!task.parent_id) return 0;
  
  const parent = allTasks.find(t => t.id === task.parent_id);
  if (!parent) return 0;
  
  return 1 + getTaskDepth(parent, allTasks);
}

// Search and filtering utilities
export function searchTasks(tasks: Task[], query: string): Task[] {
  if (!query.trim()) return tasks;
  
  const lowerQuery = query.toLowerCase();
  
  return tasks.filter(task => 
    task.title.toLowerCase().includes(lowerQuery) ||
    task.description.toLowerCase().includes(lowerQuery) ||
    task.id.toLowerCase().includes(lowerQuery)
  );
}

export function filterTasksByState(tasks: Task[], states: TaskState[]): Task[] {
  if (states.length === 0) return tasks;
  return tasks.filter(task => states.includes(task.state));
}

export function filterTasksByComplexity(
  tasks: Task[], 
  minComplexity: number, 
  maxComplexity: number
): Task[] {
  return tasks.filter(task => 
    task.complexity >= minComplexity && task.complexity <= maxComplexity
  );
}

// Validation utilities
export function validateTaskTitle(title: string): string | null {
  if (!title.trim()) return 'Title is required';
  if (title.length > 200) return 'Title must be less than 200 characters';
  return null;
}

export function validateTaskDescription(description: string): string | null {
  if (description.length > 2000) return 'Description must be less than 2000 characters';
  return null;
}

export function validateTaskComplexity(complexity: number): string | null {
  if (complexity < 1 || complexity > 10) return 'Complexity must be between 1 and 10';
  return null;
}

// Progress calculation utilities
export function calculateProjectProgress(tasks: Task[]): {
  total: number;
  completed: number;
  inProgress: number;
  pending: number;
  blocked: number;
  cancelled: number;
  progress: number;
} {
  const total = tasks.length;
  const completed = tasks.filter(t => t.state === 'completed').length;
  const inProgress = tasks.filter(t => t.state === 'in-progress').length;
  const pending = tasks.filter(t => t.state === 'pending').length;
  const blocked = tasks.filter(t => t.state === 'blocked').length;
  const cancelled = tasks.filter(t => t.state === 'cancelled').length;
  
  const progress = total > 0 ? (completed / total) * 100 : 0;
  
  return {
    total,
    completed,
    inProgress,
    pending,
    blocked,
    cancelled,
    progress: Math.round(progress * 100) / 100
  };
}

// ID generation utility
export function generateId(): string {
  return Math.random().toString(36).substring(2) + Date.now().toString(36);
}
