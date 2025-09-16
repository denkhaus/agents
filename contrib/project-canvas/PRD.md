# Project Canvas - Product Requirements Document (PRD)

## 1. Project Overview

### 1.1 Vision
Eine React-basierte Webanwendung zur Visualisierung und Echtzeit-Bearbeitung von Projekten, Tasks und deren Abhängigkeiten in einem interaktiven Flow-Diagramm.

### 1.2 Core Features
- **Multi-Workspace Architecture**: Erweiterbare Sidebar für Projects, Agents, Settings und zukünftige Datentypen
- **Collapsible Sidebar**: Icon-only Modus für maximalen Canvas-Platz
- **Projekt-Visualisierung**: Hierarchische Darstellung von Projekten und Tasks als ReactFlow Nodes
- **Dependency Mapping**: Abhängigkeiten zwischen Tasks als Edges visualisiert
- **Real-time Updates**: Echtzeit-Synchronisation über Convex
- **Interactive Canvas**: Drag & Drop, Zoom, Pan Funktionalitäten
- **Auto-Layout**: Intelligente Positionierung zur Vermeidung von Überlappungen
- **Multi-Project Support**: Navigation zwischen verschiedenen Projekten
- **Dark Mode**: Theme-Switching Unterstützung
- **Structured State Management**: Zustand-basierte Store-Architektur für saubere Datentrennung

## 2. Technical Architecture

### 2.1 Tech Stack
- **Frontend**: React 18 + TypeScript
- **Build Tool**: Vite
- **UI Framework**: Tailwind CSS 3 + shadcn/ui
- **Flow Visualization**: ReactFlow
- **Real-time Data**: Convex
- **State Management**: Zustand + Convex Hooks
- **Routing**: React Router (falls Multi-Page benötigt)

### 2.2 Data Model (TypeScript)

```typescript
// Core Types (UUID-based)
type UUID = string; // UUID v4 format validation

enum TaskState {
  PENDING = "pending",
  IN_PROGRESS = "in-progress", 
  COMPLETED = "completed",
  BLOCKED = "blocked",
  CANCELLED = "cancelled"
}

interface Task {
  id: UUID;
  projectId: UUID;
  parentId?: UUID;
  title: string;
  description: string;
  state: TaskState;
  complexity: number; // 1-10
  depth: number;
  estimate?: number; // minutes
  assignedAgent?: UUID;
  dependencies: UUID[];
  dependents: UUID[];
  createdAt: Date;
  updatedAt: Date;
  completedAt?: Date;
  // UI-specific
  position?: { x: number; y: number };
}

interface Project {
  id: UUID;
  title: string;
  description: string;
  createdAt: Date;
  updatedAt: Date;
  totalTasks: number;
  completedTasks: number;
  progress: number; // 0-100
}

interface ProjectProgress {
  projectId: UUID;
  totalTasks: number;
  completedTasks: number;
  inProgressTasks: number;
  pendingTasks: number;
  blockedTasks: number;
  cancelledTasks: number;
  overallProgress: number;
  tasksByDepth: Record<number, number>;
}
```

### 2.3 Zustand Store Architecture

```typescript
// stores/projectStore.ts
interface ProjectStore {
  projects: Project[];
  currentProject: Project | null;
  loading: boolean;
  error: string | null;
  
  // Actions
  setProjects: (projects: Project[]) => void;
  setCurrentProject: (project: Project) => void;
  addProject: (project: Project) => void;
  updateProject: (id: UUID, updates: Partial<Project>) => void;
  deleteProject: (id: UUID) => void;
}

// stores/taskStore.ts
interface TaskStore {
  tasks: Task[];
  tasksByProject: Record<UUID, Task[]>;
  loading: boolean;
  error: string | null;
  
  // Actions
  setTasks: (tasks: Task[]) => void;
  addTask: (task: Task) => void;
  updateTask: (id: UUID, updates: Partial<Task>) => void;
  deleteTask: (id: UUID) => void;
  updateTaskPosition: (id: UUID, position: { x: number; y: number }) => void;
}

// stores/uiStore.ts
interface UIStore {
  sidebarCollapsed: boolean;
  currentWorkspace: 'projects' | 'agents' | 'settings';
  darkMode: boolean;
  selectedNodes: string[];
  
  // Actions
  toggleSidebar: () => void;
  setWorkspace: (workspace: string) => void;
  toggleDarkMode: () => void;
  setSelectedNodes: (nodes: string[]) => void;
}

// stores/agentStore.ts (für zukünftige Erweiterung)
interface AgentStore {
  agents: Agent[];
  loading: boolean;
  error: string | null;
  
  // Actions
  setAgents: (agents: Agent[]) => void;
  addAgent: (agent: Agent) => void;
  updateAgent: (id: UUID, updates: Partial<Agent>) => void;
}
```

