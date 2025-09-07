import { AgentId } from "../constants/agents";

export type AgentStatus = "online" | "offline" | "busy";

export interface AgentResponse {
  id: AgentId;
  name: string;
  description: string;
  role: string;
  is_streaming: boolean;
}

export interface AppInfoResponse {
  applicationName: string;
  agents: AgentResponse[];
}

export interface AgentInfo {
  name: string;
  role: string;
  id: AgentId;
  applicationName: string;
  status: AgentStatus;
  description: string;
  lastActivity?: Date;
}
