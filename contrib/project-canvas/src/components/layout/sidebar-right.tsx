/**
 * Right Sidebar Component
 * Displays properties of selected workspace items
 */

import React from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { useUIStore, useProjectStore, useTaskStore } from "@/stores";
import { useRealTimeData } from "@/hooks/use-real-time-data";
import { cn } from "@/lib/utils";
import { X, Info, Calendar, User, Tag } from "lucide-react";
import { Button } from "@/components/ui/button";
import type { PropertyPanelNode, PropertyUpdateCallback } from "@/types";

interface PropertyItem {
  label: string;
  value: string;
  type?: "badge" | "status" | "progress";
  icon?: React.ComponentType<{ className?: string }>;
}

export const SidebarRight: React.FC = () => {
  const { 
    rightSidebarCollapsed, 
    currentWorkspace, 
    toggleRightSidebar,
    selection 
  } = useUIStore();
  const { currentProject, updateProject } = useProjectStore();
  const { tasks, updateTask } = useTaskStore();
  const { agents, nodes } = useRealTimeData();

  // Property update callback
  const handlePropertyUpdate: PropertyUpdateCallback = React.useCallback(
    (nodeId, updates) => {
      // Find the node type and update accordingly
      const task = tasks.find(t => t.id === nodeId);
      const project = currentProject && currentProject.id === nodeId ? currentProject : null;
      
      if (task) {
        updateTask(nodeId, updates);
      } else if (project) {
        updateProject(nodeId, updates);
      }
    },
    [tasks, currentProject, updateTask, updateProject]
  );

  const getSelectedNodePropertyPanel = (): PropertyPanelNode | null => {
    // Get the first selected node
    const selectedNodeId = selection.selectedNodes[0];
    if (!selectedNodeId) return null;

    // Find the node in the ReactFlow nodes (this would come from the canvas)
    // For now, we'll reconstruct the property panel from our data
    const task = tasks.find(t => t.id === selectedNodeId);
    const project = currentProject && currentProject.id === selectedNodeId ? currentProject : null;

    if (task) {
      return {
        id: task.id,
        type: "Task",
        getPropertyInfo: () => ({
          id: task.id,
          type: "Task",
          title: task.title,
          description: task.description,
          component: React.createElement(
            require("@/components/property-panels").TaskPropertyPanel,
            { task, onUpdate: handlePropertyUpdate }
          )
        })
      };
    }

    if (project) {
      return {
        id: project.id,
        type: "Project", 
        getPropertyInfo: () => ({
          id: project.id,
          type: "Project",
          title: project.title,
          description: project.description,
          component: React.createElement(
            require("@/components/property-panels").ProjectPropertyPanel,
            { project, onUpdate: handlePropertyUpdate }
          )
        })
      };
    }

    return null;
  };

  const getSelectedItemProperties = () => {
    // If we're in projects workspace and have a selected node, show its property panel
    if (currentWorkspace === "projects") {
      const propertyPanelNode = getSelectedNodePropertyPanel();
      if (propertyPanelNode) {
        return propertyPanelNode.getPropertyInfo();
      }
      // Fallback to current project if no node selected
      return currentProject ? getProjectProperties(currentProject) : null;
    }

    // Handle other workspaces
    switch (currentWorkspace) {
      case "agents":
        return agents.length > 0 ? getAgentProperties(agents[0]) : null;
      case "settings":
        return getSettingsProperties();
      default:
        return null;
    }
  };

  const getProjectProperties = (project: any): { title: string; type: string; properties: PropertyItem[] } => ({
    title: project.title,
    type: "Project",
    properties: [
      { label: "Description", value: project.description || "No description" },
      { label: "Status", value: project.status || "Active", type: "badge" },
      { label: "Created", value: new Date(project.createdAt).toLocaleDateString(), icon: Calendar },
      { label: "Tasks", value: `${project.tasks?.length || 0} tasks`, icon: Tag },
      { label: "Priority", value: project.priority || "Medium", type: "badge" },
    ]
  });

  const getAgentProperties = (agent: any): { title: string; type: string; properties: PropertyItem[] } => ({
    title: agent.name || "Agent",
    type: "Agent",
    properties: [
      { label: "Status", value: agent.status || "online", type: "status" },
      { label: "Role", value: agent.role || "Developer", icon: User },
      { label: "Current Task", value: agent.currentTask || "Idle" },
      { label: "Efficiency", value: `${agent.efficiency || 85}%`, type: "progress" },
    ]
  });

  const getSettingsProperties = (): { title: string; type: string; properties: PropertyItem[] } => ({
    title: "Application Settings",
    type: "Settings",
    properties: [
      { label: "Theme", value: "Light Mode" },
      { label: "Language", value: "English" },
      { label: "Auto-save", value: "Enabled", type: "badge" },
      { label: "Notifications", value: "All enabled" },
    ]
  });

  const selectedItem = getSelectedItemProperties();

  return (
    <aside
      className={cn(
        "fixed right-0 top-14 h-[calc(100vh-3.5rem)] bg-background border-l border-border transition-all duration-300 ease-in-out z-40",
        rightSidebarCollapsed ? "w-0 opacity-0" : "w-80 opacity-100"
      )}
    >
      {!rightSidebarCollapsed && (
        <div className="flex h-full flex-col">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b">
            <div className="flex items-center gap-2">
              <Info className="h-4 w-4" />
              <h2 className="text-sm font-semibold">Properties</h2>
            </div>
            <Button
              variant="ghost"
              size="sm"
              className="h-6 w-6 p-0"
              onClick={toggleRightSidebar}
            >
              <X className="h-3 w-3" />
            </Button>
          </div>

          {/* Content */}
          <ScrollArea className="flex-1">
            <div className="p-4">
              {selectedItem ? (
                <div>
                  {/* Check if this is a new property panel component */}
                  {selectedItem.component ? (
                    // Render the unified property panel component
                    selectedItem.component
                  ) : (
                    // Fallback to legacy property display
                    <div className="space-y-4">
                      {/* Item Header */}
                      <div>
                        <h3 className="font-medium text-sm mb-1">{selectedItem.title}</h3>
                        <Badge variant="outline" className="text-xs">
                          {selectedItem.type}
                        </Badge>
                      </div>

                      <Separator />

                      {/* Properties List */}
                      <div className="space-y-3">
                        {selectedItem.properties?.map((property, index) => {
                          const IconComponent = property.icon;
                          return (
                            <div key={index} className="space-y-1">
                              <div className="flex items-center gap-2">
                                {IconComponent && (
                                  <IconComponent className="h-3 w-3 text-muted-foreground" />
                                )}
                                <span className="text-xs font-medium text-muted-foreground">
                                  {property.label}
                                </span>
                              </div>
                            <div className="ml-5">
                              {property.type === "badge" ? (
                                <Badge variant="secondary" className="text-xs">
                                  {property.value}
                                </Badge>
                              ) : property.type === "status" ? (
                                <div className="flex items-center gap-2">
                                  <div
                                    className={cn(
                                      "h-2 w-2 rounded-full",
                                      property.value === "online" && "bg-green-500",
                                      property.value === "busy" && "bg-yellow-500",
                                      property.value === "idle" && "bg-gray-400"
                                    )}
                                  />
                                  <span className="text-xs capitalize">{property.value}</span>
                                </div>
                              ) : property.type === "progress" ? (
                                <div className="space-y-1">
                                  <span className="text-xs">{property.value}</span>
                                  <div className="w-full bg-muted rounded-full h-1">
                                    <div
                                      className="bg-primary h-1 rounded-full"
                                      style={{ width: property.value }}
                                    />
                                  </div>
                                </div>
                              ) : (
                                <span className="text-xs">{property.value}</span>
                              )}
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  )}
                </div>
              ) : (
                <div className="text-center text-muted-foreground py-8">
                  <Info className="h-8 w-8 mx-auto mb-2 opacity-50" />
                  <p className="text-sm">No item selected</p>
                  <p className="text-xs mt-1">
                    Select a node to view its properties
                  </p>
                </div>
              )}
            </div>
          </ScrollArea>
        </div>
      )}
    </aside>
  );
};