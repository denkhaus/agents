import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import type { 
  WorkspaceType, 
  FilterState, 
  SortConfig, 
  ViewConfig, 
  SelectionState, 
  DragDropState,
  NotificationState,
  ContextMenuState,
  ThemeConfig,
  Notification
} from '../types';

interface UIState {
  // Workspace and navigation
  currentWorkspace: WorkspaceType;
  sidebarCollapsed: boolean;
  
  // View configuration
  viewConfig: ViewConfig;
  sortConfig: SortConfig;
  filterState: FilterState;
  
  // Selection and interaction
  selectionState: SelectionState;
  dragDropState: DragDropState;
  
  // UI state
  notificationState: NotificationState;
  contextMenuState: ContextMenuState;
  themeConfig: ThemeConfig;
  
  // Modal and dialog states
  isCreateProjectModalOpen: boolean;
  isCreateTaskModalOpen: boolean;
  isEditProjectModalOpen: boolean;
  isEditTaskModalOpen: boolean;
  isBulkEditModalOpen: boolean;
  isDeleteConfirmModalOpen: boolean;
  
  // Chat state
  isChatPanelOpen: boolean;
  activeChatEntityId: string | null;
  activeChatEntityType: 'project' | 'task' | null;
  
  // Actions
  setCurrentWorkspace: (workspace: WorkspaceType) => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  toggleSidebar: () => void;
  
  setViewConfig: (config: Partial<ViewConfig>) => void;
  setSortConfig: (config: SortConfig) => void;
  setFilterState: (filter: Partial<FilterState>) => void;
  clearFilters: () => void;
  
  setSelectedProject: (projectId: string | undefined) => void;
  setSelectedTasks: (taskIds: string[]) => void;
  addSelectedTask: (taskId: string) => void;
  removeSelectedTask: (taskId: string) => void;
  clearSelection: () => void;
  setExpandedTasks: (taskIds: string[]) => void;
  toggleTaskExpanded: (taskId: string) => void;
  
  setDragDropState: (state: Partial<DragDropState>) => void;
  clearDragDropState: () => void;
  
  addNotification: (notification: Omit<Notification, 'id'>) => void;
  removeNotification: (id: string) => void;
  markNotificationRead: (id: string) => void;
  clearAllNotifications: () => void;
  
  setContextMenu: (state: Partial<ContextMenuState>) => void;
  closeContextMenu: () => void;
  
  setThemeConfig: (config: Partial<ThemeConfig>) => void;
  
  // Modal actions
  openCreateProjectModal: () => void;
  closeCreateProjectModal: () => void;
  openCreateTaskModal: () => void;
  closeCreateTaskModal: () => void;
  openEditProjectModal: () => void;
  closeEditProjectModal: () => void;
  openEditTaskModal: () => void;
  closeEditTaskModal: () => void;
  openBulkEditModal: () => void;
  closeBulkEditModal: () => void;
  openDeleteConfirmModal: () => void;
  closeDeleteConfirmModal: () => void;
  
  // Chat actions
  openChatPanel: (entityType: 'project' | 'task', entityId: string) => void;
  closeChatPanel: () => void;
}

const initialFilterState: FilterState = {
  search: '',
  states: [],
  complexityRange: [1, 10],
  depthRange: [0, 5],
  assignedAgent: undefined,
  showOnlyRootTasks: false,
};

const initialViewConfig: ViewConfig = {
  type: 'kanban',
  groupBy: 'state',
  showCompleted: true,
  compactMode: false,
};

const initialSortConfig: SortConfig = {
  field: 'updated_at',
  direction: 'desc',
};

const initialSelectionState: SelectionState = {
  selectedProjectId: undefined,
  selectedTaskIds: new Set(),
  selectedKanbanColumn: undefined,
  expandedTaskIds: new Set(),
};

const initialDragDropState: DragDropState = {
  isDragging: false,
  draggedItem: undefined,
  dropTarget: undefined,
};

const initialNotificationState: NotificationState = {
  notifications: [],
  unreadCount: 0,
};

const initialContextMenuState: ContextMenuState = {
  isOpen: false,
  position: { x: 0, y: 0 },
  target: undefined,
  items: [],
};

const initialThemeConfig: ThemeConfig = {
  mode: 'system',
  primaryColor: '#3b82f6',
  fontSize: 'medium',
  compactMode: false,
};

