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
import { calculateLayout } from "@/utils/layout";
import { TaskNodeData } from "@/types/reactflow.types";

import { useRealTimeData } from "@/hooks/use-real-time-data";

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
}

export const ProjectsCanvas: React.FC<ReactFlowCanvasProps> = ({
  className,
}) => {
  const { currentProject } = useProjectStore();
  const { tasks } = useTaskStore();
  const { updateTaskPosition, updateProjectPosition } = useRealTimeData();
  const { setSelectedNodes, setRightSidebarCollapsed, setReactFlowNodes } =
    useUIStore();
  const { fitView } = useReactFlow();

  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  // Update nodes and edges when project or tasks change
  useEffect(() => {
    if (!currentProject) {
      setNodes([]);
      setEdges([]);
      return;
    }

    // Tasks are already filtered by currentProject in useRealTimeData
    const projectTasks = tasks;

    // Always use project layout when a project is selected
    const { nodes: newNodes, edges: newEdges } = calculateLayout(
      currentProject,
      projectTasks
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
              updateTaskPosition(change.id, change.position);
            } else if (node.type === "project") {
              updateProjectPosition(change.id, change.position);
            }
          }
        }
      });
    },
    [onNodesChange, updateTaskPosition, updateProjectPosition, nodes]
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

  return (
    <div className={`h-full w-full ${className}`}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={handleNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onSelectionChange={handleSelectionChange}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        connectionMode={ConnectionMode.Loose}
        fitView
        attributionPosition="bottom-left"
        className="bg-background"
        minZoom={0.1}
        maxZoom={2}
        defaultViewport={{ x: 0, y: 0, zoom: 1 }}
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
