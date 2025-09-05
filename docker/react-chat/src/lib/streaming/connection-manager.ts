import { AgentEvent, AgentRunRequest } from '@/lib/types'
import { 
  StreamingConnection, 
  ConnectionStatus, 
  StreamingHandlers, 
  StreamingMessageManagerConfig
} from '@/lib/types/streaming'
import { apiClient } from '@/lib/api'
import { debug } from '@/lib/utils/debug'

export class ConnectionManager {
  private connections = new Map<string, StreamingConnection>()
  private config: StreamingMessageManagerConfig
  
  constructor(config: StreamingMessageManagerConfig) {
    this.config = config
  }

  createAgentConnection(
    agentId: string, 
    sessionId: string, 
    handlers: StreamingHandlers
  ): StreamingConnection {
    const connectionId = `agent-${agentId}-${sessionId}`
    
    // Check if existing connection exists and is still valid
    const existingConnection = this.connections.get(connectionId)
    if (existingConnection && existingConnection.status.isConnected) {
      debug.connection(`Reusing existing agent connection: ${connectionId}`)
      // Update handlers for the existing connection
      existingConnection.handlers = handlers
      return existingConnection
    }
    
    // Only close existing connection if it's not connected
    if (existingConnection && existingConnection.status.isConnected) {
      debug.connection(`Not closing active connection ${connectionId}, reusing it`)
      existingConnection.handlers = handlers
      return existingConnection
    }
    
    // Close existing connection if any and it's not connected
    if (existingConnection) {
      this.closeConnection(connectionId)
    }
    
    const connection: StreamingConnection = {
      id: connectionId,
      type: 'agent_run',
      agentId,
      sessionId,
      eventSource: null,
      status: {
        connectionId,
        isConnected: false,
        errorCount: 0,
        reconnectAttempts: 0
      },
      handlers
    }

    this.connections.set(connectionId, connection)
    return connection
  }

  createInterAgentConnection(
    agents: string[], 
    sessionId: string, 
    userId: string, 
    handlers: StreamingHandlers
  ): StreamingConnection {
    const connectionId = `inter-agent-${agents.join('-')}-${sessionId}`
    
    // Check if existing connection exists and is still valid
    const existingConnection = this.connections.get(connectionId)
    if (existingConnection && existingConnection.status.isConnected) {
      debug.connection(`Reusing existing inter-agent connection: ${connectionId}`)
      // Update handlers for the existing connection
      existingConnection.handlers = handlers
      return existingConnection
    }
    
    // Only close existing connection if it's not connected
    if (existingConnection && existingConnection.status.isConnected) {
      debug.connection(`Not closing active inter-agent connection ${connectionId}, reusing it`)
      existingConnection.handlers = handlers
      return existingConnection
    }
    
    // Close existing connection if any and it's not connected
    if (existingConnection) {
      this.closeConnection(connectionId)
    }
    
    const connection: StreamingConnection = {
      id: connectionId,
      type: 'inter_agent',
      sessionId,
      eventSource: null,
      status: {
        connectionId,
        isConnected: false,
        errorCount: 0,
        reconnectAttempts: 0
      },
      handlers
    }

    this.connections.set(connectionId, connection)
    this.establishInterAgentConnection(connection, agents, sessionId, userId)
    return connection
  }

