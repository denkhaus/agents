import { create } from "zustand";
import { Agent } from "@/lib/types";

interface AgentsStore {
  // State
  agents: Agent[];
  selectedAgentId: string | null;
  loading: boolean;
  error: string | null;

  // Actions
  setAgents: (agents: Agent[]) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  updateAgent: (agent: Agent) => void;
  removeAgent: (agentId: string) => void;
  setSelectedAgentId: (agentId: string | null) => void;
}

export const useAgentsStore = create<AgentsStore>((set, get) => ({
  // Initial state
  agents: [],
  selectedAgentId:
    typeof window !== "undefined"
      ? localStorage.getItem("selectedAgentId")
      : null,
  loading: false,
  error: null,

  // Actions
  setAgents: (agents) => {
    const currentSelectedAgentId = get().selectedAgentId;
    const newSelectedAgentId = agents.some(
      (agent) => agent.id === currentSelectedAgentId
    )
      ? currentSelectedAgentId
      : null;
    set({ agents, error: null, selectedAgentId: newSelectedAgentId });
    if (typeof window !== "undefined") {
      if (newSelectedAgentId) {
        localStorage.setItem("selectedAgentId", newSelectedAgentId);
      } else {
        localStorage.removeItem("selectedAgentId");
      }
    }
  },

  setLoading: (loading) => set({ loading }),

  setError: (error) => set({ error, loading: false }),

  updateAgent: (updatedAgent) => {
    const agents = get().agents;
    const updatedAgents = agents.map((agent) =>
      agent.id === updatedAgent.id ? updatedAgent : agent
    );
    set({ agents: updatedAgents });
  },

  removeAgent: (agentId) => {
    const agents = get().agents;
    const filteredAgents = agents.filter((agent) => agent.id !== agentId);
    set({ agents: filteredAgents });
    if (get().selectedAgentId === agentId) {
      set({ selectedAgentId: null });
      if (typeof window !== "undefined") {
        localStorage.removeItem("selectedAgentId");
      }
    }
  },

  setSelectedAgentId: (agentId) => {
    set({ selectedAgentId: agentId });
    if (typeof window !== "undefined") {
      if (agentId) {
        localStorage.setItem("selectedAgentId", agentId);
      } else {
        localStorage.removeItem("selectedAgentId");
      }
    }
  },
}));
