"use client";

import { useChatStore } from "@/lib/store";
import { MessageItem } from "./message-item";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useEffect, useRef } from "react";

interface MessageListProps {
  agentId: string;
}

export function MessageList({ agentId }: MessageListProps) {
  const { getSession } = useChatStore();
  const session = getSession(agentId);
  const scrollAreaRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    const scrollToBottom = () => {
      if (scrollAreaRef.current) {
        const scrollContainer = scrollAreaRef.current.querySelector(
          "[data-radix-scroll-area-viewport]"
        );
        if (scrollContainer) {
          // Use requestAnimationFrame to ensure DOM has updated
          requestAnimationFrame(() => {
            scrollContainer.scrollTop = scrollContainer.scrollHeight;
          });
        }
      }
    };

    scrollToBottom();
  }, [session?.messages]);

  // Also scroll when message content changes (for streaming)
  useEffect(() => {
    const scrollToBottom = () => {
      if (scrollAreaRef.current) {
        const scrollContainer = scrollAreaRef.current.querySelector(
          "[data-radix-scroll-area-viewport]"
        );
        if (scrollContainer) {
          requestAnimationFrame(() => {
            scrollContainer.scrollTop = scrollContainer.scrollHeight;
          });
        }
      }
    };

    // Scroll when any message content changes
    if (session?.messages && session.messages.length > 0) {
      scrollToBottom();
    }
  }, [session?.messages?.map((m) => m.content).join("")]);

  if (!session || session.messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center text-muted-foreground">
          <p>No messages yet</p>
          <p className="text-sm">Start a conversation with this agent</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col">
      <ScrollArea ref={scrollAreaRef} className="flex-1 p-4 flex-grow">
        <div className="space-y-4 min-h-full">
          {session.messages.map((message) => (
            <MessageItem key={message.id} message={message} />
          ))}
        </div>
      </ScrollArea>
    </div>
  );
}