### 2.4 Convex Schema

```typescript
// convex/schema.ts
import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  projects: defineTable({
    title: v.string(),
    description: v.string(),
    totalTasks: v.number(),
    completedTasks: v.number(),
    progress: v.number(),
  }),
  
  tasks: defineTable({
    projectId: v.id("projects"),
    parentId: v.optional(v.id("tasks")),
    title: v.string(),
    description: v.string(),
    state: v.union(
      v.literal("pending"),
      v.literal("in-progress"),
      v.literal("completed"),
      v.literal("blocked"),
      v.literal("cancelled")
    ),
    complexity: v.number(),
    depth: v.number(),
    estimate: v.optional(v.number()),
    assignedAgent: v.optional(v.string()),
    dependencies: v.array(v.id("tasks")),
    dependents: v.array(v.id("tasks")),
    positionX: v.optional(v.number()),
    positionY: v.optional(v.number()),
  })
  .index("by_project", ["projectId"])
  .index("by_parent", ["parentId"])
  .index("by_state", ["state"]),
});
```

## 3. UI/UX Design

### 3.1 Layout Structure

#### Expanded Sidebar
```
┌─────────────────────────────────────────────────────────┐
│ Navigation Bar (Dark Mode Toggle, Title)               │
├─────────────┬───────────────────────────────────────────┤
│             │                                           │
│  Sidebar    │           ReactFlow Canvas                │
│             │                                           │
│ Workspaces: │  ┌─────┐    ┌─────┐    ┌─────┐           │
│ ┌─────────┐ │  │Task │───▶│Task │───▶│Task │           │
│ │Projects │ │  │  A  │    │  B  │    │  C  │           │
│ │ Agents  │ │  └─────┘    └─────┘    └─────┘           │
│ │Settings │ │                                           │
│ └─────────┘ │  ┌─────┐                                  │
│             │  │Task │                                  │
│ Content:    │  │  D  │                                  │
│ - List      │  └─────┘                                  │
│ - Filters   │                                           │
│ - Controls  │                                           │
└─────────────┴───────────────────────────────────────────┘
```

#### Collapsed Sidebar
```
┌─────────────────────────────────────────────────────────┐
│ Navigation Bar (Dark Mode Toggle, Title)               │
├───┬─────────────────────────────────────────────────────┤
│   │                                                     │
│ ⚡ │              ReactFlow Canvas                       │
│ 👥 │                                                     │
│ ⚙️ │  ┌─────┐    ┌─────┐    ┌─────┐                     │
│   │  │Task │───▶│Task │───▶│Task │                     │
│   │  │  A  │    │  B  │    │  C  │                     │
│   │  └─────┘    └─────┘    └─────┘                     │
│   │                                                     │
│   │  ┌─────┐                                            │
│   │  │Task │                                            │
│   │  │  D  │                                            │
│   │  └─────┘                                            │
└───┴─────────────────────────────────────────────────────┘
```

### 3.2 Node Design
- **Task Nodes**: Rechteckige Karten mit State-abhängigen Farben
- **Project Nodes**: Größere Container für Root-Level Visualisierung
- **State Colors**:
  - Pending: Gray
  - In-Progress: Blue
  - Completed: Green
  - Blocked: Red
  - Cancelled: Dark Gray

### 3.3 Edge Design
- **Dependencies**: Gerichtete Pfeile von Dependency zu Task
- **Hierarchy**: Gestrichelte Linien für Parent-Child Beziehungen
- **State-based Styling**: Verschiedene Farben je nach Task-Status

## 4. Auto-Layout Algorithm

### 4.1 Layout Strategy
1. **Hierarchical Layout**: Dagre-basiertes Layout für Task-Hierarchien
2. **Dependency Layout**: Zusätzliche Berücksichtigung von Dependencies
3. **Collision Detection**: Überlappungsvermeidung
4. **Smooth Transitions**: Animierte Positionsänderungen

### 4.2 Layout Implementation
```typescript
interface LayoutOptions {
  direction: 'TB' | 'LR' | 'BT' | 'RL';
  nodeSpacing: number;
  rankSpacing: number;
  edgeSpacing: number;
}

function calculateLayout(
  tasks: Task[], 
  dependencies: Edge[], 
  options: LayoutOptions
): { nodes: Node[], edges: Edge[] }
```

## 5. Real-time Features

### 5.1 Convex Integration
- **Live Queries**: Automatische UI-Updates bei Datenänderungen
- **Optimistic Updates**: Sofortige UI-Reaktion mit Server-Synchronisation
- **Conflict Resolution**: Handling von gleichzeitigen Änderungen
- **Offline Support**: Lokale Änderungen mit Sync bei Reconnect

