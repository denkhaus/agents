"use client";

import { useChatStore } from "@/lib/store";
import { MessageItem } from "./message-item";
import { useEffect, useRef, useMemo } from "react";
import { Message } from "@/lib/types";

interface MessageListProps {
  agentId: string;
}

export function MessageList({ agentId }: MessageListProps) {
  const { sessions, currentSessionId, isLoadingMessages } = useChatStore();
  const session =
    sessions[agentId]?.sessionId === currentSessionId
      ? sessions[agentId]
      : undefined;
  const scrollAreaRef = useRef<HTMLDivElement>(null);

  // Memoize the messages content hash to prevent unnecessary re-renders
  const messagesContentHash = useMemo(() => {
    if (!session?.messages || session.messages.length === 0) return '';
    return session.messages.map((m: Message) => m.content).join("");
  }, [session?.messages]);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    const scrollToBottom = () => {
      const currentScrollArea = scrollAreaRef.current;
      if (currentScrollArea) {
        // Use requestAnimationFrame to ensure DOM has updated
        requestAnimationFrame(() => {
          currentScrollArea.scrollTop = currentScrollArea.scrollHeight;
        });
      }
    };

    scrollToBottom();
  }, [session?.messages?.length]);

  // Also scroll when message content changes (for streaming)
  useEffect(() => {
    const scrollToBottom = () => {
      const currentScrollArea = scrollAreaRef.current;
      if (currentScrollArea) {
        requestAnimationFrame(() => {
          currentScrollArea.scrollTop = currentScrollArea.scrollHeight;
        });
      }
    };

    // Scroll when any message content changes
    if (session?.messages && session.messages.length > 0) {
      scrollToBottom();
    }
  }, [messagesContentHash, session?.messages?.length]);

  if (isLoadingMessages) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="text-center text-muted-foreground">
          <p>Loading messages...</p>
        </div>
      </div>
    );
  }

  if (!session || session.messages.length === 0) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="text-center text-muted-foreground">
          <p>No messages yet</p>
          <p className="text-sm">Start a conversation with this agent</p>
        </div>
      </div>
    );
  }

  return (
    <div ref={scrollAreaRef} className="flex-1 p-4 overflow-y-auto">
      <div className="space-y-4">
        {session.messages.map((message: Message) => (
          <MessageItem key={message.id} message={message} />
        ))}
      </div>
    </div>
  );
}
