/**
 * E-Commerce Project Tasks Part 2 (Restliche 12 Tasks)
 * Go-Model konform: pkg/tools/project/shared/shared.go
 * Max 500 Zeilen pro Datei!
 */

import { Task, TaskState } from '../types';
import { ecommerceProject } from './projects';
import { designerAgent, frontendAgent, backendAgent, qaAgent, devopsAgent } from './agents';
import { ECOM_TASK_IDS } from './ecommerce-tasks';

export const ecommerceTasksPart2: Task[] = [
  // Phase 2: Architecture
  {
    id: ECOM_TASK_IDS.API_DESIGN,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "API Architecture & Documentation",
    description: `## API Architecture & Documentation

### Scope
Design RESTful API endpoints, define data models, and create comprehensive API documentation.

### Key Components
- **Authentication API**: Login, register, password reset
- **Product API**: CRUD operations, search, filtering
- **Order API**: Cart management, checkout, order history
- **User API**: Profile management, preferences
- **Admin API**: Dashboard, analytics, user management

### Deliverables
- OpenAPI 3.0 specification
- API documentation site
- Data model diagrams
- Authentication flow documentation`,
    state: TaskState.COMPLETED,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: backendAgent.id,
    dependencies: [ECOM_TASK_IDS.RESEARCH],
    dependents: [ECOM_TASK_IDS.DATABASE_SCHEMA, ECOM_TASK_IDS.BACKEND_SETUP],
    createdAt: new Date('2024-01-15T10:00:00Z'),
    updatedAt: new Date('2024-01-20T11:30:00Z'),
    completedAt: new Date('2024-01-20T11:30:00Z'),
    position: { x: 400, y: 300 }
  },

  {
    id: ECOM_TASK_IDS.DATABASE_SCHEMA,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Database Schema Design",
    description: `## Database Schema Design

### Tables to Design
- **Users**: Authentication, profiles, preferences
- **Products**: Catalog, variants, inventory
- **Orders**: Cart, checkout, order history
- **Categories**: Product organization
- **Reviews**: User feedback system
- **Analytics**: Tracking and metrics

### Requirements
- PostgreSQL optimization
- Proper indexing strategy
- Migration scripts
- Data validation rules
- Performance considerations`,
    state: TaskState.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: backendAgent.id,
    dependencies: [ECOM_TASK_IDS.API_DESIGN],
    dependents: [ECOM_TASK_IDS.BACKEND_SETUP],
    createdAt: new Date('2024-01-15T10:15:00Z'),
    updatedAt: new Date('2024-01-20T11:30:00Z'),
    position: { x: 700, y: 300 }
  },

  {
    id: ECOM_TASK_IDS.FRONTEND_SETUP,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Frontend Project Setup",
    description: `## Frontend Project Setup

### Technology Stack
- **React 18** with TypeScript
- **Vite** for build tooling
- **Tailwind CSS** for styling
- **React Query** for data fetching
- **React Router** for navigation
- **Zustand** for state management

### Configuration
- ESLint + Prettier setup
- Husky pre-commit hooks
- Jest + Testing Library
- Storybook for components
- CI/CD pipeline integration`,
    state: TaskState.COMPLETED,
    complexity: 5,
    depth: 0,
    estimate: 480, // 8 hours
    assignedAgent: frontendAgent.id,
    dependencies: [ECOM_TASK_IDS.WIREFRAMES, ECOM_TASK_IDS.DESIGN_SYSTEM],
    dependents: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.PRODUCT_CATALOG],
    createdAt: new Date('2024-01-16T09:00:00Z'),
    updatedAt: new Date('2024-01-18T17:00:00Z'),
    completedAt: new Date('2024-01-18T17:00:00Z'),
    position: { x: 1000, y: 100 }
  },

  {
    id: ECOM_TASK_IDS.BACKEND_SETUP,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Backend Infrastructure Setup",
    description: `## Backend Infrastructure Setup

### Technology Stack
- **Node.js** with Express.js
- **TypeScript** for type safety
- **PostgreSQL** database
- **Redis** for caching
- **JWT** for authentication
- **Docker** for containerization

### Infrastructure
- Database connection pooling
- Middleware configuration
- Error handling setup
- Logging and monitoring
- Security headers and CORS`,
    state: TaskState.COMPLETED,
    complexity: 6,
    depth: 0,
    estimate: 720, // 12 hours
    assignedAgent: backendAgent.id,
    dependencies: [ECOM_TASK_IDS.API_DESIGN, ECOM_TASK_IDS.DATABASE_SCHEMA],
    dependents: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.PRODUCT_CATALOG],
    createdAt: new Date('2024-01-16T09:30:00Z'),
    updatedAt: new Date('2024-01-19T16:30:00Z'),
    completedAt: new Date('2024-01-19T16:30:00Z'),
    position: { x: 1000, y: 300 }
  },

  // Phase 3: Core Development
  {
    id: ECOM_TASK_IDS.USER_AUTH,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "User Authentication System",
    description: `## User Authentication System

### Features
- **Registration**: Email verification, password strength
- **Login**: JWT tokens, remember me option
- **Password Reset**: Secure token-based reset
- **Profile Management**: Update details, preferences
- **Session Management**: Auto-logout, concurrent sessions

### Security
- Password hashing with bcrypt
- Rate limiting for auth endpoints
- CSRF protection
- Secure cookie handling
- Two-factor authentication (optional)`,
    state: TaskState.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: backendAgent.id,
    dependencies: [ECOM_TASK_IDS.FRONTEND_SETUP, ECOM_TASK_IDS.BACKEND_SETUP],
    dependents: [ECOM_TASK_IDS.PRODUCT_CATALOG, ECOM_TASK_IDS.SHOPPING_CART],
    createdAt: new Date('2024-01-17T10:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1300, y: 200 }
  },

  {
    id: ECOM_TASK_IDS.PRODUCT_CATALOG,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Product Catalog & Search",
    description: `## Product Catalog & Search

### Core Features
- **Product Listing**: Grid/list views, pagination
- **Advanced Search**: Filters, sorting, faceted search
- **Product Details**: Images, descriptions, variants
- **Categories**: Hierarchical navigation
- **Recommendations**: Related products, recently viewed

### Technical Implementation
- Elasticsearch for search
- Image optimization and CDN
- SEO-friendly URLs
- Performance optimization
- Mobile-responsive design`,
    state: TaskState.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: frontendAgent.id,
    dependencies: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.FRONTEND_SETUP],
    dependents: [ECOM_TASK_IDS.SHOPPING_CART],
    createdAt: new Date('2024-01-17T11:00:00Z'),
    updatedAt: new Date('2024-01-20T14:30:00Z'),
    position: { x: 1600, y: 100 }
  }
];