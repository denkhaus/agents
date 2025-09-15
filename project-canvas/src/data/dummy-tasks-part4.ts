import { Task, TaskState } from '../types/task.types';
import { dummyProject } from './dummy-project';
import { generateUUID } from '../utils/uuid';

// Feature Development Task IDs
const TASK_PRODUCT_CATALOG = generateUUID();
const TASK_SHOPPING_CART = generateUUID();
const TASK_PAYMENT_INTEGRATION = generateUUID();
const TASK_TESTING_SETUP = generateUUID();
const TASK_UNIT_TESTS = generateUUID();
const TASK_INTEGRATION_TESTS = generateUUID();
const TASK_DEPLOYMENT = generateUUID();

// Agent UUIDs
const AGENT_FRONTEND_DEV = "frontend-dev-agent-id";
const AGENT_BACKEND_DEV = "backend-dev-agent-id";
const AGENT_QA_ENGINEER = "qa-engineer-agent-id";
const AGENT_DEVOPS = "devops-agent-id";

export const dummyTasksPart4: Task[] = [
  {
    id: TASK_PRODUCT_CATALOG,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Product Catalog & Search",
    description: "Build product listing, filtering, search functionality, and detailed product pages.",
    state: TaskState.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: AGENT_FRONTEND_DEV,
    dependencies: ["user-auth-task-id", "frontend-setup-task-id"],
    dependents: [TASK_SHOPPING_CART],
    createdAt: new Date('2024-01-17T11:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1600, y: 100 }
  },

  {
    id: TASK_SHOPPING_CART,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Shopping Cart & Checkout",
    description: "Implement shopping cart functionality, checkout process, and order management.",
    state: TaskState.PENDING,
    complexity: 9,
    depth: 0,
    estimate: 2160, // 36 hours
    assignedAgent: AGENT_FRONTEND_DEV,
    dependencies: ["user-auth-task-id", TASK_PRODUCT_CATALOG],
    dependents: [TASK_PAYMENT_INTEGRATION],
    createdAt: new Date('2024-01-17T12:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1900, y: 200 }
  },

  {
    id: TASK_PAYMENT_INTEGRATION,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Payment Gateway Integration",
    description: "Integrate secure payment processing with Stripe/PayPal, handle transactions and refunds.",
    state: TaskState.BLOCKED,
    complexity: 8,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: AGENT_BACKEND_DEV,
    dependencies: [TASK_SHOPPING_CART],
    dependents: [TASK_INTEGRATION_TESTS],
    createdAt: new Date('2024-01-18T09:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 2200, y: 300 }
  }
];