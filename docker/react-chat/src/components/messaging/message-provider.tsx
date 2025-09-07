"use client";

import { useEffect, useCallback, useMemo } from "react";
import { Message } from "@/lib/types";
import { ConnectionStatus } from "@/lib/types/streaming";
import { MessageStreaming } from "./message-streaming";
import { useInterAgentCommunication } from "@/hooks/use-inter-agent-communication";
import { useChatStore } from "@/lib/store";
import { debug } from "@/lib/utils/debug";
import { AgentId } from "@/lib/constants/agents";

interface MessageProviderProps {
  agentId: AgentId;
  sessionId: string;
  agents?: string[]; // For inter-agent communication
  children?: React.ReactNode;
  onMessageUpdate?: (message: Message) => void;
  onStatusChange?: (status: ConnectionStatus) => void;
}

export function MessageProvider({
  agentId,
  sessionId,
  agents = [],
  children,
  onMessageUpdate,
  onStatusChange,
}: MessageProviderProps) {
  const { setConnected } = useChatStore();

  // Input validation
  const validateInputs = useMemo(() => {
    const errors: string[] = [];

    if (!agentId?.trim()) errors.push("Agent ID is required");
    if (!sessionId?.trim()) errors.push("Session ID is required");
    if (!Array.isArray(agents)) errors.push("Agents must be an array");

    return {
      isValid: errors.length === 0,
      errors,
    };
  }, [agentId, sessionId, JSON.stringify(agents)]); // Use deep comparison for agents

  // Memoize agents array to prevent unnecessary re-renders in inter-agent communication
  // CRITICAL: This must be declared before useInterAgentCommunication to avoid temporal dead zone
  const agentsKey = useMemo(() => agents?.join(",") || "", [agents?.join(",")]);
  const stableAgents = useMemo(() => {
    try {
      if (!Array.isArray(agents)) {
        console.warn(
          "MessageProvider: agents prop is not an array, defaulting to empty array"
        );
        return [];
      }

      return agents.length > 1 ? agents : [];
    } catch (error) {
      console.error("MessageProvider: Error processing agents array:", error);
      return [];
    }
  }, [agentsKey]);

  // Set up inter-agent communication if multiple agents
  const interAgentComm = useInterAgentCommunication({
    agents: stableAgents,
    sessionId,
    userId: "user",
  });

  useEffect(() => {
    // Log validation errors in development
    if (!validateInputs.isValid) {
      console.warn("MessageProvider validation errors:", validateInputs.errors);
    }
  }, [validateInputs]);

  useEffect(() => {
    // Update global connection status based on inter-agent connection
    setConnected(interAgentComm.isConnected);
  }, [interAgentComm.isConnected, setConnected]);

  const handleStatusChange = useCallback(
    (status: ConnectionStatus) => {
      debug.connection(`Connection status changed for ${agentId}:`, status);
      setConnected(status.isConnected);
      onStatusChange?.(status);
    },
    [agentId, setConnected, onStatusChange]
  );

  const handleMessageUpdate = useCallback(
    (message: Message) => {
      debug.streaming(`Message update for ${agentId}:`, message);
      onMessageUpdate?.(message);
    },
    [agentId, onMessageUpdate]
  );

  return (
    <>
      {/* Main agent message streaming */}
      <MessageStreaming
        agentId={agentId}
        sessionId={sessionId}
        onMessageUpdate={handleMessageUpdate}
        onStatusChange={handleStatusChange}
      />

      {/* Render children */}
      {children}
    </>
  );
}
