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
        // Handle connection status changes gracefully
        if (!connected) {
          console.log('Agent SSE connection closed normally')
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
        // Extract text content from parts
        const textParts = event.content.parts.filter(part => part.text)
        content = textParts.map(part => part.text).join('')
        
        // Convert parts to MessagePart format
        parts = event.content.parts.map(part => {
          const messagePart: MessagePart = {}
          
          if (part.text) {
            messagePart.text = part.text
          }
          
          if (part.functionCall) {
            messagePart.functionCall = {
              name: part.functionCall.name,
              args: part.functionCall.args,
              id: part.functionCall.id
            }
          }
          
          if (part.functionResponse) {
            messagePart.functionResponse = {
              name: part.functionResponse.name,
              response: part.functionResponse.response,
              id: part.functionResponse.id
            }
          }
          
          return messagePart
        })
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

    // Determine message type based on event properties
    let messageType: 'user' | 'agent' | 'inter_agent' | 'system' = 'system'
    
    if (event.type === 'inter_agent' || event.type === 'communication') {
      messageType = 'inter_agent'
    } else if (event.object === 'tool_call' || event.object === 'tool_response') {
      messageType = 'system'
    } else if (event.author && event.author !== 'system') {
      messageType = 'agent'
    }

    return {
      id: event.id || event.invocationId || Date.now().toString(),
      content,
      timestamp,
      sender: event.fromAgent || event.author || 'system',
      type: messageType,
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