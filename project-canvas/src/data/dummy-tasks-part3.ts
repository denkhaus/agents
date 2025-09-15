import { Task, TaskState } from '../types/task.types';
import { dummyProject } from './dummy-project';
import { generateUUID } from '../utils/uuid';

// Development Phase Task IDs
const TASK_FRONTEND_SETUP = generateUUID();
const TASK_BACKEND_SETUP = generateUUID();
const TASK_USER_AUTH = generateUUID();
const TASK_PRODUCT_CATALOG = generateUUID();
const TASK_SHOPPING_CART = generateUUID();
const TASK_PAYMENT_INTEGRATION = generateUUID();

// Agent UUIDs
const AGENT_FRONTEND_DEV = "frontend-dev-agent-id";
const AGENT_BACKEND_DEV = "backend-dev-agent-id";

export const dummyTasksPart3: Task[] = [
  // Phase 3: Development Setup
  {
    id: TASK_FRONTEND_SETUP,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Frontend Project Setup",
    description: "Initialize React project with TypeScript, configure build tools, linting, and testing framework.",
    state: TaskState.COMPLETED,
    complexity: 5,
    depth: 0,
    estimate: 480, // 8 hours
    assignedAgent: AGENT_FRONTEND_DEV,
    dependencies: ["wireframes-task-id", "design-system-task-id"],
    dependents: [TASK_USER_AUTH, TASK_PRODUCT_CATALOG],
    createdAt: new Date('2024-01-16T09:00:00Z'),
    updatedAt: new Date('2024-01-18T17:00:00Z'),
    completedAt: new Date('2024-01-18T17:00:00Z'),
    position: { x: 1000, y: 100 }
  },

  {
    id: TASK_BACKEND_SETUP,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Backend Infrastructure Setup",
    description: "Set up Node.js/Express server, configure database connections, middleware, and basic security.",
    state: TaskState.COMPLETED,
    complexity: 6,
    depth: 0,
    estimate: 720, // 12 hours
    assignedAgent: AGENT_BACKEND_DEV,
    dependencies: ["api-design-task-id", "database-schema-task-id"],
    dependents: [TASK_USER_AUTH, TASK_PRODUCT_CATALOG],
    createdAt: new Date('2024-01-16T09:30:00Z'),
    updatedAt: new Date('2024-01-19T16:30:00Z'),
    completedAt: new Date('2024-01-19T16:30:00Z'),
    position: { x: 1000, y: 300 }
  },

  // Phase 4: Core Features
  {
    id: TASK_USER_AUTH,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "User Authentication System",
    description: "Implement secure user registration, login, password reset, and session management.",
    state: TaskState.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: AGENT_BACKEND_DEV,
    dependencies: [TASK_FRONTEND_SETUP, TASK_BACKEND_SETUP],
    dependents: [TASK_PRODUCT_CATALOG, TASK_SHOPPING_CART],
    createdAt: new Date('2024-01-17T10:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1300, y: 200 }
  }
];