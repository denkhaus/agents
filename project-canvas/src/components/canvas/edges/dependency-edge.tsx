/**
 * Dependency Edge Component
 * Custom ReactFlow edge for task dependencies
 */

import React from 'react';
import {
  EdgeProps,
  getBezierPath,
  EdgeLabelRenderer,
  BaseEdge,
} from 'reactflow';
import { Badge } from '@/components/ui/badge';
import { DependencyEdgeData } from '@/types/reactflow.types';
import { cn } from '@/lib/utils';

export const DependencyEdge: React.FC<EdgeProps<DependencyEdgeData>> = ({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  style = {},
  data,
  selected,
}) => {
  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  const isBlocking = data?.isBlocking || false;
  const dependencyType = data?.dependencyType || 'finish-to-start';

  // Edge styling based on dependency type and state
  const edgeStyle = {
    ...style,
    strokeWidth: selected ? 3 : 2,
    stroke: isBlocking ? '#ef4444' : '#6b7280',
    strokeDasharray: dependencyType === 'start-to-start' ? '5,5' : undefined,
  };

  // Arrow marker styling
  const markerId = `arrow-${id}`;

  return (
    <>
      {/* SVG Marker Definition */}
      <defs>
        <marker
          id={markerId}
          markerWidth="12"
          markerHeight="12"
          refX="6"
          refY="3"
          orient="auto"
          markerUnits="strokeWidth"
        >
          <path
            d="M0,0 L0,6 L9,3 z"
            fill={isBlocking ? '#ef4444' : '#6b7280'}
          />
        </marker>
      </defs>

      {/* Main Edge Path */}
      <BaseEdge
        path={edgePath}
        style={{
          ...edgeStyle,
          markerEnd: `url(#${markerId})`,
        }}
      />

      {/* Edge Label */}
      <EdgeLabelRenderer>
        <div
          style={{
            position: 'absolute',
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            fontSize: 12,
            pointerEvents: 'all',
          }}
          className="nodrag nopan"
        >
          {(selected || isBlocking) && (
            <Badge
              variant={isBlocking ? "destructive" : "secondary"}
              className={cn(
                "text-xs px-2 py-1 shadow-sm",
                "bg-background border border-border"
              )}
            >
              {isBlocking ? 'Blocked' : getDependencyTypeLabel(dependencyType)}
            </Badge>
          )}
        </div>
      </EdgeLabelRenderer>
    </>
  );
};

// Helper function to get human-readable dependency type labels
function getDependencyTypeLabel(type: string): string {
  switch (type) {
    case 'finish-to-start':
      return 'FS';
    case 'start-to-start':
      return 'SS';
    case 'finish-to-finish':
      return 'FF';
    default:
      return 'DEP';
  }
}