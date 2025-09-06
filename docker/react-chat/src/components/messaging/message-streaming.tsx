"use client";

import { useEffect, useRef, useCallback, useMemo } from "react";
import { Message } from "@/lib/types";
import { ConnectionStatus } from "@/lib/types/streaming";
import { useStreamingManager } from "@/hooks/use-streaming-manager";
import { useChatStore } from "@/lib/store";
import { parseStructuredThoughts } from "@/lib/parsing";

interface MessageStreamingProps {
  agentId: string;
  sessionId: string;
  onMessageUpdate?: (message: Message) => void;
  onStatusChange?: (status: ConnectionStatus) => void;
}

import { parseStructuredThoughts } from "@/lib/parsing";

export function MessageStreaming({
  agentId,
  sessionId,
  onMessageUpdate,
  onStatusChange,
}: MessageStreamingProps) {
  const streamingManager = useStreamingManager();
  const {
    addStreamingMessage,
    updateStreamingMessage,
    finalizeMessage,
    getSession,
  } = useChatStore();
  const connectionRef = useRef<string | null>(null);
  const streamingMessagesRef = useRef<Map<string, Message>>(new Map());
  const finalizationTimeouts = useRef<Map<string, NodeJS.Timeout>>(new Map());

  const stableDepsRef = useRef({
    addStreamingMessage,
    updateStreamingMessage,
    finalizeMessage,
    getSession,
    onMessageUpdate,
    onStatusChange,
  });

  useEffect(() => {
    stableDepsRef.current = {
      addStreamingMessage,
      updateStreamingMessage,
      finalizeMessage,
      getSession,
      onMessageUpdate,
      onStatusChange,
    };
  }, [
    addStreamingMessage,
    updateStreamingMessage,
    finalizeMessage,
    getSession,
    onMessageUpdate,
    onStatusChange,
  ]);

  const handleIncomingMessage = useCallback(
    (message: Message) => {
      const streamingKey =
        message.metadata?.streamingKey || message.metadata?.invocationId;
      const {
        addStreamingMessage,
        updateStreamingMessage,
        finalizeMessage,
        getSession,
        onMessageUpdate,
      } = stableDepsRef.current;

      if (message.metadata?.partial && streamingKey) {
        let existingMessage = streamingMessagesRef.current.get(streamingKey);

        if (!existingMessage) {
          const session = getSession(agentId);
          const storeMessage = session?.messages.find(
            (msg) => msg.id === streamingKey
          );
          if (storeMessage) {
            existingMessage = storeMessage;
            streamingMessagesRef.current.set(streamingKey, storeMessage);
          }
        }

        if (!existingMessage) {
          const newMessage = {
            ...message,
            id: streamingKey,
            content: message.content || "",
          };
          streamingMessagesRef.current.set(streamingKey, newMessage);
          addStreamingMessage(
            agentId,
            streamingKey,
            newMessage.content,
            newMessage
          );

          const timeoutId = setTimeout(() => {
            if (streamingMessagesRef.current.has(streamingKey)) {
              finalizeMessage(streamingKey);
              streamingMessagesRef.current.delete(streamingKey);
              finalizationTimeouts.current.delete(streamingKey);
            }
          }, 2000);
          finalizationTimeouts.current.set(streamingKey, timeoutId);
        } else {
          const updatedContent =
            existingMessage.content + (message.content || "");
          const structuredParts = parseStructuredThoughts(updatedContent);

          const hasStructuredThoughts = structuredParts.length > 0;

          const updatedMessage = {
            ...existingMessage,
            content: updatedContent,
            metadata: {
              ...existingMessage.metadata,
              partial:
                existingMessage.metadata?.partial === false
                  ? false
                  : message.metadata?.partial ?? false,
              done:
                existingMessage.metadata?.done === true
                  ? true
                  : message.metadata?.done ?? false,
              chunkIndex: message.metadata?.chunkIndex,
              hasStructuredThoughts:
                hasStructuredThoughts ||
                existingMessage.metadata?.hasStructuredThoughts,
            },
            parts: hasStructuredThoughts
              ? structuredParts
              : message.parts || existingMessage.parts,
          };

          streamingMessagesRef.current.set(streamingKey, updatedMessage);
          updateStreamingMessage(
            streamingKey,
            updatedContent,
            updatedMessage.metadata,
            updatedMessage.parts
          );
        }

        const hasContent = message.content && message.content.trim().length > 0;
        const isDoneExplicitly = message.metadata?.done === true;
        const isNotExplicitlyPartial = message.metadata?.partial !== true;
        const isCompletion = message.metadata?.isCompletion === true;

        const shouldFinalize =
          isDoneExplicitly ||
          (isCompletion && message.metadata?.done === true) ||
          (hasContent && isNotExplicitlyPartial);

        if (shouldFinalize) {
          const timeoutId = finalizationTimeouts.current.get(streamingKey);
          if (timeoutId) {
            clearTimeout(timeoutId);
            finalizationTimeouts.current.delete(streamingKey);
          }
          finalizeMessage(streamingKey);
          streamingMessagesRef.current.delete(streamingKey);
        }

        const currentMessage = streamingMessagesRef.current.get(streamingKey);
        if (currentMessage) {
          onMessageUpdate?.(currentMessage);
        }
      } else if (streamingKey && !message.metadata?.partial) {
        if (!message.content || message.content.trim().length === 0) {
          return;
        }

        let existingMessage = streamingMessagesRef.current.get(streamingKey);

        if (!existingMessage) {
          const session = getSession(agentId);
          const storeMessage = session?.messages.find(
            (msg) => msg.id === streamingKey
          );
          if (storeMessage) {
            existingMessage = storeMessage;
            streamingMessagesRef.current.set(streamingKey, storeMessage);
          }
        }

        // If message is already finalized, don't process it again
        if (existingMessage && existingMessage.metadata?.done === true) {
          return;
        }

        const finalContent =
          (existingMessage?.content || "") + (message.content || "");
        const structuredParts = parseStructuredThoughts(finalContent);
        const hasStructuredThoughts = structuredParts.length > 0;

        if (existingMessage) {
          updateStreamingMessage(
            streamingKey,
            finalContent,
            {
              ...existingMessage.metadata,
              partial: false,
              done: true,
              hasStructuredThoughts:
                hasStructuredThoughts ||
                existingMessage.metadata?.hasStructuredThoughts,
            },
            hasStructuredThoughts ? structuredParts : message.parts
          );
        } else {
          const newMessage = {
            ...message,
            content: finalContent,
            parts: hasStructuredThoughts ? structuredParts : message.parts,
            metadata: { ...message.metadata, hasStructuredThoughts },
          };
          addStreamingMessage(agentId, streamingKey, finalContent, newMessage);
        }

        finalizeMessage(streamingKey);
        streamingMessagesRef.current.delete(streamingKey);
        // Don't call onMessageUpdate here - the finalized message is already in the store
      } else {
        const messageId =
          message.id ||
          `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        const finalMessage = { ...message, id: messageId };

        addStreamingMessage(
          agentId,
          messageId,
          finalMessage.content,
          finalMessage
        );
        finalizeMessage(messageId);
        onMessageUpdate?.(finalMessage);
      }
    },
    [agentId]
  );

  const agentSessionKey = useMemo(
    () => `${agentId}-${sessionId}`,
    [agentId, sessionId]
  );

  useEffect(() => {
    if (!agentId || !sessionId) return;

    const connection = streamingManager.establishAgentConnection(
      agentId,
      sessionId
    );

    if (connection) {
      connectionRef.current = connection.id;

      const unsubscribeMessage = streamingManager.onMessage(
        (message: Message) => {
          if (message.sender === agentId || message.sender === "user") {
            handleIncomingMessage(message);
          }
        }
      );

      const unsubscribeStatus = streamingManager.onConnectionChange(
        (status: ConnectionStatus) => {
          stableDepsRef.current.onStatusChange?.(status);
          if (
            !status.isConnected &&
            connectionRef.current === status.connectionId
          ) {
            streamingMessagesRef.current.forEach((_message, streamingKey) => {
              stableDepsRef.current.finalizeMessage(streamingKey);
            });
            finalizationTimeouts.current.forEach((timeoutId) =>
              clearTimeout(timeoutId)
            );
            finalizationTimeouts.current.clear();
            streamingMessagesRef.current.clear();
          }
        }
      );

      const unsubscribeError = streamingManager.onError((error: Error) => {
        streamingMessagesRef.current.forEach((_message, streamingKey) => {
          stableDepsRef.current.finalizeMessage(streamingKey);
        });
        finalizationTimeouts.current.forEach((timeoutId) =>
          clearTimeout(timeoutId)
        );
        finalizationTimeouts.current.clear();
        streamingMessagesRef.current.clear();
      });

      return () => {
        unsubscribeMessage();
        unsubscribeStatus();
        unsubscribeError();

        streamingMessagesRef.current.forEach((_message, streamingKey) => {
          stableDepsRef.current.finalizeMessage(streamingKey);
        });

        finalizationTimeouts.current.forEach((timeoutId) =>
          clearTimeout(timeoutId)
        );
        finalizationTimeouts.current.clear();
        streamingMessagesRef.current.clear();

        if (connectionRef.current) {
          connectionRef.current = null;
        }
      };
    }
  }, [agentSessionKey, streamingManager]);

  return null;
}
