/**
 * Main Application Layout
 * Provides the overall structure with sidebar and main content area
 */

import React from "react";
import { SidebarLeft } from "./sidebar-left";
import { SidebarRight } from "./sidebar-right";
import { TopNavigation } from "./top-navigation";
import { useUIStore } from "@/stores";
import { cn } from "@/lib/utils";

interface AppLayoutProps {
  children: React.ReactNode;
}

export const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const { sidebarCollapsed, rightSidebarCollapsed } = useUIStore();

  return (
    <div className="h-screen flex flex-col bg-background">
      {/* Top Navigation */}
      <TopNavigation />

      {/* Main Content Area */}
      <div className="flex flex-1 overflow-hidden">
        <SidebarLeft />
        {/* Main Content */}
        <main
          className={cn(
            "flex-1 transition-all duration-300 ease-in-out",
            "bg-background",
            sidebarCollapsed ? "ml-16" : "ml-64",
            rightSidebarCollapsed ? "mr-0" : "mr-80"
          )}
        >
          {children}
        </main>
        <SidebarRight />
      </div>
    </div>
  );
};
