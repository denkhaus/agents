import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { Agent, Message, ChatSession, InterAgentEvent, ADKSession } from '@/lib/types'
import { apiClient } from '@/lib/api'

interface ChatStore {
  // State
  agents: Agent[]
  sessions: Record<string, ChatSession>
  activeAgentId: string | null
  isConnected: boolean
  interAgentEvents: InterAgentEvent[]
  availableSessions: Record<string, ADKSession[]> // agentId -> sessions
  currentSessionId: string | null
  
  // Actions
  setAgents: (agents: Agent[]) => void
  addMessage: (agentId: string, message: Message) => void
  updateMessage: (agentId: string, messageId: string, updatedMessage: Partial<Message>) => void
  setActiveAgent: (agentId: string | null) => void
  setConnected: (connected: boolean) => void
  addInterAgentEvent: (event: InterAgentEvent) => void
  clearInterAgentEvents: () => void
  createSession: (agentId: string) => Promise<void>
  getSession: (agentId: string) => ChatSession | undefined
  updateAgentStatus: (agentId: string, status: Agent['status']) => void
  loadSessions: (agentId: string) => Promise<void>
  loadSessionMessages: (agentId: string, sessionId: string) => Promise<void>
  setCurrentSession: (sessionId: string | null) => void
  deleteSession: (agentId: string, sessionId: string) => Promise<void>
}

// Helper function to convert ADK events to messages
const convertADKEventToMessage = (event: any, agentId: string): Message => {
  let content = ''
  let parts = undefined

  if (event.content?.parts && Array.isArray(event.content.parts)) {
    const textParts = event.content.parts.filter((part: any) => part.text)
    if (textParts.length > 0) {
      content = textParts.map((part: any) => part.text).join('')
    }
    parts = event.content.parts
  }

  let messageType: 'user' | 'agent' | 'inter_agent' | 'system' = 'agent'
  if (event.object === 'tool_call' || event.object === 'tool_response') {
    messageType = 'system'
  } else if (event.author === 'user') {
    messageType = 'user'
  }

  return {
    id: event.id || event.invocationId || Date.now().toString(),
    content,
    timestamp: new Date(event.timestamp * 1000),
    sender: event.author || agentId,
    type: messageType,
    metadata: {
      invocationId: event.invocationId,
      partial: event.partial,
      done: event.done
    },
    parts
  }
}

