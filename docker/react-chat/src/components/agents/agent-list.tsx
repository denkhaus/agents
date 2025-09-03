'use client'

import { useChatStore } from '@/lib/store'
import { AgentCard } from './agent-card'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'

export function AgentList() {
  const { agents, isConnected } = useChatStore()

  return (
    <div className="h-full flex flex-col">
      <div className="p-4 border-b">
        <h2 className="font-semibold">Available Agents</h2>
        <div className="flex items-center gap-2 mt-1">
          <div className={`h-2 w-2 rounded-full ${
            isConnected ? 'bg-green-500' : 'bg-red-500'
          }`} />
          <span className="text-sm text-muted-foreground">
            {isConnected ? 'Connected' : 'Disconnected'} • {agents.length} agents
          </span>
        </div>
      </div>
      
      <ScrollArea className="flex-1">
        <div className="p-4 space-y-3">
          {agents.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-muted-foreground">No agents available</p>
              <p className="text-sm text-muted-foreground mt-1">
                Agents will appear here when they come online
              </p>
            </div>
          ) : (
            agents.map((agent, index) => (
              <div key={agent.id}>
                <AgentCard agent={agent} />
                {index < agents.length - 1 && <Separator className="mt-3" />}
              </div>
            ))
          )}
        </div>
      </ScrollArea>
    </div>
  )
}