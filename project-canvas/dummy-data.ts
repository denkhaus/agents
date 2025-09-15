// dummy-data.ts
// Comprehensive dummy data for testing the Project Canvas application

import { v4 as uuidv4 } from 'uuid';

export type UUID = string;

export enum TaskState {
  PENDING = "pending",
  IN_PROGRESS = "in-progress", 
  COMPLETED = "completed",
  BLOCKED = "blocked",
  CANCELLED = "cancelled"
}

export interface Task {
  id: UUID;
  projectId: UUID;
  parentId?: UUID;
  title: string;
  description: string;
  state: TaskState;
  complexity: number; // 1-10
  depth: number;
  estimate?: number; // minutes
  assignedAgent?: UUID;
  dependencies: UUID[];
  dependents: UUID[];
  createdAt: Date;
  updatedAt: Date;
  completedAt?: Date;
  // UI-specific
  position?: { x: number; y: number };
}

export interface Project {
  id: UUID;
  title: string;
  description: string;
  createdAt: Date;
  updatedAt: Date;
  totalTasks: number;
  completedTasks: number;
  progress: number; // 0-100
}

export interface Agent {
  id: UUID;
  name: string;
  role: string;
  description: string;
  isActive: boolean;
  createdAt: Date;
}

// Generate UUIDs for consistency
const PROJECT_ID = uuidv4();
const AGENT_IDS = {
  supervisor: uuidv4(),
  projectManager: uuidv4(),
  coder: uuidv4(),
  researcher: uuidv4(),
};

// Task IDs for dependency management
const TASK_IDS = {
  // Root tasks
  planning: uuidv4(),
  development: uuidv4(),
  testing: uuidv4(),
  deployment: uuidv4(),
  
  // Planning subtasks
  requirements: uuidv4(),
  architecture: uuidv4(),
  wireframes: uuidv4(),
  
  // Development subtasks
  backend: uuidv4(),
  frontend: uuidv4(),
  database: uuidv4(),
  api: uuidv4(),
  
  // Backend subtasks
  userAuth: uuidv4(),
  dataModels: uuidv4(),
  businessLogic: uuidv4(),
  
  // Frontend subtasks
  components: uuidv4(),
  routing: uuidv4(),
  stateManagement: uuidv4(),
  
  // Testing subtasks
  unitTests: uuidv4(),
  integrationTests: uuidv4(),
  e2eTests: uuidv4(),
  
  // Deployment subtasks
  cicd: uuidv4(),
  monitoring: uuidv4(),
  documentation: uuidv4(),
};

export const dummyAgents: Agent[] = [
  {
    id: AGENT_IDS.supervisor,
    name: "Supervisor Agent",
    role: "supervisor",
    description: "Team lead with authority over all specialized agents",
    isActive: true,
    createdAt: new Date('2024-01-01T09:00:00Z'),
  },
  {
    id: AGENT_IDS.projectManager,
    name: "Project Manager Agent",
    role: "project-manager",
    description: "Operational right hand of supervisor for project management",
    isActive: true,
    createdAt: new Date('2024-01-01T09:15:00Z'),
  },
  {
    id: AGENT_IDS.coder,
    name: "Coder Agent",
    role: "coder",
    description: "Skilled software engineer specialized in Golang development",
    isActive: true,
    createdAt: new Date('2024-01-01T09:30:00Z'),
  },
  {
    id: AGENT_IDS.researcher,
    name: "Researcher Agent",
    role: "researcher",
    description: "Research assistant specialized in gathering and analyzing information",
    isActive: true,
    createdAt: new Date('2024-01-01T09:45:00Z'),
  },
];

