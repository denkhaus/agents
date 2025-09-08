import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import {
  AgentInfo,
  Message,
  ChatSession,
  LLMEvent,
  ADKSession,
  EventType,
} from "@/lib/types";
import { ConnectionStatus } from "@/lib/types/streaming";
import { apiClient } from "@/lib/api";
import { AGENT_IDS, AgentId, normalizeToAgentId } from "@/lib/constants/agents";
import { parseStructuredThoughts } from "@/lib/parsing";

interface ChatStore {
  // State
  agents: AgentInfo[];
  sessions: Record<AgentId, ChatSession>;
  activeAgentId: AgentId | null;
  isConnected: boolean;
  interAgentEvents: LLMEvent[];
  availableSessions: Record<AgentId, ADKSession[]>; // agentId -> sessions
  currentSessionId: string | null;
  applicationName: string | null;
  isLoadingMessages: boolean;

  // Actions
  setAgents: (agents: AgentInfo[]) => void;
  addMessage: (agentId: AgentId, message: Message) => void;
  updateMessage: (
    agentId: AgentId,
    messageId: string,
    updatedMessage: Partial<Message>
  ) => void;
  setActiveAgent: (agentId: AgentId | null) => void;
  setConnected: (connected: boolean) => void;
  addInterAgentEvent: (event: LLMEvent) => void;
  clearInterAgentEvents: () => void;
  createSession: (agentId: AgentId) => Promise<void>;
  getSession: (agentId: AgentId) => ChatSession | undefined;
  updateAgentStatus: (agentId: AgentId, status: AgentInfo["status"]) => void;
  loadSessions: (agentId: AgentId) => Promise<void>;
  loadSessionMessages: (agentId: AgentId, sessionId: string) => Promise<void>;
  setCurrentSession: (sessionId: string | null) => void;
  deleteSession: (agentId: AgentId, sessionId: string) => Promise<void>;

  // New Simplified Operations for Streaming
  addStreamingMessage: (
    agentId: AgentId,
    messageId: string,
    initialContent: string,
    metadata?: Partial<Message>
  ) => void;
  updateStreamingMessage: (
    messageId: string,
    content: string,
    metadata?: Partial<Message["metadata"]>,
    parts?: Message["parts"]
  ) => void;
  finalizeMessage: (messageId: string) => void;
  setConnectionStatus: (agentId: string, status: ConnectionStatus) => void;
}

// Helper function to convert ADK events to messages
const convertADKEventToMessage = (
  event: LLMEvent,
  // event: {
  //   id?: string;
  //   invocationId?: string;
  //   author?: string;
  //   timestamp?: number;
  //   content?: {
  //     parts?: Array<{ text?: string }>;
  //   };
  //   done?: boolean;
  //   partial?: boolean;
  //   object?: string;
  // },
  agentId: string
): Message => {
  let content = "";
  let parts = undefined;

  if (event.content?.parts && Array.isArray(event.content.parts)) {
    const textParts = event.content.parts.filter(
      (part: { text?: string }) => part.text
    );
    if (textParts.length > 0) {
      content = textParts.map((part: { text?: string }) => part.text).join("");
    }
    parts = event.content.parts;
  }

  let messageType: "user" | "agent" | "inter_agent" | "system" = "agent";
  if (event.type === EventType.TOOL_CALL || event.type === EventType.TOOL_RESPONSE) {
    messageType = "system";
  } else if (
    event.author === AGENT_IDS.HUMAN ||
    (event.author && normalizeToAgentId(event.author) === AGENT_IDS.HUMAN)
  ) {
    messageType = "user";
  }

  // Normalize sender to agent ID
  const rawSender = event.author || agentId;
  const normalizedSender = normalizeToAgentId(rawSender) || rawSender;

  return {
    id: event.id || event.invocationId || Date.now().toString(),
    content,
    timestamp: new Date((event.timestamp || Date.now() / 1000) * 1000),
    sender: normalizedSender as AgentId,
    type: messageType,
    metadata: {
      invocationId: event.invocationId,
      partial: event.partial,
      done: event.done,
      fromAgent: event.fromAgent
        ? normalizeToAgentId(event.fromAgent) || event.fromAgent
        : undefined,
      toAgent: event.toAgent
        ? normalizeToAgentId(event.toAgent) || event.toAgent
        : undefined,
    },
    parts,
  };
};

