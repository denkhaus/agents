/**
 * Main Application Component
 * Updated to use the new layout system
 */

import React from "react";
import { AppLayout } from "@/components/layout/app-layout";
import { CanvasContainer } from "@/components/canvas/canvas-container";
import { ReactFlowCanvas } from "@/components/canvas/reactflow-canvas";
import { Toaster } from "@/components/ui/toaster";
import { useUIStore } from "@/stores";
import { useRealTimeData } from "@/hooks/use-real-time-data";
import { useDummyData } from "@/hooks/use-dummy-data";
import { useSettingsSync } from "@/hooks/use-settings-sync";

function App() {
  const { currentWorkspace } = useUIStore();
  
  // Real-time data integration
  const { loading } = useRealTimeData();
  
  // Load dummy data for development (fallback when Convex is empty)
  useDummyData();
  
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

  return (
    <AppLayout>
      {currentWorkspace === "projects" && (
        <CanvasContainer>
          <ReactFlowCanvas />
        </CanvasContainer>
      )}

      {currentWorkspace === "agents" && (
        <div className="h-full flex items-center justify-center">
          <div className="text-center">
            <h3 className="text-lg font-semibold mb-2">Agents Dashboard</h3>
            <p className="text-muted-foreground">
              Agent management interface will be here
            </p>
          </div>
        </div>
      )}

      {currentWorkspace === "settings" && (
        <div className="h-full flex items-center justify-center">
          <div className="text-center">
            <h3 className="text-lg font-semibold mb-2">Settings Panel</h3>
            <p className="text-muted-foreground">
              Application settings will be here
            </p>
          </div>
        </div>
      )}
      <Toaster />
    </AppLayout>
  );
}

export default App;
