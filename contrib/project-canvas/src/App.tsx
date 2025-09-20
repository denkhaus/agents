/**
 * Main Application Component
 * Updated to use the new layout system
 */

import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AppLayout } from "@/components/layout/app-layout";
import { CanvasContainer } from "@/components/canvas/canvas-container";
import { ProjectsCanvas } from "@/components/canvas/projects-canvas";
import { AgentsCanvas } from "@/components/canvas/agents-canvas";
import { SettingsCanvas } from "@/components/canvas/settings-canvas";
import { Toaster } from "@/components/ui/toaster";
import { useSettingsSync } from "@/hooks/use-settings-sync";
import { useConvexProjects, useConvexTasks } from "@/hooks/use-convex-data";
import { useProjectStore } from "@/stores";

function App() {
  // Sync settings with Convex
  useSettingsSync();

  // Load projects and tasks from Convex
  useConvexProjects();
  const { currentProject } = useProjectStore();
  useConvexTasks(currentProject?.id);

  return (
    <BrowserRouter>
      <AppLayout>
        <Routes>
          <Route path="/" element={<Navigate to="/projects" replace />} />
          <Route
            path="/projects"
            element={
              <CanvasContainer
                key="projects-workspace"
                title="Project Canvas"
                subtitle="Interactive project and task visualization"
              >
                <ProjectsCanvas />
              </CanvasContainer>
            }
          />
          <Route
            path="/agents"
            element={
              <CanvasContainer
                key="agents-workspace"
                title="Agent Canvas"
                subtitle="Agent flow and collaboration visualization"
              >
                <AgentsCanvas />
              </CanvasContainer>
            }
          />
          <Route
            path="/settings"
            element={
              <CanvasContainer
                key="settings-workspace"
                title="Application Settings"
                subtitle="Configure your workspace preferences"
              >
                <SettingsCanvas />
              </CanvasContainer>
            }
          />
        </Routes>
        <Toaster />
      </AppLayout>
    </BrowserRouter>
  );
}

export default App;
