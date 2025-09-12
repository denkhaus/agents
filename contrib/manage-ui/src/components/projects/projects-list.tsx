"use client";

import { MessageSquare, Calendar, BarChart3 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useProjectStore } from "@/lib/stores/project-store";
import { useUIStore } from "@/lib/stores/ui-store";
import { ProjectStats } from "./project-stats";
import { formatRelativeTime, cn } from "@/lib/utils";
import type { Project } from "@/lib/types";

interface ProjectCardProps {
  project: Project;
  onSelect: (projectId: string) => void;
  onChat: (projectId: string) => void;
}

function ProjectCard({ project, onSelect, onChat }: ProjectCardProps) {
  const getProgressColor = (progress: number) => {
    if (progress >= 80) return "bg-green-500";
    if (progress >= 50) return "bg-blue-500";
    if (progress >= 25) return "bg-yellow-500";
    return "bg-gray-300";
  };

  return (
    <div className="bg-card rounded-lg border border-border p-6 hover:shadow-md transition-shadow cursor-pointer">
      <div className="flex items-start justify-between mb-4">
        <div className="flex-1" onClick={() => onSelect(project.id)}>
          <h3 className="text-lg font-semibold text-card-foreground mb-2">
            {project.title}
          </h3>
          <p className="text-muted-foreground text-sm line-clamp-2">
            {project.description}
          </p>
        </div>

        <Button
          variant="ghost"
          size="sm"
          onClick={(e) => {
            e.stopPropagation();
            onChat(project.id);
          }}
          className="h-8 w-8 p-0 ml-2"
        >
          <MessageSquare className="h-4 w-4" />
        </Button>
      </div>

      <div className="space-y-3">
        {/* Progress bar */}
        <div>
          <div className="flex items-center justify-between text-sm mb-1">
            <span className="text-muted-foreground">Progress</span>
            <span className="font-medium">{Math.round(project.progress)}%</span>
          </div>
          <div className="w-full bg-muted rounded-full h-2">
            <div
              className={cn(
                "h-2 rounded-full transition-all",
                getProgressColor(project.progress)
              )}
              style={{ width: `${project.progress}%` }}
            />
          </div>
        </div>

        {/* Stats */}
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4 text-sm text-muted-foreground">
            <div className="flex items-center">
              <BarChart3 className="h-4 w-4 mr-1" />
              <span>
                {project.completed_tasks}/{project.total_tasks} tasks
              </span>
            </div>
            <div className="flex items-center">
              <Calendar className="h-4 w-4 mr-1" />
              <span>{formatRelativeTime(project.updated_at)}</span>
            </div>
          </div>

          <Badge variant="secondary" className="text-[0.6rem] opacity-70">
            ID: {project.id}
          </Badge>
        </div>
      </div>
    </div>
  );
}

export function ProjectsList() {
  const { projects, tasks, isLoading } = useProjectStore();
  const { setSelectedProject, openChatPanel } = useUIStore();

  const handleSelectProject = (projectId: string) => {
    setSelectedProject(projectId);
  };

  const handleChatProject = (projectId: string) => {
    openChatPanel("project", projectId);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading projects...</p>
        </div>
      </div>
    );
  }

  if (projects.length === 0) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            No projects found
          </h3>
          <p className="text-gray-600 mb-4">
            Get started by creating your first project.
          </p>
          <Button>Create Project</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      {/* Statistics Overview */}
      <ProjectStats projects={projects} allTasks={tasks} />

      {/* Projects Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {projects.map((project) => (
          <ProjectCard
            key={project.id}
            project={project}
            onSelect={handleSelectProject}
            onChat={handleChatProject}
          />
        ))}
      </div>
    </div>
  );
}
