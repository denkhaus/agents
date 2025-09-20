/**
 * Workspace Content Switcher
 * Routes to appropriate workspace component based on current workspace
 */

import React from "react";
import { WorkspaceType } from "@/types";
import { useProjectStore } from "@/stores";
import { ProjectsWorkspace } from "./projects-workspace";
import { SettingsWorkspace } from "./settings-workspace";
import { AgentsWorkspace } from "./agents-workspace";

interface WorkspaceContentProps {
  workspace: WorkspaceType;
}

export const WorkspaceContent: React.FC<WorkspaceContentProps> = ({
  workspace,
}) => {
  const { projects, setCurrentProject, currentProject } = useProjectStore(); // Get projects directly from useProjectStore

  switch (workspace) {
    case "projects":
      return (
        <ProjectsWorkspace
          projects={projects}
          setCurrentProject={setCurrentProject}
          currentProject={currentProject}
        />
      );
    case "agents":
      return <AgentsWorkspace />;
    case "settings":
      return <SettingsWorkspace />;
    default:
      return null;
  }
};
