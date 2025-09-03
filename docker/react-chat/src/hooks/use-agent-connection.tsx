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
      console.log('React Query: Fetching agents from backend...')
      try {
        const result = await agentApi.getAgents()
        console.log('React Query: Agents fetched successfully:', result)
        return result
      } catch (error) {
        console.error('React Query: Failed to fetch agents:', error)
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
      console.log('Setting agents in store:', agentsData)
      setAgents(agentsData)
      
      // Auto-select first agent if none is selected
      if (!activeAgentId) {
        console.log('Auto-selecting first agent:', agentsData[0].id)
        setActiveAgent(agentsData[0].id)
      }
    }
  }, [agentsData, setAgents, setActiveAgent, activeAgentId])

  // Set up SSE connection when agents are available
  useEffect(() => {
    if (agents.length > 0) {
      console.log('Setting up SSE connection for agents:', agents.map(a => a.id))
      const agentIds = agents.map(agent => agent.id)
      const sessionId = `session-${Date.now()}`
      const userId = 'user'

      let reconnectAttempts = 0
      const maxReconnectAttempts = 5

      const connectWithRetry = () => {
        sseService.connect(agentIds, sessionId, userId, {
          onConnectionStatusChange: (connected) => {
            console.log('SSE connection status changed:', connected)
            setConnected(connected)
            if (connected) {
              reconnectAttempts = 0 // Reset on successful connection
            }
          },
          onInterAgentEvent: (event) => {
            console.log('Inter-agent event received:', event)
            addInterAgentEvent(event)
          },
          onMessage: (event) => {
            console.log('Message event received:', event)
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
              console.log(`Retrying SSE connection in ${delay}ms (attempt ${reconnectAttempts}/${maxReconnectAttempts})`)
              setTimeout(connectWithRetry, delay)
            } else {
              console.error('Max reconnection attempts reached')
            }
          }
        })
      }

      connectWithRetry()

      return () => {
        console.log('Cleaning up SSE connection')
        sseService.disconnect()
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