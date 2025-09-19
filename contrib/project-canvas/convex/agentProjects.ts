/**
 * Agent Projects - Convex Queries and Mutations
 * Handles team configurations and agent canvas layouts
 */

import { v } from "convex/values";
import { mutation, query } from "./_generated/server";

// Query: List all agent projects
export const list = query({
  args: {},
  handler: async (ctx) => {
    const agentProjects = await ctx.db
      .query("agentProjects")
      .order("desc")
      .collect();
    return agentProjects;
  },
});

// Query: Get agent project by ID
export const get = query({
  args: { id: v.string() },
  handler: async (ctx, args) => {
    const agentProject = await ctx.db
      .query("agentProjects")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();
    return agentProject;
  },
});

// Query: Get agent project with full agent data
export const getWithAgents = query({
  args: { id: v.string() },
  handler: async (ctx, args) => {
    const agentProject = await ctx.db
      .query("agentProjects")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();

    if (!agentProject) {
      return null;
    }

    // Fetch full agent data for each agent in the project
    const agents = await Promise.all(
      agentProject.agentIds.map(async (agentId) => {
        return await ctx.db
          .query("agents")
          .filter((q) => q.eq(q.field("id"), agentId))
          .first();
      })
    );

    return {
      ...agentProject,
      agents: agents.filter(Boolean), // Filter out any null results
    };
  },
});

// Query: Search agent projects by name
export const searchByName = query({
  args: { searchTerm: v.string() },
  handler: async (ctx, args) => {
    const agentProjects = await ctx.db
      .query("agentProjects")
      .order("desc")
      .collect();

    if (!args.searchTerm) {
      return agentProjects;
    }

    const searchLower = args.searchTerm.toLowerCase();
    return agentProjects.filter(
      (project) =>
        project.name.toLowerCase().includes(searchLower) ||
        project.description.toLowerCase().includes(searchLower)
    );
  },
});

// Mutation: Create new agent project
export const create = mutation({
  args: {
    id: v.string(),
    name: v.string(),
    description: v.string(),
    agentIds: v.array(v.string()),
    agentNodes: v.array(v.object({
      id: v.string(),
      type: v.literal("agent"),
      position: v.object({
        x: v.number(),
        y: v.number(),
      }),
      data: v.object({
        isSelected: v.optional(v.boolean()),
      }),
    })),
    connections: v.optional(v.array(v.object({
      id: v.string(),
      source: v.string(),
      target: v.string(),
      type: v.union(
        v.literal("communication"),
        v.literal("hierarchy"),
        v.literal("collaboration")
      ),
      label: v.optional(v.string()),
      data: v.optional(v.object({
        frequency: v.optional(v.number()),
        protocol: v.optional(v.string()),
      })),
    }))),
  },
  handler: async (ctx, args) => {
    const now = Date.now();
    const agentProject = {
      id: args.id,
      name: args.name,
      description: args.description,
      agentIds: args.agentIds,
      agentNodes: args.agentNodes,
      connections: args.connections || [],
      createdAt: now,
      updatedAt: now,
    };

    const result = await ctx.db.insert("agentProjects", agentProject);
    return { _id: result, ...agentProject };
  },
});

// Mutation: Update agent project
export const updateEditableFields = mutation({
  args: {
    id: v.string(),
    name: v.optional(v.string()),
    description: v.optional(v.string()),
  },
  handler: async (ctx, args) => {
    const agentProject = await ctx.db
      .query("agentProjects")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();

    if (!agentProject) {
      throw new Error(`Agent project with id ${args.id} not found`);
    }

    const updates: any = {
      updatedAt: Date.now(),
    };

    if (args.name !== undefined) updates.name = args.name;
    if (args.description !== undefined) updates.description = args.description;

    await ctx.db.patch(agentProject._id, updates);
    return { ...agentProject, ...updates };
  },
});

