/**
 * Convex Settings Functions
 * Manage user preferences and application settings
 */

import { query, mutation } from "./_generated/server";
import { v } from "convex/values";

export const SCHEMA = {
  userId: v.string(), // Unique identifier for the user
  theme: v.optional(v.union(v.literal("light"), v.literal("dark"))),
  // Application settings
  notifications: v.optional(v.boolean()),
  autoSave: v.optional(v.boolean()),
  language: v.optional(v.string()),
  // UI state persistence
  leftSidebarCollapsed: v.optional(v.boolean()),
  rightSidebarCollapsed: v.optional(v.boolean()),
  currentWorkspace: v.optional(v.string()),
  selectedProjectId: v.optional(v.string()),
  selectedNodeIds: v.optional(v.array(v.string())),
  // Canvas settings
  showMiniMap: v.optional(v.boolean()),
  showBackground: v.optional(v.boolean()),
  autoLayout: v.optional(v.boolean()),
  projectCanvasViewport: v.optional(
    v.object({ x: v.number(), y: v.number(), zoom: v.number() })
  ),
  agentsCanvasViewport: v.optional(
    v.object({ x: v.number(), y: v.number(), zoom: v.number() })
  ),
  createdAt: v.number(), // Unix timestamp
  updatedAt: v.number(), // Unix timestamp
};

// Get user settings
export const getSettings = query({
  args: {
    userId: v.string(),
  },
  handler: async (ctx, args) => {
    const settings = await ctx.db
      .query("settings")
      .withIndex("by_user", (q) => q.eq("userId", args.userId))
      .unique();

    return settings;
  },
});

// Create or update user settings
export const updateSettings = mutation({
  args: {
    userId: v.string(),
    theme: v.optional(v.union(v.literal("light"), v.literal("dark"))),
    notifications: v.optional(v.boolean()),
    autoSave: v.optional(v.boolean()),
    language: v.optional(v.string()),
    leftSidebarCollapsed: v.optional(v.boolean()),
    rightSidebarCollapsed: v.optional(v.boolean()),
    currentWorkspace: v.optional(v.string()),
    selectedProjectId: v.optional(v.string()),
    selectedNodeIds: v.optional(v.array(v.string())),
    showMiniMap: v.optional(v.boolean()),
    showBackground: v.optional(v.boolean()),
    autoLayout: v.optional(v.boolean()),
    projectCanvasViewport: v.optional(
      v.object({ x: v.number(), y: v.number(), zoom: v.number() })
    ),
    agentsCanvasViewport: v.optional(
      v.object({ x: v.number(), y: v.number(), zoom: v.number() })
    ),
  },
  handler: async (ctx, args) => {
    // Check if settings already exist for this user
    const existingSettings = await ctx.db
      .query("settings")
      .withIndex("by_user", (q) => q.eq("userId", args.userId))
      .unique();

    const timestamp = Date.now();

    if (existingSettings) {
      // Update existing settings
      await ctx.db.patch(existingSettings._id, {
        ...(args.theme !== undefined && { theme: args.theme }),
        ...(args.notifications !== undefined && {
          notifications: args.notifications,
        }),
        ...(args.autoSave !== undefined && { autoSave: args.autoSave }),
        ...(args.language !== undefined && { language: args.language }),
        ...(args.leftSidebarCollapsed !== undefined && {
          leftSidebarCollapsed: args.leftSidebarCollapsed,
        }),
        ...(args.rightSidebarCollapsed !== undefined && {
          rightSidebarCollapsed: args.rightSidebarCollapsed,
        }),
        ...(args.currentWorkspace !== undefined && {
          currentWorkspace: args.currentWorkspace,
        }),
        ...(args.selectedProjectId !== undefined && {
          selectedProjectId: args.selectedProjectId,
        }),
        ...(args.selectedNodeIds !== undefined && {
          selectedNodeIds: args.selectedNodeIds,
        }),
        ...(args.showMiniMap !== undefined && {
          showMiniMap: args.showMiniMap,
        }),
        ...(args.showBackground !== undefined && {
          showBackground: args.showBackground,
        }),
        ...(args.autoLayout !== undefined && { autoLayout: args.autoLayout }),
        ...(args.projectCanvasViewport !== undefined && {
          projectCanvasViewport: args.projectCanvasViewport,
        }),
        ...(args.agentsCanvasViewport !== undefined && {
          agentsCanvasViewport: args.agentsCanvasViewport,
        }),
        updatedAt: timestamp,
      });
      return existingSettings._id;
    } else {
      // Create new settings
      const newSettingsId = await ctx.db.insert("settings", {
        userId: args.userId,
        ...(args.theme !== undefined && { theme: args.theme }),
        ...(args.notifications !== undefined && {
          notifications: args.notifications,
        }),
        ...(args.autoSave !== undefined && { autoSave: args.autoSave }),
        ...(args.language !== undefined && { language: args.language }),
        ...(args.leftSidebarCollapsed !== undefined && {
          leftSidebarCollapsed: args.leftSidebarCollapsed,
        }),
        ...(args.rightSidebarCollapsed !== undefined && {
          rightSidebarCollapsed: args.rightSidebarCollapsed,
        }),
        ...(args.currentWorkspace !== undefined && {
          currentWorkspace: args.currentWorkspace,
        }),
        ...(args.selectedProjectId !== undefined && {
          selectedProjectId: args.selectedProjectId,
        }),
        ...(args.selectedNodeIds !== undefined && {
          selectedNodeIds: args.selectedNodeIds,
        }),
        ...(args.showMiniMap !== undefined && {
          showMiniMap: args.showMiniMap,
        }),
        ...(args.showBackground !== undefined && {
          showBackground: args.showBackground,
        }),
        ...(args.autoLayout !== undefined && { autoLayout: args.autoLayout }),
        ...(args.projectCanvasViewport !== undefined && {
          projectCanvasViewport: args.projectCanvasViewport,
        }),
        ...(args.agentsCanvasViewport !== undefined && {
          agentsCanvasViewport: args.agentsCanvasViewport,
        }),
        createdAt: timestamp,
        updatedAt: timestamp,
      });
      return newSettingsId;
    }
  },
});

// Update only the theme setting
export const updateTheme = mutation({
  args: {
    userId: v.string(),
    theme: v.union(v.literal("light"), v.literal("dark")),
  },
  handler: async (ctx, args) => {
    // Check if settings already exist for this user
    const existingSettings = await ctx.db
      .query("settings")
      .withIndex("by_user", (q) => q.eq("userId", args.userId))
      .unique();

    const timestamp = Date.now();

    if (existingSettings) {
      // Update existing settings
      await ctx.db.patch(existingSettings._id, {
        theme: args.theme,
        updatedAt: timestamp,
      });
      return existingSettings._id;
    } else {
      // Create new settings with only theme
      const newSettingsId = await ctx.db.insert("settings", {
        userId: args.userId,
        theme: args.theme,
        createdAt: timestamp,
        updatedAt: timestamp,
      });
      return newSettingsId;
    }
  },
});
