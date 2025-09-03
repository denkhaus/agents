export interface Agent {
  id: string
  name: string
  status: 'online' | 'offline' | 'busy'
  capabilities: string[]
  avatar?: string
  lastActivity?: Date
}

export interface AgentInfo {
  name: string
  role?: string
  id?: string
  capabilities?: string[]
  status?: Agent['status']
}

export type AgentStatus = 'online' | 'offline' | 'busy'