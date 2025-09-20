/**
 * Agent-specific ReactFlow types
 */

import { Node, Edge } from "@xyflow/react";
import { Agent } from "./agent.types";

export interface AgentNodeData extends Record<string, unknown> {
  agent: Agent;
  isSelected: boolean;
  isHighlighted: boolean;
}

export interface AgentNode extends Node {
  type: "agent";
  data: AgentNodeData;
}

export interface AgentEdgeData extends Record<string, unknown> {
  type: "hierarchy" | "communication" | "collaboration";
  frequency?: number;
  protocol?: string;
}

export interface AgentEdge extends Edge {
  data?: AgentEdgeData;
}

export type AgentCustomNode = AgentNode;
export type AgentCustomEdge = AgentEdge;