/**
 * Agent Project Data Hook
 * Loads agents from Convex and creates agent project visualization
 */

import { useEffect, useRef } from "react";
import { useAgentProjectStore, useAgentStore } from "@/stores";
import { useConvexAgents, useConvexAgentProjects } from "./use-convex-data";


export const useAgentProjectData = () => {
  const { agentProjects, setAgentProjects, setCurrentAgentProject, currentAgentProject } = useAgentProjectStore();
  const { setAgents } = useAgentStore();
  
  // Load agents and agent projects from Convex
  const { agents: convexAgents, loading: agentsLoading } = useConvexAgents();
  const { agentProjects: convexAgentProjects, loading: agentProjectsLoading } = useConvexAgentProjects();
  
  // Track if we've already initialized to prevent loops
  const initializedRef = useRef(false);

  useEffect(() => {
    const loading = agentsLoading || agentProjectsLoading;
    
    if (!loading) {
      if (!initializedRef.current) {
        // Update agent store with Convex data
        if (convexAgents.length > 0) {
          setAgents(convexAgents);
        }
        
        // Load agent projects from Convex
        if (convexAgentProjects.length > 0) {
          setAgentProjects(convexAgentProjects);
          
          // Set the first agent project as current if none is selected
          if (!currentAgentProject) {
            setCurrentAgentProject(convexAgentProjects[0]);
          }
        } else {
          // Fallback to dummy data if no Convex agent projects
          const { masterAgentProjects } = require('@/data/master-dummy-data');
          
          if (masterAgentProjects && masterAgentProjects.length > 0) {
            setAgentProjects(masterAgentProjects);
            
            if (!currentAgentProject) {
              setCurrentAgentProject(masterAgentProjects[0]);
            }
          }
        }
        
        initializedRef.current = true;
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [agentsLoading, agentProjectsLoading, convexAgents.length, convexAgentProjects.length]); // Store setters are stable and don't need to be in deps

  return {
    agentProjects,
    agents: convexAgents,
    loading: agentsLoading || agentProjectsLoading,
  };
};