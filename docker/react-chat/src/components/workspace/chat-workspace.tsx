'use client'

import { useChatStore } from '@/lib/store'
import { ChatInterface } from '@/components/chat/chat-interface'
import { AgentList } from '@/components/agents/agent-list'

export function ChatWorkspace() {
  const { activeAgentId, agents } = useChatStore()

  return (
    <div className="flex h-full flex-col">
      {agents.length === 0 ? (
        <div className="flex h-full items-center justify-center">
          <div className="text-center">
            <h3 className="text-lg font-semibold">No agents available</h3>
            <p className="text-muted-foreground">
              Waiting for agents to come online...
            </p>
          </div>
        </div>
      ) : (
        <div className="flex h-full">
          {!activeAgentId ? (
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center">
                <h3 className="text-lg font-semibold">Select an agent</h3>
                <p className="text-muted-foreground">
                  Choose an agent from the bottom navigation to start chatting
                </p>
              </div>
            </div>
          ) : (
            <div className="flex-1">
              <ChatInterface agentId={activeAgentId} />
            </div>
          )}
          
          <div className="w-80 border-l">
            <AgentList />
          </div>
        </div>
      )}
    </div>
  )
}