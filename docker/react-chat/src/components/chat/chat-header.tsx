'use client'

import { useChatStore } from '@/lib/store'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { AgentStatus } from '@/components/agents/agent-status'

interface ChatHeaderProps {
  agentId: string
}

export function ChatHeader({ agentId }: ChatHeaderProps) {
  const { agents } = useChatStore()
  const agent = agents.find(a => a.id === agentId)

  if (!agent) {
    return (
      <div className="border-b p-4">
        <div className="flex items-center gap-3">
          <div className="text-lg font-semibold">Agent not found</div>
        </div>
      </div>
    )
  }

  return (
    <div className="border-b p-4">
      <div className="flex items-center gap-3">
        <Avatar className="h-10 w-10">
          <AvatarImage src={agent.avatar} alt={agent.name} />
          <AvatarFallback>
            {agent.name.slice(0, 2).toUpperCase()}
          </AvatarFallback>
        </Avatar>
        
        <div className="flex-1">
          <h2 className="text-lg font-semibold">{agent.name}</h2>
          <div className="flex items-center gap-2">
            <AgentStatus status={agent.status} size="sm" />
            {agent.capabilities.length > 0 && (
              <Badge variant="outline" className="text-xs">
                {agent.capabilities.length} capabilities
              </Badge>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}