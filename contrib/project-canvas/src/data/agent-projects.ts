/**
 * Agent Projects Data
 * Go-Model konform: AgentProject structure for seeding Convex database
 * Diese Daten dienen als Vorlage für die Agent-Canvas-Visualisierung
 */

import type { AgentProject, AgentNode, AgentConnection } from "../types";
import { AGENT_IDS } from "./agents";
import { PROJECT_IDS } from "./projects";

// Feste AgentProject-IDs für konsistente Referenzen (UUID-Format)
export const AGENT_PROJECT_IDS = {
  ECOMMERCE_TEAM: "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
  MOBILE_TEAM: "b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e",
  FULL_DEVELOPMENT_TEAM: "c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f",
};

// Agent Node Positionen für Canvas-Layout
const ECOMMERCE_AGENT_POSITIONS = [
  { x: 150, y: 100 }, // Designer
  { x: 400, y: 50 }, // Frontend Dev
  { x: 400, y: 200 }, // Backend Dev
  { x: 650, y: 125 }, // QA Engineer
  { x: 150, y: 300 }, // DevOps
];

const MOBILE_AGENT_POSITIONS = [
  { x: 200, y: 150 }, // Designer
  { x: 450, y: 100 }, // Frontend Dev (Mobile)
  { x: 450, y: 250 }, // Backend Dev
  { x: 700, y: 175 }, // QA Engineer
];

const FULL_TEAM_POSITIONS = [
  { x: 100, y: 100 }, // Designer
  { x: 300, y: 50 }, // Frontend Dev
  { x: 500, y: 50 }, // Backend Dev
  { x: 300, y: 200 }, // QA Engineer
  { x: 500, y: 200 }, // DevOps
];

// Erstelle Agent Nodes
const createAgentNodes = (
  agentIds: string[],
  positions: { x: number; y: number }[]
): AgentNode[] => {
  return agentIds.map((agentId, index) => ({
    id: agentId,
    type: "agent" as const,
    position: positions[index] || { x: 100 + index * 200, y: 100 },
    data: {
      agent: {
        id: agentId,
        // Note: Vollständige Agent-Daten werden zur Laufzeit aus agents.ts geladen
      } as any,
      isSelected: false,
    },
  }));
};

// Agent Connections für verschiedene Projekttypen
const createEcommerceConnections = (): AgentConnection[] => [
  {
    id: "d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a",
    source: AGENT_IDS.DESIGNER,
    target: AGENT_IDS.FRONTEND_DEV,
    type: "collaboration",
    label: "Design Handoff",
    data: { frequency: 8, protocol: "figma-handoff" },
  },
  {
    id: "e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b",
    source: AGENT_IDS.FRONTEND_DEV,
    target: AGENT_IDS.BACKEND_DEV,
    type: "collaboration",
    label: "API Integration",
    data: { frequency: 9, protocol: "rest-api" },
  },
  {
    id: "f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c",
    source: AGENT_IDS.BACKEND_DEV,
    target: AGENT_IDS.QA_ENGINEER,
    type: "hierarchy",
    label: "Code Review",
    data: { frequency: 7, protocol: "pull-request" },
  },
  {
    id: "a7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d",
    source: AGENT_IDS.QA_ENGINEER,
    target: AGENT_IDS.DEVOPS,
    type: "collaboration",
    label: "Testing Results",
    data: { frequency: 6, protocol: "test-reports" },
  },
  {
    id: "b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e",
    source: AGENT_IDS.DESIGNER,
    target: AGENT_IDS.QA_ENGINEER,
    type: "communication",
    label: "UX Feedback",
    data: { frequency: 4, protocol: "user-testing" },
  },
];

