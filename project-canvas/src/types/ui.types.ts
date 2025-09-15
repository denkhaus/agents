/**
 * UI-specific TypeScript types
 * For React components and state management
 */

import { UUID } from './project.types';

export type WorkspaceType = 'projects' | 'agents' | 'settings';

export interface Position {
  x: number;
  y: number;
}

export interface Dimensions {
  width: number;
  height: number;
}

export interface ViewportState {
  x: number;
  y: number;
  zoom: number;
}

export interface SelectionState {
  selectedNodes: UUID[];
  selectedEdges: UUID[];
  multiSelect: boolean;
}

export interface LayoutOptions {
  direction: 'TB' | 'LR' | 'BT' | 'RL';
  nodeSpacing: number;
  rankSpacing: number;
  edgeSpacing: number;
  animate: boolean;
}

export interface ThemeConfig {
  mode: 'light' | 'dark';
  primaryColor: string;
  accentColor: string;
}

export interface SidebarConfig {
  collapsed: boolean;
  width: number;
  collapsedWidth: number;
}

export interface CanvasConfig {
  minZoom: number;
  maxZoom: number;
  defaultZoom: number;
  snapToGrid: boolean;
  gridSize: number;
  showGrid: boolean;
}

export interface NotificationState {
  id: UUID;
  type: 'success' | 'error' | 'warning' | 'info';
  title: string;
  message: string;
  duration?: number;
  timestamp: Date;
}