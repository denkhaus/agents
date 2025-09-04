import { apiClient } from './client'
import { Message } from '@/lib/types'

export interface BackendSession {
  id: string
  appName: string
  userId: string
  createTime: number
  lastUpdateTime: number
  state: Record<string, any>
  events: BackendEvent[]
}

export interface BackendEvent {
  id: string
  invocationId: string
  author: string
  timestamp: number
  content: {
    role: string
    parts: Array<{
      text?: string
      functionCall?: {
        name: string
        args: any
        id?: string
      }
      functionResponse?: {
        name: string
        response: any
        id?: string
      }
    }>
  }
  done: boolean
  partial: boolean
  object?: string
}

export const sessionApi = {
  async getSessions(appName: string, userId: string = 'user'): Promise<BackendSession[]> {
    try {
      const response = await apiClient.getSessions(appName, userId)
      return response.data || []
    } catch (error) {
      console.error('Error fetching sessions:', error)
      return []
    }
  },

  async getSession(appName: string, sessionId: string, userId: string = 'user'): Promise<BackendSession | null> {
    try {
      const response = await apiClient.getSession(appName, userId, sessionId)
      return response.data || null
    } catch (error) {
      console.error('Error fetching session:', error)
      return null
    }
  },

  async createSession(appName: string, userId: string = 'user'): Promise<BackendSession | null> {
    try {
      const response = await apiClient.createSession(appName, userId)
      return response.data || null
    } catch (error) {
      console.error('Error creating session:', error)
      return null
    }
  },

  convertBackendEventToMessage(event: BackendEvent, agentId: string): Message {
    // Extract content from event
    let content = ''
    let parts = undefined

    if (event.content?.parts && Array.isArray(event.content.parts)) {
      const textParts = event.content.parts.filter(part => part.text)
      if (textParts.length > 0) {
        content = textParts.map(part => part.text).join('')
      }
      parts = event.content.parts
    }

    // Determine message type
    let messageType: 'user' | 'agent' | 'inter_agent' | 'system' = 'agent'
    if (event.object === 'tool_call' || event.object === 'tool_response') {
      messageType = 'system'
    } else if (event.author === 'user') {
      messageType = 'user'
    }

    return {
      id: event.id || event.invocationId || Date.now().toString(),
      content,
      timestamp: new Date(event.timestamp * 1000),
      sender: event.author || agentId,
      type: messageType,
      metadata: {
        invocationId: event.invocationId,
        partial: event.partial,
        done: event.done
      },
      parts
    }
  },

  convertBackendSessionToMessages(session: BackendSession, agentId: string): Message[] {
    if (!session.events) return []
    
    return session.events.map(event => 
      this.convertBackendEventToMessage(event, agentId)
    )
  }
}