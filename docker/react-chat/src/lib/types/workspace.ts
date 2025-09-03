export type WorkspaceType = 'chat' | 'settings'

export interface WorkspaceState {
  activeWorkspace: WorkspaceType
  sidebarOpen: boolean
}

export interface ChatWorkspaceState {
  activeAgentId: string | null
  isConnected: boolean
}

export interface SettingsWorkspaceState {
  theme: 'light' | 'dark'
  serverUrl: string
  autoConnect: boolean
  soundEnabled: boolean
}