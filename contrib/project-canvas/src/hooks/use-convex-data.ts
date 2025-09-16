/**
 * Convex Data Hooks
 * Custom hooks for real-time data fetching with Convex
 */

import React from "react";
import { useQuery, useMutation } from "convex/react";
import { api } from "../../convex/_generated/api";
import { useProjectStore, useTaskStore, useAgentStore } from "@/stores";
import {
  convexProjectToProject,
  convexTaskToTask,
  convexAgentToAgent,
} from "@/utils/convex-helpers";

// Hook for real-time projects
export const useConvexProjects = () => {
  const convexProjects = useQuery(api.projects.list);
  const { setProjects, setCurrentProject, currentProject } = useProjectStore();

  React.useEffect(() => {
    if (convexProjects) {
      const projects = convexProjects.map(convexProjectToProject);
      setProjects(projects);

      // Set first project as current if none selected, and only if it's not already set.
      // This prevents a re-render loop.
      if (!useProjectStore.getState().currentProject && projects.length > 0) {
        setCurrentProject(projects[0]);
      }
    }
  }, [convexProjects, setProjects, setCurrentProject]);

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

  // eslint-disable-next-line react-hooks/exhaustive-deps
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

  // eslint-disable-next-line react-hooks/exhaustive-deps
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

  return {
    updateProjectFields,
    updateTaskFields,
    updateTaskPosition,
  };
};
