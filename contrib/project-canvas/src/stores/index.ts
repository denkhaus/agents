/**
 * Zustand Stores Index
 * Central export for all application stores
 */

export { useProjectStore } from './projectStore';
export { useTaskStore } from './taskStore';
export { useUIStore } from './uiStore';
export { useAgentStore } from './agentStore';

// Store types are exported from individual store files

// Store initialization helper
export const initializeStores = () => {
  // Initialize stores with default values if needed
  // This can be called on app startup
  console.log('Stores initialized');
};

// Store reset helper for testing
export const resetAllStores = () => {
  // Reset functionality will be implemented in individual stores
  console.log('Reset all stores');
};