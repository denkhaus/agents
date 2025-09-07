"use client"

import { useEffect, useRef, useCallback, useMemo } from 'react'
import { useChatStore } from '@/lib/store'
import { useStreamingManager } from './use-streaming-manager'
import { LLMEvent } from '@/lib/types'
import { ConnectionStatus } from '@/lib/types/streaming'
import { debug } from '@/lib/utils/debug'

interface UseInterAgentCommunicationProps {
  agents: string[]
  sessionId: string
  userId?: string
}

export function useInterAgentCommunication({
  agents,
  sessionId,
  userId = 'user'
}: UseInterAgentCommunicationProps) {
  const { addInterAgentEvent } = useChatStore()
  const streamingManager = useStreamingManager()
  const connectionRef = useRef<string | null>(null)

  // Memoize agents array to prevent unnecessary re-renders
  const agentsKey = useMemo(() => agents.join(','), [JSON.stringify(agents)])
  // Create a stable key for all dependencies
  const stableKey = useMemo(() => `${agentsKey}-${sessionId}-${userId}`, [agentsKey, sessionId, userId]);

  // Create stable references to prevent infinite loops
  const stableDepsRef = useRef({
    addInterAgentEvent
  })

  // Update stable references when dependencies change
  useEffect(() => {
    stableDepsRef.current = {
      addInterAgentEvent
    }
  }, [addInterAgentEvent])

  // Memoize event handlers to prevent re-creation on every render
  const handleInterAgentEvent = useCallback((event: LLMEvent) => {
    debug.streaming('Inter-agent event received:', event)
    stableDepsRef.current.addInterAgentEvent(event)
  }, [])

  const handleConnectionChange = useCallback((status: ConnectionStatus) => {
    debug.connection('Inter-agent connection status changed:', status)
  }, [])

  const handleError = useCallback((error: Error) => {
    debug.error('Inter-agent connection error:', error)
  }, [])

  useEffect(() => {
    if (agents.length === 0 || !sessionId) return

    debug.critical('Setting up inter-agent communication for:', { agents, sessionId })

    // Check if we already have an active connection for these agents/session
    const connectionId = `inter-agent-${agents.join('-')}-${sessionId}`
    const connectionStatus = streamingManager.getConnectionStatus()
    const existingConnection = connectionStatus[connectionId]
    
    if (existingConnection?.isConnected) {
      debug.connection(`Reusing existing inter-agent connection for ${connectionId}`)
      connectionRef.current = connectionId
      return
    }

    // Establish inter-agent connection
    const connection = streamingManager.establishInterAgentConnection(agents, sessionId, userId)
    
    if (connection) {
      connectionRef.current = connection.id

      // Set up event listeners with memoized handlers
      const unsubscribeInterAgentEvent = streamingManager.onInterAgentEvent(handleInterAgentEvent)
      const unsubscribeConnectionChange = streamingManager.onConnectionChange(handleConnectionChange)
      const unsubscribeError = streamingManager.onError(handleError)

      // Cleanup function
      return () => {
        debug.critical('Cleaning up inter-agent communication')
        unsubscribeInterAgentEvent()
        unsubscribeConnectionChange()
        unsubscribeError()
        
        // Don't close the connection immediately to allow for reuse
        connectionRef.current = null
      }
    }
  }, [stableKey]) // Only depend on the stable key, not on streamingManager

  // Return connection status and utility functions
  return {
    isConnected: streamingManager.isConnected(connectionRef.current || undefined),
    connectionId: connectionRef.current,
    sendInterAgentMessage: streamingManager.sendInterAgentMessage
  }
}