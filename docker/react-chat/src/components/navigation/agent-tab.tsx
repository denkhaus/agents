'use client'

import { Agent } from '@/lib/types'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { TabsTrigger } from '@/components/ui/tabs'

interface AgentTabProps {
  agent: Agent
}

export function AgentTab({ agent }: AgentTabProps) {
  return (
    <TabsTrigger
      value={agent.id}
      className="flex items-center gap-2 px-3 py-2"
    >
      <div className="flex items-center gap-2 min-w-0">
        <Avatar className="h-6 w-6">
          <AvatarImage src={agent.avatar} alt={agent.name} />
          <AvatarFallback className="text-xs">
            {agent.name.slice(0, 2).toUpperCase()}
          </AvatarFallback>
        </Avatar>
        
        <span className="truncate text-sm font-medium">
          {agent.name}
        </span>
        
        <Badge 
          variant={agent.status === 'online' ? 'default' : 'secondary'}
          className="h-2 w-2 p-0 rounded-full"
        />
      </div>
    </TabsTrigger>
  )
}