"use client";

import { useEffect, useCallback, useMemo } from "react";
import { MessageList } from "./message-list";
import { MessageInput } from "./message-input";
import { MessageProvider } from "@/components/messaging";
import { useChatStore } from "@/lib/store";
import { AgentNavigationBar } from "../navigation";
import { debug } from "@/lib/utils/debug";
import { Message } from "@/lib/types";
import { ConnectionStatus } from "@/lib/types/streaming";

interface ChatInterfaceProps {
  agentId: string;
}

export function ChatInterface({ agentId }: ChatInterfaceProps) {
  const { loadSessions, currentSessionId, agents } = useChatStore();

  // Load sessions when agent changes
  useEffect(() => {
    const initializeAgent = async () => {
      await loadSessions(agentId);
      // Create session will be called by setActiveAgent in the parent component
    };

    initializeAgent();
  }, [agentId, loadSessions]);

  // Memoize agents list to prevent unnecessary re-renders
  // Use deep comparison based on agent IDs to prevent constant re-creation
  const agentIds = useMemo(() => {
    return agents.map(agent => agent.id).sort();
  }, [JSON.stringify(agents.map(agent => agent.id).sort())]); // Use deep comparison

  // Memoize callback functions to prevent infinite re-renders
  const handleMessageUpdate = useCallback((message: Message) => {
    debug.streaming('Chat Interface: Message updated', message);
  }, []);

  const handleStatusChange = useCallback((status: ConnectionStatus) => {
    debug.connection('Chat Interface: Connection status changed', status);
  }, []);

  if (!currentSessionId) {
    return (
      <div className="flex flex-col h-full">
        <div className="border-b">
          <AgentNavigationBar />
        </div>
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center text-muted-foreground">
            <p className="text-lg">No active session</p>
            <p className="text-sm">Create a new session to start chatting</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      <div className="border-b">
        <AgentNavigationBar />
      </div>

      <MessageProvider
        agentId={agentId}
        sessionId={currentSessionId}
        agents={agentIds}
        onMessageUpdate={handleMessageUpdate}
        onStatusChange={handleStatusChange}
      >
        <div className="overflow-auto h-[80vh]">
          <MessageList agentId={agentId} />
        </div>
        <div className="sticky bg-background border-t">
          <MessageInput agentId={agentId} />
        </div>
      </MessageProvider>
    </div>
  );
}