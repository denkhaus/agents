import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import {
  Agent,
  Message,
  ChatSession,
  InterAgentEvent,
  ADKSession,
} from "@/lib/types";
import { apiClient } from "@/lib/api";

interface ChatStore {
  // State
  agents: Agent[];
  sessions: Record<string, ChatSession>;
  activeAgentId: string | null;
  isConnected: boolean;
  interAgentEvents: InterAgentEvent[];
  availableSessions: Record<string, ADKSession[]>; // agentId -> sessions
  currentSessionId: string | null;
  isLoadingMessages: boolean;

  // Actions
  setAgents: (agents: Agent[]) => void;
  addMessage: (agentId: string, message: Message) => void;
  updateMessage: (
    agentId: string,
    messageId: string,
    updatedMessage: Partial<Message>
  ) => void;
  setActiveAgent: (agentId: string | null) => void;
  setConnected: (connected: boolean) => void;
  addInterAgentEvent: (event: InterAgentEvent) => void;
  clearInterAgentEvents: () => void;
  createSession: (agentId: string) => Promise<void>;
  getSession: (agentId: string) => ChatSession | undefined;
  updateAgentStatus: (agentId: string, status: Agent["status"]) => void;
  loadSessions: (agentId: string) => Promise<void>;
  loadSessionMessages: (agentId: string, sessionId: string) => Promise<void>;
  setCurrentSession: (sessionId: string | null) => void;
  deleteSession: (agentId: string, sessionId: string) => Promise<void>;
}

