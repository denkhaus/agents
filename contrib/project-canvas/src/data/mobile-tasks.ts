/**
 * Mobile App Project Tasks (15 Tasks)
 * Go-Model konform: pkg/tools/project/shared/shared.go
 * Max 500 Zeilen pro Datei!
 */

import type { Task } from '../types';
import { generateUUID } from '../utils/uuid';
import { mobileProject } from './projects';
import { designerAgent, frontendAgent } from './agents';

// Import the actual values for runtime usage
import { TaskState as TaskStateValue } from '../types/task.types';

// Task IDs für Dependencies
export const MOBILE_TASK_IDS = {
  // Phase 1: Research & Planning (3 Tasks)
  RESEARCH: generateUUID(),
  WIREFRAMES: generateUUID(),
  DESIGN_SYSTEM: generateUUID(),
  
  // Phase 2: Architecture (4 Tasks)
  TECH_STACK: generateUUID(),
  API_INTEGRATION: generateUUID(),
  STATE_MANAGEMENT: generateUUID(),
  NAVIGATION: generateUUID(),
  
  // Phase 3: Core Development (5 Tasks)
  AUTH_MODULE: generateUUID(),
  HOME_SCREEN: generateUUID(),
  PROFILE_SCREEN: generateUUID(),
  NOTIFICATIONS: generateUUID(),
  SETTINGS: generateUUID(),
  
  // Phase 4: Testing & Deployment (3 Tasks)
  TESTING: generateUUID(),
  APP_STORE: generateUUID(),
  LAUNCH: generateUUID()
};

export const mobileTasks: Task[] = [
  // Phase 1: Research & Planning
  {
    id: MOBILE_TASK_IDS.RESEARCH,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Market Research & User Analysis",
    description: `## Mobile App Market Research

**Objective**: Understand mobile app market trends and user behavior patterns.

### Activities
- **Competitor Analysis**: Study 10+ mobile apps in similar domains
- **User Surveys**: Conduct 20+ user interviews with target demographics
- **Platform Analysis**: iOS vs Android market share and user preferences
- **Feature Prioritization**: Identify must-have vs nice-to-have features

### Deliverables
- Market research report
- User persona documentation
- Competitive analysis matrix
- Feature prioritization matrix`,
    state: TaskStateValue.COMPLETED,
    complexity: 6,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: designerAgent.id,
    dependencies: [],
    dependents: [MOBILE_TASK_IDS.WIREFRAMES, MOBILE_TASK_IDS.DESIGN_SYSTEM],
    createdAt: new Date('2024-01-15T09:00:00Z'),
    updatedAt: new Date('2024-01-17T16:30:00Z'),
    completedAt: new Date('2024-01-17T16:30:00Z'),
    position: { x: 100, y: 100 }
  },

  {
    id: MOBILE_TASK_IDS.WIREFRAMES,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Create Wireframes & Prototypes",
    description: `## Mobile App Wireframes & Prototypes

### Scope
Design low-fidelity wireframes and interactive prototypes for all major user flows.

### User Flows to Cover
- **Onboarding Process**: First-time user experience
- **Authentication Flow**: Login, registration, password reset
- **Main Navigation**: Tab-based navigation structure
- **Core Features**: Key functionality screens
- **Settings & Profile**: User configuration screens

### Deliverables
- Figma wireframes (mobile-specific)
- Interactive prototypes
- User flow diagrams
- Responsive design guidelines`,
    state: TaskStateValue.COMPLETED,
    complexity: 7,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: designerAgent.id,
    dependencies: [MOBILE_TASK_IDS.RESEARCH],
    dependents: [MOBILE_TASK_IDS.DESIGN_SYSTEM, MOBILE_TASK_IDS.TECH_STACK],
    createdAt: new Date('2024-01-15T09:15:00Z'),
    updatedAt: new Date('2024-01-18T14:20:00Z'),
    completedAt: new Date('2024-01-18T14:20:00Z'),
    position: { x: 400, y: 100 }
  },

  {
    id: MOBILE_TASK_IDS.DESIGN_SYSTEM,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Design System & Component Library",
    description: `## Mobile Design System

### Components to Create
- **Typography Scale**: Mobile-optimized text hierarchy
- **Color Palette**: Dark mode compatible colors
- **Spacing System**: Platform-specific spacing guidelines
- **Component Library**: Buttons, forms, cards, navigation
- **Icon System**: Platform-consistent iconography
- **Gesture System**: Touch interaction patterns

### Deliverables
- Complete Figma design system
- Design token specifications
- Component documentation
- Platform guidelines (iOS/Android)`,
    state: TaskStateValue.COMPLETED,
    complexity: 8,
    depth: 0,
    estimate: 2160, // 36 hours
    assignedAgent: designerAgent.id,
    dependencies: [MOBILE_TASK_IDS.RESEARCH, MOBILE_TASK_IDS.WIREFRAMES],
    dependents: [MOBILE_TASK_IDS.TECH_STACK],
    createdAt: new Date('2024-01-15T09:30:00Z'),
    updatedAt: new Date('2024-01-19T10:15:00Z'),
    completedAt: new Date('2024-01-19T10:15:00Z'),
    position: { x: 700, y: 100 }
  },

  // Phase 2: Architecture
  {
    id: MOBILE_TASK_IDS.TECH_STACK,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Technology Stack Selection",
    description: `## Mobile App Technology Stack

### Framework Selection
- **Cross-platform**: React Native vs Flutter vs Native
- **UI Library**: React Navigation vs React Native Navigation
- **State Management**: Redux vs Context API vs Zustand
- **Backend Integration**: REST vs GraphQL
- **Testing Framework**: Jest vs Detox

### Infrastructure
- **Development Environment**: IDE setup and tooling
- **CI/CD Pipeline**: Fastlane and GitHub Actions
- **Analytics**: Firebase Analytics integration
- **Push Notifications**: Firebase Cloud Messaging

### Deliverables
- Tech stack documentation
- Development environment setup guide
- CI/CD pipeline configuration
- Performance benchmarks`,
    state: TaskStateValue.IN_PROGRESS,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.WIREFRAMES, MOBILE_TASK_IDS.DESIGN_SYSTEM],
    dependents: [MOBILE_TASK_IDS.API_INTEGRATION, MOBILE_TASK_IDS.STATE_MANAGEMENT],
    createdAt: new Date('2024-01-16T09:00:00Z'),
    updatedAt: new Date('2024-01-20T10:15:00Z'),
    position: { x: 1000, y: 100 }
  }
];