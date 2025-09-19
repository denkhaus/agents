/**
 * Canvas Container Component
 * Main container for ReactFlow visualization
 */

import React, { useState, useCallback, useMemo } from "react";
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
import { useReactFlow, ReactFlowProvider } from "@xyflow/react";
import { useUIStore, useSettingsStore } from "@/stores";
import {
  WORKSPACE_TYPES,
  getNodeTypeForWorkspace,
  isValidNodeType,
  type NodeType,
  type WorkspaceType,
} from "@/constants";
import { StoreRegistry } from "@/stores/store-registry";
import { NodeTypeRegistry } from "@/registry";
import "@/stores/registry-initializer"; // Auto-initialize registries

interface CanvasContainerProps {
  children: React.ReactNode;
  title?: string;
  subtitle?: string;
  className?: string;
}

export const CanvasContainer: React.FC<CanvasContainerProps> = ({
  children,
  title = "Canvas Visualization",
  subtitle = "Interactive flow diagram",
  className,
}) => {
  return (
    <ReactFlowProvider>
      <CanvasContainerInner
        title={title}
        subtitle={subtitle}
        className={className}
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
}) => {
  const reactFlowInstance = useReactFlow();
  const { canvasZoom, updateCanvasSettings } = useSettingsStore();
  const { currentWorkspace } = useUIStore();

  // Zustand für Zoomstufe aus Settings Store
  const [zoomLevel, setZoomLevel] = useState(canvasZoom);

  // Get current workspace configuration
  const currentNodeType = useMemo(() => {
    return getNodeTypeForWorkspace(currentWorkspace as WorkspaceType);
  }, [currentWorkspace]);

  // Get node count and display name using registry
  const { nodeCount, displayName } = useMemo(() => {
    const config = NodeTypeRegistry.getConfigForWorkspace(
      currentWorkspace as WorkspaceType
    );
    if (config) {
      return {
        nodeCount: config.operations.getCount(),
        displayName: config.displayName,
      };
    }
    return { nodeCount: 0, displayName: "Items" };
  }, [currentWorkspace]);

  // Update zoom level when reactFlowInstance changes or when zoom changes
  React.useEffect(() => {
    if (reactFlowInstance) {
      const currentZoom = Math.round(reactFlowInstance.getZoom() * 100);
      setZoomLevel(currentZoom);

      // Save to settings store if different
      if (currentZoom !== canvasZoom) {
        updateCanvasSettings({ canvasZoom: currentZoom });
      }
    }
  }, [reactFlowInstance, canvasZoom, updateCanvasSettings]);

  // Add event listener for zoom changes
  React.useEffect(() => {
    if (!reactFlowInstance) return;

    // Use a timer to periodically check for zoom changes
    // This is a workaround since onViewportChange might not be available
    const interval = setInterval(() => {
      const currentZoom = Math.round(reactFlowInstance.getZoom() * 100);
      if (currentZoom !== zoomLevel) {
        setZoomLevel(currentZoom);
        updateCanvasSettings({ canvasZoom: currentZoom });
      }
    }, 500);

    return () => {
      clearInterval(interval);
    };
  }, [reactFlowInstance, updateCanvasSettings, zoomLevel]);

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
    const config = NodeTypeRegistry.getConfigForWorkspace(
      currentWorkspace as WorkspaceType
    );
    if (!config) {
      console.warn("No configuration found for current workspace");
      return;
    }

    // Create export data based on current workspace
    const exportData = {
      workspace: currentWorkspace,
      nodeType: currentNodeType,
      timestamp: new Date().toISOString(),
      nodeCount: config.operations.getCount(),
    };

    // Convert to JSON
    const jsonData = JSON.stringify(exportData, null, 2);

    // Create blob and download link
    const blob = new Blob([jsonData], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${currentWorkspace}_export_${Date.now()}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  // Settings Funktion - Aktiviert den Settings Workspace
  const { setWorkspace } = useUIStore();
  const handleSettings = () => {
    setWorkspace(WORKSPACE_TYPES.SETTINGS);
  };

  // UNIFIED Auto-layout Funktion - Node-agnostic
  const handleAutoLayout = useCallback(() => {
    if (!reactFlowInstance) {
      console.warn("ReactFlow instance not available for auto-layout.");
      return;
    }

    const nodes = reactFlowInstance.getNodes();
    if (nodes.length === 0) {
      console.warn("No nodes available for auto-layout.");
      return;
    }

    // Use existing layout calculation or simple grid layout
    const gridSize = Math.ceil(Math.sqrt(nodes.length));
    const spacing = 250;

    const newNodes = nodes.map((node, index) => {
      const row = Math.floor(index / gridSize);
      const col = index % gridSize;

      return {
        ...node,
        position: {
          x: col * spacing + 100,
          y: row * spacing + 100,
        },
      };
    });

    // Update node positions in ReactFlow
    reactFlowInstance.setNodes(newNodes);

    // UNIFIED POSITION UPDATE - Node-agnostic approach
    newNodes.forEach((node: any) => {
      if (node.position && isValidNodeType(node.type)) {
        const nodeType = node.type as NodeType;
        const storeOperations = StoreRegistry.getStoreForNodeType(nodeType);

        if (storeOperations?.updatePosition) {
          storeOperations
            .updatePosition(node.id, node.position)
            .catch((error) => {
              console.error(
                `Error updating position for ${nodeType} node ${node.id}:`,
                error
              );
            });
        } else {
          console.warn(
            `No position update operation found for node type: ${nodeType}`
          );
        }
      }
    });

    // Fit view after layout
    setTimeout(() => {
      reactFlowInstance.fitView({ padding: 0.2 });
    }, 100);
  }, [reactFlowInstance]);

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
            {nodeCount} {displayName}
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
            {nodeCount} nodes, {Math.max(0, nodeCount - 1)} edges
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
