/**
 * Master Dummy Data - Index
 * Führt alle aufgeteilten Dummy-Daten zusammen
 * Jede Datei max 500 Zeilen!
 */

import { Project } from '../types/project.types';
import { Agent, AgentRole, AgentStatus } from '../types/agent.types';
import { Task, TaskState } from '../types/task.types';
import { generateUUID } from '../utils/uuid';

export * from './projects';
export * from './agents';
export * from './ecommerce-tasks';
export * from './ecommerce-tasks-part2';
export * from './ecommerce-tasks-part3';
export * from './mobile-tasks';
export * from './mobile-tasks-part2';

// Kombinierte Arrays für einfachen Zugriff
import { ecommerceTasks } from './ecommerce-tasks';
import { ecommerceTasksPart2 } from './ecommerce-tasks-part2';
import { ecommerceTasksPart3 } from './ecommerce-tasks-part3';
import { mobileTasks } from './mobile-tasks';
import { mobileTasksPart2 } from './mobile-tasks-part2';

export const allEcommerceTasks = [
  ...ecommerceTasks,
  ...ecommerceTasksPart2,
  ...ecommerceTasksPart3
];

export const allMobileTasks = [
  ...mobileTasks,
  ...mobileTasksPart2
];

export const allTasks = [
  ...allEcommerceTasks,
  ...allMobileTasks
];

// ============================================================================
// PROJECTS (2 Projekte wie im Go-Model definiert)
// ============================================================================

// Feste Project-IDs für konsistente Referenzen
const PROJECT_IDS = {
  ECOMMERCE: "53fd7995-1606-4360-65ae-9b057525c12e",
  MOBILE: "980c9337-0630-41d0-3b95-b71caca7635b"
};

