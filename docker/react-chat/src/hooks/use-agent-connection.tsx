"use client";

import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { agentApi, sseService } from "@/lib/api";
import { useChatStore } from "@/lib/store";
import { messageApi } from "@/lib/api";

export function useAgentConnection() {
  const {
    setAgents,
    setConnected,
    addInterAgentEvent,
    addMessage,
    agents,
    activeAgentId,
    setActiveAgent,
  } = useChatStore();

  // Fetch agents with better error handling
  const {
    data: agentsData,
    isLoading,
    error,
    isError,
    isFetching,
  } = useQuery({
    queryKey: ["agents"],
    queryFn: async () => {
      try {
        const result = await agentApi.getAgents();
        return result;
      } catch (error) {
        console.error("Failed to fetch agents:", error);
        throw error;
      }
    },
    refetchInterval: 30000, // Refresh every 30 seconds
    retry: 3,
    retryDelay: 1000,
    staleTime: 0, // Always consider data stale
    gcTime: 0, // Don't cache data
    refetchOnMount: true,
    refetchOnWindowFocus: true,
  });

  // Update store when agents change
  useEffect(() => {
    if (agentsData && agentsData.length > 0) {
      setAgents(agentsData);

      // Auto-select first agent if none is selected
      if (!activeAgentId) {
        // Use async function to handle the promise
        const selectAgent = async () => {
          try {
            await setActiveAgent(agentsData[0].id);
          } catch (error) {
            console.error("Failed to set active agent:", error);
          }
        };
        selectAgent();
      }
    }
  }, [agentsData, setAgents, setActiveAgent, activeAgentId]);

  // Set up SSE connection when agents are available
  useEffect(() => {
    if (agents.length > 0) {
      const agentIds = agents.map((agent) => agent.id);
      const sessionId = `session-${Date.now()}`;
      const userId = "user";

      let reconnectAttempts = 0;
      const maxReconnectAttempts = 5;

      const connectWithRetry = () => {
        sseService.connectMainChat(agentIds, sessionId, userId, {
          onConnectionStatusChange: (connected) => {
            setConnected(connected);
            if (connected) {
              reconnectAttempts = 0; // Reset on successful connection
            }
          },
          onInterAgentEvent: (event) => {
            addInterAgentEvent(event);
          },
          onMessage: (event) => {
            try {
              // Handle ONLY inter-agent and system messages
              // User-initiated messages are handled by MessageInput component
              console.log(
                "use-agent-connection: SSE Main Chat Event:",
                JSON.stringify({
                  type: event.type,
                  object: event.object,
                  fromAgent: event.fromAgent,
                  toAgent: event.toAgent,
                  partial: event.partial,
                })
              );

              // Handle inter-agent messages based on object type as well as event type
              if (
                event.type === "inter_agent" ||
                event.type === "communication" ||
                event.object === "inter_agent"
              ) {
                const message = messageApi.convertEventToMessage(event);
                // Display inter-agent messages in the INITIATING agent's chat
                const targetAgentId =
                  event.fromAgent || event.toAgent || event.author;
                if (
                  targetAgentId &&
                  agents.some((a) => a.id === targetAgentId)
                ) {
                  addMessage(targetAgentId, {
                    ...message,
                    type: "inter_agent" as const,
                    metadata: {
                      ...message.metadata,
                      fromAgent: event.fromAgent || event.author,
                      toAgent: event.toAgent,
                      eventType: event.type || "inter_agent",
                    },
                  });
                  console.log(
                    "use-agent-connection: Added inter-agent message:",
                    JSON.stringify({
                      targetAgentId,
                      fromAgent: event.fromAgent || event.author,
                      toAgent: event.toAgent,
                      content: message.content,
                    })
                  );
                } else {
                  console.warn(
                    "use-agent-connection: No valid target agent for inter-agent message:",
                    event
                  );
                }
              } else if (
                event.object === "message" &&
                (event.fromAgent || event.author)
              ) {
                // Handle regular messages from other agents
                const message = messageApi.convertEventToMessage(event);
                const sourceAgent = event.fromAgent || event.author;
                if (sourceAgent && agents.some((a) => a.id === sourceAgent)) {
                  addMessage(sourceAgent, message);
                  console.log(
                    "use-agent-connection: Added regular message from agent:",
                    JSON.stringify(sourceAgent)
                  );
                }
              } else {
                // Log other events for debugging (these should be handled by MessageInput)
                console.log(
                  "use-agent-connection: SSE event (should be handled by MessageInput):",
                  JSON.stringify({
                    type: event.type,
                    object: event.object,
                    partial: event.partial,
                    fromAgent: event.fromAgent,
                    author: event.author,
                  })
                );
              }
            } catch (error) {
              console.error(
                "use-agent-connection: Error processing inter-agent message:",
                error,
                event
              );
            }
          },
          onError: (error) => {
            console.error("use-agent-connection: SSE error:", error);
            setConnected(false);

            // Retry connection with exponential backoff
            if (reconnectAttempts < maxReconnectAttempts) {
              reconnectAttempts++;
              const delay = Math.min(
                1000 * Math.pow(2, reconnectAttempts),
                30000
              );
              setTimeout(connectWithRetry, delay);
            } else {
              console.error(
                "use-agent-connection: Max reconnection attempts reached"
              );
            }
          },
        });
      };

      connectWithRetry();

      return () => {
        sseService.disconnect("mainChat");
      };
    }
  }, [agents, setConnected, addInterAgentEvent, addMessage]);

  return {
    agents: agentsData || [],
    isLoading,
    error,
    isError,
    isFetching,
  };
}
