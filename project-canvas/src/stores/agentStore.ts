/**
 * Agent Store - Zustand State Management
 * Handles agent data for future workspace extension
 */

import { create } from "zustand";
import { devtools, subscribeWithSelector } from "zustand/middleware";
import type {
  Agent,
  AgentRole,
  AgentStatus,
  AgentFilter,
  AgentUpdateInput,
  UUID,
} from "@/types";

interface AgentState {
  // State
  agents: Agent[];
  loading: boolean;
  error: string | null;
  filter: AgentFilter;

  // Actions
  setAgents: (agents: Agent[]) => void;
  addAgent: (agent: Agent) => void;
  updateAgent: (id: UUID, updates: AgentUpdateInput) => void;
  deleteAgent: (id: UUID) => void;
  updateAgentStatus: (id: UUID, status: AgentStatus) => void;
  assignTaskToAgent: (agentId: UUID, taskId: UUID) => void;
  unassignTaskFromAgent: (agentId: UUID, taskId: UUID) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setFilter: (filter: Partial<AgentFilter>) => void;
  clearFilter: () => void;

  // Computed
  getAgentById: (id: UUID) => Agent | undefined;
  getAgentsByRole: (role: AgentRole) => Agent[];
  getAgentsByStatus: (status: AgentStatus) => Agent[];
  getAvailableAgents: () => Agent[];
  getFilteredAgents: () => Agent[];
  getAgentWorkload: (id: UUID) => number;
}

const initialFilter: AgentFilter = {
  searchTerm: "",
  role: undefined,
  status: undefined,
  hasActiveTasks: undefined,
};

export const useAgentStore = create<AgentState>()(
  devtools(
    subscribeWithSelector((set, get) => ({
      // Initial State
      agents: [],
      loading: false,
      error: null,
      filter: initialFilter,

      // Actions
      setAgents: (agents) => set({ agents }, false, "setAgents"),

      addAgent: (agent) =>
        set(
          (state) => ({ agents: [...state.agents, agent] }),
          false,
          "addAgent"
        ),

      updateAgent: (id, updates) =>
        set(
          (state) => ({
            agents: state.agents.map((agent) =>
              agent.id === id
                ? { ...agent, ...updates, updatedAt: new Date() }
                : agent
            ),
          }),
          false,
          "updateAgent"
        ),

      deleteAgent: (id) =>
        set(
          (state) => ({
            agents: state.agents.filter((agent) => agent.id !== id),
          }),
          false,
          "deleteAgent"
        ),

      updateAgentStatus: (id, status) => {
        const updates: AgentUpdateInput = { status };
        if (status === AgentStatus.ONLINE) {
          // Update lastActiveAt when going online
          (updates as any).lastActiveAt = new Date();
        }
        get().updateAgent(id, updates);
      },

      assignTaskToAgent: (agentId, taskId) =>
        set(
          (state) => ({
            agents: state.agents.map((agent) =>
              agent.id === agentId
                ? {
                    ...agent,
                    currentTasks: [...agent.currentTasks, taskId],
                    updatedAt: new Date(),
                  }
                : agent
            ),
          }),
          false,
          "assignTaskToAgent"
        ),

      unassignTaskFromAgent: (agentId, taskId) =>
        set(
          (state) => ({
            agents: state.agents.map((agent) =>
              agent.id === agentId
                ? {
                    ...agent,
                    currentTasks: agent.currentTasks.filter(
                      (id) => id !== taskId
                    ),
                    updatedAt: new Date(),
                  }
                : agent
            ),
          }),
          false,
          "unassignTaskFromAgent"
        ),

      setLoading: (loading) => set({ loading }, false, "setLoading"),
      setError: (error) => set({ error }, false, "setError"),
      setFilter: (newFilter) =>
        set(
          (state) => ({ filter: { ...state.filter, ...newFilter } }),
          false,
          "setFilter"
        ),
      clearFilter: () => set({ filter: initialFilter }, false, "clearFilter"),

      // Computed
      getAgentById: (id) => get().agents.find((agent) => agent.id === id),

      getAgentsByRole: (role) =>
        get().agents.filter((agent) => agent.role === role),

      getAgentsByStatus: (status) =>
        get().agents.filter((agent) => agent.status === status),

      getAvailableAgents: () => {
        return get().agents.filter(
          (agent) =>
            agent.status === AgentStatus.ONLINE ||
            agent.status === AgentStatus.IDLE
        );
      },

      getFilteredAgents: () => {
        const { agents, filter } = get();
        let filtered = [...agents];

        // Role filter
        if (filter.role) {
          filtered = filtered.filter((agent) => agent.role === filter.role);
        }

        // Status filter
        if (filter.status) {
          filtered = filtered.filter((agent) => agent.status === filter.status);
        }

        // Active tasks filter
        if (filter.hasActiveTasks !== undefined) {
          filtered = filtered.filter((agent) =>
            filter.hasActiveTasks
              ? agent.currentTasks.length > 0
              : agent.currentTasks.length === 0
          );
        }

        // Search filter
        if (filter.searchTerm) {
          const searchTerm = filter.searchTerm.toLowerCase();
          filtered = filtered.filter(
            (agent) =>
              agent.name.toLowerCase().includes(searchTerm) ||
              agent.description.toLowerCase().includes(searchTerm) ||
              agent.role.toLowerCase().includes(searchTerm)
          );
        }

        return filtered;
      },

      getAgentWorkload: (id) => {
        const agent = get().getAgentById(id);
        return agent?.currentTasks.length || 0;
      },
    })),
    { name: "agent-store" }
  )
);
