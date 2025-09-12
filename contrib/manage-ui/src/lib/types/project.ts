// TypeScript types based on Go structs from pkg/tools/project/shared/shared.go

export type TaskState = 
  | 'pending'
  | 'in-progress' 
  | 'completed'
  | 'blocked'
  | 'cancelled';

export interface Task {
  id: string;
  project_id: string;
  parent_id?: string | null;
  title: string;
  description: string;
  state: TaskState;
  complexity: number;
  depth: number;
  estimate?: number | null; // Time estimate in minutes
  assigned_agent?: string | null;
  dependencies?: string[];
  dependents?: string[];
  created_at: string;
  updated_at: string;
  completed_at?: string | null;
}

export interface Project {
  id: string;
  title: string;
  description: string;
  created_at: string;
  updated_at: string;
  total_tasks: number;
  completed_tasks: number;
  progress: number; // Percentage (0-100)
}

export interface ProjectProgress {
  project_id: string;
  total_tasks: number;
  completed_tasks: number;
  in_progress_tasks: number;
  pending_tasks: number;
  blocked_tasks: number;
  cancelled_tasks: number;
  overall_progress: number;
  tasks_by_depth: Record<number, number>;
}

export interface TaskFilter {
  project_id?: string;
  parent_id?: string;
  state?: TaskState;
  min_depth?: number;
  max_depth?: number;
  min_complexity?: number;
  max_complexity?: number;
}

export interface TaskUpdates {
  state?: TaskState;
  complexity?: number;
  title?: string;
  description?: string;
  estimate?: number;
  assigned_agent?: string;
}

// UI-specific types
export interface TaskHierarchy {
  task: Task;
  children: TaskHierarchy[];
  level: number;
}

export interface KanbanColumn {
  id: string;
  title: string;
  level: number;
  tasks: Task[];
  parentTaskId?: string;
}

export interface ProjectWithTasks extends Project {
  tasks: Task[];
  rootTasks: Task[];
  taskHierarchy: TaskHierarchy[];
}