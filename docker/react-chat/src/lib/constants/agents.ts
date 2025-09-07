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

/**
 * Get agent ID from name (reverse lookup)
 */
export function getAgentIdFromName(name: string): AgentId | null {
  const entry = Object.entries(AGENT_NAMES).find(([_, agentName]) => agentName === name);
  return entry ? entry[0] as AgentId : null;
}

/**
 * Check if a name is a valid agent name
 */
export function isValidAgentName(name: string): boolean {
  return Object.values(AGENT_NAMES).includes(name as any);
}

/**
 * Convert agent identifier (name or ID) to ID
 */
export function normalizeToAgentId(identifier: string): AgentId | null {
  // If it's already a valid ID, return it
  if (isValidAgentId(identifier)) {
    return identifier;
  }
  
  // If it's a valid name, convert to ID
  if (isValidAgentName(identifier)) {
    return getAgentIdFromName(identifier);
  }
  
  return null;
}

/**
 * Convert agent identifier (name or ID) to name
 */
export function normalizeToAgentName(identifier: string): string | null {
  // If it's a valid ID, convert to name
  if (isValidAgentId(identifier)) {
    return AGENT_NAMES[identifier];
  }
  
  // If it's already a valid name, return it
  if (isValidAgentName(identifier)) {
    return identifier;
  }
  
  return null;
}