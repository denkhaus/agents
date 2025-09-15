/**
 * Mobile App Project Tasks (10 Tasks)
 * Go-Model konform: pkg/tools/project/shared/shared.go
 * Max 500 Zeilen pro Datei!
 */

import { Task, TaskState } from '../types';
import { generateUUID } from '../utils/uuid';
import { mobileProject } from './projects';
import { designerAgent, frontendAgent, backendAgent, qaAgent, devopsAgent } from './agents';

// Task IDs für Mobile App Dependencies
export const MOBILE_TASK_IDS = {
  // Phase 1: Planning (2 Tasks)
  REQUIREMENTS: generateUUID(),
  UI_DESIGN: generateUUID(),
  
  // Phase 2: Development (5 Tasks)
  PROJECT_SETUP: generateUUID(),
  AUTHENTICATION: generateUUID(),
  CORE_FEATURES: generateUUID(),
  OFFLINE_SYNC: generateUUID(),
  PUSH_NOTIFICATIONS: generateUUID(),
  
  // Phase 3: Testing & Deployment (3 Tasks)
  TESTING: generateUUID(),
  APP_STORE_PREP: generateUUID(),
  DEPLOYMENT: generateUUID()
};

export const mobileTasks: Task[] = [
  // Phase 1: Planning
  {
    id: MOBILE_TASK_IDS.REQUIREMENTS,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Mobile App Requirements Analysis",
    description: `## Mobile App Requirements Analysis

### Platform Requirements
- **iOS**: Minimum iOS 14.0 support
- **Android**: Minimum API level 21 (Android 5.0)
- **Cross-Platform**: React Native 0.72+
- **Performance**: 60fps animations, <3s startup time

### Feature Requirements
- Offline-first architecture
- Real-time data synchronization
- Push notifications
- Biometric authentication
- Dark mode support
- Accessibility compliance

### Technical Constraints
- App size under 50MB
- Battery optimization
- Network efficiency
- Security compliance`,
    state: TaskState.COMPLETED,
    complexity: 5,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: designerAgent.id,
    dependencies: [],
    dependents: [MOBILE_TASK_IDS.UI_DESIGN, MOBILE_TASK_IDS.PROJECT_SETUP],
    createdAt: new Date('2024-01-10T10:00:00Z'),
    updatedAt: new Date('2024-01-12T16:00:00Z'),
    completedAt: new Date('2024-01-12T16:00:00Z'),
    position: { x: 100, y: 700 }
  },

  {
    id: MOBILE_TASK_IDS.UI_DESIGN,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Mobile UI/UX Design",
    description: `## Mobile UI/UX Design

### Design System
- **Typography**: Mobile-optimized font scales
- **Colors**: Dark/light mode palettes
- **Components**: Native-feeling UI components
- **Icons**: Consistent iconography
- **Spacing**: Touch-friendly spacing system

### Screen Designs
- Onboarding flow
- Authentication screens
- Main navigation
- Core feature screens
- Settings and profile
- Error and loading states

### Interaction Design
- Gesture navigation
- Haptic feedback
- Micro-animations
- Accessibility features`,
    state: TaskState.COMPLETED,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: designerAgent.id,
    dependencies: [MOBILE_TASK_IDS.REQUIREMENTS],
    dependents: [MOBILE_TASK_IDS.PROJECT_SETUP, MOBILE_TASK_IDS.CORE_FEATURES],
    createdAt: new Date('2024-01-10T11:00:00Z'),
    updatedAt: new Date('2024-01-15T17:00:00Z'),
    completedAt: new Date('2024-01-15T17:00:00Z'),
    position: { x: 400, y: 700 }
  },

  // Phase 2: Development
  {
    id: MOBILE_TASK_IDS.PROJECT_SETUP,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "React Native Project Setup",
    description: `## React Native Project Setup

### Development Environment
- **React Native CLI**: Latest stable version
- **TypeScript**: Strict configuration
- **Metro**: Bundler optimization
- **Flipper**: Debugging and profiling
- **Detox**: E2E testing framework

### Project Structure
- Feature-based folder structure
- Shared components library
- Navigation setup (React Navigation)
- State management (Redux Toolkit)
- API client configuration

### Build Configuration
- iOS: Xcode project setup
- Android: Gradle configuration
- Code signing setup
- Environment variables
- CI/CD pipeline foundation`,
    state: TaskState.COMPLETED,
    complexity: 6,
    depth: 0,
    estimate: 720, // 12 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.REQUIREMENTS, MOBILE_TASK_IDS.UI_DESIGN],
    dependents: [MOBILE_TASK_IDS.AUTHENTICATION, MOBILE_TASK_IDS.CORE_FEATURES],
    createdAt: new Date('2024-01-11T09:00:00Z'),
    updatedAt: new Date('2024-01-13T18:00:00Z'),
    completedAt: new Date('2024-01-13T18:00:00Z'),
    position: { x: 700, y: 700 }
  },

  {
    id: MOBILE_TASK_IDS.AUTHENTICATION,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Mobile Authentication System",
    description: `## Mobile Authentication System

### Authentication Methods
- **Email/Password**: Traditional login
- **Biometric**: Face ID, Touch ID, Fingerprint
- **Social Login**: Google, Apple, Facebook
- **Phone Number**: SMS verification
- **Guest Mode**: Limited functionality without account

### Security Features
- Secure token storage (Keychain/Keystore)
- Certificate pinning
- Jailbreak/root detection
- Session management
- Auto-logout on inactivity

### User Experience
- Smooth onboarding flow
- Remember login preferences
- Quick re-authentication
- Error handling and recovery`,
    state: TaskState.IN_PROGRESS,
    complexity: 8,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.PROJECT_SETUP],
    dependents: [MOBILE_TASK_IDS.CORE_FEATURES, MOBILE_TASK_IDS.OFFLINE_SYNC],
    createdAt: new Date('2024-01-12T10:00:00Z'),
    updatedAt: new Date('2024-01-20T15:00:00Z'),
    position: { x: 1000, y: 700 }
  }
];