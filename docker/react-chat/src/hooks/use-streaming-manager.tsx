"use client"

import { useEffect, useRef, useCallback } from 'react'
import { StreamingMessageManager } from '@/lib/streaming'
import { 
  MessageCallback, 
  InterAgentCallback, 
  ConnectionCallback,
  ErrorCallback
} from '@/lib/types/streaming'

export function useStreamingManager() {
  const managerRef = useRef<StreamingMessageManager | null>(null)

  // Initialize the manager
  useEffect(() => {
    if (!managerRef.current) {
      managerRef.current = StreamingMessageManager.getInstance({
        maxReconnectAttempts: 5,
        reconnectInterval: 1000,
        backoffMultiplier: 2,
        connectionTimeout: 30000
      })
    }

    return () => {
      // Cleanup on unmount
      if (managerRef.current) {
        managerRef.current.destroy()
        managerRef.current = null
      }
    }
  }, [])

  const establishAgentConnection = useCallback((agentId: string, sessionId: string) => {
    return managerRef.current?.establishAgentConnection(agentId, sessionId)
  }, [])

  const establishInterAgentConnection = useCallback((
    agents: string[], 
    sessionId: string, 
    userId?: string
  ) => {
    return managerRef.current?.establishInterAgentConnection(agents, sessionId, userId)
  }, [])

  const sendUserMessage = useCallback(async (
    agentId: string, 
    content: string, 
    options?: { sessionId?: string; userId?: string }
  ) => {
    if (!managerRef.current) {
      throw new Error('StreamingMessageManager not initialized')
    }
    return managerRef.current.sendUserMessage(agentId, content, options)
  }, [])

  const sendInterAgentMessage = useCallback(async (request: {
    fromAgent: string
    toAgent: string
    message: string
    sessionId: string
    userId: string
  }) => {
    if (!managerRef.current) {
      throw new Error('StreamingMessageManager not initialized')
    }
    return managerRef.current.sendInterAgentMessage(request)
  }, [])

  const closeConnection = useCallback((connectionId: string) => {
    managerRef.current?.closeConnection(connectionId)
  }, [])

  const onMessage = useCallback((callback: MessageCallback) => {
    return managerRef.current?.onMessage(callback) || (() => {})
  }, [])

  const onInterAgentEvent = useCallback((callback: InterAgentCallback) => {
    return managerRef.current?.onInterAgentEvent(callback) || (() => {})
  }, [])

  const onConnectionChange = useCallback((callback: ConnectionCallback) => {
    return managerRef.current?.onConnectionChange(callback) || (() => {})
  }, [])

  const onError = useCallback((callback: ErrorCallback) => {
    return managerRef.current?.onError(callback) || (() => {})
  }, [])

  const getConnectionStatus = useCallback(() => {
    return managerRef.current?.getConnectionStatus() || {}
  }, [])

  const isConnected = useCallback((connectionId?: string) => {
    return managerRef.current?.isConnected(connectionId) || false
  }, [])

  return {
    establishAgentConnection,
    establishInterAgentConnection,
    sendUserMessage,
    sendInterAgentMessage,
    closeConnection,
    onMessage,
    onInterAgentEvent,
    onConnectionChange,
    onError,
    getConnectionStatus,
    isConnected
  }
}