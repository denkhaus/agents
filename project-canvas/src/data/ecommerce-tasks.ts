/**
 * E-Commerce Project Tasks (15 Tasks)
 * Go-Model konform: pkg/tools/project/shared/shared.go
 * Max 500 Zeilen pro Datei!
 */

import type { Task } from '../types';
import { generateUUID } from '../utils/uuid';
import { ecommerceProject } from './projects';
import { designerAgent } from './agents';

// Import the actual values for runtime usage
import { TaskState as TaskStateValue } from '../types/task.types';

// Task IDs für Dependencies (Go-Model: uuid.UUID)
export const ECOM_TASK_IDS = {
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

export const ecommerceTasks: Task[] = [
  // Phase 1: Research & Planning
  {
    id: ECOM_TASK_IDS.RESEARCH,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Market Research & User Analysis",
    description: `## Market Research & User Analysis

**Objective**: Conduct comprehensive market research and analyze user behavior patterns.

### Key Activities
- **Competitor Analysis**: Research 5-10 major e-commerce platforms
- **User Interviews**: Conduct 15-20 user interviews with target demographics  
- **Analytics Review**: Analyze current platform usage data
- **Trend Analysis**: Identify emerging UX/UI trends in e-commerce

### Deliverables
- Market research report (PDF)
- User persona documentation  
- Competitive analysis matrix
- Recommendations summary`,
    state: TaskStateValue.COMPLETED,
    complexity: 6,
    depth: 0,
    estimate: 2400, // 40 hours in minutes
    assignedAgent: designerAgent.id,
    dependencies: [],
    dependents: [ECOM_TASK_IDS.WIREFRAMES, ECOM_TASK_IDS.DESIGN_SYSTEM],
    createdAt: new Date('2024-01-15T09:00:00Z'),
    updatedAt: new Date('2024-01-17T16:30:00Z'),
    completedAt: new Date('2024-01-17T16:30:00Z'),
    position: { x: 100, y: 100 }
  },

  {
    id: ECOM_TASK_IDS.WIREFRAMES,
    projectId: ecommerceProject.id,
    parentId: undefined,
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

### Deliverables
- Figma wireframes (low-fi)
- Interactive prototypes
- User flow diagrams
- Mobile breakpoint designs`,
    state: TaskStateValue.COMPLETED,
    complexity: 7,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: designerAgent.id,
    dependencies: [ECOM_TASK_IDS.RESEARCH],
    dependents: [ECOM_TASK_IDS.DESIGN_SYSTEM, ECOM_TASK_IDS.FRONTEND_SETUP],
    createdAt: new Date('2024-01-15T09:15:00Z'),
    updatedAt: new Date('2024-01-19T14:20:00Z'),
    completedAt: new Date('2024-01-19T14:20:00Z'),
    position: { x: 400, y: 100 }
  },

  {
    id: ECOM_TASK_IDS.DESIGN_SYSTEM,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Design System & Component Library",
    description: `## Design System Development

### Components to Create
- **Typography Scale**: Headings, body text, captions
- **Color Palette**: Primary, secondary, semantic colors
- **Spacing System**: Consistent margins and padding
- **Component Library**: Buttons, forms, cards, navigation
- **Icon System**: Consistent iconography
- **Layout Grid**: Responsive grid system

### Deliverables
- Complete Figma design system
- Design token specifications
- Component documentation
- Usage guidelines`,
    state: TaskStateValue.IN_PROGRESS,
    complexity: 8,
    depth: 0,
    estimate: 2160, // 36 hours
    assignedAgent: designerAgent.id,
    dependencies: [ECOM_TASK_IDS.RESEARCH, ECOM_TASK_IDS.WIREFRAMES],
    dependents: [ECOM_TASK_IDS.FRONTEND_SETUP],
    createdAt: new Date('2024-01-15T09:30:00Z'),
    updatedAt: new Date('2024-01-20T10:15:00Z'),
    position: { x: 700, y: 100 }
  }

  // Weitere Tasks werden in ecommerce-tasks-part2.ts fortgesetzt
];