const createMobileConnections = (): AgentConnection[] => [
  {
    id: "c9d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f",
    source: AGENT_IDS.DESIGNER,
    target: AGENT_IDS.FRONTEND_DEV,
    type: "collaboration",
    label: "Mobile Design System",
    data: { frequency: 9, protocol: "design-tokens" },
  },
  {
    id: "d0e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a",
    source: AGENT_IDS.FRONTEND_DEV,
    target: AGENT_IDS.BACKEND_DEV,
    type: "collaboration",
    label: "Mobile API",
    data: { frequency: 8, protocol: "graphql-api" },
  },
  {
    id: "e1f2a3b4-c5d6-4e7f-8a9b-0c1d2e3f4a5b",
    source: AGENT_IDS.BACKEND_DEV,
    target: AGENT_IDS.QA_ENGINEER,
    type: "hierarchy",
    label: "Mobile Testing",
    data: { frequency: 7, protocol: "device-testing" },
  },
];

const createFullTeamConnections = (): AgentConnection[] => [
  {
    id: "f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c",
    source: AGENT_IDS.DESIGNER,
    target: AGENT_IDS.FRONTEND_DEV,
    type: "collaboration",
    label: "Design System",
    data: { frequency: 8, protocol: "design-handoff" },
  },
  {
    id: "a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d",
    source: AGENT_IDS.DESIGNER,
    target: AGENT_IDS.BACKEND_DEV,
    type: "communication",
    label: "UX Requirements",
    data: { frequency: 5, protocol: "requirements-doc" },
  },
  {
    id: "b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e",
    source: AGENT_IDS.FRONTEND_DEV,
    target: AGENT_IDS.BACKEND_DEV,
    type: "collaboration",
    label: "Full-Stack Integration",
    data: { frequency: 9, protocol: "api-contracts" },
  },
  {
    id: "c5d6e7f8-a9b0-4c1d-2e3f-4a5b6c7d8e9f",
    source: AGENT_IDS.FRONTEND_DEV,
    target: AGENT_IDS.QA_ENGINEER,
    type: "hierarchy",
    label: "Frontend Testing",
    data: { frequency: 7, protocol: "e2e-tests" },
  },
  {
    id: "d6e7f8a9-b0c1-4d2e-3f4a-5b6c7d8e9f0a",
    source: AGENT_IDS.BACKEND_DEV,
    target: AGENT_IDS.QA_ENGINEER,
    type: "hierarchy",
    label: "Backend Testing",
    data: { frequency: 6, protocol: "integration-tests" },
  },
  {
    id: "e7f8a9b0-c1d2-4e3f-4a5b-6c7d8e9f0a1b",
    source: AGENT_IDS.QA_ENGINEER,
    target: AGENT_IDS.DEVOPS,
    type: "collaboration",
    label: "Deployment Pipeline",
    data: { frequency: 5, protocol: "ci-cd" },
  },
  {
    id: "f8a9b0c1-d2e3-4f4a-5b6c-7d8e9f0a1b2c",
    source: AGENT_IDS.BACKEND_DEV,
    target: AGENT_IDS.DEVOPS,
    type: "collaboration",
    label: "Infrastructure",
    data: { frequency: 6, protocol: "infrastructure-as-code" },
  },
];

