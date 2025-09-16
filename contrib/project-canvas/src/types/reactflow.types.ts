/**
 * ReactFlow-specific TypeScript types
 * For nodes, edges, and flow diagram components
 */

import { Node, Edge } from 'reactflow';
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
  completedTaskCount: number;
}

// Milestone Node Data
export interface MilestoneNodeData {
  title: string;
  date: Date;
  description?: string;
  isSelected: boolean;
}

// Custom Node Types
export interface TaskNode extends Node<TaskNodeData, 'task'> {
  type: 'task';
}

export interface ProjectNode extends Node<ProjectNodeData, 'project'> {
  type: 'project';
}

export interface MilestoneNode extends Node<MilestoneNodeData, 'milestone'> {
  type: 'milestone';
}

export type CustomNode = TaskNode | ProjectNode | MilestoneNode;

// Custom Edge Types
export type CustomEdgeType = 'dependency' | 'hierarchy';

// Dependency Edge Data
export interface DependencyEdgeData {
  sourceTaskId: UUID;
  targetTaskId: UUID;
  isBlocking: boolean;
  dependencyType: 'finish-to-start' | 'start-to-start' | 'finish-to-finish';
}

// Hierarchy Edge Data
export interface HierarchyEdgeData {
  parentTaskId: UUID;
  childTaskId: UUID;
  depth: number;
}

// Custom Edge Interface
export type CustomEdge = 
  | Edge<DependencyEdgeData>
  | Edge<HierarchyEdgeData>;

// Flow State
export interface FlowState {
  nodes: CustomNode[];
  edges: CustomEdge[];
  viewport: { x: number; y: number; zoom: number };
  selectedElements: { nodes: UUID[]; edges: UUID[] };
}

// Node Style Configurations
export type NodeStyleConfig = {
  [key in TaskState]: {
    backgroundColor: string;
    borderColor: string;
    textColor: string;
  };
};
