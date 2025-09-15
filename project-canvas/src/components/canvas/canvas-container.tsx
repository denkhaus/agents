/**
 * Canvas Container Component
 * Main container for ReactFlow visualization
 */

import React from 'react';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { 
  ZoomIn, 
  ZoomOut, 
  Maximize, 
  RotateCcw,
  Download,
  Settings
} from 'lucide-react';
import { cn } from '@/lib/utils';

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
  className
}) => {
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
            15 Tasks
          </Badge>
        </div>

        {/* Canvas Controls */}
        <div className="flex items-center gap-2">
          {/* Zoom Controls */}
          <div className="flex items-center gap-1 border border-border rounded-md">
            <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
              <ZoomOut className="h-4 w-4" />
            </Button>
            <div className="px-2 text-xs font-medium border-x border-border">
              100%
            </div>
            <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
              <ZoomIn className="h-4 w-4" />
            </Button>
          </div>

          {/* Action Buttons */}
          <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
            <RotateCcw className="h-4 w-4" />
          </Button>
          
          <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
            <Maximize className="h-4 w-4" />
          </Button>

          <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
            <Download className="h-4 w-4" />
          </Button>

          <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
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
          <span>15 nodes, 12 edges</span>
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
          <p className="text-sm text-muted-foreground">Loading project data...</p>
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