export const useUIStore = create<UIState>()(
  devtools(
    (set, get) => ({
      // Initial state
      currentWorkspace: 'projects',
      sidebarCollapsed: false,
      viewConfig: initialViewConfig,
      sortConfig: initialSortConfig,
      filterState: initialFilterState,
      selectionState: initialSelectionState,
      dragDropState: initialDragDropState,
      notificationState: initialNotificationState,
      contextMenuState: initialContextMenuState,
      themeConfig: initialThemeConfig,
      
      // Modal states
      isCreateProjectModalOpen: false,
      isCreateTaskModalOpen: false,
      isEditProjectModalOpen: false,
      isEditTaskModalOpen: false,
      isBulkEditModalOpen: false,
      isDeleteConfirmModalOpen: false,
      
      // Chat state
      isChatPanelOpen: false,
      activeChatEntityId: null,
      activeChatEntityType: null,

      // Actions
      setCurrentWorkspace: (workspace) => set({ currentWorkspace: workspace }),
      
      setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
      
      toggleSidebar: () => 
        set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),

      setViewConfig: (config) =>
        set((state) => ({
          viewConfig: { ...state.viewConfig, ...config }
        })),

      setSortConfig: (config) => set({ sortConfig: config }),

      setFilterState: (filter) =>
        set((state) => ({
          filterState: { ...state.filterState, ...filter }
        })),

      clearFilters: () => set({ filterState: initialFilterState }),

      setSelectedProject: (projectId) =>
        set((state) => ({
          selectionState: { 
            ...state.selectionState, 
            selectedProjectId: projectId 
          }
        })),

      setSelectedTasks: (taskIds) =>
        set((state) => ({
          selectionState: { 
            ...state.selectionState, 
            selectedTaskIds: new Set(taskIds) 
          }
        })),

      addSelectedTask: (taskId) =>
        set((state) => {
          const newSet = new Set(state.selectionState.selectedTaskIds);
          newSet.add(taskId);
          return {
            selectionState: { 
              ...state.selectionState, 
              selectedTaskIds: newSet 
            }
          };
        }),

      removeSelectedTask: (taskId) =>
        set((state) => {
          const newSet = new Set(state.selectionState.selectedTaskIds);
          newSet.delete(taskId);
          return {
            selectionState: { 
              ...state.selectionState, 
              selectedTaskIds: newSet 
            }
          };
        }),

      clearSelection: () =>
        set((state) => ({
          selectionState: { 
            ...state.selectionState, 
            selectedTaskIds: new Set(),
            selectedKanbanColumn: undefined
          }
        })),

      setExpandedTasks: (taskIds) =>
        set((state) => ({
          selectionState: { 
            ...state.selectionState, 
            expandedTaskIds: new Set(taskIds) 
          }
        })),

      toggleTaskExpanded: (taskId) =>
        set((state) => {
          const newSet = new Set(state.selectionState.expandedTaskIds);
          if (newSet.has(taskId)) {
            newSet.delete(taskId);
          } else {
            newSet.add(taskId);
          }
          return {
            selectionState: { 
              ...state.selectionState, 
              expandedTaskIds: newSet 
            }
          };
        }),

      setDragDropState: (dragState) =>
        set((state) => ({
          dragDropState: { ...state.dragDropState, ...dragState }
        })),

      clearDragDropState: () => set({ dragDropState: initialDragDropState }),

      addNotification: (notification) =>
        set((state) => {
          const id = Date.now().toString();
          const newNotification = { ...notification, id, read: false };
          return {
            notificationState: {
              notifications: [newNotification, ...state.notificationState.notifications],
              unreadCount: state.notificationState.unreadCount + 1
            }
          };
        }),

      removeNotification: (id) =>
        set((state) => {
          const notification = state.notificationState.notifications.find(n => n.id === id);
          const wasUnread = notification && !notification.read;
          return {
            notificationState: {
              notifications: state.notificationState.notifications.filter(n => n.id !== id),
              unreadCount: wasUnread 
                ? state.notificationState.unreadCount - 1 
                : state.notificationState.unreadCount
            }
          };
        }),

      markNotificationRead: (id) =>
        set((state) => {
          const notification = state.notificationState.notifications.find(n => n.id === id);
          const wasUnread = notification && !notification.read;
          return {
            notificationState: {
              notifications: state.notificationState.notifications.map(n =>
                n.id === id ? { ...n, read: true } : n
              ),
              unreadCount: wasUnread 
                ? state.notificationState.unreadCount - 1 
                : state.notificationState.unreadCount
            }
          };
        }),

      clearAllNotifications: () =>
        set({ notificationState: initialNotificationState }),

      setContextMenu: (contextState) =>
        set((state) => ({
          contextMenuState: { ...state.contextMenuState, ...contextState }
        })),

      closeContextMenu: () => set({ contextMenuState: initialContextMenuState }),

      setThemeConfig: (config) =>
        set((state) => ({
          themeConfig: { ...state.themeConfig, ...config }
        })),

      // Modal actions
      openCreateProjectModal: () => set({ isCreateProjectModalOpen: true }),
      closeCreateProjectModal: () => set({ isCreateProjectModalOpen: false }),
      openCreateTaskModal: () => set({ isCreateTaskModalOpen: true }),
      closeCreateTaskModal: () => set({ isCreateTaskModalOpen: false }),
      openEditProjectModal: () => set({ isEditProjectModalOpen: true }),
      closeEditProjectModal: () => set({ isEditProjectModalOpen: false }),
      openEditTaskModal: () => set({ isEditTaskModalOpen: true }),
      closeEditTaskModal: () => set({ isEditTaskModalOpen: false }),
      openBulkEditModal: () => set({ isBulkEditModalOpen: true }),
      closeBulkEditModal: () => set({ isBulkEditModalOpen: false }),
      openDeleteConfirmModal: () => set({ isDeleteConfirmModalOpen: true }),
      closeDeleteConfirmModal: () => set({ isDeleteConfirmModalOpen: false }),

      // Chat actions
      openChatPanel: (entityType, entityId) =>
        set({
          isChatPanelOpen: true,
          activeChatEntityType: entityType,
          activeChatEntityId: entityId
        }),

      closeChatPanel: () =>
        set({
          isChatPanelOpen: false,
          activeChatEntityType: null,
          activeChatEntityId: null
        }),
    }),
    { name: 'ui-store' }
  )
);