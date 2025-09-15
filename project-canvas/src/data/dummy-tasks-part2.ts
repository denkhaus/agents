import { Task, TaskState } from '../types/task.types';
import { dummyProject } from './dummy-project';
import { generateUUID } from '../utils/uuid';

// Import task IDs from part 1
const TASK_RESEARCH = "research-task-id"; // These would be imported from part1
const TASK_WIREFRAMES = "wireframes-task-id";

// Continue with more task IDs
const TASK_DESIGN_SYSTEM = generateUUID();
const TASK_API_DESIGN = generateUUID();
const TASK_DATABASE_SCHEMA = generateUUID();
const TASK_FRONTEND_SETUP = generateUUID();
const TASK_BACKEND_SETUP = generateUUID();

// Agent UUIDs
const AGENT_DESIGNER = "designer-agent-id";
const AGENT_FRONTEND_DEV = "frontend-dev-agent-id";
const AGENT_BACKEND_DEV = "backend-dev-agent-id";

export const dummyTasksPart2: Task[] = [
  // Phase 2: Design System & Architecture
  {
    id: TASK_DESIGN_SYSTEM,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Design System & Component Library",
    description: "Create comprehensive design system with reusable components, color schemes, typography, and interaction patterns.",
    state: TaskState.IN_PROGRESS,
    complexity: 8,
    depth: 0,
    estimate: 2160, // 36 hours
    assignedAgent: AGENT_DESIGNER,
    dependencies: [TASK_RESEARCH, TASK_WIREFRAMES],
    dependents: [TASK_FRONTEND_SETUP],
    createdAt: new Date('2024-01-15T09:30:00Z'),
    updatedAt: new Date('2024-01-20T10:15:00Z'),
    position: { x: 700, y: 100 }
  },

  {
    id: TASK_API_DESIGN,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "API Architecture & Documentation",
    description: "Design RESTful API endpoints, define data models, and create comprehensive API documentation.",
    state: TaskState.IN_PROGRESS,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: AGENT_BACKEND_DEV,
    dependencies: [TASK_RESEARCH],
    dependents: [TASK_DATABASE_SCHEMA, TASK_BACKEND_SETUP],
    createdAt: new Date('2024-01-15T10:00:00Z'),
    updatedAt: new Date('2024-01-20T11:30:00Z'),
    position: { x: 400, y: 300 }
  },

  {
    id: TASK_DATABASE_SCHEMA,
    projectId: dummyProject.id,
    parentId: undefined,
    title: "Database Schema Design",
    description: "Design optimized database schema with proper indexing, relationships, and migration scripts.",
    state: TaskState.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: AGENT_BACKEND_DEV,
    dependencies: [TASK_API_DESIGN],
    dependents: [TASK_BACKEND_SETUP],
    createdAt: new Date('2024-01-15T10:15:00Z'),
    updatedAt: new Date('2024-01-20T11:30:00Z'),
    position: { x: 700, y: 300 }
  }
];