export const useChatStore = create<ChatStore>()(
  persist(
    (set, get) => ({
      // Initial state
      agents: [],
      sessions: {} as Record<AgentId, ChatSession>,
      activeAgentId: null,
      isConnected: false,
      interAgentEvents: [],
      availableSessions: {} as Record<AgentId, ADKSession[]>,
      currentSessionId: null,
      isLoadingMessages: false,
      applicationName: null,

      // Actions
      setAgents: (agents: AgentInfo[]) => {
        const sortedAgents = [...agents].sort((a, b) =>
          a.name.localeCompare(b.name)
        );
        set({ agents: sortedAgents });

        if (sortedAgents.length) {
          // since all agents have the same ApplicationName we pick the first
          set({ applicationName: sortedAgents[0].applicationName });
        }
      },

      addMessage: (agentId: AgentId, message: Message) => {
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
        agentId: AgentId,
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

      setActiveAgent: async (agentId: AgentId | null) => {
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
              availableSessions.some(
                (s: ADKSession) => s.id === currentSessionId
              )
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
              console.warn(
                "No session ID set after agent activation, force creating session"
              );
              await get().createSession(agentId);
            }
          } catch (error) {
            console.error("Error in setActiveAgent:", error);
            // Fallback: try to create a session anyway
            try {
              await get().createSession(agentId);
            } catch (createError) {
              console.error("Failed to create fallback session:", createError);
            }
          }
        }
      },

      setConnected: (connected: boolean) => set({ isConnected: connected }),

      addInterAgentEvent: (event: LLMEvent) => {
        const currentEvents = get().interAgentEvents;
        set({
          interAgentEvents: [...currentEvents, event],
        });
      },

      clearInterAgentEvents: () => set({ interAgentEvents: [] }),

      createSession: async (agentId: AgentId) => {
        try {
          const appName = get().applicationName;
          const newSession = await apiClient.createSession(appName!, agentId);

          const chatSession: ChatSession = {
            agentId: agentId,
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

      getSession: (agentId: AgentId) => get().sessions[agentId],

      updateAgentStatus: (agentId: AgentId, status: AgentInfo["status"]) => {
        const agents = get().agents;
        const updatedAgents = agents.map((agent: AgentInfo) =>
          agent.id === agentId ? { ...agent, status } : agent
        );
        set({ agents: updatedAgents });
      },

      loadSessions: async (agentId: AgentId) => {
        try {
          const appName = get().applicationName;
          const sessions = await apiClient.getSessions(appName!, agentId);

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

      loadSessionMessages: async (agentId: AgentId, sessionId: string) => {
        set({ isLoadingMessages: true });
        try {
          const appName = get().applicationName;
          const session = await apiClient.getSession(
            appName!,
            agentId,
            sessionId
          );

          if (session && session.events) {
            const backendMessages = session.events
              .map((event: unknown) =>
                convertADKEventToMessage(
                  event as Parameters<typeof convertADKEventToMessage>[0],
                  agentId
                )
              )
              .map((message) => {
                const structuredParts = parseStructuredThoughts(
                  message.content
                );
                if (structuredParts.length > 0) {
                  return {
                    ...message,
                    parts: structuredParts,
                    metadata: {
                      ...message.metadata,
                      hasStructuredThoughts: true,
                    },
                  };
                }
                return message;
              });

            // Get existing session to preserve local messages (like reasoning messages)
            const existingSession = get().sessions[agentId];
            let finalMessages = backendMessages;

            if (existingSession && existingSession.sessionId === sessionId) {
              // Preserve local messages that don't exist in backend
              const localMessages = existingSession.messages.filter(
                (localMsg) => {
                  // Keep messages that are:
                  // 1. Reasoning messages (hasStructuredThoughts)
                  // 2. User messages that might not be synced yet
                  // 3. Messages with partial/streaming state
                  return (
                    localMsg.metadata?.hasStructuredThoughts ||
                    localMsg.sender === AGENT_IDS.HUMAN ||
                    localMsg.metadata?.partial ||
                    !backendMessages.some(
                      (backendMsg) => backendMsg.id === localMsg.id
                    )
                  );
                }
              );

              // Merge local and backend messages, sort by timestamp
              finalMessages = [...localMessages, ...backendMessages].sort(
                (a, b) => a.timestamp.getTime() - b.timestamp.getTime()
              );
            }

            const chatSession: ChatSession = {
              agentId,
              messages: finalMessages,
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

      deleteSession: async (agentId: AgentId, sessionId: string) => {
        console.log("ChatStore: Starting deleteSession for", {
          agentId,
          sessionId,
        });

        try {
          // Try to delete from backend first
          try {
            console.log("ChatStore: Calling apiClient.deleteSession...");
            const appName = get().applicationName;
            await apiClient.deleteSession(appName!, agentId, sessionId);
            console.log("ChatStore: Session deleted from backend successfully");
          } catch (backendError) {
            console.warn(
              "ChatStore: Backend delete failed, proceeding with local deletion:",
              backendError
            );

            // Check if it's a CORS error and provide specific handling
            if (
              backendError instanceof Error &&
              backendError.message.includes("CORS")
            ) {
              console.log(
                "ChatStore: CORS error detected - session will be removed locally only"
              );
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

      // New Simplified Operations for Streaming
      addStreamingMessage: (
        agentId: AgentId,
        messageId: string,
        initialContent: string,
        metadata?: Partial<Message>
      ) => {
        // console.log(`[CHAT STORE] addStreamingMessage called`, { agentId, messageId, contentLength: initialContent.length });

        const sessions = get().sessions;
        const session = sessions[agentId];

        if (!session) {
          // console.warn(`Cannot add streaming message: no session found for agent ${agentId}`);
          return;
        }

        // Check if message already exists (prevent duplicates)
        const existingMessageIndex = session.messages.findIndex(
          (msg) => msg.id === messageId
        );
        if (existingMessageIndex !== -1) {
          // Message exists - this is expected for streaming, just return silently
          // console.log(`[CHAT STORE] Message already exists, skipping add`, { messageId });
          return;
        }

        const newMessage: Message = {
          id: messageId,
          content: initialContent,
          timestamp: new Date(),
          sender: metadata?.sender || agentId,
          type: metadata?.type || "agent",
          metadata: {
            partial: true,
            done: false,
            streamingKey: messageId,
            ...metadata?.metadata,
          },
          ...metadata,
        };

        // console.log(`[CHAT STORE] Adding new streaming message`, {
        //   messageId,
        //   contentLength: newMessage.content.length,
        //   sender: newMessage.sender,
        //   type: newMessage.type
        // });

        const updatedSession = {
          ...session,
          messages: [...session.messages, newMessage],
          lastActivity: new Date(),
        };

        set({
          sessions: {
            ...sessions,
            [agentId]: updatedSession,
          },
        });
      },

      updateStreamingMessage: (
        messageId: string,
        content: string,
        metadata?: Partial<Message["metadata"]>,
        parts?: Message["parts"]
      ) => {
        // console.log(`[CHAT STORE] updateStreamingMessage called`, { messageId, contentLength: content.length });

        const sessions = get().sessions;
        let messageFound = false;

        // Find the session containing the message
        for (const [agentId, session] of Object.entries(sessions)) {
          const messageIndex = session.messages.findIndex(
            (msg) => msg.id === messageId
          );

          if (messageIndex !== -1) {
            messageFound = true;
            const updatedMessages = [...session.messages];
            const currentMessage = updatedMessages[messageIndex];

            // Update the message with accumulated content
            updatedMessages[messageIndex] = {
              ...currentMessage,
              content, // This should be the full accumulated content
              parts: parts || currentMessage.parts, // Update parts if provided
              timestamp: new Date(), // Update timestamp for last update
              metadata: {
                ...currentMessage.metadata,
                ...metadata,
                streamingKey: messageId, // Ensure streamingKey is preserved
              },
            };

            // console.log(`[CHAT STORE] Updating existing message`, {
            //   messageId,
            //   oldContentLength: currentMessage.content.length,
            //   newContentLength: content.length,
            //   isPartial: updatedMessages[messageIndex].metadata?.partial,
            //   isDone: updatedMessages[messageIndex].metadata?.done
            // });

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
            break;
          }
        }

        if (!messageFound) {
          // console.warn(`UpdateStreamingMessage: Message with ID ${messageId} not found`);
        }
      },

      finalizeMessage: (messageId: string) => {
        // console.log(`[CHAT STORE] finalizeMessage called`, { messageId });

        const sessions = get().sessions;
        let messageFound = false;

        // Find and finalize the message
        for (const [agentId, session] of Object.entries(sessions)) {
          const messageIndex = session.messages.findIndex(
            (msg) => msg.id === messageId
          );

          if (messageIndex !== -1) {
            messageFound = true;
            const updatedMessages = [...session.messages];
            const currentMessage = updatedMessages[messageIndex];

            updatedMessages[messageIndex] = {
              ...currentMessage,
              metadata: {
                ...currentMessage.metadata,
                partial: false,
                done: true,
                finalizedAt: new Date().toISOString(),
              },
            };

            // console.log(`[CHAT STORE] Finalizing message`, {
            //   messageId,
            //   wasPartial: currentMessage.metadata?.partial,
            //   contentLength: currentMessage.content.length
            // });

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
            break;
          }
        }

        if (!messageFound) {
          // console.warn(`FinalizeMessage: Message with ID ${messageId} not found`);
        }
      },

      setConnectionStatus: (agentId: string, status: ConnectionStatus) => {
        // Update agent connection status if needed
        // This could be expanded to track per-agent connection states
        console.log(`Connection status for ${agentId}:`, status);
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
