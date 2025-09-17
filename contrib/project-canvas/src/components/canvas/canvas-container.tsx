/**
 * Canvas Container Component
 * Main container for ReactFlow visualization
 */

import React, { useState, useCallback } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  ZoomIn,
  ZoomOut,
  Maximize,
  Download,
  Settings,
  Layout,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useReactFlow, ReactFlowProvider } from "reactflow";
import { useProjectStore, useUIStore } from "@/stores";
import { useTaskStore } from "@/stores";
import { useRealTimeData } from "@/hooks/use-real-time-data";
import { calculateTaskLayout, defaultLayoutOptions } from "@/utils/layout";
import { Project } from "@/types/project.types";
import { WorkspaceType } from "@/types";

interface CanvasContainerProps {
  children: React.ReactNode;
  title?: string;
  subtitle?: string;
  className?: string;
}

export const CanvasContainer: React.FC<CanvasContainerProps> = ({
  children,
  title = "Project Visualization",
  subtitle = "Interactive task flow diagram",
  className,
}) => {
  const { currentProject } = useProjectStore();
  const { setWorkspace } = useUIStore();
  const { tasksByProject } = useTaskStore();

  // Task-Anzahl für aktuelles Projekt berechnen
  const taskCount = currentProject
    ? tasksByProject[currentProject.id]?.length || 0
    : 0;

  return (
    <ReactFlowProvider>
      <CanvasContainerInner
        title={title}
        subtitle={subtitle}
        className={className}
        taskCount={taskCount}
        setWorkspace={setWorkspace}
        currentProject={currentProject}
        tasksByProject={tasksByProject}
      >
        {children}
      </CanvasContainerInner>
    </ReactFlowProvider>
  );
};

interface CanvasContainerInnerProps extends CanvasContainerProps {
  taskCount: number;
  setWorkspace: (workspace: WorkspaceType) => void;
  currentProject: Project | null;
  tasksByProject: Record<string, any[]>;
}

