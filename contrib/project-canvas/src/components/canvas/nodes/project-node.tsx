/**
 * Project Node Component
 * Custom ReactFlow node for displaying projects (root level)
 */

import React from "react";
import { Handle, Position, NodeProps } from "@xyflow/react";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { MarkdownRenderer } from "@/components/ui/markdown-renderer";
import { ProjectNodeData } from "@/types/reactflow.types";
import { FolderOpen, Calendar, CheckCircle2, Clock } from "lucide-react";
import { cn } from "@/lib/utils";
import { ProjectPropertyPanel } from "@/components/property-panels";
import type { PropertyPanelNode, PropertyInfo, PropertyUpdateCallback } from "@/types";

// Create a class that implements the PropertyPanelNode interface
class ProjectNodeClass implements PropertyPanelNode {
  constructor(
    public id: string,
    public type: string,
    private project: any,
    private onUpdate: PropertyUpdateCallback
  ) {}

  getPropertyInfo(): PropertyInfo {
    return {
      id: this.id,
      type: this.type,
      title: this.project.title,
      description: this.project.description,
      component: (
        <ProjectPropertyPanel
          project={this.project}
          onUpdate={this.onUpdate}
        />
      )
    };
  }
}

export const ProjectNode: React.FC<NodeProps> = ({ data, selected }) => {
  const { project, taskCount, completionRate } = data as ProjectNodeData;

  // Create a property panel node instance - this will be used by the sidebar
  const createPropertyPanelNode = (onUpdate: PropertyUpdateCallback): PropertyPanelNode => {
    return new ProjectNodeClass(project.id, "Project", project, onUpdate);
  };

  // Store the factory function on the node data for the sidebar to access
  React.useEffect(() => {
    (data as any).getPropertyPanelNode = createPropertyPanelNode;
  }, [project, data]);

  // Use completionRate in the UI
  const displayCompletionRate =
    completionRate > 0 ? `(${Math.round(completionRate)}%)` : "";

  return (
    <div className="project-node">
      <Card
        className={cn(
          "w-96 bg-gradient-to-br from-primary/5 to-primary/10 transition-all duration-200 cursor-pointer", // Removed border here
          "dark:from-primary/10 dark:to-primary/20", // Removed dark border here
          selected && "ring-2 ring-primary ring-offset-2 shadow-lg"
        )}
      >
        <CardHeader className="pb-3 border-x border-t border-b border-border/20 bg-background/80 dark:bg-background/80 rounded-t-lg">
          <div className="flex items-start gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/20">
              <FolderOpen className="h-5 w-5 text-primary" />
            </div>

            <div className="flex-1 min-w-0">
              <h2 className="font-semibold text-lg leading-tight text-foreground">
                {project.title}
              </h2>
              <div className="flex items-center gap-2 mt-1">
                <Badge variant="secondary" className="text-xs">
                  Project
                </Badge>
                <Badge variant="outline" className="text-xs">
                  {taskCount} Tasks {displayCompletionRate}
                </Badge>
              </div>
            </div>
          </div>
        </CardHeader>

        <CardContent className="space-y-4 border-x border-b border-border/40 rounded-b-lg">
          {/* Description */}
          <div>
            <MarkdownRenderer
              content={project.description}
              maxLength={150}
              showPreview={true}
              className="text-sm leading-relaxed"
            />
          </div>

          {/* Progress Section */}
          <div className="space-y-2">
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">Progress</span>
              <span className="font-medium">
                {Math.round(project.progress)}%
              </span>
            </div>
            <Progress value={project.progress} className="h-2" />
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              <span>{project.completedTasks} completed</span>
              <span>{project.totalTasks} total</span>
            </div>
          </div>

          {/* Project Metadata */}
          <div className="grid grid-cols-2 gap-4 pt-2 border-t border-border/50">
            <div className="space-y-2">
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <Calendar className="w-3 h-3" />
                <span>Created</span>
              </div>
              <p className="text-xs font-medium">
                {project.createdAt.toLocaleDateString()}
              </p>
            </div>

            <div className="space-y-2">
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <Clock className="w-3 h-3" />
                <span>Updated</span>
              </div>
              <p className="text-xs font-medium">
                {project.updatedAt.toLocaleDateString()}
              </p>
            </div>
          </div>

          {/* Quick Stats */}
          <div className="flex items-center justify-between pt-2 border-t border-border/50">
            <div className="flex items-center gap-4 text-xs text-muted-foreground">
              <div className="flex items-center gap-1">
                <CheckCircle2 className="w-3 h-3 text-green-500" />
                <span>{project.completedTasks} done</span>
              </div>
              <div className="flex items-center gap-1">
                <Clock className="w-3 h-3 text-blue-500" />
                <span>
                  {project.totalTasks - project.completedTasks} remaining
                </span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Output Handle for child tasks */}
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-4 h-4 border-2 border-background"
        style={{ background: "#8b5cf6" }}
      />
    </div>
  );
};
