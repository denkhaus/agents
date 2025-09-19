/**
 * Agents Canvas Component
 * ReactFlow canvas for agent visualization
 */

import React, { useCallback, useEffect, useMemo } from "react";
import {
  ReactFlow,
  Background,
  MiniMap,
  useNodesState,
  useEdgesState,
  addEdge,
  Connection,
  Edge,
  Node,
  NodeTypes,
  EdgeTypes,
  useReactFlow,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { AgentNode } from "@/components/canvas/nodes";
import { useAgentProjectStore, useUIStore } from "@/stores";
import { useAgentProjectData } from "@/hooks/use-agent-project-data";

// Define node types for ReactFlow
const nodeTypes: NodeTypes = {
  agent: AgentNode,
};

// Define edge types (using default for now)
const edgeTypes: EdgeTypes = {};

interface AgentsCanvasProps {
  className?: string;
}

export const AgentsCanvas: React.FC<AgentsCanvasProps> = ({ className }) => {
  // Initialize data
  const { loading: agentDataLoading } = useAgentProjectData();

  const { currentAgentProject } = useAgentProjectStore();
  const { setSelectedNodes, setRightSidebarCollapsed } = useUIStore();
  const { fitView } = useReactFlow();



  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  // Convert agent project data to ReactFlow nodes and edges
  const { reactFlowNodes, reactFlowEdges } = useMemo(() => {
    if (!currentAgentProject || !currentAgentProject.agentNodes) {
      return { reactFlowNodes: [], reactFlowEdges: [] };
    }

    // Convert agent nodes - ensure they are agent type
    const reactFlowNodes: Node[] = currentAgentProject.agentNodes
      .filter(agentNode => agentNode.type === "agent")
      .map((agentNode) => ({
        id: agentNode.id, // Use agentNode.id directly
        type: "agent",
        position: agentNode.position,
        data: agentNode.data,
      }));

    // Convert connections to edges - only for agent connections
    const reactFlowEdges: Edge[] = currentAgentProject.connections
      ? currentAgentProject.connections
          .filter(connection => {
            // Ensure we only process valid connection types
            const validTypes = ["hierarchy", "communication", "collaboration"];
            return validTypes.includes(connection.type);
          })
          .map((connection) => ({
            id: connection.id,
            source: connection.source,
            target: connection.target,
            type: "default", // Always use default edge type
            label: connection.label,
            data: connection.data,
            style: {
              stroke:
                connection.type === "hierarchy"
                  ? "#8b5cf6"
                  : connection.type === "communication"
                    ? "#3b82f6"
                    : "#10b981",
              strokeWidth: 2,
            },
          }))
      : [];
    
    return { reactFlowNodes, reactFlowEdges };
  }, [currentAgentProject]);

  // Update nodes and edges when agent project changes
  useEffect(() => {
    // Additional safety check - ensure all nodes are agent type
    const safeNodes = reactFlowNodes.filter(node => node.type === "agent");
    
    // Additional safety check - ensure all edges use default type
    const safeEdges = reactFlowEdges.map(edge => ({
      ...edge,
      type: "default" // Force all edges to use default type
    }));
    
    setNodes(safeNodes);
    setEdges(safeEdges);
  }, [reactFlowNodes, reactFlowEdges, setNodes, setEdges]);

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

  // Handle edge creation
  const onConnect = useCallback(
    (params: Connection) => {
      const newEdge: Edge = {
        id: crypto.randomUUID(),
        source: params.source!,
        target: params.target!,
        sourceHandle: params.sourceHandle,
        targetHandle: params.targetHandle,
        type: "default",
        style: { stroke: "#3b82f6", strokeWidth: 2 },
      };
      setEdges((eds) => addEdge(newEdge, eds));
    },
    [setEdges]
  );

  // Auto-fit view when data changes
  useEffect(() => {
    if (nodes.length > 0) {
      setTimeout(() => fitView({ padding: 0.2 }), 100);
    }
  }, [nodes.length, fitView]);

  // Show loading state
  if (agentDataLoading) {
    return (
      <div className={`h-full w-full ${className} flex items-center justify-center`}>
        <div className="text-center">
          <div className="animate-spin h-8 w-8 border-2 border-primary border-t-transparent rounded-full mx-auto mb-4" />
          <p className="text-muted-foreground">Loading agents from Convex...</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`h-full w-full ${className}`}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onSelectionChange={handleSelectionChange}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
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
            if (node.type === "agent") {
              return "#8b5cf6"; // Purple for agent nodes
            }
            return "#8b5cf6";
          }}
          position="bottom-right"
        />
      </ReactFlow>
    </div>
  );
};