// Helper function to convert ADK events to messages
const convertADKEventToMessage = (event: { 
  id?: string; 
  invocationId?: string; 
  author?: string; 
  timestamp?: number; 
  content?: { 
    parts?: Array<{ text?: string }> 
  }; 
  done?: boolean; 
  partial?: boolean; 
  object?: string 
}, agentId: string): Message => {
  let content = "";
  let parts = undefined;

  if (event.content?.parts && Array.isArray(event.content.parts)) {
    const textParts = event.content.parts.filter((part: { text?: string }) => part.text);
    if (textParts.length > 0) {
      content = textParts.map((part: { text?: string }) => part.text).join("");
    }
    parts = event.content.parts;
  }

  let messageType: "user" | "agent" | "inter_agent" | "system" = "agent";
  if (event.object === "tool_call" || event.object === "tool_response") {
    messageType = "system";
  } else if (event.author === "user") {
    messageType = "user";
  }

  return {
    id: event.id || event.invocationId || Date.now().toString(),
    content,
    timestamp: new Date((event.timestamp || Date.now() / 1000) * 1000),
    sender: event.author || agentId,
    type: messageType,
    metadata: {
      invocationId: event.invocationId,
      partial: event.partial,
      done: event.done,
    },
    parts,
  };
};

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
      isLoadingMessages: false,

      // Actions
      setAgents: (agents: Agent[]) => {
        const sortedAgents = [...agents].sort((a, b) =>
          a.name.localeCompare(b.name)
        );
        set({ agents: sortedAgents });
      },

      addMessage: (agentId: string, message: Message) => {
        const sessions = get().sessions;
        const session = sessions[agentId];

        if (session) {
          const updatedSession = {
            ...session,
            messages: [...session.messages, message],
            lastActivity: new Date(),
          };

          set({
            sessions: {
              ...sessions,
              [agentId]: updatedSession,
            },
          });
        }
      },

      updateMessage: (
        agentId: string,
        messageId: string,
        updatedMessage: Partial<Message>
      ) => {
        const sessions = get().sessions;
        const session = sessions[agentId];

        if (session) {
          const updatedMessages = session.messages.map((msg: Message) =>
            msg.id === messageId ? { ...msg, ...updatedMessage } : msg
          );

          const updatedSession = {
            ...session,
            messages: updatedMessages,
            lastActivity: new Date(),
          };

          set({
            sessions: {
              ...sessions,
              [agentId]: updatedSession,
            },
          });
        }
      },

      setActiveAgent: async (agentId: string | null) => {
        set({ activeAgentId: agentId });

        if (agentId) {
          try {
            // Load available sessions for this agent
            await get().loadSessions(agentId);

            // If we have a persisted currentSessionId, check if it exists in available sessions
            const currentSessionId = get().currentSessionId;
            const availableSessions = get().availableSessions[agentId] || [];

            if (
              currentSessionId &&
              availableSessions.some((s: ADKSession) => s.id === currentSessionId)
            ) {
              try {
                await get().loadSessionMessages(agentId, currentSessionId);
              } catch (error) {
                console.warn(
                  "Failed to load persisted session, clearing current session:",
                  error
                );
                set({ currentSessionId: null });
              }
            } else if (currentSessionId) {
              // If currentSessionId is set but not found in availableSessions (e.g., deleted on backend)
              console.warn(
                `Persisted session ID ${currentSessionId} not found for agent ${agentId}, clearing.`
              );
              set({ currentSessionId: null });
            }

            // If no persisted session was loaded or found, load the most recent available session
            if (!get().currentSessionId && availableSessions.length > 0) {
              const mostRecentSession = availableSessions.sort(
                (a: ADKSession, b: ADKSession) =>
                  b.lastUpdateTime - a.lastUpdateTime
              )[0];
              await get().loadSessionMessages(agentId, mostRecentSession.id);
            } else if (
              !get().currentSessionId &&
              availableSessions.length === 0
            ) {
              // If no sessions are available, create a new one
              await get().createSession(agentId);
            }
            
            // Final check to ensure we have a current session
            const finalSessionId = get().currentSessionId;
            
            // If we still don't have a session, force create one
            if (!finalSessionId) {
              console.warn('No session ID set after agent activation, force creating session');
              await get().createSession(agentId);
            }
          } catch (error) {
            console.error('Error in setActiveAgent:', error);
            // Fallback: try to create a session anyway
            try {
              await get().createSession(agentId);
            } catch (createError) {
              console.error('Failed to create fallback session:', createError);
            }
          }
        }
      },

      setConnected: (connected: boolean) => set({ isConnected: connected }),

      addInterAgentEvent: (event: InterAgentEvent) => {
        const currentEvents = get().interAgentEvents;
        set({
          interAgentEvents: [...currentEvents, event],
        });
      },

      clearInterAgentEvents: () => set({ interAgentEvents: [] }),

      createSession: async (agentId: string) => {
        try {
          const newSession = await apiClient.createSession(agentId, "user");

          const chatSession: ChatSession = {
            agentId,
            messages: [],
            isActive: true,
            lastActivity: new Date(),
            sessionId: newSession.id,
          };

          set({
            sessions: {
              ...get().sessions,
              [agentId]: chatSession,
            },
            currentSessionId: newSession.id,
          });

          // Refresh available sessions
          await get().loadSessions(agentId);
        } catch (error) {
          console.error("Error creating session:", error);
          // Fallback to local session
          const sessions = get().sessions;
          if (!sessions[agentId]) {
            const newSession: ChatSession = {
              agentId,
              messages: [],
              isActive: false,
              lastActivity: new Date(),
            };

            set({
              sessions: {
                ...sessions,
                [agentId]: newSession,
              },
            });
          }
        }
      },

      getSession: (agentId: string) => get().sessions[agentId],

      updateAgentStatus: (agentId: string, status: Agent["status"]) => {
        const agents = get().agents;
        const updatedAgents = agents.map((agent: Agent) =>
          agent.id === agentId ? { ...agent, status } : agent
        );
        set({ agents: updatedAgents });
      },

      loadSessions: async (agentId: string) => {
        try {
          const sessions = await apiClient.getSessions(agentId, "user");
          set({
            availableSessions: {
              ...get().availableSessions,
              [agentId]: sessions,
            },
          });
        } catch (error) {
          console.error("Error loading sessions:", error);
        }
      },

      loadSessionMessages: async (agentId: string, sessionId: string) => {
        set({ isLoadingMessages: true });
        try {
          const session = await apiClient.getSession(
            agentId,
            "user",
            sessionId
          );

          if (session && session.events) {
            const messages = session.events.map((event: unknown) =>
              convertADKEventToMessage(event as Parameters<typeof convertADKEventToMessage>[0], agentId)
            );

            const chatSession: ChatSession = {
              agentId,
              messages,
              isActive: true,
              lastActivity: new Date(session.lastUpdateTime * 1000),
              sessionId: session.id,
            };

            set({
              sessions: {
                ...get().sessions,
                [agentId]: chatSession,
              },
              currentSessionId: sessionId,
            });
          }
        } catch (error) {
          console.error("Error loading session messages:", error);
        } finally {
          set({ isLoadingMessages: false });
        }
      },

      setCurrentSession: (sessionId: string | null) => {
        set({ currentSessionId: sessionId });
      },

      deleteSession: async (agentId: string, sessionId: string) => {
        console.log('ChatStore: Starting deleteSession for', { agentId, sessionId });
        
        try {
          // Try to delete from backend first
          try {
            console.log('ChatStore: Calling apiClient.deleteSession...');
            await apiClient.deleteSession(agentId, "user", sessionId);
            console.log("ChatStore: Session deleted from backend successfully");
          } catch (backendError) {
            console.warn(
              "ChatStore: Backend delete failed, proceeding with local deletion:",
              backendError
            );
            
            // Check if it's a CORS error and provide specific handling
            if (backendError instanceof Error && backendError.message.includes('CORS')) {
              console.log('ChatStore: CORS error detected - session will be removed locally only');
              // For CORS errors, we'll just proceed with local deletion
              // The session might still exist on the server, but we can't delete it from the browser
            }
            
            // Continue with local deletion even if backend fails
          }

          // Remove from available sessions (always do this)
          const availableSessions = get().availableSessions;
          const agentSessions = availableSessions[agentId] || [];
          const updatedSessions = agentSessions.filter(
            (session: ADKSession) => session.id !== sessionId
          );

          set({
            availableSessions: {
              ...availableSessions,
              [agentId]: updatedSessions,
            },
          });

          // If deleted session was current, clear it and create new one
          if (get().currentSessionId === sessionId) {
            set({ currentSessionId: null });

            // Remove from local sessions
            const sessions = get().sessions;
            const updatedLocalSessions = { ...sessions };
            delete updatedLocalSessions[agentId];
            set({ sessions: updatedLocalSessions });

            // Create new session if no others exist
            if (updatedSessions.length === 0) {
              await get().createSession(agentId);
            } else {
              // Switch to the most recent session
              const mostRecent = updatedSessions.sort(
                (a: ADKSession, b: ADKSession) =>
                  b.lastUpdateTime - a.lastUpdateTime
              )[0];
              await get().loadSessionMessages(agentId, mostRecent.id);
            }
          }
        } catch (error) {
          console.error("Error in deleteSession:", error);
          // Don't throw error to prevent UI from breaking
          // The session will be removed from local state even if backend fails
        }
      },
    }),
    {
      name: "chat-session-storage",
      storage: createJSONStorage(() => localStorage),
      // Only persist the currentSessionId and activeAgentId
      partialize: (state: ChatStore) => ({
        currentSessionId: state.currentSessionId,
        activeAgentId: state.activeAgentId,
      }),
    }
  )
);
