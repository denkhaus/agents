/**
 * Agent Project Data Hook
 * Initializes agent project store with sample data
 */

import { useEffect } from "react";
import { useAgentProjectStore, useAgentStore } from "@/stores";
import { masterAgents } from "@/data/agents";
import { AgentProject, AgentNode } from "@/types/agent.types";

// Create agent nodes from real agents
const createAgentNodes = (): AgentNode[] => [
  {
    id: "node-1",
    type: "agent",
    position: { x: 100, y: 100 },
    data: { agent: masterAgents[0] }, // Designer
  },
  {
    id: "node-2", 
    type: "agent",
    position: { x: 300, y: 200 },
    data: { agent: masterAgents[1] }, // Frontend Dev
  },
  {
    id: "node-3",
    type: "agent",
    position: { x: 500, y: 100 },
    data: { agent: masterAgents[2] }, // Backend Dev
  },
  {
    id: "node-4",
    type: "agent",
    position: { x: 300, y: 350 },
    data: { agent: masterAgents[3] }, // QA Engineer
  },
  {
    id: "node-5",
    type: "agent",
    position: { x: 500, y: 300 },
    data: { agent: masterAgents[4] }, // DevOps
  },
];

// Create sample agent project with real agents
const createSampleAgentProject = (): AgentProject => ({
  id: "project-1",
  name: "Development Team",
  description: "Main development team with all core roles",
  agents: masterAgents,
  agentNodes: createAgentNodes(),
  connections: [
    {
      id: "conn-1",
      source: masterAgents[0].id, // Designer
      target: masterAgents[1].id, // Frontend Dev
      type: "collaboration",
      label: "Design Handoff",
    },
    {
      id: "conn-2", 
      source: masterAgents[1].id, // Frontend Dev
      target: masterAgents[2].id, // Backend Dev
      type: "communication",
      label: "API Integration",
    },
    {
      id: "conn-3",
      source: masterAgents[2].id, // Backend Dev
      target: masterAgents[3].id, // QA Engineer
      type: "collaboration",
      label: "Testing",
    },
    {
      id: "conn-4",
      source: masterAgents[3].id, // QA Engineer
      target: masterAgents[4].id, // DevOps
      type: "collaboration", 
      label: "Deployment",
    },
  ],
  createdAt: new Date("2024-01-15"),
  updatedAt: new Date("2024-01-20"),
});

export const useAgentProjectData = () => {
  const { agentProjects, setAgentProjects, setCurrentAgentProject } = useAgentProjectStore();
  const { agents, setAgents } = useAgentStore();

  useEffect(() => {
    // Initialize with real agent data if empty
    if (agentProjects.length === 0) {
      const sampleProject = createSampleAgentProject();
      setAgentProjects([sampleProject]);
      setCurrentAgentProject(sampleProject);
    }
    
    if (agents.length === 0) {
      setAgents(masterAgents);
    }
  }, [agentProjects.length, agents.length, setAgentProjects, setCurrentAgentProject, setAgents]);

  return {
    agentProjects,
    agents,
  };
};