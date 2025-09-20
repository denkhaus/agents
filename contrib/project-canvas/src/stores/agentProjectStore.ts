/**
 * Agent Project Store - Zustand State Management
 * Handles agent project configurations and canvas layouts
 */

import { create } from "zustand";
import { devtools, subscribeWithSelector } from "zustand/middleware";
import type {
  AgentProject,
  AgentProjectFilter,
  AgentProjectUpdateInput,
  AgentNode,
  AgentConnection,
  UUID,
} from "@/types";

interface AgentProjectState {
  // State
  agentProjects: AgentProject[];
  currentAgentProject: AgentProject | null;
  loading: boolean;
  error: string | null;
  filter: AgentProjectFilter;

  // Actions
  setAgentProjects: (projects: AgentProject[]) => void;
  setCurrentAgentProject: (project: AgentProject | null) => void;
  addAgentProject: (project: AgentProject) => void;
  updateAgentProject: (id: UUID, updates: AgentProjectUpdateInput) => void;
  deleteAgentProject: (id: UUID) => void;
  duplicateAgentProject: (id: UUID) => void;

  // Agent Node Management
  addAgentToProject: (projectId: UUID, agentNode: AgentNode) => void;
  removeAgentFromProject: (projectId: UUID, agentId: UUID) => void;

  // Connection Management
  addConnection: (projectId: UUID, connection: AgentConnection) => void;
  removeConnection: (projectId: UUID, connectionId: UUID) => void;
  updateConnection: (
    projectId: UUID,
    connectionId: UUID,
    updates: Partial<AgentConnection>
  ) => void;

  // Utility
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setFilter: (filter: Partial<AgentProjectFilter>) => void;
  clearFilter: () => void;

  // Computed
  getAgentProjectById: (id: UUID) => AgentProject | undefined;
  getFilteredAgentProjects: () => AgentProject[];
  getAgentProjectByName: (name: string) => AgentProject | undefined;
}

const initialFilter: AgentProjectFilter = {
  searchTerm: "",
  hasAgents: undefined,
};

