/**
 * Project Store - Zustand
 * Manages project state and operations
 */

import { create } from 'zustand';
import { devtools, subscribeWithSelector } from 'zustand/middleware';
import { Project, ProjectFilter, ProjectCreateInput, ProjectUpdateInput, UUID } from '../types/project.types';

interface ProjectStore {
  // State
  projects: Project[];
  currentProject: Project | null;
  loading: boolean;
  error: string | null;
  filter: ProjectFilter;
  
  // Computed
  filteredProjects: Project[];
  
  // Actions
  setProjects: (projects: Project[]) => void;
  setCurrentProject: (project: Project | null) => void;
  addProject: (project: Project) => void;
  updateProject: (id: UUID, updates: ProjectUpdateInput) => void;
  deleteProject: (id: UUID) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setFilter: (filter: Partial<ProjectFilter>) => void;
  clearFilter: () => void;
  
  // Async Actions
  fetchProjects: () => Promise<void>;
  createProject: (input: ProjectCreateInput) => Promise<Project>;
  updateProjectAsync: (id: UUID, updates: ProjectUpdateInput) => Promise<void>;
  deleteProjectAsync: (id: UUID) => Promise<void>;
}

const defaultFilter: ProjectFilter = {
  searchTerm: '',
  sortBy: 'updatedAt',
  sortOrder: 'desc'
};

export const useProjectStore = create<ProjectStore>()(
  devtools(
    subscribeWithSelector((set, get) => ({
      // Initial State
      projects: [],
      currentProject: null,
      loading: false,
      error: null,
      filter: defaultFilter,
      
      // Computed
      get filteredProjects() {
        const { projects, filter } = get();
        let filtered = [...projects];
        
        // Search filter
        if (filter.searchTerm) {
          const searchLower = filter.searchTerm.toLowerCase();
          filtered = filtered.filter(project => 
            project.title.toLowerCase().includes(searchLower) ||
            project.description.toLowerCase().includes(searchLower)
          );
        }
        
        // Sort
        if (filter.sortBy) {
          filtered.sort((a, b) => {
            const aVal = a[filter.sortBy!];
            const bVal = b[filter.sortBy!];
            
            let comparison = 0;
            if (aVal < bVal) comparison = -1;
            if (aVal > bVal) comparison = 1;
            
            return filter.sortOrder === 'desc' ? -comparison : comparison;
          });
        }
        
        return filtered;
      },
      
      // Sync Actions
      setProjects: (projects) => set({ projects }),
      
      setCurrentProject: (project) => set({ currentProject: project }),
      
      addProject: (project) => set((state) => ({
        projects: [...state.projects, project]
      })),
      
      updateProject: (id, updates) => set((state) => ({
        projects: state.projects.map(project =>
          project.id === id 
            ? { ...project, ...updates, updatedAt: new Date() }
            : project
        ),
        currentProject: state.currentProject?.id === id
          ? { ...state.currentProject, ...updates, updatedAt: new Date() }
          : state.currentProject
      })),
      
      deleteProject: (id) => set((state) => ({
        projects: state.projects.filter(project => project.id !== id),
        currentProject: state.currentProject?.id === id ? null : state.currentProject
      })),
      
      setLoading: (loading) => set({ loading }),
      
      setError: (error) => set({ error }),
      
      setFilter: (newFilter) => set((state) => ({
        filter: { ...state.filter, ...newFilter }
      })),
      
      clearFilter: () => set({ filter: defaultFilter }),
      
      // Async Actions (to be implemented with Convex)
      fetchProjects: async () => {
        set({ loading: true, error: null });
        try {
          // TODO: Implement with Convex
          // const projects = await convex.query(api.projects.list);
          // set({ projects, loading: false });
          console.log('fetchProjects - to be implemented with Convex');
          set({ loading: false });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to fetch projects',
            loading: false 
          });
        }
      },
      
      createProject: async (input) => {
        set({ loading: true, error: null });
        try {
          // TODO: Implement with Convex
          // const project = await convex.mutation(api.projects.create, input);
          // get().addProject(project);
          console.log('createProject - to be implemented with Convex', input);
          set({ loading: false });
          return {} as Project; // Temporary
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to create project',
            loading: false 
          });
          throw error;
        }
      },
      
      updateProjectAsync: async (id, updates) => {
        set({ loading: true, error: null });
        try {
          // TODO: Implement with Convex
          // await convex.mutation(api.projects.update, { id, ...updates });
          get().updateProject(id, updates);
          console.log('updateProjectAsync - to be implemented with Convex', id, updates);
          set({ loading: false });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to update project',
            loading: false 
          });
          throw error;
        }
      },
      
      deleteProjectAsync: async (id) => {
        set({ loading: true, error: null });
        try {
          // TODO: Implement with Convex
          // await convex.mutation(api.projects.delete, { id });
          get().deleteProject(id);
          console.log('deleteProjectAsync - to be implemented with Convex', id);
          set({ loading: false });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to delete project',
            loading: false 
          });
          throw error;
        }
      }
    })),
    { name: 'project-store' }
  )
);