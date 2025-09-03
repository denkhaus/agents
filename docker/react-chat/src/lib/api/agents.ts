import { apiClient } from './client'
import { Agent } from '@/lib/types'

export const agentApi = {
  async getAgents(): Promise<Agent[]> {
    console.log('agentApi.getAgents: Starting request')
    const agentNames = await apiClient.getAgents()
    console.log('agentApi.getAgents: Received agent names:', agentNames)
    
    // Convert agent names to Agent objects with default properties
    const agents = agentNames.map((name) => ({
      id: name,
      name: name,
      status: 'online' as const,
      capabilities: [],
      avatar: undefined,
      lastActivity: new Date()
    }))
    
    console.log('agentApi.getAgents: Converted to agents:', agents)
    return agents
  }
}