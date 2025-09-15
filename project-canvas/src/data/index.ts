/**
 * Dummy Data Index
 * Combines all dummy data parts into complete datasets
 */

import { dummyProject, dummyProjects } from './dummy-project';
import { dummyTasks } from './dummy-tasks';
import { dummyTasksPart2 } from './dummy-tasks-part2';
import { dummyTasksPart3 } from './dummy-tasks-part3';
import { dummyTasksPart4 } from './dummy-tasks-part4';
import { dummyTasksPart5 } from './dummy-tasks-part5';
import { Task } from '../types/task.types';
import { Project } from '../types/project.types';
import { Agent, AgentRole, AgentStatus } from '../types/agent.types';
import { generateUUID } from '../utils/uuid';

// Combine all task parts
export const allDummyTasks: Task[] = [
  ...dummyTasks,
  ...dummyTasksPart2,
  ...dummyTasksPart3,
  ...dummyTasksPart4,
  ...dummyTasksPart5
];

// Export projects
export const allDummyProjects: Project[] = dummyProjects;
export { dummyProject as mainDummyProject };

// Create dummy agents
export const dummyAgents: Agent[] = [
  {
    id: generateUUID(),
    name: "Design Lead",
    role: AgentRole.DESIGNER,
    description: "Senior UX/UI designer specializing in e-commerce platforms",
    status: AgentStatus.ONLINE,
    isStreaming: false,
    capabilities: ["wireframing", "prototyping", "design-systems", "user-research"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T09:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T14:30:00Z')
  },
  {
    id: generateUUID(),
    name: "Frontend Developer",
    role: AgentRole.CODER,
    description: "React/TypeScript specialist with expertise in modern frontend development",
    status: AgentStatus.BUSY,
    isStreaming: true,
    capabilities: ["react", "typescript", "tailwind", "testing", "performance-optimization"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T09:15:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T14:25:00Z')
  },
  {
    id: generateUUID(),
    name: "Backend Developer",
    role: AgentRole.CODER,
    description: "Node.js/Express expert with database and API design experience",
    status: AgentStatus.ONLINE,
    isStreaming: false,
    capabilities: ["nodejs", "express", "postgresql", "api-design", "security"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T09:30:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T14:20:00Z')
  },
  {
    id: generateUUID(),
    name: "QA Engineer",
    role: AgentRole.QA_ENGINEER,
    description: "Quality assurance specialist with automated testing expertise",
    status: AgentStatus.IDLE,
    isStreaming: false,
    capabilities: ["jest", "cypress", "test-automation", "performance-testing", "accessibility"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T10:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T13:45:00Z')
  },
  {
    id: generateUUID(),
    name: "DevOps Engineer",
    role: AgentRole.DEVOPS,
    description: "Infrastructure and deployment automation specialist",
    status: AgentStatus.ONLINE,
    isStreaming: false,
    capabilities: ["docker", "kubernetes", "ci-cd", "monitoring", "cloud-infrastructure"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T10:15:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T14:10:00Z')
  }
];

// Helper function to get tasks by project
export const getTasksByProject = (projectId: string): Task[] => {
  return allDummyTasks.filter(task => task.projectId === projectId);
};

// Helper function to get root tasks (no parent)
export const getRootTasks = (projectId: string): Task[] => {
  return allDummyTasks.filter(task => 
    task.projectId === projectId && !task.parentId
  );
};

// Helper function to get task dependencies
export const getTaskDependencies = (taskId: string): Task[] => {
  const task = allDummyTasks.find(t => t.id === taskId);
  if (!task) return [];
  
  return allDummyTasks.filter(t => task.dependencies.includes(t.id));
};

// Helper function to get dependent tasks
export const getDependentTasks = (taskId: string): Task[] => {
  return allDummyTasks.filter(task => task.dependencies.includes(taskId));
};