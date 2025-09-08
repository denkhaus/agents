'use client'

import type { LLMEvent } from '@/lib/types'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { ArrowRight } from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'
import { AGENT_IDS, getAgentDisplayName } from '@/lib/constants/agents'

interface InterAgentEventProps {
  event: LLMEvent
}

export function InterAgentEvent({ event }: InterAgentEventProps) {
  const isInterAgentCommunication = event.type === 'inter_agent'
  
  const getEventContent = () => {
    if (event.parts && event.parts.length > 0) {
      return event.parts
        .map(part => {
          if ('content' in part) return part.content
          return ''
        })
        .filter(Boolean)
        .join(' ')
    }
    
    return 'No content'
  }

  const getEventType = () => {
    return event.type || 'unknown'
  }

  const getEventTypeVariant = () => {
    const eventType = getEventType()
    switch (eventType) {
      case 'inter_agent':
      case 'communication':
        return 'default'
      case 'heartbeat':
        return 'outline'
      case 'agent_list':
        return 'secondary'
      case 'message':
        return 'default'
      case 'tool_code':
        return 'destructive'
      default:
        return 'outline'
    }
  }

  const getEventTypeColor = () => {
    const eventType = getEventType()
    switch (eventType) {
      case 'inter_agent':
      case 'communication':
        return 'border-l-blue-500'
      case 'heartbeat':
        return 'border-l-green-500'
      case 'agent_list':
        return 'border-l-purple-500'
      case 'message':
        return 'border-l-blue-400'
      case 'tool_code':
        return 'border-l-orange-500'
      default:
        return 'border-l-gray-500'
    }
  }

  return (
    <Card className={`mb-3 border-l-4 ${getEventTypeColor()}`}>
      <CardContent className="p-3">
        <div className="flex items-start justify-between mb-2">
          <Badge variant={getEventTypeVariant()} className="text-xs">
            {getEventType().replace('_', ' ').toUpperCase()}
          </Badge>
          <span className="text-xs text-muted-foreground">
            {formatDistanceToNow(new Date(event.timestamp * 1000), { addSuffix: true })}
          </span>
        </div>

        {isInterAgentCommunication && event.inter_agent && (
          <div className="flex items-center gap-2 mb-3">
            <div className="flex items-center gap-1">
              <Avatar className="h-6 w-6">
                <AvatarFallback className="text-xs">
                  {getAgentDisplayName(event.inter_agent.from_agent).slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <span className="text-sm font-medium">
                {event.inter_agent.from_agent === AGENT_IDS.HUMAN 
                  ? "You" 
                  : getAgentDisplayName(event.inter_agent.from_agent)}
              </span>
            </div>
            
            <ArrowRight className="h-4 w-4 text-muted-foreground" />
            
            <div className="flex items-center gap-1">
              <Avatar className="h-6 w-6">
                <AvatarFallback className="text-xs">
                  {getAgentDisplayName(event.inter_agent.to_agent).slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <span className="text-sm font-medium">
                {event.inter_agent.to_agent === AGENT_IDS.HUMAN 
                  ? "You" 
                  : getAgentDisplayName(event.inter_agent.to_agent)}
              </span>
            </div>
          </div>
        )}

        <div className="text-sm">
          {getEventContent()}
        </div>


        {event.usage && (
          <div className="mt-2 pt-2 border-t text-xs text-muted-foreground">
            <div className="flex gap-4">
              {event.usage.prompt_token_count && (
                <span>Prompt: {event.usage.prompt_token_count} tokens</span>
              )}
              {event.usage.candidates_token_count && (
                <span>Response: {event.usage.candidates_token_count} tokens</span>
              )}
              {event.usage.total_token_count && (
                <span>Total: {event.usage.total_token_count} tokens</span>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}