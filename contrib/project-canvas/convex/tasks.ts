/**
 * Convex Task Queries and Mutations
 * Handles task CRUD operations and real-time updates
 */

import { query, mutation } from "./_generated/server";
import { v } from "convex/values";
import { generateUUID } from "../src/utils/uuid";

// Get all tasks for a project
export const listByProject = query({
  args: { 
    projectId: v.string(), // UUID string instead of v.id("projects")
    includeCompleted: v.optional(v.boolean()),
  },
  handler: async (ctx, args) => {
    let query = ctx.db
      .query("tasks")
      .withIndex("by_project", (q) => q.eq("projectId", args.projectId));

    const tasks = await query.collect();

    // Filter out completed tasks if requested
    if (args.includeCompleted === false) {
      return tasks.filter(task => task.state !== "completed");
    }

    return tasks;
  },
});

// Get a specific task by ID
export const get = query({
  args: { id: v.string() }, // UUID string
  handler: async (ctx, args) => {
    const task = await ctx.db
      .query("tasks")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();
    return task;
  },
});

// Get tasks assigned to a specific agent
export const listByAgent = query({
  args: { agentId: v.string() },
  handler: async (ctx, args) => {
    const tasks = await ctx.db
      .query("tasks")
      .withIndex("by_assigned_agent", (q) => q.eq("assignedAgent", args.agentId))
      .collect();

    return tasks.filter(task => 
      task.state === "pending" || task.state === "in-progress"
    );
  },
});

// Get root tasks (no parent) for a project
export const listRootTasks = query({
  args: { projectId: v.string() }, // UUID string instead of v.id("projects")
  handler: async (ctx, args) => {
    const tasks = await ctx.db
      .query("tasks")
      .withIndex("by_project", (q) => q.eq("projectId", args.projectId))
      .collect();

    return tasks.filter(task => task.parentId === undefined);
  },
});

// Create a new task (LLM only)
export const create = mutation({
  args: {
    projectId: v.string(), // UUID string instead of v.id("projects")
    parentId: v.optional(v.string()), // UUID string
    title: v.string(),
    description: v.string(),
    complexity: v.number(),
    depth: v.number(),
    estimate: v.optional(v.number()),
    assignedAgent: v.optional(v.string()),
    dependencies: v.optional(v.array(v.string())), // Array of task UUIDs
  },
  handler: async (ctx, args) => {
    // Generate a unique ID for the task
    const id = generateUUID();
    
    const taskId = await ctx.db.insert("tasks", {
      id: id,
      projectId: args.projectId,
      parentId: args.parentId,
      title: args.title,
      description: args.description,
      state: "pending",
      complexity: args.complexity,
      depth: args.depth,
      estimate: args.estimate,
      assignedAgent: args.assignedAgent,
      dependencies: args.dependencies || [],
      dependents: [],
      createdAt: Date.now(),
      updatedAt: Date.now(),
    });

    // Update dependent tasks
    if (args.dependencies && args.dependencies.length > 0) {
      for (const depUuid of args.dependencies) {
        const depTask = await ctx.db
          .query("tasks")
          .filter((q) => q.eq(q.field("id"), depUuid))
          .first();
        if (depTask) {
          await ctx.db.patch(depTask._id, {
            dependents: [...depTask.dependents, id],
            updatedAt: Date.now(),
          });
        }
      }
    }

    // Update project statistics
    await updateProjectStats(ctx, args.projectId);

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "task_created",
      entityId: taskId,
      data: { 
        action: "created", 
        title: args.title, 
        projectId: args.projectId 
      },
      timestamp: Date.now(),
    });

    return taskId;
  },
});

// Update task editable fields (frontend)
export const updateEditableFields = mutation({
  args: {
    id: v.string(), // UUID string
    title: v.optional(v.string()),
    description: v.optional(v.string()),
  },
  handler: async (ctx, args) => {
    const { id, ...updates } = args;
    
    // Find task by UUID
    const task = await ctx.db
      .query("tasks")
      .filter((q) => q.eq(q.field("id"), id))
      .first();
    if (!task) throw new Error("Task not found");
    
    const updateData: any = { updatedAt: Date.now() };
    if (updates.title !== undefined) updateData.title = updates.title;
    if (updates.description !== undefined) updateData.description = updates.description;

    await ctx.db.patch(task._id, updateData);

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "task_updated",
      entityId: id,
      data: { action: "updated", fields: Object.keys(updates) },
      timestamp: Date.now(),
    });

    return id;
  },
});

