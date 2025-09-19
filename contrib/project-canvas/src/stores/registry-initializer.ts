/**
 * Registry Initializer
 * Initializes the store registry with all store operations
 */

import { StoreRegistry } from './store-registry';
import { TaskStoreOperations, ProjectStoreOperations, AgentStoreOperations } from './operations';
import { STORE_NAMES } from '@/constants';

// Initialize store registry
let isInitialized = false;

export const initializeStoreRegistry = () => {
  if (isInitialized) {
    return;
  }

  // Register all store operations
  StoreRegistry.register(STORE_NAMES.TASK, new TaskStoreOperations());
  StoreRegistry.register(STORE_NAMES.PROJECT, new ProjectStoreOperations());
  StoreRegistry.register(STORE_NAMES.AGENT, new AgentStoreOperations());

  isInitialized = true;
  console.log('Store registry initialized with const-based store names');
};

// Auto-initialize on import
initializeStoreRegistry();