/**
 * Project Store Operations
 * Implements StoreOperations interface for project-related functionality
 */

import { StoreOperations } from '../store-registry';
import { useProjectStore } from '@/stores';
import { useRealTimeData } from '@/hooks/use-real-time-data';

export class ProjectStoreOperations implements StoreOperations {
  getCount(): number {
    try {
      const projectStore = useProjectStore();
      const { projects } = projectStore;
      return projects.length;
    } catch (error) {
      console.warn('Error getting project count:', error);
      return 0;
    }
  }
  
  async updatePosition(id: string, position: { x: number; y: number }): Promise<void> {
    try {
      const realTimeData = useRealTimeData();
      const { updateProjectPosition } = realTimeData;
      if (updateProjectPosition) {
        await updateProjectPosition(id, position);
      }
    } catch (error) {
      console.error('Error updating project position:', error);
      throw error;
    }
  }
  
  getDisplayName(): string {
    return 'Projects';
  }
}