// Mutation: Update agent nodes (positions, selections)
export const updateAgentNodes = mutation({
  args: {
    id: v.string(),
    agentNodes: v.array(v.object({
      id: v.string(),
      type: v.literal("agent"),
      position: v.object({
        x: v.number(),
        y: v.number(),
      }),
      data: v.object({
        isSelected: v.optional(v.boolean()),
      }),
    })),
  },
  handler: async (ctx, args) => {
    const agentProject = await ctx.db
      .query("agentProjects")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();

    if (!agentProject) {
      throw new Error(`Agent project with id ${args.id} not found`);
    }

    await ctx.db.patch(agentProject._id, {
      agentNodes: args.agentNodes,
      updatedAt: Date.now(),
    });

    return { ...agentProject, agentNodes: args.agentNodes };
  },
});

// Mutation: Update connections
export const updateConnections = mutation({
  args: {
    id: v.string(),
    connections: v.array(v.object({
      id: v.string(),
      source: v.string(),
      target: v.string(),
      type: v.union(
        v.literal("communication"),
        v.literal("hierarchy"),
        v.literal("collaboration")
      ),
      label: v.optional(v.string()),
      data: v.optional(v.object({
        frequency: v.optional(v.number()),
        protocol: v.optional(v.string()),
      })),
    })),
  },
  handler: async (ctx, args) => {
    const agentProject = await ctx.db
      .query("agentProjects")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();

    if (!agentProject) {
      throw new Error(`Agent project with id ${args.id} not found`);
    }

    await ctx.db.patch(agentProject._id, {
      connections: args.connections,
      updatedAt: Date.now(),
    });

    return { ...agentProject, connections: args.connections };
  },
});

// Mutation: Add agent to project
export const addAgent = mutation({
  args: {
    id: v.string(),
    agentId: v.string(),
    position: v.object({
      x: v.number(),
      y: v.number(),
    }),
  },
  handler: async (ctx, args) => {
    const agentProject = await ctx.db
      .query("agentProjects")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();

    if (!agentProject) {
      throw new Error(`Agent project with id ${args.id} not found`);
    }

    // Check if agent already exists in project
    if (agentProject.agentIds.includes(args.agentId)) {
      throw new Error(`Agent ${args.agentId} already exists in project`);
    }

    const newAgentNode = {
      id: args.agentId,
      type: "agent" as const,
      position: args.position,
      data: { isSelected: false },
    };

    const updatedAgentIds = [...agentProject.agentIds, args.agentId];
    const updatedAgentNodes = [...agentProject.agentNodes, newAgentNode];

    await ctx.db.patch(agentProject._id, {
      agentIds: updatedAgentIds,
      agentNodes: updatedAgentNodes,
      updatedAt: Date.now(),
    });

    return {
      ...agentProject,
      agentIds: updatedAgentIds,
      agentNodes: updatedAgentNodes,
    };
  },
});

// Mutation: Remove agent from project
export const removeAgent = mutation({
  args: {
    id: v.string(),
    agentId: v.string(),
  },
  handler: async (ctx, args) => {
    const agentProject = await ctx.db
      .query("agentProjects")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();

    if (!agentProject) {
      throw new Error(`Agent project with id ${args.id} not found`);
    }

    const updatedAgentIds = agentProject.agentIds.filter(
      (id) => id !== args.agentId
    );
    const updatedAgentNodes = agentProject.agentNodes.filter(
      (node) => node.id !== args.agentId
    );
    const updatedConnections = agentProject.connections.filter(
      (conn) => conn.source !== args.agentId && conn.target !== args.agentId
    );

    await ctx.db.patch(agentProject._id, {
      agentIds: updatedAgentIds,
      agentNodes: updatedAgentNodes,
      connections: updatedConnections,
      updatedAt: Date.now(),
    });

    return {
      ...agentProject,
      agentIds: updatedAgentIds,
      agentNodes: updatedAgentNodes,
      connections: updatedConnections,
    };
  },
});

// Mutation: Delete agent project
export const remove = mutation({
  args: { id: v.string() },
  handler: async (ctx, args) => {
    const agentProject = await ctx.db
      .query("agentProjects")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();

    if (!agentProject) {
      throw new Error(`Agent project with id ${args.id} not found`);
    }

    await ctx.db.delete(agentProject._id);
    return { success: true, deletedId: args.id };
  },
});