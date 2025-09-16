/**
 * Collapsible Sidebar Component
 * Multi-workspace navigation (Projects, Agents, Settings)
 */

import React from 'react';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { useUIStore, useProjectStore } from '@/stores';
import { WorkspaceType } from '@/types';
import { useRealTimeData } from '@/hooks/use-real-time-data';
import { Project } from '@/types/project.types';
import { 
  FolderOpen,
  Users,
  Settings,
  ChevronRight,
  Plus
} from 'lucide-react';
import { cn } from '@/lib/utils';

const getWorkspaces = (projectCount: number, agentCount: number): Array<{
  id: WorkspaceType;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  description: string;
  count?: number;
}> => [
  {
    id: 'projects',
    label: 'Projects',
    icon: FolderOpen,
    description: 'Manage and visualize projects',
    count: projectCount
  },
  {
    id: 'agents',
    label: 'Agents',
    icon: Users,
    description: 'View agent status and assignments',
    count: agentCount
  },
  {
    id: 'settings',
    label: 'Settings',
    icon: Settings,
    description: 'Application preferences'
  }
];

export const Sidebar: React.FC = () => {
  const { 
    sidebarCollapsed, 
    currentWorkspace, 
    setWorkspace 
  } = useUIStore();
  
  const { setCurrentProject } = useProjectStore();
  const { projects, agents } = useRealTimeData();
  
  const workspaces = getWorkspaces(projects.length, agents.length);
  
  // Auto-select first project when clicking Projects workspace
  const handleProjectsClick = () => {
    setWorkspace('projects');
    if (projects.length > 0) {
      setCurrentProject(projects[0]);
    }
  };

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
              const isActive = currentWorkspace === workspace.id;
              
              return (
                <Button
                  key={workspace.id}
                  variant={isActive ? "secondary" : "ghost"}
                  className={cn(
                    "w-full justify-start h-10",
                    sidebarCollapsed ? "px-2" : "px-3"
                  )}
                  onClick={() => workspace.id === 'projects' ? handleProjectsClick() : setWorkspace(workspace.id)}
                >
                  <Icon className={cn(
                    "h-4 w-4 shrink-0",
                    sidebarCollapsed ? "mx-auto" : "mr-3"
                  )} />
                  
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
                </Button>
              );
            })}
          </div>
        </div>

        <Separator />

        {/* Workspace Content */}
        <div className="flex-1 overflow-hidden">
          {!sidebarCollapsed && (
            <ScrollArea className="h-full">
              <div className="p-2">
                <WorkspaceContent workspace={currentWorkspace} />
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

// Workspace-specific content
const WorkspaceContent: React.FC<{ workspace: WorkspaceType }> = ({ workspace }) => {
  const { projects } = useRealTimeData();
  const { setCurrentProject, currentProject } = useProjectStore();
  
  switch (workspace) {
    case 'projects':
      return <ProjectsWorkspace projects={projects} setCurrentProject={setCurrentProject} currentProject={currentProject} />;
    case 'agents':
      return <AgentsWorkspace />;
    case 'settings':
      return <SettingsWorkspace />;
    default:
      return null;
  }
};

interface ProjectsWorkspaceProps {
  projects: Project[];
  setCurrentProject: (project: Project) => void;
  currentProject: Project | null;
}

const ProjectsWorkspace: React.FC<ProjectsWorkspaceProps> = ({ projects, setCurrentProject, currentProject }) => (
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
          onClick={() => setCurrentProject(project)}
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

const AgentsWorkspace: React.FC = () => (
  <div className="space-y-2">
    <h3 className="text-sm font-medium">Active Agents</h3>
    <div className="space-y-1">
      {[
        { name: 'Designer', status: 'online' },
        { name: 'Frontend Dev', status: 'busy' },
        { name: 'Backend Dev', status: 'online' },
        { name: 'QA Engineer', status: 'idle' },
        { name: 'DevOps', status: 'online' }
      ].map((agent) => (
        <div key={agent.name} className="flex items-center gap-2 p-2 rounded-md hover:bg-muted">
          <div className={cn(
            "h-2 w-2 rounded-full",
            agent.status === 'online' && "bg-green-500",
            agent.status === 'busy' && "bg-yellow-500",
            agent.status === 'idle' && "bg-gray-400"
          )} />
          <span className="text-xs">{agent.name}</span>
        </div>
      ))}
    </div>
  </div>
);

const SettingsWorkspace: React.FC = () => (
  <div className="space-y-2">
    <h3 className="text-sm font-medium">Preferences</h3>
    <div className="space-y-1">
      {['Display', 'Notifications', 'Keyboard', 'Advanced'].map((setting) => (
        <Button
          key={setting}
          variant="ghost"
          className="w-full justify-between h-8 px-2"
        >
          <span className="text-xs">{setting}</span>
          <ChevronRight className="h-3 w-3" />
        </Button>
      ))}
    </div>
  </div>
);