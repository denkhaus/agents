/**
 * Clear Database Mutation
 * Efficiently clears all data from the database with optimized batching
 */

import { mutation } from "./_generated/server";

export const clearDatabase = mutation({
  args: {},
  handler: async (ctx) => {
    let totalDeleted = {
      tasks: 0,
      agents: 0,
      projects: 0,
      events: 0,
    };

    // Optimized approach - larger batches but still under limits
    const batchSize = 100;
    const maxReadsPerTable = 2000; // Higher limit for faster clearing
    
    // Helper function to delete in batches
    const deleteBatch = async (tableName: string) => {
      let deletedCount = 0;
      let totalReads = 0;
      
      while (totalReads < maxReadsPerTable) {
        const batch = await ctx.db.query(tableName as any).take(batchSize);
        totalReads += batchSize;
        
        if (batch.length === 0) break;
        
        // Delete all items in parallel for better performance
        await Promise.all(batch.map(item => ctx.db.delete(item._id)));
        deletedCount += batch.length;
        
        // If we got fewer items than requested, we're done with this table
        if (batch.length < batchSize) break;
      }
      
      return deletedCount;
    };

    // Clear events first (likely the largest table)
    totalDeleted.events = await deleteBatch("events");
    
    // Then clear other tables
    totalDeleted.tasks = await deleteBatch("tasks");
    totalDeleted.agents = await deleteBatch("agents");
    totalDeleted.projects = await deleteBatch("projects");

    // Check if there might be more data remaining
    const remainingTasks = await ctx.db.query("tasks").take(1);
    const remainingAgents = await ctx.db.query("agents").take(1);
    const remainingProjects = await ctx.db.query("projects").take(1);
    const remainingEvents = await ctx.db.query("events").take(1);
    
    const hasRemainingData = 
      remainingTasks.length > 0 || 
      remainingAgents.length > 0 || 
      remainingProjects.length > 0 || 
      remainingEvents.length > 0;

    return {
      message: hasRemainingData 
        ? "Database partially cleared - run again to continue" 
        : "Database completely cleared",
      deletedTasks: totalDeleted.tasks,
      deletedAgents: totalDeleted.agents,
      deletedProjects: totalDeleted.projects,
      deletedEvents: totalDeleted.events,
      hasRemainingData,
      totalDeleted: totalDeleted.tasks + totalDeleted.agents + totalDeleted.projects + totalDeleted.events
    };
  },
});