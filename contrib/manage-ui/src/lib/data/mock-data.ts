import type { Project, Task } from '../types';

// Helper function to generate realistic dates
const getRandomDate = (daysAgo: number) => {
  const date = new Date();
  date.setDate(date.getDate() - Math.floor(Math.random() * daysAgo));
  return date.toISOString();
};

// Helper function to generate UUIDs (v4)
const generateUUID = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
};

// Pre-generated UUIDs for consistent references
const PROJECT_IDS = {
  frontend: 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d',
  backend: 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e',
  orchestration: 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f'
};

const AGENT_IDS = {
  frontend1: 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a',
  frontend2: 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b',
  design1: 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c',
  backend1: 'a7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d',
  backend2: 'b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e',
  ai1: 'c9d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f',
  ai2: 'd0e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a'
};

// Mock Projects with proper UUIDs
export const mockProjects: Project[] = [
  {
    id: PROJECT_IDS.frontend,
    title: 'Multi-Agent System Frontend',
    description: 'Entwicklung einer modernen Admin-Oberfläche für das Multi-Agent System mit NextJS, TypeScript und Tailwind CSS. Ziel ist es, eine intuitive Benutzeroberfläche für die Verwaltung von Projekten und Tasks zu schaffen.',
    created_at: getRandomDate(30),
    updated_at: getRandomDate(2),
    total_tasks: 24,
    completed_tasks: 8,
    progress: 33.33
  },
  {
    id: PROJECT_IDS.backend,
    title: 'Backend API Enhancement',
    description: 'Verbesserung und Erweiterung der Backend-API für bessere Performance, Skalierbarität und neue Features. Implementierung von SSE, Chat-Integration und erweiterte Task-Management-Funktionen.',
    created_at: getRandomDate(45),
    updated_at: getRandomDate(1),
    total_tasks: 18,
    completed_tasks: 14,
    progress: 77.78
  },
  {
    id: PROJECT_IDS.orchestration,
    title: 'AI Agent Orchestration Platform',
    description: 'Entwicklung einer Plattform zur Orchestrierung und Verwaltung von KI-Agenten. Umfasst Agent-Lifecycle-Management, Task-Zuweisung, Performance-Monitoring und Inter-Agent-Kommunikation.',
    created_at: getRandomDate(60),
    updated_at: getRandomDate(3),
    total_tasks: 31,
    completed_tasks: 12,
    progress: 38.71
  }
];

