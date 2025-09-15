import { Task, TaskState } from '../types/task.types';
import { dummyProject } from './dummy-project';
import { generateUUID } from '../utils/uuid';

// Testing & Deployment Task IDs
const TASK_TESTING_SETUP = generateUUID();
const TASK_UNIT_TESTS = generateUUID();
const TASK_INTEGRATION_TESTS = generateUUID();
const TASK_DEPLOYMENT = generateUUID();

// Agent UUIDs
const AGENT_QA_ENGINEER = "qa-engineer-agent-id";
const AGENT_DEVOPS = "devops-agent-id";

export const dummyTasksPart5: Task[] = [
  // Phase 5: Testing
  {
    id: TASK_TESTING_SETUP,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Testing Framework Setup",
    description: "Configure Jest, Cypress, and testing utilities for comprehensive test coverage.",
    state: TaskState.PENDING,
    complexity: 4,
    depth: 0,
    estimate: 360, // 6 hours
    assignedAgent: AGENT_QA_ENGINEER,
    dependencies: ["frontend-setup-task-id"],
    dependents: [TASK_UNIT_TESTS, TASK_INTEGRATION_TESTS],
    createdAt: new Date('2024-01-18T10:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1300, y: 500 }
  },

  {
    id: TASK_UNIT_TESTS,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Unit Tests Implementation",
    description: "Write comprehensive unit tests for all components, utilities, and business logic.",
    state: TaskState.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: AGENT_QA_ENGINEER,
    dependencies: [TASK_TESTING_SETUP],
    dependents: [TASK_DEPLOYMENT],
    createdAt: new Date('2024-01-18T11:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1600, y: 500 }
  },

  {
    id: TASK_INTEGRATION_TESTS,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Integration & E2E Tests",
    description: "Create end-to-end tests for critical user flows and API integration tests.",
    state: TaskState.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: AGENT_QA_ENGINEER,
    dependencies: [TASK_TESTING_SETUP, "payment-integration-task-id"],
    dependents: [TASK_DEPLOYMENT],
    createdAt: new Date('2024-01-18T12:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1900, y: 500 }
  },

  // Phase 6: Deployment
  {
    id: TASK_DEPLOYMENT,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Production Deployment & CI/CD",
    description: "Set up production environment, configure CI/CD pipeline, monitoring, and deployment automation.",
    state: TaskState.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: AGENT_DEVOPS,
    dependencies: [TASK_UNIT_TESTS, TASK_INTEGRATION_TESTS],
    dependents: [],
    createdAt: new Date('2024-01-19T09:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 2200, y: 600 }
  }
];