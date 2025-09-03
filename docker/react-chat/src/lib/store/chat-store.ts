import { create } from 'zustand'
import { Agent, Message, ChatSession, InterAgentEvent } from '@/lib/types'

interface ChatStore {
  // State
  agents: Agent[]
  sessions: Record<string, ChatSession>
  activeAgentId: string | null
  isConnected: boolean
  interAgentEvents: InterAgentEvent[]
  
  // Actions
  setAgents: (agents: Agent[]) => void
  addMessage: (agentId: string, message: Message) => void
  updateMessage: (agentId: string, messageId: string, updatedMessage: Partial<Message>) => void
  setActiveAgent: (agentId: string | null) => void
  setConnected: (connected: boolean) => void
  addInterAgentEvent: (event: InterAgentEvent) => void
  clearInterAgentEvents: () => void
  createSession: (agentId: string) => void
  getSession: (agentId: string) => ChatSession | undefined
  updateAgentStatus: (agentId: string, status: Agent['status']) => void
}

export const useChatStore = create<ChatStore>((set, get) => ({
  // Initial state
  agents: [],
  sessions: {},
  activeAgentId: null,
  isConnected: false,
  interAgentEvents: [],
  
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
  
  setActiveAgent: (agentId) => {
    set({ activeAgentId: agentId })
    
    // Create session if it doesn't exist
    if (agentId && !get().sessions[agentId]) {
      get().createSession(agentId)
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
  
  createSession: (agentId) => {
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
  },
  
  getSession: (agentId) => get().sessions[agentId],
  
  updateAgentStatus: (agentId, status) => {
    const agents = get().agents
    const updatedAgents = agents.map(agent =>
      agent.id === agentId ? { ...agent, status } : agent
    )
    set({ agents: updatedAgents })
  }
}))