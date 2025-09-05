import { AgentInfo } from './agent'

export interface ApiResponse<T = unknown> {
  success: boolean
  data?: T
  error?: string
}

export interface AgentRunRequest {
  appName: string
  userID: string
  sessionID: string
  streaming: boolean
  newMessage: {
    role: string
    parts: Array<{
      text: string
    }>
  }
}

export interface MultiChatRequest {
  fromAgent: string
  toAgent: string
  message: string
  sessionId: string
  userId: string
}

export interface MultiChatResponse {
  success: boolean
  message?: string
  error?: string
}

export interface ADKSession {
  appName: string
  userID: string
  id: string
  createTime: number
  lastUpdateTime: number
  state: Record<string, unknown>
  events: unknown[]
}

/**
 * Represents all types of events in the agent system.
 * This includes:
 * - Regular agent responses to user messages (object: 'message', 'tool_call', etc.)
 * - Inter-agent communication (type: 'inter_agent', 'communication')
 * - System events (type: 'heartbeat', 'agent_list', 'system')
 * 
 * Despite the legacy name, this type handles ALL agent events, not just inter-agent ones.
 */
export interface AgentEvent {
  id?: string
  type: 'communication' | 'heartbeat' | 'agent_list' | 'inter_agent' | 'system'
  fromAgent?: string
  toAgent?: string
  content?: string | {
    role?: string
    parts?: Array<{
      text?: string
      functionCall?: {
        name: string
        args: unknown
        id?: string
      }
      functionResponse?: {
        name: string
        response: unknown
        id?: string
      }
    }>
  }
  message?: string
  timestamp: number
  agents?: AgentInfo[]
  interAgent?: {
    fromAgent: string
    toAgent: string
    type: string
  }
  invocationId?: string
  author?: string
  done?: boolean
  partial?: boolean
  usageMetadata?: {
    promptTokenCount?: number
    candidatesTokenCount?: number
    totalTokenCount?: number
  }
  // Additional backend fields
  object?: string
  model?: string
  created?: number
}

/**
 * @deprecated Use AgentEvent instead. This alias is kept for backward compatibility.
 */
export type InterAgentEvent = AgentEvent

export interface SSEEventData {
  type: string
  data: AgentEvent
}