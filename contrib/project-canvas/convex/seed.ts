/**
 * Convex Seed Data
 * Verwendet die master-dummy-data.ts (Go-Model konform)
 * Keine eigenen Daten-Definitionen!
 */

import { mutation } from "./_generated/server";
import { v } from "convex/values";
import {
  masterProjects,
  masterAgents,
  allTasks,
} from "../src/data/master-dummy-data";
import {
  projectToConvexProject,
  taskToConvexTask,
} from "../src/utils/convex-helpers";

export const seedDatabase = mutation({
  args: {},
  handler: async (ctx) => {
    // Check if data already exists
    const existingProjects = await ctx.db.query("projects").collect();
    if (existingProjects.length > 0) {
      return { message: "Database already seeded" };
    }

    // Convert and create projects from master dummy data
    const projectIdMapping: Record<string, string> = {};
    for (const project of masterProjects) {
      const convexProject = projectToConvexProject(project);
      const projectId = await ctx.db.insert("projects", convexProject);
      projectIdMapping[project.id] = projectId;
      console.log(`Created project: ${project.title} -> ${projectId}`);
    }

    // Convert and create agents from master dummy data
    const agentIdMapping: Record<string, string> = {};
    for (const agent of masterAgents) {
      const convexAgent = {
        name: agent.name,
        role: agent.role,
        description: agent.description,
        status: agent.status,
        isStreaming: agent.isStreaming,
        capabilities: agent.capabilities,
        currentTasks: [], // Will be updated after tasks are created
        lastActiveAt: agent.lastActiveAt
          ? agent.lastActiveAt.getTime()
          : undefined,
        id: agent.id, // Keep original ID for compatibility
      };
      const agentId = await ctx.db.insert("agents", convexAgent);
      agentIdMapping[agent.id] = agentId;
    }

    // Convert and create tasks from master dummy data
    const taskIdMapping: Record<string, string> = {};
    for (const task of allTasks) {
      // Ensure projectId is mapped correctly
      const mappedProjectId = projectIdMapping[task.projectId];
      if (!mappedProjectId) {
        console.error(`Project ID ${task.projectId} not found in mapping for task ${task.title}`);
        continue;
      }

      const convexTask = {
        ...taskToConvexTask(task),
        projectId: mappedProjectId,
        parentId: task.parentId ? taskIdMapping[task.parentId] : undefined,
        assignedAgent: task.assignedAgent ? agentIdMapping[task.assignedAgent] : undefined,
        dependencies: task.dependencies.map(depId => taskIdMapping[depId]).filter(Boolean),
        dependents: task.dependents.map(depId => taskIdMapping[depId]).filter(Boolean),
      };
      
      console.log(`Creating task: ${task.title} with projectId: ${mappedProjectId}`);
      const taskId = await ctx.db.insert("tasks", convexTask);
      taskIdMapping[task.id] = taskId;
    }

    // Update task dependencies and dependents with correct Convex IDs
    for (const task of allTasks) {
      const convexTaskId = taskIdMapping[task.id];
      if (convexTaskId) {
        const dependencies = task.dependencies
          .map((depId) => taskIdMapping[depId])
          .filter(Boolean);
        const dependents = task.dependents
          .map((depId) => taskIdMapping[depId])
          .filter(Boolean);

        if (dependencies.length > 0 || dependents.length > 0) {
          await ctx.db.patch(convexTaskId, {
            dependencies,
            dependents,
          });
        }
      }
    }

    // Update agent current tasks with correct Convex task IDs
    for (const agent of masterAgents) {
      const convexAgentId = agentIdMapping[agent.id];
      if (convexAgentId && agent.currentTasks.length > 0) {
        const currentTasks = agent.currentTasks
          .map((taskId) => taskIdMapping[taskId])
          .filter(Boolean);
        await ctx.db.patch(convexAgentId, { currentTasks });
      }
    }

    return {
      message: "Database seeded with master dummy data (Go-Model konform)",
      projectCount: masterProjects.length,
      agentCount: masterAgents.length,
      taskCount: allTasks.length,
    };
  },
});
