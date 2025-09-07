"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Send, Loader2 } from "lucide-react";
import { useAgentsStore, useChatStore } from "@/lib/store";
import { useStreamingManager } from "@/hooks/use-streaming-manager";
import { toast } from "sonner";
import { Message, SendMessageOptions } from "@/lib/types";
import { AGENT_IDS, AgentId, normalizeToAgentId } from "@/lib/constants/agents";
import { debug } from "@/lib/utils/debug";

interface MessageInputProps {
  agentId: AgentId;
}

export function MessageInput({ agentId }: MessageInputProps) {
  const [message, setMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { addMessage, currentSessionId, applicationName } = useChatStore();

  const streamingManager = useStreamingManager();

  const agentName = normalizeToAgentId(agentId);
  const handleSend = async () => {
    if (!message.trim() || !currentSessionId) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      content: message.trim(),
      timestamp: new Date(),
      sender: AGENT_IDS.HUMAN,
      type: "user",
    };

    // Add user message immediately
    addMessage(agentId, userMessage);
    const messageContent = message.trim();
    setMessage("");
    setIsLoading(true);

    try {
      // Send message using StreamingMessageManager
      debug.log(
        `[MESSAGE INPUT] Sending message to agent ${agentName}-[${agentId}]:`,
        {
          content: messageContent.substring(0, 50) + "...",
          sessionId: currentSessionId,
          messageLength: messageContent.length,
        }
      );

      const options: SendMessageOptions = {
        onError: (error: Error) => {
          debug.error(
            `[MESSAGE INPUT] Error in message sending callback for agent ${agentId}:`,
            error
          );
          toast.error(
            `Failed to send message to ${agentName}-[${agentId}]: ${error.message}`
          );
        },
      };

      console.log(`[MESSAGE INPUT] Calling streamingManager.sendUserMessage`);
      await streamingManager.sendUserMessage(
        applicationName!,
        agentId,
        currentSessionId,
        messageContent,
        options
      );

      debug.log(`[MESSAGE INPUT] Completed streamingManager.sendUserMessage`);

      debug.log(
        `[MESSAGE INPUT] Message sent to agent ${agentName}-[${agentId}] successfully`
      );
    } catch (error) {
      debug.error("[MESSAGE INPUT] Error sending message:", {
        agentId,
        error: error instanceof Error ? error.message : error,
        stack: error instanceof Error ? error.stack : undefined,
      });
      toast.error(`Failed to send message to ${agentName}`);
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
          placeholder={`Type a message to ${agentName}...`}
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
    </div>
  );
}
