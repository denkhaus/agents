/**
 * Auto-layout utilities for ReactFlow
 * Implements intelligent positioning to avoid overlaps
 */

import { Node, Edge, Position } from 'reactflow';
import { Task } from '../types/task.types';
import { CustomNode, CustomEdge } from '../types/reactflow.types';

export interface LayoutOptions {
  direction: 'TB' | 'LR' | 'BT' | 'RL';
  nodeSpacing: number;
  rankSpacing: number;
  edgeSpacing: number;
  animate: boolean;
}

export const defaultLayoutOptions: LayoutOptions = {
  direction: 'TB',
  nodeSpacing: 150,
  rankSpacing: 200,
  edgeSpacing: 50,
  animate: true
};

/**
 * Calculate optimal layout for tasks with dependencies
 */
export function calculateTaskLayout(
  tasks: Task[],
  options: LayoutOptions = defaultLayoutOptions
): { nodes: CustomNode[]; edges: CustomEdge[] } {
  
  // Create a dependency graph
  const dependencyGraph = buildDependencyGraph(tasks);
  
  // Perform topological sort to determine levels
  const levels = topologicalSort(dependencyGraph);
  
  // Calculate positions based on levels
  const positions = calculatePositions(levels, options);
  
  // Create nodes with calculated positions
  const nodes = createNodesWithPositions(tasks, positions);
  
  // Create edges for dependencies
  const edges = createDependencyEdges(tasks);
  
  return { nodes, edges };
}

/**
 * Build dependency graph from tasks
 */
function buildDependencyGraph(tasks: Task[]): Map<string, string[]> {
  const graph = new Map<string, string[]>();
  
  tasks.forEach(task => {
    graph.set(task.id, task.dependencies);
  });
  
  return graph;
}

/**
 * Topological sort to determine task levels
 */
function topologicalSort(graph: Map<string, string[]>): Map<string, number> {
  const levels = new Map<string, number>();
  const visited = new Set<string>();
  const visiting = new Set<string>();
  
  function visit(nodeId: string): number {
    if (visiting.has(nodeId)) {
      throw new Error(`Circular dependency detected involving task ${nodeId}`);
    }
    
    if (visited.has(nodeId)) {
      return levels.get(nodeId) || 0;
    }
    
    visiting.add(nodeId);
    
    const dependencies = graph.get(nodeId) || [];
    let maxLevel = 0;
    
    dependencies.forEach(depId => {
      if (graph.has(depId)) {
        maxLevel = Math.max(maxLevel, visit(depId) + 1);
      }
    });
    
    visiting.delete(nodeId);
    visited.add(nodeId);
    levels.set(nodeId, maxLevel);
    
    return maxLevel;
  }
  
  // Visit all nodes
  graph.forEach((_, nodeId) => {
    if (!visited.has(nodeId)) {
      visit(nodeId);
    }
  });
  
  return levels;
}

/**
 * Calculate positions based on levels
 */
function calculatePositions(
  levels: Map<string, number>,
  options: LayoutOptions
): Map<string, { x: number; y: number }> {
  const positions = new Map<string, { x: number; y: number }>();
  
  // Group tasks by level
  const tasksByLevel = new Map<number, string[]>();
  levels.forEach((level, taskId) => {
    if (!tasksByLevel.has(level)) {
      tasksByLevel.set(level, []);
    }
    tasksByLevel.get(level)!.push(taskId);
  });
  
  // Calculate positions for each level
  tasksByLevel.forEach((taskIds, level) => {
    const tasksInLevel = taskIds.length;
    const startOffset = -(tasksInLevel - 1) * options.nodeSpacing / 2;
    
    taskIds.forEach((taskId, index) => {
      let x: number, y: number;
      
      switch (options.direction) {
        case 'TB': // Top to Bottom
          x = startOffset + index * options.nodeSpacing;
          y = level * options.rankSpacing;
          break;
        case 'LR': // Left to Right
          x = level * options.rankSpacing;
          y = startOffset + index * options.nodeSpacing;
          break;
        case 'BT': // Bottom to Top
          x = startOffset + index * options.nodeSpacing;
          y = -level * options.rankSpacing;
          break;
        case 'RL': // Right to Left
          x = -level * options.rankSpacing;
          y = startOffset + index * options.nodeSpacing;
          break;
      }
      
      positions.set(taskId, { x, y });
    });
  });
  
  return positions;
}

/**
 * Create nodes with calculated positions
 */
function createNodesWithPositions(
  tasks: Task[],
  positions: Map<string, { x: number; y: number }>
): CustomNode[] {
  return tasks.map(task => {
    const position = positions.get(task.id) || { x: 0, y: 0 };
    
    return {
      id: task.id,
      type: 'task',
      position,
      data: {
        task,
        isSelected: false,
        isHighlighted: false,
        showDetails: false
      }
    } as CustomNode;
  });
}

/**
 * Create dependency edges
 */
function createDependencyEdges(tasks: Task[]): CustomEdge[] {
  const edges: CustomEdge[] = [];
  
  tasks.forEach(task => {
    task.dependencies.forEach(depId => {
      edges.push({
        id: `${depId}-${task.id}`,
        source: depId,
        target: task.id,
        type: 'dependency',
        data: {
          sourceTaskId: depId,
          targetTaskId: task.id,
          isBlocking: task.state === 'blocked',
          dependencyType: 'finish-to-start'
        }
      } as CustomEdge);
    });
  });
  
  return edges;
}

/**
 * Detect and resolve node overlaps
 */
export function resolveOverlaps(nodes: CustomNode[]): CustomNode[] {
  const resolvedNodes = [...nodes];
  const nodeSize = { width: 200, height: 100 }; // Approximate node size
  const minDistance = 50; // Minimum distance between nodes
  
  // Simple overlap resolution using force-based approach
  for (let iteration = 0; iteration < 10; iteration++) {
    let hasOverlap = false;
    
    for (let i = 0; i < resolvedNodes.length; i++) {
      for (let j = i + 1; j < resolvedNodes.length; j++) {
        const nodeA = resolvedNodes[i];
        const nodeB = resolvedNodes[j];
        
        const dx = nodeB.position.x - nodeA.position.x;
        const dy = nodeB.position.y - nodeA.position.y;
        const distance = Math.sqrt(dx * dx + dy * dy);
        
        const requiredDistance = nodeSize.width + minDistance;
        
        if (distance < requiredDistance) {
          hasOverlap = true;
          
          // Calculate repulsion force
          const force = (requiredDistance - distance) / 2;
          const angle = Math.atan2(dy, dx);
          
          const forceX = Math.cos(angle) * force;
          const forceY = Math.sin(angle) * force;
          
          // Apply force to both nodes
          nodeA.position.x -= forceX;
          nodeA.position.y -= forceY;
          nodeB.position.x += forceX;
          nodeB.position.y += forceY;
        }
      }
    }
    
    if (!hasOverlap) break;
  }
  
  return resolvedNodes;
}