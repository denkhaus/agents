// API response types and request interfaces

export interface ApiResponse<T> {
  data: T;
  success: boolean;
  message?: string;
  error?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface CreateProjectRequest {
  title: string;
  description: string;
}

export interface UpdateProjectRequest {
  title?: string;
  description?: string;
}

export interface CreateTaskRequest {
  project_id: string;
  parent_id?: string;
  title: string;
  description: string;
  complexity: number;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  complexity?: number;
  state?: import('./project').TaskState;
  estimate?: number;
  assigned_agent?: string;
}

export interface BulkUpdateTasksRequest {
  task_ids: string[];
  updates: import('./project').TaskUpdates;
}

export interface AddDependencyRequest {
  task_id: string;
  depends_on_task_id: string;
}

// SSE Event types
export interface SSEEvent<T = unknown> {
  type: string;
  data: T;
  timestamp: string;
}

export interface ProjectUpdateEvent {
  project_id: string;
  action: 'created' | 'updated' | 'deleted';
  project?: import('./project').Project;
}

export interface TaskUpdateEvent {
  task_id: string;
  project_id: string;
  action: 'created' | 'updated' | 'deleted' | 'state_changed';
  task?: import('./project').Task;
  old_state?: import('./project').TaskState;
  new_state?: import('./project').TaskState;
}

// Chat integration types
export interface ChatSession {
  id: string;
  entity_type: 'project' | 'task';
  entity_id: string;
  created_at: string;
  updated_at: string;
}

export interface ChatMessage {
  id: string;
  session_id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: string;
}

export interface ChatRequest {
  entity_type: 'project' | 'task';
  entity_id: string;
  message: string;
}