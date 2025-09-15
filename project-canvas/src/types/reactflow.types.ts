/**
 * ReactFlow-specific TypeScript types
 * For nodes, edges, and flow diagram components
 */

import { Node, Edge, NodeTypes, EdgeTypes } from 'reactflow';
import { Task, TaskState } from './task.types';
import { Project, UUID } from './project.types';

// Custom Node Types
export type CustomNodeType = 'task' | 'project' | 'milestone';

// Task Node Data
export interface TaskNodeData {
  task: Task;
  isSelected: boolean;
  isHighlighted: boolean;
  showDetails: boolean;
}

// Project Node Data  
export interface ProjectNodeData {
  project: Project;
  isSelected: boolean;
  taskCount: number;
  completionRate: number;
}

// Milestone Node Data
export interface MilestoneNodeData {
  id: UUID;
  title: string;
  description: string;
  dueDate: Date;
  isCompleted: boolean;
  dependentTasks: UUID[];
}

// Custom Node Types
export interface TaskNode extends Node {
  type: 'task';
  data: TaskNodeData;
}

export interface ProjectNode extends Node {
  type: 'project';
  data: ProjectNodeData;
}

export interface MilestoneNode extends Node {
  type: 'milestone';
  data: MilestoneNodeData;
}

export type CustomNode = TaskNode | ProjectNode | MilestoneNode;

// Custom Edge Types
export type CustomEdgeType = 'dependency' | 'hierarchy' | 'milestone';

export interface DependencyEdgeData {
  sourceTaskId: UUID;
  targetTaskId: UUID;
  isBlocking: boolean;
  dependencyType: 'finish-to-start' | 'start-to-start' | 'finish-to-finish';
}

export interface HierarchyEdgeData {
  parentTaskId: UUID;
  childTaskId: UUID;
  depth: number;
}

export interface CustomEdge extends Edge {
  type: CustomEdgeType;
  data: DependencyEdgeData | HierarchyEdgeData;
}

// Flow State
export interface FlowState {
  nodes: CustomNode[];
  edges: CustomEdge[];
  viewport: { x: number; y: number; zoom: number };
  selectedElements: { nodes: UUID[]; edges: UUID[] };
}

// Node Style Configurations
export interface NodeStyleConfig {
  [TaskState.PENDING]: {
    backgroundColor: string;
    borderColor: string;
    textColor: string;
  };
  [TaskState.IN_PROGRESS]: {
    backgroundColor: string;
    borderColor: string;
    textColor: string;
  };
  [TaskState.COMPLETED]: {
    backgroundColor: string;
    borderColor: string;
    textColor: string;
  };
  [TaskState.BLOCKED]: {
    backgroundColor: string;
    borderColor: string;
    textColor: string;
  };
  [TaskState.CANCELLED]: {
    backgroundColor: string;
    borderColor: string;
    textColor: string;
  };
}