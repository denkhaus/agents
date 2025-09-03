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
    console.log('messageApi.sendMessage:', { agentId, content, sessionId, userId })
    
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

    console.log('Connecting to agent via SSE for streaming response:', request)
    
    // Use SSE service for streaming responses
    sseService.connectForAgentRun(request, {
      onMessage: (event) => {
        console.log('Received agent response:', event)
        onResponse?.(event)
      },
      onError: (error) => {
        console.error('SSE error during agent communication:', error)
        onError?.(error)
      },
      onConnectionStatusChange: (connected) => {
        console.log('Agent SSE connection status:', connected)
      }
    })
  },

  async sendMultiAgentMessage(request: MultiChatRequest) {
    return apiClient.sendMessage(request)
  },

  convertEventToMessage(event: InterAgentEvent): Message {
    const content = typeof event.content === 'string' 
      ? event.content 
      : event.content.parts?.[0]?.text || ''

    return {
      id: event.id || event.invocationId || Date.now().toString(),
      content,
      timestamp: new Date(event.timestamp * 1000),
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
      parts: typeof event.content === 'object' ? event.content.parts as MessagePart[] : undefined,
      usageMetadata: event.usageMetadata
    }
  }
}