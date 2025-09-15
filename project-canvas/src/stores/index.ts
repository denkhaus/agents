/**
 * Zustand Stores Index
 * Central export for all application stores
 */

export { useProjectStore } from './projectStore';
export { useTaskStore } from './taskStore';
export { useUIStore } from './uiStore';
export { useAgentStore } from './agentStore';

// Re-export store types for convenience
export type { ProjectStore } from './projectStore';
export type { TaskStore } from './taskStore';
export type { UIStore } from './uiStore';
export type { AgentStore } from './agentStore';

// Store initialization helper
export const initializeStores = () => {
  // Initialize stores with default values if needed
  // This can be called on app startup
  console.log('Stores initialized');
};

// Store reset helper for testing
export const resetAllStores = () => {
  useProjectStore.getState().reset();
  useTaskStore.getState().reset();
  useUIStore.getState().reset();
  useAgentStore.getState().reset();
};