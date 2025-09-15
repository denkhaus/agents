import { Task, TaskState } from '../types/task.types';
import { dummyProject } from './dummy-project';
import { generateUUID } from '../utils/uuid';

/**
 * Dummy Task Data with Complex Dependencies
 * Realistic task hierarchy with various states and dependencies
 */

// Agent UUIDs for assignment
const AGENT_DESIGNER = generateUUID();
const AGENT_FRONTEND_DEV = generateUUID();
const AGENT_BACKEND_DEV = generateUUID();
const AGENT_QA_ENGINEER = generateUUID();
const AGENT_DEVOPS = generateUUID();

// Task IDs for dependency management
const TASK_RESEARCH = generateUUID();
const TASK_WIREFRAMES = generateUUID();
const TASK_DESIGN_SYSTEM = generateUUID();
const TASK_API_DESIGN = generateUUID();
const TASK_DATABASE_SCHEMA = generateUUID();
const TASK_FRONTEND_SETUP = generateUUID();
const TASK_BACKEND_SETUP = generateUUID();
const TASK_USER_AUTH = generateUUID();
const TASK_PRODUCT_CATALOG = generateUUID();
const TASK_SHOPPING_CART = generateUUID();
const TASK_PAYMENT_INTEGRATION = generateUUID();
const TASK_TESTING_SETUP = generateUUID();
const TASK_UNIT_TESTS = generateUUID();
const TASK_INTEGRATION_TESTS = generateUUID();
const TASK_DEPLOYMENT = generateUUID();

export const dummyTasks: Task[] = [
  // Phase 1: Research & Planning
  {
    id: TASK_RESEARCH,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Market Research & User Analysis",
    description: "Conduct comprehensive market research and analyze user behavior patterns to inform design decisions.",
    state: TaskState.COMPLETED,
    complexity: 6,
    depth: 0,
    estimate: 2400, // 40 hours
    assignedAgent: AGENT_DESIGNER,
    dependencies: [],
    dependents: [TASK_WIREFRAMES, TASK_DESIGN_SYSTEM],
    createdAt: new Date('2024-01-15T09:00:00Z'),
    updatedAt: new Date('2024-01-17T16:30:00Z'),
    completedAt: new Date('2024-01-17T16:30:00Z'),
    position: { x: 100, y: 100 }
  },

  {
    id: TASK_WIREFRAMES,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Create Wireframes & Prototypes",
    description: "Design low-fidelity wireframes and interactive prototypes for all major user flows.",
    state: TaskState.COMPLETED,
    complexity: 7,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: AGENT_DESIGNER,
    dependencies: [TASK_RESEARCH],
    dependents: [TASK_DESIGN_SYSTEM, TASK_FRONTEND_SETUP],
    createdAt: new Date('2024-01-15T09:15:00Z'),
    updatedAt: new Date('2024-01-19T14:20:00Z'),
    completedAt: new Date('2024-01-19T14:20:00Z'),
    position: { x: 400, y: 100 }
  }
];