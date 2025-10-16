/**
 * Agents Canvas Component
 * ReactFlow canvas for agent visualization
 */

import React, { useCallback, useEffect, useMemo, useRef } from "react";
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
  NodeChange,
  NodePositionChange,
  useReactFlow, // Import useReactFlow
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { AgentNode } from "@/components/canvas/nodes";
import { useAgentProjectStore, useUIStore } from "@/stores";
import { useAgentProjectData } from "@/hooks/use-agent-project-data";
import { useCanvasNodePositionSync } from "@/hooks/use-canvas-node-position-sync"; // Import the new hook
import {
  calculateAgentLayout,
  calculateGridLayout,
  defaultAgentLayoutOptions,
} from "@/utils/agent-layout";

// Define node types for ReactFlow
const nodeTypes: NodeTypes = {
  agent: AgentNode,
};

// Define edge types (using default for now)
const edgeTypes: EdgeTypes = {};

interface AgentsCanvasProps {
  className?: string;
  onMove?: (
    event: MouseEvent | TouchEvent | null,
    viewport: { x: number; y: number; zoom: number }
  ) => void;
  defaultViewport?: { x: number; y: number; zoom: number };
  onAutoLayoutRequest?: () => void; // Add auto-layout callback
}

export const AgentsCanvas: React.FC<AgentsCanvasProps> = ({
  className,
  onMove,
  defaultViewport,
  onAutoLayoutRequest,
}) => {
  // Initialize data
  const { loading: agentDataLoading } = useAgentProjectData();

  const { currentAgentProject, setAgentProjects, setCurrentAgentProject } = useAgentProjectStore();
  const { setSelectedNodes, setRightSidebarCollapsed } = useUIStore();
  const { updateNodePositionAndPersist } = useCanvasNodePositionSync(
    "agentProjects",
    currentAgentProject?.id || null
  );
  const reactFlowInstance = useReactFlow();
  const { setViewport } = reactFlowInstance;

  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  
  // Debounce timer for position updates
  const positionUpdateTimers = useRef<Map<string, NodeJS.Timeout>>(new Map());

  // Custom onNodesChange handler to sync position changes with Convex store
  const handleNodesChange = useCallback(
    (changes: NodeChange[]) => {
      // Apply changes to local state immediately for smooth UI
      onNodesChange(changes);

      // Process position changes with debouncing
      changes.forEach((change) => {
        if (change.type === "position" && currentAgentProject) {
          const positionChange = change as NodePositionChange;
          if (positionChange.position) {
            // Clear existing timer for this node
            const existingTimer = positionUpdateTimers.current.get(positionChange.id);
            if (existingTimer) {
              clearTimeout(existingTimer);
            }

            // Set new debounced timer (500ms delay)
            const newTimer = setTimeout(() => {
              updateNodePositionAndPersist(
                positionChange.id,
                positionChange.position!
              );
              positionUpdateTimers.current.delete(positionChange.id);
            }, 500);

            positionUpdateTimers.current.set(positionChange.id, newTimer);
          }
        }
      });
    },
    [onNodesChange, currentAgentProject, updateNodePositionAndPersist]
  );

  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  // Auto-layout function triggered by canvas container
  const handleAutoLayout = useCallback(() => {
    if (!currentAgentProject || !currentAgentProject.agentNodes) return;

    // Extract agents from agent nodes
    const agents = currentAgentProject.agentNodes
      .filter((agentNode) => agentNode.type === "agent")
      .map((agentNode) => ({
        ...agentNode.data.agent,
        position: agentNode.position,
      }));

    // Get connections
    const connections = currentAgentProject.connections || [];

    // Force new layout calculation
    const { nodes: layoutNodes, edges: layoutEdges } = calculateAgentLayout(
      agents,
      connections,
      {
        ...defaultAgentLayoutOptions,
        force: true, // Force new layout
      }
    );

    // Update nodes with new positions
    const newNodes: Node[] = layoutNodes.map((node) => ({
      id: node.id,
      type: "agent",
      position: node.position,
      data: node.data as Record<string, unknown>,
    }));

    // Update edges
    const newEdges: Edge[] = layoutEdges.map((edge) => ({
      id: edge.id,
      source: edge.source,
      target: edge.target,
      type: "default",
      label: edge.label,
      data: edge.data as Record<string, unknown>,
      style: edge.style,
    }));

    setNodes(newNodes);
    setEdges(newEdges);

    // First update local store immediately to prevent race conditions
    if (currentAgentProject) {
      const updatedAgentNodes = currentAgentProject.agentNodes.map((node) => {
        const layoutNode = layoutNodes.find(ln => ln.id === node.id);
        return layoutNode ? { ...node, position: layoutNode.position } : node;
      });

      // Update local store immediately
      setAgentProjects(
        useAgentProjectStore.getState().agentProjects.map((project) =>
          project.id === currentAgentProject.id
            ? { ...project, agentNodes: updatedAgentNodes }
            : project
        )
      );

      // Also update current project
      setCurrentAgentProject({
        ...currentAgentProject,
        agentNodes: updatedAgentNodes
      });

    }

    // Then persist to Convex
    Promise.all(
      layoutNodes.map((node) =>
        updateNodePositionAndPersist(node.id, node.position)
      )
    ).catch((error) => {
      console.error("Error persisting agent positions:", error);
    });
  }, [currentAgentProject, setNodes, setEdges, updateNodePositionAndPersist, setAgentProjects, setCurrentAgentProject]);

  // Expose auto-layout function to parent
  useEffect(() => {
    if (onAutoLayoutRequest) {
      // Store the handler so it can be called from parent
      (window as any).triggerAgentAutoLayout = handleAutoLayout;
    }
    return () => {
      delete (window as any).triggerAgentAutoLayout;
    };
  }, [handleAutoLayout, onAutoLayoutRequest]);

  // Convert agent project data to ReactFlow nodes and edges with auto-layout
  const { reactFlowNodes, reactFlowEdges } = useMemo(() => {
    if (!currentAgentProject || !currentAgentProject.agentNodes) {
      return { reactFlowNodes: [], reactFlowEdges: [] };
    }

    // Extract agents from agent nodes
    const agents = currentAgentProject.agentNodes
      .filter((agentNode) => agentNode.type === "agent")
      .map((agentNode) => ({
        ...agentNode.data.agent,
        position: agentNode.position,
      }));

    // Use auto-layout if no positions are set or if force layout is requested
    const hasPositions = agents.some(
      (agent) =>
        agent.position && agent.position.x !== 0 && agent.position.y !== 0
    );

    if (!hasPositions || onAutoLayoutRequest) {
      // Use hierarchical layout if connections exist, otherwise use grid layout
      const connections = currentAgentProject.connections || [];

      let layoutResult;
      if (connections.length > 0) {
        layoutResult = calculateAgentLayout(agents, connections, {
          ...defaultAgentLayoutOptions,
          force: !hasPositions,
        });
      } else {
        layoutResult = calculateGridLayout(agents, {
          force: !hasPositions,
        });
      }

      return {
        reactFlowNodes: layoutResult.nodes,
        reactFlowEdges: layoutResult.edges,
      };
    }

    // Use existing positions
    const reactFlowNodes: Node[] = currentAgentProject.agentNodes
      .filter((agentNode) => agentNode.type === "agent")
      .map((agentNode) => {
        console.log(`USING existing position - Agent ${agentNode.id}:`, agentNode.position);
        return {
          id: agentNode.id,
          type: "agent",
          position: agentNode.position,
          data: agentNode.data,
        };
      });

    // Convert connections to edges
    const reactFlowEdges: Edge[] = currentAgentProject.connections
      ? currentAgentProject.connections
          .filter((connection) => {
            const validTypes = ["hierarchy", "communication", "collaboration"];
            return validTypes.includes(connection.type);
          })
          .map((connection) => ({
            id: connection.id,
            source: connection.source,
            target: connection.target,
            type: "default",
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
  }, [currentAgentProject, onAutoLayoutRequest]);

  // Update nodes and edges when agent project changes
  useEffect(() => {
    // Additional safety check - ensure all nodes are agent type
    const safeNodes = reactFlowNodes.filter((node) => node.type === "agent");
    if (defaultViewport) {
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
    }

    // Additional safety check - ensure all edges use default type
    const safeEdges = reactFlowEdges.map((edge) => ({
      ...edge,
      type: "default", // Force all edges to use default type
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

  // Show loading state
  if (agentDataLoading) {
    return (
      <div
        className={`h-full w-full ${className} flex items-center justify-center`}
      >
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
        onNodesChange={handleNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onSelectionChange={handleSelectionChange}
        onMove={onMove} // Pass the onMove prop here
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
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
