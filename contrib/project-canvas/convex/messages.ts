/**
 * Convex Message Queries and Mutations
 * Handles message CRUD operations and real-time updates
 */

import { query, mutation } from "./_generated/server";
import { v } from "convex/values";
import { generateUUID } from "../src/utils/uuid";

export const MessageEventTypes = {
  ASSISTANT: "assistant",
  TOOL_CALL: "tool.call",
  TOOL_RESPONSE: "tool.response",
  REASONING: "reasoning",
  INTER_AGENT: "inter_agent",
} as const;

export const MessageEventType = v.union(
  v.literal(MessageEventTypes.ASSISTANT),
  v.literal(MessageEventTypes.TOOL_CALL),
  v.literal(MessageEventTypes.TOOL_RESPONSE),
  v.literal(MessageEventTypes.REASONING),
  v.literal(MessageEventTypes.INTER_AGENT)
);

export const InterAgentEventTypes = {
  COMMUNICATION: "communication",
  RECEIVED: "received",
} as const;

export const MessageInteragentType = v.union(
  v.literal(InterAgentEventTypes.COMMUNICATION),
  v.literal(InterAgentEventTypes.RECEIVED)
);

export const MessageRoles = {
  SYSTEM: "system",
  USER: "user",
  ASSISTANT: "assistant",
  TOOL: "tool",
} as const;

export const MessageRoleType = v.union(
  v.literal(MessageRoles.ASSISTANT),
  v.literal(MessageRoles.USER),
  v.literal(MessageRoles.SYSTEM),
  v.literal(MessageRoles.TOOL)
);

export const MessageRoutingInfo = v.object({
  fromAgentId: v.string(), // is UUID
  toAgentId: v.string(), // is UUID
  sessionId: v.string(), // is UUID
  streaming: v.optional(v.boolean()),
});

export const MessageUsageInfo = v.object({
  promptTokenCount: v.number(), // is int
  candidatesTokenCount: v.number(), // is int
  totalTokenCount: v.number(), // is int
});

// Get all messages for a specific invocation
export const listByInvocation = query({
  args: {
    invocationId: v.string(),
    limit: v.number(),
  },
  handler: async (ctx, args) => {
    return ctx.db
      .query("messages")
      .withIndex("by_invocation_id", (q) =>
        q.eq("invocationId", args.invocationId)
      )
      .order("desc")
      .take(args.limit);
  },
});

// Get a specific message by ID
export const get = query({
  args: { id: v.string() },
  handler: async (ctx, args) => {
    return await ctx.db
      .query("messages")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();
  },
});

// Get messages by role
export const listByRole = query({
  args: {
    role: MessageRoleType,
    limit: v.number(),
  },
  handler: async (ctx, args) => {
    return ctx.db
      .query("messages")
      .withIndex("by_role", (q) => q.eq("role", args.role))
      .order("desc")
      .take(args.limit);
  },
});

// Get messages by type
export const listByType = query({
  args: {
    type: MessageEventType,
    limit: v.number(),
  },
  handler: async (ctx, args) => {
    return ctx.db
      .query("messages")
      .withIndex("by_type", (q) => q.eq("type", args.type))
      .order("desc")
      .take(args.limit);
  },
});

// Get recent messages across all invocations
// export const getRecentMessages = query({
//   args: {
//     limit: v.number(),
//     since: v.optional(v.number()), // Unix timestamp
//   },
//   handler: async (ctx, args) => {
// //     return ctx.db.query("messages").order("desc").filter((m)-> m.t).take(args.limit);
//  let messages = await query.collect();

//     // Filter by timestamp if provided
//     if (args.since) {
//       messages = messages.filter((msg) => msg.timestamp >= args.since!);
//     }
// //   },
// // });

// Get inter-agent messages
// export const getInterAgentMessages = query({
//   args: {
//     fromAgentId: v.optional(v.string()),
//     toAgentId: v.optional(v.string()),
//     sessionId: v.optional(v.string()),
//     limit: v.optional(v.number()),
//   },
//   handler: async (ctx, args) => {
//     let query = ctx.db
//       .query("messages")
//       .filter((q) => q.neq(q.field("routing"), undefined))
//       .order("desc");

//     if (args.limit) {
//       query = query.take(args.limit);
//     }

//     return query;
//   },
// });

