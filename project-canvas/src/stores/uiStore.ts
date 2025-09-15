/**
 * UI Store - Zustand
 * Manages UI state and user interface interactions
 */

import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import { 
  WorkspaceType, 
  ViewportState, 
  SelectionState, 
  LayoutOptions, 
  ThemeConfig, 
  SidebarConfig, 
  CanvasConfig,
  NotificationState,
  UUID 
} from '../types/ui.types';

interface UIStore {
  // Sidebar State
  sidebar: SidebarConfig;
  currentWorkspace: WorkspaceType;
  
  // Theme State
  theme: ThemeConfig;
  
  // Canvas State
  canvas: CanvasConfig;
  viewport: ViewportState;
  selection: SelectionState;
  layout: LayoutOptions;
  
  // Notifications
  notifications: NotificationState[];
  
  // Loading States
  isLayouting: boolean;
  isLoading: boolean;
  
  // Sidebar Actions
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  setWorkspace: (workspace: WorkspaceType) => void;
  
  // Theme Actions
  toggleDarkMode: () => void;
  setTheme: (theme: Partial<ThemeConfig>) => void;
  
  // Canvas Actions
  setViewport: (viewport: Partial<ViewportState>) => void;
  resetViewport: () => void;
  setCanvasConfig: (config: Partial<CanvasConfig>) => void;
  
  // Selection Actions
  setSelectedNodes: (nodes: UUID[]) => void;
  setSelectedEdges: (edges: UUID[]) => void;
  addToSelection: (nodeIds: UUID[], edgeIds?: UUID[]) => void;
  removeFromSelection: (nodeIds: UUID[], edgeIds?: UUID[]) => void;
  clearSelection: () => void;
  toggleMultiSelect: (enabled: boolean) => void;
  
  // Layout Actions
  setLayoutOptions: (options: Partial<LayoutOptions>) => void;
  setIsLayouting: (isLayouting: boolean) => void;
  
  // Notification Actions
  addNotification: (notification: Omit<NotificationState, 'id' | 'timestamp'>) => void;
  removeNotification: (id: UUID) => void;
  clearNotifications: () => void;
  
  // Loading Actions
  setIsLoading: (isLoading: boolean) => void;
}

const defaultSidebar: SidebarConfig = {
  collapsed: false,
  width: 320,
  collapsedWidth: 64
};

const defaultTheme: ThemeConfig = {
  mode: 'light',
  primaryColor: '#3b82f6',
  accentColor: '#10b981'
};

const defaultCanvas: CanvasConfig = {
  minZoom: 0.1,
  maxZoom: 2.0,
  defaultZoom: 1.0,
  snapToGrid: false,
  gridSize: 20,
  showGrid: true
};

const defaultViewport: ViewportState = {
  x: 0,
  y: 0,
  zoom: 1.0
};

const defaultSelection: SelectionState = {
  selectedNodes: [],
  selectedEdges: [],
  multiSelect: false
};

const defaultLayout: LayoutOptions = {
  direction: 'LR',
  nodeSpacing: 100,
  rankSpacing: 150,
  edgeSpacing: 50,
  animate: true
};

export const useUIStore = create<UIStore>()(
  devtools(
    persist(
      (set, get) => ({
        // Initial State
        sidebar: defaultSidebar,
        currentWorkspace: 'projects',
        theme: defaultTheme,
        canvas: defaultCanvas,
        viewport: defaultViewport,
        selection: defaultSelection,
        layout: defaultLayout,
        notifications: [],
        isLayouting: false,
        isLoading: false,
        
        // Sidebar Actions
        toggleSidebar: () => set((state) => ({
          sidebar: { ...state.sidebar, collapsed: !state.sidebar.collapsed }
        })),
        
        setSidebarCollapsed: (collapsed) => set((state) => ({
          sidebar: { ...state.sidebar, collapsed }
        })),
        
        setWorkspace: (workspace) => set({ currentWorkspace: workspace }),
        
        // Theme Actions
        toggleDarkMode: () => set((state) => ({
          theme: { 
            ...state.theme, 
            mode: state.theme.mode === 'light' ? 'dark' : 'light' 
          }
        })),
        
        setTheme: (newTheme) => set((state) => ({
          theme: { ...state.theme, ...newTheme }
        })),
        
        // Canvas Actions
        setViewport: (newViewport) => set((state) => ({
          viewport: { ...state.viewport, ...newViewport }
        })),
        
        resetViewport: () => set({ viewport: defaultViewport }),
        
        setCanvasConfig: (config) => set((state) => ({
          canvas: { ...state.canvas, ...config }
        })),
        
        // Selection Actions
        setSelectedNodes: (nodes) => set((state) => ({
          selection: { ...state.selection, selectedNodes: nodes }
        })),
        
        setSelectedEdges: (edges) => set((state) => ({
          selection: { ...state.selection, selectedEdges: edges }
        })),
        
        addToSelection: (nodeIds, edgeIds = []) => set((state) => ({
          selection: {
            ...state.selection,
            selectedNodes: [...new Set([...state.selection.selectedNodes, ...nodeIds])],
            selectedEdges: [...new Set([...state.selection.selectedEdges, ...edgeIds])]
          }
        })),
        
        removeFromSelection: (nodeIds, edgeIds = []) => set((state) => ({
          selection: {
            ...state.selection,
            selectedNodes: state.selection.selectedNodes.filter(id => !nodeIds.includes(id)),
            selectedEdges: state.selection.selectedEdges.filter(id => !edgeIds.includes(id))
          }
        })),
        
        clearSelection: () => set({ selection: defaultSelection }),
        
        toggleMultiSelect: (enabled) => set((state) => ({
          selection: { ...state.selection, multiSelect: enabled }
        })),
        
        // Layout Actions
        setLayoutOptions: (options) => set((state) => ({
          layout: { ...state.layout, ...options }
        })),
        
        setIsLayouting: (isLayouting) => set({ isLayouting }),
        
        // Notification Actions
        addNotification: (notification) => {
          const id = crypto.randomUUID();
          const newNotification: NotificationState = {
            ...notification,
            id,
            timestamp: new Date()
          };
          
          set((state) => ({
            notifications: [...state.notifications, newNotification]
          }));
          
          // Auto-remove notification after duration
          if (notification.duration) {
            setTimeout(() => {
              get().removeNotification(id);
            }, notification.duration);
          }
        },
        
        removeNotification: (id) => set((state) => ({
          notifications: state.notifications.filter(n => n.id !== id)
        })),
        
        clearNotifications: () => set({ notifications: [] }),
        
        // Loading Actions
        setIsLoading: (isLoading) => set({ isLoading })
      }),
      {
        name: 'ui-store',
        // Only persist certain UI preferences
        partialize: (state) => ({
          sidebar: state.sidebar,
          theme: state.theme,
          canvas: state.canvas,
          layout: state.layout
        })
      }
    ),
    { name: 'ui-store' }
  )
);