  async establishAgentRunConnection(
    connection: StreamingConnection, 
    request: { agentId: string; sessionId: string; userId: string; content: string }
  ): Promise<void> {
    if (connection.eventSource) {
      connection.eventSource.close()
    }

    const { agentId, sessionId, userId, content } = request
    
    debug.critical(`ConnectionManager: Establishing agent run connection for ${agentId}`, {
      content: content.substring(0, 50) + '...',
      sessionId,
      userId
    })
    
    try {
      const url = apiClient.getRunSSEUrl()
      const requestBody: AgentRunRequest = {
        appName: agentId,
        userID: userId,
        sessionID: sessionId,
        streaming: true,
        newMessage: {
          role: 'user',
          parts: [{ text: content }]
        }
      }

      debug.critical(`ConnectionManager: Sending SSE request to ${url}`, {
        requestBody: {
          ...requestBody,
          newMessage: {
            role: requestBody.newMessage.role,
            parts: [{ text: requestBody.newMessage.parts[0].text.substring(0, 50) + '...' }]
          }
        }
      })

      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'text/event-stream',
          'Cache-Control': 'no-cache',
          'Connection': 'keep-alive',
        },
        body: JSON.stringify(requestBody),
      })

      debug.critical(`ConnectionManager: SSE Response received`, {
        status: response.status,
        statusText: response.statusText,
        headers: Object.fromEntries(response.headers.entries())
      })

      if (!response.ok) {
        let errorDetails = '';
        try {
          const errorText = await response.text();
          errorDetails = errorText ? ` - ${errorText}` : '';
        } catch (e) {
          // Ignore error text parsing issues
        }
        
        console.error('Agent run request failed:', {
          status: response.status,
          statusText: response.statusText,
          url,
          requestBody,
          errorDetails
        });
        
        debug.error(`ConnectionManager: Agent run request failed for ${agentId}`, {
          status: response.status,
          statusText: response.statusText,
          url,
          errorDetails
        });
        
        throw new Error(`Agent run request failed: ${response.status} ${response.statusText}${errorDetails}`);
      }

      if (!response.body) {
        debug.error(`ConnectionManager: No response body for ${agentId}`)
        throw new Error('No response body')
      }

      debug.critical(`ConnectionManager: Starting to handle streaming response for ${agentId}`)
      await this.handleStreamingResponse(connection, response)
      
      debug.critical(`ConnectionManager: Completed handling streaming response for ${agentId}`)
    } catch (error) {
      console.error(`Error establishing agent run connection for ${agentId}:`, error)
      debug.error(`ConnectionManager: Error in establishAgentRunConnection for ${agentId}:`, error)
      this.handleConnectionError(connection, error as Error)
    }
  }

  private establishInterAgentConnection(
    connection: StreamingConnection, 
    agents: string[], 
    sessionId: string, 
    userId: string
  ): void {
    try {
      const url = apiClient.getSSEUrl(agents, sessionId, userId)
      
      const eventSource = new EventSource(url)
      connection.eventSource = eventSource

      eventSource.onopen = () => {
        this.updateConnectionStatus(connection, { 
          isConnected: true, 
          lastConnected: new Date(),
          errorCount: 0,
          reconnectAttempts: 0
        })
        connection.handlers.onConnectionChange?.(connection.status)
      }

      eventSource.onmessage = (event) => {
        try {
          const data: AgentEvent = JSON.parse(event.data)
          connection.handlers.onMessage?.(data)
        } catch (error) {
          console.error('Error parsing inter-agent SSE data:', error, 'Raw data:', event.data)
        }
      }

      eventSource.onerror = () => {
        this.handleConnectionError(connection, new Error('EventSource error'))
      }

    } catch (error) {
      this.handleConnectionError(connection, error as Error)
    }
  }

  private async handleStreamingResponse(
    connection: StreamingConnection, 
    response: Response
  ): Promise<void> {
    debug.critical(`ConnectionManager: Starting handleStreamingResponse for ${connection.id}`)
    
    if (!response.body) {
      debug.error(`ConnectionManager: No response body in handleStreamingResponse for ${connection.id}`)
      return
    }

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let eventCount = 0

    this.updateConnectionStatus(connection, { 
      isConnected: true, 
      lastConnected: new Date(),
      errorCount: 0,
      reconnectAttempts: 0
    })
    connection.handlers.onConnectionChange?.(connection.status)
    
    debug.critical(`ConnectionManager: Connection established successfully for ${connection.id}`)

    try {
      while (true) {
        const { done, value } = await reader.read()
        debug.streaming(`ConnectionManager: Read chunk for ${connection.id}`, {
          done,
          valueLength: value?.length || 0
        })
        
        if (done) {
          debug.critical(`ConnectionManager: Stream ended for ${connection.id} after ${eventCount} events`)
          break
        }

        const chunk = decoder.decode(value, { stream: true })
        debug.streaming(`ConnectionManager: Decoded chunk for ${connection.id}`, {
          chunkLength: chunk.length,
          chunkPreview: chunk.substring(0, 100)
        })
        
        const lines = chunk.split('\n')

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            try {
              const data: AgentEvent = JSON.parse(line.slice(6))
              eventCount++
              
              debug.streaming(`ConnectionManager: Processing SSE event #${eventCount} for ${connection.id}`, {
                type: data.type,
                object: data.object,
                hasContent: !!data.content,
                // contentPreview: data.content ? String(data.content).substring(0, 100) : 'null',  // Remove verbose content preview
                invocationId: data.invocationId,
                done: data.done,
                partial: data.partial
              })
              
              // Always forward the event to handlers, even if it seems empty
              connection.handlers.onMessage?.(data)
            } catch (error) {
              console.error('Error parsing agent run SSE data:', error, 'Raw data:', line)
              debug.error(`ConnectionManager: Error parsing SSE data for ${connection.id}:`, {
                error: error instanceof Error ? error.message : String(error),
                rawLine: line
              })
            }
          } else if (line.trim() === '' || line.startsWith(':')) {
            // Skip empty lines and comments
            continue
          } else if (line.trim() !== '' && !line.startsWith('event:') && !line.startsWith('id:')) {
            debug.warn('Unexpected SSE line format:', line)
          }
        }
      }
    } catch (error) {
      console.error(`Stream read error for ${connection.id}:`, error)
      debug.error(`ConnectionManager: Stream read error for ${connection.id}:`, error)
      this.handleConnectionError(connection, error as Error)
    } finally {
      debug.critical(`ConnectionManager: Closing connection ${connection.id} (processed ${eventCount} events)`)
      this.updateConnectionStatus(connection, { isConnected: false })
      connection.handlers.onConnectionChange?.(connection.status)
    }
  }

  private handleConnectionError(connection: StreamingConnection, error: Error): void {
    console.error(`Connection error for ${connection.id}:`, error)
    
    this.updateConnectionStatus(connection, { 
      isConnected: false, 
      errorCount: connection.status.errorCount + 1 
    })
    
    connection.handlers.onError?.(error)
    connection.handlers.onConnectionChange?.(connection.status)

    // Attempt reconnection for inter-agent connections
    if (connection.type === 'inter_agent' && 
        connection.status.reconnectAttempts < this.config.maxReconnectAttempts) {
      this.scheduleReconnection(connection)
    }
  }

  private scheduleReconnection(connection: StreamingConnection): void {
    const delay = Math.min(
      this.config.reconnectInterval * Math.pow(this.config.backoffMultiplier, connection.status.reconnectAttempts),
      30000 // Max 30 seconds
    )

    setTimeout(() => {
      if (this.connections.has(connection.id) && !connection.status.isConnected) {
        this.updateConnectionStatus(connection, { 
          reconnectAttempts: connection.status.reconnectAttempts + 1 
        })
        
        // Reconstruct the connection based on type
        // This is a simplified approach - in production you'd store the original parameters
        console.log(`Attempting to reconnect ${connection.id}...`)
      }
    }, delay)
  }

  private updateConnectionStatus(connection: StreamingConnection, updates: Partial<ConnectionStatus>): void {
    connection.status = {
      ...connection.status,
      ...updates
    }
  }

  closeConnection(connectionId: string): void {
    const connection = this.connections.get(connectionId)
    if (connection) {
      if (connection.eventSource) {
        connection.eventSource.close()
        connection.eventSource = null
      }
      
      this.updateConnectionStatus(connection, { isConnected: false })
      this.connections.delete(connectionId)
      
      debug.connection(`Connection ${connectionId} closed`)
    }
  }

  closeAllConnections(): void {
    for (const connectionId of this.connections.keys()) {
      this.closeConnection(connectionId)
    }
  }

  getConnection(connectionId: string): StreamingConnection | undefined {
    return this.connections.get(connectionId)
  }

  getAllConnections(): StreamingConnection[] {
    return Array.from(this.connections.values())
  }

  getConnectionStatus(connectionId: string): ConnectionStatus | undefined {
    return this.connections.get(connectionId)?.status
  }
}