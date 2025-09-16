/**
 * ReactFlow Canvas Component
 * Main visualization component for projects and tasks
 */

import React, { useCallback, useMemo } from 'react';
import ReactFlow, {
  useNodesState,
  useEdgesState,
  useReactFlow,
  ReactFlowProvider,
  Background,
  Controls,
  MiniMap,
  ConnectionMode,
  OnConnect,
  OnNodesChange,
  Node,
  Edge,
} from 'reactflow';
import 'reactflow/dist/style.css';

import { TaskNode } from './nodes/task-node';
import { ProjectNode } from './nodes/project-node';
import { DependencyEdge } from './edges/dependency-edge';
import { useProjectStore, useTaskStore } from '@/stores';
import { calculateTaskLayout } from '@/utils/layout';
import { TaskNodeData } from '@/types/reactflow.types';

// Define custom node types
const nodeTypes = {
  task: TaskNode,
  project: ProjectNode,
};

// Define custom edge types
const edgeTypes = {
  dependency: DependencyEdge,
};

interface ReactFlowCanvasProps {
  className?: string;
}

const ReactFlowCanvasInner: React.FC<ReactFlowCanvasProps> = ({ className }) => {
  const { currentProject } = useProjectStore();
  const { tasks, updateTaskPosition } = useTaskStore();
  const { fitView } = useReactFlow();

  // Generate nodes and edges from current project data
  const { initialNodes, initialEdges } = useMemo(() => {
    if (!currentProject || tasks.length === 0) {
      return { initialNodes: [], initialEdges: [] };
    }

    const projectTasks = tasks.filter(task => task.projectId === currentProject.id);
    const { nodes, edges } = calculateTaskLayout(projectTasks);

    return {
      initialNodes: nodes,
      initialEdges: edges
    };
  }, [currentProject, tasks]);

  const [nodes, , onNodesChange] = useNodesState(initialNodes as Node[]);
  const [edges, , onEdgesChange] = useEdgesState(initialEdges as Edge[]);

  // Handle new connections (for future dependency management)
  const onConnect: OnConnect = useCallback(
    (params) => console.log('Connection made:', params),
    []
  );

  // Handle node position changes
  const handleNodesChange: OnNodesChange = useCallback(
    (changes) => {
      onNodesChange(changes);
      
      // Update task positions in store for persistence
      changes.forEach((change) => {
        if (change.type === 'position' && change.position) {
          updateTaskPosition(change.id, change.position);
        }
      });
    },
    [onNodesChange, updateTaskPosition]
  );

  // Auto-fit view when data changes
  React.useEffect(() => {
    if (nodes.length > 0) {
      setTimeout(() => fitView({ padding: 0.2 }), 100);
    }
  }, [nodes.length, fitView, fitView]);

  return (
    <div className={`h-full w-full ${className}`}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={handleNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
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
        <Controls 
          className="bg-background border border-border rounded-md shadow-sm"
          showInteractive={false}
        />
        <MiniMap 
          className="bg-background border border-border rounded-md"
          maskColor="rgba(0, 0, 0, 0.1)"
          nodeColor={(node) => {
            if (node.type === 'task') {
              const taskNodeData = node.data as TaskNodeData;
              switch (taskNodeData.task.state) {
                case 'completed': return '#22c55e';
                case 'in-progress': return '#3b82f6';
                case 'blocked': return '#ef4444';
                case 'cancelled': return '#6b7280';
                default: return '#94a3b8';
              }
            }
            return '#8b5cf6';
          }}
        />
      </ReactFlow>
    </div>
  );
};

// Wrapper component with ReactFlowProvider
export const ReactFlowCanvas: React.FC<ReactFlowCanvasProps> = (props) => {
  return (
    <ReactFlowProvider>
      <ReactFlowCanvasInner {...props} />
    </ReactFlowProvider>
  );
};