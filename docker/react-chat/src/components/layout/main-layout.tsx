'use client'

import { useWorkspaceStore, useChatStore } from '@/lib/store'
import { WorkspaceSidebar } from '@/components/workspace/workspace-sidebar'
import { TopNavigation } from '@/components/navigation/top-navigation'
import { BottomNavigation } from '@/components/navigation/bottom-navigation'
import { ChatWorkspace } from '@/components/workspace/chat-workspace'
import { SettingsPanel } from '@/components/workspace/settings-panel'
import { ApiMonitor } from '@/components/debug/api-monitor'
import { useEffect } from 'react'

export function MainLayout() {
  const { activeWorkspace } = useWorkspaceStore()
  const { agents } = useChatStore()

  useEffect(() => {
    console.log('MainLayout: activeWorkspace:', activeWorkspace)
    console.log('MainLayout: agents count:', agents.length)
  }, [activeWorkspace, agents])

  const renderWorkspace = () => {
    switch (activeWorkspace) {
      case 'chat':
        return <ChatWorkspace />
      case 'settings':
        return <SettingsPanel />
      default:
        return <ChatWorkspace />
    }
  }

  const shouldShowBottomNav = activeWorkspace === 'chat'
  console.log('MainLayout: shouldShowBottomNav:', shouldShowBottomNav)

  return (
    <div className="h-screen flex flex-col">
      <TopNavigation />
      
      <div className="flex flex-1 overflow-hidden">
        <WorkspaceSidebar />
        
        <main className="flex-1 flex flex-col">
          <div className="flex-1 overflow-hidden">
            {renderWorkspace()}
          </div>
          
          {shouldShowBottomNav && (
            <div>
              <BottomNavigation />
            </div>
          )}
        </main>
      </div>
      
      <ApiMonitor />
    </div>
  )
}