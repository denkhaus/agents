"use client"

import { useEffect, useRef, useCallback, useMemo } from 'react'
import { Message } from '@/lib/types'
import { ConnectionStatus } from '@/lib/types/streaming'
import { useStreamingManager } from '@/hooks/use-streaming-manager'
import { useChatStore } from '@/lib/store'
import { debug } from '@/lib/utils/debug'

interface MessageStreamingProps {
  agentId: string
  sessionId: string
  onMessageUpdate?: (message: Message) => void
  onStatusChange?: (status: ConnectionStatus) => void
}

export function MessageStreaming({
  agentId,
  sessionId,
  onMessageUpdate,
  onStatusChange
}: MessageStreamingProps) {
  const streamingManager = useStreamingManager()
  const { addStreamingMessage, updateStreamingMessage, finalizeMessage, getSession } = useChatStore()
  const connectionRef = useRef<string | null>(null)
  // Track streaming messages by their streaming key (invocationId)
  const streamingMessagesRef = useRef<Map<string, Message>>(new Map())
  // Track finalization timeouts for messages
  const finalizationTimeouts = useRef<Map<string, NodeJS.Timeout>>(new Map())

  // Create stable references using useRef to prevent infinite re-renders
  const stableDepsRef = useRef({
    addStreamingMessage,
    updateStreamingMessage,
    finalizeMessage,
    getSession,
    onMessageUpdate,
    onStatusChange
  })

  // Update stable references when dependencies change
  useEffect(() => {
    stableDepsRef.current = {
      addStreamingMessage,
      updateStreamingMessage,
      finalizeMessage,
      getSession,
      onMessageUpdate,
      onStatusChange
    }
  }, [addStreamingMessage, updateStreamingMessage, finalizeMessage, getSession, onMessageUpdate, onStatusChange])

  const handleIncomingMessage = useCallback((message: Message) => {
    const streamingKey = message.metadata?.streamingKey || message.metadata?.invocationId
    const {
      addStreamingMessage,
      updateStreamingMessage,
      finalizeMessage,
      getSession,
      onMessageUpdate
    } = stableDepsRef.current
    
    debug.streaming(`MessageStreaming: Incoming message for agent ${agentId}`, {
      messageId: message.id,
      streamingKey,
      hasContent: !!message.content,
      contentLength: message.content?.length || 0,
      isPartial: message.metadata?.partial,
      isDone: message.metadata?.done,
      sender: message.sender,
      contentPreview: message.content?.substring(0, 100)
    })
    
    // Handle streaming message updates
    if (message.metadata?.partial && streamingKey) {
      let existingMessage = streamingMessagesRef.current.get(streamingKey)
      
      // If not in local cache, check if it exists in store
      if (!existingMessage) {
        const session = getSession(agentId)
        const storeMessage = session?.messages.find(msg => msg.id === streamingKey)
        if (storeMessage) {
          existingMessage = storeMessage
          streamingMessagesRef.current.set(streamingKey, storeMessage)
        }
      }
      
      if (!existingMessage) {
        // New streaming message - create initial message
        const newMessage = {
          ...message,
          id: streamingKey,
          content: message.content || ''
        }
        streamingMessagesRef.current.set(streamingKey, newMessage)
        addStreamingMessage(agentId, streamingKey, newMessage.content, newMessage)
        
        debug.streaming(`MessageStreaming: Created new streaming message ${streamingKey} with content: "${newMessage.content.substring(0, 50)}..."`)
        
        // Set up finalization timeout as fallback
        const timeoutId = setTimeout(() => {
          if (streamingMessagesRef.current.has(streamingKey)) {
            debug.streaming(`Auto-finalizing message ${streamingKey} after timeout (even if marked as partial)`)
            finalizeMessage(streamingKey)
            streamingMessagesRef.current.delete(streamingKey)
            finalizationTimeouts.current.delete(streamingKey)
          }
        }, 2000) // Reduced to 2 second timeout for more responsive finalization
        
        finalizationTimeouts.current.set(streamingKey, timeoutId)
      } else {
        // Update existing streaming message by accumulating content
        const updatedContent = existingMessage.content + (message.content || '')
        const updatedMessage = {
          ...existingMessage,
          content: updatedContent,
          metadata: {
            ...existingMessage.metadata,
            // Only update metadata if it indicates progression (don't regress from finalized to partial)
            partial: existingMessage.metadata?.partial === false ? false : (message.metadata?.partial ?? false),
            done: existingMessage.metadata?.done === true ? true : (message.metadata?.done ?? false),
            chunkIndex: message.metadata?.chunkIndex
          },
          parts: message.parts || existingMessage.parts
        }
        
        streamingMessagesRef.current.set(streamingKey, updatedMessage)
        updateStreamingMessage(streamingKey, updatedContent, updatedMessage.metadata)
        
        debug.streaming(`MessageStreaming: Updated streaming message ${streamingKey}, new content length: ${updatedContent.length}`)
      }

      // Finalize message when complete
      const hasContent = message.content && message.content.trim().length > 0
      const isDoneExplicitly = message.metadata?.done === true
      const isNotExplicitlyPartial = message.metadata?.partial !== true
      const isCompletion = message.metadata?.isCompletion === true
      
      // More inclusive finalization conditions:
      // 1. Explicitly marked as done
      // 2. Completion event marked as done
      // 3. Has content and is not explicitly marked as partial
      // 4. Has content and has been accumulating for a while (fallback)
      const shouldFinalize = (
        isDoneExplicitly || 
        (isCompletion && message.metadata?.done === true) || // Handle completion properly
        // If message has content and is not explicitly marked as partial
        (hasContent && isNotExplicitlyPartial)
      )
      
      if (shouldFinalize) {
        debug.streaming(`MessageStreaming: Finalizing message ${streamingKey} (reason: isDone=${isDoneExplicitly}, isCompletion=${isCompletion}, hasContent=${hasContent}, notPartial=${isNotExplicitlyPartial})`)
        
        // Clear any pending finalization timeout
        const timeoutId = finalizationTimeouts.current.get(streamingKey)
        if (timeoutId) {
          clearTimeout(timeoutId)
          finalizationTimeouts.current.delete(streamingKey)
        }
        
        finalizeMessage(streamingKey)
        streamingMessagesRef.current.delete(streamingKey)
      }
      
      // Notify parent of message update
      const currentMessage = streamingMessagesRef.current.get(streamingKey)
      if (currentMessage) {
        onMessageUpdate?.(currentMessage)
      }
    } else if (streamingKey && !message.metadata?.partial) {
      // Handle non-partial messages with streamingKey (complete messages)
      // But only if they have content to avoid creating empty duplicate messages
      if (!message.content || message.content.trim().length === 0) {
        debug.streaming(`MessageStreaming: Skipping empty non-partial message ${streamingKey}`);
        return;
      }
      
      let existingMessage = streamingMessagesRef.current.get(streamingKey);
      
      debug.streaming(`MessageStreaming: Processing complete message ${streamingKey}`);
      
      if (!existingMessage) {
        // Check store for existing message
        const session = getSession(agentId);
        const storeMessage = session?.messages.find(msg => msg.id === streamingKey);
        if (storeMessage) {
          existingMessage = storeMessage;
        }
      }
      
      if (existingMessage) {
        // Update existing message with final content
        const finalContent = message.content || existingMessage.content;
        updateStreamingMessage(streamingKey, finalContent, {
          ...existingMessage.metadata,
          partial: false,
          done: true
        });
      } else {
        // Create new complete message
        addStreamingMessage(agentId, streamingKey, message.content || '', message);
      }
      
      finalizeMessage(streamingKey);
      streamingMessagesRef.current.delete(streamingKey);
      onMessageUpdate?.(message);
    } else {
      // Non-streaming message - add directly with unique ID
      const messageId = message.id || `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
      const finalMessage = { ...message, id: messageId }
      
      debug.streaming(`MessageStreaming: Adding direct message ${messageId}`)
      
      addStreamingMessage(agentId, messageId, finalMessage.content, finalMessage)
      finalizeMessage(messageId)
      onMessageUpdate?.(finalMessage)
    }
  }, [agentId]) // Only agentId as dependency

  // Create a stable key for agentId and sessionId to prevent unnecessary effect runs
  const agentSessionKey = useMemo(() => `${agentId}-${sessionId}`, [agentId, sessionId]);

  useEffect(() => {
    if (!agentId || !sessionId) return

    debug.critical(`Establishing message streaming for agent: ${agentId}, session: ${sessionId}`)

    // Check if we already have an active connection for this agent/session
    const connectionId = `agent-${agentId}-${sessionId}`
    const connectionStatus = streamingManager.getConnectionStatus()
    const existingConnection = connectionStatus[connectionId]
    
    if (existingConnection?.isConnected) {
      debug.connection(`Reusing existing connection for ${agentId}-${sessionId}`)
      // Just update the handlers without creating a new connection
      connectionRef.current = connectionId
      return
    }

    // Establish agent connection
    const connection = streamingManager.establishAgentConnection(agentId, sessionId)
    
    if (connection) {
      connectionRef.current = connection.id

      // Set up message listener with stable handler
      const unsubscribeMessage = streamingManager.onMessage((message: Message) => {
        debug.streaming(`MessageStreaming: Message received for agent ${agentId}`, {
          messageId: message.id,
          messageSender: message.sender,
          messageType: message.type,
          hasContent: !!message.content,
          contentLength: message.content?.length || 0,
          isPartial: message.metadata?.partial,
          streamingKey: message.metadata?.streamingKey,
          invocationId: message.metadata?.invocationId
          // contentPreview: message.content?.substring(0, 100)  // Remove verbose content preview
        })
        
        // Check if this message belongs to our agent
        if (message.sender === agentId || message.sender === 'user') {
          debug.streaming(`MessageStreaming: Processing message for our agent ${agentId}`)
          handleIncomingMessage(message)
        } 
        // else {
        //   debug.streaming(`MessageStreaming: Skipping message - sender ${message.sender} does not match agent ${agentId}`)
        // }
      })

      // Set up connection status listener
      const unsubscribeStatus = streamingManager.onConnectionChange((status: ConnectionStatus) => {
        debug.connection(`Connection status changed for agent ${agentId}:`, status);
        stableDepsRef.current.onStatusChange?.(status);
        
        // If connection is closing, finalize all pending streaming messages
        if (!status.isConnected && connectionRef.current === status.connectionId) {
          debug.streaming(`Connection closed for agent ${agentId}, finalizing all pending messages`);
          
          // Finalize all pending messages
          streamingMessagesRef.current.forEach((message, streamingKey) => {
            debug.streaming(`Auto-finalizing message ${streamingKey} due to connection close`);
            stableDepsRef.current.finalizeMessage(streamingKey);
          });
          
          // Clear all finalization timeouts
          finalizationTimeouts.current.forEach(timeoutId => clearTimeout(timeoutId));
          finalizationTimeouts.current.clear();
          
          // Clear all streaming message references
          streamingMessagesRef.current.clear();
        }
      });

      // Set up error listener
      const unsubscribeError = streamingManager.onError((error: Error) => {
        debug.error(`Connection error for agent ${agentId}:`, error);
        
        // On error, also finalize all pending messages to prevent hanging partial messages
        streamingMessagesRef.current.forEach((message, streamingKey) => {
          debug.streaming(`Auto-finalizing message ${streamingKey} due to connection error`);
          stableDepsRef.current.finalizeMessage(streamingKey);
        });
        
        // Clear all finalization timeouts
        finalizationTimeouts.current.forEach(timeoutId => clearTimeout(timeoutId));
        finalizationTimeouts.current.clear();
        
        // Clear all streaming message references
        streamingMessagesRef.current.clear();
      });

      // Cleanup function
      return () => {
        debug.critical(`Cleaning up message streaming for agent: ${agentId}`)
        unsubscribeMessage()
        unsubscribeStatus()
        unsubscribeError()
        
        // Finalize all pending messages on cleanup
        streamingMessagesRef.current.forEach((message, streamingKey) => {
          debug.streaming(`Auto-finalizing message ${streamingKey} due to component cleanup`);
          stableDepsRef.current.finalizeMessage(streamingKey);
        });
        
        // Clear all finalization timeouts
        finalizationTimeouts.current.forEach(timeoutId => clearTimeout(timeoutId))
        finalizationTimeouts.current.clear()
        
        // Clear all streaming message references
        streamingMessagesRef.current.clear()
        
        // Only close connection if it's not being used by another component
        if (connectionRef.current) {
          // Don't close the connection immediately, let the connection manager handle reuse
          connectionRef.current = null
        }
      }
    }
  }, [agentSessionKey]) // Only depend on the stable key, not on streamingManager or callbacks

  // This is a logic-only component, no UI
  return null
}