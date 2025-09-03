'use client'

import { useEffect } from 'react'
import { ChatHeader } from './chat-header'
import { MessageList } from './message-list'
import { MessageInput } from './message-input'
import { InterAgentEventDisplay } from '@/components/inter-agent/inter-agent-event-display'
import { useChatStore } from '@/lib/store'

interface ChatInterfaceProps {
  agentId: string
}

export function ChatInterface({ agentId }: ChatInterfaceProps) {
  const { createSession } = useChatStore()

  // Ensure session exists for this agent
  useEffect(() => {
    createSession(agentId)
  }, [agentId, createSession])

  return (
    <div className="flex h-full flex-col">
      <ChatHeader agentId={agentId} />
      
      <div className="flex-1 flex">
        <div className="flex-1 flex flex-col">
          <MessageList agentId={agentId} />
          <MessageInput agentId={agentId} />
        </div>
        
        <div className="w-80 border-l">
          <InterAgentEventDisplay />
        </div>
      </div>
    </div>
  )
}