/**
 * Agent Store Operations
 * Implements StoreOperations interface for agent-related functionality
 */

import { StoreOperations } from '../store-registry';
import { useAgentProjectStore } from '@/stores';

export class AgentStoreOperations implements StoreOperations {
  getCount(): number {
    try {
      const agentProjectStore = useAgentProjectStore();
      const { currentAgentProject } = agentProjectStore;
      return currentAgentProject?.agentNodes.length || 0;
    } catch (error) {
      console.warn('Error getting agent count:', error);
      return 0;
    }
  }
  
  async updatePosition(id: string, position: { x: number; y: number }): Promise<void> {
    try {
      const agentProjectStore = useAgentProjectStore();
      const { updateAgentNodePosition, currentAgentProject } = agentProjectStore;
      if (currentAgentProject && updateAgentNodePosition) {
        updateAgentNodePosition(currentAgentProject.id, id, position);
      }
    } catch (error) {
      console.error('Error updating agent position:', error);
      throw error;
    }
  }
  
  getDisplayName(): string {
    return 'Agents';
  }
}