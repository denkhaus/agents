'use client'

import { useChatStore } from '@/lib/store'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { useEffect } from 'react'

export function BottomNavigation() {
  const { agents, activeAgentId, setActiveAgent } = useChatStore()

  useEffect(() => {
    console.log('BottomNavigation: agents updated:', agents)
    console.log('BottomNavigation: activeAgentId:', activeAgentId)
  }, [agents, activeAgentId])

  console.log('BottomNavigation rendering with agents:', agents.length)

  if (agents.length === 0) {
    console.log('BottomNavigation: No agents available, not rendering')
    return (
      <div className="border-t bg-background p-2">
        <div className="text-center text-sm text-muted-foreground py-2">
          No agents available - waiting for connection...
        </div>
      </div>
    )
  }

  return (
    <div className="border-t bg-background p-2">
      <Tabs value={activeAgentId || ''} onValueChange={setActiveAgent}>
        <TabsList className="grid h-12 w-full" style={{ gridTemplateColumns: `repeat(${agents.length}, 1fr)` }}>
          {agents.map((agent) => {
            console.log('Rendering agent tab:', agent.name)
            return (
              <TabsTrigger
                key={agent.id}
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
          })}
        </TabsList>
      </Tabs>
    </div>
  )
}