"use client";

import { useEffect } from "react";
import { ChatHeader } from "./chat-header";
import { MessageList } from "./message-list";
import { MessageInput } from "./message-input";
import { useChatStore } from "@/lib/store";

interface ChatInterfaceProps {
  agentId: string;
}

export function ChatInterface({ agentId }: ChatInterfaceProps) {
  const { loadSessions } = useChatStore();

  // Load sessions when agent changes
  useEffect(() => {
    const initializeAgent = async () => {
      await loadSessions(agentId);
      // Create session will be called by setActiveAgent in the parent component
    };

    initializeAgent();
  }, [agentId, loadSessions]);

  return (
    <div className="flex-1 flex flex-col h-full">
      <ChatHeader agentId={agentId} />
      <div className="flex-1 flex">
        <div className="flex-1 flex flex-col">
          <MessageList agentId={agentId} />
          <MessageInput agentId={agentId} />
        </div>
      </div>
    </div>
  );
}