export const masterAgentProjects: AgentProject[] = [
  // E-Commerce Team Configuration
  {
    id: AGENT_PROJECT_IDS.ECOMMERCE_TEAM,
    name: "E-Commerce Development Team",
    description: `# E-Commerce Development Team

## Team Overview
Specialized team configuration for the e-commerce platform redesign project. Optimized workflow for web-focused development with strong design-to-development handoff processes.

### Team Composition
- **Designer**: UX/UI specialist for e-commerce platforms
- **Frontend Developer**: React/TypeScript expert
- **Backend Developer**: Node.js/Express specialist
- **QA Engineer**: E-commerce testing and automation
- **DevOps Engineer**: Web deployment and infrastructure

### Workflow Highlights
- Design-first approach with Figma handoffs
- Component-based frontend development
- RESTful API architecture
- Automated testing pipeline
- Containerized deployment strategy

### Communication Protocols
- Daily standups via Slack
- Design reviews in Figma
- Code reviews via GitHub
- Testing reports in Jira
- Deployment notifications in Teams`,
    agentNodes: createAgentNodes(
      [
        AGENT_IDS.DESIGNER,
        AGENT_IDS.FRONTEND_DEV,
        AGENT_IDS.BACKEND_DEV,
        AGENT_IDS.QA_ENGINEER,
        AGENT_IDS.DEVOPS,
      ],
      ECOMMERCE_AGENT_POSITIONS
    ),
    connections: createEcommerceConnections(),
    createdAt: new Date("2024-01-15T10:00:00Z"),
    updatedAt: new Date("2024-01-20T15:30:00Z"),
  },

  // Mobile Team Configuration
  {
    id: AGENT_PROJECT_IDS.MOBILE_TEAM,
    name: "Mobile Development Team",
    description: `# Mobile Development Team

## Team Overview
Lean team configuration optimized for mobile app development with React Native. Focus on cross-platform development and mobile-specific testing.

### Team Composition
- **Designer**: Mobile UX/UI specialist
- **Frontend Developer**: React Native expert
- **Backend Developer**: Mobile API specialist
- **QA Engineer**: Mobile testing and device compatibility

### Mobile-Specific Focus
- Native mobile design patterns
- Cross-platform React Native development
- Mobile-optimized APIs with offline support
- Device-specific testing across iOS/Android
- App store deployment processes

### Collaboration Model
- Design tokens for consistent mobile UI
- GraphQL APIs for efficient mobile data loading
- Automated device testing pipeline
- Beta testing via TestFlight/Play Console`,
    agentNodes: createAgentNodes(
      [
        AGENT_IDS.DESIGNER,
        AGENT_IDS.FRONTEND_DEV,
        AGENT_IDS.BACKEND_DEV,
        AGENT_IDS.QA_ENGINEER,
      ],
      MOBILE_AGENT_POSITIONS
    ),
    connections: createMobileConnections(),
    createdAt: new Date("2024-01-10T11:00:00Z"),
    updatedAt: new Date("2024-01-18T17:15:00Z"),
  },

  // Full Development Team Configuration
  {
    id: AGENT_PROJECT_IDS.FULL_DEVELOPMENT_TEAM,
    name: "Full Development Team",
    description: `# Full Development Team

## Team Overview
Complete development team configuration supporting both web and mobile development. Maximum collaboration and cross-functional capabilities.

### Team Composition
- **Designer**: Full-stack design (web + mobile)
- **Frontend Developer**: Multi-platform specialist
- **Backend Developer**: Full-stack backend architect
- **QA Engineer**: Comprehensive testing lead
- **DevOps Engineer**: Multi-environment infrastructure

### Multi-Project Capabilities
- Simultaneous web and mobile development
- Shared design system across platforms
- Unified backend services
- Cross-platform testing strategies
- Multi-environment deployment pipelines

### Advanced Collaboration
- Cross-functional pair programming
- Shared component libraries
- Integrated API contracts
- End-to-end automation
- Infrastructure as code`,
    agentNodes: createAgentNodes(
      [
        AGENT_IDS.DESIGNER,
        AGENT_IDS.FRONTEND_DEV,
        AGENT_IDS.BACKEND_DEV,
        AGENT_IDS.QA_ENGINEER,
        AGENT_IDS.DEVOPS,
      ],
      FULL_TEAM_POSITIONS
    ),
    connections: createFullTeamConnections(),
    createdAt: new Date("2024-01-08T09:00:00Z"),
    updatedAt: new Date("2024-01-22T12:45:00Z"),
  },
];

// Export individual agent projects for easy access
export const ecommerceTeam = masterAgentProjects[0];
export const mobileTeam = masterAgentProjects[1];
export const fullDevelopmentTeam = masterAgentProjects[2];

// Helper function to get agent project by project ID (for linking with regular projects)
export const getAgentProjectForProject = (
  projectId: string
): AgentProject | undefined => {
  switch (projectId) {
    case PROJECT_IDS.ECOMMERCE:
      return ecommerceTeam;
    case PROJECT_IDS.MOBILE:
      return mobileTeam;
    default:
      return fullDevelopmentTeam;
  }
};

// Export connection types for validation
export const CONNECTION_TYPES = {
  COMMUNICATION: "communication",
  HIERARCHY: "hierarchy",
  COLLABORATION: "collaboration",
} as const;
