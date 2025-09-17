/**
 * Mobile App Project Tasks Part 2 (15 Tasks)
 * Go-Model konform: pkg/tools/project/shared/shared.go
 * Max 500 Zeilen pro Datei!
 */

import type { Task } from '../types';
import { generateUUID } from '../utils/uuid';
import { mobileProject } from './projects';
import { frontendAgent, qaAgent, devopsAgent } from './agents';
import { MOBILE_TASK_IDS } from './mobile-tasks'; // Import the correct task IDs

// Import the actual values for runtime usage
import { TaskState as TaskStateValue } from '../types/task.types';

// Additional Task IDs
export const MOBILE_TASK_IDS_PART2 = {
  // Phase 2: Architecture (continued)
  DATABASE: generateUUID(),
  CI_CD: generateUUID(),
  
  // Phase 3: Core Development (continued)
  AUTH_IMPLEMENTATION: generateUUID(),
  HOME_IMPLEMENTATION: generateUUID(),
  PROFILE_IMPLEMENTATION: generateUUID(),
  NOTIFICATIONS_IMPLEMENTATION: generateUUID(),
  SETTINGS_IMPLEMENTATION: generateUUID(),
  
  // Phase 4: Testing & Deployment (continued)
  PERFORMANCE_TESTING: generateUUID(),
  SECURITY_AUDIT: generateUUID(),
  MONITORING: generateUUID()
};

