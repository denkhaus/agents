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
}

export const defaultLayoutOptions: LayoutOptions = {
  direction: "TB",
  nodeSpacing: 70,   // Increased spacing between nodes
  rankSpacing: 100,  // Increased spacing between ranks (rows)
  edgeSpacing: 20,   // Increased spacing between edges
  animate: true,
};

/**
 * Calculate optimal layout for tasks with dependencies using Dagre
 */
export function calculateDagreLayout(
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

  // Create nodes with calculated positions
  const nodes: CustomNode[] = [];
  g.nodes().forEach((nodeId: string) => {
    const node = g.node(nodeId);
    const task = taskMap.get(nodeId);
    
    if (task) {
      nodes.push({
        id: nodeId,
        type: "task",
        position: { x: node.x - node.width / 2, y: node.y - node.height / 2 },
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