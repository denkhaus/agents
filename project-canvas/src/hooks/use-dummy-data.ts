/**
 * Dummy Data Hook
 * Loads dummy data into stores for development
 */

import { useEffect } from 'react';
import { useProjectStore, useTaskStore, useAgentStore } from '@/stores';
import { 
  masterProjects, 
  allTasks, 
  masterAgents 
} from '@/data/master-dummy-data';

export const useDummyData = () => {
  const { setProjects, setCurrentProject } = useProjectStore();
  const { setTasks } = useTaskStore();
  const { setAgents } = useAgentStore();

  useEffect(() => {
    // Load master dummy data into stores
    setProjects(masterProjects);
    setCurrentProject(masterProjects[0]); // E-Commerce project as default
    setTasks(allTasks);
    setAgents(masterAgents);
  }, [setProjects, setCurrentProject, setTasks, setAgents]);

  return {
    projects: masterProjects,
    tasks: allTasks,
    agents: masterAgents
  };
};