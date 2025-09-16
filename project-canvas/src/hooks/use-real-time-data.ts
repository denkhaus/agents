/**
 * Real-time Data Hook
 * Combines Convex real-time data with local stores
 */

import { useProjectStore, useTaskStore } from '@/stores';
import { 
  useConvexProjects, 
  useConvexTasks, 
  useConvexAgents,
  useConvexMutations 
} from './use-convex-data';

export const useRealTimeData = () => {
  const { currentProject } = useProjectStore();
  const { updateTaskPosition: updateTaskPositionStore } = useTaskStore();
  
  // Real-time data hooks
  const { projects, loading: projectsLoading } = useConvexProjects();
  const { tasks, loading: tasksLoading } = useConvexTasks(currentProject?.id);
  const { agents, loading: agentsLoading } = useConvexAgents();
  
  // Mutation hooks
  const { updateTaskPosition } = useConvexMutations();

  // Enhanced task position update with Convex sync
  const handleTaskPositionUpdate = async (taskId: string, position: { x: number; y: number }) => {
    // Optimistic update in local store
    updateTaskPositionStore(taskId, position);
    
    try {
      // Sync to Convex
      await updateTaskPosition({
        id: taskId as any,
        positionX: position.x,
        positionY: position.y,
      });
    } catch (error) {
      console.error('Failed to sync task position:', error);
      // Could implement retry logic or rollback here
    }
  };

  return {
    // Data
    projects,
    tasks,
    agents,
    currentProject,
    
    // Loading states
    loading: projectsLoading || tasksLoading || agentsLoading,
    projectsLoading,
    tasksLoading,
    agentsLoading,
    
    // Actions
    updateTaskPosition: handleTaskPositionUpdate,
  };
};