### 5.2 Real-time Operations
- Task State Updates
- Position Changes (Drag & Drop)
- New Task/Project Creation
- Dependency Modifications
- Progress Updates

## 6. Development Roadmap

### Phase 1: Foundation (Week 1)
- [ ] Project Setup (Vite + TypeScript + Tailwind)
- [ ] Convex Integration & Schema Definition
- [ ] Zustand Store Architecture Setup
- [ ] Multi-Workspace Sidebar (Collapsible)
- [ ] Basic UI Layout (Navbar + Sidebar + Canvas)
- [ ] shadcn/ui Component Setup
- [ ] Dark Mode Implementation
- [ ] UUID Type System Setup
- [ ] Dummy Data Generation

### Phase 2: Core Visualization (Week 2)
- [ ] ReactFlow Integration
- [ ] Basic Node Components (Task, Project)
- [ ] Edge Components (Dependencies, Hierarchy)
- [ ] Data Fetching with Convex Hooks
- [ ] Project Sidebar with Navigation
- [ ] Basic CRUD Operations

### Phase 3: Advanced Features (Week 3)
- [ ] Auto-Layout Algorithm Implementation
- [ ] Drag & Drop Functionality
- [ ] Real-time Position Updates
- [ ] Task State Management
- [ ] Dependency Management UI
- [ ] Progress Visualization

### Phase 4: Polish & Optimization (Week 4)
- [ ] Performance Optimization
- [ ] Advanced Filtering & Search
- [ ] Keyboard Shortcuts
- [ ] Export/Import Functionality
- [ ] Error Handling & Loading States
- [ ] Testing & Documentation

## 7. Technical Considerations

### 7.1 Performance
- **Virtualization**: Für große Task-Mengen
- **Memoization**: React.memo für Node-Components
- **Debouncing**: Für Position-Updates
- **Lazy Loading**: Für große Projekt-Hierarchien

### 7.2 Type Safety
- **UUID Validation**: Runtime-Validierung von UUID-Formaten
- **Strict TypeScript**: Keine `any` Types
- **Schema Validation**: Convex Schema als Single Source of Truth
- **Type Guards**: Für API-Responses

### 7.3 Error Handling
- **Network Errors**: Retry-Mechanismen
- **Data Conflicts**: Optimistic Update Rollbacks
- **Validation Errors**: User-friendly Error Messages
- **Fallback States**: Graceful Degradation

## 8. Questions & Clarifications Needed

### 8.1 Business Logic
1. **Agent Assignment**: Sollen Agent-Informationen visualisiert werden? Wenn ja, wie?
2. **Task Complexity**: Wie soll Complexity (1-10) visuell dargestellt werden?
3. **Time Estimates**: Sollen Schätzungen in Nodes angezeigt werden?
4. **Permissions**: Gibt es verschiedene User-Rollen mit unterschiedlichen Berechtigungen?

### 8.2 UI/UX
1. **Canvas Size**: Soll es eine maximale Canvas-Größe geben?
2. **Zoom Levels**: Min/Max Zoom-Faktoren?
3. **Mobile Support**: Ist Responsive Design erforderlich?
4. **Keyboard Navigation**: Welche Shortcuts sind gewünscht?

### 8.3 Data Integration
1. **API Integration**: Soll die App mit der bestehenden Go-API kommunizieren?
2. **Data Migration**: Gibt es bestehende Daten, die importiert werden müssen?
3. **Sync Strategy**: Wie oft sollen Daten synchronisiert werden?
4. **Offline Mode**: Ist Offline-Funktionalität erforderlich?

### 8.4 Technical
1. **Authentication**: Wird User-Authentication benötigt?
2. **Multi-tenancy**: Sollen verschiedene Teams/Organisationen unterstützt werden?
3. **Export Formats**: Welche Export-Formate sind gewünscht (PNG, SVG, JSON)?
4. **Browser Support**: Welche Browser-Versionen müssen unterstützt werden?

## 9. Success Metrics

### 9.1 Performance Metrics
- Initial Load Time < 2s
- Smooth 60fps Animations
- Real-time Update Latency < 100ms
- Memory Usage < 100MB für 1000+ Tasks

### 9.2 User Experience
- Intuitive Drag & Drop
- Responsive UI (< 16ms frame time)
- Zero Data Loss
- Consistent Dark/Light Mode

### 9.3 Technical Quality
- 100% TypeScript Coverage
- Zero Runtime Type Errors
- Comprehensive Error Handling
- Automated Testing Coverage > 80%

---

**Next Steps**: Bitte die Fragen in Abschnitt 8 klären, dann kann mit der Implementierung begonnen werden.