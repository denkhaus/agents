/**
 * Task Node Component
 * Custom ReactFlow node for displaying tasks
 */

import React from "react";
import { Handle, Position, NodeProps } from "@xyflow/react";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { MarkdownRenderer } from "@/components/ui/markdown-renderer";
import { TaskState } from "@/types/task.types";
import { Clock, User } from "lucide-react";
import { cn } from "@/lib/utils";
import { TaskNodeData } from "@/types/reactflow.types";
import { TaskPropertyPanel } from "@/components/property-panels";
import { taskStateConfig } from "@/config/task-state-config";
import type {
  PropertyPanelNode,
  PropertyInfo,
  PropertyUpdateCallback,
  Task,
} from "@/types";

// Create a class that implements the PropertyPanelNode interface
class TaskNodeClass implements PropertyPanelNode {
  constructor(
    public id: string,
    public type: string,
    private task: Task,
    private onUpdate: PropertyUpdateCallback
  ) {}

  getPropertyInfo(): PropertyInfo {
    return {
      id: this.id,
      type: this.type,
      title: this.task.title,
      description: this.task.description,
      component: (
        <TaskPropertyPanel task={this.task} onUpdate={this.onUpdate} />
      ),
    };
  }
}

export const TaskNode: React.FC<NodeProps> = ({ data, selected }) => {
  const { task, isHighlighted, showDetails } = data as TaskNodeData;
  const config = taskStateConfig[task.state as TaskState];
  const Icon = config.icon;

  // Create a property panel node instance - this will be used by the sidebar
  const createPropertyPanelNode = (
    onUpdate: PropertyUpdateCallback
  ): PropertyPanelNode => {
    return new TaskNodeClass(task.id, "Task", task, onUpdate);
  };

  // Store the factory function on the node data for the sidebar to access
  React.useEffect(() => {
    (data as any).getPropertyPanelNode = createPropertyPanelNode;
  }, [task, data]);

  return (
    <div className="task-node">
      {/* Input Handle */}
      <Handle
        type="target"
        position={Position.Left}
        className="w-3 h-3 border-2 border-background"
        style={{ background: "#6b7280" }}
      />

      <Card
        className={cn(
          "w-80 transition-all duration-200 cursor-pointer", // Removed border here
          config.color,
          config.darkColor,
          selected && "ring-2 ring-primary ring-offset-2",
          isHighlighted && "shadow-lg scale-105"
        )}
      >
        <CardHeader
          className={cn(
            "pb-6 border-x border-t border-b bg-background/80 dark:bg-background/80 rounded-t-lg",
            config.color
              .split(" ")
              .filter((c) => c.startsWith("border-"))
              .join(" "),
            config.darkColor
              .split(" ")
              .filter((c) => c.startsWith("dark:border-"))
              .join(" ")
          )}
        >
          <div className="flex items-start justify-between gap-2">
            <div className="flex-1 min-w-0">
              <h3 className="font-medium text-sm leading-tight truncate">
                {task.title}
              </h3>
              <div className="flex items-center gap-2 mt-1">
                <Badge variant={config.badge as any} className="text-xs">
                  <Icon className="w-3 h-3 mr-1" />
                  {task.state.replace("-", " ")}
                </Badge>
                <Badge variant="outline" className="text-xs">
                  C{task.complexity}
                </Badge>
              </div>
            </div>

            {task.assignedAgent && (
              <Avatar className="w-6 h-6">
                <AvatarFallback className="text-xs">
                  {task.assignedAgent.slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
            )}
          </div>
        </CardHeader>

        <CardContent className="pt-6 border-x border-b border-border/20 rounded-b-lg">
          {/* Description Preview */}
          <div className="mb-3">
            <MarkdownRenderer
              content={task.description}
              maxLength={showDetails ? undefined : 100}
              showPreview={!showDetails}
              className="text-xs leading-relaxed"
            />
          </div>

          {/* Task Metadata */}
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <div className="flex items-center gap-3">
              {task.estimate && (
                <div className="flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  <span>{Math.round(task.estimate / 60)}h</span>
                </div>
              )}

              {task.dependencies.length > 0 && (
                <div className="flex items-center gap-1">
                  <span>Deps: {task.dependencies.length}</span>
                </div>
              )}
            </div>

            <div className="flex items-center gap-1">
              <span>D{task.depth}</span>
            </div>
          </div>

          {/* Expanded Details */}
          {showDetails && (
            <div className="mt-3 pt-3 border-t border-border/50">
              <div className="space-y-2 text-xs">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Created:</span>
                  <span>{task.createdAt.toLocaleDateString()}</span>
                </div>

                {task.completedAt && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Completed:</span>
                    <span>{task.completedAt.toLocaleDateString()}</span>
                  </div>
                )}

                {task.assignedAgent && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Assigned:</span>
                    <span className="flex items-center gap-1">
                      <User className="w-3 h-3" />
                      Agent
                    </span>
                  </div>
                )}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Output Handle */}
      <Handle
        type="source"
        position={Position.Right}
        className="w-3 h-3 border-2 border-background"
        style={{ background: "#6b7280" }}
      />
    </div>
  );
};
