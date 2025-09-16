/**
 * Test Script für Seed-Funktion
 * Testet die Konvertierung der Dummy-Daten ohne Convex
 */

// Simuliere Convex-Umgebung
const mockConvexContext = {
  db: {
    insert: (table, data) => {
      console.log(`📝 INSERT into ${table}:`, JSON.stringify(data, null, 2));
      return `mock_id_${Math.random().toString(36).substr(2, 9)}`;
    },
    patch: (id, updates) => {
      console.log(`🔄 PATCH ${id}:`, JSON.stringify(updates, null, 2));
    },
    query: (table) => ({
      collect: () => {
        console.log(`🔍 QUERY ${table}`);
        return []; // Leere Datenbank simulieren
      }
    })
  }
};

// Importiere die Dummy-Daten (simuliert)
const masterProjects = [
  {
    id: "proj-1",
    title: "E-Commerce Platform Redesign",
    description: "Complete redesign and modernization",
    totalTasks: 15,
    completedTasks: 4,
    progress: 26.7
  },
  {
    id: "proj-2", 
    title: "Mobile App Development",
    description: "Native mobile application",
    totalTasks: 10,
    completedTasks: 3,
    progress: 30.0
  }
];

const masterAgents = [
  {
    id: "agent-1",
    name: "Design Lead",
    role: "designer",
    description: "Senior UX/UI designer",
    status: "online",
    isStreaming: false,
    capabilities: ["wireframing", "prototyping"],
    currentTasks: [],
    lastActiveAt: new Date()
  },
  {
    id: "agent-2",
    name: "Frontend Developer", 
    role: "coder",
    description: "React/TypeScript specialist",
    status: "busy",
    isStreaming: true,
    capabilities: ["react", "typescript"],
    currentTasks: ["task-1"],
    lastActiveAt: new Date()
  }
];

const allTasks = [
  {
    id: "task-1",
    projectId: "proj-1",
    parentId: undefined,
    title: "Market Research & User Analysis",
    description: "Conduct comprehensive market research",
    state: "completed",
    complexity: 6,
    depth: 0,
    estimate: 2400,
    assignedAgent: "agent-1",
    dependencies: [],
    dependents: ["task-2"],
    position: { x: 100, y: 100 },
    updatedAt: new Date(),
    completedAt: new Date()
  },
  {
    id: "task-2",
    projectId: "proj-1", 
    parentId: undefined,
    title: "Create Wireframes & Prototypes",
    description: "Design wireframes and prototypes",
    state: "in-progress",
    complexity: 7,
    depth: 0,
    estimate: 1800,
    assignedAgent: "agent-2",
    dependencies: ["task-1"],
    dependents: [],
    position: { x: 400, y: 100 },
    updatedAt: new Date()
  }
];

// Simuliere die Seed-Funktion
async function testSeedFunction() {
  console.log("🚀 Testing Seed Function...\n");
  
  const ctx = mockConvexContext;
  
  // Check if data already exists
  const existingProjects = await ctx.db.query("projects").collect();
  if (existingProjects.length > 0) {
    console.log("❌ Database already seeded");
    return { message: "Database already seeded" };
  }

  console.log("📊 Converting and creating projects from master dummy data...");
  
  // Convert and create projects
  const projectIdMapping = {};
  for (const project of masterProjects) {
    const convexProject = {
      title: project.title,
      description: project.description,
      totalTasks: project.totalTasks,
      completedTasks: project.completedTasks,
      progress: project.progress,
    };
    const projectId = await ctx.db.insert("projects", convexProject);
    projectIdMapping[project.id] = projectId;
  }

  console.log("\n👥 Converting and creating agents from master dummy data...");
  
  // Convert and create agents
  const agentIdMapping = {};
  for (const agent of masterAgents) {
    const convexAgent = {
      name: agent.name,
      role: agent.role,
      description: agent.description,
      status: agent.status,
      isStreaming: agent.isStreaming,
      capabilities: agent.capabilities,
      currentTasks: [], // Will be updated after tasks are created
      lastActiveAt: agent.lastActiveAt ? agent.lastActiveAt.getTime() : undefined,
      id: agent.id, // Keep original ID for compatibility
    };
    const agentId = await ctx.db.insert("agents", convexAgent);
    agentIdMapping[agent.id] = agentId;
  }

  console.log("\n📋 Converting and creating tasks from master dummy data...");
  
  // Convert and create tasks
  const taskIdMapping = {};
  for (const task of allTasks) {
    const convexTask = {
      projectId: projectIdMapping[task.projectId],
      parentId: task.parentId ? taskIdMapping[task.parentId] : undefined,
      title: task.title,
      description: task.description,
      state: task.state,
      complexity: task.complexity,
      depth: task.depth,
      estimate: task.estimate,
      assignedAgent: task.assignedAgent ? agentIdMapping[task.assignedAgent] : undefined,
      dependencies: task.dependencies.map(depId => taskIdMapping[depId]).filter(Boolean),
      dependents: task.dependents.map(depId => taskIdMapping[depId]).filter(Boolean),
      positionX: task.position?.x,
      positionY: task.position?.y,
      updatedAt: task.updatedAt.getTime(),
      completedAt: task.completedAt ? task.completedAt.getTime() : undefined,
    };
    const taskId = await ctx.db.insert("tasks", convexTask);
    taskIdMapping[task.id] = taskId;
  }

  console.log("\n🔗 Updating task dependencies and dependents...");
  
  // Update task dependencies and dependents with correct Convex IDs
  for (const task of allTasks) {
    const convexTaskId = taskIdMapping[task.id];
    if (convexTaskId) {
      const dependencies = task.dependencies.map(depId => taskIdMapping[depId]).filter(Boolean);
      const dependents = task.dependents.map(depId => taskIdMapping[depId]).filter(Boolean);
      
      if (dependencies.length > 0 || dependents.length > 0) {
        await ctx.db.patch(convexTaskId, {
          dependencies,
          dependents,
        });
      }
    }
  }

  console.log("\n👤 Updating agent current tasks...");
  
  // Update agent current tasks with correct Convex task IDs
  for (const agent of masterAgents) {
    const convexAgentId = agentIdMapping[agent.id];
    if (convexAgentId && agent.currentTasks.length > 0) {
      const currentTasks = agent.currentTasks.map(taskId => taskIdMapping[taskId]).filter(Boolean);
      await ctx.db.patch(convexAgentId, { currentTasks });
    }
  }

  const result = {
    message: "Database seeded with master dummy data (Go-Model konform)",
    projectCount: masterProjects.length,
    agentCount: masterAgents.length,
    taskCount: allTasks.length,
  };
  
  console.log("\n✅ Seed function completed successfully!");
  console.log("📈 Result:", JSON.stringify(result, null, 2));
  
  return result;
}

// Führe den Test aus
testSeedFunction().catch(console.error);