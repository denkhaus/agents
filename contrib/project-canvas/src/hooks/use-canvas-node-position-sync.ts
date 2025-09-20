import { useCallback } from "react";
import { useConvexMutations } from "./use-convex-data";
import { useAgentProjectStore, useTaskStore } from "@/stores";
import { UUID } from "@/types";

interface Position {
  x: number;
  y: number;
}

type CanvasType = "agentProjects" | "projects" | "tasks";

export const useCanvasNodePositionSync = (
  canvasType: CanvasType,
  projectId: UUID | null
) => {
  const { updateAgentNodes, updateTaskPosition, updateProjectPosition } =
    useConvexMutations();
  const { updateTaskPosition: updateTaskPositionLocal } = useTaskStore();
  const { currentAgentProject, setAgentProjects } = useAgentProjectStore(); // To get current agent project's nodes for updateAgentNodes and update local state

  const updateNodePositionAndPersist = useCallback(
    async (nodeId: UUID, position: Position) => {
      if (!projectId) {
        console.warn("No project ID provided for node position sync.");
        return;
      }

      try {
        switch (canvasType) {
          case "agentProjects":
            // Optimistic update in local Zustand store
            if (currentAgentProject) {
              const updatedAgentNodes = currentAgentProject.agentNodes.map(
                (node) => (node.id === nodeId ? { ...node, position } : node)
              );

              setAgentProjects(
                useAgentProjectStore
                  .getState()
                  .agentProjects.map((project) =>
                    project.id === projectId
                      ? { ...project, agentNodes: updatedAgentNodes }
                      : project
                  )
              );

              // Prepare agentNodes for Convex by ensuring only expected fields are sent
              const agentNodesForConvex = updatedAgentNodes.map((node) => ({
                id: node.id,
                type: node.type,
                position: node.position,
                data: {
                  isSelected: node.data?.isSelected || false,
                },
              }));

              await updateAgentNodes({
                id: projectId,
                agentNodes: agentNodesForConvex,
              });
            }
            break;
          case "projects":
            // For project nodes, we only update in Convex, no local Zustand store for individual node positions
            await updateProjectPosition({
              id: projectId, // Corrected: use projectId for project canvas
              positionX: position.x,
              positionY: position.y,
            });
            break;
          case "tasks":
            // Optimistic update in local Zustand store
            updateTaskPositionLocal(nodeId, position);
            await updateTaskPosition({
              id: nodeId,
              positionX: position.x,
              positionY: position.y,
            });
            break;
          default:
            console.warn(`Unknown canvas type: ${canvasType}`);
            break;
        }
      } catch (error) {
        console.error(`Failed to sync ${canvasType} node position:`, error);
        // TODO: Implement retry logic or rollback if necessary
      }
    },
    [
      canvasType,
      projectId,
      updateAgentNodes,
      updateTaskPosition,
      updateProjectPosition,
      setAgentProjects, // Add setAgentProjects to dependencies
      updateTaskPositionLocal,
      currentAgentProject, // Dependency for currentAgentProject
    ]
  );

  return { updateNodePositionAndPersist };
};
