import { create } from 'zustand'
import { WorkspaceType } from '@/lib/types'

interface WorkspaceStore {
  // State
  activeWorkspace: WorkspaceType
  sidebarOpen: boolean
  
  // Actions
  setActiveWorkspace: (workspace: WorkspaceType) => void
  toggleSidebar: () => void
  setSidebarOpen: (open: boolean) => void
}

export const useWorkspaceStore = create<WorkspaceStore>((set) => ({
  // Initial state
  activeWorkspace: 'chat',
  sidebarOpen: true,
  
  // Actions
  setActiveWorkspace: (workspace) => set({ activeWorkspace: workspace }),
  
  toggleSidebar: () => set((state) => ({ 
    sidebarOpen: !state.sidebarOpen 
  })),
  
  setSidebarOpen: (open) => set({ sidebarOpen: open })
}))