export const useChatStore = create<ChatStore>()(
  persist(
    (set, get) => ({
      // Initial state
      agents: [],
      sessions: {},
      activeAgentId: null,
      isConnected: false,
      interAgentEvents: [],
      availableSessions: {},
      currentSessionId: null,
  
  // Actions
  setAgents: (agents) => set({ agents }),
  
  addMessage: (agentId, message) => {
    const sessions = get().sessions
    const session = sessions[agentId]
    
    if (session) {
      const updatedSession = {
        ...session,
        messages: [...session.messages, message],
        lastActivity: new Date()
      }
      
      set({
        sessions: {
          ...sessions,
          [agentId]: updatedSession
        }
      })
    }
  },

  updateMessage: (agentId, messageId, updatedMessage) => {
    const sessions = get().sessions
    const session = sessions[agentId]
    
    if (session) {
      const updatedMessages = session.messages.map(msg =>
        msg.id === messageId ? { ...msg, ...updatedMessage } : msg
      )
      
      const updatedSession = {
        ...session,
        messages: updatedMessages,
        lastActivity: new Date()
      }
      
      set({
        sessions: {
          ...sessions,
          [agentId]: updatedSession
        }
      })
    }
  },
  
  setActiveAgent: async (agentId) => {
    set({ activeAgentId: agentId })
    
    if (agentId) {
      // Load available sessions for this agent
      await get().loadSessions(agentId)
      
      // If we have a persisted currentSessionId, try to load that session
      const currentSessionId = get().currentSessionId
      if (currentSessionId) {
        try {
          await get().loadSessionMessages(agentId, currentSessionId)
        } catch (error) {
          console.warn('Failed to load persisted session, clearing current session:', error)
          set({ currentSessionId: null })
        }
      }
      // Don't automatically create a new session - let user create one via Session Panel
    }
  },
  
  setConnected: (connected) => set({ isConnected: connected }),
  
  addInterAgentEvent: (event) => {
    const currentEvents = get().interAgentEvents
    set({
      interAgentEvents: [...currentEvents, event]
    })
  },
  
  clearInterAgentEvents: () => set({ interAgentEvents: [] }),
  
  createSession: async (agentId) => {
    try {
      const newSession = await apiClient.createSession(agentId, 'user')
      
      const chatSession: ChatSession = {
        agentId,
        messages: [],
        isActive: true,
        lastActivity: new Date(),
        sessionId: newSession.id
      }
      
      set({
        sessions: {
          ...get().sessions,
          [agentId]: chatSession
        },
        currentSessionId: newSession.id
      })
      
      // Refresh available sessions
      await get().loadSessions(agentId)
    } catch (error) {
      console.error('Error creating session:', error)
      // Fallback to local session
      const sessions = get().sessions
      if (!sessions[agentId]) {
        const newSession: ChatSession = {
          agentId,
          messages: [],
          isActive: false,
          lastActivity: new Date()
        }
        
        set({
          sessions: {
            ...sessions,
            [agentId]: newSession
          }
        })
      }
    }
  },
  
  getSession: (agentId) => get().sessions[agentId],
  
  updateAgentStatus: (agentId, status) => {
    const agents = get().agents
    const updatedAgents = agents.map(agent =>
      agent.id === agentId ? { ...agent, status } : agent
    )
    set({ agents: updatedAgents })
  },

  loadSessions: async (agentId) => {
    try {
      const sessions = await apiClient.getSessions(agentId, 'user')
      set({
        availableSessions: {
          ...get().availableSessions,
          [agentId]: sessions
        }
      })
    } catch (error) {
      console.error('Error loading sessions:', error)
    }
  },

  loadSessionMessages: async (agentId, sessionId) => {
    try {
      const session = await apiClient.getSession(agentId, 'user', sessionId)
      
      if (session && session.events) {
        const messages = session.events.map(event => 
          convertADKEventToMessage(event, agentId)
        )
        
        const chatSession: ChatSession = {
          agentId,
          messages,
          isActive: true,
          lastActivity: new Date(session.lastUpdateTime * 1000),
          sessionId: session.id
        }
        
        set({
          sessions: {
            ...get().sessions,
            [agentId]: chatSession
          },
          currentSessionId: sessionId
        })
      }
    } catch (error) {
      console.error('Error loading session messages:', error)
    }
  },

  setCurrentSession: (sessionId) => {
    set({ currentSessionId: sessionId })
  },

  deleteSession: async (agentId, sessionId) => {
    try {
      // Try to delete from backend first
      try {
        await apiClient.deleteSession(agentId, 'user', sessionId)
        console.log('Session deleted from backend successfully')
      } catch (backendError) {
        console.warn('Backend delete failed, proceeding with local deletion:', backendError)
        // Continue with local deletion even if backend fails
      }
      
      // Remove from available sessions (always do this)
      const availableSessions = get().availableSessions
      const agentSessions = availableSessions[agentId] || []
      const updatedSessions = agentSessions.filter(session => session.id !== sessionId)
      
      set({
        availableSessions: {
          ...availableSessions,
          [agentId]: updatedSessions
        }
      })
      
      // If deleted session was current, clear it and create new one
      if (get().currentSessionId === sessionId) {
        set({ currentSessionId: null })
        
        // Remove from local sessions
        const sessions = get().sessions
        const updatedLocalSessions = { ...sessions }
        delete updatedLocalSessions[agentId]
        set({ sessions: updatedLocalSessions })
        
        // Create new session if no others exist
        if (updatedSessions.length === 0) {
          await get().createSession(agentId)
        } else {
          // Switch to the most recent session
          const mostRecent = updatedSessions.sort((a, b) => b.lastUpdateTime - a.lastUpdateTime)[0]
          await get().loadSessionMessages(agentId, mostRecent.id)
        }
      }
    } catch (error) {
      console.error('Error in deleteSession:', error)
      // Don't throw error to prevent UI from breaking
      // The session will be removed from local state even if backend fails
    }
  }
    }),
    {
      name: 'chat-session-storage',
      storage: createJSONStorage(() => localStorage),
      // Only persist the currentSessionId and activeAgentId
      partialize: (state) => ({
        currentSessionId: state.currentSessionId,
        activeAgentId: state.activeAgentId,
      }),
    }
  )
)