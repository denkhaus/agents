/**
 * Task Store Operations
 * Implements StoreOperations interface for task-related functionality
 */

import { StoreOperations } from '../store-registry';
import { useTaskStore, useProjectStore } from '@/stores';
import { useRealTimeData } from '@/hooks/use-real-time-data';

export class TaskStoreOperations implements StoreOperations {
  getCount(): number {
    try {
      const taskStore = useTaskStore();
      const projectStore = useProjectStore();
      const { tasksByProject } = taskStore;
      const { currentProject } = projectStore;
      return currentProject ? tasksByProject[currentProject.id]?.length || 0 : 0;
    } catch (error) {
      console.warn('Error getting task count:', error);
      return 0;
    }
  }
  
  async updatePosition(id: string, position: { x: number; y: number }): Promise<void> {
    try {
      const realTimeData = useRealTimeData();
      const { updateTaskPosition } = realTimeData;
      if (updateTaskPosition) {
        await updateTaskPosition(id, position);
      }
    } catch (error) {
      console.error('Error updating task position:', error);
      throw error;
    }
  }
  
  getDisplayName(): string {
    return 'Tasks';
  }
}