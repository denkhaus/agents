import { AgentEvent, Message, MultiChatRequest } from '@/lib/types'
import { 
  StreamingConnection, 
  ConnectionStatus,
  MessageCallback, 
  InterAgentCallback, 
  ConnectionCallback, 
  ErrorCallback,
  StreamingMessageManagerConfig,
  SendMessageOptions,
  MessageProcessingContext
} from '@/lib/types/streaming'
import { ConnectionManager } from './connection-manager'
import { MessageEventRouter } from './message-router'
import { debug } from '@/lib/utils/debug'

export class StreamingMessageManager {
  private connectionManager: ConnectionManager
  private messageRouter: MessageEventRouter
  private messageCallbacks: Set<MessageCallback> = new Set()
  private interAgentCallbacks: Set<InterAgentCallback> = new Set()
  private connectionCallbacks: Set<ConnectionCallback> = new Set()
  private errorCallbacks: Set<ErrorCallback> = new Set()

  private static instance: StreamingMessageManager | null = null

  constructor(config?: Partial<StreamingMessageManagerConfig>) {
    const defaultConfig: StreamingMessageManagerConfig = {
      maxReconnectAttempts: 5,
      reconnectInterval: 1000,
      backoffMultiplier: 2,
      connectionTimeout: 30000,
      ...config
    }

    this.connectionManager = new ConnectionManager(defaultConfig)
    this.messageRouter = new MessageEventRouter()
  }

  static getInstance(config?: Partial<StreamingMessageManagerConfig>): StreamingMessageManager {
    if (!StreamingMessageManager.instance) {
      StreamingMessageManager.instance = new StreamingMessageManager(config)
    }
    return StreamingMessageManager.instance
  }

  // Connection Management
  establishAgentConnection(agentId: string, sessionId: string): StreamingConnection {
    debug.connection(`Establishing agent connection for ${agentId}, session: ${sessionId}`)
    
    const connection = this.connectionManager.createAgentConnection(
      agentId, 
      sessionId, 
      {
        onMessage: (event: AgentEvent) => {
          debug.streaming(`Received event in StreamingMessageManager for agent ${agentId}:`, {
            type: event.type,
            object: event.object,
            content: event.content ? this.formatContentPreview(event.content) : 'null'
          })
          this.handleMessage(event, agentId, sessionId, 'agent_run')
        },
        onConnectionChange: (status) => {
          debug.connection(`Connection status change for ${agentId}:`, status)
          this.notifyConnectionChange(status)
        },
        onError: (error) => {
          debug.error(`Connection error for ${agentId}:`, error)
          this.notifyError(error)
        }
      }
    )

    debug.connection(`Established agent connection: ${connection.id}`)
    return connection
  }

  establishInterAgentConnection(
    agents: string[], 
    sessionId: string, 
    userId: string = 'user'
  ): StreamingConnection {
    const connection = this.connectionManager.createInterAgentConnection(
      agents, 
      sessionId, 
      userId, 
      {
        onMessage: (event: AgentEvent) => this.handleInterAgentMessage(event, agents[0], sessionId),
        onInterAgentEvent: (event: AgentEvent) => this.notifyInterAgentEvent(event),
        onConnectionChange: (status) => this.notifyConnectionChange(status),
        onError: (error) => this.notifyError(error)
      }
    )

    debug.connection(`Established inter-agent connection: ${connection.id}`)
    return connection
  }

  closeConnection(connectionId: string): void {
    this.connectionManager.closeConnection(connectionId)
  }

  closeAllConnections(): void {
    this.connectionManager.closeAllConnections()
  }

  // Message Sending
  async sendUserMessage(
    agentId: string, 
    content: string, 
    options?: SendMessageOptions
  ): Promise<void> {
    if (!content.trim()) {
      throw new Error('Message content cannot be empty')
    }

    const sessionId = options?.sessionId || `session-${Date.now()}`
    const userId = options?.userId || 'user'

    debug.critical(`StreamingMessageManager: sendUserMessage called for agent ${agentId}`, {
      content: content.substring(0, 50) + '...',
      sessionId,
      userId,
      options
    })

    try {
      // Establish connection for this specific message
      const connection = this.establishAgentConnection(agentId, sessionId)
      
      debug.critical(`StreamingMessageManager: Established agent connection for ${agentId}`, {
        connectionId: connection.id,
        sessionId: connection.sessionId
      })
      
      // Send the message through the connection
      debug.critical(`StreamingMessageManager: Calling establishAgentRunConnection for ${agentId}`)
      await this.connectionManager.establishAgentRunConnection(connection, {
        agentId,
        sessionId,
        userId,
        content
      })
      
      debug.critical(`StreamingMessageManager: Completed establishAgentRunConnection for ${agentId}`)

    } catch (error) {
      console.error('Error sending user message:', error)
      debug.error(`StreamingMessageManager: Error in sendUserMessage for agent ${agentId}:`, error)
      options?.onError?.(error as Error)
      throw error
    }
  }

  async sendInterAgentMessage(request: MultiChatRequest): Promise<void> {
    // This would be implemented based on your inter-agent communication API
    // For now, we'll establish the connection to listen for responses
    // Convert single-agent request to multi-agent by using fromAgent and toAgent
    const agents = [request.fromAgent, request.toAgent].filter(Boolean)
    this.establishInterAgentConnection(agents, request.sessionId, request.userId)
  }

