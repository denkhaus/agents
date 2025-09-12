import { api, ApiError } from './client';
import type {
  Project,
  Task,
  ProjectProgress,
  TaskFilter,
  ApiResponse,
  CreateProjectRequest,
  UpdateProjectRequest,
  CreateTaskRequest,
  UpdateTaskRequest,
  BulkUpdateTasksRequest,
  AddDependencyRequest,
} from '../types';

// Project API functions
export const projectsApi = {
  // Project CRUD operations
  async getProjects(): Promise<Project[]> {
    const response = await api.get<ApiResponse<Project[]>>('/projects');
    return response.data;
  },

  async getProject(projectId: string): Promise<Project> {
    const response = await api.get<ApiResponse<Project>>(`/projects/${projectId}`);
    return response.data;
  },

  async createProject(data: CreateProjectRequest): Promise<Project> {
    const response = await api.post<ApiResponse<Project>>('/projects', data);
    return response.data;
  },

  async updateProject(projectId: string, data: UpdateProjectRequest): Promise<Project> {
    const response = await api.put<ApiResponse<Project>>(`/projects/${projectId}`, data);
    return response.data;
  },

  async deleteProject(projectId: string): Promise<void> {
    await api.delete(`/projects/${projectId}`);
  },

  async getProjectProgress(projectId: string): Promise<ProjectProgress> {
    const response = await api.get<ApiResponse<ProjectProgress>>(`/projects/${projectId}/progress`);
    return response.data;
  },

  // Task operations
  async getProjectTasks(projectId: string): Promise<Task[]> {
    const response = await api.get<ApiResponse<Task[]>>(`/projects/${projectId}/tasks`);
    return response.data;
  },

  async getRootTasks(projectId: string): Promise<Task[]> {
    const response = await api.get<ApiResponse<Task[]>>(`/projects/${projectId}/tasks/root`);
    return response.data;
  },

  async getTask(taskId: string): Promise<Task> {
    const response = await api.get<ApiResponse<Task>>(`/tasks/${taskId}`);
    return response.data;
  },

  async createTask(data: CreateTaskRequest): Promise<Task> {
    const response = await api.post<ApiResponse<Task>>('/tasks', data);
    return response.data;
  },

  async updateTask(taskId: string, data: UpdateTaskRequest): Promise<Task> {
    const response = await api.put<ApiResponse<Task>>(`/tasks/${taskId}`, data);
    return response.data;
  },

  async deleteTask(taskId: string): Promise<void> {
    await api.delete(`/tasks/${taskId}`);
  },

  async deleteTaskSubtree(taskId: string): Promise<void> {
    await api.delete(`/tasks/${taskId}/subtree`);
  },

  // Task queries
  async getChildTasks(taskId: string): Promise<Task[]> {
    const response = await api.get<ApiResponse<Task[]>>(`/tasks/${taskId}/children`);
    return response.data;
  },

  async getParentTask(taskId: string): Promise<Task> {
    const response = await api.get<ApiResponse<Task>>(`/tasks/${taskId}/parent`);
    return response.data;
  },

  async listTasks(filter: TaskFilter): Promise<Task[]> {
    const params = new URLSearchParams();
    
    Object.entries(filter).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        params.append(key, value.toString());
      }
    });

    const queryString = params.toString();
    const endpoint = queryString ? `/tasks?${queryString}` : '/tasks';
    
    const response = await api.get<ApiResponse<Task[]>>(endpoint);
    return response.data;
  },

  async getTasksByState(projectId: string, state: string): Promise<Task[]> {
    const response = await api.get<ApiResponse<Task[]>>(`/projects/${projectId}/tasks/state/${state}`);
    return response.data;
  },

  async findNextActionableTask(projectId: string): Promise<Task | null> {
    try {
      const response = await api.get<ApiResponse<Task>>(`/projects/${projectId}/tasks/next-actionable`);
      return response.data;
    } catch (error: unknown) {
      if (error instanceof ApiError && error.status === 404) {
        return null;
      }
      throw error;
    }
  },

  async findTasksNeedingBreakdown(projectId: string): Promise<Task[]> {
    const response = await api.get<ApiResponse<Task[]>>(`/projects/${projectId}/tasks/needs-breakdown`);
    return response.data;
  },

  // Bulk operations
  async bulkUpdateTasks(data: BulkUpdateTasksRequest): Promise<void> {
    await api.patch('/tasks/bulk', data);
  },

  async duplicateTask(taskId: string, newProjectId: string): Promise<Task> {
    const response = await api.post<ApiResponse<Task>>(`/tasks/${taskId}/duplicate`, {
      new_project_id: newProjectId
    });
    return response.data;
  },

  async setTaskEstimate(taskId: string, estimate: number): Promise<Task> {
    const response = await api.patch<ApiResponse<Task>>(`/tasks/${taskId}/estimate`, {
      estimate
    });
    return response.data;
  },

  // Agent assignment
  async assignTaskToAgent(taskId: string, agentId: string): Promise<Task> {
    const response = await api.patch<ApiResponse<Task>>(`/tasks/${taskId}/assign`, {
      assigned_agent: agentId
    });
    return response.data;
  },

  async unassignTaskFromAgent(taskId: string): Promise<Task> {
    const response = await api.patch<ApiResponse<Task>>(`/tasks/${taskId}/unassign`);
    return response.data;
  },

  async getTasksByAgent(projectId: string, agentId: string): Promise<Task[]> {
    const response = await api.get<ApiResponse<Task[]>>(`/projects/${projectId}/tasks/agent/${agentId}`);
    return response.data;
  },

  async getUnassignedTasks(projectId: string): Promise<Task[]> {
    const response = await api.get<ApiResponse<Task[]>>(`/projects/${projectId}/tasks/unassigned`);
    return response.data;
  },

  // Dependencies
  async addTaskDependency(data: AddDependencyRequest): Promise<Task> {
    const response = await api.post<ApiResponse<Task>>('/tasks/dependencies', data);
    return response.data;
  },

  async removeTaskDependency(taskId: string, dependsOnTaskId: string): Promise<Task> {
    const response = await api.delete<ApiResponse<Task>>(`/tasks/${taskId}/dependencies/${dependsOnTaskId}`);
    return response.data;
  },

  async getTaskDependencies(taskId: string): Promise<Task[]> {
    const response = await api.get<ApiResponse<Task[]>>(`/tasks/${taskId}/dependencies`);
    return response.data;
  },

  async getDependentTasks(taskId: string): Promise<Task[]> {
    const response = await api.get<ApiResponse<Task[]>>(`/tasks/${taskId}/dependents`);
    return response.data;
  },
};