export const mobileTasksPart2: Task[] = [
  // Phase 2: Architecture (continued)
  {
    id: MOBILE_TASK_IDS_PART2.DATABASE,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "Local Database Integration",
    description: `## Mobile Database Implementation

### Database Selection
- **SQLite**: Relational database for structured data
- **Realm**: Object database for complex relationships
- **WatermelonDB**: Reactive database for large datasets

### Implementation Strategy
- **Data Models**: Define schema and relationships
- **CRUD Operations**: Create, read, update, delete
- **Migrations**: Schema versioning and updates
- **Performance**: Indexing and query optimization

### Features
- **Offline Storage**: Persistent local data
- **Sync Mechanism**: Server synchronization
- **Encryption**: Data at rest protection
- **Backup**: Automatic backup strategies`,
    state: TaskStateValue.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.TECH_STACK],
    dependents: [MOBILE_TASK_IDS_PART2.CI_CD],
    createdAt: new Date('2024-01-16T10:00:00Z'),
    updatedAt: new Date('2024-01-16T10:00:00Z'),
    position: { x: 1300, y: 100 }
  },

  {
    id: MOBILE_TASK_IDS_PART2.CI_CD,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "CI/CD Pipeline Setup",
    description: `## Mobile CI/CD Implementation

### Pipeline Stages
- **Code Quality**: ESLint, Prettier, TypeScript checks
- **Unit Testing**: Jest and React Native Testing Library
- **Build Process**: iOS and Android build automation
- **Deployment**: TestFlight, Google Play internal testing

### Tools & Services
- **GitHub Actions**: Workflow automation
- **Fastlane**: iOS/Android deployment tools
- **CodeMagic**: Alternative CI/CD platform
- **HockeyApp/App Center**: Distribution platform

### Deliverables
- Automated build scripts
- Test automation configuration
- Deployment workflows
- Release documentation`,
    state: TaskStateValue.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: devopsAgent.id,
    dependencies: [MOBILE_TASK_IDS.TECH_STACK, MOBILE_TASK_IDS_PART2.DATABASE],
    dependents: [],
    createdAt: new Date('2024-01-16T10:15:00Z'),
    updatedAt: new Date('2024-01-16T10:15:00Z'),
    position: { x: 1600, y: 100 }
  },

  // Phase 3: Core Development (continued)
  {
    id: MOBILE_TASK_IDS_PART2.AUTH_IMPLEMENTATION,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "Authentication Module Implementation",
    description: `## Mobile Authentication System

### Authentication Methods
- **Email/Password**: Traditional authentication
- **Social Login**: Google, Apple, Facebook integration
- **Biometric**: Face ID, Touch ID, fingerprint
- **Passwordless**: Magic links and SMS codes

### Security Features
- **Token Management**: Secure storage and refresh
- **Session Handling**: Auto-logout and session expiry
- **Encryption**: Data protection at rest and in transit
- **Rate Limiting**: Brute force protection

### User Experience
- **Onboarding Flow**: Smooth user registration
- **Recovery Options**: Password reset and account recovery
- **Multi-factor**: Optional 2FA support
- **Guest Mode**: Limited functionality without account`,
    state: TaskStateValue.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.TECH_STACK],
    dependents: [MOBILE_TASK_IDS.HOME_SCREEN, MOBILE_TASK_IDS.PROFILE_SCREEN],
    createdAt: new Date('2024-01-17T09:00:00Z'),
    updatedAt: new Date('2024-01-17T09:00:00Z'),
    position: { x: 100, y: 400 }
  },

  {
    id: MOBILE_TASK_IDS_PART2.HOME_IMPLEMENTATION,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "Home Screen Development",
    description: `## Mobile Home Screen Implementation

### Core Components
- **Dashboard Layout**: Personalized content display
- **Navigation Menu**: Quick access to main features
- **Notifications Panel**: Recent alerts and updates
- **Search Functionality**: Quick content discovery
- **Quick Actions**: Frequently used functions

### Implementation Details
- **Responsive Design**: Adaptive layouts for all devices
- **Performance Optimization**: Efficient rendering
- **Accessibility**: VoiceOver and TalkBack support
- **Dark Mode**: Automatic theme switching

### Features
- **Personalization**: User-specific content
- **Widgets**: Customizable home screen elements
- **Offline Support**: Cached content display
- **Analytics**: User interaction tracking`,
    state: TaskStateValue.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.TECH_STACK, MOBILE_TASK_IDS_PART2.AUTH_IMPLEMENTATION],
    dependents: [],
    createdAt: new Date('2024-01-17T09:15:00Z'),
    updatedAt: new Date('2024-01-17T09:15:00Z'),
    position: { x: 400, y: 400 }
  },

  {
    id: MOBILE_TASK_IDS_PART2.PROFILE_IMPLEMENTATION,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "Profile & Settings Implementation",
    description: `## User Profile & Settings System

### Profile Features
- **User Information**: Name, email, avatar
- **Activity History**: Recent actions and interactions
- **Preferences**: Personalization settings
- **Account Management**: Security and privacy options

### Settings Categories
- **General**: Language, theme, notifications
- **Privacy**: Data sharing and visibility
- **Security**: Password, 2FA, login activity
- **Accessibility**: Font size, contrast, voice control

### Implementation
- **Form Validation**: Real-time input validation
- **Image Upload**: Avatar and media handling
- **Data Sync**: Cross-device preference synchronization
- **Export/Import**: User data portability`,
    state: TaskStateValue.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.TECH_STACK, MOBILE_TASK_IDS_PART2.AUTH_IMPLEMENTATION],
    dependents: [],
    createdAt: new Date('2024-01-17T09:30:00Z'),
    updatedAt: new Date('2024-01-17T09:30:00Z'),
    position: { x: 700, y: 400 }
  },

  {
    id: MOBILE_TASK_IDS_PART2.NOTIFICATIONS_IMPLEMENTATION,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "Notifications System Implementation",
    description: `## Mobile Notifications Framework

### Notification Types
- **Push Notifications**: Real-time alerts
- **In-app Alerts**: Contextual notifications
- **Badges**: Visual indicators for unread items
- **Sounds**: Audible notification triggers

### Management Features
- **Category Filtering**: Group notifications by type
- **Priority Settings**: Critical vs informational
- **Scheduling**: Time-based notification delivery
- **History**: Notification archive and search

### Technical Implementation
- **Firebase Integration**: FCM for iOS and Android
- **Local Notifications**: Offline and scheduled alerts
- **Deep Linking**: Navigation from notifications
- **Analytics**: Delivery and engagement tracking`,
    state: TaskStateValue.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.TECH_STACK],
    dependents: [],
    createdAt: new Date('2024-01-17T09:45:00Z'),
    updatedAt: new Date('2024-01-17T09:45:00Z'),
    position: { x: 1000, y: 400 }
  },

  {
    id: MOBILE_TASK_IDS_PART2.SETTINGS_IMPLEMENTATION,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "Advanced Settings Implementation",
    description: `## Advanced Settings & Configuration

### Advanced Features
- **Developer Options**: Debugging and testing tools
- **Network Settings**: Proxy, bandwidth controls
- **Storage Management**: Cache and data cleanup
- **Performance Tuning**: Animation and rendering controls

### Customization Options
- **Theme Editor**: Custom color schemes
- **Layout Preferences**: Grid vs list views
- **Behavior Settings**: Gestures and shortcuts
- **Integration Setup**: Third-party service connections

### Implementation Details
- **Preference Persistence**: Local storage and sync
- **Reset Options**: Factory defaults and partial resets
- **Import/Export**: Configuration sharing
- **User Guidance**: Tooltips and help documentation`,
    state: TaskStateValue.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: frontendAgent.id,
    dependencies: [MOBILE_TASK_IDS.TECH_STACK],
    dependents: [],
    createdAt: new Date('2024-01-17T10:00:00Z'),
    updatedAt: new Date('2024-01-17T10:00:00Z'),
    position: { x: 1300, y: 400 }
  },

  // Phase 4: Testing & Deployment (continued)
  {
    id: MOBILE_TASK_IDS_PART2.PERFORMANCE_TESTING,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "Performance & Battery Testing",
    description: `## Mobile Performance Validation

### Testing Scenarios
- **Load Testing**: Simulate high user activity
- **Battery Usage**: Monitor power consumption
- **Memory Leaks**: Identify resource leaks
- **Network Efficiency**: Optimize data usage
- **Startup Time**: Reduce app launch duration

### Tools & Metrics
- **Xcode Instruments**: iOS profiling tools
- **Android Profiler**: Memory and CPU analysis
- **Firebase Performance**: Real-time monitoring
- **Appetize**: Browser-based device testing

### Optimization Targets
- **Frame Rate**: Maintain 60fps animations
- **Memory Usage**: Under 100MB baseline
- **Battery Impact**: Minimal background activity
- **Network Usage**: Efficient data transfer`,
    state: TaskStateValue.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: qaAgent.id,
    dependencies: [MOBILE_TASK_IDS.LAUNCH],
    dependents: [MOBILE_TASK_IDS_PART2.SECURITY_AUDIT],
    createdAt: new Date('2024-01-18T09:00:00Z'),
    updatedAt: new Date('2024-01-18T09:00:00Z'),
    position: { x: 1600, y: 400 }
  },

  {
    id: MOBILE_TASK_IDS_PART2.SECURITY_AUDIT,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "Security Audit & Compliance",
    description: `## Mobile App Security Assessment

### Audit Areas
- **OWASP Mobile Top 10**: Vulnerability assessment
- **Data Protection**: Encryption and secure storage
- **Authentication**: Token security and session management
- **Network Security**: SSL pinning and certificate validation
- **Code Security**: Obfuscation and tamper detection

### Compliance Requirements
- **GDPR**: Data privacy compliance
- **CCPA**: California privacy rights
- **App Store Guidelines**: Apple and Google policy compliance
- **Industry Standards**: SOC 2, ISO 27001

### Deliverables
- Security audit report
- Remediation plan
- Compliance documentation
- Penetration testing results`,
    state: TaskStateValue.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: qaAgent.id,
    dependencies: [MOBILE_TASK_IDS.LAUNCH, MOBILE_TASK_IDS_PART2.PERFORMANCE_TESTING],
    dependents: [MOBILE_TASK_IDS_PART2.MONITORING],
    createdAt: new Date('2024-01-18T09:15:00Z'),
    updatedAt: new Date('2024-01-18T09:15:00Z'),
    position: { x: 1900, y: 400 }
  },

  {
    id: MOBILE_TASK_IDS_PART2.MONITORING,
    projectId: mobileProject.id,
    parentId: MOBILE_TASK_IDS.TECH_STACK,
    title: "Production Monitoring Setup",
    description: `## Mobile App Monitoring Infrastructure

### Monitoring Layers
- **Crash Reporting**: Real-time crash detection
- **Performance Metrics**: Load times and responsiveness
- **User Experience**: Touch interactions and navigation
- **Business Metrics**: Feature usage and conversion

### Tools & Services
- **Sentry**: Crash reporting and error tracking
- **New Relic**: Mobile application monitoring
- **Firebase Analytics**: User behavior tracking
- **AppDynamics**: End-to-end performance monitoring

### Alerting System
- **Thresholds**: Configurable alert triggers
- **Channels**: Email, Slack, PagerDuty integrations
- **Escalation**: Multi-level alert routing
- **Dashboards**: Real-time operational views

### Deliverables
- Monitoring infrastructure
- Alert configuration
- Operational dashboards
- Incident response procedures`,
    state: TaskStateValue.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: devopsAgent.id,
    dependencies: [MOBILE_TASK_IDS.LAUNCH, MOBILE_TASK_IDS_PART2.SECURITY_AUDIT],
    dependents: [],
    createdAt: new Date('2024-01-18T09:30:00Z'),
    updatedAt: new Date('2024-01-18T09:30:00Z'),
    position: { x: 2200, y: 400 }
  }
];