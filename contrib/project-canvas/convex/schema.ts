/**
 * Convex Database Schema
 * Defines the structure for projects, tasks, and agents
 */

import { defineSchema, defineTable } from "convex/server";
import {
  MessageEventType,
  MessageInteragentType,
  MessageRoleType,
  MessageRoutingInfo,
  MessageUsageInfo,
} from "./messages.ts";
import { v } from "convex/values";
import { SCHEMA as SettingsSchema } from "./settings.ts";
import { SCHEMA as AgentProjectsSchema } from "./agentProjects.ts";
import { SCHEMA as EventsSchema } from "./events.ts";

export default defineSchema({
  // Projects table
  projects: defineTable({
    id: v.string(),
    title: v.string(),
    description: v.string(),
    totalTasks: v.number(),
    completedTasks: v.number(),
    progress: v.number(), // 0-100 percentage
    // UI-specific fields
    positionX: v.optional(v.number()),
    positionY: v.optional(v.number()),
    createdAt: v.number(), // Unix timestamp
    updatedAt: v.number(), // Unix timestamp
    // Convex automatically adds _id and _creationTime
  }).index("by_title", ["title"]),

  // Tasks table
  tasks: defineTable({
    id: v.string(),
    projectId: v.string(), // UUID string to match projects.id
    parentId: v.optional(v.string()), // UUID string
    title: v.string(),
    description: v.string(),
    state: v.union(
      v.literal("pending"),
      v.literal("in-progress"),
      v.literal("completed"),
      v.literal("blocked"),
      v.literal("cancelled")
    ),
    complexity: v.number(), // 1-10
    depth: v.number(), // 0 for root tasks
    estimate: v.optional(v.number()), // minutes
    assignedAgent: v.optional(v.string()), // UUID string
    dependencies: v.array(v.string()), // Array of task UUIDs
    dependents: v.array(v.string()), // Array of task UUIDs
    // UI-specific fields
    positionX: v.optional(v.number()),
    positionY: v.optional(v.number()),
    // Timestamps
    createdAt: v.number(), // Unix timestamp
    completedAt: v.optional(v.number()), // Unix timestamp
    updatedAt: v.number(), // Unix timestamp
  })
    .index("by_project", ["projectId"])
    .index("by_parent", ["parentId"])
    .index("by_state", ["state"])
    .index("by_assigned_agent", ["assignedAgent"])
    .index("by_updated", ["updatedAt"]),

  // messages table
  messages: defineTable({
    id: v.string(), // is UUID
    invocationId: v.string(), // is UUID
    timestamp: v.number(), // Unix timestamp
    routing: v.optional(MessageRoutingInfo),
    usage: v.optional(MessageUsageInfo),
    interagentType: v.optional(MessageInteragentType),
    done: v.boolean(),
    partial: v.boolean(),
    content: v.string(),
    author: v.string(),
    type: MessageEventType,
    role: MessageRoleType,
  })
    .index("by_role", ["role"])
    .index("by_type", ["type"])
    .index("by_invocation_id", ["invocationId"]),

  // Agents table
  agents: defineTable({
    name: v.string(),
    role: v.union(
      v.literal("supervisor"),
      v.literal("project-manager"),
      v.literal("coder"),
      v.literal("researcher"),
      v.literal("qa-engineer"),
      v.literal("devops"),
      v.literal("designer"),
      v.literal("human") // Add human role for existing data
    ),
    description: v.string(),
    status: v.optional(
      v.union(
        v.literal("online"),
        v.literal("offline"),
        v.literal("busy"),
        v.literal("idle")
      )
    ),
    isStreaming: v.optional(v.boolean()),
    capabilities: v.optional(v.array(v.string())),
    currentTasks: v.optional(v.array(v.string())), // Array of task UUIDs
    lastActiveAt: v.optional(v.number()), // Unix timestamp
    id: v.string(),
  })
    .index("by_role", ["role"])
    .index("by_status", ["status"])
    .index("by_last_active", ["lastActiveAt"]),

  // Real-time events for live updates
  events: defineTable(EventsSchema)
    .index("by_type", ["type"])
    .index("by_timestamp", ["timestamp"])
    .index("by_entity", ["entityId"]),

  // Agent Projects table - for team configurations and canvas layouts
  agentProjects: defineTable(AgentProjectsSchema)
    .index("by_name", ["name"])
    .index("by_updated", ["updatedAt"]),

  // User settings table
  settings: defineTable(SettingsSchema).index("by_user", ["userId"]),
});
