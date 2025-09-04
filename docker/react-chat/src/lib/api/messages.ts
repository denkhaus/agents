import { apiClient } from './client'
import { sseService } from './sse'
import { AgentRunRequest, Message, MessagePart, InterAgentEvent, MultiChatRequest } from '@/lib/types'

export const messageApi = {
  async sendMessage(
    agentId: string, 
    content: string, 
    sessionId: string, 
    userId: string = 'user',
    onResponse?: (event: InterAgentEvent) => void,
    onError?: (error: any) => void
  ): Promise<void> {
    const request: AgentRunRequest = {
      appName: agentId,
      userID: userId,
      sessionID: sessionId,
      streaming: true,
      newMessage: {
        role: 'user',
        parts: [{ text: content }]
      }
    }
    
    // Use SSE service for streaming responses
    sseService.connectAgentRun(request, {
      onMessage: (event) => {
        onResponse?.(event)
      },
      onError: (error) => {
        console.error('SSE error during agent communication:', error)
        onError?.(error)
      },
      onConnectionStatusChange: (connected) => {
        // Only log connection failures, not status changes
        if (!connected) {
          console.error('Agent SSE connection lost')
        }
      }
    })
  },

  async sendMultiAgentMessage(request: MultiChatRequest) {
    return apiClient.sendMessage(request)
  },

  convertEventToMessage(event: InterAgentEvent): Message {
    
    // Safe content extraction with null checks
    let content = ''
    let parts: MessagePart[] | undefined = undefined
    
    if (typeof event.content === 'string') {
      content = event.content
    } else if (event.content && typeof event.content === 'object') {
      // Check if content has parts array
      if (Array.isArray(event.content.parts) && event.content.parts.length > 0) {
        content = event.content.parts[0]?.text || ''
        parts = event.content.parts as MessagePart[]
      } else {
        // Fallback: try to extract text from content object
        content = JSON.stringify(event.content)
      }
    }
    
    // Handle timestamp - could be Unix timestamp or already a Date
    let timestamp: Date
    if (typeof event.timestamp === 'number') {
      // Unix timestamp (seconds or milliseconds)
      timestamp = event.timestamp > 1000000000000 
        ? new Date(event.timestamp) 
        : new Date(event.timestamp * 1000)
    } else {
      timestamp = new Date()
    }

    return {
      id: event.id || event.invocationId || Date.now().toString(),
      content,
      timestamp,
      sender: event.fromAgent || event.author || 'system',
      type: event.type === 'inter_agent' ? 'inter_agent' : 
            event.type === 'communication' ? 'inter_agent' : 'system',
      metadata: {
        fromAgent: event.fromAgent,
        toAgent: event.toAgent,
        eventType: event.type,
        partial: event.partial,
        done: event.done
      },
      parts,
      usageMetadata: event.usageMetadata
    }
  }
}