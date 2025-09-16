/**
 * Convex Settings Functions
 * Manage user preferences and application settings
 */

import { query, mutation } from "./_generated/server";
import { v } from "convex/values";

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
        ...args,
        updatedAt: timestamp,
      });
      return existingSettings._id;
    } else {
      // Create new settings
      const newSettingsId = await ctx.db.insert("settings", {
        userId: args.userId,
        theme: args.theme,
        notifications: args.notifications,
        autoSave: args.autoSave,
        language: args.language,
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