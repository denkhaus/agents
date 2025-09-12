"use client";

import { Search, Bell, User, Sun, Moon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useUIStore } from "@/lib/stores/ui-store";
import { Badge } from "@/components/ui/badge";

export function TopNavigation() {
  const { currentWorkspace, notificationState, themeConfig, setThemeConfig } =
    useUIStore();

  const toggleTheme = () => {
    const newMode = themeConfig.mode === "light" ? "dark" : "light";
    setThemeConfig({ mode: newMode });

    // Apply dark mode to document immediately
    const htmlElement = document.documentElement;
    if (newMode === "dark") {
      htmlElement.classList.add("dark");
    } else {
      htmlElement.classList.remove("dark");
    }

    // Store preference
    localStorage.setItem("theme", newMode);

    console.log(
      "Theme toggled to:",
      newMode,
      "Dark class present:",
      htmlElement.classList.contains("dark")
    );
  };

  const getWorkspaceTitle = () => {
    switch (currentWorkspace) {
      case "projects":
        return "Projects & Tasks";
      case "agents":
        return "Agents Management";
      case "monitoring":
        return "System Monitoring";
      case "settings":
        return "Settings";
      default:
        return "Dashboard";
    }
  };

  return (
    <header className="bg-background border-b border-border px-6 py-4">
      <div className="flex items-center justify-between">
        {/* Left side - Title and breadcrumb */}
        <div className="flex items-center space-x-4">
          <h1 className="text-xl font-semibold text-gray-900">
            {getWorkspaceTitle()}
          </h1>
        </div>

        {/* Right side - Search, notifications, theme toggle, user */}
        <div className="flex items-center space-x-4">
          {/* Search */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
            <Input
              placeholder="Search projects, tasks..."
              className="pl-10 w-64"
            />
          </div>

          {/* Dark/Light Mode toggle */}
          <Button
            variant="ghost"
            size="sm"
            onClick={toggleTheme}
            className="h-9 w-9 p-0"
          >
            {themeConfig.mode === "light" ? (
              <Moon className="h-4 w-4" />
            ) : (
              <Sun className="h-4 w-4" />
            )}
          </Button>

          {/* Notifications */}
          <Button variant="ghost" size="sm" className="h-9 w-9 p-0 relative">
            <Bell className="h-4 w-4" />
            {notificationState.unreadCount > 0 && (
              <Badge
                variant="destructive"
                className="absolute -top-1 -right-1 h-5 w-5 text-xs p-0 flex items-center justify-center"
              >
                {notificationState.unreadCount > 9
                  ? "9+"
                  : notificationState.unreadCount}
              </Badge>
            )}
          </Button>

          {/* User menu */}
          <Button variant="ghost" size="sm" className="h-9 w-9 p-0">
            <User className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </header>
  );
}
