"use client";

import { useState, useEffect } from "react";
import { useProjectStore } from "@/lib/stores/project-store";
import { useUIStore } from "@/lib/stores/ui-store";
import { ProjectsHeader } from "./projects-header";
import { ProjectsList } from "./projects-list";
import { KanbanView } from "@/components/kanban/kanban-view";
import { ChatPanel } from "@/components/chat/chat-panel";
import { cn } from "@/lib/utils";

export function ProjectsWorkspace() {
  const {
    currentProject,
    projects,
    isLoading,
    setProjects,
    setCurrentProject,
    setLoading,
    setTasks,
  } = useProjectStore();

  const { selectionState, isChatPanelOpen, viewConfig } = useUIStore();

  const [isInitialized, setIsInitialized] = useState(false);

  // Initialize data on mount
  useEffect(() => {
    const initializeData = async () => {
      if (isInitialized) return;

      setLoading(true);
      try {
        // Import corrected mock data with proper UUIDs
        const { mockProjects, mockTasks } = await import(
          "@/lib/data/mock-data-corrected"
        );

        setProjects(mockProjects);
        setTasks(mockTasks); // Load all tasks for statistics
        setIsInitialized(true);
      } catch (error) {
        console.error("Failed to initialize projects:", error);
      } finally {
        setLoading(false);
      }
    };

    initializeData();
  }, [isInitialized, setProjects, setLoading]);

  if (isLoading && !isInitialized) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading projects...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full">
      {/* Main content area */}
      <div
        className={cn(
          "flex-1 flex flex-col",
          isChatPanelOpen && "mr-96" // Make room for chat panel
        )}
      >
        <ProjectsHeader />

        <div className="flex-1 overflow-hidden">
          {!selectionState?.selectedProjectId ? (
            <ProjectsList />
          ) : (
            <KanbanView />
          )}
        </div>
      </div>

      {/* Chat panel */}
      {isChatPanelOpen && (
        <div className="w-96 border-l border-border bg-card">
          <ChatPanel />
        </div>
      )}
    </div>
  );
}
