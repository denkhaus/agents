/**
 * Agent-specific layout utilities for ReactFlow
 * Implements hierarchical auto-layout for agent networks
 */

import { Agent, AgentConnection } from "../types/agent.types";
import { AgentCustomNode, AgentCustomEdge } from "../types/agent-reactflow.types";
import dagre from "dagre";

export interface AgentLayoutOptions {
  direction: "TB" | "LR" | "BT" | "RL";
  nodeSpacing: number;
  rankSpacing: number;
  edgeSpacing: number;
  animate: boolean;
  force: boolean;
}

export const defaultAgentLayoutOptions: AgentLayoutOptions = {
  direction: "TB", // Top to Bottom for hierarchy
  nodeSpacing: 80, // Spacing between nodes in same rank
  rankSpacing: 120, // Spacing between hierarchy levels
  edgeSpacing: 30, // Spacing between edges
  animate: true,
  force: false, // Don't force layout if positions exist
};

// AgentConnection is now imported from agent.types.ts

/**
 * Calculate hierarchical layout for agents using Dagre
 */
export function calculateAgentLayout(
  agents: (Agent & { position?: { x: number; y: number } })[],
  connections: AgentConnection[] = [],
  options: AgentLayoutOptions = defaultAgentLayoutOptions
): { nodes: AgentCustomNode[]; edges: AgentCustomEdge[] } {
  const g = createAgentGraph(options);

  // Add agent nodes to the graph
  agents.forEach((agent) => {
    const nodeWidth = 280;
    const nodeHeight = 180;
    g.setNode(agent.id, {
      width: nodeWidth,
      height: nodeHeight,
    });
  });

  // Add edges based on connections
  connections.forEach((connection) => {
    if (g.hasNode(connection.source) && g.hasNode(connection.target)) {
      g.setEdge(connection.source, connection.target);
    }
  });

  // If no explicit connections, try to infer hierarchy from agent roles
  if (connections.length === 0) {
    inferAgentHierarchy(agents, g);
  }

  // Run the layout algorithm
  dagre.layout(g);

  // Create nodes with calculated positions
  const nodes: AgentCustomNode[] = [];
  const agentMap = new Map<string, Agent & { position?: { x: number; y: number } }>();
  agents.forEach((agent) => agentMap.set(agent.id, agent));

  // Add agent nodes
  g.nodes().forEach((nodeId: string) => {
    const dagreNode = g.node(nodeId);
    const agent = agentMap.get(nodeId);

    if (agent) {
      const useDagrePosition = options.force || !agent.position;
      const position = useDagrePosition
        ? {
            x: dagreNode.x - dagreNode.width / 2,
            y: dagreNode.y - dagreNode.height / 2,
          }
        : { x: agent.position!.x, y: agent.position!.y };

      nodes.push({
        id: nodeId,
        type: "agent",
        position,
        data: {
          agent,
          isSelected: false,
          isHighlighted: false,
        },
      } as AgentCustomNode);
    }
  });

  // Create edges
  const edges: AgentCustomEdge[] = connections.map((connection) => ({
    id: connection.id,
    source: connection.source,
    target: connection.target,
    type: "default",
    label: connection.label,
    data: {
      type: connection.type,
      frequency: connection.data?.frequency,
      protocol: connection.data?.protocol,
    },
    style: {
      stroke: getEdgeColor(connection.type),
      strokeWidth: 2,
    },
  }));

  return { nodes, edges };
}

/**
 * Helper function to initialize a Dagre graph for agents
 */
function createAgentGraph(options: AgentLayoutOptions): dagre.graphlib.Graph {
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
 * Infer agent hierarchy based on roles and create connections
 */
function inferAgentHierarchy(agents: Agent[], graph: dagre.graphlib.Graph): void {
  // Define role hierarchy (higher index = higher in hierarchy)
  const roleHierarchy = [
    "researcher",
    "coder", 
    "project-manager",
    "supervisor"
  ];

  // Group agents by role
  const agentsByRole = new Map<string, any[]>();
  agents.forEach((agent) => {
    const role = agent.role || "researcher";
    if (!agentsByRole.has(role)) {
      agentsByRole.set(role, []);
    }
    agentsByRole.get(role)!.push(agent);
  });

  // Create hierarchy connections
  for (let i = 0; i < roleHierarchy.length - 1; i++) {
    const lowerRole = roleHierarchy[i];
    const higherRole = roleHierarchy[i + 1];

    const lowerAgents = agentsByRole.get(lowerRole) || [];
    const higherAgents = agentsByRole.get(higherRole) || [];

    // Connect each lower role agent to the first higher role agent
    if (higherAgents.length > 0) {
      const supervisor = higherAgents[0];
      lowerAgents.forEach((agent) => {
        if (graph.hasNode(agent.id) && graph.hasNode(supervisor.id)) {
          graph.setEdge(supervisor.id, agent.id);
        }
      });
    }
  }
}

/**
 * Get edge color based on connection type
 */
function getEdgeColor(type: string): string {
  switch (type) {
    case "hierarchy":
      return "#8b5cf6"; // Purple
    case "communication":
      return "#3b82f6"; // Blue
    case "collaboration":
      return "#10b981"; // Green
    default:
      return "#6b7280"; // Gray
  }
}

/**
 * Auto-arrange agents in a grid layout (fallback when no hierarchy exists)
 */
export function calculateGridLayout(
  agents: (Agent & { position?: { x: number; y: number } })[],
  options: Partial<AgentLayoutOptions> = {}
): { nodes: AgentCustomNode[]; edges: AgentCustomEdge[] } {
  const opts = { ...defaultAgentLayoutOptions, ...options };
  const nodes: AgentCustomNode[] = [];
  
  const cols = Math.ceil(Math.sqrt(agents.length));
  const nodeWidth = 280;
  const nodeHeight = 180;
  const spacingX = nodeWidth + opts.nodeSpacing;
  const spacingY = nodeHeight + opts.rankSpacing;

  agents.forEach((agent, index) => {
    const row = Math.floor(index / cols);
    const col = index % cols;
    
    const useDagrePosition = opts.force || !agent.position;
    const position = useDagrePosition
      ? {
          x: col * spacingX,
          y: row * spacingY,
        }
      : { x: agent.position!.x, y: agent.position!.y };

    nodes.push({
      id: agent.id,
      type: "agent",
      position,
      data: {
        agent,
        isSelected: false,
        isHighlighted: false,
      },
    } as AgentCustomNode);
  });

  return { nodes, edges: [] };
}