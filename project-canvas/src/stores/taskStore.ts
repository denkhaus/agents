/**
 * Task Store - Zustand
 * Manages task state and operations
 */

import { create } from 'zustand';
import { devtools, subscribeWithSelector } from 'zustand/middleware';
import { Task, TaskFilter, TaskCreateInput, TaskUpdates, TaskState, UUID } from '../types/task.types';
import { Position } from '../types/ui.types';

interface TaskStore {
  // State
  tasks: Task[];
  tasksByProject: Record<UUID, Task[]>;
  loading: boolean;
  error: string | null;
  filter: TaskFilter;
  
  // Computed
  filteredTasks: Task[];
  rootTasks: Task[];
  taskHierarchy: Record<UUID, Task[]>;
  
  // Actions
  setTasks: (tasks: Task[]) => void;
  addTask: (task: Task) => void;
  updateTask: (id: UUID, updates: TaskUpdates) => void;
  deleteTask: (id: UUID) => void;
  updateTaskPosition: (id: UUID, position: Position) => void;
  updateTaskState: (id: UUID, state: TaskState) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setFilter: (filter: Partial<TaskFilter>) => void;
  clearFilter: () => void;
  
  // Dependency Management
  addDependency: (taskId: UUID, dependsOnTaskId: UUID) => void;
  removeDependency: (taskId: UUID, dependsOnTaskId: UUID) => void;
  getDependencies: (taskId: UUID) => Task[];
  getDependents: (taskId: UUID) => Task[];
  
  // Hierarchy Management
  getChildTasks: (parentId: UUID) => Task[];
  getParentTask: (taskId: UUID) => Task | null;
  getTaskDepth: (taskId: UUID) => number;
  
  // Async Actions
  fetchTasksByProject: (projectId: UUID) => Promise<void>;
  createTask: (input: TaskCreateInput) => Promise<Task>;
  updateTaskAsync: (id: UUID, updates: TaskUpdates) => Promise<void>;
  deleteTaskAsync: (id: UUID) => Promise<void>;
}

const defaultFilter: TaskFilter = {
  searchTerm: '',
  state: undefined,
  minDepth: undefined,
  maxDepth: undefined
};

