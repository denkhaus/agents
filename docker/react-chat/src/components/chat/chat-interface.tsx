"use client";

import { useEffect } from "react";
import { MessageList } from "./message-list";
import { MessageInput } from "./message-input";
import { useChatStore } from "@/lib/store";
import { AgentNavigationBar } from "../navigation";

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
    <div className="flex flex-col h-full">
      <div className="border-b">
        <AgentNavigationBar />
      </div>

      <div className="overflow-auto h-[80vh]">
        <MessageList agentId={agentId} />
      </div>
      <div className="sticky bg-background border-t">
        <MessageInput agentId={agentId} />
      </div>
    </div>
  );
}