  // Event Handlers
  private handleMessage(
    event: AgentEvent, 
    agentId: string, 
    sessionId: string, 
    connectionType: StreamingConnection['type']
  ): void {
    debug.streaming(`StreamingMessageManager: Processing event for agent ${agentId}`, {
      eventType: event.type,
      eventObject: event.object,
      hasContent: !!event.content,
      invocationId: event.invocationId
    })

    const context: MessageProcessingContext = {
      agentId,
      sessionId,
      connectionType,
      timestamp: new Date()
    }

    const message = this.messageRouter.processEvent(event, context)
    
    if (message) {
      debug.streaming(`StreamingMessageManager: Message created successfully`, {
        messageId: message.id,
        messageType: message.type,
        hasContent: !!message.content,
        isPartial: message.metadata?.partial
      })
      this.notifyMessage(message)
    } else {
      debug.warn(`StreamingMessageManager: No message generated from event`, {
        eventType: event.type,
        eventObject: event.object
      })
    }
  }

  private handleInterAgentMessage(
    event: AgentEvent, 
    primaryAgentId: string, 
    sessionId: string
  ): void {
    // Handle inter-agent messages - route to appropriate agent based on fromAgent/toAgent
    const targetAgentId = event.fromAgent || event.toAgent || primaryAgentId
    
    const context: MessageProcessingContext = {
      agentId: targetAgentId,
      sessionId,
      connectionType: 'inter_agent',
      timestamp: new Date()
    }

    const message = this.messageRouter.processEvent(event, context)
    
    if (message) {
      this.notifyMessage(message)
    }

    // Also notify inter-agent event listeners
    this.notifyInterAgentEvent(event)
  }

  // Callback Management
  onMessage(callback: MessageCallback): () => void {
    this.messageCallbacks.add(callback)
    return () => this.messageCallbacks.delete(callback)
  }

  onInterAgentEvent(callback: InterAgentCallback): () => void {
    this.interAgentCallbacks.add(callback)
    return () => this.interAgentCallbacks.delete(callback)
  }

  onConnectionChange(callback: ConnectionCallback): () => void {
    this.connectionCallbacks.add(callback)
    return () => this.connectionCallbacks.delete(callback)
  }

  onError(callback: ErrorCallback): () => void {
    this.errorCallbacks.add(callback)
    return () => this.errorCallbacks.delete(callback)
  }

  // Notification Methods
  private notifyMessage(message: Message): void {
    for (const callback of this.messageCallbacks) {
      try {
        callback(message)
      } catch (error) {
        console.error('Error in message callback:', error)
      }
    }
  }

  private notifyInterAgentEvent(event: AgentEvent): void {
    for (const callback of this.interAgentCallbacks) {
      try {
        callback(event)
      } catch (error) {
        console.error('Error in inter-agent event callback:', error)
      }
    }
  }

  private notifyConnectionChange(status: ConnectionStatus): void {
    for (const callback of this.connectionCallbacks) {
      try {
        callback(status)
      } catch (error) {
        console.error('Error in connection change callback:', error)
      }
    }
  }

  private notifyError(error: Error): void {
    for (const callback of this.errorCallbacks) {
      try {
        callback(error)
      } catch (error) {
        console.error('Error in error callback:', error)
      }
    }
  }

  // Utility Methods
  getConnectionStatus(): Record<string, ConnectionStatus> {
    const connections = this.connectionManager.getAllConnections()
    const status: Record<string, ConnectionStatus> = {}
    
    for (const connection of connections) {
      status[connection.id] = connection.status
    }
    
    return status
  }

  isConnected(connectionId?: string): boolean {
    if (connectionId) {
      const status = this.connectionManager.getConnectionStatus(connectionId)
      return status?.isConnected || false
    }
    
    // Check if any connection is active
    const connections = this.connectionManager.getAllConnections()
    return connections.some(conn => conn.status.isConnected)
  }

  // Cleanup
  destroy(): void {
    this.closeAllConnections()
    this.messageCallbacks.clear()
    this.interAgentCallbacks.clear()
    this.connectionCallbacks.clear()
    this.errorCallbacks.clear()
    StreamingMessageManager.instance = null
  }

  // Add this helper method to safely format content preview
  private formatContentPreview(content: unknown): string {
    try {
      // Handle different content structures
      if (typeof content === 'string') {
        return content.substring(0, 100) + '...';
      }
      
      // Handle object content with parts array
      if (content && typeof content === 'object') {
        if ('parts' in content && Array.isArray((content as { parts: unknown[] }).parts)) {
          // Extract text from parts
          const textParts = (content as { parts: Array<{ text?: string }> }).parts
            .filter((part: { text?: string }) => part && part.text)
            .map((part: { text?: string }) => part.text)
            .join('');
          return textParts.substring(0, 100) + '...';
        }
        
        // Handle other object structures
        return JSON.stringify(content).substring(0, 100) + '...';
      }
      
      // Fallback for other types
      return String(content).substring(0, 100) + '...';
    } catch (error) {
      return '[Error formatting content preview]';
    }
  }
}