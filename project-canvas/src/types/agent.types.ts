/**
 * Agent-related TypeScript types
 * For future workspace extension
 */

import { UUID } from './project.types';

export const AgentRole = {
  SUPERVISOR: "supervisor",
  PROJECT_MANAGER: "project-manager", 
  CODER: "coder",
  RESEARCHER: "researcher",
  QA_ENGINEER: "qa-engineer",
  DEVOPS: "devops",
  DESIGNER: "designer"
} as const;

export type AgentRole = typeof AgentRole[keyof typeof AgentRole];

export const AgentStatus = {
  ONLINE: "online",
  OFFLINE: "offline",
  BUSY: "busy",
  IDLE: "idle"
} as const;

export type AgentStatus = typeof AgentStatus[keyof typeof AgentStatus];

export interface Agent {
  id: UUID;
  name: string;
  role: AgentRole;
  description: string;
  status: AgentStatus;
  isStreaming: boolean;
  capabilities: string[];
  currentTasks: UUID[];
  createdAt: Date;
  updatedAt: Date;
  lastActiveAt?: Date;
}

export interface AgentFilter {
  role?: AgentRole;
  status?: AgentStatus;
  searchTerm?: string;
  hasActiveTasks?: boolean;
}

export interface AgentUpdateInput {
  name?: string;
  description?: string;
  status?: AgentStatus;
  capabilities?: string[];
}