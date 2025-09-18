/**
 * Agent Projects Panel Component
 * Displays agent projects list and management controls for the left sidebar
 */

import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Plus,
  Search,
  MoreVertical,
  Users,
  Copy,
  Trash2,
  Edit,
  FolderOpen,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useAgentProjectStore } from "@/stores";
import { useAgentProjectData } from "@/hooks/use-agent-project-data";
import { AgentProject } from "@/types/agent.types";

interface AgentProjectsPanelProps {
  className?: string;
}

export const AgentsWorkspace: React.FC<AgentProjectsPanelProps> = ({
  className,
}) => {
  // Initialize data
  useAgentProjectData();

  const {
    currentAgentProject,
    setCurrentAgentProject,
    addAgentProject,
    updateAgentProject,
    deleteAgentProject,
    duplicateAgentProject,
    filter,
    setFilter,
    getFilteredAgentProjects,
  } = useAgentProjectStore();

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<AgentProject | null>(
    null
  );
  const [newProjectName, setNewProjectName] = useState("");
  const [newProjectDescription, setNewProjectDescription] = useState("");

  const filteredProjects = getFilteredAgentProjects();

  const handleCreateProject = () => {
    if (newProjectName.trim()) {
      const newProject: AgentProject = {
        id: crypto.randomUUID(),
        name: newProjectName.trim(),
        description: newProjectDescription.trim(),
        agents: [],
        agentNodes: [],
        connections: [],
        createdAt: new Date(),
        updatedAt: new Date(),
      };

      addAgentProject(newProject);
      setCurrentAgentProject(newProject);
      setNewProjectName("");
      setNewProjectDescription("");
      setIsCreateDialogOpen(false);
    }
  };

  const handleEditProject = () => {
    if (editingProject && newProjectName.trim()) {
      updateAgentProject(editingProject.id, {
        name: newProjectName.trim(),
        description: newProjectDescription.trim(),
      });
      setNewProjectName("");
      setNewProjectDescription("");
      setEditingProject(null);
      setIsEditDialogOpen(false);
    }
  };

  const handleDeleteProject = (project: AgentProject) => {
    deleteAgentProject(project.id);
    if (currentAgentProject?.id === project.id) {
      setCurrentAgentProject(null);
    }
  };

  const handleDuplicateProject = (project: AgentProject) => {
    duplicateAgentProject(project.id);
  };

  const openEditDialog = (project: AgentProject) => {
    setEditingProject(project);
    setNewProjectName(project.name);
    setNewProjectDescription(project.description);
    setIsEditDialogOpen(true);
  };

  return (
    <div className={cn("space-y-2", className)}>
      {/* Header */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h3 className="text-sm font-medium">Agent Projects</h3>
          <Dialog
            open={isCreateDialogOpen}
            onOpenChange={setIsCreateDialogOpen}
          >
            <DialogTrigger asChild>
              <Button size="sm" className="h-8 w-8 p-0">
                <Plus className="h-4 w-4" />
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Create Agent Project</DialogTitle>
                <DialogDescription>
                  Create a new agent project to organize your agent
                  configurations.
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="project-name">Project Name</Label>
                  <Input
                    id="project-name"
                    value={newProjectName}
                    onChange={(e) => setNewProjectName(e.target.value)}
                    placeholder="Enter project name"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="project-description">Description</Label>
                  <Textarea
                    id="project-description"
                    value={newProjectDescription}
                    onChange={(e) => setNewProjectDescription(e.target.value)}
                    placeholder="Enter project description"
                    rows={3}
                  />
                </div>
              </div>
              <DialogFooter>
                <Button
                  variant="outline"
                  onClick={() => setIsCreateDialogOpen(false)}
                >
                  Cancel
                </Button>
                <Button onClick={handleCreateProject}>Create Project</Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>

        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search projects..."
            value={filter.searchTerm || ""}
            onChange={(e) => setFilter({ searchTerm: e.target.value })}
            className="pl-9"
          />
        </div>
      </div>

      {/* Projects List */}
      <div className="space-y-1">
        {filteredProjects.length === 0 ? (
          <div className="text-center py-4">
            <FolderOpen className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
            <p className="text-xs text-muted-foreground mb-2">
              No agent projects
            </p>
            <Button
              size="sm"
              onClick={() => setIsCreateDialogOpen(true)}
              className="h-7 text-xs"
            >
              <Plus className="h-3 w-3 mr-1" />
              Create
            </Button>
          </div>
        ) : (
          filteredProjects.map((project) => (
            <Card
              key={project.id}
              className={cn(
                "cursor-pointer transition-all duration-200 hover:shadow-md",
                currentAgentProject?.id === project.id &&
                  "ring-2 ring-blue-500 ring-offset-2"
              )}
              onClick={() => setCurrentAgentProject(project)}
            >
              <CardContent className="p-2">
                <div className="flex items-start justify-between">
                  <div className="flex-1 min-w-0">
                    <h3
                      className="font-medium text-sm truncate"
                      title={project.name}
                    >
                      {project.name}
                    </h3>
                    {project.description && (
                      <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                        {project.description}
                      </p>
                    )}
                    <div className="flex items-center gap-2 mt-2">
                      <Badge variant="secondary" className="text-xs">
                        <Users className="h-3 w-3 mr-1" />
                        {project.agentNodes.length} agents
                      </Badge>
                      {project.connections.length > 0 && (
                        <Badge variant="outline" className="text-xs">
                          {project.connections.length} connections
                        </Badge>
                      )}
                    </div>
                  </div>

                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <MoreVertical className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => openEditDialog(project)}>
                        <Edit className="h-4 w-4 mr-2" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => handleDuplicateProject(project)}
                      >
                        <Copy className="h-4 w-4 mr-2" />
                        Duplicate
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        onClick={() => handleDeleteProject(project)}
                        className="text-destructive"
                      >
                        <Trash2 className="h-4 w-4 mr-2" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>

                <div className="text-xs text-muted-foreground mt-2">
                  Updated {new Date(project.updatedAt).toLocaleDateString()}
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Agent Project</DialogTitle>
            <DialogDescription>
              Update the project name and description.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="edit-project-name">Project Name</Label>
              <Input
                id="edit-project-name"
                value={newProjectName}
                onChange={(e) => setNewProjectName(e.target.value)}
                placeholder="Enter project name"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="edit-project-description">Description</Label>
              <Textarea
                id="edit-project-description"
                value={newProjectDescription}
                onChange={(e) => setNewProjectDescription(e.target.value)}
                placeholder="Enter project description"
                rows={3}
              />
            </div>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setIsEditDialogOpen(false)}
            >
              Cancel
            </Button>
            <Button onClick={handleEditProject}>Save Changes</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};
