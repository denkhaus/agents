/**
 * Main Application Component
 * Updated to use the new layout system
 */

import React from "react";
import { AppLayout } from "@/components/layout/app-layout";
import { CanvasContainer } from "@/components/canvas/canvas-container";
import { ProjectsCanvas } from "@/components/canvas/projects-canvas";
import { AgentsCanvas } from "@/components/canvas/agents-canvas";
import { SettingsCanvas } from "@/components/canvas/settings-canvas";
import { Toaster } from "@/components/ui/toaster";
import { useUIStore } from "@/stores";
import { useRealTimeData } from "@/hooks/use-real-time-data";
import { useSettingsSync } from "@/hooks/use-settings-sync";

function App() {
  const { currentWorkspace } = useUIStore();

  // Real-time data integration
  const { loading } = useRealTimeData();

  // Sync settings with Convex
  useSettingsSync();

  // Show loading only for a short time, then show app even if no data
  const [showLoading, setShowLoading] = React.useState(true);

  React.useEffect(() => {
    const timer = setTimeout(() => setShowLoading(false), 2000);
    return () => clearTimeout(timer);
  }, []);

  if (loading && showLoading) {
    return (
      <div className="h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin h-8 w-8 border-2 border-primary border-t-transparent rounded-full mx-auto mb-4" />
          <p className="text-muted-foreground">Loading project data...</p>
        </div>
      </div>
    );
  }

  const renderWorkspaceContent = () => {
    switch (currentWorkspace) {
      case "projects":
        return (
          <CanvasContainer
            title="Project Canvas"
            subtitle="Interactive project and task visualization"
          >
            <ProjectsCanvas />
          </CanvasContainer>
        );
      
      case "agents":
        return (
          <CanvasContainer
            title="Agent Canvas"
            subtitle="Agent flow and collaboration visualization"
          >
            <AgentsCanvas />
          </CanvasContainer>
        );
      
      case "settings":
        return (
          <CanvasContainer
            title="Application Settings"
            subtitle="Configure your workspace preferences"
          >
            <SettingsCanvas />
          </CanvasContainer>
        );
      
      default:
        return null;
    }
  };

  return (
    <AppLayout>
      {renderWorkspaceContent()}
      <Toaster />
    </AppLayout>
  );
}

export default App;
