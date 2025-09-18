/**
 * Agent-related TypeScript types
 * For future workspace extension
 */

import { UUID } from "./project.types";

export const AgentRole = {
  SUPERVISOR: "supervisor",
  PROJECT_MANAGER: "project-manager",
  CODER: "coder",
  RESEARCHER: "researcher",
  QA_ENGINEER: "qa-engineer",
  DEVOPS: "devops",
  DESIGNER: "designer",
} as const;

export type AgentRoleType = (typeof AgentRole)[keyof typeof AgentRole];

export const AgentStatus = {
  ONLINE: "online",
  OFFLINE: "offline",
  BUSY: "busy",
  IDLE: "idle",
} as const;

export type AgentStatusType = (typeof AgentStatus)[keyof typeof AgentStatus];

export interface Agent {
  id: UUID;
  name: string;
  role: AgentRoleType;
  description: string;
  status: AgentStatusType;
  isStreaming: boolean;
  capabilities: string[];
  currentTasks: UUID[];
  currentTask?: string;
  efficiency?: number;
  createdAt: Date;
  updatedAt: Date;
  lastActiveAt?: Date;
}

export interface AgentFilter {
  role?: AgentRoleType;
  status?: AgentStatusType;
  searchTerm?: string;
  hasActiveTasks?: boolean;
}

export interface AgentUpdateInput {
  name?: string;
  description?: string;
  status?: AgentStatusType;
  capabilities?: string[];
  lastActiveAt?: Date;
}

// Agent Canvas/Flow specific types
export interface AgentNode {
  id: UUID;
  type: 'agent';
  position: { x: number; y: number };
  data: {
    agent: Agent;
    isSelected?: boolean;
  };
}

export interface AgentProject {
  id: UUID;
  name: string;
  description: string;
  agents: Agent[];
  agentNodes: AgentNode[];
  connections: AgentConnection[];
  createdAt: Date;
  updatedAt: Date;
}

export interface AgentConnection {
  id: UUID;
  source: UUID; // Agent ID
  target: UUID; // Agent ID
  type: 'communication' | 'hierarchy' | 'collaboration';
  label?: string;
  data?: {
    frequency?: number;
    protocol?: string;
  };
}

export interface AgentProjectFilter {
  searchTerm?: string;
  hasAgents?: boolean;
}

export interface AgentProjectUpdateInput {
  name?: string;
  description?: string;
  agents?: Agent[];
  agentNodes?: AgentNode[];
  connections?: AgentConnection[];
}
