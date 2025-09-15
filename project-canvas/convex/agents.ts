/**
 * Convex Agent Queries and Mutations
 * Handles agent data and status updates
 */

import { query, mutation } from "./_generated/server";
import { v } from "convex/values";

// Get all agents
export const list = query({
  args: {},
  handler: async (ctx) => {
    const agents = await ctx.db
      .query("agents")
      .order("desc")
      .collect();
    
    return agents;
  },
});

// Get agents by role
export const listByRole = query({
  args: { 
    role: v.union(
      v.literal("supervisor"),
      v.literal("project-manager"),
      v.literal("coder"),
      v.literal("researcher"),
      v.literal("qa-engineer"),
      v.literal("devops"),
      v.literal("designer")
    )
  },
  handler: async (ctx, args) => {
    const agents = await ctx.db
      .query("agents")
      .withIndex("by_role", (q) => q.eq("role", args.role))
      .collect();
    
    return agents;
  },
});

// Get agents by status
export const listByStatus = query({
  args: { 
    status: v.union(
      v.literal("online"),
      v.literal("offline"),
      v.literal("busy"),
      v.literal("idle")
    )
  },
  handler: async (ctx, args) => {
    const agents = await ctx.db
      .query("agents")
      .withIndex("by_status", (q) => q.eq("status", args.status))
      .collect();
    
    return agents;
  },
});

// Get a specific agent by ID
export const get = query({
  args: { id: v.id("agents") },
  handler: async (ctx, args) => {
    const agent = await ctx.db.get(args.id);
    return agent;
  },
});

// Get agent with current tasks
export const getWithTasks = query({
  args: { id: v.id("agents") },
  handler: async (ctx, args) => {
    const agent = await ctx.db.get(args.id);
    if (!agent) return null;

    // Get current tasks
    const tasks = await Promise.all(
      agent.currentTasks.map(taskId => ctx.db.get(taskId))
    );

    const validTasks = tasks.filter(task => task !== null);

    return {
      ...agent,
      tasks: validTasks,
    };
  },
});

// Create a new agent (system only)
export const create = mutation({
  args: {
    name: v.string(),
    role: v.union(
      v.literal("supervisor"),
      v.literal("project-manager"),
      v.literal("coder"),
      v.literal("researcher"),
      v.literal("qa-engineer"),
      v.literal("devops"),
      v.literal("designer")
    ),
    description: v.string(),
    capabilities: v.array(v.string()),
  },
  handler: async (ctx, args) => {
    const agentId = await ctx.db.insert("agents", {
      name: args.name,
      role: args.role,
      description: args.description,
      status: "offline",
      isStreaming: false,
      capabilities: args.capabilities,
      currentTasks: [],
      lastActiveAt: Date.now(),
    });

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "agent_status_changed",
      entityId: agentId,
      data: { action: "created", name: args.name, role: args.role },
      timestamp: Date.now(),
    });

    return agentId;
  },
});

// Update agent status (system only)
export const updateStatus = mutation({
  args: {
    id: v.id("agents"),
    status: v.union(
      v.literal("online"),
      v.literal("offline"),
      v.literal("busy"),
      v.literal("idle")
    ),
    isStreaming: v.optional(v.boolean()),
  },
  handler: async (ctx, args) => {
    const updateData: any = {
      status: args.status,
      lastActiveAt: Date.now(),
    };

    if (args.isStreaming !== undefined) {
      updateData.isStreaming = args.isStreaming;
    }

    await ctx.db.patch(args.id, updateData);

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "agent_status_changed",
      entityId: args.id,
      data: { 
        action: "status_updated", 
        status: args.status,
        isStreaming: args.isStreaming 
      },
      timestamp: Date.now(),
    });

    return args.id;
  },
});

// Assign task to agent (system only)
export const assignTask = mutation({
  args: {
    agentId: v.id("agents"),
    taskId: v.id("tasks"),
  },
  handler: async (ctx, args) => {
    const agent = await ctx.db.get(args.agentId);
    const task = await ctx.db.get(args.taskId);

    if (!agent || !task) {
      throw new Error("Agent or task not found");
    }

    // Add task to agent's current tasks
    await ctx.db.patch(args.agentId, {
      currentTasks: [...agent.currentTasks, args.taskId],
      lastActiveAt: Date.now(),
    });

    // Update task assignment
    await ctx.db.patch(args.taskId, {
      assignedAgent: args.agentId,
      updatedAt: Date.now(),
    });

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "task_updated",
      entityId: args.taskId,
      data: { 
        action: "assigned", 
        agentId: args.agentId,
        agentName: agent.name 
      },
      timestamp: Date.now(),
    });

    return { agentId: args.agentId, taskId: args.taskId };
  },
});

// Unassign task from agent (system only)
export const unassignTask = mutation({
  args: {
    agentId: v.id("agents"),
    taskId: v.id("tasks"),
  },
  handler: async (ctx, args) => {
    const agent = await ctx.db.get(args.agentId);
    const task = await ctx.db.get(args.taskId);

    if (!agent || !task) {
      throw new Error("Agent or task not found");
    }

    // Remove task from agent's current tasks
    await ctx.db.patch(args.agentId, {
      currentTasks: agent.currentTasks.filter(id => id !== args.taskId),
      lastActiveAt: Date.now(),
    });

    // Remove assignment from task
    await ctx.db.patch(args.taskId, {
      assignedAgent: undefined,
      updatedAt: Date.now(),
    });

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "task_updated",
      entityId: args.taskId,
      data: { 
        action: "unassigned", 
        agentId: args.agentId,
        agentName: agent.name 
      },
      timestamp: Date.now(),
    });

    return { agentId: args.agentId, taskId: args.taskId };
  },
});