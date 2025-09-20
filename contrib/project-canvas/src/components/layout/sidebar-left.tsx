/**
 * Collapsible Sidebar Component
 * Multi-workspace navigation (Projects, Agents, Settings)
 */

import React from "react";
import { Link, useLocation } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { useUIStore, useProjectStore, useAgentStore } from "@/stores"; // Added useAgentStore
import { WorkspaceType } from "@/types";
import { FolderOpen, Users, Settings, Plus } from "lucide-react";
import { cn } from "@/lib/utils";

const getWorkspaces = (
  projectCount: number,
  agentCount: number
): Array<{
  id: WorkspaceType;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  description: string;
  count?: number;
  path: string; // Added path for routing
}> => [
  {
    id: "projects",
    label: "Projects",
    icon: FolderOpen,
    description: "Manage and visualize projects",
    count: projectCount,
    path: "/projects",
  },
  {
    id: "agents",
    label: "Agents",
    icon: Users,
    description: "View agent status and assignments",
    count: agentCount,
    path: "/agents",
  },
  {
    id: "settings",
    label: "Settings",
    icon: Settings,
    description: "Application preferences",
    path: "/settings",
  },
];

export const SidebarLeft: React.FC = () => {
  const { sidebarCollapsed } = useUIStore();
  const { setCurrentProject, projects } = useProjectStore();
  const { agents } = useAgentStore();
  const location = useLocation();

  const workspaces = getWorkspaces(projects.length, agents.length);

  // Auto-select first project when navigating to Projects workspace
  React.useEffect(() => {
    if (location.pathname === "/projects" && projects.length > 0) {
      setCurrentProject(projects[0]);
    }
  }, [location.pathname, projects, setCurrentProject]);

  return (
    <aside
      className={cn(
        "fixed left-0 top-14 h-[calc(100vh-3.5rem)] bg-background border-r border-border transition-all duration-300 ease-in-out z-40",
        sidebarCollapsed ? "w-16" : "w-64"
      )}
    >
      <div className="flex h-full flex-col">
        {/* Workspace Navigation */}
        <div className="p-2">
          <div className="space-y-1">
            {workspaces.map((workspace) => {
              const Icon = workspace.icon;
              const isActive = location.pathname === workspace.path;

              return (
                <Button
                  key={workspace.id}
                  variant={isActive ? "secondary" : "ghost"}
                  className={cn(
                    "w-full justify-start h-10",
                    sidebarCollapsed ? "px-2" : "px-3"
                  )}
                  asChild
                >
                  <Link to={workspace.path}>
                    <Icon
                      className={cn(
                        "h-4 w-4 shrink-0",
                        sidebarCollapsed ? "mx-auto" : "mr-3"
                      )}
                    />

                    {!sidebarCollapsed && (
                      <>
                        <span className="flex-1 text-left">
                          {workspace.label}
                        </span>
                        {workspace.count && (
                          <Badge variant="secondary" className="ml-auto">
                            {workspace.count}
                          </Badge>
                        )}
                      </>
                    )}
                  </Link>
                </Button>
              );
            })}
          </div>
        </div>

        <Separator />

        {/* Workspace Content - no longer rendered here directly */}
        <div className="flex-1 overflow-hidden">
          {!sidebarCollapsed && (
            <ScrollArea className="h-full">
              <div className="p-2 max-w-full overflow-hidden">
                {/* Content will be rendered by router */}
              </div>
            </ScrollArea>
          )}
        </div>

        {/* Bottom Actions */}
        {!sidebarCollapsed && (
          <>
            <Separator />
            <div className="p-2">
              <Button variant="outline" className="w-full" size="sm">
                <Plus className="h-4 w-4 mr-2" />
                Quick Action
              </Button>
            </div>
          </>
        )}
      </div>
    </aside>
  );
};