// Pre-generated Task UUIDs for consistent references
const TASK_IDS = {
  // Project 1: Frontend
  p1_setup: 'e1f2a3b4-c5d6-4e7f-8a9b-0c1d2e3f4a5b',
  p1_setup_nextjs: 'f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c',
  p1_setup_tailwind: 'a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d',
  p1_setup_structure: 'b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e',
  
  p1_design: 'c5d6e7f8-a9b0-4c1d-2e3f-4a5b6c7d8e9f',
  p1_design_shadcn: 'd6e7f8a9-b0c1-4d2e-3f4a-5b6c7d8e9f0a',
  p1_design_layout: 'e7f8a9b0-c1d2-4e3f-4a5b-6c7d8e9f0a1b',
  p1_design_theme: 'f8a9b0c1-d2e3-4f4a-5b6c-7d8e9f0a1b2c',
  
  p1_state: 'a9b0c1d2-e3f4-4a5b-6c7d-8e9f0a1b2c3d',
  p1_state_stores: 'b0c1d2e3-f4a5-4b6c-7d8e-9f0a1b2c3d4e',
  p1_state_api: 'c1d2e3f4-a5b6-4c7d-8e9f-0a1b2c3d4e5f',
  
  p1_kanban: 'd2e3f4a5-b6c7-4d8e-9f0a-1b2c3d4e5f6a',
  p1_kanban_layout: 'e3f4a5b6-c7d8-4e9f-0a1b-2c3d4e5f6a7b',
  p1_kanban_dragdrop: 'f4a5b6c7-d8e9-4f0a-1b2c-3d4e5f6a7b8c',
  
  p1_chat: 'a5b6c7d8-e9f0-4a1b-2c3d-4e5f6a7b8c9d',
  p1_testing: 'b6c7d8e9-f0a1-4b2c-3d4e-5f6a7b8c9d0e',
  
  // Project 2: Backend
  p2_performance: 'c7d8e9f0-a1b2-4c3d-4e5f-6a7b8c9d0e1f',
  p2_perf_db: 'd8e9f0a1-b2c3-4d4e-5f6a-7b8c9d0e1f2a',
  p2_perf_cache: 'e9f0a1b2-c3d4-4e5f-6a7b-8c9d0e1f2a3b',
  
  p2_sse: 'f0a1b2c3-d4e5-4f6a-7b8c-9d0e1f2a3b4c',
  p2_chat_api: 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d',
  p2_migration: 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e',
  p2_docs: 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f',
  p2_security: 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a',
  
  // Project 3: AI Orchestration
  p3_lifecycle: 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b',
  p3_lifecycle_registry: 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c',
  p3_lifecycle_health: 'a7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d',
  p3_lifecycle_scaling: 'b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e',
  
  p3_assignment: 'c9d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f',
  p3_assignment_matching: 'd0e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a',
  p3_assignment_queue: 'e1f2a3b4-c5d6-4e7f-8a9b-0c1d2e3f4a5b',
  
  p3_monitoring: 'f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c',
  p3_monitoring_metrics: 'a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d',
  p3_monitoring_dashboard: 'b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e',
  
  p3_communication: 'c5d6e7f8-a9b0-4c1d-2e3f-4a5b6c7d8e9f',
  p3_security: 'd6e7f8a9-b0c1-4d2e-3f4a-5b6c7d8e9f0a',
  p3_loadbalancing: 'e7f8a9b0-c1d2-4e3f-4a5b-6c7d8e9f0a1b',
  p3_documentation: 'f8a9b0c1-d2e3-4f4a-5b6c-7d8e9f0a1b2c'
};

