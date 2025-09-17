import React from "react";
import { BaseEdge, EdgeProps, getBezierPath } from "@xyflow/react";

export const HierarchyEdge: React.FC<EdgeProps> = ({
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  style = {},

  selected,
}) => {
  const [edgePath] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  // You can add custom styling or labels here if needed
  const edgeStyle = {
    ...style,
    strokeWidth: selected ? 2 : 1,
    stroke: "#a3a3a3", // A neutral color for hierarchy edges
  };

  return (
    <>
      <BaseEdge path={edgePath} style={edgeStyle} />
      {/* Optional: Add a label for hierarchy edges */}
      {/* <EdgeLabelRenderer>
        <div
          style={{
            position: 'absolute',
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            fontSize: 10,
            pointerEvents: 'all',
          }}
          className="nodrag nopan"
        >
          Hierarchy
        </div>
      </EdgeLabelRenderer> */}
    </>
  );
};
