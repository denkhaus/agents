/**
 * Store Operations Types
 * Strong typing for all store operation interfaces
 */

import { Project } from "./project.types";
import { AgentProject } from "./agent.types";
import { WorkspaceType } from "./index";

// Task Store State Interface
export interface TaskStoreState {
  tasksByProject: Record<string, any[]>; // TODO: Replace with proper Task type
  loading: boolean;
  error: string | null;
}

// Project Store State Interface
export interface ProjectStoreState {
  projects: Project[];
  currentProject: Project | null;
  loading: boolean;
  error: string | null;
}

// Agent Project Store State Interface
export interface AgentProjectStoreState {
  agentProjects: AgentProject[];
  currentAgentProject: AgentProject | null;
  loading: boolean;
  error: string | null;
  updateAgentNodePosition: (
    projectId: string,
    agentId: string,
    position: { x: number; y: number }
  ) => void;
}

// UI Store State Interface
export interface UIStoreState {
  selection: {
    selectedNodes: string[];
    selectedEdges: string[];
    multiSelect: boolean;
  };
  viewport: {
    x: number;
    y: number;
    zoom: number;
  };
  setViewport: (
    viewport: Partial<{ x: number; y: number; zoom: number }>
  ) => void;
  loading: boolean;
  error: string | null;
}

// Settings Store State Interface
export interface SettingsStoreState {
  selectedNodeIds: string[];
  currentWorkspace: WorkspaceType;
  updateSelectedNodes: (nodeIds: string[]) => void;
  updateWorkspace: (workspace: WorkspaceType) => void;
  loading: boolean;
  error: string | null;
}

// Real Time Data Interface
export interface RealTimeDataState {
  updateTaskPosition: (
    id: string,
    position: { x: number; y: number }
  ) => Promise<void>;
  updateProjectPosition: (
    id: string,
    position: { x: number; y: number }
  ) => Promise<void>;
  loading: boolean;
  error: string | null;
}

// Store Getter Functions
// Store Getter Functions - strongly typed
export type TaskStoreGetter = () => TaskStoreState;
export type ProjectStoreGetter = () => ProjectStoreState;
export type AgentProjectStoreGetter = () => AgentProjectStoreState;
export type UIGetter = () => UIStoreState;
export type SettingsGetter = () => SettingsStoreState;
export type RealTimeDataGetter = () => RealTimeDataState;

// Position type for consistency
export interface Position {
  x: number;
  y: number;
}
