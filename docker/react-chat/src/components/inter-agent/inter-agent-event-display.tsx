'use client'

import { useChatStore } from '@/lib/store'
import { InterAgentEvent } from './inter-agent-event'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Trash2, Activity } from 'lucide-react'

export function InterAgentEventDisplay() {
  const { interAgentEvents, clearInterAgentEvents, isConnected } = useChatStore()

  return (
    <div className="h-full flex flex-col">
      <div className="p-4 border-b flex-shrink-0">
        <div className="flex items-center justify-between mb-2">
          <h3 className="font-semibold">Inter-Agent Communication</h3>
          <Button
            variant="ghost"
            size="sm"
            onClick={clearInterAgentEvents}
            disabled={interAgentEvents.length === 0}
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
        
        <div className="flex items-center gap-2">
          <div className={`h-2 w-2 rounded-full ${
            isConnected ? 'bg-green-500' : 'bg-red-500'
          }`} />
          <span className="text-sm text-muted-foreground">
            {isConnected ? 'Connected' : 'Disconnected'}
          </span>
          <Badge variant="outline" className="ml-auto text-xs">
            {interAgentEvents.length} events
          </Badge>
        </div>
      </div>

      <div className="flex-1 min-h-0">
        <ScrollArea className="h-full">
          <div className="p-4">
            {interAgentEvents.length === 0 ? (
              <div className="text-center py-8">
                <Activity className="h-8 w-8 mx-auto text-muted-foreground mb-2" />
                <p className="text-muted-foreground">No inter-agent events yet</p>
                <p className="text-sm text-muted-foreground mt-1">
                  Communication between agents will appear here
                </p>
              </div>
            ) : (
              <div className="space-y-1">
                {interAgentEvents.map((event, index) => (
                  <InterAgentEvent
                    key={event.id || event.invocationId || `${event.timestamp}-${index}`}
                    event={event}
                  />
                ))}
              </div>
            )}
          </div>
        </ScrollArea>
      </div>
    </div>
  )
}