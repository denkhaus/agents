/**
 * Convex Integration Helpers
 * Utilities for working with Convex real-time data
 */

import { Task } from "../types/task.types";
import { Project, UUID } from "../types/project.types";
import type {
  Agent,
  AgentRoleType,
  AgentStatusType,
} from "../types/agent.types";
import { Doc } from "../../convex/_generated/dataModel";
import { AgentProject, AgentNode, AgentConnection } from "../types/agent.types";

// Use Convex-generated types for type safety
export type ConvexProject = Doc<"projects">;
export type ConvexTask = Doc<"tasks">;
export type ConvexAgent = Doc<"agents">;
export type ConvexEvent = Doc<"events">;
export type ConvexSettings = Doc<"settings">;
export type ConvexAgentProject = Doc<"agentProjects">;

/**
 * Convert Convex project to frontend Project type
 */
export function convexProjectToProject(convexProject: ConvexProject): Project {
  return {
    id: convexProject.id as UUID, // Use the custom UUID, not the Convex _id
    title: convexProject.title,
    description: convexProject.description,
    createdAt: new Date(convexProject.createdAt),
    updatedAt: new Date(convexProject.updatedAt),
    totalTasks: convexProject.totalTasks,
    completedTasks: convexProject.completedTasks,
    progress: convexProject.progress,
    positionX: convexProject.positionX,
    positionY: convexProject.positionY,
  };
}

/**
 * Convert Convex task to frontend Task type
 */
export function convexTaskToTask(convexTask: ConvexTask): Task {
  return {
    id: convexTask.id as UUID,
    projectId: convexTask.projectId as UUID,
    parentId: convexTask.parentId as UUID | undefined,
    title: convexTask.title,
    description: convexTask.description,
    state: convexTask.state,
    complexity: convexTask.complexity,
    depth: convexTask.depth,
    estimate: convexTask.estimate,
    assignedAgent: convexTask.assignedAgent as UUID | undefined,
    dependencies: convexTask.dependencies as UUID[],
    dependents: convexTask.dependents as UUID[],
    createdAt: new Date(convexTask.createdAt),
    updatedAt: new Date(convexTask.updatedAt),
    completedAt: convexTask.completedAt
      ? new Date(convexTask.completedAt)
      : undefined,
    position:
      convexTask.positionX !== undefined && convexTask.positionY !== undefined
        ? { x: convexTask.positionX, y: convexTask.positionY }
        : undefined,
  };
}

/**
 * Convert Convex agent to frontend Agent type
 */
export function convexAgentToAgent(convexAgent: ConvexAgent): Agent {
  // Provide defaults for required fields that might be missing
  const role = convexAgent.role || "coder"; // Default role
  const description = convexAgent.description || ""; // Default empty description
  const status = convexAgent.status || "offline"; // Default status
  const isStreaming = convexAgent.isStreaming || false; // Default streaming status
  const capabilities = convexAgent.capabilities || []; // Default empty capabilities
  const currentTasks = convexAgent.currentTasks || []; // Default empty tasks

  return {
    id: convexAgent.id as UUID,
    name: convexAgent.name,
    role: role as AgentRoleType,
    description: description,
    status: status as AgentStatusType,
    isStreaming: isStreaming,
    capabilities: capabilities,
    currentTasks: currentTasks as UUID[],
    createdAt: new Date(convexAgent._creationTime),
    updatedAt: new Date(convexAgent._creationTime),
    lastActiveAt: convexAgent.lastActiveAt
      ? new Date(convexAgent.lastActiveAt)
      : undefined,
  };
}

/**
 * Convert Convex AgentProject to frontend AgentProject type
 */
export function convexAgentProjectToAgentProject(
  convexAgentProject: ConvexAgentProject,
  agents: Agent[] // Pass agents to enrich agentNodes
): AgentProject {
  const agentNodes: AgentNode[] = convexAgentProject.agentNodes.map((node) => {
    const agent = agents.find((a) => a.id === node.id);
    return {
      id: node.id as UUID,
      type: "agent",
      position: node.position,
      data: {
        agent: agent || ({} as Agent), // Ensure agent is populated, fallback to empty Agent if not found
        isSelected: node.data?.isSelected || false,
      },
    };
  });

  const connections: AgentConnection[] = convexAgentProject.connections.map(
    (conn) => ({
      id: conn.id as UUID,
      source: conn.source as UUID,
      target: conn.target as UUID,
      type: conn.type as "communication" | "hierarchy" | "collaboration",
      label: conn.label,
      data: conn.data,
    })
  );

  return {
    id: convexAgentProject.id as UUID,
    name: convexAgentProject.name,
    description: convexAgentProject.description,
    agentNodes: agentNodes,
    connections: connections,
    createdAt: new Date(convexAgentProject.createdAt),
    updatedAt: new Date(convexAgentProject.updatedAt),
  };
}