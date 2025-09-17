/**
 * Central type exports
 * Re-exports all types for easy importing
 */

// Project types
export type {
  UUID,
  Project,
  ProjectProgress,
  ProjectFilter,
  ProjectEditableFields,
} from "./project.types";

// Task types
export type {
  Task,
  TaskState,
  TaskFilter,
  TaskEditableFields,
  TaskUIUpdates,
} from "./task.types";

// Agent types
export type {
  Agent,
  AgentRole,
  AgentStatus,
  AgentFilter,
  AgentUpdateInput,
} from "./agent.types";

// UI types
export type {
  WorkspaceType,
  Position,
  Dimensions,
  ViewportState,
  SelectionState,
  LayoutOptions,
  ThemeConfig,
  SidebarConfig,
  CanvasConfig,
  NotificationState,
} from "./ui.types";

export type { UIStore } from "../stores/uiStore";

export type { ProjectStore } from "../stores/projectStore";

// ReactFlow types
export type {
  CustomNodeType,
  TaskNodeData,
  ProjectNodeData,
  MilestoneNodeData,
  TaskNode,
  ProjectNode,
  MilestoneNode,
  CustomNode,
  CustomEdgeType,
  DependencyEdgeData,
  HierarchyEdgeData,
  CustomEdge,
  FlowState,
  NodeStyleConfig,
} from "./reactflow.types";

// Property Panel types
export type {
  PropertyPanelNode,
  PropertyInfo,
  EditableNodeProperties,
  PropertyUpdateCallback,
  PropertyPanelState,
} from "./property-panel.types";
