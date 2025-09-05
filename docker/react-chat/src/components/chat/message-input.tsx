"use client";

import { useState, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Send, Loader2 } from "lucide-react";
import { useChatStore } from "@/lib/store";
import { messageApi } from "@/lib/api";
import { toast } from "sonner";
import { Message } from "@/lib/types";

interface MessageInputProps {
  agentId: string;
}

export function MessageInput({ agentId }: MessageInputProps) {
  const [message, setMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { addMessage, updateMessage, currentSessionId } = useChatStore();
  const currentAgentMessageRef = useRef<Message | null>(null);

  const handleSend = async () => {
    if (!message.trim() || !currentSessionId) return;

    const userMessage = {
      id: Date.now().toString(),
      content: message.trim(),
      timestamp: new Date(),
      sender: "user" as const,
      type: "user" as const,
    };

    // Add user message immediately
    addMessage(agentId, userMessage);
    setMessage("");
    setIsLoading(true);

    try {
      // Use the current session ID from the store
      const sessionId = currentSessionId;

      await messageApi.sendMessage(
        agentId,
        userMessage.content,
        sessionId,
        "user",
        (responseEvent) => {
          // Handle streaming response from agent

          // Handle the final stream termination event
          if (
            responseEvent.done === true &&
            !responseEvent.id &&
            !responseEvent.content
          ) {
            currentAgentMessageRef.current = null;
            return;
          }

          const messageId =
            responseEvent.id ||
            responseEvent.invocationId ||
            `${agentId}-${Date.now()}`;

          // Extract content from event
          let newContent = "";
          let parts = undefined;

          if (
            typeof responseEvent.content === "object" &&
            responseEvent.content !== null &&
            "parts" in responseEvent.content &&
            Array.isArray(responseEvent.content.parts)
          ) {
            // Handle structured content with parts
            const textParts = responseEvent.content.parts.filter(
              (part: { text?: string }) => part.text
            );
            if (textParts.length > 0) {
              newContent = textParts
                .map((part: { text?: string }) => part.text)
                .join("");
            }
            parts = responseEvent.content.parts;
          } else if (typeof responseEvent.content === "string") {
            newContent = responseEvent.content;
          }

          // Determine message type based on object type
          let messageType: "agent" | "system" = "agent";
          if (
            responseEvent.object === "tool_call" ||
            responseEvent.object === "tool_response"
          ) {
            messageType = "system";
          }

          // Check if this is a streaming continuation of the same message
          // For streaming responses, we should continue the same message if:
          // 1. We have a current message reference
          // 2. The message ID matches OR it's a partial response from the same invocation
          const isSameMessage =
            currentAgentMessageRef.current &&
            (currentAgentMessageRef.current.id === messageId ||
              (responseEvent.partial &&
                currentAgentMessageRef.current.metadata?.invocationId ===
                  responseEvent.invocationId));

          if (!isSameMessage) {
            // Create new message for new invocation (like Angular)
            const newAgentMessage: Message = {
              id: messageId,
              content: newContent,
              timestamp: new Date(
                (responseEvent.timestamp || Date.now() / 1000) * 1000
              ),
              sender: agentId,
              type: messageType,
              metadata: {
                invocationId: responseEvent.invocationId,
                partial: responseEvent.partial || false,
                done: responseEvent.done || false,
              },
              parts: parts,
            };
            currentAgentMessageRef.current = newAgentMessage;
            console.log(
              `[DEBUG] Adding NEW message: "${newContent.substring(
                0,
                50
              )}..." (partial: ${responseEvent.partial})`
            );
            addMessage(agentId, newAgentMessage);
          } else {
            // Accumulate streaming content (like Angular: streamingTextMessage.text += newChunk)
            const existingContent =
              currentAgentMessageRef.current?.content || "";
            const updatedContent = existingContent + newContent;

            // Update the reference
            currentAgentMessageRef.current = {
              ...currentAgentMessageRef.current!,
              content: updatedContent,
              metadata: {
                ...currentAgentMessageRef.current!.metadata,
                partial: responseEvent.partial || false,
                done: responseEvent.done || false,
              },
              parts: parts || currentAgentMessageRef.current!.parts,
            };

            // Update the message in the store (like Angular: streamingTextMessageSubject.next)
            console.log(
              `[DEBUG] UPDATING message: "${updatedContent.substring(
                0,
                50
              )}..." (length: ${updatedContent.length})`
            );
            updateMessage(agentId, currentAgentMessageRef.current!.id, {
              content: updatedContent,
              metadata: currentAgentMessageRef.current!.metadata,
              parts: currentAgentMessageRef.current!.parts,
            });
          }

          // Reset reference when done or connection closes
          if (responseEvent.done) {
            // Clear any pending updates and force final update
            if (updateTimeoutRef.current) {
              clearTimeout(updateTimeoutRef.current);
              updateTimeoutRef.current = null;
            }
            currentAgentMessageRef.current = null;
          }
        },
        (error) => {
          console.error("Error during agent communication:", error);
          toast.error("Failed to communicate with agent");
        }
      );
    } catch (error) {
      console.error("Error sending message:", error);
      toast.error("Failed to send message");
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) {
      e.preventDefault();
      handleSend();
    }
  };

  if (!currentSessionId) {
    return (
      <div className="border-t p-4">
        <div className="text-center text-muted-foreground">
          <p className="text-sm">No active session</p>
          <p className="text-xs">Create a new session to start chatting</p>
        </div>
      </div>
    );
  }

  return (
    <div className="border-t p-4">
      <div className="flex gap-2">
        <Textarea
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={`Type a message to ${agentId}...`}
          className="min-h-[60px] resize-none"
          disabled={isLoading}
        />
        <Button
          onClick={handleSend}
          disabled={!message.trim() || isLoading}
          size="icon"
          className="h-[60px] w-[60px]"
        >
          {isLoading ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Send className="h-4 w-4" />
          )}
        </Button>
      </div>
      <p className="text-xs text-muted-foreground mt-2">
        Press Cmd+Enter (Mac) or Ctrl+Enter (Windows) to send
      </p>
    </div>
  );
}
