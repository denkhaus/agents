import { create } from 'zustand'
import { Agent } from '@/lib/types'

interface AgentsStore {
  // State
  agents: Agent[]
  loading: boolean
  error: string | null
  
  // Actions
  setAgents: (agents: Agent[]) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  updateAgent: (agent: Agent) => void
  removeAgent: (agentId: string) => void
}

export const useAgentsStore = create<AgentsStore>((set, get) => ({
  // Initial state
  agents: [],
  loading: false,
  error: null,
  
  // Actions
  setAgents: (agents) => set({ agents, error: null }),
  
  setLoading: (loading) => set({ loading }),
  
  setError: (error) => set({ error, loading: false }),
  
  updateAgent: (updatedAgent) => {
    const agents = get().agents
    const updatedAgents = agents.map(agent =>
      agent.id === updatedAgent.id ? updatedAgent : agent
    )
    set({ agents: updatedAgents })
  },
  
  removeAgent: (agentId) => {
    const agents = get().agents
    const filteredAgents = agents.filter(agent => agent.id !== agentId)
    set({ agents: filteredAgents })
  }
}))