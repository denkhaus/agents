/**
 * Mobile App Project Tasks Part 2 (Restliche 6 Tasks)
 * Go-Model konform: pkg/tools/project/shared/shared.go
 * Max 500 Zeilen pro Datei!
 */

import { Task, TaskState } from '../types';
import { mobileProject } from './projects';
import { frontendAgent, backendAgent, qaAgent, devopsAgent } from './agents';
import { MOBILE_TASK_IDS } from './mobile-tasks';

export const mobileTasksPart2: Task[] = [
  {
    id: MOBILE_TASK_IDS.CORE_FEATURES,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Core App Features Implementation",
    description: `## Core App Features Implementation

### Main Features
- **Dashboard**: User overview, quick actions
- **Data Management**: CRUD operations for main entities
- **Search & Filter**: Advanced search capabilities
- **Media Handling**: Image/video upload and display
- **Settings**: User preferences, app configuration

### Navigation
- Tab-based navigation
- Stack navigation for details
- Deep linking support
- Universal links (iOS) / App links (Android)
- Navigation state persistence

### Performance
- Lazy loading of screens
- Image optimization
- Memory management
- Battery optimization`,
    state: TaskState.PENDING,
    complexity: 9,
    depth: 0,
    estimate: 2400, // 40 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.AUTHENTICATION, MOBILE_TASK_IDS.UI_DESIGN],
    dependents: [MOBILE_TASK_IDS.OFFLINE_SYNC, MOBILE_TASK_IDS.PUSH_NOTIFICATIONS],
    createdAt: new Date('2024-01-13T10:00:00Z'),
    updatedAt: new Date('2024-01-20T15:00:00Z'),
    position: { x: 1300, y: 700 }
  },

  {
    id: MOBILE_TASK_IDS.OFFLINE_SYNC,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Offline Mode & Data Synchronization",
    description: `## Offline Mode & Data Synchronization

### Offline Capabilities
- **Local Storage**: SQLite database for offline data
- **Conflict Resolution**: Handle data conflicts on sync
- **Queue Management**: Offline actions queue
- **Background Sync**: Sync when app becomes active
- **Partial Sync**: Incremental data updates

### Data Strategy
- Critical data always cached
- Smart prefetching
- Compression for large datasets
- Encryption for sensitive data
- Cleanup of old cached data

### User Experience
- Offline indicators
- Sync progress feedback
- Graceful degradation
- Error recovery mechanisms`,
    state: TaskState.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.AUTHENTICATION, MOBILE_TASK_IDS.CORE_FEATURES],
    dependents: [MOBILE_TASK_IDS.TESTING],
    createdAt: new Date('2024-01-14T09:00:00Z'),
    updatedAt: new Date('2024-01-20T15:00:00Z'),
    position: { x: 1600, y: 700 }
  },

  {
    id: MOBILE_TASK_IDS.PUSH_NOTIFICATIONS,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Push Notifications System",
    description: `## Push Notifications System

### Notification Types
- **Transactional**: Order updates, account changes
- **Marketing**: Promotions, feature announcements
- **Behavioral**: Reminders, re-engagement
- **Real-time**: Chat messages, live updates
- **Location-based**: Geofenced notifications

### Implementation
- Firebase Cloud Messaging (FCM)
- Apple Push Notification Service (APNs)
- Token management and registration
- Deep linking from notifications
- Rich notifications with images/actions

### User Control
- Notification preferences
- Granular opt-in/opt-out
- Quiet hours settings
- Channel-based management`,
    state: TaskState.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: backendAgent.id,
    dependencies: [MOBILE_TASK_IDS.CORE_FEATURES],
    dependents: [MOBILE_TASK_IDS.TESTING],
    createdAt: new Date('2024-01-15T10:00:00Z'),
    updatedAt: new Date('2024-01-20T15:00:00Z'),
    position: { x: 1900, y: 700 }
  },

  // Phase 3: Testing & Deployment
  {
    id: MOBILE_TASK_IDS.TESTING,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Mobile App Testing Suite",
    description: `## Mobile App Testing Suite

### Testing Strategy
- **Unit Tests**: Business logic, utilities
- **Integration Tests**: API integration, data flow
- **E2E Tests**: Critical user journeys (Detox)
- **Device Testing**: Multiple devices and OS versions
- **Performance Testing**: Memory, battery, network

### Platform-Specific Testing
- iOS: Xcode Instruments, TestFlight
- Android: Android Studio Profiler, Internal Testing
- Cross-platform: Flipper debugging
- Accessibility testing
- Security testing

### Automated Testing
- CI/CD integration
- Automated builds and tests
- Screenshot testing
- Crash reporting setup`,
    state: TaskState.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: qaAgent.id,
    dependencies: [MOBILE_TASK_IDS.OFFLINE_SYNC, MOBILE_TASK_IDS.PUSH_NOTIFICATIONS],
    dependents: [MOBILE_TASK_IDS.APP_STORE_PREP],
    createdAt: new Date('2024-01-16T09:00:00Z'),
    updatedAt: new Date('2024-01-20T15:00:00Z'),
    position: { x: 1600, y: 900 }
  },

  {
    id: MOBILE_TASK_IDS.APP_STORE_PREP,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "App Store Preparation",
    description: `## App Store Preparation

### App Store Assets
- **Icons**: All required sizes and formats
- **Screenshots**: Multiple device sizes
- **App Preview Videos**: Showcase key features
- **Metadata**: Descriptions, keywords, categories
- **Privacy Policy**: GDPR and platform compliance

### iOS App Store
- App Store Connect setup
- TestFlight beta testing
- App Review Guidelines compliance
- In-App Purchase setup (if needed)
- Privacy nutrition labels

### Google Play Store
- Play Console setup
- Internal/Alpha/Beta testing tracks
- Play Store policies compliance
- App signing key management
- Store listing optimization`,
    state: TaskState.PENDING,
    complexity: 5,
    depth: 0,
    estimate: 720, // 12 hours
    assignedAgent: devopsAgent.id,
    dependencies: [MOBILE_TASK_IDS.TESTING],
    dependents: [MOBILE_TASK_IDS.DEPLOYMENT],
    createdAt: new Date('2024-01-17T10:00:00Z'),
    updatedAt: new Date('2024-01-20T15:00:00Z'),
    position: { x: 1900, y: 900 }
  },

  {
    id: MOBILE_TASK_IDS.DEPLOYMENT,
    projectId: mobileProject.id,
    parentId: undefined,
    title: "Mobile App Deployment",
    description: `## Mobile App Deployment

### Release Strategy
- **Phased Rollout**: Gradual release to users
- **A/B Testing**: Feature flags and experiments
- **Monitoring**: Crash reporting, analytics
- **Rollback Plan**: Quick revert if issues arise
- **Communication**: Release notes, user notifications

### CI/CD Pipeline
- Automated builds for releases
- Code signing automation
- Store upload automation
- Version management
- Release branch strategy

### Post-Launch
- Performance monitoring
- User feedback collection
- Crash analysis and fixes
- Feature usage analytics
- Update planning`,
    state: TaskState.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: devopsAgent.id,
    dependencies: [MOBILE_TASK_IDS.APP_STORE_PREP],
    dependents: [],
    createdAt: new Date('2024-01-18T09:00:00Z'),
    updatedAt: new Date('2024-01-20T15:00:00Z'),
    position: { x: 2200, y: 900 }
  }
];