export const dummyProjects: Project[] = [
  {
    id: PROJECT_ID,
    title: "E-Commerce Platform Development",
    description: "Complete development of a modern e-commerce platform with real-time features, user management, and payment processing.",
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-15T14:30:00Z'),
    totalTasks: 19,
    completedTasks: 8,
    progress: 42.1,
  },
  {
    id: uuidv4(),
    title: "Mobile App Redesign",
    description: "Complete UI/UX redesign of the mobile application with improved user experience and performance optimizations.",
    createdAt: new Date('2024-01-10T08:00:00Z'),
    updatedAt: new Date('2024-01-20T16:45:00Z'),
    totalTasks: 12,
    completedTasks: 3,
    progress: 25.0,
  },
  {
    id: uuidv4(),
    title: "Data Analytics Dashboard",
    description: "Development of a comprehensive analytics dashboard for business intelligence and reporting.",
    createdAt: new Date('2024-01-05T11:30:00Z'),
    updatedAt: new Date('2024-01-18T13:15:00Z'),
    totalTasks: 15,
    completedTasks: 12,
    progress: 80.0,
  },
];

export const dummyTasks: Task[] = [
  // ROOT LEVEL TASKS (Depth 0)
  {
    id: TASK_IDS.planning,
    projectId: PROJECT_ID,
    title: "Project Planning & Design",
    description: "Complete project planning phase including requirements gathering, architecture design, and wireframe creation.",
    state: TaskState.COMPLETED,
    complexity: 8,
    depth: 0,
    estimate: 2400, // 40 hours
    assignedAgent: AGENT_IDS.projectManager,
    dependencies: [],
    dependents: [TASK_IDS.development],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-05T17:00:00Z'),
    completedAt: new Date('2024-01-05T17:00:00Z'),
    position: { x: 100, y: 100 },
  },
  {
    id: TASK_IDS.development,
    projectId: PROJECT_ID,
    title: "Core Development",
    description: "Implementation of all core features including backend services, frontend components, and database integration.",
    state: TaskState.IN_PROGRESS,
    complexity: 10,
    depth: 0,
    estimate: 7200, // 120 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.planning],
    dependents: [TASK_IDS.testing],
    createdAt: new Date('2024-01-05T18:00:00Z'),
    updatedAt: new Date('2024-01-15T14:30:00Z'),
    position: { x: 400, y: 100 },
  },
  {
    id: TASK_IDS.testing,
    projectId: PROJECT_ID,
    title: "Quality Assurance & Testing",
    description: "Comprehensive testing including unit tests, integration tests, and end-to-end testing scenarios.",
    state: TaskState.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 2880, // 48 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.development],
    dependents: [TASK_IDS.deployment],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-01T10:00:00Z'),
    position: { x: 700, y: 100 },
  },
  {
    id: TASK_IDS.deployment,
    projectId: PROJECT_ID,
    title: "Deployment & Launch",
    description: "Production deployment, monitoring setup, and documentation finalization for project launch.",
    state: TaskState.BLOCKED,
    complexity: 6,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: AGENT_IDS.projectManager,
    dependencies: [TASK_IDS.testing],
    dependents: [],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-01T10:00:00Z'),
    position: { x: 1000, y: 100 },
  },

  // PLANNING SUBTASKS (Depth 1)
  {
    id: TASK_IDS.requirements,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.planning,
    title: "Requirements Gathering",
    description: "Detailed analysis of business requirements, user stories, and technical specifications.",
    state: TaskState.COMPLETED,
    complexity: 5,
    depth: 1,
    estimate: 960, // 16 hours
    assignedAgent: AGENT_IDS.researcher,
    dependencies: [],
    dependents: [TASK_IDS.architecture, TASK_IDS.wireframes],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-02T16:00:00Z'),
    completedAt: new Date('2024-01-02T16:00:00Z'),
    position: { x: 50, y: 250 },
  },
  {
    id: TASK_IDS.architecture,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.planning,
    title: "System Architecture Design",
    description: "Design of system architecture, database schema, and API specifications.",
    state: TaskState.COMPLETED,
    complexity: 7,
    depth: 1,
    estimate: 720, // 12 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.requirements],
    dependents: [TASK_IDS.backend, TASK_IDS.database],
    createdAt: new Date('2024-01-02T17:00:00Z'),
    updatedAt: new Date('2024-01-04T15:00:00Z'),
    completedAt: new Date('2024-01-04T15:00:00Z'),
    position: { x: 100, y: 250 },
  },
  {
    id: TASK_IDS.wireframes,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.planning,
    title: "UI/UX Wireframes",
    description: "Creation of detailed wireframes and user interface mockups for all application screens.",
    state: TaskState.COMPLETED,
    complexity: 4,
    depth: 1,
    estimate: 720, // 12 hours
    assignedAgent: AGENT_IDS.researcher,
    dependencies: [TASK_IDS.requirements],
    dependents: [TASK_IDS.frontend, TASK_IDS.components],
    createdAt: new Date('2024-01-02T17:00:00Z'),
    updatedAt: new Date('2024-01-05T12:00:00Z'),
    completedAt: new Date('2024-01-05T12:00:00Z'),
    position: { x: 150, y: 250 },
  },

  // DEVELOPMENT SUBTASKS (Depth 1)
  {
    id: TASK_IDS.backend,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.development,
    title: "Backend Development",
    description: "Implementation of server-side logic, APIs, and business rules.",
    state: TaskState.IN_PROGRESS,
    complexity: 9,
    depth: 1,
    estimate: 2880, // 48 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.architecture],
    dependents: [TASK_IDS.api],
    createdAt: new Date('2024-01-05T18:00:00Z'),
    updatedAt: new Date('2024-01-15T14:30:00Z'),
    position: { x: 350, y: 250 },
  },
  {
    id: TASK_IDS.frontend,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.development,
    title: "Frontend Development",
    description: "Implementation of user interface components and client-side functionality.",
    state: TaskState.PENDING,
    complexity: 8,
    depth: 1,
    estimate: 2400, // 40 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.wireframes],
    dependents: [],
    createdAt: new Date('2024-01-05T18:00:00Z'),
    updatedAt: new Date('2024-01-05T18:00:00Z'),
    position: { x: 400, y: 250 },
  },
  {
    id: TASK_IDS.database,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.development,
    title: "Database Implementation",
    description: "Database setup, schema implementation, and data migration scripts.",
    state: TaskState.COMPLETED,
    complexity: 6,
    depth: 1,
    estimate: 1440, // 24 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.architecture],
    dependents: [TASK_IDS.dataModels],
    createdAt: new Date('2024-01-05T18:00:00Z'),
    updatedAt: new Date('2024-01-10T16:00:00Z'),
    completedAt: new Date('2024-01-10T16:00:00Z'),
    position: { x: 450, y: 250 },
  },
  {
    id: TASK_IDS.api,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.development,
    title: "API Development",
    description: "RESTful API endpoints implementation and documentation.",
    state: TaskState.PENDING,
    complexity: 7,
    depth: 1,
    estimate: 1920, // 32 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.backend],
    dependents: [TASK_IDS.integrationTests],
    createdAt: new Date('2024-01-05T18:00:00Z'),
    updatedAt: new Date('2024-01-05T18:00:00Z'),
    position: { x: 500, y: 250 },
  },

  // BACKEND SUBTASKS (Depth 2)
  {
    id: TASK_IDS.userAuth,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.backend,
    title: "User Authentication System",
    description: "Implementation of user registration, login, JWT tokens, and session management.",
    state: TaskState.COMPLETED,
    complexity: 6,
    depth: 2,
    estimate: 960, // 16 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [],
    dependents: [],
    createdAt: new Date('2024-01-06T09:00:00Z'),
    updatedAt: new Date('2024-01-09T17:00:00Z'),
    completedAt: new Date('2024-01-09T17:00:00Z'),
    position: { x: 300, y: 400 },
  },
  {
    id: TASK_IDS.dataModels,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.backend,
    title: "Data Models & ORM",
    description: "Implementation of data models, ORM configuration, and database relationships.",
    state: TaskState.COMPLETED,
    complexity: 5,
    depth: 2,
    estimate: 720, // 12 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.database],
    dependents: [],
    createdAt: new Date('2024-01-10T17:00:00Z'),
    updatedAt: new Date('2024-01-12T15:00:00Z'),
    completedAt: new Date('2024-01-12T15:00:00Z'),
    position: { x: 350, y: 400 },
  },
  {
    id: TASK_IDS.businessLogic,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.backend,
    title: "Business Logic Implementation",
    description: "Core business rules, validation logic, and service layer implementation.",
    state: TaskState.IN_PROGRESS,
    complexity: 8,
    depth: 2,
    estimate: 1200, // 20 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.dataModels],
    dependents: [],
    createdAt: new Date('2024-01-12T16:00:00Z'),
    updatedAt: new Date('2024-01-15T14:30:00Z'),
    position: { x: 400, y: 400 },
  },

  // FRONTEND SUBTASKS (Depth 2)
  {
    id: TASK_IDS.components,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.frontend,
    title: "React Components",
    description: "Development of reusable React components and UI elements.",
    state: TaskState.PENDING,
    complexity: 6,
    depth: 2,
    estimate: 960, // 16 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.wireframes],
    dependents: [],
    createdAt: new Date('2024-01-05T18:00:00Z'),
    updatedAt: new Date('2024-01-05T18:00:00Z'),
    position: { x: 450, y: 400 },
  },
  {
    id: TASK_IDS.routing,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.frontend,
    title: "Application Routing",
    description: "Setup of React Router and navigation between different application screens.",
    state: TaskState.PENDING,
    complexity: 3,
    depth: 2,
    estimate: 480, // 8 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [],
    dependents: [],
    createdAt: new Date('2024-01-05T18:00:00Z'),
    updatedAt: new Date('2024-01-05T18:00:00Z'),
    position: { x: 500, y: 400 },
  },
  {
    id: TASK_IDS.stateManagement,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.frontend,
    title: "State Management",
    description: "Implementation of global state management using Zustand and context providers.",
    state: TaskState.PENDING,
    complexity: 5,
    depth: 2,
    estimate: 720, // 12 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [],
    dependents: [],
    createdAt: new Date('2024-01-05T18:00:00Z'),
    updatedAt: new Date('2024-01-05T18:00:00Z'),
    position: { x: 550, y: 400 },
  },

  // TESTING SUBTASKS (Depth 1)
  {
    id: TASK_IDS.unitTests,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.testing,
    title: "Unit Testing",
    description: "Comprehensive unit tests for all components and business logic functions.",
    state: TaskState.PENDING,
    complexity: 6,
    depth: 1,
    estimate: 960, // 16 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [],
    dependents: [],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-01T10:00:00Z'),
    position: { x: 650, y: 250 },
  },
  {
    id: TASK_IDS.integrationTests,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.testing,
    title: "Integration Testing",
    description: "Integration tests for API endpoints and database interactions.",
    state: TaskState.PENDING,
    complexity: 7,
    depth: 1,
    estimate: 1200, // 20 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [TASK_IDS.api],
    dependents: [],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-01T10:00:00Z'),
    position: { x: 700, y: 250 },
  },
  {
    id: TASK_IDS.e2eTests,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.testing,
    title: "End-to-End Testing",
    description: "Complete user journey testing using automated testing frameworks.",
    state: TaskState.PENDING,
    complexity: 5,
    depth: 1,
    estimate: 720, // 12 hours
    assignedAgent: AGENT_IDS.coder,
    dependencies: [],
    dependents: [],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-01T10:00:00Z'),
    position: { x: 750, y: 250 },
  },

  // DEPLOYMENT SUBTASKS (Depth 1)
  {
    id: TASK_IDS.cicd,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.deployment,
    title: "CI/CD Pipeline",
    description: "Setup of continuous integration and deployment pipeline with automated testing.",
    state: TaskState.PENDING,
    complexity: 6,
    depth: 1,
    estimate: 720, // 12 hours
    assignedAgent: AGENT_IDS.projectManager,
    dependencies: [],
    dependents: [],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-01T10:00:00Z'),
    position: { x: 950, y: 250 },
  },
  {
    id: TASK_IDS.monitoring,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.deployment,
    title: "Monitoring & Logging",
    description: "Implementation of application monitoring, logging, and alerting systems.",
    state: TaskState.PENDING,
    complexity: 4,
    depth: 1,
    estimate: 480, // 8 hours
    assignedAgent: AGENT_IDS.projectManager,
    dependencies: [],
    dependents: [],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-01T10:00:00Z'),
    position: { x: 1000, y: 250 },
  },
  {
    id: TASK_IDS.documentation,
    projectId: PROJECT_ID,
    parentId: TASK_IDS.deployment,
    title: "Documentation & Handover",
    description: "Complete project documentation, API docs, and deployment guides.",
    state: TaskState.PENDING,
    complexity: 3,
    depth: 1,
    estimate: 240, // 4 hours
    assignedAgent: AGENT_IDS.researcher,
    dependencies: [],
    dependents: [],
    createdAt: new Date('2024-01-01T10:00:00Z'),
    updatedAt: new Date('2024-01-01T10:00:00Z'),
    position: { x: 1050, y: 250 },
  },
];

