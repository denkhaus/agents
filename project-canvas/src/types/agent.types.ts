/**
 * Agent-related TypeScript types
 * For future workspace extension
 */

import { UUID } from './project.types';

export enum AgentRole {
  SUPERVISOR = "supervisor",
  PROJECT_MANAGER = "project-manager", 
  CODER = "coder",
  RESEARCHER = "researcher",
  QA_ENGINEER = "qa-engineer",
  DEVOPS = "devops",
  DESIGNER = "designer"
}

export enum AgentStatus {
  ONLINE = "online",
  OFFLINE = "offline",
  BUSY = "busy",
  IDLE = "idle"
}

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