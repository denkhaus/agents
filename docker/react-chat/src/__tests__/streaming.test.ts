import { StreamingMessageManager } from '@/lib/streaming'
import { MessageEventRouter } from '@/lib/streaming/message-router'
import { AgentEvent } from '@/lib/types'

// Mock the API client
jest.mock('@/lib/api', () => ({
  apiClient: {
    getRunSSEUrl: () => 'http://localhost:3001/run_sse',
    getSSEUrl: () => 'http://localhost:3001/sse'
  }
}))

// Mock fetch globally
global.fetch = jest.fn()

describe('StreamingMessageManager', () => {
  let manager: StreamingMessageManager
  
  beforeEach(() => {
    manager = new StreamingMessageManager({
      maxReconnectAttempts: 3,
      reconnectInterval: 1000,
      backoffMultiplier: 2,
      connectionTimeout: 5000
    })
    
    // Reset fetch mock
    ;(global.fetch as jest.Mock).mockReset()
  })

  afterEach(() => {
    manager.destroy()
  })

  describe('Connection Management', () => {
    test('should establish agent connection', () => {
      const connection = manager.establishAgentConnection('agent-1', 'session-1')
      
      expect(connection).toBeDefined()
      expect(connection.id).toBe('agent-agent-1-session-1')
      expect(connection.type).toBe('agent_run')
      expect(connection.agentId).toBe('agent-1')
      expect(connection.sessionId).toBe('session-1')
    })

    test('should establish inter-agent connection', () => {
      const connection = manager.establishInterAgentConnection(['agent-1', 'agent-2'], 'session-1')
      
      expect(connection).toBeDefined()
      expect(connection.id).toBe('inter-agent-agent-1-agent-2-session-1')
      expect(connection.type).toBe('inter_agent')
      expect(connection.sessionId).toBe('session-1')
    })

    test('should close connection', () => {
      const connection = manager.establishAgentConnection('agent-1', 'session-1')
      const connectionId = connection.id
      
      manager.closeConnection(connectionId)
      
      expect(manager.isConnected(connectionId)).toBe(false)
    })
  })

  describe('Message Callbacks', () => {
    test('should register and unregister message callbacks', () => {
      const callback = jest.fn()
      const unsubscribe = manager.onMessage(callback)
      
      expect(typeof unsubscribe).toBe('function')
      
      unsubscribe()
      // After unsubscribing, callback should not be called
    })

    test('should register inter-agent event callbacks', () => {
      const callback = jest.fn()
      const unsubscribe = manager.onInterAgentEvent(callback)
      
      expect(typeof unsubscribe).toBe('function')
      
      unsubscribe()
    })
  })

  describe('Connection Status', () => {
    test('should return connection status', () => {
      manager.establishAgentConnection('agent-1', 'session-1')
      
      const status = manager.getConnectionStatus()
      expect(status).toBeDefined()
      expect(typeof status).toBe('object')
    })

    test('should check if connected', () => {
      const connection = manager.establishAgentConnection('agent-1', 'session-1')
      
      // Initially not connected until EventSource is established
      expect(manager.isConnected(connection.id)).toBe(false)
    })
  })
})

describe('MessageEventRouter', () => {
  let router: MessageEventRouter

  beforeEach(() => {
    router = new MessageEventRouter()
  })

  describe('Event Processing', () => {
    test('should process user response events', () => {
      const event: AgentEvent = {
        id: 'msg-1',
        content: 'Hello world',
        object: 'message',
        type: 'user_response',
        timestamp: Date.now() / 1000,
        author: 'agent-1'
      }

      const context = {
        agentId: 'agent-1',
        sessionId: 'session-1',
        connectionType: 'agent_run' as const,
        timestamp: new Date()
      }

      const result = router.processEvent(event, context)
      
      expect(result).toBeDefined()
      expect(result?.content).toBe('Hello world')
      expect(result?.type).toBe('agent')
      expect(result?.sender).toBe('agent-1')
    })

    test('should process inter-agent events', () => {
      const event: AgentEvent = {
        id: 'msg-2',
        content: 'Inter-agent message',
        type: 'inter_agent',
        fromAgent: 'agent-1',
        toAgent: 'agent-2',
        timestamp: Date.now() / 1000
      }

      const context = {
        agentId: 'agent-1',
        sessionId: 'session-1',
        connectionType: 'inter_agent' as const,
        timestamp: new Date()
      }

      const result = router.processEvent(event, context)
      
      expect(result).toBeDefined()
      expect(result?.content).toBe('Inter-agent message')
      expect(result?.type).toBe('inter_agent')
      expect(result?.metadata?.fromAgent).toBe('agent-1')
      expect(result?.metadata?.toAgent).toBe('agent-2')
    })

    test('should process tool call events', () => {
      const event: AgentEvent = {
        id: 'tool-1',
        object: 'tool_call',
        content: 'Tool execution',
        timestamp: Date.now() / 1000
      }

      const context = {
        agentId: 'agent-1',
        sessionId: 'session-1',
        connectionType: 'agent_run' as const,
        timestamp: new Date()
      }

      const result = router.processEvent(event, context)
      
      expect(result).toBeDefined()
      expect(result?.type).toBe('system')
      expect(result?.metadata?.toolCallType).toBe('tool_call')
    })

    test('should handle events without content', () => {
      const event: AgentEvent = {
        id: 'empty-1',
        type: 'system',
        timestamp: Date.now() / 1000
      }

      const context = {
        agentId: 'agent-1',
        sessionId: 'session-1',
        connectionType: 'agent_run' as const,
        timestamp: new Date()
      }

      const result = router.processEvent(event, context)
      
      // System events without content should return null (no message to display)
      expect(result).toBeNull()
    })
  })
})