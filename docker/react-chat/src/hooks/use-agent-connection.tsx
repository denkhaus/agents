'use client'

import { useQuery } from '@tanstack/react-query'
import { useEffect } from 'react'
import { agentApi, sseService } from '@/lib/api'
import { useChatStore } from '@/lib/store'
import { messageApi } from '@/lib/api'

export function useAgentConnection() {
  const { 
    setAgents, 
    setConnected, 
    addInterAgentEvent, 
    addMessage,
    agents,
    activeAgentId,
    setActiveAgent
  } = useChatStore()

  // Fetch agents with better error handling
  const { data: agentsData, isLoading, error, isError, isFetching } = useQuery({
    queryKey: ['agents'],
    queryFn: async () => {
      try {
        const result = await agentApi.getAgents()
        return result
      } catch (error) {
        console.error('Failed to fetch agents:', error)
        throw error
      }
    },
    refetchInterval: 30000, // Refresh every 30 seconds
    retry: 3,
    retryDelay: 1000,
    staleTime: 0, // Always consider data stale
    gcTime: 0, // Don't cache data
    refetchOnMount: true,
    refetchOnWindowFocus: true,
  })

  // Update store when agents change
  useEffect(() => {
    if (agentsData && agentsData.length > 0) {
      setAgents(agentsData)
      
      // Auto-select first agent if none is selected
      if (!activeAgentId) {
        // Use async function to handle the promise
        const selectAgent = async () => {
          try {
            await setActiveAgent(agentsData[0].id)
          } catch (error) {
            console.error('Failed to set active agent:', error)
          }
        }
        selectAgent()
      }
    }
  }, [agentsData, setAgents, setActiveAgent, activeAgentId])

  // Set up SSE connection when agents are available
  useEffect(() => {
    if (agents.length > 0) {
      const agentIds = agents.map(agent => agent.id)
      const sessionId = `session-${Date.now()}`
      const userId = 'user'

      let reconnectAttempts = 0
      const maxReconnectAttempts = 5

      const connectWithRetry = () => {
        sseService.connectMainChat(agentIds, sessionId, userId, {
          onConnectionStatusChange: (connected) => {
            setConnected(connected)
            if (connected) {
              reconnectAttempts = 0 // Reset on successful connection
            }
          },
          onInterAgentEvent: (event) => {
            addInterAgentEvent(event)
          },
          onMessage: (event) => {
            try {
              // Convert event to message and add to appropriate agent session
              const message = messageApi.convertEventToMessage(event)
              if (event.fromAgent && agents.some(a => a.id === event.fromAgent)) {
                addMessage(event.fromAgent, message)
              }
            } catch (error) {
              console.error('Error converting event to message:', error, event)
            }
          },
          onError: (error) => {
            console.error('SSE error:', error)
            setConnected(false)
            
            // Retry connection with exponential backoff
            if (reconnectAttempts < maxReconnectAttempts) {
              reconnectAttempts++
              const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000)
              setTimeout(connectWithRetry, delay)
            } else {
              console.error('Max reconnection attempts reached')
            }
          }
        })
      }

      connectWithRetry()

      return () => {
        sseService.disconnect('mainChat')
      }
    }
  }, [agents, setConnected, addInterAgentEvent, addMessage])

  return {
    agents: agentsData || [],
    isLoading,
    error,
    isError,
    isFetching
  }
}