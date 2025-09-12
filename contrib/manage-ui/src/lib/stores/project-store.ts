import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import type { 
  Project, 
  Task, 
  ProjectProgress, 
  TaskFilter, 
  ProjectWithTasks,
  TaskHierarchy 
} from '../types';

interface ProjectState {
  // Data
  projects: Project[];
  currentProject: ProjectWithTasks | null;
  tasks: Task[];
  projectProgress: Record<string, ProjectProgress>;
  
  // Loading states
  isLoading: boolean;
  isLoadingTasks: boolean;
  isLoadingProgress: boolean;
  
  // Error states
  error: string | null;
  taskError: string | null;
  
  // Filters and search
  filter: TaskFilter;
  searchQuery: string;
  
  // Actions
  setProjects: (projects: Project[]) => void;
  setCurrentProject: (project: ProjectWithTasks | null) => void;
  addProject: (project: Project) => void;
  updateProject: (projectId: string, updates: Partial<Project>) => void;
  removeProject: (projectId: string) => void;
  
  setTasks: (tasks: Task[]) => void;
  addTask: (task: Task) => void;
  updateTask: (taskId: string, updates: Partial<Task>) => void;
  removeTask: (taskId: string) => void;
  bulkUpdateTasks: (taskIds: string[], updates: Partial<Task>) => void;
  
  setProjectProgress: (projectId: string, progress: ProjectProgress) => void;
  
  setFilter: (filter: Partial<TaskFilter>) => void;
  clearFilter: () => void;
  setSearchQuery: (query: string) => void;
  
  setLoading: (loading: boolean) => void;
  setTasksLoading: (loading: boolean) => void;
  setProgressLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setTaskError: (error: string | null) => void;
  
  // Computed getters
  getFilteredTasks: () => Task[];
  getTasksByParent: (parentId: string | null) => Task[];
  getTaskHierarchy: (projectId: string) => TaskHierarchy[];
  getRootTasks: (projectId: string) => Task[];
  getTaskById: (taskId: string) => Task | undefined;
  getProjectById: (projectId: string) => Project | undefined;
}

const initialFilter: TaskFilter = {
  project_id: undefined,
  parent_id: undefined,
  state: undefined,
  min_depth: undefined,
  max_depth: undefined,
  min_complexity: undefined,
  max_complexity: undefined,
};

export const useProjectStore = create<ProjectState>()(
  devtools(
    (set, get) => ({
      // Initial state
      projects: [],
      currentProject: null,
      tasks: [],
      projectProgress: {},
      isLoading: false,
      isLoadingTasks: false,
      isLoadingProgress: false,
      error: null,
      taskError: null,
      filter: initialFilter,
      searchQuery: '',

      // Actions
      setProjects: (projects) => set({ projects }),
      
      setCurrentProject: (project) => set({ currentProject: project }),
      
      addProject: (project) => 
        set((state) => ({ 
          projects: [...state.projects, project] 
        })),
      
      updateProject: (projectId, updates) =>
        set((state) => ({
          projects: state.projects.map((p) =>
            p.id === projectId ? { ...p, ...updates } : p
          ),
          currentProject: state.currentProject?.id === projectId 
            ? { ...state.currentProject, ...updates }
            : state.currentProject
        })),
      
      removeProject: (projectId) =>
        set((state) => ({
          projects: state.projects.filter((p) => p.id !== projectId),
          currentProject: state.currentProject?.id === projectId 
            ? null 
            : state.currentProject,
          tasks: state.tasks.filter((t) => t.project_id !== projectId)
        })),

      setTasks: (tasks) => set({ tasks }),
      
      addTask: (task) =>
        set((state) => ({ 
          tasks: [...state.tasks, task] 
        })),
      
      updateTask: (taskId, updates) =>
        set((state) => ({
          tasks: state.tasks.map((t) =>
            t.id === taskId ? { ...t, ...updates } : t
          )
        })),
      
      removeTask: (taskId) =>
        set((state) => ({
          tasks: state.tasks.filter((t) => t.id !== taskId)
        })),
      
      bulkUpdateTasks: (taskIds, updates) =>
        set((state) => ({
          tasks: state.tasks.map((t) =>
            taskIds.includes(t.id) ? { ...t, ...updates } : t
          )
        })),

      setProjectProgress: (projectId, progress) =>
        set((state) => ({
          projectProgress: {
            ...state.projectProgress,
            [projectId]: progress
          }
        })),

      setFilter: (newFilter) =>
        set((state) => ({
          filter: { ...state.filter, ...newFilter }
        })),
      
      clearFilter: () => set({ filter: initialFilter }),
      
      setSearchQuery: (query) => set({ searchQuery: query }),

      setLoading: (loading) => set({ isLoading: loading }),
      setTasksLoading: (loading) => set({ isLoadingTasks: loading }),
      setProgressLoading: (loading) => set({ isLoadingProgress: loading }),
      setError: (error) => set({ error }),
      setTaskError: (error) => set({ taskError: error }),

      // Computed getters
      getFilteredTasks: () => {
        const { tasks, filter, searchQuery } = get();
        
        return tasks.filter((task) => {
          // Search filter
          if (searchQuery && !task.title.toLowerCase().includes(searchQuery.toLowerCase()) &&
              !task.description.toLowerCase().includes(searchQuery.toLowerCase())) {
            return false;
          }
          
          // Apply filters
          if (filter.project_id && task.project_id !== filter.project_id) return false;
          if (filter.parent_id !== undefined && task.parent_id !== filter.parent_id) return false;
          if (filter.state && task.state !== filter.state) return false;
          if (filter.min_depth !== undefined && task.depth < filter.min_depth) return false;
          if (filter.max_depth !== undefined && task.depth > filter.max_depth) return false;
          if (filter.min_complexity !== undefined && task.complexity < filter.min_complexity) return false;
          if (filter.max_complexity !== undefined && task.complexity > filter.max_complexity) return false;
          
          return true;
        });
      },

      getTasksByParent: (parentId) => {
        const { tasks } = get();
        return tasks.filter((task) => task.parent_id === parentId);
      },

      getTaskHierarchy: (projectId) => {
        const { tasks } = get();
        const projectTasks = tasks.filter((task) => task.project_id === projectId);
        
        const buildHierarchy = (parentId: string | null, level: number = 0): TaskHierarchy[] => {
          const children = projectTasks.filter((task) => task.parent_id === parentId);
          
          return children.map((task) => ({
            task,
            children: buildHierarchy(task.id, level + 1),
            level
          }));
        };
        
        return buildHierarchy(null);
      },

      getRootTasks: (projectId) => {
        const { tasks } = get();
        return tasks.filter((task) => 
          task.project_id === projectId && task.parent_id === null
        );
      },

      getTaskById: (taskId) => {
        const { tasks } = get();
        return tasks.find((task) => task.id === taskId);
      },

      getProjectById: (projectId) => {
        const { projects } = get();
        return projects.find((project) => project.id === projectId);
      },
    }),
    { name: 'project-store' }
  )
);