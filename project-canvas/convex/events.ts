/**
 * Convex Real-time Events
 * Handles live updates and event streaming
 */

import { query, mutation } from "./_generated/server";
import { v } from "convex/values";

// Get recent events for real-time updates
export const getRecentEvents = query({
  args: { 
    since: v.optional(v.number()), // Unix timestamp
    types: v.optional(v.array(v.string())), // Filter by event types
    limit: v.optional(v.number()),
  },
  handler: async (ctx, args) => {
    let query = ctx.db
      .query("events")
      .withIndex("by_timestamp");

    // Filter by timestamp if provided
    if (args.since) {
      query = query.filter((q) => q.gte(q.field("timestamp"), args.since));
    }

    let events = await query
      .order("desc")
      .take(args.limit || 50);

    // Filter by event types if provided
    if (args.types && args.types.length > 0) {
      events = events.filter(event => args.types!.includes(event.type));
    }

    return events;
  },
});

// Subscribe to events for a specific project
export const subscribeToProject = query({
  args: { 
    projectId: v.string(),
    since: v.optional(v.number()),
  },
  handler: async (ctx, args) => {
    const events = await ctx.db
      .query("events")
      .withIndex("by_timestamp")
      .filter((q) => {
        let filter = q.or(
          q.eq(q.field("entityId"), args.projectId),
          q.eq(q.field("data.projectId"), args.projectId)
        );
        
        if (args.since) {
          filter = q.and(filter, q.gte(q.field("timestamp"), args.since));
        }
        
        return filter;
      })
      .order("desc")
      .take(100);

    return events;
  },
});

// Clean up old events (maintenance)
export const cleanupOldEvents = mutation({
  args: { 
    olderThan: v.number(), // Unix timestamp
  },
  handler: async (ctx, args) => {
    const oldEvents = await ctx.db
      .query("events")
      .withIndex("by_timestamp")
      .filter((q) => q.lt(q.field("timestamp"), args.olderThan))
      .collect();

    let deletedCount = 0;
    for (const event of oldEvents) {
      await ctx.db.delete(event._id);
      deletedCount++;
    }

    return { deletedCount };
  },
});

// Emit a custom event
export const emitEvent = mutation({
  args: {
    type: v.string(),
    entityId: v.string(),
    data: v.any(),
    userId: v.optional(v.string()),
  },
  handler: async (ctx, args) => {
    const eventId = await ctx.db.insert("events", {
      type: args.type as any,
      entityId: args.entityId,
      data: args.data,
      userId: args.userId,
      timestamp: Date.now(),
    });

    return eventId;
  },
});