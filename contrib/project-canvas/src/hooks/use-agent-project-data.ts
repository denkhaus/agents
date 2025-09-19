/**
 * Agent Project Data Hook
 * Loads agents from Convex and creates agent project visualization
 */

import { useEffect, useCallback, useRef } from "react";
import { useAgentProjectStore, useAgentStore } from "@/stores";
import { useConvexAgents, useConvexAgentProjects } from "./use-convex-data";
import { AgentProject, AgentNode } from "@/types/agent.types";

// Create agent nodes from Convex agents
const createAgentNodes = (agents: any[]): AgentNode[] => {
  const positions = [
    { x: 100, y: 100 },   // Designer
    { x: 300, y: 200 },   // Frontend Dev  
    { x: 500, y: 100 },   // Backend Dev
    { x: 300, y: 350 },   // QA Engineer
    { x: 500, y: 300 },   // DevOps
  ];

  return agents.map((agent, index) => ({
    id: agent.id, // Use agent ID directly
    type: "agent" as const,
    position: positions[index] || { x: 100 + (index * 200), y: 100 },
    data: { agent },
  }));
};

// Create agent project from Convex agents
const createAgentProject = (agents: any[]): AgentProject => {
  if (agents.length === 0) {
    return {
      id: "project-1",
      name: "Development Team",
      description: "Main development team with all core roles",
      agents: [],
      agentNodes: [],
      connections: [],
      createdAt: new Date(),
      updatedAt: new Date(),
    };
  }

  return {
    id: "project-1",
    name: "Development Team", 
    description: "Main development team with all core roles",
    agents: agents,
    agentNodes: createAgentNodes(agents),
    connections: [
      // Create connections based on agent roles
      ...agents.slice(0, -1).map((agent, index) => ({
        id: `conn-${index + 1}`,
        source: agent.id,
        target: agents[index + 1].id,
        type: "collaboration" as const,
        label: getConnectionLabel(agent.role, agents[index + 1].role),
      })),
    ],
    createdAt: new Date("2024-01-15"),
    updatedAt: new Date("2024-01-20"),
  };
};

// Helper function to get connection labels based on agent roles
const getConnectionLabel = (sourceRole: string, targetRole: string): string => {
  const roleMap: Record<string, Record<string, string>> = {
    designer: {
      coder: "Design Handoff",
      "frontend-dev": "Design Handoff",
    },
    coder: {
      "qa-engineer": "Code Review",
      devops: "Deployment",
    },
    "qa-engineer": {
      devops: "Testing Results",
    },
  };

  return roleMap[sourceRole]?.[targetRole] || "Collaboration";
};

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