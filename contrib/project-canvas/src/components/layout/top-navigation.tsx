/**
 * Top Navigation Bar
 * Contains app title, theme toggle, and global actions
 */

import React from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "./theme-toggle";
import { useUIStore } from "@/stores";
import { Menu, X, Zap, Settings, HelpCircle, PanelRight } from "lucide-react";

export const TopNavigation: React.FC = () => {
  const {
    sidebarCollapsed,
    toggleLeftSidebar,
    rightSidebarCollapsed,
    toggleRightSidebar,
  } = useUIStore();
  const navigate = useNavigate();

  return (
    <header className="h-14 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-full items-center justify-between px-4">
        {/* Left Section */}
        <div className="flex items-center gap-4">
          {/* Sidebar Toggle */}
          <Button
            variant="ghost"
            size="sm"
            onClick={toggleLeftSidebar}
            className="h-8 w-8 p-0"
          >
            {sidebarCollapsed ? (
              <Menu className="h-4 w-4" />
            ) : (
              <X className="h-4 w-4" />
            )}
          </Button>

          {/* App Title */}
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary">
              <Zap className="h-4 w-4 text-primary-foreground" />
            </div>
            <div className="flex flex-col">
              <h1 className="text-sm font-semibold leading-none">
                Project Canvas
              </h1>
              <p className="text-xs text-muted-foreground">
                Real-time Visualization
              </p>
            </div>
          </div>
        </div>

        {/* Right Section */}
        <div className="flex items-center gap-2">
          {/* Right Sidebar Toggle */}
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={toggleRightSidebar}
            title={
              rightSidebarCollapsed ? "Show Properties" : "Hide Properties"
            }
          >
            <PanelRight className="h-4 w-4" />
          </Button>

          {/* Help Button */}
          <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
            <HelpCircle className="h-4 w-4" />
          </Button>

          {/* Settings Button */}
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={() => navigate("/settings")}
          >
            <Settings className="h-4 w-4" />
          </Button>

          {/* Theme Toggle */}
          <ThemeToggle />
        </div>
      </div>
    </header>
  );
};
