/**
 * Agent Project Data Hook
 * Loads agents from Convex and creates agent project visualization
 */

import { useEffect, useRef } from "react";
import { useAgentProjectStore, useAgentStore } from "@/stores";
import { useConvexAgents, useConvexAgentProjects } from "./use-convex-data";
import { masterAgentProjects } from "@/data/master-dummy-data";
import { convexAgentProjectToAgentProject } from "@/utils/convex-helpers";

export const useAgentProjectData = () => {
  const {
    agentProjects,
    setAgentProjects,
    setCurrentAgentProject,
    currentAgentProject,
  } = useAgentProjectStore();
  const { setAgents } = useAgentStore(); // Get agents from agentStore

  // Load agents and agent projects from Convex
  const { agents: convexAgents, loading: agentsLoading } = useConvexAgents();
  const { agentProjects: convexAgentProjects, loading: agentProjectsLoading } =
    useConvexAgentProjects();

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

        // Load agent projects from Convex, converting them to the frontend type
        if (convexAgentProjects.length > 0 && convexAgents.length > 0) {
          const convertedAgentProjects = convexAgentProjects.map((project) =>
            convexAgentProjectToAgentProject(project, convexAgents)
          );
          
          // Only update if we don't have agent projects yet, or if the current project changed
          const shouldUpdate = agentProjects.length === 0 || 
            !currentAgentProject || 
            !convertedAgentProjects.find(p => p.id === currentAgentProject.id);
            
          if (shouldUpdate) {
            setAgentProjects(convertedAgentProjects);

            // Set the first agent project as current if none is selected
            if (!currentAgentProject && convertedAgentProjects.length > 0) {
              setCurrentAgentProject(convertedAgentProjects[0]);
            }
          } else {
            // Merge positions from Convex with current local state
            // Prioritize local positions to avoid race conditions
            const mergedProjects = agentProjects.map(localProject => {
              const convexProject = convertedAgentProjects.find(p => p.id === localProject.id);
              if (convexProject) {
                const mergedAgentNodes = localProject.agentNodes.map(localNode => {
                  const convexNode = convexProject.agentNodes.find(n => n.id === localNode.id);
                  if (convexNode) {
                    // Always prefer local position if it exists and is not at origin
                    const useLocalPosition = localNode.position && 
                      (localNode.position.x !== 0 || localNode.position.y !== 0);
                    
                    
                    return {
                      ...convexNode,
                      position: useLocalPosition ? localNode.position : convexNode.position
                    };
                  }
                  return localNode;
                });
                
                return {
                  ...convexProject,
                  agentNodes: mergedAgentNodes
                };
              }
              return localProject;
            });
            
            setAgentProjects(mergedProjects);
            
            // Update current project if it's affected
            if (currentAgentProject) {
              const updatedCurrentProject = mergedProjects.find(p => p.id === currentAgentProject.id);
              if (updatedCurrentProject) {
                setCurrentAgentProject(updatedCurrentProject);
              }
            }
          }
        } else if (
          masterAgentProjects &&
          masterAgentProjects.length > 0 &&
          convexAgents.length > 0
        ) {
          // Fallback to dummy data if no Convex agent projects and agents are loaded
          const convertedMasterAgentProjects = masterAgentProjects.map(
            (project) =>
              convexAgentProjectToAgentProject(project as any, convexAgents) // Cast to any because masterAgentProjects might not strictly match ConvexAgentProject
          );
          setAgentProjects(convertedMasterAgentProjects);

          if (!currentAgentProject && convertedMasterAgentProjects.length > 0) {
            setCurrentAgentProject(convertedMasterAgentProjects[0]);
          }
        }

        initializedRef.current = true;
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    agentsLoading,
    agentProjectsLoading,
    convexAgents.length,
    convexAgentProjects.length,
  ]); // Store setters are stable and don't need to be in deps

  return {
    agentProjects,
    agents: convexAgents,
    loading: agentsLoading || agentProjectsLoading,
  };
};