export const masterProjects: Project[] = [
  // Projekt 1: E-Commerce Platform Redesign (Hauptprojekt)
  {
    id: PROJECT_IDS.ECOMMERCE,
    title: "E-Commerce Platform Redesign",
    description: `# E-Commerce Platform Redesign

## Project Overview
Complete redesign and modernization of the existing e-commerce platform with improved UX, performance, and mobile responsiveness.

### Key Objectives
- **User Experience**: Streamlined checkout process and intuitive navigation
- **Performance**: 50% improvement in page load times  
- **Mobile**: Responsive design for all screen sizes
- **Accessibility**: WCAG 2.1 AA compliance

### Success Metrics
- Conversion rate increase by 25%
- Mobile traffic engagement up 40%
- Customer satisfaction score > 4.5/5

### Technical Stack
- Frontend: React 18 + TypeScript
- Backend: Node.js + Express
- Database: PostgreSQL
- Deployment: Docker + Kubernetes`,
    createdAt: new Date('2024-01-15T09:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    totalTasks: 15,
    completedTasks: 4,
    progress: 26.7
  },

  // Projekt 2: Mobile App Development
  {
    id: PROJECT_IDS.MOBILE, 
    title: "Mobile App Development",
    description: `# Mobile App Development

## Project Overview
Native mobile application for iOS and Android platforms with offline capabilities and real-time synchronization.

### Key Features
- **Cross-Platform**: React Native for iOS and Android
- **Offline Mode**: Local data storage with sync
- **Push Notifications**: Real-time user engagement
- **Biometric Auth**: Secure login with fingerprint/face ID

### Technical Requirements
- React Native 0.72+
- TypeScript
- Redux Toolkit
- SQLite for offline storage
- Firebase for push notifications

### Timeline
- Phase 1: Core features (8 weeks)
- Phase 2: Advanced features (4 weeks)
- Phase 3: Testing & deployment (2 weeks)`,
    createdAt: new Date('2024-01-10T10:00:00Z'),
    updatedAt: new Date('2024-01-18T16:45:00Z'),
    totalTasks: 10,
    completedTasks: 3,
    progress: 30.0
  }
];

// ============================================================================
// AGENTS (5 Agents wie im Go-Model definiert)
// ============================================================================

// Feste Agent-IDs für konsistente Referenzen
const AGENT_IDS = {
  DESIGNER: "agent-designer-001",
  FRONTEND_DEV: "agent-frontend-002", 
  BACKEND_DEV: "agent-backend-003",
  QA_ENGINEER: "agent-qa-004",
  DEVOPS: "agent-devops-005"
};

export const masterAgents: Agent[] = [
  {
    id: AGENT_IDS.DESIGNER,
    name: "Design Lead",
    role: AgentRole.DESIGNER,
    description: "Senior UX/UI designer specializing in e-commerce platforms and design systems",
    status: AgentStatus.ONLINE,
    isStreaming: false,
    capabilities: ["wireframing", "prototyping", "design-systems", "user-research", "figma", "accessibility"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T09:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T14:30:00Z')
  },
  {
    id: AGENT_IDS.FRONTEND_DEV,
    name: "Frontend Developer",
    role: AgentRole.CODER,
    description: "React/TypeScript specialist with expertise in modern frontend development and performance optimization",
    status: AgentStatus.BUSY,
    isStreaming: true,
    capabilities: ["react", "typescript", "tailwind", "testing", "performance-optimization", "accessibility"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T09:15:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T14:25:00Z')
  },
  {
    id: AGENT_IDS.BACKEND_DEV,
    name: "Backend Developer",
    role: AgentRole.CODER,
    description: "Node.js/Express expert with database design and API architecture experience",
    status: AgentStatus.ONLINE,
    isStreaming: false,
    capabilities: ["nodejs", "express", "postgresql", "api-design", "security", "performance"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T09:30:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T14:20:00Z')
  },
  {
    id: AGENT_IDS.QA_ENGINEER,
    name: "QA Engineer", 
    role: AgentRole.QA_ENGINEER,
    description: "Quality assurance specialist with automated testing expertise and mobile testing experience",
    status: AgentStatus.IDLE,
    isStreaming: false,
    capabilities: ["jest", "cypress", "test-automation", "performance-testing", "accessibility", "mobile-testing"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T10:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T13:45:00Z')
  },
  {
    id: AGENT_IDS.DEVOPS,
    name: "DevOps Engineer",
    role: AgentRole.DEVOPS,
    description: "Infrastructure and deployment automation specialist with cloud and mobile deployment expertise",
    status: AgentStatus.ONLINE,
    isStreaming: false,
    capabilities: ["docker", "kubernetes", "ci-cd", "monitoring", "cloud-infrastructure", "mobile-deployment"],
    currentTasks: [],
    createdAt: new Date('2024-01-10T10:15:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    lastActiveAt: new Date('2024-01-20T14:10:00Z')
  }
];

// ============================================================================
// TASKS - Projekt 1: E-Commerce Platform (15 Tasks)
// ============================================================================

// Task IDs f√ºr Dependencies (Go-Model: uuid.UUID)
const ECOM_TASK_IDS = {
  // Phase 1: Research & Planning (3 Tasks)
  RESEARCH: generateUUID(),
  WIREFRAMES: generateUUID(), 
  DESIGN_SYSTEM: generateUUID(),
  
  // Phase 2: Architecture (4 Tasks)
  API_DESIGN: generateUUID(),
  DATABASE_SCHEMA: generateUUID(),
  FRONTEND_SETUP: generateUUID(),
  BACKEND_SETUP: generateUUID(),
  
  // Phase 3: Core Development (5 Tasks)
  USER_AUTH: generateUUID(),
  PRODUCT_CATALOG: generateUUID(),
  SHOPPING_CART: generateUUID(),
  PAYMENT_INTEGRATION: generateUUID(),
  ADMIN_PANEL: generateUUID(),
  
  // Phase 4: Testing & Deployment (3 Tasks)
  TESTING_SETUP: generateUUID(),
  E2E_TESTS: generateUUID(),
  DEPLOYMENT: generateUUID()
};

export const ecommerceProjectTasks: Task[] = [
  // Phase 1: Research & Planning (3 Tasks)
  {
    id: ECOM_TASK_IDS.RESEARCH,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: undefined,
    title: "Market Research & User Analysis",
    description: `## Market Research & User Analysis

**Objective**: Conduct comprehensive market research and analyze user behavior patterns to inform design decisions.

### Key Activities
- **Competitor Analysis**: Research 5-10 major e-commerce platforms
- **User Interviews**: Conduct 15-20 user interviews with target demographics  
- **Analytics Review**: Analyze current platform usage data
- **Trend Analysis**: Identify emerging UX/UI trends in e-commerce

### Deliverables
- Market research report (PDF)
- User persona documentation  
- Competitive analysis matrix
- Recommendations summary

### Success Criteria
- Complete analysis of top 10 competitors
- 20+ user interviews conducted
- Clear user personas defined
- Actionable recommendations documented`,
    state: TaskState.COMPLETED,
    complexity: 6,
    depth: 0,
    estimate: 2400, // 40 hours in minutes
    assignedAgent: AGENT_IDS.DESIGNER, // Designer
    dependencies: [],
    dependents: [ECOM_TASK_IDS.WIREFRAMES, ECOM_TASK_IDS.DESIGN_SYSTEM],
    createdAt: new Date('2024-01-15T09:00:00Z'),
    updatedAt: new Date('2024-01-17T16:30:00Z'),
    completedAt: new Date('2024-01-17T16:30:00Z'),
    position: { x: 100, y: 100 }
  },

  {
    id: ECOM_TASK_IDS.WIREFRAMES,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.RESEARCH,
    title: "Create Wireframes & Prototypes",
    description: `## Wireframes & Interactive Prototypes

### Scope
Design low-fidelity wireframes and interactive prototypes for all major user flows.

### User Flows to Cover
- **Homepage & Navigation**
- **Product Discovery & Search**
- **Product Detail Pages**
- **Shopping Cart & Checkout**
- **User Account Management**
- **Mobile Responsive Layouts**

### Tools & Deliverables
- Figma wireframes (low-fi)
- Interactive prototypes
- User flow diagrams
- Mobile breakpoint designs

### Acceptance Criteria
- All major user flows wireframed
- Interactive prototypes for key flows
- Mobile-first responsive design
- Stakeholder approval obtained`,
    state: TaskState.COMPLETED,
    complexity: 7,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: AGENT_IDS.DESIGNER, // Designer
    dependencies: [ECOM_TASK_IDS.RESEARCH],
    dependents: [ECOM_TASK_IDS.DESIGN_SYSTEM, ECOM_TASK_IDS.FRONTEND_SETUP],
    createdAt: new Date('2024-01-15T09:15:00Z'),
    updatedAt: new Date('2024-01-19T14:20:00Z'),
    completedAt: new Date('2024-01-19T14:20:00Z'),
    position: { x: 400, y: 100 }
  },

  {
    id: ECOM_TASK_IDS.DESIGN_SYSTEM,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.WIREFRAMES,
    title: "Design System & Component Library",
    description: `## Design System Development

### Components to Create
- **Typography Scale**: Headings, body text, captions
- **Color Palette**: Primary, secondary, semantic colors
- **Spacing System**: Consistent margins and padding
- **Component Library**: Buttons, forms, cards, navigation
- **Icon System**: Consistent iconography
- **Layout Grid**: Responsive grid system

### Technical Requirements
- Figma component library
- Design tokens (JSON)
- Documentation site
- Accessibility guidelines

### Deliverables
- Complete Figma design system
- Design token specifications
- Component documentation
- Usage guidelines`,
    state: TaskState.IN_PROGRESS,
    complexity: 8,
    depth: 0,
    estimate: 2160, // 36 hours
    assignedAgent: AGENT_IDS.DESIGNER, // Designer
    dependencies: [ECOM_TASK_IDS.RESEARCH, ECOM_TASK_IDS.WIREFRAMES],
    dependents: [ECOM_TASK_IDS.FRONTEND_SETUP],
    createdAt: new Date('2024-01-15T09:30:00Z'),
    updatedAt: new Date('2024-01-20T10:15:00Z'),
    position: { x: 700, y: 100 }
  },

  // Phase 2: Architecture (4 Tasks)
  {
    id: ECOM_TASK_IDS.API_DESIGN,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.RESEARCH,
    title: "API Architecture & Documentation",
    description: `## RESTful API Design

### API Endpoints to Design
- **Authentication**: Login, register, refresh tokens
- **Products**: CRUD operations, search, filtering
- **Cart**: Add/remove items, quantity updates
- **Orders**: Create, track, history
- **Users**: Profile management, preferences
- **Admin**: Dashboard, analytics, management

### Technical Specifications
- OpenAPI 3.0 documentation
- Rate limiting strategies
- Caching mechanisms
- Error handling standards
- Security considerations

### Deliverables
- Complete API specification
- Postman collection
- Security documentation
- Performance guidelines`,
    state: TaskState.IN_PROGRESS,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: AGENT_IDS.BACKEND_DEV, // Backend Developer
    dependencies: [ECOM_TASK_IDS.RESEARCH],
    dependents: [ECOM_TASK_IDS.DATABASE_SCHEMA, ECOM_TASK_IDS.BACKEND_SETUP],
    createdAt: new Date('2024-01-15T10:00:00Z'),
    updatedAt: new Date('2024-01-20T11:30:00Z'),
    position: { x: 400, y: 300 }
  },

  {
    id: ECOM_TASK_IDS.DATABASE_SCHEMA,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.API_DESIGN,
    title: "Database Schema Design",
    description: `## PostgreSQL Database Design

### Core Tables
- **Users**: Authentication, profiles, preferences
- **Products**: Catalog, variants, inventory
- **Categories**: Hierarchical product organization
- **Orders**: Transactions, line items, status
- **Cart**: Session-based shopping cart
- **Reviews**: Product ratings and comments

### Technical Requirements
- Normalized schema design
- Proper indexing strategy
- Foreign key constraints
- Migration scripts
- Performance optimization

### Deliverables
- Complete ERD diagram
- SQL migration files
- Indexing strategy
- Data seeding scripts`,
    state: TaskState.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: AGENT_IDS.BACKEND_DEV, // Backend Developer
    dependencies: [ECOM_TASK_IDS.API_DESIGN],
    dependents: [ECOM_TASK_IDS.BACKEND_SETUP],
    createdAt: new Date('2024-01-15T10:15:00Z'),
    updatedAt: new Date('2024-01-20T11:30:00Z'),
    position: { x: 700, y: 300 }
  },

  {
    id: ECOM_TASK_IDS.FRONTEND_SETUP,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.DESIGN_SYSTEM,
    title: "Frontend Project Setup",
    description: `## React + TypeScript Setup

### Technology Stack
- **Framework**: React 18 with TypeScript
- **Styling**: Tailwind CSS + shadcn/ui
- **State Management**: Zustand + React Query
- **Routing**: React Router v6
- **Testing**: Jest + React Testing Library
- **Build**: Vite with optimizations

### Configuration
- ESLint + Prettier setup
- Husky pre-commit hooks
- CI/CD pipeline configuration
- Environment management
- Performance monitoring

### Deliverables
- Complete project scaffold
- Development environment
- Build and deployment scripts
- Code quality tools`,
    state: TaskState.COMPLETED,
    complexity: 5,
    depth: 0,
    estimate: 480, // 8 hours
    assignedAgent: AGENT_IDS.FRONTEND_DEV, // Frontend Developer
    dependencies: [ECOM_TASK_IDS.WIREFRAMES, ECOM_TASK_IDS.DESIGN_SYSTEM],
    dependents: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.PRODUCT_CATALOG],
    createdAt: new Date('2024-01-16T09:00:00Z'),
    updatedAt: new Date('2024-01-18T17:00:00Z'),
    completedAt: new Date('2024-01-18T17:00:00Z'),
    position: { x: 1000, y: 100 }
  },

  {
    id: ECOM_TASK_IDS.BACKEND_SETUP,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.DATABASE_SCHEMA,
    title: "Backend Infrastructure Setup",
    description: `## Node.js + Express Setup

### Technology Stack
- **Runtime**: Node.js 18+ with TypeScript
- **Framework**: Express.js with middleware
- **Database**: PostgreSQL with Prisma ORM
- **Authentication**: JWT with refresh tokens
- **Validation**: Zod schema validation
- **Documentation**: Swagger/OpenAPI

### Infrastructure
- Docker containerization
- Environment configuration
- Logging and monitoring
- Error handling middleware
- Security headers and CORS

### Deliverables
- Complete server setup
- Database connection
- Authentication middleware
- API documentation
- Docker configuration`,
    state: TaskState.COMPLETED,
    complexity: 6,
    depth: 0,
    estimate: 720, // 12 hours
    assignedAgent: AGENT_IDS.BACKEND_DEV, // Backend Developer
    dependencies: [ECOM_TASK_IDS.API_DESIGN, ECOM_TASK_IDS.DATABASE_SCHEMA],
    dependents: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.PRODUCT_CATALOG],
    createdAt: new Date('2024-01-16T09:30:00Z'),
    updatedAt: new Date('2024-01-19T16:30:00Z'),
    completedAt: new Date('2024-01-19T16:30:00Z'),
    position: { x: 1000, y: 300 }
  },

  // Phase 3: Core Development (5 Tasks)
  {
    id: ECOM_TASK_IDS.USER_AUTH,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.BACKEND_SETUP,
    title: "User Authentication System",
    description: `## Secure Authentication Implementation

### Features to Implement
- **Registration**: Email verification, password strength
- **Login**: Secure session management
- **Password Reset**: Email-based recovery
- **Profile Management**: User data updates
- **Session Handling**: JWT tokens, refresh logic
- **Social Login**: Google, Facebook integration

### Security Requirements
- Password hashing (bcrypt)
- Rate limiting on auth endpoints
- CSRF protection
- Secure cookie handling
- Account lockout policies

### Deliverables
- Complete auth API endpoints
- Frontend auth components
- Security testing
- Documentation`,
    state: TaskState.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: AGENT_IDS.BACKEND_DEV, // Backend Developer
    dependencies: [ECOM_TASK_IDS.FRONTEND_SETUP, ECOM_TASK_IDS.BACKEND_SETUP],
    dependents: [ECOM_TASK_IDS.PRODUCT_CATALOG, ECOM_TASK_IDS.SHOPPING_CART],
    createdAt: new Date('2024-01-17T10:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1300, y: 200 }
  },

  {
    id: ECOM_TASK_IDS.PRODUCT_CATALOG,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.FRONTEND_SETUP,
    title: "Product Catalog & Search",
    description: `## Product Management System

### Core Features
- **Product Listing**: Grid/list views with pagination
- **Search & Filtering**: Full-text search, faceted filters
- **Product Details**: Images, descriptions, variants
- **Categories**: Hierarchical navigation
- **Inventory**: Stock tracking, availability
- **Reviews**: Customer ratings and comments

### Technical Implementation
- Elasticsearch for search
- Image optimization and CDN
- Caching strategies
- SEO optimization
- Mobile responsiveness

### Deliverables
- Product catalog API
- Search functionality
- Frontend components
- Admin management interface`,
    state: TaskState.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: AGENT_IDS.FRONTEND_DEV, // Frontend Developer
    dependencies: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.FRONTEND_SETUP],
    dependents: [ECOM_TASK_IDS.SHOPPING_CART],
    createdAt: new Date('2024-01-17T11:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1600, y: 100 }
  },

  {
    id: ECOM_TASK_IDS.SHOPPING_CART,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.PRODUCT_CATALOG,
    title: "Shopping Cart & Checkout",
    description: `## E-commerce Cart System

### Cart Features
- **Add/Remove Items**: Product variants, quantities
- **Persistent Cart**: Session and user-based storage
- **Price Calculation**: Taxes, discounts, shipping
- **Checkout Flow**: Multi-step process
- **Guest Checkout**: No registration required
- **Save for Later**: Wishlist functionality

### Checkout Process
- Shipping address management
- Payment method selection
- Order review and confirmation
- Email notifications
- Order tracking

### Deliverables
- Cart management API
- Checkout flow components
- Payment integration prep
- Order confirmation system`,
    state: TaskState.PENDING,
    complexity: 9,
    depth: 0,
    estimate: 2160, // 36 hours
    assignedAgent: AGENT_IDS.FRONTEND_DEV, // Frontend Developer
    dependencies: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.PRODUCT_CATALOG],
    dependents: [ECOM_TASK_IDS.PAYMENT_INTEGRATION],
    createdAt: new Date('2024-01-17T12:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1900, y: 200 }
  },

  {
    id: ECOM_TASK_IDS.PAYMENT_INTEGRATION,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.SHOPPING_CART,
    title: "Payment Gateway Integration",
    description: `## Secure Payment Processing

### Payment Methods
- **Credit Cards**: Visa, MasterCard, Amex
- **Digital Wallets**: PayPal, Apple Pay, Google Pay
- **Bank Transfers**: ACH, wire transfers
- **Buy Now Pay Later**: Klarna, Afterpay
- **Cryptocurrency**: Bitcoin, Ethereum (optional)

### Security & Compliance
- PCI DSS compliance
- 3D Secure authentication
- Fraud detection
- Webhook handling
- Refund processing

### Technical Implementation
- Stripe/PayPal SDK integration
- Secure tokenization
- Transaction logging
- Error handling
- Testing with sandbox

### Deliverables
- Payment API endpoints
- Frontend payment forms
- Webhook handlers
- Security documentation`,
    state: TaskState.BLOCKED,
    complexity: 8,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: AGENT_IDS.BACKEND_DEV, // Backend Developer
    dependencies: [ECOM_TASK_IDS.SHOPPING_CART],
    dependents: [ECOM_TASK_IDS.E2E_TESTS],
    createdAt: new Date('2024-01-18T09:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 2200, y: 300 }
  },

  {
    id: ECOM_TASK_IDS.ADMIN_PANEL,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.USER_AUTH,
    title: "Admin Dashboard & Management",
    description: `## Administrative Interface

### Dashboard Features
- **Analytics**: Sales, traffic, conversion metrics
- **Product Management**: CRUD operations, inventory
- **Order Management**: Processing, fulfillment, refunds
- **User Management**: Customer support, account management
- **Content Management**: Pages, banners, promotions
- **Settings**: Configuration, integrations

### Technical Requirements
- Role-based access control
- Real-time data updates
- Export functionality
- Audit logging
- Mobile-responsive design

### Deliverables
- Admin dashboard interface
- Management API endpoints
- User permission system
- Analytics integration`,
    state: TaskState.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1680, // 28 hours
    assignedAgent: AGENT_IDS.FRONTEND_DEV, // Frontend Developer
    dependencies: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.PRODUCT_CATALOG],
    dependents: [ECOM_TASK_IDS.E2E_TESTS],
    createdAt: new Date('2024-01-18T10:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1600, y: 400 }
  },

  // Phase 4: Testing & Deployment (3 Tasks)
  {
    id: ECOM_TASK_IDS.TESTING_SETUP,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.ADMIN_PANEL,
    title: "Testing Framework Setup",
    description: `## Comprehensive Testing Strategy

### Testing Types
- **Unit Tests**: Jest for components and functions
- **Integration Tests**: API endpoint testing
- **E2E Tests**: Cypress for user flows
- **Performance Tests**: Load testing with Artillery
- **Security Tests**: OWASP vulnerability scanning
- **Accessibility Tests**: axe-core integration

### Test Infrastructure
- CI/CD pipeline integration
- Test data management
- Mock services setup
- Coverage reporting
- Automated test execution

### Deliverables
- Complete testing framework
- Test utilities and helpers
- CI/CD integration
- Coverage reporting setup`,
    state: TaskState.PENDING,
    complexity: 4,
    depth: 0,
    estimate: 360, // 6 hours
    assignedAgent: AGENT_IDS.QA_ENGINEER, // QA Engineer
    dependencies: [ECOM_TASK_IDS.FRONTEND_SETUP],
    dependents: [ECOM_TASK_IDS.E2E_TESTS],
    createdAt: new Date('2024-01-18T10:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1300, y: 500 }
  },

  {
    id: ECOM_TASK_IDS.E2E_TESTS,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.TESTING_SETUP,
    title: "End-to-End Testing Implementation",
    description: `## Critical User Flow Testing

### Test Scenarios
- **User Registration & Login**: Account creation, authentication
- **Product Discovery**: Search, filtering, navigation
- **Shopping Flow**: Add to cart, checkout, payment
- **Order Management**: Tracking, history, returns
- **Admin Functions**: Product management, order processing
- **Error Handling**: Network failures, validation errors

### Cross-Browser Testing
- Chrome, Firefox, Safari, Edge
- Mobile browsers (iOS Safari, Chrome Mobile)
- Different screen sizes and resolutions
- Accessibility compliance testing

### Performance Testing
- Page load times
- API response times
- Database query optimization
- CDN effectiveness

### Deliverables
- Complete E2E test suite
- Performance benchmarks
- Cross-browser compatibility report
- Accessibility audit`,
    state: TaskState.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: AGENT_IDS.QA_ENGINEER, // QA Engineer
    dependencies: [ECOM_TASK_IDS.TESTING_SETUP, ECOM_TASK_IDS.PAYMENT_INTEGRATION, ECOM_TASK_IDS.ADMIN_PANEL],
    dependents: [ECOM_TASK_IDS.DEPLOYMENT],
    createdAt: new Date('2024-01-18T12:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1900, y: 500 }
  },

  {
    id: ECOM_TASK_IDS.DEPLOYMENT,
    projectId: PROJECT_IDS.ECOMMERCE,
    parentId: ECOM_TASK_IDS.E2E_TESTS,
    title: "Production Deployment & CI/CD",
    description: `## Production Infrastructure

### Deployment Strategy
- **Containerization**: Docker for all services
- **Orchestration**: Kubernetes cluster setup
- **Load Balancing**: NGINX with SSL termination
- **Database**: PostgreSQL with read replicas
- **CDN**: CloudFront for static assets
- **Monitoring**: Prometheus + Grafana

### CI/CD Pipeline
- GitHub Actions workflows
- Automated testing on PR
- Staging environment deployment
- Production deployment with rollback
- Database migration automation

### Security & Monitoring
- SSL certificates (Let's Encrypt)
- Security headers configuration
- Log aggregation (ELK stack)
- Error tracking (Sentry)
- Performance monitoring (New Relic)

### Deliverables
- Production infrastructure
- CI/CD pipeline
- Monitoring dashboards
- Deployment documentation`,
    state: TaskState.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: AGENT_IDS.DEVOPS, // DevOps Engineer
    dependencies: [ECOM_TASK_IDS.E2E_TESTS],
    dependents: [],
    createdAt: new Date('2024-01-19T09:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 2200, y: 600 }
  }
];