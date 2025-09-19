/**
 * Master Dummy Data - Index
 * Führt alle aufgeteilten Dummy-Daten zusammen
 */

// Import projects and agents from their respective files
import { PROJECT_IDS, masterProjects } from "./projects";
import { AGENT_IDS, masterAgents } from "./agents";
import { AGENT_PROJECT_IDS, masterAgentProjects } from "./agent-projects";

// Kombinierte Arrays für einfachen Zugriff
import { ecommerceTasks } from "./ecommerce-tasks";
import { ecommerceTasksPart2 } from "./ecommerce-tasks-part2";
import { ecommerceTasksPart3 } from "./ecommerce-tasks-part3";
import { mobileTasks } from "./mobile-tasks";
import { mobileTasksPart2 } from "./mobile-tasks-part2";

export const allEcommerceTasks = [
  ...ecommerceTasks,
  ...ecommerceTasksPart2,
  ...ecommerceTasksPart3,
];

export const allMobileTasks = [...mobileTasks, ...mobileTasksPart2];

export const allTasks = [...allEcommerceTasks, ...allMobileTasks];

// Re-export master projects, agents, and agent projects for consistency
export { masterProjects, PROJECT_IDS };
export { masterAgents, AGENT_IDS };
export { masterAgentProjects, AGENT_PROJECT_IDS };

// Task IDs for Dependencies (Go-Model: uuid.UUID)
// These should ideally be moved to their respective task files (ecommerce-tasks.ts, mobile-tasks.ts)
// if they are only used within those contexts. For now, keeping them here as they are referenced
// from multiple task data files.
export const ECOM_TASK_IDS = {
  // Root Phase Tasks
  PHASE_1_PLANNING: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
  PHASE_2_ARCHITECTURE: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
  PHASE_3_DEVELOPMENT: "c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f",
  PHASE_4_TESTING: "d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a",

  // Phase 1: Research & Planning (3 Tasks)
  RESEARCH: "e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b",
  WIREFRAMES: "f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c",
  DESIGN_SYSTEM: "a7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d",

  // Phase 2: Architecture (4 Tasks)
  API_DESIGN: "b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e",
  DATABASE_SCHEMA: "c9d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f",
  FRONTEND_SETUP: "d0e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a",
  BACKEND_SETUP: "e1f2a3b4-c5d6-4e7f-8a9b-0c1d2e3f4a5b",

  // Phase 3: Core Development (5 Tasks)
  USER_AUTH: "f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c",
  PRODUCT_CATALOG: "a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d",
  SHOPPING_CART: "b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e",
  PAYMENT_INTEGRATION: "c5d6e7f8-a9b0-4c1d-2e3f-4a5b6c7d8e9f",
  ADMIN_PANEL: "d6e7f8a9-b0c1-4d2e-3f4a-5b6c7d8e9f0a",

  // Phase 4: Testing & Deployment (3 Tasks)
  TESTING_SETUP: "e7f8a9b0-c1d2-4e3f-4a5b-6c7d8e9f0a1b",
  E2E_TESTS: "f8a9b0c1-d2e3-4f4a-5b6c-7d8e9f0a1b2c",
  DEPLOYMENT: "a9b0c1d2-e3f4-4a5b-6c7d-8e9f0a1b2c3d",
};

export const MOBILE_TASK_IDS = {
  // Phase 1: Research & Planning (3 Tasks)
  RESEARCH: "1a2b3c4d-5e6f-4789-abcd-ef1234567890",
  WIREFRAMES: "2b3c4d5e-6f7a-4890-bcde-f12345678901",
  DESIGN_SYSTEM: "3c4d5e6f-7a8b-4901-cdef-123456789012",

  // Phase 2: Architecture (4 Tasks)
  TECH_STACK: "4d5e6f7a-8b9c-4012-def1-234567890123",
  API_INTEGRATION: "5e6f7a8b-9c0d-4123-ef12-345678901234",
  STATE_MANAGEMENT: "6f7a8b9c-0d1e-4234-f123-456789012345",
  NAVIGATION: "7a8b9c0d-1e2f-4345-1234-567890123456",

  // Phase 3: Core Development (5 Tasks)
  AUTH_MODULE: "8b9c0d1e-2f3a-4456-2345-678901234567",
  HOME_SCREEN: "9c0d1e2f-3a4b-4567-3456-789012345678",
  PROFILE_SCREEN: "0d1e2f3a-4b5c-4678-4567-890123456789",
  NOTIFICATIONS: "1e2f3a4b-5c6d-4789-5678-901234567890",
  SETTINGS: "2f3a4b5c-6d7e-4890-6789-012345678901",

  // Phase 4: Testing & Deployment (3 Tasks)
  TESTING: "3a4b5c6d-7e8f-4901-7890-123456789012",
  APP_STORE: "4b5c6d7e-8f9a-4012-8901-234567890123",
  LAUNCH: "5c6d7e8f-9a0b-4123-9012-345678901234",
};