// Mock Tasks with proper UUIDs and dependencies
export const mockTasks: Task[] = [
  // PROJECT 1: Multi-Agent System Frontend
  // Root Tasks
  {
    id: TASK_IDS.p1_setup,
    project_id: PROJECT_IDS.frontend,
    parent_id: null,
    title: 'Projekt Setup und Architektur',
    description: 'Initialisierung des NextJS Projekts, Setup der Entwicklungsumgebung und Definition der Grundarchitektur.',
    state: 'completed',
    complexity: 5,
    depth: 0,
    estimate: 480,
    assigned_agent: 'agent-frontend-1',
    dependencies: [],
    dependents: ['task-1-2', 'task-1-3'],
    created_at: getRandomDate(25),
    updated_at: getRandomDate(20),
    completed_at: getRandomDate(20)
  },
  {
    id: 'task-1-2',
    project_id: 'proj-1',
    parent_id: null,
    title: 'UI/UX Design System',
    description: 'Entwicklung eines konsistenten Design Systems mit Tailwind CSS und shadcn/ui Komponenten.',
    state: 'completed',
    complexity: 7,
    depth: 0,
    estimate: 720,
    assigned_agent: 'agent-design-1',
    dependencies: ['task-1-1'],
    dependents: ['task-1-4'],
    created_at: getRandomDate(22),
    updated_at: getRandomDate(15),
    completed_at: getRandomDate(15)
  },
  {
    id: 'task-1-3',
    project_id: 'proj-1',
    parent_id: null,
    title: 'State Management Implementation',
    description: 'Implementierung von Zustand für globales State Management und API-Integration.',
    state: 'in-progress',
    complexity: 8,
    depth: 0,
    estimate: 600,
    assigned_agent: 'agent-frontend-2',
    dependencies: ['task-1-1'],
    dependents: ['task-1-5'],
    created_at: getRandomDate(20),
    updated_at: getRandomDate(2)
  },
  {
    id: 'task-1-4',
    project_id: 'proj-1',
    parent_id: null,
    title: 'Kanban Task Management',
    description: 'Entwicklung der hierarchischen Kanban-Ansicht für Task-Management mit Drag & Drop.',
    state: 'in-progress',
    complexity: 9,
    depth: 0,
    estimate: 960,
    assigned_agent: 'agent-frontend-1',
    dependencies: ['task-1-2'],
    dependents: [],
    created_at: getRandomDate(18),
    updated_at: getRandomDate(1)
  },
  {
    id: 'task-1-5',
    project_id: 'proj-1',
    parent_id: null,
    title: 'Chat Integration',
    description: 'Integration der LLM-Chat-Funktionalität für Task- und Projekt-Assistenz.',
    state: 'pending',
    complexity: 8,
    depth: 0,
    estimate: 720,
    dependencies: ['task-1-3'],
    dependents: [],
    created_at: getRandomDate(15),
    updated_at: getRandomDate(15)
  },

  // Subtasks for Task 1-1 (Projekt Setup)
  {
    id: 'task-1-1-1',
    project_id: 'proj-1',
    parent_id: 'task-1-1',
    title: 'NextJS Projekt initialisieren',
    description: 'Setup des NextJS 15 Projekts mit TypeScript und App Router.',
    state: 'completed',
    complexity: 3,
    depth: 1,
    estimate: 120,
    assigned_agent: 'agent-frontend-1',
    dependencies: [],
    dependents: ['task-1-1-2'],
    created_at: getRandomDate(25),
    updated_at: getRandomDate(24),
    completed_at: getRandomDate(24)
  },
  {
    id: 'task-1-1-2',
    project_id: 'proj-1',
    parent_id: 'task-1-1',
    title: 'Tailwind CSS v4 konfigurieren',
    description: 'Installation und Konfiguration von Tailwind CSS v4 mit CSS-Variablen.',
    state: 'completed',
    complexity: 4,
    depth: 1,
    estimate: 180,
    assigned_agent: 'agent-frontend-1',
    dependencies: ['task-1-1-1'],
    dependents: ['task-1-1-3'],
    created_at: getRandomDate(24),
    updated_at: getRandomDate(23),
    completed_at: getRandomDate(23)
  },
  {
    id: 'task-1-1-3',
    project_id: 'proj-1',
    parent_id: 'task-1-1',
    title: 'Projektstruktur definieren',
    description: 'Erstellung der Ordnerstruktur und Definition der Architektur-Patterns.',
    state: 'completed',
    complexity: 3,
    depth: 1,
    estimate: 180,
    assigned_agent: 'agent-frontend-1',
    dependencies: ['task-1-1-2'],
    dependents: [],
    created_at: getRandomDate(23),
    updated_at: getRandomDate(22),
    completed_at: getRandomDate(22)
  },

  // Subtasks for Task 1-2 (UI/UX Design)
  {
    id: 'task-1-2-1',
    project_id: 'proj-1',
    parent_id: 'task-1-2',
    title: 'shadcn/ui Setup',
    description: 'Installation und Konfiguration der shadcn/ui Komponentenbibliothek.',
    state: 'completed',
    complexity: 4,
    depth: 1,
    estimate: 240,
    assigned_agent: 'agent-design-1',
    dependencies: [],
    dependents: ['task-1-2-2'],
    created_at: getRandomDate(22),
    updated_at: getRandomDate(18),
    completed_at: getRandomDate(18)
  },
  {
    id: 'task-1-2-2',
    project_id: 'proj-1',
    parent_id: 'task-1-2',
    title: 'Layout Komponenten',
    description: 'Entwicklung der Basis-Layout-Komponenten (Sidebar, Navigation, etc.).',
    state: 'completed',
    complexity: 6,
    depth: 1,
    estimate: 360,
    assigned_agent: 'agent-design-1',
    dependencies: ['task-1-2-1'],
    dependents: ['task-1-2-3'],
    created_at: getRandomDate(20),
    updated_at: getRandomDate(16),
    completed_at: getRandomDate(16)
  },
  {
    id: 'task-1-2-3',
    project_id: 'proj-1',
    parent_id: 'task-1-2',
    title: 'Theme System',
    description: 'Implementierung des Dark/Light Theme Systems mit CSS-Variablen.',
    state: 'completed',
    complexity: 5,
    depth: 1,
    estimate: 120,
    assigned_agent: 'agent-design-1',
    dependencies: ['task-1-2-2'],
    dependents: [],
    created_at: getRandomDate(18),
    updated_at: getRandomDate(15),
    completed_at: getRandomDate(15)
  },

  // Subtasks for Task 1-3 (State Management)
  {
    id: 'task-1-3-1',
    project_id: 'proj-1',
    parent_id: 'task-1-3',
    title: 'Zustand Stores definieren',
    description: 'Definition der Zustand Stores für Project, UI und Chat Management.',
    state: 'completed',
    complexity: 6,
    depth: 1,
    estimate: 300,
    assigned_agent: 'agent-frontend-2',
    dependencies: [],
    dependents: ['task-1-3-2'],
    created_at: getRandomDate(20),
    updated_at: getRandomDate(10),
    completed_at: getRandomDate(10)
  },
  {
    id: 'task-1-3-2',
    project_id: 'proj-1',
    parent_id: 'task-1-3',
    title: 'API Client implementieren',
    description: 'Entwicklung des API Clients mit Error Handling und TypeScript Types.',
    state: 'in-progress',
    complexity: 7,
    depth: 1,
    estimate: 300,
    assigned_agent: 'agent-frontend-2',
    dependencies: ['task-1-3-1'],
    dependents: [],
    created_at: getRandomDate(15),
    updated_at: getRandomDate(2)
  },

  // Subtasks for Task 1-4 (Kanban)
  {
    id: 'task-1-4-1',
    project_id: 'proj-1',
    parent_id: 'task-1-4',
    title: 'Kanban Board Layout',
    description: 'Grundlegendes Layout für die Kanban-Ansicht mit Spalten und Cards.',
    state: 'in-progress',
    complexity: 6,
    depth: 1,
    estimate: 480,
    assigned_agent: 'agent-frontend-1',
    dependencies: [],
    dependents: ['task-1-4-2'],
    created_at: getRandomDate(18),
    updated_at: getRandomDate(1)
  },
  {
    id: 'task-1-4-2',
    project_id: 'proj-1',
    parent_id: 'task-1-4',
    title: 'Drag & Drop Funktionalität',
    description: 'Implementierung von Drag & Drop für Task-State-Management.',
    state: 'pending',
    complexity: 8,
    depth: 1,
    estimate: 480,
    dependencies: ['task-1-4-1'],
    dependents: [],
    created_at: getRandomDate(15),
    updated_at: getRandomDate(15)
  },

  // PROJECT 2: Backend API Enhancement
  // Root Tasks
  {
    id: 'task-2-1',
    project_id: 'proj-2',
    parent_id: null,
    title: 'API Performance Optimierung',
    description: 'Optimierung der bestehenden API-Endpunkte für bessere Performance und Skalierbarkeit.',
    state: 'completed',
    complexity: 8,
    depth: 0,
    estimate: 960,
    assigned_agent: 'agent-backend-1',
    dependencies: [],
    dependents: ['task-2-2'],
    created_at: getRandomDate(40),
    updated_at: getRandomDate(25),
    completed_at: getRandomDate(25)
  },
  {
    id: 'task-2-2',
    project_id: 'proj-2',
    parent_id: null,
    title: 'SSE Implementation',
    description: 'Implementierung von Server-Sent Events für Real-time Updates.',
    state: 'completed',
    complexity: 7,
    depth: 0,
    estimate: 720,
    assigned_agent: 'agent-backend-2',
    dependencies: ['task-2-1'],
    dependents: ['task-2-4'],
    created_at: getRandomDate(35),
    updated_at: getRandomDate(20),
    completed_at: getRandomDate(20)
  },
  {
    id: 'task-2-3',
    project_id: 'proj-2',
    parent_id: null,
    title: 'Chat API Integration',
    description: 'Entwicklung der Chat-API für LLM-Integration und Session-Management.',
    state: 'completed',
    complexity: 9,
    depth: 0,
    estimate: 1200,
    assigned_agent: 'agent-backend-1',
    dependencies: [],
    dependents: ['task-2-5'],
    created_at: getRandomDate(30),
    updated_at: getRandomDate(15),
    completed_at: getRandomDate(15)
  },
  {
    id: 'task-2-4',
    project_id: 'proj-2',
    parent_id: null,
    title: 'Database Schema Migration',
    description: 'Migration der Datenbankschemas für neue Features und Performance-Verbesserungen.',
    state: 'in-progress',
    complexity: 6,
    depth: 0,
    estimate: 480,
    assigned_agent: 'agent-backend-2',
    dependencies: ['task-2-2'],
    dependents: [],
    created_at: getRandomDate(25),
    updated_at: getRandomDate(3)
  },
  {
    id: 'task-2-5',
    project_id: 'proj-2',
    parent_id: null,
    title: 'API Documentation Update',
    description: 'Aktualisierung der API-Dokumentation mit OpenAPI/Swagger Specs.',
    state: 'pending',
    complexity: 4,
    depth: 0,
    estimate: 360,
    assigned_agent: 'agent-backend-1',
    dependencies: ['task-2-3'],
    dependents: [],
    created_at: getRandomDate(20),
    updated_at: getRandomDate(20)
  },

  // Subtasks for Task 2-1 (Performance)
  {
    id: 'task-2-1-1',
    project_id: 'proj-2',
    parent_id: 'task-2-1',
    title: 'Database Query Optimierung',
    description: 'Optimierung der Datenbankabfragen und Indizierung.',
    state: 'completed',
    complexity: 7,
    depth: 1,
    estimate: 480,
    assigned_agent: 'agent-backend-1',
    dependencies: [],
    dependents: ['task-2-1-2'],
    created_at: getRandomDate(40),
    updated_at: getRandomDate(30),
    completed_at: getRandomDate(30)
  },
  {
    id: 'task-2-1-2',
    project_id: 'proj-2',
    parent_id: 'task-2-1',
    title: 'Caching Strategy',
    description: 'Implementierung einer effizienten Caching-Strategie.',
    state: 'completed',
    complexity: 6,
    depth: 1,
    estimate: 360,
    assigned_agent: 'agent-backend-1',
    dependencies: ['task-2-1-1'],
    dependents: [],
    created_at: getRandomDate(35),
    updated_at: getRandomDate(25),
    completed_at: getRandomDate(25)
  },

  // PROJECT 3: AI Agent Orchestration Platform
  // Root Tasks
  {
    id: 'task-3-1',
    project_id: 'proj-3',
    parent_id: null,
    title: 'Agent Lifecycle Management',
    description: 'Entwicklung des Systems für Agent-Erstellung, -Verwaltung und -Überwachung.',
    state: 'completed',
    complexity: 9,
    depth: 0,
    estimate: 1440,
    assigned_agent: 'agent-ai-1',
    dependencies: [],
    dependents: ['task-3-2', 'task-3-3'],
    created_at: getRandomDate(55),
    updated_at: getRandomDate(35),
    completed_at: getRandomDate(35)
  },
  {
    id: 'task-3-2',
    project_id: 'proj-3',
    parent_id: null,
    title: 'Task Assignment Engine',
    description: 'Intelligente Task-Zuweisung basierend auf Agent-Fähigkeiten und Verfügbarkeit.',
    state: 'completed',
    complexity: 8,
    depth: 0,
    estimate: 1200,
    assigned_agent: 'agent-ai-2',
    dependencies: ['task-3-1'],
    dependents: ['task-3-4'],
    created_at: getRandomDate(50),
    updated_at: getRandomDate(30),
    completed_at: getRandomDate(30)
  },
  {
    id: 'task-3-3',
    project_id: 'proj-3',
    parent_id: null,
    title: 'Performance Monitoring',
    description: 'Real-time Monitoring und Analytics für Agent-Performance.',
    state: 'in-progress',
    complexity: 7,
    depth: 0,
    estimate: 960,
    assigned_agent: 'agent-ai-1',
    dependencies: ['task-3-1'],
    dependents: [],
    created_at: getRandomDate(45),
    updated_at: getRandomDate(5)
  },
  {
    id: 'task-3-4',
    project_id: 'proj-3',
    parent_id: null,
    title: 'Inter-Agent Communication',
    description: 'Protokoll und Infrastruktur für Kommunikation zwischen Agenten.',
    state: 'in-progress',
    complexity: 9,
    depth: 0,
    estimate: 1440,
    assigned_agent: 'agent-ai-2',
    dependencies: ['task-3-2'],
    dependents: ['task-3-5'],
    created_at: getRandomDate(40),
    updated_at: getRandomDate(2)
  },
  {
    id: 'task-3-5',
    project_id: 'proj-3',
    parent_id: null,
    title: 'Security Framework',
    description: 'Sicherheitsframework für Agent-Authentifizierung und -Autorisierung.',
    state: 'pending',
    complexity: 8,
    depth: 0,
    estimate: 960,
    dependencies: ['task-3-4'],
    dependents: [],
    created_at: getRandomDate(35),
    updated_at: getRandomDate(35)
  },
  {
    id: 'task-3-6',
    project_id: 'proj-3',
    parent_id: null,
    title: 'Load Balancing',
    description: 'Implementierung von Load Balancing für Agent-Workloads.',
    state: 'blocked',
    complexity: 7,
    depth: 0,
    estimate: 720,
    dependencies: [],
    dependents: [],
    created_at: getRandomDate(30),
    updated_at: getRandomDate(10)
  },

  // Subtasks for Task 3-1 (Agent Lifecycle)
  {
    id: 'task-3-1-1',
    project_id: 'proj-3',
    parent_id: 'task-3-1',
    title: 'Agent Registry',
    description: 'Zentrales Registry für alle verfügbaren Agenten und ihre Capabilities.',
    state: 'completed',
    complexity: 6,
    depth: 1,
    estimate: 480,
    assigned_agent: 'agent-ai-1',
    dependencies: [],
    dependents: ['task-3-1-2'],
    created_at: getRandomDate(55),
    updated_at: getRandomDate(45),
    completed_at: getRandomDate(45)
  },
  {
    id: 'task-3-1-2',
    project_id: 'proj-3',
    parent_id: 'task-3-1',
    title: 'Agent Health Monitoring',
    description: 'System zur Überwachung des Gesundheitszustands von Agenten.',
    state: 'completed',
    complexity: 7,
    depth: 1,
    estimate: 600,
    assigned_agent: 'agent-ai-1',
    dependencies: ['task-3-1-1'],
    dependents: ['task-3-1-3'],
    created_at: getRandomDate(50),
    updated_at: getRandomDate(40),
    completed_at: getRandomDate(40)
  },
  {
    id: 'task-3-1-3',
    project_id: 'proj-3',
    parent_id: 'task-3-1',
    title: 'Auto-scaling Logic',
    description: 'Automatische Skalierung von Agenten basierend auf Workload.',
    state: 'completed',
    complexity: 8,
    depth: 1,
    estimate: 360,
    assigned_agent: 'agent-ai-1',
    dependencies: ['task-3-1-2'],
    dependents: [],
    created_at: getRandomDate(45),
    updated_at: getRandomDate(35),
    completed_at: getRandomDate(35)
  },

  // Subtasks for Task 3-2 (Task Assignment)
  {
    id: 'task-3-2-1',
    project_id: 'proj-3',
    parent_id: 'task-3-2',
    title: 'Capability Matching Algorithm',
    description: 'Algorithmus zum Matching von Tasks mit Agent-Capabilities.',
    state: 'completed',
    complexity: 8,
    depth: 1,
    estimate: 720,
    assigned_agent: 'agent-ai-2',
    dependencies: [],
    dependents: ['task-3-2-2'],
    created_at: getRandomDate(50),
    updated_at: getRandomDate(35),
    completed_at: getRandomDate(35)
  },
  {
    id: 'task-3-2-2',
    project_id: 'proj-3',
    parent_id: 'task-3-2',
    title: 'Priority Queue System',
    description: 'Prioritäts-basiertes Queue-System für Task-Assignment.',
    state: 'completed',
    complexity: 6,
    depth: 1,
    estimate: 480,
    assigned_agent: 'agent-ai-2',
    dependencies: ['task-3-2-1'],
    dependents: [],
    created_at: getRandomDate(45),
    updated_at: getRandomDate(30),
    completed_at: getRandomDate(30)
  },

  // Subtasks for Task 3-3 (Performance Monitoring)
  {
    id: 'task-3-3-1',
    project_id: 'proj-3',
    parent_id: 'task-3-3',
    title: 'Metrics Collection',
    description: 'System zur Sammlung von Performance-Metriken von Agenten.',
    state: 'completed',
    complexity: 6,
    depth: 1,
    estimate: 480,
    assigned_agent: 'agent-ai-1',
    dependencies: [],
    dependents: ['task-3-3-2'],
    created_at: getRandomDate(45),
    updated_at: getRandomDate(25),
    completed_at: getRandomDate(25)
  },
  {
    id: 'task-3-3-2',
    project_id: 'proj-3',
    parent_id: 'task-3-3',
    title: 'Dashboard Implementation',
    description: 'Real-time Dashboard für Performance-Monitoring.',
    state: 'in-progress',
    complexity: 7,
    depth: 1,
    estimate: 480,
    assigned_agent: 'agent-ai-1',
    dependencies: ['task-3-3-1'],
    dependents: [],
    created_at: getRandomDate(30),
    updated_at: getRandomDate(5)
  },

  // Additional tasks for variety
  {
    id: 'task-1-6',
    project_id: 'proj-1',
    parent_id: null,
    title: 'Testing & Quality Assurance',
    description: 'Umfassende Tests und Qualitätssicherung für das Frontend.',
    state: 'pending',
    complexity: 6,
    depth: 0,
    estimate: 720,
    dependencies: ['task-1-4', 'task-1-5'],
    dependents: [],
    created_at: getRandomDate(10),
    updated_at: getRandomDate(10)
  },
  {
    id: 'task-2-6',
    project_id: 'proj-2',
    parent_id: null,
    title: 'Security Audit',
    description: 'Sicherheitsaudit der API-Endpunkte und Authentifizierung.',
    state: 'cancelled',
    complexity: 5,
    depth: 0,
    estimate: 480,
    dependencies: [],
    dependents: [],
    created_at: getRandomDate(15),
    updated_at: getRandomDate(8)
  },
  {
    id: 'task-3-7',
    project_id: 'proj-3',
    parent_id: null,
    title: 'Documentation & Training',
    description: 'Erstellung von Dokumentation und Trainingsmaterialien.',
    state: 'pending',
    complexity: 4,
    depth: 0,
    estimate: 600,
    dependencies: ['task-3-5'],
    dependents: [],
    created_at: getRandomDate(20),
    updated_at: getRandomDate(20)
  }
];

// Calculate and update project progress based on tasks
export const updateProjectProgress = () => {
  mockProjects.forEach(project => {
    const projectTasks = mockTasks.filter(task => task.project_id === project.id);
    const completedTasks = projectTasks.filter(task => task.state === 'completed');
    
    project.total_tasks = projectTasks.length;
    project.completed_tasks = completedTasks.length;
    project.progress = projectTasks.length > 0 ? (completedTasks.length / projectTasks.length) * 100 : 0;
  });
};

// Initialize project progress
updateProjectProgress();