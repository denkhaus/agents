/**
 * Dagre-based layout utilities for ReactFlow
 * Implements intelligent positioning using the Dagre graph layout library
 */

import { Task } from "../types/task.types";
import { CustomNode, CustomEdge } from "../types/reactflow.types";
import dagre from "dagre";

export interface LayoutOptions {
  direction: "TB" | "LR" | "BT" | "RL";
  nodeSpacing: number;
  rankSpacing: number;
  edgeSpacing: number;
  animate: boolean;
  force: boolean; // Add force option
}

export const defaultLayoutOptions: LayoutOptions = {
  direction: "TB",
  nodeSpacing: 70,   // Increased spacing between nodes
  rankSpacing: 100,  // Increased spacing between ranks (rows)
  edgeSpacing: 20,   // Increased spacing between edges
  animate: true,
  force: false, // Default to not forcing layout
};

/**
 * Calculate optimal layout for tasks with dependencies using Dagre
 */
export function calculateDagreLayout(
  tasks: Task[],
  options: LayoutOptions = defaultLayoutOptions
): { nodes: CustomNode[]; edges: CustomEdge[] } {
  return calculateTaskOnlyLayout(tasks, options);
}

/**
 * Calculate layout including project node and its tasks
 */
export function calculateProjectLayout(
  project: any, // Project type
  tasks: Task[],
  options: LayoutOptions = defaultLayoutOptions
): { nodes: CustomNode[]; edges: CustomEdge[] } {
  // Create a new Dagre graph
  const g = new dagre.graphlib.Graph();
  
  // Configure the graph layout
  g.setGraph({
    rankdir: options.direction,
    nodesep: options.nodeSpacing,
    ranksep: options.rankSpacing,
    edgesep: options.edgeSpacing,
    marginx: 50,
    marginy: 50,
  });
  
  // Set default edge label
  g.setDefaultEdgeLabel(() => ({}));

  // Add project node (larger dimensions)
  const projectNodeWidth = 400;
  const projectNodeHeight = 200;
  g.setNode(project.id, { 
    width: projectNodeWidth, 
    height: projectNodeHeight,
  });

  // Add task nodes to the graph
  tasks.forEach((task) => {
    const nodeWidth = 320;
    const nodeHeight = 160;
    
    g.setNode(task.id, { 
      width: nodeWidth, 
      height: nodeHeight,
    });
  });

  // Only add edges if we have tasks
  if (tasks.length > 0) {
    // Add hierarchy edges from project to root tasks (tasks with no parent)
    const rootTasks = tasks.filter(task => task.parentId === undefined);
    rootTasks.forEach((task) => {
      g.setEdge(project.id, task.id);
    });

    // Add parent-child hierarchy edges
    tasks.forEach((task) => {
      if (task.parentId && g.hasNode(task.parentId) && g.hasNode(task.id)) {
        g.setEdge(task.parentId, task.id);
      }
    });

    // Add task dependency edges (separate from hierarchy)
    tasks.forEach((task) => {
      task.dependencies.forEach((depId) => {
        if (g.hasNode(depId) && g.hasNode(task.id)) {
          g.setEdge(depId, task.id);
        }
      });
    });
  }

  // Run the layout algorithm
  dagre.layout(g);

  // Create nodes with calculated positions
  const nodes: CustomNode[] = [];
  const taskMap = new Map<string, Task>();
  tasks.forEach(task => taskMap.set(task.id, task));

  // Add project node
  const projectDagreNode = g.node(project.id);
  const useProjectDagrePosition = options.force || !project.positionX || !project.positionY;
  const projectPosition = useProjectDagrePosition
    ? { x: projectDagreNode.x - projectDagreNode.width / 2, y: projectDagreNode.y - projectDagreNode.height / 2 }
    : { x: project.positionX, y: project.positionY };
    
  nodes.push({
    id: project.id,
    type: "project",
    position: projectPosition,
    data: {
      project,
      isSelected: false,
      taskCount: tasks.length,
      completionRate: project.progress || 0,
      completedTaskCount: project.completedTasks || 0,
    },
  });

  // Add task nodes
  g.nodes().forEach((nodeId: string) => {
    if (nodeId === project.id) return; // Skip project node, already added
    
    const dagreNode = g.node(nodeId);
    const task = taskMap.get(nodeId);
    
    if (task) {
      const useDagrePosition = options.force || !task.position;
      const position = useDagrePosition
        ? { x: dagreNode.x - dagreNode.width / 2, y: dagreNode.y - dagreNode.height / 2 }
        : { x: task.position!.x, y: task.position!.y };

      nodes.push({
        id: nodeId,
        type: "task",
        position,
        data: {
          task,
          isSelected: false,
          isHighlighted: false,
          showDetails: false,
        },
      });
    }
  });

  // Create edges
  const edges: CustomEdge[] = [];

  // Add hierarchy edges from project to root tasks (only if tasks exist)
  if (tasks.length > 0) {
    const rootTasks = tasks.filter(task => task.parentId === undefined);
    rootTasks.forEach((task) => {
      edges.push({
        id: `${project.id}-${task.id}`,
        source: project.id,
        target: task.id,
        type: "dependency", // We can create a "hierarchy" type later if needed
        data: {
          sourceTaskId: project.id,
          targetTaskId: task.id,
          isBlocking: false,
          dependencyType: "hierarchy",
        },
      } as CustomEdge);
    });
  }

  // Add parent-child hierarchy edges
  tasks.forEach((task) => {
    if (task.parentId) {
      edges.push({
        id: `${task.parentId}-${task.id}`,
        source: task.parentId,
        target: task.id,
        type: "dependency", // We can create a "hierarchy" type later if needed
        data: {
          sourceTaskId: task.parentId,
          targetTaskId: task.id,
          isBlocking: false,
          dependencyType: "hierarchy",
        },
      } as CustomEdge);
    }
  });

  // Add task dependency edges (separate from hierarchy)
  tasks.forEach((task) => {
    task.dependencies.forEach((depId) => {
      edges.push({
        id: `${depId}-${task.id}-dep`,
        source: depId,
        target: task.id,
        type: "dependency",
        data: {
          sourceTaskId: depId,
          targetTaskId: task.id,
          isBlocking: task.state === "blocked",
          dependencyType: "finish-to-start",
        },
      } as CustomEdge);
    });
  });

  return { nodes, edges };
}

