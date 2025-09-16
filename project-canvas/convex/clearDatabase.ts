/**
 * Clear Database Mutation
 * Löscht alle Daten aus der Datenbank für einen sauberen Neustart
 */

import { mutation } from "./_generated/server";

export const clearDatabase = mutation({
  args: {},
  handler: async (ctx) => {
    // Lösche alle Tasks
    const tasks = await ctx.db.query("tasks").collect();
    for (const task of tasks) {
      await ctx.db.delete(task._id);
    }

    // Lösche alle Agents
    const agents = await ctx.db.query("agents").collect();
    for (const agent of agents) {
      await ctx.db.delete(agent._id);
    }

    // Lösche alle Projects
    const projects = await ctx.db.query("projects").collect();
    for (const project of projects) {
      await ctx.db.delete(project._id);
    }

    // Lösche alle Events
    const events = await ctx.db.query("events").collect();
    for (const event of events) {
      await ctx.db.delete(event._id);
    }

    return {
      message: "Database cleared successfully",
      deletedTasks: tasks.length,
      deletedAgents: agents.length,
      deletedProjects: projects.length,
      deletedEvents: events.length,
    };
  },
});