// Update task position (UI only)
export const updatePosition = mutation({
  args: {
    id: v.string(), // UUID string
    positionX: v.number(),
    positionY: v.number(),
  },
  handler: async (ctx, args) => {
    // Find task by UUID
    const task = await ctx.db
      .query("tasks")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();
    if (!task) throw new Error("Task not found");
    
    await ctx.db.patch(task._id, {
      positionX: args.positionX,
      positionY: args.positionY,
      updatedAt: Date.now(),
    });

    // Emit real-time event (throttled)
    await ctx.db.insert("events", {
      type: "task_position_changed",
      entityId: args.id,
      data: { 
        positionX: args.positionX, 
        positionY: args.positionY 
      },
      timestamp: Date.now(),
    });

    return args.id;
  },
});

// Update task state (LLM only)
export const updateState = mutation({
  args: {
    id: v.string(), // UUID string
    state: v.union(
      v.literal("pending"),
      v.literal("in-progress"),
      v.literal("completed"),
      v.literal("blocked"),
      v.literal("cancelled")
    ),
  },
  handler: async (ctx, args) => {
    // Find task by UUID
    const task = await ctx.db
      .query("tasks")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();
    if (!task) throw new Error("Task not found");

    const updateData: any = {
      state: args.state,
      updatedAt: Date.now(),
    };

    // Set completion timestamp
    if (args.state === "completed") {
      updateData.completedAt = Date.now();
    } else if (task.completedAt) {
      updateData.completedAt = undefined;
    }

    await ctx.db.patch(task._id, updateData);

    // Update project statistics
    await updateProjectStats(ctx, task.projectId);

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "task_updated",
      entityId: args.id,
      data: { 
        action: "state_changed", 
        oldState: task.state, 
        newState: args.state 
      },
      timestamp: Date.now(),
    });

    return args.id;
  },
});

// Delete task (LLM only)
export const remove = mutation({
  args: { id: v.string() }, // UUID string
  handler: async (ctx, args) => {
    // Find task by UUID
    const task = await ctx.db
      .query("tasks")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();
    if (!task) throw new Error("Task not found");

    // Remove from dependent tasks
    for (const depUuid of task.dependents) {
      const depTask = await ctx.db
        .query("tasks")
        .filter((q) => q.eq(q.field("id"), depUuid))
        .first();
      if (depTask) {
        await ctx.db.patch(depTask._id, {
          dependencies: depTask.dependencies.filter(uuid => uuid !== args.id),
          updatedAt: Date.now(),
        });
      }
    }

    // Remove from dependency tasks
    for (const depUuid of task.dependencies) {
      const depTask = await ctx.db
        .query("tasks")
        .filter((q) => q.eq(q.field("id"), depUuid))
        .first();
      if (depTask) {
        await ctx.db.patch(depTask._id, {
          dependents: depTask.dependents.filter(uuid => uuid !== args.id),
          updatedAt: Date.now(),
        });
      }
    }

    await ctx.db.delete(task._id);

    // Update project statistics
    await updateProjectStats(ctx, task.projectId);

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "task_deleted",
      entityId: args.id,
      data: { action: "deleted", title: task.title },
      timestamp: Date.now(),
    });

    return args.id;
  },
});

// Helper function to update project statistics
async function updateProjectStats(ctx: any, projectId: string) {
  const tasks = await ctx.db
    .query("tasks")
    .withIndex("by_project", (q: any) => q.eq("projectId", projectId))
    .collect();

  const totalTasks = tasks.length;
  const completedTasks = tasks.filter((task: any) => task.state === "completed").length;

  // Find the project by UUID and update stats
  const project = await ctx.db
    .query("projects")
    .filter((q: any) => q.eq(q.field("id"), projectId))
    .first();
    
  if (project) {
    await ctx.db.patch(project._id, {
      totalTasks,
      completedTasks,
      progress: totalTasks > 0 ? Math.round((completedTasks / totalTasks) * 1000) / 10 : 0,
    });
  }
}