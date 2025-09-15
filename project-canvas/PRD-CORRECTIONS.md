# PRD Corrections - Data Model Alignment

## ✅ **CRITICAL CORRECTIONS MADE**

### **1. Data Model Alignment with Go Backend**
- ✅ **UUID Types**: All IDs are now properly typed as UUID, not strings
- ✅ **Agent Assignment**: `assignedAgent` field uses UUID, not agent names
- ✅ **Dependency Arrays**: All dependency/dependent fields use UUID arrays
- ✅ **Date Handling**: Proper Date objects with Go time.Time conversion utilities

### **2. Removed CRUD Operations (LLM-Only)**
- ❌ **No Project Creation**: Projects created only via LLM
- ❌ **No Task Creation**: Tasks created only via LLM  
- ❌ **No Deletion**: No delete operations in frontend
- ✅ **Limited Editing**: Only title/description editable in frontend
- ✅ **UI Updates**: Position changes for drag & drop allowed

### **3. Markdown Support**
- ✅ **Markdown Renderer**: Safe rendering of task/project descriptions
- ✅ **Editable Markdown**: In-place editing with live preview
- ✅ **XSS Protection**: Sanitized HTML output
- ✅ **Plain Text Extraction**: For previews and search

### **4. Consistent Dummy Data**
- ✅ **Valid UUIDs**: All agent assignments use proper UUIDs
- ✅ **Realistic Markdown**: Task descriptions contain actual markdown
- ✅ **Proper Dependencies**: Correct dependency chains
- ✅ **Go Model Structure**: Exactly matches backend structure

## 📋 **UPDATED ARCHITECTURE**

### **Frontend Responsibilities (READ-ONLY + LIMITED EDIT)**
```typescript
// ✅ ALLOWED Operations
- Display projects and tasks
- Edit title/description fields (markdown support)
- Update task positions (drag & drop)
- Filter and search
- Real-time updates via Convex

// ❌ FORBIDDEN Operations  
- Create projects/tasks
- Delete projects/tasks
- Change task states
- Modify dependencies
- Assign/unassign agents
```

### **LLM Responsibilities (FULL CRUD)**
```typescript
// ✅ LLM-ONLY Operations
- Create new projects
- Create new tasks
- Set task dependencies
- Assign agents to tasks
- Update task states
- Modify task complexity/estimates
- Delete projects/tasks
```

### **Data Flow Architecture**
```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│   LLM       │───▶│  Go Backend  │───▶│  Convex     │
│ (Creates)   │    │ (Processes)  │    │ (Syncs)     │
└─────────────┘    └──────────────┘    └─────────────┘
                                              │
                                              ▼
                                    ┌─────────────┐
                                    │  Frontend   │
                                    │ (Displays)  │
                                    └─────────────┘
```

## 🎯 **CORRECTED TYPE DEFINITIONS**

### **Task Model (Go-Aligned)**
```typescript
interface Task {
  id: UUID;                    // ✅ UUID type
  projectId: UUID;             // ✅ UUID type  
  parentId?: UUID;             // ✅ Optional UUID
  title: string;               // ✅ Editable
  description: string;         // ✅ Editable (Markdown)
  state: TaskState;            // ❌ Read-only (LLM sets)
  complexity: number;          // ❌ Read-only (LLM sets)
  depth: number;               // ❌ Read-only (calculated)
  estimate?: number;           // ❌ Read-only (LLM sets)
  assignedAgent?: UUID;        // ❌ Read-only (LLM assigns)
  dependencies: UUID[];        // ❌ Read-only (LLM manages)
  dependents: UUID[];          // ❌ Read-only (calculated)
  createdAt: Date;             // ❌ Read-only
  updatedAt: Date;             // ❌ Read-only
  completedAt?: Date;          // ❌ Read-only
  position?: { x: number; y: number }; // ✅ UI-only
}
```

### **Project Model (Go-Aligned)**
```typescript
interface Project {
  id: UUID;                    // ✅ UUID type
  title: string;               // ✅ Editable
  description: string;         // ✅ Editable (Markdown)
  createdAt: Date;             // ❌ Read-only
  updatedAt: Date;             // ❌ Read-only
  totalTasks: number;          // ❌ Read-only (calculated)
  completedTasks: number;      // ❌ Read-only (calculated)
  progress: number;            // ❌ Read-only (calculated)
}
```

## 🔧 **IMPLEMENTATION STATUS**

### **✅ COMPLETED**
- [x] Corrected TypeScript types
- [x] Removed CRUD operations from stores
- [x] Added markdown rendering components
- [x] Created consistent dummy data
- [x] Added data validation utilities
- [x] UUID type safety throughout

### **🚧 NEXT STEPS**
- [ ] React components (Navbar, Sidebar, Canvas)
- [ ] ReactFlow integration with corrected types
- [ ] Convex schema matching Go model
- [ ] Package.json and Vite setup
- [ ] shadcn/ui configuration

## 🎯 **KEY PRINCIPLES ESTABLISHED**

1. **Frontend = Display Layer**: No business logic, only presentation
2. **LLM = Business Logic**: All creation, modification, state changes
3. **Go Backend = Source of Truth**: Data structure authority
4. **Convex = Real-time Sync**: Efficient data synchronization
5. **Type Safety = Critical**: UUID types prevent runtime errors

**The foundation is now solid and aligned with the Go backend model!**