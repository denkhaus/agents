import { apiClient } from './client'
import { Agent } from '@/lib/types'

export const agentApi = {
  async getAgents(): Promise<Agent[]> {
    try {
      const agentNames = await apiClient.getAgents()
      
      if (!Array.isArray(agentNames)) {
        console.error('Expected array, got:', typeof agentNames, agentNames)
        return []
      }
      
      // Convert agent names to Agent objects with default properties
      const agents = agentNames.map((name) => ({
        id: name,
        name: name,
        status: 'online' as const,
        capabilities: [],
        avatar: undefined,
        lastActivity: new Date()
      }))
      
      return agents
    } catch (error) {
      console.error('Error fetching agents:', error)
      throw error
    }
  }
}