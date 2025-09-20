/**
 * Canvas Container Component
 * Main container for ReactFlow visualization
 */

import React, { useCallback, useMemo } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  ZoomIn,
  ZoomOut,
  Maximize,
  Download,
  Settings,
  Layout,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useReactFlow, ReactFlowProvider, Viewport } from "@xyflow/react"; // Import Viewport type
import { useSettingsStore } from "@/stores";
import { getNodeTypeForWorkspace, type WorkspaceType } from "@/constants";
import { NodeTypeRegistry } from "@/registry";

interface CanvasContainerProps {
  children: React.ReactNode;
  title?: string;
  subtitle?: string;
  className?: string;
  onViewportChange?: (viewport: { x: number; y: number; zoom: number }) => void;
  onAutoLayoutRequest?: () => void; // New prop for auto-layout
}

export const CanvasContainer: React.FC<CanvasContainerProps> = ({
  children,
  title = "Canvas Visualization",
  subtitle = "Interactive flow diagram",
  className,
  onAutoLayoutRequest,
}) => {
  return (
    <ReactFlowProvider>
      <CanvasContainerInner
        title={title}
        subtitle={subtitle}
        className={className}
        onAutoLayoutRequest={onAutoLayoutRequest} // Pass the prop down
      >
        {children}
      </CanvasContainerInner>
    </ReactFlowProvider>
  );
};

const CanvasContainerInner: React.FC<CanvasContainerProps> = ({
  children,
  title = "Canvas Visualization",
  subtitle = "Interactive flow diagram",
  className,
  onViewportChange,
  onAutoLayoutRequest, // Destructure the new prop
}) => {
  const reactFlowInstance = useReactFlow();
  const {
    updateCanvasSettings,
    projectCanvasViewport,
    agentsCanvasViewport,
    canvasZoom, // Read canvasZoom from the store
  } = useSettingsStore();
  const location = useLocation();
  const navigate = useNavigate();

  // Determine the current viewport based on the active route
  const initialViewport = useMemo(() => {
    if (location.pathname === "/projects") {
      return projectCanvasViewport;
    } else if (location.pathname === "/agents") {
      return agentsCanvasViewport;
    }
    return { x: 0, y: 0, zoom: 1 }; // Default or fallback
  }, [location.pathname, projectCanvasViewport, agentsCanvasViewport]);

  // Handler for viewport changes (pan, zoom)
  const handleOnViewportChange = useCallback(
    (_event: MouseEvent | TouchEvent | null, viewport: Viewport) => {
      const newPersistedViewport = {
        x: viewport.x,
        y: viewport.y,
        zoom: viewport.zoom,
      };
      updateCanvasSettings(
        newPersistedViewport,
        location.pathname === "/projects" ? "project" : "agents"
      );

      // Call external onViewportChange if provided
      if (onViewportChange) {
        onViewportChange({ x: viewport.x, y: viewport.y, zoom: viewport.zoom });
      }
    },
    [updateCanvasSettings, location.pathname, onViewportChange]
  );

  // Get current workspace configuration
  const currentNodeType = useMemo(() => {
    // Map pathname to WorkspaceType for getNodeTypeForWorkspace
    const workspaceType = location.pathname.substring(1) as WorkspaceType;
    return getNodeTypeForWorkspace(workspaceType);
  }, [location.pathname]);

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

  // UNIFIED Download Funktion - Workspace-agnostic
  const handleDownload = () => {
    // Map pathname to WorkspaceType for NodeTypeRegistry
    const workspaceType = location.pathname.substring(1) as WorkspaceType;
    const config = NodeTypeRegistry.getConfigForWorkspace(workspaceType);
    if (!config) {
      console.warn("No configuration found for current workspace");
      return;
    }

    // Create export data based on current workspace
    const exportData = {
      workspace: workspaceType,
      nodeType: currentNodeType,
      timestamp: new Date().toISOString(),
      nodeCount: 0, // Will be calculated by the component
    };

    // Convert to JSON
    const jsonData = JSON.stringify(exportData, null, 2);

    // Create blob and download link
    const blob = new Blob([jsonData], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${location.pathname.substring(1)}_export_${Date.now()}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  // Settings Funktion - Navigiert zum Settings Workspace
  const handleSettings = () => {
    navigate("/settings");
  };

  // UNIFIED Auto-layout Funktion - Node-agnostic
  const handleAutoLayout = useCallback(() => {
    if (onAutoLayoutRequest) {
      onAutoLayoutRequest();
    }
    
    // Also trigger the specific canvas auto-layout based on current workspace
    const currentPath = location.pathname;
    if (currentPath.includes('/agents')) {
      // Trigger agent canvas auto-layout
      if ((window as any).triggerAgentAutoLayout) {
        (window as any).triggerAgentAutoLayout();
      }
    } else if (currentPath.includes('/projects')) {
      // Trigger project canvas auto-layout
      if ((window as any).triggerProjectAutoLayout) {
        (window as any).triggerProjectAutoLayout();
      }
    }
  }, [onAutoLayoutRequest, location.pathname]);

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
              {canvasZoom}%
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
        {React.Children.map(children, (child) =>
          React.isValidElement(child)
            ? React.cloneElement(child as React.ReactElement<any>, {
                onMove: handleOnViewportChange, // Use onMove for general viewport changes
                defaultViewport: initialViewport, // Pass initialViewport to child canvas
                onAutoLayoutRequest: handleAutoLayout, // Pass auto-layout handler to canvas
              })
            : child
        )}

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
          <span>Canvas View</span>
          <span>•</span>
          <span>Connected</span>
        </div>
      </div>
    </div>
  );
};

// Canvas overlay for states
export const CanvasOverlay: React.FC = () => {
  const [isLoading] = React.useState(false);
  const [isEmpty] = React.useState(false);

  if (isLoading) {
    return (
      <div className="absolute inset-0 flex items-center justify-center bg-background/80 backdrop-blur-sm">
        <Card className="p-6 text-center">
          <div className="animate-spin h-8 w-8 border-2 border-primary border-t-transparent rounded-full mx-auto mb-4" />
          <p className="text-sm text-muted-foreground">Loading data...</p>
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
            Select a item from the sidebar to view its visualization.
          </p>
        </Card>
      </div>
    );
  }

  return null;
};
