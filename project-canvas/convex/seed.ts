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
      const convexProject = {
        id: project.id,
        title: project.title,
        description: project.description,
        totalTasks: project.totalTasks,
        completedTasks: project.completedTasks,
        progress: project.progress,
      };
      const projectId = await ctx.db.insert("projects", convexProject);
      projectIdMapping[project.id] = projectId;
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
      const convexTask = {
        id: task.id,
        projectId: projectIdMapping[task.projectId],
        parentId: task.parentId ? taskIdMapping[task.parentId] : undefined,
        title: task.title,
        description: task.description,
        state: task.state,
        complexity: task.complexity,
        depth: task.depth,
        estimate: task.estimate,
        assignedAgent: task.assignedAgent
          ? agentIdMapping[task.assignedAgent]
          : undefined,
        dependencies: task.dependencies
          .map((depId) => taskIdMapping[depId])
          .filter(Boolean),
        dependents: task.dependents
          .map((depId) => taskIdMapping[depId])
          .filter(Boolean),
        positionX: task.position?.x,
        positionY: task.position?.y,
        updatedAt: task.updatedAt.getTime(),
        completedAt: task.completedAt ? task.completedAt.getTime() : undefined,
      };
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
