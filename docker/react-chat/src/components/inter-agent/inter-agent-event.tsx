'use client'

import type { AgentEvent } from '@/lib/types'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { ArrowRight } from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'
import { AGENT_IDS, getAgentDisplayName } from '@/lib/constants/agents'

interface InterAgentEventProps {
  event: AgentEvent
}

export function InterAgentEvent({ event }: InterAgentEventProps) {
  const isInterAgentCommunication = event.object === 'inter_agent' || event.type === 'communication'
  
  const getEventContent = () => {
    if (typeof event.content === 'string') {
      return event.content
    }
    
    if (event.content && typeof event.content === 'object' && event.content.parts) {
      return event.content.parts
        .map(part => part.text)
        .filter(Boolean)
        .join(' ')
    }
    
    return event.message || 'No content'
  }

  const getEventType = () => {
    // Backend uses "object" for main classification, "type" for system events
    return event.object || event.type || 'unknown'
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

        {isInterAgentCommunication && (event.fromAgent || event.interAgent?.fromAgent) && (event.toAgent || event.interAgent?.toAgent) && (
          <div className="flex items-center gap-2 mb-3">
            <div className="flex items-center gap-1">
              <Avatar className="h-6 w-6">
                <AvatarFallback className="text-xs">
                  {getAgentDisplayName(event.fromAgent || event.interAgent?.fromAgent || 'Unknown').slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <span className="text-sm font-medium">
                {event.fromAgent === AGENT_IDS.HUMAN || event.interAgent?.fromAgent === AGENT_IDS.HUMAN 
                  ? "You" 
                  : getAgentDisplayName(event.fromAgent || event.interAgent?.fromAgent || 'Unknown')}
              </span>
            </div>
            
            <ArrowRight className="h-4 w-4 text-muted-foreground" />
            
            <div className="flex items-center gap-1">
              <Avatar className="h-6 w-6">
                <AvatarFallback className="text-xs">
                  {getAgentDisplayName(event.toAgent || event.interAgent?.toAgent || 'Unknown').slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <span className="text-sm font-medium">
                {event.toAgent === AGENT_IDS.HUMAN || event.interAgent?.toAgent === AGENT_IDS.HUMAN 
                  ? "You" 
                  : getAgentDisplayName(event.toAgent || event.interAgent?.toAgent || 'Unknown')}
              </span>
            </div>
          </div>
        )}

        <div className="text-sm">
          {getEventContent()}
        </div>

        {event.type === 'agent_list' && event.agents && (
          <div className="mt-2 pt-2 border-t">
            <div className="text-xs text-muted-foreground mb-1">Available Agents:</div>
            <div className="flex flex-wrap gap-1">
              {event.agents.map((agent, index) => (
                <Badge key={index} variant="outline" className="text-xs">
                  {agent.name}
                </Badge>
              ))}
            </div>
          </div>
        )}

        {event.usageMetadata && (
          <div className="mt-2 pt-2 border-t text-xs text-muted-foreground">
            <div className="flex gap-4">
              {event.usageMetadata.promptTokenCount && (
                <span>Prompt: {event.usageMetadata.promptTokenCount} tokens</span>
              )}
              {event.usageMetadata.candidatesTokenCount && (
                <span>Response: {event.usageMetadata.candidatesTokenCount} tokens</span>
              )}
              {event.usageMetadata.totalTokenCount && (
                <span>Total: {event.usageMetadata.totalTokenCount} tokens</span>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}