/**
 * Projects Workspace Component
 * Displays and manages project navigation
 */

import React from "react";
import { Button } from "@/components/ui/button";
import { Project } from "@/types/project.types";
import { ChevronRight, Plus } from "lucide-react";

interface ProjectsWorkspaceProps {
  projects: Project[];
  setCurrentProject: (project: Project) => void;
  currentProject: Project | null;
}

export const ProjectsWorkspace: React.FC<ProjectsWorkspaceProps> = ({
  projects,
  setCurrentProject,
  currentProject,
}) => (
  <div className="space-y-2">
    <div className="flex items-center justify-between">
      <h3 className="text-sm font-medium">Recent Projects</h3>
      <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
        <Plus className="h-3 w-3" />
      </Button>
    </div>
    <div className="space-y-1">
      {projects.map((project) => (
        <Button
          key={project.id}
          variant={currentProject?.id === project.id ? "secondary" : "ghost"}
          className="w-full justify-between h-8 px-2"
          onClick={() => {
            // eslint-disable-next-line no-console
            console.log("Setting current project:", project);
            setCurrentProject(project);
          }}
        >
          <span className="text-xs truncate">{project.title}</span>
          <ChevronRight className="h-3 w-3 shrink-0" />
        </Button>
      ))}

      {projects.length === 0 && (
        <div className="text-xs text-muted-foreground p-2">
          No projects found
        </div>
      )}
    </div>
  </div>
);