// Helper functions for data manipulation
export const getTasksByProject = (projectId: UUID): Task[] => {
  return dummyTasks.filter(task => task.projectId === projectId);
};

export const getTasksByDepth = (projectId: UUID, depth: number): Task[] => {
  return dummyTasks.filter(task => task.projectId === projectId && task.depth === depth);
};

export const getRootTasks = (projectId: UUID): Task[] => {
  return getTasksByDepth(projectId, 0);
};

export const getTaskDependencies = (taskId: UUID): Task[] => {
  const task = dummyTasks.find(t => t.id === taskId);
  if (!task) return [];
  
  return dummyTasks.filter(t => task.dependencies.includes(t.id));
};

export const getTaskDependents = (taskId: UUID): Task[] => {
  const task = dummyTasks.find(t => t.id === taskId);
  if (!task) return [];
  
  return dummyTasks.filter(t => task.dependents.includes(t.id));
};

export const getTasksByAgent = (agentId: UUID): Task[] => {
  return dummyTasks.filter(task => task.assignedAgent === agentId);
};

export const getProjectProgress = (projectId: UUID) => {
  const tasks = getTasksByProject(projectId);
  const totalTasks = tasks.length;
  const completedTasks = tasks.filter(t => t.state === TaskState.COMPLETED).length;
  const inProgressTasks = tasks.filter(t => t.state === TaskState.IN_PROGRESS).length;
  const pendingTasks = tasks.filter(t => t.state === TaskState.PENDING).length;
  const blockedTasks = tasks.filter(t => t.state === TaskState.BLOCKED).length;
  const cancelledTasks = tasks.filter(t => t.state === TaskState.CANCELLED).length;
  
  const tasksByDepth = tasks.reduce((acc, task) => {
    acc[task.depth] = (acc[task.depth] || 0) + 1;
    return acc;
  }, {} as Record<number, number>);
  
  return {
    projectId,
    totalTasks,
    completedTasks,
    inProgressTasks,
    pendingTasks,
    blockedTasks,
    cancelledTasks,
    overallProgress: totalTasks > 0 ? (completedTasks / totalTasks) * 100 : 0,
    tasksByDepth,
  };
};

// Export all data for easy import
export const dummyData = {
  projects: dummyProjects,
  tasks: dummyTasks,
  agents: dummyAgents,
  mainProject: dummyProjects[0],
  mainProjectTasks: getTasksByProject(PROJECT_ID),
  mainProjectProgress: getProjectProgress(PROJECT_ID),
};