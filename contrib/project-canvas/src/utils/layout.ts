/**
 * Auto-layout utilities for ReactFlow
 * Implements intelligent positioning to avoid overlaps
 * 
 * NOTE: This implementation has been replaced with Dagre-based layout.
 * The function signature is maintained for backward compatibility.
 */

import { Task } from "../types/task.types";
import { CustomNode, CustomEdge } from "../types/reactflow.types";
import { calculateDagreLayout, calculateProjectLayout, LayoutOptions, defaultLayoutOptions } from "./dagre-layout";

export type { LayoutOptions };
export { defaultLayoutOptions, calculateProjectLayout };

/**
 * Calculate optimal layout for tasks with dependencies using Dagre
 * This replaces the previous custom implementation with a more robust solution
 */
export function calculateTaskLayout(
  tasks: Task[],
  options: LayoutOptions = defaultLayoutOptions
): { nodes: CustomNode[]; edges: CustomEdge[] } {
  // Use the new Dagre-based layout implementation
  return calculateDagreLayout(tasks, options);
}

/**
 * Detect and resolve node overlaps
 * 
 * NOTE: This function is now deprecated as Dagre handles positioning automatically.
 * It is kept for backward compatibility but returns nodes unchanged.
 */
export function resolveOverlaps(nodes: CustomNode[]): CustomNode[] {
  // Dagre layout already resolves overlaps, so we don't need this function anymore
  // Keeping it for backward compatibility
  return nodes;
}
