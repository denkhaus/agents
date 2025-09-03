'use client'

import { Agent } from '@/lib/types'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useChatStore } from '@/lib/store'
import { MessageSquare } from 'lucide-react'

interface AgentCardProps {
  agent: Agent
}

export function AgentCard({ agent }: AgentCardProps) {
  const { activeAgentId, setActiveAgent } = useChatStore()
  const isActive = activeAgentId === agent.id

  const getStatusColor = (status: Agent['status']) => {
    switch (status) {
      case 'online':
        return 'bg-green-500'
      case 'busy':
        return 'bg-yellow-500'
      case 'offline':
        return 'bg-gray-500'
      default:
        return 'bg-gray-500'
    }
  }

  const getStatusVariant = (status: Agent['status']) => {
    switch (status) {
      case 'online':
        return 'default'
      case 'busy':
        return 'secondary'
      case 'offline':
        return 'outline'
      default:
        return 'outline'
    }
  }

  return (
    <Card className={`cursor-pointer transition-colors hover:bg-muted/50 ${
      isActive ? 'ring-2 ring-primary' : ''
    }`}>
      <CardContent className="p-4">
        <div className="flex items-center space-x-3">
          <div className="relative">
            <Avatar className="h-10 w-10">
              <AvatarImage src={agent.avatar} alt={agent.name} />
              <AvatarFallback>
                {agent.name.slice(0, 2).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div
              className={`absolute -bottom-1 -right-1 h-3 w-3 rounded-full border-2 border-background ${getStatusColor(
                agent.status
              )}`}
            />
          </div>
          
          <div className="flex-1 min-w-0">
            <div className="flex items-center justify-between">
              <h3 className="font-medium truncate">{agent.name}</h3>
              <Badge variant={getStatusVariant(agent.status)} className="text-xs">
                {agent.status}
              </Badge>
            </div>
            
            {agent.capabilities.length > 0 && (
              <p className="text-sm text-muted-foreground truncate">
                {agent.capabilities.slice(0, 2).join(', ')}
                {agent.capabilities.length > 2 && '...'}
              </p>
            )}
            
            {agent.lastActivity && (
              <p className="text-xs text-muted-foreground">
                Last active: {agent.lastActivity.toLocaleTimeString()}
              </p>
            )}
          </div>
        </div>
        
        <div className="mt-3 flex gap-2">
          <Button
            size="sm"
            variant={isActive ? 'default' : 'outline'}
            className="flex-1"
            onClick={() => setActiveAgent(agent.id)}
          >
            <MessageSquare className="h-4 w-4 mr-2" />
            {isActive ? 'Active' : 'Chat'}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}