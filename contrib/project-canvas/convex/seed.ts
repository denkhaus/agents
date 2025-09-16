/**
 * Convex Seed Data
 * Verwendet die master-dummy-data.ts (Go-Model konform)
 * Keine eigenen Daten-Definitionen!
 */

import { mutation } from "./_generated/server";
import {
  masterProjects,
  masterAgents,
  allTasks,
} from "../src/data/master-dummy-data";

export const seedDatabase = mutation({
  args: {},
  handler: async (ctx) => {
    // Check if data already exists
    const existingProjects = await ctx.db.query("projects").collect();
    if (existingProjects.length > 0) {
      return { message: "Database already seeded" };
    }

    // Convert and create projects from master dummy data (using UUIDs directly)
    for (const project of masterProjects) {
      const convexProject = {
        id: project.id,
        title: project.title,
        description: project.description,
        totalTasks: project.totalTasks,
        completedTasks: project.completedTasks,
        progress: project.progress,
        createdAt: project.createdAt.getTime(),
        updatedAt: project.updatedAt.getTime(),
      };
      await ctx.db.insert("projects", convexProject);
      console.log(`Created project: ${project.title} with UUID: ${project.id}`);
    }

    // Convert and create agents from master dummy data (using UUIDs directly)
    for (const agent of masterAgents) {
      const convexAgent = {
        name: agent.name,
        role: agent.role,
        description: agent.description,
        status: agent.status,
        isStreaming: agent.isStreaming,
        capabilities: agent.capabilities,
        currentTasks: agent.currentTasks, // Keep UUIDs directly
        lastActiveAt: agent.lastActiveAt
          ? agent.lastActiveAt.getTime()
          : undefined,
        id: agent.id, // Keep original UUID
      };
      await ctx.db.insert("agents", convexAgent);
      console.log(`Created agent: ${agent.name} with UUID: ${agent.id}`);
    }

    // Convert and create tasks from master dummy data (using UUIDs directly)
    for (const task of allTasks) {
      const convexTask = {
        id: task.id,
        projectId: task.projectId,
        parentId: task.parentId,
        title: task.title,
        description: task.description,
        state: task.state,
        complexity: task.complexity,
        depth: task.depth,
        estimate: task.estimate,
        assignedAgent: task.assignedAgent,
        dependencies: task.dependencies,
        dependents: task.dependents,
        positionX: task.position?.x,
        positionY: task.position?.y,
        createdAt: task.createdAt.getTime(),
        completedAt: task.completedAt?.getTime(),
        updatedAt: task.updatedAt.getTime(),
      };
      
      console.log(`Creating task: ${task.title} with UUID: ${task.id}`);
      await ctx.db.insert("tasks", convexTask);
    }

    // No need for ID mapping updates since we use UUIDs directly

    return {
      message: "Database seeded with master dummy data (Go-Model konform)",
      projectCount: masterProjects.length,
      agentCount: masterAgents.length,
      taskCount: allTasks.length,
    };
  },
});
