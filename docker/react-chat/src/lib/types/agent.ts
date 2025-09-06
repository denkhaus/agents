import { AgentId } from '../constants/agents';

export interface Agent {
  id: AgentId
  name: string
  status: 'online' | 'offline' | 'busy'
  capabilities: string[]
  avatar?: string
  lastActivity?: Date
}

export interface AgentInfo {
  name: string
  role?: string
  id?: AgentId
  capabilities?: string[]
  status?: Agent['status']
}

export type AgentStatus = 'online' | 'offline' | 'busy'