export const useAgentProjectStore = create<AgentProjectState>()(
  devtools(
    subscribeWithSelector((set, get) => ({
      // Initial State
      agentProjects: [],
      currentAgentProject: null,
      loading: false,
      error: null,
      filter: initialFilter,

      // Actions
      setAgentProjects: (agentProjects) =>
        set({ agentProjects }, false, "agentProject/setAgentProjects"),

      setCurrentAgentProject: (currentAgentProject) =>
        set(
          { currentAgentProject },
          false,
          "agentProject/setCurrentAgentProject"
        ),

      addAgentProject: (project) =>
        set(
          (state) => ({
            agentProjects: [...state.agentProjects, project],
          }),
          false,
          "agentProject/addAgentProject"
        ),

      updateAgentProject: (id, updates) =>
        set(
          (state) => ({
            agentProjects: state.agentProjects.map((project) =>
              project.id === id
                ? { ...project, ...updates, updatedAt: new Date() }
                : project
            ),
            currentAgentProject:
              state.currentAgentProject?.id === id
                ? {
                    ...state.currentAgentProject,
                    ...updates,
                    updatedAt: new Date(),
                  }
                : state.currentAgentProject,
          }),
          false,
          "agentProject/updateAgentProject"
        ),

      deleteAgentProject: (id) =>
        set(
          (state) => ({
            agentProjects: state.agentProjects.filter(
              (project) => project.id !== id
            ),
            currentAgentProject:
              state.currentAgentProject?.id === id
                ? null
                : state.currentAgentProject,
          }),
          false,
          "agentProject/deleteAgentProject"
        ),

      duplicateAgentProject: (id) =>
        set(
          (state) => {
            const project = state.agentProjects.find((p) => p.id === id);
            if (!project) return state;

            const duplicatedProject: AgentProject = {
              ...project,
              id: crypto.randomUUID(),
              name: `${project.name} (Copy)`,
              createdAt: new Date(),
              updatedAt: new Date(),
              agentNodes: project.agentNodes.map((node) => ({
                ...node,
                id: crypto.randomUUID(),
              })),
              connections: project.connections.map((conn) => ({
                ...conn,
                id: crypto.randomUUID(),
              })),
            };

            return {
              agentProjects: [...state.agentProjects, duplicatedProject],
            };
          },
          false,
          "agentProject/duplicateAgentProject"
        ),

      // Agent Node Management
      addAgentToProject: (projectId, agentNode) =>
        set(
          (state) => ({
            agentProjects: state.agentProjects.map((project) =>
              project.id === projectId
                ? {
                    ...project,
                    agentNodes: [...project.agentNodes, agentNode],
                    updatedAt: new Date(),
                  }
                : project
            ),
            currentAgentProject:
              state.currentAgentProject?.id === projectId
                ? {
                    ...state.currentAgentProject,
                    agentNodes: [
                      ...state.currentAgentProject.agentNodes,
                      agentNode,
                    ],
                    updatedAt: new Date(),
                  }
                : state.currentAgentProject,
          }),
          false,
          "agentProject/addAgentToProject"
        ),

      removeAgentFromProject: (projectId, agentId) =>
        set(
          (state) => ({
            agentProjects: state.agentProjects.map((project) =>
              project.id === projectId
                ? {
                    ...project,
                    agentNodes: project.agentNodes.filter(
                      (node) => node.data.agent.id !== agentId
                    ),
                    connections: project.connections.filter(
                      (conn) =>
                        conn.source !== agentId && conn.target !== agentId
                    ),
                    updatedAt: new Date(),
                  }
                : project
            ),
            currentAgentProject:
              state.currentAgentProject?.id === projectId
                ? {
                    ...state.currentAgentProject,
                    agentNodes: state.currentAgentProject.agentNodes.filter(
                      (node) => node.data.agent.id !== agentId
                    ),
                    connections: state.currentAgentProject.connections.filter(
                      (conn) =>
                        conn.source !== agentId && conn.target !== agentId
                    ),
                    updatedAt: new Date(),
                  }
                : state.currentAgentProject,
          }),
          false,
          "agentProject/removeAgentFromProject"
        ),

      // Connection Management
      addConnection: (projectId, connection) =>
        set(
          (state) => ({
            agentProjects: state.agentProjects.map((project) =>
              project.id === projectId
                ? {
                    ...project,
                    connections: [...project.connections, connection],
                    updatedAt: new Date(),
                  }
                : project
            ),
            currentAgentProject:
              state.currentAgentProject?.id === projectId
                ? {
                    ...state.currentAgentProject,
                    connections: [
                      ...state.currentAgentProject.connections,
                      connection,
                    ],
                    updatedAt: new Date(),
                  }
                : state.currentAgentProject,
          }),
          false,
          "agentProject/addConnection"
        ),

      removeConnection: (projectId, connectionId) =>
        set(
          (state) => ({
            agentProjects: state.agentProjects.map((project) =>
              project.id === projectId
                ? {
                    ...project,
                    connections: project.connections.filter(
                      (conn) => conn.id !== connectionId
                    ),
                    updatedAt: new Date(),
                  }
                : project
            ),
            currentAgentProject:
              state.currentAgentProject?.id === projectId
                ? {
                    ...state.currentAgentProject,
                    connections: state.currentAgentProject.connections.filter(
                      (conn) => conn.id !== connectionId
                    ),
                    updatedAt: new Date(),
                  }
                : state.currentAgentProject,
          }),
          false,
          "agentProject/removeConnection"
        ),

      updateConnection: (projectId, connectionId, updates) =>
        set(
          (state) => ({
            agentProjects: state.agentProjects.map((project) =>
              project.id === projectId
                ? {
                    ...project,
                    connections: project.connections.map((conn) =>
                      conn.id === connectionId ? { ...conn, ...updates } : conn
                    ),
                    updatedAt: new Date(),
                  }
                : project
            ),
            currentAgentProject:
              state.currentAgentProject?.id === projectId
                ? {
                    ...state.currentAgentProject,
                    connections: state.currentAgentProject.connections.map(
                      (conn) =>
                        conn.id === connectionId
                          ? { ...conn, ...updates }
                          : conn
                    ),
                    updatedAt: new Date(),
                  }
                : state.currentAgentProject,
          }),
          false,
          "agentProject/updateConnection"
        ),

      // Utility
      setLoading: (loading) =>
        set({ loading }, false, "agentProject/setLoading"),
      setError: (error) => set({ error }, false, "agentProject/setError"),
      setFilter: (filter) =>
        set(
          (state) => ({ filter: { ...state.filter, ...filter } }),
          false,
          "agentProject/setFilter"
        ),
      clearFilter: () =>
        set({ filter: initialFilter }, false, "agentProject/clearFilter"),

      // Computed
      getAgentProjectById: (id) => {
        const state = get();
        return state.agentProjects.find((project) => project.id === id);
      },

      getFilteredAgentProjects: () => {
        const state = get();
        const { agentProjects, filter } = state;

        return agentProjects.filter((project) => {
          // Search term filter
          if (filter.searchTerm) {
            const searchLower = filter.searchTerm.toLowerCase();
            const matchesName = project.name
              .toLowerCase()
              .includes(searchLower);
            const matchesDescription = project.description
              .toLowerCase()
              .includes(searchLower);
            if (!matchesName && !matchesDescription) return false;
          }

          // Has agents filter
          if (filter.hasAgents !== undefined) {
            const hasAgents = project.agentNodes.length > 0;
            if (filter.hasAgents !== hasAgents) return false;
          }

          return true;
        });
      },

      getAgentProjectByName: (name) => {
        const state = get();
        return state.agentProjects.find((project) => project.name === name);
      },
    })),
    { name: "agent-project-store" }
  )
);