const CanvasContainerInner: React.FC<CanvasContainerInnerProps> = ({
  children,
  title = "Project Visualization",
  subtitle = "Interactive task flow diagram",
  className,
  taskCount,
  setWorkspace,
  currentProject,
  tasksByProject,
}) => {
  const reactFlowInstance = useReactFlow();

  // Zustand für Zoomstufe
  const [zoomLevel, setZoomLevel] = useState(100);

  // Update zoom level when reactFlowInstance changes
  React.useEffect(() => {
    if (reactFlowInstance) {
      setZoomLevel(Math.round(reactFlowInstance.getZoom() * 100));
    }
  }, [reactFlowInstance]);

  // Zoom-In Funktion
  const handleZoomIn = () => {
    if (reactFlowInstance) {
      const newZoom = reactFlowInstance.getZoom() * 1.2;
      reactFlowInstance.zoomTo(newZoom, { duration: 200 });
    }
  };

  // Zoom-Out Funktion
  const handleZoomOut = () => {
    if (reactFlowInstance) {
      const newZoom = reactFlowInstance.getZoom() * 0.8;
      reactFlowInstance.zoomTo(newZoom, { duration: 200 });
    }
  };

  // Fit-to-View Funktion
  const handleFitView = () => {
    if (reactFlowInstance) {
      reactFlowInstance.fitView({ padding: 0.2, duration: 300 });
    }
  };

  // Download Funktion - Exportiert Projekt als JSON
  const handleDownload = () => {
    if (!currentProject) {
      console.warn("No project selected for download");
      return;
    }

    // Erstelle ein Datenobjekt mit Projekt und Tasks
    const projectData = {
      project: currentProject,
      tasks: tasksByProject[currentProject.id] || [],
    };

    // Konvertiere in JSON
    const jsonData = JSON.stringify(projectData, null, 2);

    // Erstelle Blob und Download-Link
    const blob = new Blob([jsonData], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${currentProject.title.replace(/\s+/g, "_")}_export.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  // Settings Funktion - Aktiviert den Settings Workspace
  const handleSettings = () => {
    setWorkspace("settings");
  };

  // Auto-layout Funktion
  const { tasks } = useTaskStore();
  const { updateTaskPosition } = useRealTimeData();

  const handleAutoLayout = useCallback(() => {
    const { nodes: newNodes } = calculateTaskLayout(tasks, {
      ...defaultLayoutOptions,
      force: true,
    });

    // Update node positions in ReactFlow
    if (reactFlowInstance) {
      reactFlowInstance.setNodes((prevNodes) => {
        return prevNodes.map((node) => {
          const updatedNode = newNodes.find((n) => n.id === node.id);
          if (updatedNode) {
            return {
              ...node,
              position: updatedNode.position,
            };
          }
          return node;
        });
      });
    }

    // Persist the new positions for all tasks
    newNodes.forEach((node) => {
      if (node.position) {
        updateTaskPosition(node.id, node.position);
      }
    });
  }, [tasks, updateTaskPosition, reactFlowInstance]);

  return (
    <div className={cn("h-full flex flex-col", className)}>
      {/* Canvas Header */}
      <div className="flex items-center justify-between p-4 border-b border-border bg-background/95 backdrop-blur">
        {/* Title Section */}
        <div className="flex items-center gap-3">
          <div>
            <h2 className="text-lg font-semibold">{title}</h2>
            <p className="text-sm text-muted-foreground">{subtitle}</p>
          </div>
          <Badge variant="secondary" className="ml-2">
            {taskCount} Tasks
          </Badge>
        </div>

        {/* Canvas Controls */}
        <div className="flex items-center gap-2">
          {/* Zoom Controls */}
          <div className="flex items-center gap-1 border border-border rounded-md">
            <Button
              variant="ghost"
              size="sm"
              className="h-8 w-8 p-0"
              onClick={handleZoomOut}
            >
              <ZoomOut className="h-4 w-4" />
            </Button>
            <div className="px-2 text-xs font-medium border-x border-border">
              {zoomLevel}%
            </div>
            <Button
              variant="ghost"
              size="sm"
              className="h-8 w-8 p-0"
              onClick={handleZoomIn}
            >
              <ZoomIn className="h-4 w-4" />
            </Button>
          </div>
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={handleFitView}
          >
            <Maximize className="h-4 w-4" />
          </Button>

          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={handleAutoLayout}
          >
            <Layout className="h-4 w-4" />
          </Button>

          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={handleDownload}
          >
            <Download className="h-4 w-4" />
          </Button>

          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={handleSettings}
          >
            <Settings className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Canvas Content */}
      <div className="flex-1 relative overflow-hidden bg-muted/20">
        {children}

        {/* Canvas Overlay - Loading/Empty States */}
        <CanvasOverlay />
      </div>

      {/* Canvas Footer */}
      <div className="flex items-center justify-between p-2 border-t border-border bg-background/95 backdrop-blur text-xs text-muted-foreground">
        <div className="flex items-center gap-4">
          <span>Last updated: 2 minutes ago</span>
          <span>•</span>
          <span>Auto-save enabled</span>
        </div>
        <div className="flex items-center gap-4">
          <span>
            {taskCount} nodes, {Math.max(0, taskCount - 1)} edges
          </span>
          <span>•</span>
          <span>Connected</span>
        </div>
      </div>
    </div>
  );
};

// Canvas overlay for states
const CanvasOverlay: React.FC = () => {
  const [isLoading] = React.useState(false);
  const [isEmpty] = React.useState(false);

  if (isLoading) {
    return (
      <div className="absolute inset-0 flex items-center justify-center bg-background/80 backdrop-blur-sm">
        <Card className="p-6 text-center">
          <div className="animate-spin h-8 w-8 border-2 border-primary border-t-transparent rounded-full mx-auto mb-4" />
          <p className="text-sm text-muted-foreground">
            Loading project data...
          </p>
        </Card>
      </div>
    );
  }

  if (isEmpty) {
    return (
      <div className="absolute inset-0 flex items-center justify-center">
        <Card className="p-8 text-center max-w-md">
          <div className="h-12 w-12 rounded-full bg-muted flex items-center justify-center mx-auto mb-4">
            <Settings className="h-6 w-6 text-muted-foreground" />
          </div>
          <h3 className="text-lg font-semibold mb-2">No Project Selected</h3>
          <p className="text-sm text-muted-foreground mb-4">
            Select a project from the sidebar to view its task visualization.
          </p>
          <Button>Browse Projects</Button>
        </Card>
      </div>
    );
  }

  return null;
};
