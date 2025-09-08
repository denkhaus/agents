import { StreamingMessageManager } from '@/lib/streaming'
import { MessageEventRouter } from '@/lib/streaming/message-router'
import { LLMEvent } from '@/lib/types'
import { AgentId } from '@/lib/constants/agents'

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
      const connection = manager.establishAgentConnection("agent-1" as AgentId, "session-1")
      
      expect(connection).toBeDefined()
      expect(connection.id).toBe('agent-agent-1-session-1')
      expect(connection.type).toBe('agent_run')
      expect(connection.agentId).toBe('agent-1' as AgentId)
      expect(connection.sessionId).toBe('session-1')
    })

    test('should establish inter-agent connection', () => {
      const connection = manager.establishAgentConnection("agent-1" as AgentId, "session-1")
      
      expect(connection).toBeDefined()
      expect(connection.id).toBe('inter-agent-agent-1-agent-2-session-1')
      expect(connection.type).toBe('inter_agent')
      expect(connection.sessionId).toBe('session-1')
    })

    test('should close connection', () => {
      const connection = manager.establishAgentConnection("agent-1" as AgentId, "session-1")
      const connectionId = connection.id
      
      manager.closeConnection(connectionId)
      
      expect(manager.isConnected(connectionId)).toBe(false)
    })
  })

  describe('Message Callbacks', () => {
    test('should register and unregister message callbacks', () => {
      const callback = jest.fn()
      const unsubscribe = manager.onMessage(callback)

      // @ts-ignore
      manager.notifyMessage({ id: '1', content: 'test', timestamp: new Date(), sender: 'agent-1', type: 'agent' })
      expect(callback).toHaveBeenCalledTimes(1)

      unsubscribe()

      // @ts-ignore
      manager.notifyMessage({ id: '2', content: 'test2', timestamp: new Date(), sender: 'agent-1', type: 'agent' })
      expect(callback).toHaveBeenCalledTimes(1)
    })

    test('should register inter-agent event callbacks', () => {
      const callback = jest.fn()
      const unsubscribe = manager.onInterAgentEvent(callback)

      // @ts-ignore
      manager.notifyInterAgentEvent({ id: '1', type: 'inter_agent', parts: [], timestamp: 0 })
      expect(callback).toHaveBeenCalledTimes(1)

      unsubscribe()

      // @ts-ignore
      manager.notifyInterAgentEvent({ id: '2', type: 'inter_agent', parts: [], timestamp: 0 })
      expect(callback).toHaveBeenCalledTimes(1)
    })
  })

  describe('Connection Status', () => {
    test('should return connection status', () => {
      manager.establishAgentConnection("agent-1" as AgentId, "session-1")
      
      const status = manager.getConnectionStatus()
      expect(status).toBeDefined()
      expect(typeof status).toBe('object')
    })

    test('should check if connected', () => {
      const connection = manager.establishAgentConnection("agent-1" as AgentId, "session-1")
      
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
      const event: LLMEvent = {
        id: 'msg-1',
        parts: [{ content: "Test content" }],
        type: 'assistant',
  
        timestamp: Date.now() / 1000,
        author: 'agent-1' as AgentId
      }

      const context = {
        agentId: 'agent-1' as AgentId,
        sessionId: 'session-1',
        connectionType: 'agent_run' as const,
        timestamp: new Date()
      }

      const result = router.processEvent(event, context)
      
      expect(result).toBeDefined()
      expect(result?.content).toBe('Hello world')
      expect(result?.type).toBe('agent')
      expect(result?.sender).toBe('agent-1' as AgentId)
    })

    test('should process inter-agent events', () => {
      const event: LLMEvent = {
        id: 'msg-2',
        parts: [{ content: "Test content" }],
        type: 'inter_agent',
        inter_agent: { from_agent: 'agent-1', to_agent: 'agent-2', type: 'communication' },
        timestamp: Date.now() / 1000
      }

      const context = {
        agentId: 'agent-1' as AgentId,
        sessionId: 'session-1',
        connectionType: 'inter_agent' as const,
        timestamp: new Date()
      }

      const result = router.processEvent(event, context)
      
      expect(result).toBeDefined()
      expect(result?.content).toBe('Inter-agent message')
      expect(result?.type).toBe('inter_agent')
      expect(result?.metadata?.fromAgent).toBe('agent-1' as AgentId)
      expect(result?.metadata?.toAgent).toBe('agent-2')
    })

    test('should process tool call events', () => {
      const event: LLMEvent = {
        id: 'tool-1',
        type: 'tool.call',
        parts: [{ content: "Test content" }],
        timestamp: Date.now() / 1000
      }

      const context = {
        agentId: 'agent-1' as AgentId,
        sessionId: 'session-1',
        connectionType: 'agent_run' as const,
        timestamp: new Date()
      }

      const result = router.processEvent(event, context)
      
      expect(result).toBeDefined()
      expect(result?.type).toBe('system')
      expect(result?.metadata?.object).toBe('tool_call')
    })

    test('should handle events without content', () => {
      const event: LLMEvent = {
        id: 'empty-1',
        type: 'assistant',
        timestamp: Date.now() / 1000
      }

      const context = {
        agentId: 'agent-1' as AgentId,
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