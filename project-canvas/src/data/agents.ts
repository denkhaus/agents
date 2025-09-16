/**
 * Agents Data
 * Go-Model konform: pkg/tools/project/shared/shared.go
 */

import type { Agent } from "../types";
import { generateUUID } from "../utils/uuid";

// Import the actual values for runtime usage
import {
  AgentRole as AgentRoleValue,
  AgentStatus as AgentStatusValue,
} from "../types/agent.types";

export const masterAgents: Agent[] = [
  {
    id: generateUUID(),
    name: "Design Lead",
    role: AgentRoleValue.DESIGNER,
    description:
      "Senior UX/UI designer specializing in e-commerce platforms and design systems",
    status: AgentStatusValue.ONLINE,
    isStreaming: false,
    capabilities: [
      "wireframing",
      "prototyping",
      "design-systems",
      "user-research",
      "figma",
      "accessibility",
    ],
    currentTasks: [],
    createdAt: new Date("2024-01-10T09:00:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    lastActiveAt: new Date("2024-01-20T14:30:00Z"),
  },
  {
    id: generateUUID(),
    name: "Frontend Developer",
    role: AgentRoleValue.CODER,
    description:
      "React/TypeScript specialist with expertise in modern frontend development and performance optimization",
    status: AgentStatusValue.BUSY,
    isStreaming: true,
    capabilities: [
      "react",
      "typescript",
      "tailwind",
      "testing",
      "performance-optimization",
      "accessibility",
    ],
    currentTasks: [],
    createdAt: new Date("2024-01-10T09:15:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    lastActiveAt: new Date("2024-01-20T14:25:00Z"),
  },
  {
    id: generateUUID(),
    name: "Backend Developer",
    role: AgentRoleValue.CODER,
    description:
      "Node.js/Express expert with database design and API architecture experience",
    status: AgentStatusValue.ONLINE,
    isStreaming: false,
    capabilities: [
      "nodejs",
      "express",
      "postgresql",
      "api-design",
      "security",
      "performance",
    ],
    currentTasks: [],
    createdAt: new Date("2024-01-10T09:30:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    lastActiveAt: new Date("2024-01-20T14:20:00Z"),
  },
  {
    id: generateUUID(),
    name: "QA Engineer",
    role: AgentRoleValue.QA_ENGINEER,
    description:
      "Quality assurance specialist with automated testing expertise and mobile testing experience",
    status: AgentStatusValue.IDLE,
    isStreaming: false,
    capabilities: [
      "jest",
      "cypress",
      "test-automation",
      "performance-testing",
      "accessibility",
      "mobile-testing",
    ],
    currentTasks: [],
    createdAt: new Date("2024-01-10T10:00:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    lastActiveAt: new Date("2024-01-20T13:45:00Z"),
  },
  {
    id: generateUUID(),
    name: "DevOps Engineer",
    role: AgentRoleValue.DEVOPS,
    description:
      "Infrastructure and deployment automation specialist with cloud and mobile deployment expertise",
    status: AgentStatusValue.ONLINE,
    isStreaming: false,
    capabilities: [
      "docker",
      "kubernetes",
      "ci-cd",
      "monitoring",
      "cloud-infrastructure",
      "mobile-deployment",
    ],
    currentTasks: [],
    createdAt: new Date("2024-01-10T10:15:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    lastActiveAt: new Date("2024-01-20T14:10:00Z"),
  },
];

// Export individual agents for easy access
export const designerAgent = masterAgents[0];
export const frontendAgent = masterAgents[1];
export const backendAgent = masterAgents[2];
export const qaAgent = masterAgents[3];
export const devopsAgent = masterAgents[4];
