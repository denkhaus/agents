import { create } from "zustand";
import { AgentInfo } from "@/lib/types";
import { AgentId } from "@/lib/constants/agents";

interface AgentsStore {
  // State
  agents: AgentInfo[];
  selectedAgentId: AgentId | null;
  loading: boolean;
  error: string | null;

  // Actions
  setAgents: (agents: AgentInfo[]) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  updateAgent: (agent: AgentInfo) => void;
  removeAgent: (agentId: AgentId) => void;
  setSelectedAgentId: (agentId: AgentId | null) => void;
}

export const useAgentsStore = create<AgentsStore>((set, get) => ({
  // Initial state
  agents: [],
  selectedAgentId:
    typeof window !== "undefined"
      ? (localStorage.getItem("selectedAgentId") as AgentId)
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

  removeAgent: (agentId: AgentId) => {
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

  setSelectedAgentId: (agentId: AgentId | null) => {
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
