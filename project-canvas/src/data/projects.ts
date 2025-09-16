/**
 * Projects Data
 * Go-Model konform: pkg/tools/project/shared/shared.go
 */

import { Project } from '../types';

// Feste Project-IDs für konsistente Referenzen (müssen mit master-dummy-data.ts übereinstimmen)
const PROJECT_IDS = {
  ECOMMERCE: "53fd7995-1606-4360-65ae-9b057525c12e",
  MOBILE: "980c9337-0630-41d0-3b95-b71caca7635b"
};

export const masterProjects: Project[] = [
  // Projekt 1: E-Commerce Platform Redesign
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

// Export individual projects for easy access
export const ecommerceProject = masterProjects[0];
export const mobileProject = masterProjects[1];