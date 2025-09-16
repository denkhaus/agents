/**
 * E-Commerce Project Tasks Part 2 (15 Tasks)
 * Go-Model konform: pkg/tools/project/shared/shared.go
 * Max 500 Zeilen pro Datei!
 */

import type { Task } from '../types';
import { generateUUID } from '../utils/uuid';
import { ecommerceProject } from './projects';
import { frontendAgent, backendAgent, qaAgent, devopsAgent } from './agents';

// Import the actual values for runtime usage
import { TaskState as TaskStateValue } from '../types/task.types';

// Define ECOM_TASK_IDS locally since it's not exported from projects.ts
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

// Additional Task IDs
export const ECOM_TASK_IDS_PART2 = {
  // Phase 2: Architecture (continued)
  API_DOCUMENTATION: generateUUID(),
  CI_CD_SETUP: generateUUID(),
  
  // Phase 3: Core Development (continued)
  NOTIFICATIONS: generateUUID(),
  WISHLIST: generateUUID(),
  REVIEWS_RATINGS: generateUUID(),
  ORDER_MANAGEMENT: generateUUID(),
  INVENTORY_MANAGEMENT: generateUUID(),
  
  // Phase 4: Testing & Deployment (continued)
  PERFORMANCE_TESTING: generateUUID(),
  SECURITY_AUDIT: generateUUID(),
  MONITORING_SETUP: generateUUID()
};

