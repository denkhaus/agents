'use client'

import { useState, useRef } from 'react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Send, Loader2 } from 'lucide-react'
import { useChatStore } from '@/lib/store'
import { messageApi } from '@/lib/api'
import { toast } from 'sonner'
import { Message } from '@/lib/types'

interface MessageInputProps {
  agentId: string
}

export function MessageInput({ agentId }: MessageInputProps) {
  const [message, setMessage] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const { addMessage, updateMessage } = useChatStore()
  const currentAgentMessageRef = useRef<Message | null>(null)

  const handleSend = async () => {
    if (!message.trim()) return

    const userMessage = {
      id: Date.now().toString(),
      content: message.trim(),
      timestamp: new Date(),
      sender: 'user' as const,
      type: 'user' as const
    }

    // Add user message immediately
    addMessage(agentId, userMessage)
    setMessage('')
    setIsLoading(true)

    try {
      // Generate a session ID (in a real app, this would be managed differently)
      const sessionId = `session-${agentId}-${Date.now()}`
      
      console.log('Sending message to agent:', { agentId, content: userMessage.content, sessionId })
      
      await messageApi.sendMessage(
        agentId, 
        userMessage.content, 
        sessionId,
        'user',
        (responseEvent) => {
          // Handle streaming response from agent
          console.log('Received agent response event:', responseEvent)
          
          // Handle streaming response - always update/create message
          const messageId = responseEvent.id || responseEvent.invocationId || `${agentId}-${Date.now()}`
          const newContent = responseEvent.content?.parts?.[0]?.text || responseEvent.content || ''
          
          if (!currentAgentMessageRef.current) {
            // Create new streaming message
            const newAgentMessage: Message = {
              id: messageId,
              content: newContent,
              timestamp: new Date(responseEvent.timestamp * 1000),
              sender: agentId,
              type: 'agent',
              metadata: {
                invocationId: responseEvent.invocationId,
                partial: !responseEvent.done,
                done: responseEvent.done
              }
            }
            currentAgentMessageRef.current = newAgentMessage
            addMessage(agentId, newAgentMessage)
          } else {
            // Accumulate content for streaming
            const updatedContent = currentAgentMessageRef.current.content + newContent
            
            updateMessage(agentId, currentAgentMessageRef.current.id, {
              content: updatedContent,
              metadata: {
                ...currentAgentMessageRef.current.metadata,
                partial: !responseEvent.done,
                done: responseEvent.done
              }
            })
            
            currentAgentMessageRef.current = {
              ...currentAgentMessageRef.current,
              content: updatedContent
            }
          }
          
          // Reset reference when done
          if (responseEvent.done) {
            currentAgentMessageRef.current = null
          }
        },
        (error) => {
          console.error('Error during agent communication:', error)
          toast.error('Failed to communicate with agent')
        }
      )
      
      console.log('Message sent successfully, waiting for agent response...')
      
    } catch (error) {
      console.error('Error sending message:', error)
      toast.error('Failed to send message')
    } finally {
      setIsLoading(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) {
      e.preventDefault()
      handleSend()
    }
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
  )
}