/**
 * Convex Data Hooks
 * Custom hooks for real-time data fetching with Convex
 */

import React from "react";
import { useQuery, useMutation } from "convex/react";
import { api } from "../../convex/_generated/api";
import {
  useProjectStore,
  useTaskStore,
  useAgentStore,
  useSettingsStore,
} from "@/stores";
import {
  convexProjectToProject,
  convexTaskToTask,
  convexAgentToAgent,
} from "@/utils/convex-helpers";

// Hook for real-time projects
export const useConvexProjects = () => {
  const convexProjects = useQuery(api.projects.list);
  const { setProjects, setCurrentProject, currentProject } = useProjectStore();
  const { selectedProjectId } = useSettingsStore(); // Get selected project ID from settings

  React.useEffect(() => {
    if (convexProjects) {
      const projects = convexProjects.map(convexProjectToProject);
      setProjects(projects);

      if (projects.length > 0) {
        let projectToSet = null;

        // 1. Try to set current project from selectedProjectId in settings
        if (selectedProjectId) {
          projectToSet = projects.find((p) => p.id === selectedProjectId);
        }

        // 2. If no project found from settings or no selectedProjectId, default to first project
        if (!projectToSet && projects.length > 0) {
          projectToSet = projects[0];
        }

        // 3. Only update if the project to set is different from the current one
        if (projectToSet && projectToSet.id !== currentProject?.id) {
          setCurrentProject(projectToSet);
        } else if (!projectToSet && currentProject) {
          // If no projects available but one was selected, clear the selection
          setCurrentProject(null);
        }
      }
    }
  }, [
    convexProjects,
    setProjects, // State setters are stable, but including them is fine
    setCurrentProject, // State setters are stable, but including them is fine
    selectedProjectId,
  ]);

  return {
    projects: convexProjects?.map(convexProjectToProject) || [],
    loading: convexProjects === undefined,
  };
};

// Hook for real-time tasks
export const useConvexTasks = (projectId?: string) => {
  const convexTasks = useQuery(
    api.tasks.listByProject,
    projectId ? { projectId } : "skip" // No more casting needed, projectId is now string
  );
  const { setTasks } = useTaskStore();

  React.useEffect(() => {
    if (convexTasks) {
      const tasks = convexTasks.map(convexTaskToTask);
      setTasks(tasks);
    }
  }, [convexTasks, setTasks]);

  return {
    tasks: convexTasks?.map(convexTaskToTask) || [],
    loading: convexTasks === undefined,
  };
};

// Hook for real-time agents
export const useConvexAgents = () => {
  const convexAgents = useQuery(api.agents.list);
  const { setAgents } = useAgentStore();

  React.useEffect(() => {
    if (convexAgents) {
      const agents = convexAgents.map(convexAgentToAgent);
      setAgents(agents);
    }
  }, [convexAgents, setAgents]);

  return {
    agents: convexAgents?.map(convexAgentToAgent) || [],
    loading: convexAgents === undefined,
  };
};

// Hook for real-time agent projects
export const useConvexAgentProjects = () => {
  const convexAgentProjects = useQuery(api.agentProjects.list);

  return {
    agentProjects: convexAgentProjects || [],
    loading: convexAgentProjects === undefined,
  };
};

// Hook for real-time events
export const useConvexEvents = (projectId?: string, since?: number) => {
  const events = useQuery(
    api.events.subscribeToProject,
    projectId
      ? {
          projectId,
          since: since || Date.now() - 24 * 60 * 60 * 1000, // Last 24 hours
        }
      : "skip"
  );

  return {
    events: events || [],
    loading: events === undefined,
  };
};

// Mutation hooks for updates
export const useConvexMutations = () => {
  const updateProjectFields = useMutation(api.projects.updateEditableFields);
  const updateTaskFields = useMutation(api.tasks.updateEditableFields);
  const updateTaskPosition = useMutation(api.tasks.updatePosition);
  const updateProjectPosition = useMutation(api.projects.updatePosition);
  const updateAgentNodes = useMutation(api.agentProjects.updateAgentNodes); // Add updateAgentNodes mutation

  return {
    updateProjectFields,
    updateTaskFields,
    updateTaskPosition,
    updateProjectPosition,
    updateAgentNodes, // Expose the new mutation
  };
};
