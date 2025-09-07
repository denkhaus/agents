"use client";

import { useEffect, useRef, useCallback } from "react";
import { StreamingMessageManager } from "@/lib/streaming";
import {
  MessageCallback,
  InterAgentCallback,
  ConnectionCallback,
  ErrorCallback,
  SendMessageOptions,
} from "@/lib/types/streaming";
import { AgentId } from "@/lib/constants/agents";

export function useStreamingManager() {
  const managerRef = useRef<StreamingMessageManager | null>(null);

  // Initialize the manager
  useEffect(() => {
    if (!managerRef.current) {
      managerRef.current = StreamingMessageManager.getInstance({
        maxReconnectAttempts: 5,
        reconnectInterval: 1000,
        backoffMultiplier: 2,
        connectionTimeout: 30000,
      });
    }

    return () => {
      // Cleanup on unmount
      if (managerRef.current) {
        managerRef.current.destroy();
        managerRef.current = null;
      }
    };
  }, []);

  const establishAgentConnection = useCallback(
    (agentId: AgentId, sessionId: string) => {
      return managerRef.current?.establishAgentConnection(agentId, sessionId);
    },
    []
  );

  const sendUserMessage = useCallback(
    async (
      appName: string,
      agentId: AgentId,
      sessionId: string,
      content: string,
      options?: SendMessageOptions
    ) => {
      if (!managerRef.current) {
        throw new Error("StreamingMessageManager not initialized");
      }

      return managerRef.current.sendUserMessage(
        appName,
        agentId,
        sessionId,
        content,
        options
      );
    },
    []
  );

  const closeConnection = useCallback((connectionId: string) => {
    managerRef.current?.closeConnection(connectionId);
  }, []);

  const onMessage = useCallback((callback: MessageCallback) => {
    return managerRef.current?.onMessage(callback) || (() => {});
  }, []);

  const onInterAgentEvent = useCallback((callback: InterAgentCallback) => {
    return managerRef.current?.onInterAgentEvent(callback) || (() => {});
  }, []);

  const onConnectionChange = useCallback((callback: ConnectionCallback) => {
    return managerRef.current?.onConnectionChange(callback) || (() => {});
  }, []);

  const onError = useCallback((callback: ErrorCallback) => {
    return managerRef.current?.onError(callback) || (() => {});
  }, []);

  const getConnectionStatus = useCallback(() => {
    return managerRef.current?.getConnectionStatus() || {};
  }, []);

  const isConnected = useCallback((connectionId?: string) => {
    return managerRef.current?.isConnected(connectionId) || false;
  }, []);

  return {
    establishAgentConnection,
    sendUserMessage,
    closeConnection,
    onMessage,
    onInterAgentEvent,
    onConnectionChange,
    onError,
    getConnectionStatus,
    isConnected,
  };
}
