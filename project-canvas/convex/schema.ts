/**
 * Convex Database Schema
 * Defines the structure for projects, tasks, and agents
 */

import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  // Projects table
  projects: defineTable({
    id: v.string(),
    title: v.string(),
    description: v.string(),
    totalTasks: v.number(),
    completedTasks: v.number(),
    progress: v.number(), // 0-100 percentage
    // Convex automatically adds _id and _creationTime
  }).index("by_title", ["title"]),

  // Tasks table
  tasks: defineTable({
    id: v.string(),
    projectId: v.id("projects"),
    parentId: v.optional(v.id("tasks")),
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
    dependencies: v.array(v.id("tasks")), // Array of task IDs
    dependents: v.array(v.id("tasks")), // Array of task IDs
    // UI-specific fields
    positionX: v.optional(v.number()),
    positionY: v.optional(v.number()),
    // Timestamps
    completedAt: v.optional(v.number()), // Unix timestamp
    updatedAt: v.number(), // Unix timestamp
  })
    .index("by_project", ["projectId"])
    .index("by_parent", ["parentId"])
    .index("by_state", ["state"])
    .index("by_assigned_agent", ["assignedAgent"])
    .index("by_updated", ["updatedAt"]),

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
    currentTasks: v.optional(v.array(v.id("tasks"))),
    lastActiveAt: v.optional(v.number()), // Unix timestamp
    id: v.string(),
  })
    .index("by_role", ["role"])
    .index("by_status", ["status"])
    .index("by_last_active", ["lastActiveAt"]),

  // Real-time events for live updates
  events: defineTable({
    type: v.union(
      v.literal("project_updated"),
      v.literal("task_created"),
      v.literal("task_updated"),
      v.literal("task_deleted"),
      v.literal("agent_status_changed"),
      v.literal("task_position_changed")
    ),
    entityId: v.string(), // ID of the affected entity
    data: v.any(), // Event-specific data
    userId: v.optional(v.string()), // User who triggered the event
    timestamp: v.number(),
  })
    .index("by_type", ["type"])
    .index("by_timestamp", ["timestamp"])
    .index("by_entity", ["entityId"]),
});
