"use client";

import { useWorkspaceStore } from "@/lib/store";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
  MessageSquare,
  Settings,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { cn } from "@/lib/utils";

export function WorkspaceSidebar() {
  const {
    activeWorkspace,
    setActiveWorkspace,
    sidebarOpen,
    sidebarCollapsed,
    toggleSidebarCollapsed,
  } = useWorkspaceStore();

  if (!sidebarOpen) {
    return null;
  }

  return (
    <aside
      className={cn(
        "h-full border-r bg-muted/40 transition-all duration-300 relative",
        sidebarCollapsed ? "w-16" : "w-full"
      )}
    >
      <div className="flex h-full max-h-screen flex-col gap-2">
        <div className="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6 pr-4 relative">
          <div
            className={cn(
              "flex items-center gap-2 font-semibold transition-opacity",
              sidebarCollapsed ? "justify-center" : ""
            )}
          >
            <MessageSquare className="h-6 w-6" />
            {!sidebarCollapsed && <span>Workspaces</span>}
          </div>
          <Button
            variant="ghost"
            size="sm"
            className="absolute -right-3 top-4 h-6 w-6 rounded-full border border-white bg-muted/60 p-0 shadow-md hover:bg-accent"
            onClick={toggleSidebarCollapsed}
          >
            {sidebarCollapsed ? (
              <ChevronRight className="h-3 w-3" />
            ) : (
              <ChevronLeft className="h-3 w-3" />
            )}
          </Button>
        </div>
        <div className="flex-1">
          <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
            <Button
              variant={activeWorkspace === "chat" ? "default" : "ghost"}
              className={cn(
                "w-full justify-start gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary",
                sidebarCollapsed ? "justify-center px-2" : "",
                activeWorkspace === "chat"
                  ? "bg-muted text-primary"
                  : "text-muted-foreground"
              )}
              onClick={() => setActiveWorkspace("chat")}
            >
              <MessageSquare className="h-4 w-4" />
              {!sidebarCollapsed && "Chat"}
            </Button>

            <Separator className="my-2" />

            <Button
              variant={activeWorkspace === "settings" ? "default" : "ghost"}
              className={cn(
                "w-full justify-start gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary",
                sidebarCollapsed ? "justify-center px-2" : "",
                activeWorkspace === "settings"
                  ? "bg-muted text-primary"
                  : "text-muted-foreground"
              )}
              onClick={() => setActiveWorkspace("settings")}
            >
              <Settings className="h-4 w-4" />
              {!sidebarCollapsed && "Settings"}
            </Button>
          </nav>
        </div>
      </div>
    </aside>
  );
}