/**
 * Calculate layout for tasks only (original function)
 */
function calculateTaskOnlyLayout(
  tasks: Task[],
  options: LayoutOptions = defaultLayoutOptions
): { nodes: CustomNode[]; edges: CustomEdge[] } {
  // Create a new Dagre graph
  const g = new dagre.graphlib.Graph();
  
  // Configure the graph layout
  g.setGraph({
    rankdir: options.direction,
    nodesep: options.nodeSpacing,
    ranksep: options.rankSpacing,
    edgesep: options.edgeSpacing,
    marginx: 50,   // Increased margins
    marginy: 50,   // Increased margins
  });
  
  // Set default edge label
  g.setDefaultEdgeLabel(() => ({}));

  // Add nodes to the graph
  tasks.forEach((task) => {
    // Approximate node dimensions (these should match your actual node sizes)
    const nodeWidth = 320; // Width of task nodes
    const nodeHeight = 160; // Height of task nodes
    
    g.setNode(task.id, { 
      width: nodeWidth, 
      height: nodeHeight,
    });
  });

  // Add edges to the graph based on dependencies
  tasks.forEach((task) => {
    task.dependencies.forEach((depId) => {
      // Only add edge if both nodes exist
      if (g.hasNode(depId) && g.hasNode(task.id)) {
        g.setEdge(depId, task.id);
      }
    });
  });

  // Run the layout algorithm
  dagre.layout(g);

  // Create a map of task IDs to tasks for quick lookup
  const taskMap = new Map<string, Task>();
  tasks.forEach(task => taskMap.set(task.id, task));

  // Create nodes with calculated positions, respecting existing positions
  const nodes: CustomNode[] = [];
  g.nodes().forEach((nodeId: string) => {
    const dagreNode = g.node(nodeId);
    const task = taskMap.get(nodeId);
    
    if (task) {
      // If forcing, or if the task has no position, use Dagre.
      // Otherwise, use the existing position.
      const useDagrePosition = options.force || !task.position;
      const position = useDagrePosition
        ? { x: dagreNode.x - dagreNode.width / 2, y: dagreNode.y - dagreNode.height / 2 }
        : { x: task.position!.x, y: task.position!.y };

      nodes.push({
        id: nodeId,
        type: "task",
        position,
        data: {
          task,
          isSelected: false,
          isHighlighted: false,
          showDetails: false,
        },
      });
    }
  });

  // Create edges for dependencies
  const edges: CustomEdge[] = [];
  tasks.forEach((task) => {
    task.dependencies.forEach((depId) => {
      edges.push({
        id: `${depId}-${task.id}`,
        source: depId,
        target: task.id,
        type: "dependency",
        data: {
          sourceTaskId: depId,
          targetTaskId: task.id,
          isBlocking: task.state === "blocked",
          dependencyType: "finish-to-start",
        },
      } as CustomEdge);
    });
  });

  return { nodes, edges };
}