// Create a new message
export const create = mutation({
  args: {
    invocationId: v.string(),
    routing: v.optional(
      v.object({
        fromAgentId: v.string(),
        toAgentId: v.string(),
        sessionId: v.string(),
        streaming: v.optional(v.boolean()),
      })
    ),
    usage: v.optional(
      v.object({
        promptTokenCount: v.number(),
        candidatesTokenCount: v.number(),
        totalTokenCount: v.number(),
      })
    ),
    interagentType: v.optional(
      v.union(v.literal("communication"), v.literal("received"))
    ),
    done: v.boolean(),
    partial: v.boolean(),
    content: v.string(),
    author: v.string(),
    type: v.union(
      v.literal("assistant"),
      v.literal("tool.call"),
      v.literal("tool.response"),
      v.literal("reasoning"),
      v.literal("inter_agent")
    ),
    role: v.union(
      v.literal("assistant"),
      v.literal("user"),
      v.literal("system"),
      v.literal("tool")
    ),
  },
  handler: async (ctx, args) => {
    const messageId = generateUUID();
    const timestamp = Date.now();

    const messageData = {
      id: messageId,
      timestamp,
      ...args,
    };

    const message = await ctx.db.insert("messages", messageData);

    // Emit event for real-time updates
    await ctx.db.insert("events", {
      type: "message_created",
      entityId: messageId,
      data: { messageId, invocationId: args.invocationId },
      timestamp,
    });

    return message;
  },
});

// Update a message (useful for streaming updates)
export const update = mutation({
  args: {
    id: v.string(),
    content: v.optional(v.string()),
    done: v.optional(v.boolean()),
    partial: v.optional(v.boolean()),
    usage: v.optional(
      v.object({
        promptTokenCount: v.number(),
        candidatesTokenCount: v.number(),
        totalTokenCount: v.number(),
      })
    ),
  },
  handler: async (ctx, args) => {
    const { id, ...updates } = args;

    // Find the message
    const existingMessage = await ctx.db
      .query("messages")
      .filter((q) => q.eq(q.field("id"), id))
      .first();

    if (!existingMessage) {
      throw new Error(`Message with id ${id} not found`);
    }

    // Update the message
    const updatedMessage = await ctx.db.patch(existingMessage._id, {
      ...updates,
      timestamp: Date.now(), // Update timestamp
    });

    // Emit event for real-time updates
    await ctx.db.insert("events", {
      type: "message_updated",
      entityId: id,
      data: { messageId: id, updates },
      timestamp: Date.now(),
    });

    return updatedMessage;
  },
});

// Delete a message
export const remove = mutation({
  args: { id: v.string() },
  handler: async (ctx, args) => {
    const message = await ctx.db
      .query("messages")
      .filter((q) => q.eq(q.field("id"), args.id))
      .first();

    if (!message) {
      throw new Error(`Message with id ${args.id} not found`);
    }

    await ctx.db.delete(message._id);

    // Emit event for real-time updates
    await ctx.db.insert("events", {
      type: "message_deleted",
      entityId: args.id,
      data: { messageId: args.id },
      timestamp: Date.now(),
    });

    return { success: true };
  },
});

// Batch delete messages by invocation ID
export const removeByInvocation = mutation({
  args: { invocationId: v.string() },
  handler: async (ctx, args) => {
    const messages = await ctx.db
      .query("messages")
      .withIndex("by_invocation_id", (q) =>
        q.eq("invocationId", args.invocationId)
      )
      .collect();

    const deletedIds = [];
    for (const message of messages) {
      await ctx.db.delete(message._id);
      deletedIds.push(message.id);
    }

    // Emit event for real-time updates
    await ctx.db.insert("events", {
      type: "messages_bulk_deleted",
      entityId: args.invocationId,
      data: { invocationId: args.invocationId, deletedIds },
      timestamp: Date.now(),
    });

    return { success: true, deletedCount: deletedIds.length };
  },
});

// Get message statistics
export const getStats = query({
  args: {
    invocationId: v.optional(v.string()),
    since: v.optional(v.number()), // Unix timestamp
  },
  handler: async (ctx, args) => {
    let query;

    if (args.invocationId) {
      query = ctx.db
        .query("messages")
        .withIndex("by_invocation_id", (q) =>
          q.eq("invocationId", args.invocationId!)
        );
    } else {
      query = ctx.db.query("messages");
    }

    let messages = await query.collect();

    // Filter by timestamp if provided
    if (args.since) {
      messages = messages.filter((msg) => msg.timestamp >= args.since!);
    }

    // Calculate statistics
    const totalMessages = messages.length;
    const byRole = messages.reduce(
      (acc, msg) => {
        acc[msg.role] = (acc[msg.role] || 0) + 1;
        return acc;
      },
      {} as Record<string, number>
    );

    const byType = messages.reduce(
      (acc, msg) => {
        acc[msg.type] = (acc[msg.type] || 0) + 1;
        return acc;
      },
      {} as Record<string, number>
    );

    const totalTokens = messages.reduce((sum, msg) => {
      return sum + (msg.usage?.totalTokenCount || 0);
    }, 0);

    const interAgentMessages = messages.filter(
      (msg) => msg.routing !== undefined
    ).length;
    const completedMessages = messages.filter((msg) => msg.done).length;

    return {
      totalMessages,
      byRole,
      byType,
      totalTokens,
      interAgentMessages,
      completedMessages,
      completionRate:
        totalMessages > 0 ? (completedMessages / totalMessages) * 100 : 0,
    };
  },
});
