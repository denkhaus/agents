/**
 * ReactFlow Canvas Component
 * Main visualization component for projects and tasks
 */

import React, { useCallback, useEffect } from "react";
import {
  ReactFlow,
  useNodesState,
  useEdgesState,
  useReactFlow,
  Background,
  MiniMap,
  ConnectionMode,
  OnConnect,
  OnNodesChange,
  Node,
  Edge,
  NodeTypes,
  EdgeTypes,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";

import { TaskNode } from "./nodes/task-node";
import { ProjectNode } from "./nodes/project-node";
import { DependencyEdge } from "./edges/dependency-edge";
import { HierarchyEdge } from "./edges/hierarchy-edge";
import { useProjectStore, useTaskStore, useUIStore } from "@/stores";
import { calculateLayout, defaultLayoutOptions } from "@/utils/dagre-layout"; // Import defaultLayoutOptions
import { TaskNodeData } from "@/types/reactflow.types";

import { useCanvasNodePositionSync } from "@/hooks/use-canvas-node-position-sync"; // Import the unified hook

// Define custom node types
const nodeTypes: NodeTypes = {
  task: TaskNode,
  project: ProjectNode,
};

// Define custom edge types
const edgeTypes: EdgeTypes = {
  dependency: DependencyEdge,
  hierarchy: HierarchyEdge,
};

interface ReactFlowCanvasProps {
  className?: string;
  onMove?: (
    event: MouseEvent | TouchEvent | null,
    viewport: { x: number; y: number; zoom: number }
  ) => void;
  defaultViewport?: { x: number; y: number; zoom: number };
  onAutoLayoutRequest?: () => void; // Add this prop
}

export const ProjectsCanvas: React.FC<ReactFlowCanvasProps> = ({
  className,
  onMove,
  defaultViewport,
  onAutoLayoutRequest, // Destructure the new prop
}) => {
  const { currentProject } = useProjectStore();
  const { tasks } = useTaskStore();

  const { setSelectedNodes, setRightSidebarCollapsed, setReactFlowNodes } =
    useUIStore();
  const reactFlowInstance = useReactFlow();
  const { fitView, setViewport } = reactFlowInstance;

  const { updateNodePositionAndPersist: updateTaskNodePositionAndPersist } =
    useCanvasNodePositionSync("tasks", currentProject?.id || null);
  const { updateNodePositionAndPersist: updateProjectNodePositionAndPersist } =
    useCanvasNodePositionSync("projects", currentProject?.id || null);

  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  // Auto-layout function triggered by canvas container
  const handleAutoLayout = useCallback(() => {
    if (!currentProject) return;

    // Tasks are already filtered by currentProject
    const projectTasks = tasks;

    // Force new layout calculation
    const { nodes: newNodes, edges: newEdges } = calculateLayout(
      currentProject,
      projectTasks,
      {
        ...defaultLayoutOptions,
        force: true, // Force new layout
      }
    );

    setNodes(newNodes as Node[]);
    setEdges(newEdges as Edge[]);

    // Store nodes in UI store for property panel access
    setReactFlowNodes(newNodes as Node[]);

    // Persist new positions to store
    newNodes.forEach((node) => {
      if (node.type === "task") {
        updateTaskNodePositionAndPersist(node.id, node.position);
      } else if (node.type === "project") {
        updateProjectNodePositionAndPersist(node.id, node.position);
      }
    });

    // Auto-fit view after layout
    setTimeout(() => fitView({ padding: 0.2 }), 100);
  }, [
    currentProject,
    tasks,
    setNodes,
    setEdges,
    setReactFlowNodes,
    updateTaskNodePositionAndPersist,
    updateProjectNodePositionAndPersist,
    fitView,
  ]);

  // Expose auto-layout function to parent
  useEffect(() => {
    if (onAutoLayoutRequest) {
      // Store the handler so it can be called from parent
      (window as any).triggerProjectAutoLayout = handleAutoLayout;
    }
    return () => {
      delete (window as any).triggerProjectAutoLayout;
    };
  }, [handleAutoLayout, onAutoLayoutRequest]);

  // Update nodes and edges when project or tasks change
  useEffect(() => {
    if (!currentProject) {
      setNodes([]);
      setEdges([]);
      return;
    }

    // Tasks are already filtered by currentProject in useRealTimeData
    const projectTasks = tasks;

    // Use hierarchical layout with improved options
    const { nodes: newNodes, edges: newEdges } = calculateLayout(
      currentProject,
      projectTasks,
      {
        ...defaultLayoutOptions,
        force: false, // Don't force layout if positions exist
      }
    );

    setNodes(newNodes as Node[]);
    setEdges(newEdges as Edge[]);

    // Store nodes in UI store for property panel access
    setReactFlowNodes(newNodes as Node[]);

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentProject, tasks.length, setNodes, setEdges]);

  // Handle new connections (for future dependency management)
  const onConnect: OnConnect = useCallback(
    // eslint-disable-next-line no-console
    (params) => console.log("Connection made:", params),
    []
  );

  // Handle node position changes
  const handleNodesChange: OnNodesChange = useCallback(
    (changes) => {
      onNodesChange(changes);

      // Update positions in store for persistence
      changes.forEach((change) => {
        if (change.type === "position" && change.position) {
          // Find the node to determine its type
          const node = nodes.find((n) => n.id === change.id);
          if (node) {
            if (node.type === "task") {
              updateTaskNodePositionAndPersist(change.id, change.position);
            } else if (node.type === "project") {
              updateProjectNodePositionAndPersist(change.id, change.position);
            }
          }
        }
      });
    },
    [
      onNodesChange,
      updateTaskNodePositionAndPersist,
      updateProjectNodePositionAndPersist,
      nodes,
    ]
  );

  // Handle node selection changes
  const handleSelectionChange = useCallback(
    ({ nodes: selectedNodes }: { nodes: Node[] }) => {
      const selectedNodeIds = selectedNodes.map((node) => node.id);
      setSelectedNodes(selectedNodeIds);

      // Open the right sidebar when a node is selected
      if (selectedNodeIds.length > 0) {
        setRightSidebarCollapsed(false);
      }
    },
    [setSelectedNodes, setRightSidebarCollapsed]
  );

  // Auto-fit view when data changes
  useEffect(() => {
    if (nodes.length > 0) {
      setTimeout(() => fitView({ padding: 0.2 }), 100);
    }
  }, [nodes.length, fitView]);

  // Update viewport when defaultViewport prop changes (for persistence across views)
  useEffect(() => {
    const currentViewport = reactFlowInstance.getViewport(); // Get current viewport
    // Only update if defaultViewport is significantly different from currentViewport
    if (
      defaultViewport &&
      (Math.abs(defaultViewport.x - currentViewport.x) > 1 ||
        Math.abs(defaultViewport.y - currentViewport.y) > 1 ||
        Math.abs(defaultViewport.zoom - currentViewport.zoom) > 0.001)
    ) {
      setViewport(defaultViewport);
    }
  }, [defaultViewport, setViewport, reactFlowInstance]); // Add reactFlowInstance to dependencies

  return (
    <div className={`h-full w-full ${className}`}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={handleNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onSelectionChange={handleSelectionChange}
        onMove={onMove} // Pass the onMove prop here
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        connectionMode={ConnectionMode.Loose}
        fitView
        attributionPosition="bottom-left"
        className="bg-background"
        minZoom={0.1}
        maxZoom={2}
        defaultViewport={defaultViewport} // Use the passed defaultViewport
      >
        <Background
          color="#aaa"
          gap={16}
          className="opacity-20 dark:opacity-10"
        />

        <MiniMap
          className="bg-background border border-border rounded-md"
          maskColor="rgba(0, 0, 0, 0.1)"
          nodeColor={(node) => {
            if (node.type === "task") {
              const taskNodeData = node.data as TaskNodeData;
              switch (taskNodeData.task.state) {
                case "completed":
                  return "#22c55e";
                case "in-progress":
                  return "#3b82f6";
                case "blocked":
                  return "#ef4444";
                case "cancelled":
                  return "#6b7280";
                default:
                  return "#94a3b8";
              }
            }
            if (node.type === "project") {
              return "#8b5cf6"; // Purple for project nodes
            }
            return "#8b5cf6";
          }}
        />
      </ReactFlow>
    </div>
  );
};
