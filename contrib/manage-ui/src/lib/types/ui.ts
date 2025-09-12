// UI-specific types and interfaces

export type WorkspaceType = 'projects' | 'agents' | 'monitoring' | 'settings';

export interface WorkspaceConfig {
  id: WorkspaceType;
  title: string;
  icon: string;
  description: string;
  enabled: boolean;
}

export interface SidebarItem {
  id: string;
  icon: string;
  label: string;
  active: boolean;
  onClick: () => void;
}

export interface FilterState {
  search: string;
  states: import('./project').TaskState[];
  complexityRange: [number, number];
  depthRange: [number, number];
  assignedAgent?: string;
  showOnlyRootTasks: boolean;
}

export interface SortConfig {
  field: 'title' | 'created_at' | 'updated_at' | 'complexity' | 'state';
  direction: 'asc' | 'desc';
}

export interface ViewConfig {
  type: 'kanban' | 'list' | 'tree';
  groupBy: 'state' | 'depth' | 'assigned_agent' | 'none';
  showCompleted: boolean;
  compactMode: boolean;
}

export interface SelectionState {
  selectedProjectId?: string;
  selectedTaskIds: Set<string>;
  selectedKanbanColumn?: string;
  expandedTaskIds: Set<string>;
}

export interface DragDropState {
  isDragging: boolean;
  draggedItem?: {
    type: 'task' | 'project';
    id: string;
    data: unknown;
  };
  dropTarget?: {
    type: 'column' | 'task' | 'project';
    id: string;
    position?: 'before' | 'after' | 'inside';
  };
}

export interface NotificationState {
  notifications: Notification[];
  unreadCount: number;
}

export interface Notification {
  id: string;
  type: 'info' | 'success' | 'warning' | 'error';
  title: string;
  message: string;
  timestamp: string;
  read: boolean;
  actions?: NotificationAction[];
}

export interface NotificationAction {
  label: string;
  action: () => void;
  variant?: 'default' | 'destructive';
}

export interface ContextMenuState {
  isOpen: boolean;
  position: { x: number; y: number };
  target?: {
    type: 'project' | 'task';
    id: string;
    data: unknown;
  };
  items: ContextMenuItem[];
}

export interface ContextMenuItem {
  id: string;
  label: string;
  icon?: string;
  disabled?: boolean;
  destructive?: boolean;
  onClick: () => void;
  submenu?: ContextMenuItem[];
}

export interface KeyboardShortcut {
  key: string;
  ctrlKey?: boolean;
  shiftKey?: boolean;
  altKey?: boolean;
  action: () => void;
  description: string;
}

export interface ThemeConfig {
  mode: 'light' | 'dark' | 'system';
  primaryColor: string;
  fontSize: 'small' | 'medium' | 'large';
  compactMode: boolean;
}