'use client'

import { useWorkspaceStore } from '@/lib/store'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { MessageSquare, Settings } from 'lucide-react'
import { cn } from '@/lib/utils'

export function WorkspaceSidebar() {
  const { activeWorkspace, setActiveWorkspace, sidebarOpen } = useWorkspaceStore()

  if (!sidebarOpen) {
    return null
  }

  return (
    <aside className="hidden border-r bg-muted/40 md:block">
      <div className="flex h-full max-h-screen flex-col gap-2">
        <div className="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
          <div className="flex items-center gap-2 font-semibold">
            <MessageSquare className="h-6 w-6" />
            <span>Workspaces</span>
          </div>
        </div>
        <div className="flex-1">
          <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
            <Button
              variant={activeWorkspace === 'chat' ? 'default' : 'ghost'}
              className={cn(
                'w-full justify-start gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary',
                activeWorkspace === 'chat' 
                  ? 'bg-muted text-primary' 
                  : 'text-muted-foreground'
              )}
              onClick={() => setActiveWorkspace('chat')}
            >
              <MessageSquare className="h-4 w-4" />
              Chat
            </Button>
            
            <Separator className="my-2" />
            
            <Button
              variant={activeWorkspace === 'settings' ? 'default' : 'ghost'}
              className={cn(
                'w-full justify-start gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary',
                activeWorkspace === 'settings' 
                  ? 'bg-muted text-primary' 
                  : 'text-muted-foreground'
              )}
              onClick={() => setActiveWorkspace('settings')}
            >
              <Settings className="h-4 w-4" />
              Settings
            </Button>
          </nav>
        </div>
      </div>
    </aside>
  )
}