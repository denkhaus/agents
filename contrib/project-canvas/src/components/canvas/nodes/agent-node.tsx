/**
 * Agent Node Component
 * Displays agent information in the ReactFlow canvas
 */

import React from "react";
import { Handle, Position, NodeProps } from "@xyflow/react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";
import { Agent, AgentStatusType } from "@/types/agent.types";

export interface AgentNodeData {
  agent: Agent;
  isSelected?: boolean;
}

export type AgentNodeProps = NodeProps;

const getStatusColor = (status: AgentStatusType): string => {
  switch (status) {
    case "online":
      return "bg-green-500";
    case "busy":
      return "bg-yellow-500";
    case "idle":
      return "bg-gray-400";
    case "offline":
      return "bg-red-500";
    default:
      return "bg-gray-400";
  }
};

const getRoleColor = (role: string): string => {
  switch (role) {
    case "supervisor":
      return "bg-purple-100 text-purple-800 border-purple-200";
    case "project-manager":
      return "bg-blue-100 text-blue-800 border-blue-200";
    case "coder":
      return "bg-green-100 text-green-800 border-green-200";
    case "researcher":
      return "bg-orange-100 text-orange-800 border-orange-200";
    case "qa-engineer":
      return "bg-red-100 text-red-800 border-red-200";
    case "devops":
      return "bg-indigo-100 text-indigo-800 border-indigo-200";
    case "designer":
      return "bg-pink-100 text-pink-800 border-pink-200";
    default:
      return "bg-gray-100 text-gray-800 border-gray-200";
  }
};

const getAgentInitials = (name: string): string => {
  if (!name) return "AG";
  
  return name
    .split(" ")
    .map((word) => word.charAt(0))
    .join("")
    .toUpperCase()
    .slice(0, 2);
};

export const AgentNode: React.FC<AgentNodeProps> = ({ data, selected }) => {
  const { agent, isSelected } = data as unknown as AgentNodeData;
  
  // Fallback for missing agent data
  if (!agent) {
    return (
      <div className="agent-node">
        <Card className="min-w-[200px] max-w-[280px]">
          <CardContent className="p-4">
            <div className="flex items-center justify-center h-20">
              <span className="text-muted-foreground">Missing agent data</span>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="agent-node">
      {/* Input Handle */}
      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 !bg-blue-500 !border-2 !border-white"
      />

      <Card
        className={cn(
          "min-w-[200px] max-w-[280px] transition-all duration-200 hover:shadow-md",
          (selected || isSelected) && "ring-2 ring-blue-500 ring-offset-2",
          agent.status === "offline" && "opacity-60"
        )}
      >
        <CardContent className="p-4">
          {/* Header */}
          <div className="flex items-start gap-3 mb-3">
            <div className="relative">
              <Avatar className="h-10 w-10">
                <AvatarFallback className="text-sm font-medium">
                  {getAgentInitials(agent.name || "Unknown Agent")}
                </AvatarFallback>
              </Avatar>
              {/* Status Indicator */}
              <div
                className={cn(
                  "absolute -bottom-1 -right-1 w-3 h-3 rounded-full border-2 border-white",
                  getStatusColor(agent.status)
                )}
              />
            </div>
            
            <div className="flex-1 min-w-0">
              <h3 className="font-medium text-sm truncate" title={agent.name || "Unknown Agent"}>
                {agent.name || "Unknown Agent"}
              </h3>
              <Badge
                variant="outline"
                className={cn("text-xs mt-1", getRoleColor(agent.role))}
              >
                {agent.role ? agent.role.replace("-", " ") : "Unknown Role"}
              </Badge>
            </div>
          </div>

          {/* Description */}
          {agent.description && (
            <p className="text-xs text-muted-foreground mb-3 line-clamp-2">
              {agent.description}
            </p>
          )}

          {/* Capabilities */}
          {agent.capabilities && agent.capabilities.length > 0 && (
            <div className="mb-3">
              <div className="flex flex-wrap gap-1">
                {agent.capabilities.slice(0, 3).map((capability: string, index: number) => (
                  <Badge
                    key={index}
                    variant="secondary"
                    className="text-xs px-2 py-0.5"
                  >
                    {capability}
                  </Badge>
                ))}
                {agent.capabilities.length > 3 && (
                  <Badge variant="secondary" className="text-xs px-2 py-0.5">
                    +{agent.capabilities.length - 3}
                  </Badge>
                )}
              </div>
            </div>
          )}

          {/* Status Info */}
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span className="capitalize">{agent.status || "unknown"}</span>
            {agent.currentTasks && agent.currentTasks.length > 0 && (
              <span>{agent.currentTasks.length} tasks</span>
            )}
            {agent.efficiency && (
              <span>{Math.round(agent.efficiency * 100)}% eff.</span>
            )}
          </div>

          {/* Streaming Indicator */}
          {agent.isStreaming && (
            <div className="mt-2 flex items-center gap-1">
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse" />
              <span className="text-xs text-blue-600">Streaming</span>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Output Handle */}
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 !bg-blue-500 !border-2 !border-white"
      />
    </div>
  );
};