export const useTaskStore = create<TaskStore>()(
  devtools(
    subscribeWithSelector((set, get) => ({
      // Initial State
      tasks: [],
      tasksByProject: {},
      loading: false,
      error: null,
      filter: defaultFilter,
      
      // Computed
      get filteredTasks() {
        const { tasks, filter } = get();
        let filtered = [...tasks];
        
        // Project filter
        if (filter.projectId) {
          filtered = filtered.filter(task => task.projectId === filter.projectId);
        }
        
        // State filter
        if (filter.state) {
          filtered = filtered.filter(task => task.state === filter.state);
        }
        
        // Depth filter
        if (filter.minDepth !== undefined) {
          filtered = filtered.filter(task => task.depth >= filter.minDepth!);
        }
        if (filter.maxDepth !== undefined) {
          filtered = filtered.filter(task => task.depth <= filter.maxDepth!);
        }
        
        // Complexity filter
        if (filter.minComplexity !== undefined) {
          filtered = filtered.filter(task => task.complexity >= filter.minComplexity!);
        }
        if (filter.maxComplexity !== undefined) {
          filtered = filtered.filter(task => task.complexity <= filter.maxComplexity!);
        }
        
        // Search filter
        if (filter.searchTerm) {
          const searchLower = filter.searchTerm.toLowerCase();
          filtered = filtered.filter(task => 
            task.title.toLowerCase().includes(searchLower) ||
            task.description.toLowerCase().includes(searchLower)
          );
        }
        
        // Assigned agent filter
        if (filter.assignedAgent) {
          filtered = filtered.filter(task => task.assignedAgent === filter.assignedAgent);
        }
        
        return filtered;
      },
      
      get rootTasks() {
        const { tasks, filter } = get();
        const projectTasks = filter.projectId 
          ? tasks.filter(task => task.projectId === filter.projectId)
          : tasks;
        return projectTasks.filter(task => !task.parentId);
      },
      
      get taskHierarchy() {
        const { tasks } = get();
        const hierarchy: Record<UUID, Task[]> = {};
        
        tasks.forEach(task => {
          if (task.parentId) {
            if (!hierarchy[task.parentId]) {
              hierarchy[task.parentId] = [];
            }
            hierarchy[task.parentId].push(task);
          }
        });
        
        return hierarchy;
      },
      
      // Sync Actions
      setTasks: (tasks) => {
        const tasksByProject: Record<UUID, Task[]> = {};
        tasks.forEach(task => {
          if (!tasksByProject[task.projectId]) {
            tasksByProject[task.projectId] = [];
          }
          tasksByProject[task.projectId].push(task);
        });
        
        set({ tasks, tasksByProject });
      },
      
      addTask: (task) => set((state) => {
        const newTasks = [...state.tasks, task];
        const newTasksByProject = { ...state.tasksByProject };
        
        if (!newTasksByProject[task.projectId]) {
          newTasksByProject[task.projectId] = [];
        }
        newTasksByProject[task.projectId].push(task);
        
        return { tasks: newTasks, tasksByProject: newTasksByProject };
      }),
      
      updateTask: (id, updates) => set((state) => {
        const updatedTasks = state.tasks.map(task =>
          task.id === id 
            ? { ...task, ...updates, updatedAt: new Date() }
            : task
        );
        
        // Rebuild tasksByProject if projectId changed
        const tasksByProject: Record<UUID, Task[]> = {};
        updatedTasks.forEach(task => {
          if (!tasksByProject[task.projectId]) {
            tasksByProject[task.projectId] = [];
          }
          tasksByProject[task.projectId].push(task);
        });
        
        return { tasks: updatedTasks, tasksByProject };
      }),
      
      deleteTask: (id) => set((state) => {
        const filteredTasks = state.tasks.filter(task => task.id !== id);
        
        // Rebuild tasksByProject
        const tasksByProject: Record<UUID, Task[]> = {};
        filteredTasks.forEach(task => {
          if (!tasksByProject[task.projectId]) {
            tasksByProject[task.projectId] = [];
          }
          tasksByProject[task.projectId].push(task);
        });
        
        return { tasks: filteredTasks, tasksByProject };
      }),
      
      updateTaskPosition: (id, position) => {
        get().updateTask(id, { position });
      },
      
      updateTaskState: (id, state) => {
        const updates: TaskUpdates = { state };
        if (state === TaskState.COMPLETED) {
          updates.completedAt = new Date();
        }
        get().updateTask(id, updates);
      },
      
      setLoading: (loading) => set({ loading }),
      setError: (error) => set({ error }),
      
      setFilter: (newFilter) => set((state) => ({
        filter: { ...state.filter, ...newFilter }
      })),
      
      clearFilter: () => set({ filter: defaultFilter }),
      
      // Dependency Management
      addDependency: (taskId, dependsOnTaskId) => {
        const { tasks } = get();
        const task = tasks.find(t => t.id === taskId);
        const dependsOnTask = tasks.find(t => t.id === dependsOnTaskId);
        
        if (task && dependsOnTask && !task.dependencies.includes(dependsOnTaskId)) {
          get().updateTask(taskId, {
            dependencies: [...task.dependencies, dependsOnTaskId]
          });
          get().updateTask(dependsOnTaskId, {
            dependents: [...dependsOnTask.dependents, taskId]
          });
        }
      },
      
      removeDependency: (taskId, dependsOnTaskId) => {
        const { tasks } = get();
        const task = tasks.find(t => t.id === taskId);
        const dependsOnTask = tasks.find(t => t.id === dependsOnTaskId);
        
        if (task && dependsOnTask) {
          get().updateTask(taskId, {
            dependencies: task.dependencies.filter(id => id !== dependsOnTaskId)
          });
          get().updateTask(dependsOnTaskId, {
            dependents: dependsOnTask.dependents.filter(id => id !== taskId)
          });
        }
      },
      
      getDependencies: (taskId) => {
        const { tasks } = get();
        const task = tasks.find(t => t.id === taskId);
        if (!task) return [];
        
        return tasks.filter(t => task.dependencies.includes(t.id));
      },
      
      getDependents: (taskId) => {
        const { tasks } = get();
        const task = tasks.find(t => t.id === taskId);
        if (!task) return [];
        
        return tasks.filter(t => task.dependents.includes(t.id));
      },
      
      // Hierarchy Management
      getChildTasks: (parentId) => {
        const { tasks } = get();
        return tasks.filter(task => task.parentId === parentId);
      },
      
      getParentTask: (taskId) => {
        const { tasks } = get();
        const task = tasks.find(t => t.id === taskId);
        if (!task || !task.parentId) return null;
        
        return tasks.find(t => t.id === task.parentId) || null;
      },
      
      getTaskDepth: (taskId) => {
        const { tasks } = get();
        const task = tasks.find(t => t.id === taskId);
        return task ? task.depth : 0;
      },
      
      // Async Actions (to be implemented with Convex)
      fetchTasksByProject: async (projectId) => {
        set({ loading: true, error: null });
        try {
          // TODO: Implement with Convex
          console.log('fetchTasksByProject - to be implemented with Convex', projectId);
          set({ loading: false });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to fetch tasks',
            loading: false 
          });
        }
      },
      
      createTask: async (input) => {
        set({ loading: true, error: null });
        try {
          // TODO: Implement with Convex
          console.log('createTask - to be implemented with Convex', input);
          set({ loading: false });
          return {} as Task; // Temporary
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to create task',
            loading: false 
          });
          throw error;
        }
      },
      
      updateTaskAsync: async (id, updates) => {
        set({ loading: true, error: null });
        try {
          // TODO: Implement with Convex
          get().updateTask(id, updates);
          console.log('updateTaskAsync - to be implemented with Convex', id, updates);
          set({ loading: false });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to update task',
            loading: false 
          });
          throw error;
        }
      },
      
      deleteTaskAsync: async (id) => {
        set({ loading: true, error: null });
        try {
          // TODO: Implement with Convex
          get().deleteTask(id);
          console.log('deleteTaskAsync - to be implemented with Convex', id);
          set({ loading: false });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to delete task',
            loading: false 
          });
          throw error;
        }
      }
    })),
    { name: 'task-store' }
  )
);