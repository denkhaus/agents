import { Project } from '../types/project.types';
import { generateUUID } from '../utils/uuid';

/**
 * Dummy Project Data
 * Comprehensive project with realistic structure for development and testing
 */
export const dummyProject: Project = {
  id: generateUUID(),
  title: "E-Commerce Platform Redesign",
  description: "Complete redesign and modernization of the existing e-commerce platform with improved UX, performance, and mobile responsiveness.",
  createdAt: new Date('2024-01-15T09:00:00Z'),
  updatedAt: new Date('2024-01-20T14:30:00Z'),
  totalTasks: 15,
  completedTasks: 4,
  progress: 26.7
};

export const dummyProjects: Project[] = [
  dummyProject,
  {
    id: generateUUID(),
    title: "Mobile App Development",
    description: "Native mobile application for iOS and Android platforms",
    createdAt: new Date('2024-01-10T10:00:00Z'),
    updatedAt: new Date('2024-01-18T16:45:00Z'),
    totalTasks: 8,
    completedTasks: 2,
    progress: 25.0
  },
  {
    id: generateUUID(),
    title: "API Modernization",
    description: "Migrate legacy REST API to GraphQL with improved performance",
    createdAt: new Date('2024-01-05T08:30:00Z'),
    updatedAt: new Date('2024-01-19T11:20:00Z'),
    totalTasks: 12,
    completedTasks: 8,
    progress: 66.7
  }
];