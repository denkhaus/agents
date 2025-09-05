"use client";

import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { agentApi } from "@/lib/api";
import { useChatStore } from "@/lib/store";

export function useAgentSetup() {
  const {
    setAgents,
    setConnected,
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

  // The new system doesn't need a global SSE connection
  // Individual components will handle their own connections through StreamingMessageManager
  useEffect(() => {
    if (agents.length > 0) {
      setConnected(true); // Assume connected when agents are available
      console.log("Agent setup completed, agents available:", agents.length);
    }
  }, [agents, setConnected]);

  return {
    agents,
    isLoading,
    error,
    isError,
    isFetching,
    isConnected: agents.length > 0,
  };
}