/**
 * E-Commerce Project Tasks Part 3 (Finale 6 Tasks)
 * Go-Model konform: pkg/tools/project/shared/shared.go
 * Max 500 Zeilen pro Datei!
 */

import type { Task } from "../types";
import { ecommerceProject } from "./projects";
import { frontendAgent, backendAgent, qaAgent, devopsAgent } from "./agents";
import { ECOM_TASK_IDS } from "./ecommerce-tasks";

// Import the actual values for runtime usage
import { TaskState as TaskStateValue } from "../types/task.types";

export const ecommerceTasksPart3: Task[] = [
  {
    id: ECOM_TASK_IDS.SHOPPING_CART,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Shopping Cart & Checkout",
    description: `## Shopping Cart & Checkout

### Cart Features
- **Add/Remove Items**: Quantity management
- **Persistent Cart**: Local storage + user account sync
- **Price Calculation**: Taxes, shipping, discounts
- **Guest Checkout**: No registration required
- **Save for Later**: Wishlist functionality

### Checkout Process
- Multi-step checkout flow
- Address management
- Shipping options
- Order summary and confirmation
- Email notifications`,
    state: TaskStateValue.PENDING,
    complexity: 9,
    depth: 0,
    estimate: 2160, // 36 hours
    assignedAgent: frontendAgent.id,
    dependencies: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.PRODUCT_CATALOG],
    dependents: [ECOM_TASK_IDS.PAYMENT_INTEGRATION],
    createdAt: new Date("2024-01-17T12:00:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    position: { x: 1900, y: 200 },
  },

  {
    id: ECOM_TASK_IDS.PAYMENT_INTEGRATION,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Payment Gateway Integration",
    description: `## Payment Gateway Integration

### Payment Methods
- **Credit/Debit Cards**: Stripe integration
- **PayPal**: Express checkout
- **Digital Wallets**: Apple Pay, Google Pay
- **Bank Transfer**: SEPA for EU customers
- **Buy Now Pay Later**: Klarna integration

### Security & Compliance
- PCI DSS compliance
- 3D Secure authentication
- Fraud detection
- Secure tokenization
- Refund processing`,
    state: TaskStateValue.BLOCKED,
    complexity: 8,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: backendAgent.id,
    dependencies: [ECOM_TASK_IDS.SHOPPING_CART],
    dependents: [ECOM_TASK_IDS.E2E_TESTS],
    createdAt: new Date("2024-01-18T09:00:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    position: { x: 2200, y: 300 },
  },

  {
    id: ECOM_TASK_IDS.ADMIN_PANEL,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Admin Dashboard & Management",
    description: `## Admin Dashboard & Management

### Dashboard Features
- **Analytics**: Sales, traffic, conversion metrics
- **Product Management**: CRUD, inventory, categories
- **Order Management**: Processing, fulfillment, returns
- **User Management**: Customer support, account management
- **Content Management**: Pages, banners, promotions

### Technical Features
- Role-based access control
- Real-time notifications
- Bulk operations
- Export/import functionality
- Audit logging`,
    state: TaskStateValue.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1920, // 32 hours
    assignedAgent: frontendAgent.id,
    dependencies: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.PRODUCT_CATALOG],
    dependents: [ECOM_TASK_IDS.E2E_TESTS],
    createdAt: new Date("2024-01-18T10:00:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    position: { x: 1600, y: 400 },
  },

  // Phase 4: Testing & Deployment
  {
    id: ECOM_TASK_IDS.TESTING_SETUP,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Testing Framework Setup",
    description: `## Testing Framework Setup

### Testing Stack
- **Unit Tests**: Jest + Testing Library
- **Integration Tests**: Supertest for API
- **E2E Tests**: Cypress for user flows
- **Performance Tests**: Lighthouse CI
- **Accessibility Tests**: axe-core integration

### CI/CD Integration
- Automated test runs on PR
- Coverage reporting
- Visual regression testing
- Cross-browser testing
- Mobile device testing`,
    state: TaskStateValue.PENDING,
    complexity: 4,
    depth: 0,
    estimate: 360, // 6 hours
    assignedAgent: qaAgent.id,
    dependencies: [ECOM_TASK_IDS.FRONTEND_SETUP],
    dependents: [ECOM_TASK_IDS.E2E_TESTS],
    createdAt: new Date("2024-01-18T11:00:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    position: { x: 1300, y: 500 },
  },

  {
    id: ECOM_TASK_IDS.E2E_TESTS,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "End-to-End Testing Suite",
    description: `## End-to-End Testing Suite

### Critical User Flows
- **User Registration & Login**: Account creation, authentication
- **Product Discovery**: Search, filtering, navigation
- **Purchase Flow**: Add to cart, checkout, payment
- **Account Management**: Profile updates, order history
- **Admin Operations**: Product management, order processing

### Test Coverage
- Happy path scenarios
- Error handling flows
- Edge cases and boundary conditions
- Performance benchmarks
- Accessibility compliance`,
    state: TaskStateValue.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: qaAgent.id,
    dependencies: [
      ECOM_TASK_IDS.TESTING_SETUP,
      ECOM_TASK_IDS.PAYMENT_INTEGRATION,
      ECOM_TASK_IDS.ADMIN_PANEL,
    ],
    dependents: [ECOM_TASK_IDS.DEPLOYMENT],
    createdAt: new Date("2024-01-18T12:00:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    position: { x: 1900, y: 500 },
  },

  {
    id: ECOM_TASK_IDS.DEPLOYMENT,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Production Deployment & CI/CD",
    description: `## Production Deployment & CI/CD

### Infrastructure
- **Containerization**: Docker for all services
- **Orchestration**: Kubernetes cluster setup
- **Load Balancing**: NGINX with SSL termination
- **Database**: PostgreSQL with replication
- **Caching**: Redis cluster for sessions and cache

### CI/CD Pipeline
- Automated builds and tests
- Blue-green deployment strategy
- Database migration automation
- Monitoring and alerting setup
- Backup and disaster recovery`,
    state: TaskStateValue.PENDING,
    complexity: 6,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: devopsAgent.id,
    dependencies: [ECOM_TASK_IDS.E2E_TESTS],
    dependents: [],
    createdAt: new Date("2024-01-19T09:00:00Z"),
    updatedAt: new Date("2024-01-20T14:30:00Z"),
    position: { x: 2200, y: 600 },
  },
];
