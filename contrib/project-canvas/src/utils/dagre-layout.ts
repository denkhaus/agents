/**
 * Dagre-based layout utilities for ReactFlow
 * Implements intelligent positioning using the Dagre graph layout library
 */

import { Task } from "../types/task.types";
import { CustomNode, CustomEdge } from "../types/reactflow.types";
import { Project } from "../types/project.types"; // Import Project type
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
  nodeSpacing: 70, // Increased spacing between nodes
  rankSpacing: 100, // Increased spacing between ranks (rows)
  edgeSpacing: 20, // Increased spacing between edges
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
  return calculateLayout(null, tasks, options);
}

/**
 * Helper function to initialize a Dagre graph
 */
function createGraph(options: LayoutOptions): dagre.graphlib.Graph {
  const g = new dagre.graphlib.Graph();
  g.setGraph({
    rankdir: options.direction,
    nodesep: options.nodeSpacing,
    ranksep: options.rankSpacing,
    edgesep: options.edgeSpacing,
    marginx: 50,
    marginy: 50,
  });
  g.setDefaultEdgeLabel(() => ({}));
  return g;
}

/**
 * Calculate layout for nodes and edges using Dagre, optionally including a project node.
 */
export function calculateLayout(
  project: Project | null,
  tasks: Task[],
  options: LayoutOptions = defaultLayoutOptions
): { nodes: CustomNode[]; edges: CustomEdge[] } {
  const g = createGraph(options);

  // Add project node if provided
  if (project) {
    const projectNodeWidth = 400;
    const projectNodeHeight = 200;
    g.setNode(project.id, {
      width: projectNodeWidth,
      height: projectNodeHeight,
    });
  }

  // Add task nodes to the graph
  tasks.forEach((task) => {
    const nodeWidth = 320;
    const nodeHeight = 160;
    g.setNode(task.id, {
      width: nodeWidth,
      height: nodeHeight,
    });
  });

  // Add edges to the graph
  if (tasks.length > 0) {
    if (project) {
      // Add hierarchy edges from project to root tasks (tasks with no parent)
      const rootTasks = tasks.filter((task) => task.parentId === undefined);
      rootTasks.forEach((task) => {
        g.setEdge(project.id, task.id);
      });
    }

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
  tasks.forEach((task) => taskMap.set(task.id, task));

  // Add project node if provided
  if (project) {
    const projectDagreNode = g.node(project.id);
    const useProjectDagrePosition =
      options.force ||
      project.positionX === undefined ||
      project.positionY === undefined;
    const projectPosition = useProjectDagrePosition
      ? {
          x: projectDagreNode.x - projectDagreNode.width / 2,
          y: projectDagreNode.y - projectDagreNode.height / 2,
        }
      : { x: project.positionX!, y: project.positionY! }; // Use non-null assertion as we checked for undefined

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
  }

  // Add task nodes
  g.nodes().forEach((nodeId: string) => {
    if (project && nodeId === project.id) return; // Skip project node if present
    const dagreNode = g.node(nodeId);
    const task = taskMap.get(nodeId);

    if (task) {
      const useDagrePosition = options.force || !task.position;
      const position = useDagrePosition
        ? {
            x: dagreNode.x - dagreNode.width / 2,
            y: dagreNode.y - dagreNode.height / 2,
          }
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

  if (project && tasks.length > 0) {
    const rootTasks = tasks.filter((task) => task.parentId === undefined);
    rootTasks.forEach((task) => {
      edges.push({
        id: `${project.id}-${task.id}`,
        source: project.id,
        target: task.id,
        type: "hierarchy",
        data: {
          parentTaskId: project.id,
          childTaskId: task.id,
          depth: 1,
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
        type: "hierarchy",
        data: {
          parentTaskId: task.parentId,
          childTaskId: task.id,
          depth: 1,
        },
      } as CustomEdge);
    }
  });

  // Add task dependency edges
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
