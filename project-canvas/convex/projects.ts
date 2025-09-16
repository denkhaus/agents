/**
 * Convex Project Queries and Mutations
 * Handles project CRUD operations and real-time updates
 */

import { query, mutation } from "./_generated/server";
import { v } from "convex/values";
import { generateUUID } from "../src/utils/uuid";

// Get all projects
export const list = query({
  args: {},
  handler: async (ctx) => {
    const projects = await ctx.db
      .query("projects")
      .order("desc")
      .collect();
    
    return projects;
  },
});

// Get a specific project by ID
export const get = query({
  args: { id: v.id("projects") },
  handler: async (ctx, args) => {
    const project = await ctx.db.get(args.id);
    return project;
  },
});

// Get project with task statistics
export const getWithStats = query({
  args: { id: v.id("projects") },
  handler: async (ctx, args) => {
    const project = await ctx.db.get(args.id);
    if (!project) return null;

    // Get task statistics
    const tasks = await ctx.db
      .query("tasks")
      .withIndex("by_project", (q) => q.eq("projectId", args.id))
      .collect();

    const totalTasks = tasks.length;
    const completedTasks = tasks.filter(task => task.state === "completed").length;
    const inProgressTasks = tasks.filter(task => task.state === "in-progress").length;
    const blockedTasks = tasks.filter(task => task.state === "blocked").length;
    const pendingTasks = tasks.filter(task => task.state === "pending").length;

    const progress = totalTasks > 0 ? (completedTasks / totalTasks) * 100 : 0;

    return {
      ...project,
      totalTasks,
      completedTasks,
      inProgressTasks,
      blockedTasks,
      pendingTasks,
      progress: Math.round(progress * 10) / 10, // Round to 1 decimal
    };
  },
});

// Create a new project (LLM only)
export const create = mutation({
  args: {
    title: v.string(),
    description: v.string(),
  },
  handler: async (ctx, args) => {
    // Generate a unique ID for the project
    const id = generateUUID();
    
    const projectId = await ctx.db.insert("projects", {
      id: id,
      title: args.title,
      description: args.description,
      totalTasks: 0,
      completedTasks: 0,
      progress: 0,
    });

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "project_updated",
      entityId: projectId,
      data: { action: "created", title: args.title },
      timestamp: Date.now(),
    });

    return projectId;
  },
});

// Update project (limited fields for frontend)
export const updateEditableFields = mutation({
  args: {
    id: v.id("projects"),
    title: v.optional(v.string()),
    description: v.optional(v.string()),
  },
  handler: async (ctx, args) => {
    const { id, ...updates } = args;
    
    // Only update provided fields
    const updateData: any = {};
    if (updates.title !== undefined) updateData.title = updates.title;
    if (updates.description !== undefined) updateData.description = updates.description;

    if (Object.keys(updateData).length === 0) {
      throw new Error("No fields to update");
    }

    await ctx.db.patch(id, updateData);

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "project_updated",
      entityId: id,
      data: { action: "updated", fields: Object.keys(updateData) },
      timestamp: Date.now(),
    });

    return id;
  },
});

// Update project statistics (internal use)
export const updateStats = mutation({
  args: {
    id: v.id("projects"),
    totalTasks: v.number(),
    completedTasks: v.number(),
  },
  handler: async (ctx, args) => {
    const progress = args.totalTasks > 0 ? (args.completedTasks / args.totalTasks) * 100 : 0;

    await ctx.db.patch(args.id, {
      totalTasks: args.totalTasks,
      completedTasks: args.completedTasks,
      progress: Math.round(progress * 10) / 10,
    });

    return args.id;
  },
});

// Delete project (LLM only)
export const remove = mutation({
  args: { id: v.id("projects") },
  handler: async (ctx, args) => {
    // First delete all tasks in the project
    const tasks = await ctx.db
      .query("tasks")
      .withIndex("by_project", (q) => q.eq("projectId", args.id))
      .collect();

    for (const task of tasks) {
      await ctx.db.delete(task._id);
    }

    // Then delete the project
    await ctx.db.delete(args.id);

    // Emit real-time event
    await ctx.db.insert("events", {
      type: "project_updated",
      entityId: args.id,
      data: { action: "deleted" },
      timestamp: Date.now(),
    });

    return args.id;
  },
});