export const ecommerceTasksPart2: Task[] = [
  // Phase 2: Architecture (continued)
  {
    id: ECOM_TASK_IDS_PART2.API_DOCUMENTATION,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "API Documentation & SDKs",
    description: `## API Documentation & SDK Development

### Documentation Requirements
- **OpenAPI Specification**: Complete Swagger/OpenAPI 3.0 spec
- **Interactive Documentation**: Swagger UI integration
- **Code Examples**: JavaScript, Python, and cURL examples
- **Authentication Guide**: OAuth2 and API key usage

### SDK Development
- **JavaScript SDK**: Node.js and browser compatible
- **Python SDK**: For backend integrations
- **Mobile SDKs**: iOS and Android (React Native compatible)

### Deliverables
- Published API documentation
- SDK packages on NPM and PyPI
- Integration guides`,
    state: TaskStateValue.COMPLETED,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: backendAgent.id,
    dependencies: [ECOM_TASK_IDS.API_DESIGN, ECOM_TASK_IDS.DATABASE_SCHEMA],
    dependents: [ECOM_TASK_IDS_PART2.CI_CD_SETUP],
    createdAt: new Date('2024-01-16T10:00:00Z'),
    updatedAt: new Date('2024-01-18T17:30:00Z'),
    completedAt: new Date('2024-01-18T17:30:00Z'),
    position: { x: 700, y: 400 }
  },

  {
    id: ECOM_TASK_IDS_PART2.CI_CD_SETUP,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "CI/CD Pipeline Setup",
    description: `## Continuous Integration & Deployment Pipeline

### Pipeline Stages
- **Code Quality**: ESLint, Prettier, and type checking
- **Automated Testing**: Unit, integration, and E2E tests
- **Security Scanning**: Dependency vulnerability checks
- **Build Process**: Webpack/Rollup optimization
- **Deployment**: Staging and production environments

### Tools & Services
- **GitHub Actions**: CI/CD workflows
- **Docker**: Containerization for consistent environments
- **Kubernetes**: Orchestration for scalable deployment
- **Helm Charts**: Kubernetes deployment templates

### Deliverables
- Fully automated CI/CD pipeline
- Deployment documentation
- Monitoring integration`,
    state: TaskStateValue.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: devopsAgent.id,
    dependencies: [ECOM_TASK_IDS.BACKEND_SETUP, ECOM_TASK_IDS.FRONTEND_SETUP, ECOM_TASK_IDS_PART2.API_DOCUMENTATION],
    dependents: [],
    createdAt: new Date('2024-01-16T10:15:00Z'),
    updatedAt: new Date('2024-01-18T09:45:00Z'),
    position: { x: 1000, y: 400 }
  },

  // Phase 3: Core Development (continued)
  {
    id: ECOM_TASK_IDS_PART2.NOTIFICATIONS,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Email & Push Notifications",
    description: `## Notification System Implementation

### Notification Types
- **Order Updates**: Confirmation, shipping, delivery
- **Account Security**: Login alerts, password changes
- **Marketing**: Promotions, abandoned cart reminders
- **System Alerts**: Maintenance, outages

### Technical Implementation
- **Email Templates**: Responsive HTML email templates
- **Push Service**: Firebase Cloud Messaging integration
- **SMS Gateway**: Twilio integration for critical alerts
- **Preferences**: User notification settings management

### Deliverables
- Notification service API
- Email template system
- Mobile push integration
- User preference dashboard`,
    state: TaskStateValue.COMPLETED,
    complexity: 6,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: backendAgent.id,
    dependencies: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.ADMIN_PANEL],
    dependents: [ECOM_TASK_IDS_PART2.ORDER_MANAGEMENT],
    createdAt: new Date('2024-01-17T09:30:00Z'),
    updatedAt: new Date('2024-01-19T16:20:00Z'),
    completedAt: new Date('2024-01-19T16:20:00Z'),
    position: { x: 400, y: 700 }
  },

  {
    id: ECOM_TASK_IDS_PART2.WISHLIST,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Wishlist & Save for Later",
    description: `## Wishlist Feature Implementation

### Core Functionality
- **Add/Remove Items**: Simple wishlist management
- **Sharing**: Social sharing of wishlists
- **Price Tracking**: Price drop notifications
- **Reminder System**: Periodic wishlist check-ins

### Technical Requirements
- **Persistence**: User-specific wishlist storage
- **Sync**: Cross-device wishlist synchronization
- **Performance**: Efficient loading of large wishlists
- **Privacy**: Public/private wishlist settings

### Deliverables
- Wishlist API endpoints
- Frontend wishlist components
- Sharing functionality
- Price tracking integration`,
    state: TaskStateValue.COMPLETED,
    complexity: 5,
    depth: 0,
    estimate: 960, // 16 hours
    assignedAgent: frontendAgent.id,
    dependencies: [ECOM_TASK_IDS.USER_AUTH, ECOM_TASK_IDS.PRODUCT_CATALOG],
    dependents: [],
    createdAt: new Date('2024-01-17T09:45:00Z'),
    updatedAt: new Date('2024-01-18T15:30:00Z'),
    completedAt: new Date('2024-01-18T15:30:00Z'),
    position: { x: 700, y: 700 }
  },

  {
    id: ECOM_TASK_IDS_PART2.REVIEWS_RATINGS,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Product Reviews & Ratings",
    description: `## Product Review System

### Features
- **Star Ratings**: 1-5 star rating system
- **Written Reviews**: Text reviews with media support
- **Helpfulness Voting**: User feedback on reviews
- **Moderation**: Admin review approval system
- **Sorting/Filtering**: Various review sorting options

### Implementation Details
- **Schema Design**: Efficient review storage and retrieval
- **Spam Protection**: Automated and manual moderation
- **Rich Media**: Image/video review attachments
- **SEO Optimization**: Structured data for reviews

### Deliverables
- Review API endpoints
- Review display components
- Moderation dashboard
- Rich media handling`,
    state: TaskStateValue.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1440, // 24 hours
    assignedAgent: frontendAgent.id,
    dependencies: [ECOM_TASK_IDS.PRODUCT_CATALOG, ECOM_TASK_IDS.USER_AUTH],
    dependents: [],
    createdAt: new Date('2024-01-17T10:00:00Z'),
    updatedAt: new Date('2024-01-17T10:00:00Z'),
    position: { x: 1000, y: 700 }
  },

  {
    id: ECOM_TASK_IDS_PART2.ORDER_MANAGEMENT,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Order Management System",
    description: `## Order Processing & Management

### Core Features
- **Order Creation**: Cart to order conversion
- **Status Tracking**: Real-time order status updates
- **Payment Processing**: Integration with payment gateways
- **Fulfillment**: Shipping and delivery coordination
- **Returns**: Return request and processing system

### Technical Requirements
- **ACID Compliance**: Transactional order processing
- **Audit Trail**: Complete order history tracking
- **Notifications**: Automated status updates
- **Reports**: Sales and fulfillment analytics

### Deliverables
- Order processing API
- Admin order management dashboard
- Customer order tracking interface
- Integration with shipping providers`,
    state: TaskStateValue.PENDING,
    complexity: 9,
    depth: 0,
    estimate: 2160, // 36 hours
    assignedAgent: backendAgent.id,
    dependencies: [ECOM_TASK_IDS.SHOPPING_CART, ECOM_TASK_IDS.PAYMENT_INTEGRATION, ECOM_TASK_IDS_PART2.NOTIFICATIONS],
    dependents: [ECOM_TASK_IDS_PART2.INVENTORY_MANAGEMENT],
    createdAt: new Date('2024-01-17T10:15:00Z'),
    updatedAt: new Date('2024-01-17T10:15:00Z'),
    position: { x: 1300, y: 700 }
  },

  {
    id: ECOM_TASK_IDS_PART2.INVENTORY_MANAGEMENT,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Inventory & Stock Management",
    description: `## Real-time Inventory System

### Key Features
- **Stock Levels**: Real-time inventory tracking
- **Low Stock Alerts**: Automatic notifications
- **Supplier Integration**: Vendor management system
- **Batch Tracking**: Lot and expiration tracking
- **Multi-location**: Warehouse and store inventory

### Technical Implementation
- **Concurrency Control**: Prevent overselling
- **Caching**: Redis for fast inventory queries
- **Webhooks**: Supplier system integrations
- **Reporting**: Inventory analytics and forecasting

### Deliverables
- Inventory API endpoints
- Admin inventory dashboard
- Supplier integration system
- Reporting and analytics`,
    state: TaskStateValue.PENDING,
    complexity: 8,
    depth: 0,
    estimate: 1800, // 30 hours
    assignedAgent: backendAgent.id,
    dependencies: [ECOM_TASK_IDS_PART2.ORDER_MANAGEMENT],
    dependents: [],
    createdAt: new Date('2024-01-17T10:30:00Z'),
    updatedAt: new Date('2024-01-17T10:30:00Z'),
    position: { x: 1600, y: 700 }
  },

  // Phase 4: Testing & Deployment (continued)
  {
    id: ECOM_TASK_IDS_PART2.PERFORMANCE_TESTING,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Performance & Load Testing",
    description: `## System Performance Validation

### Testing Scenarios
- **Load Testing**: Simulate 10K concurrent users
- **Stress Testing**: Identify system breaking points
- **Soak Testing**: Long-term stability validation
- **Spike Testing**: Sudden traffic surge handling

### Tools & Metrics
- **k6/Locust**: Load testing frameworks
- **New Relic**: Performance monitoring
- **Response Times**: Target <200ms for APIs
- **Throughput**: 1000+ requests/second

### Deliverables
- Performance test reports
- Bottleneck identification
- Optimization recommendations
- Monitoring dashboard`,
    state: TaskStateValue.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: qaAgent.id,
    dependencies: [ECOM_TASK_IDS.DEPLOYMENT],
    dependents: [ECOM_TASK_IDS_PART2.SECURITY_AUDIT],
    createdAt: new Date('2024-01-18T09:00:00Z'),
    updatedAt: new Date('2024-01-18T09:00:00Z'),
    position: { x: 1000, y: 1000 }
  },

  {
    id: ECOM_TASK_IDS_PART2.SECURITY_AUDIT,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Security Audit & Compliance",
    description: `## Security Assessment & Hardening

### Audit Areas
- **OWASP Top 10**: Vulnerability assessment
- **Data Protection**: GDPR/CCPA compliance
- **Authentication**: OAuth2, JWT security
- **API Security**: Rate limiting, input validation
- **Infrastructure**: Network and container security

### Compliance Requirements
- **PCI DSS**: Payment card industry standards
- **SOC 2**: Security and availability trust principles
- **Encryption**: At-rest and in-transit data encryption

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
    dependencies: [ECOM_TASK_IDS.DEPLOYMENT, ECOM_TASK_IDS_PART2.PERFORMANCE_TESTING],
    dependents: [ECOM_TASK_IDS_PART2.MONITORING_SETUP],
    createdAt: new Date('2024-01-18T09:15:00Z'),
    updatedAt: new Date('2024-01-18T09:15:00Z'),
    position: { x: 1300, y: 1000 }
  },

  {
    id: ECOM_TASK_IDS_PART2.MONITORING_SETUP,
    projectId: ecommerceProject.id,
    parentId: undefined,
    title: "Monitoring & Alerting System",
    description: `## Production Monitoring Infrastructure

### Monitoring Layers
- **Application**: Error rates, response times, throughput
- **Infrastructure**: CPU, memory, disk, network
- **Business Metrics**: Conversion rates, sales tracking
- **User Experience**: Frontend performance, RUM

### Alerting System
- **Thresholds**: Configurable alert thresholds
- **Channels**: Email, Slack, PagerDuty integrations
- **Escalation**: Multi-level alert escalation
- **Dashboards**: Real-time operational dashboards

### Tools & Services
- **Prometheus**: Metrics collection and storage
- **Grafana**: Visualization and dashboarding
- **ELK Stack**: Log aggregation and analysis
- **Sentry**: Error tracking and reporting

### Deliverables
- Monitoring infrastructure
- Alerting rules and policies
- Operational dashboards
- Incident response procedures`,
    state: TaskStateValue.PENDING,
    complexity: 7,
    depth: 0,
    estimate: 1200, // 20 hours
    assignedAgent: devopsAgent.id,
    dependencies: [ECOM_TASK_IDS.DEPLOYMENT, ECOM_TASK_IDS_PART2.SECURITY_AUDIT],
    dependents: [],
    createdAt: new Date('2024-01-18T09:30:00Z'),
    updatedAt: new Date('2024-01-18T09:30:00Z'),
    position: { x: 1600, y: 1000 }
  }
];