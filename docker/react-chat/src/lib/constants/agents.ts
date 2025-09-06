/**
 * Agent IDs - Must match backend pkg/shared/agents.go exactly
 * These are the canonical agent identifiers used throughout the system
 */

export const AGENT_IDS = {
  HUMAN: "550e8400-e29b-41d4-a716-665544332211",
  SUPERVISOR: "550e8400-e29b-41d4-a716-446655440000", 
  CODER: "550e8400-e29b-41d4-a716-446655440001",
  DEBUGGER: "550e8400-e29b-41d4-a716-446655440002",
  PROJECT_MANAGER: "550e8400-e29b-41d4-a716-446655440003",
  RESEARCHER: "550e8400-e29b-41d4-a716-446655440004",
} as const;

export const AGENT_NAMES = {
  [AGENT_IDS.HUMAN]: "human",
  [AGENT_IDS.SUPERVISOR]: "supervisor", 
  [AGENT_IDS.CODER]: "coder",
  [AGENT_IDS.DEBUGGER]: "debugger",
  [AGENT_IDS.PROJECT_MANAGER]: "project_manager",
  [AGENT_IDS.RESEARCHER]: "researcher",
} as const;

export const AGENT_DISPLAY_NAMES = {
  [AGENT_IDS.HUMAN]: "Human",
  [AGENT_IDS.SUPERVISOR]: "Supervisor",
  [AGENT_IDS.CODER]: "Coder", 
  [AGENT_IDS.DEBUGGER]: "Debugger",
  [AGENT_IDS.PROJECT_MANAGER]: "Project Manager",
  [AGENT_IDS.RESEARCHER]: "Researcher",
} as const;

export type AgentId = typeof AGENT_IDS[keyof typeof AGENT_IDS];

/**
 * Check if an ID is a valid agent ID
 */
export function isValidAgentId(id: string): id is AgentId {
  return Object.values(AGENT_IDS).includes(id as AgentId);
}

/**
 * Get agent display name from ID
 */
export function getAgentDisplayName(id: string): string {
  if (isValidAgentId(id)) {
    return AGENT_DISPLAY_NAMES[id];
  }
  return id; // fallback to ID itself
}

/**
 * Get agent name from ID (for API calls)
 */
export function getAgentName(id: string): string {
  if (isValidAgentId(id)) {
    return AGENT_NAMES[id];
  }
  return id; // fallback to ID itself
}