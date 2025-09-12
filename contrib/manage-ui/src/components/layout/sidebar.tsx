'use client';

import { cn } from '@/lib/utils';
import { useUIStore } from '@/lib/stores/ui-store';
import { 
  FolderKanban, 
  Users, 
  Activity, 
  Settings, 
  ChevronLeft,
  ChevronRight
} from 'lucide-react';
import { Button } from '@/components/ui/button';

const workspaces = [
  {
    id: 'projects' as const,
    title: 'Projects',
    icon: FolderKanban,
    description: 'Manage projects and tasks',
    enabled: true,
  },
  {
    id: 'agents' as const,
    title: 'Agents',
    icon: Users,
    description: 'Agent management',
    enabled: false, // Future feature
  },
  {
    id: 'monitoring' as const,
    title: 'Monitoring',
    icon: Activity,
    description: 'System monitoring',
    enabled: false, // Future feature
  },
  {
    id: 'settings' as const,
    title: 'Settings',
    icon: Settings,
    description: 'Application settings',
    enabled: false, // Future feature
  },
];

export function Sidebar() {
  const { 
    currentWorkspace, 
    sidebarCollapsed, 
    setCurrentWorkspace, 
    setSidebarCollapsed 
  } = useUIStore();

  return (
    <div className={cn(
      "flex flex-col bg-sidebar border-r border-sidebar-border transition-all duration-300",
      sidebarCollapsed ? "w-16" : "w-64"
    )}>
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-sidebar-border">
        {!sidebarCollapsed && (
          <h1 className="text-lg font-semibold text-sidebar-foreground">
            Admin Panel
          </h1>
        )}
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
          className="h-8 w-8 p-0"
        >
          {sidebarCollapsed ? (
            <ChevronRight className="h-4 w-4" />
          ) : (
            <ChevronLeft className="h-4 w-4" />
          )}
        </Button>
      </div>

      {/* Navigation */}
      <nav className="flex-1 p-2">
        <ul className="space-y-1">
          {workspaces.map((workspace) => {
            const Icon = workspace.icon;
            const isActive = currentWorkspace === workspace.id;
            const isDisabled = !workspace.enabled;

            return (
              <li key={workspace.id}>
                <Button
                  variant={isActive ? "secondary" : "ghost"}
                  className={cn(
                    "w-full justify-start h-10",
                    sidebarCollapsed ? "px-2" : "px-3",
                    isDisabled && "opacity-50 cursor-not-allowed"
                  )}
                  onClick={() => {
                    if (workspace.enabled) {
                      setCurrentWorkspace(workspace.id);
                    }
                  }}
                  disabled={isDisabled}
                  title={sidebarCollapsed ? workspace.title : undefined}
                >
                  <Icon className={cn(
                    "h-4 w-4",
                    !sidebarCollapsed && "mr-3"
                  )} />
                  {!sidebarCollapsed && (
                    <span className="truncate">{workspace.title}</span>
                  )}
                </Button>
              </li>
            );
          })}
        </ul>
      </nav>

      {/* Footer */}
      {!sidebarCollapsed && (
        <div className="p-4 border-t border-sidebar-border">
          <div className="text-xs text-muted-foreground">
            Multi-Agent System
            <br />
            Admin Frontend v1.0
          </div>
        </div>
      )}
    </div>
  );
}