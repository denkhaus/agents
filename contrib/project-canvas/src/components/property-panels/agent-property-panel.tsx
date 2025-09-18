/**
 * Agent Property Panel Component
 * Displays and allows editing of agent properties
 */

import React, { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  User,
  Settings,
  Activity,
  Clock,
  Zap,
  Plus,
  X,
  Save,
  RotateCcw,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Agent, AgentRoleType, AgentStatusType, AgentRole, AgentStatus } from "@/types/agent.types";

interface AgentPropertyPanelProps {
  agent: Agent;
  onUpdate: (updates: Partial<Agent>) => void;
  onClose?: () => void;
}

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

const getAgentInitials = (name: string): string => {
  return name
    .split(" ")
    .map((word) => word.charAt(0))
    .join("")
    .toUpperCase()
    .slice(0, 2);
};

export const AgentPropertyPanel: React.FC<AgentPropertyPanelProps> = ({
  agent,
  onUpdate,
  onClose,
}) => {
  const [editedAgent, setEditedAgent] = useState<Agent>(agent);
  const [newCapability, setNewCapability] = useState("");
  const [hasChanges, setHasChanges] = useState(false);

  const handleFieldChange = (field: keyof Agent, value: any) => {
    setEditedAgent(prev => ({ ...prev, [field]: value }));
    setHasChanges(true);
  };

  const handleSave = () => {
    onUpdate(editedAgent);
    setHasChanges(false);
  };

  const handleReset = () => {
    setEditedAgent(agent);
    setHasChanges(false);
  };

  const addCapability = () => {
    if (newCapability.trim() && !editedAgent.capabilities.includes(newCapability.trim())) {
      handleFieldChange("capabilities", [...editedAgent.capabilities, newCapability.trim()]);
      setNewCapability("");
    }
  };

  const removeCapability = (capability: string) => {
    handleFieldChange(
      "capabilities",
      editedAgent.capabilities.filter(c => c !== capability)
    );
  };

  return (
    <div className="w-full max-w-md space-y-4">
      {/* Header */}
      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg flex items-center gap-2">
              <User className="h-5 w-5" />
              Agent Properties
            </CardTitle>
            {onClose && (
              <Button variant="ghost" size="sm" onClick={onClose}>
                <X className="h-4 w-4" />
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Agent Avatar and Basic Info */}
          <div className="flex items-center gap-3">
            <div className="relative">
              <Avatar className="h-12 w-12">
                <AvatarFallback className="text-lg font-medium">
                  {getAgentInitials(editedAgent.name)}
                </AvatarFallback>
              </Avatar>
              <div
                className={cn(
                  "absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-white",
                  getStatusColor(editedAgent.status)
                )}
              />
            </div>
            <div className="flex-1">
              <h3 className="font-medium">{editedAgent.name}</h3>
              <p className="text-sm text-muted-foreground capitalize">
                {editedAgent.role.replace("-", " ")}
              </p>
            </div>
          </div>

          {/* Action Buttons */}
          {hasChanges && (
            <div className="flex gap-2">
              <Button onClick={handleSave} size="sm" className="flex-1">
                <Save className="h-4 w-4 mr-2" />
                Save Changes
              </Button>
              <Button onClick={handleReset} variant="outline" size="sm">
                <RotateCcw className="h-4 w-4" />
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Basic Information */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base flex items-center gap-2">
            <Settings className="h-4 w-4" />
            Basic Information
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="agent-name">Name</Label>
            <Input
              id="agent-name"
              value={editedAgent.name}
              onChange={(e) => handleFieldChange("name", e.target.value)}
              placeholder="Agent name"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="agent-role">Role</Label>
            <Select
              value={editedAgent.role}
              onValueChange={(value: AgentRoleType) => handleFieldChange("role", value)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select role" />
              </SelectTrigger>
              <SelectContent>
                {Object.values(AgentRole).map((role) => (
                  <SelectItem key={role} value={role}>
                    {role.replace("-", " ").replace(/\b\w/g, l => l.toUpperCase())}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="agent-status">Status</Label>
            <Select
              value={editedAgent.status}
              onValueChange={(value: AgentStatusType) => handleFieldChange("status", value)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select status" />
              </SelectTrigger>
              <SelectContent>
                {Object.values(AgentStatus).map((status) => (
                  <SelectItem key={status} value={status}>
                    <div className="flex items-center gap-2">
                      <div className={cn("w-2 h-2 rounded-full", getStatusColor(status))} />
                      {status.charAt(0).toUpperCase() + status.slice(1)}
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="agent-description">Description</Label>
            <Textarea
              id="agent-description"
              value={editedAgent.description}
              onChange={(e) => handleFieldChange("description", e.target.value)}
              placeholder="Agent description"
              rows={3}
            />
          </div>
        </CardContent>
      </Card>

      {/* Capabilities */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base flex items-center gap-2">
            <Zap className="h-4 w-4" />
            Capabilities
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex flex-wrap gap-2">
            {editedAgent.capabilities.map((capability, index) => (
              <Badge
                key={index}
                variant="secondary"
                className="flex items-center gap-1"
              >
                {capability}
                <button
                  onClick={() => removeCapability(capability)}
                  className="ml-1 hover:bg-destructive hover:text-destructive-foreground rounded-full p-0.5"
                >
                  <X className="h-3 w-3" />
                </button>
              </Badge>
            ))}
          </div>

          <div className="flex gap-2">
            <Input
              value={newCapability}
              onChange={(e) => setNewCapability(e.target.value)}
              placeholder="Add capability"
              onKeyPress={(e) => e.key === "Enter" && addCapability()}
            />
            <Button onClick={addCapability} size="sm" variant="outline">
              <Plus className="h-4 w-4" />
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Status & Activity */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base flex items-center gap-2">
            <Activity className="h-4 w-4" />
            Status & Activity
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <Label htmlFor="streaming">Streaming</Label>
            <Switch
              id="streaming"
              checked={editedAgent.isStreaming}
              onCheckedChange={(checked) => handleFieldChange("isStreaming", checked)}
            />
          </div>

          <Separator />

          <div className="space-y-2">
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">Current Tasks</span>
              <span className="font-medium">{editedAgent.currentTasks.length}</span>
            </div>
            
            {editedAgent.efficiency !== undefined && (
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Efficiency</span>
                <span className="font-medium">{Math.round(editedAgent.efficiency * 100)}%</span>
              </div>
            )}

            {editedAgent.lastActiveAt && (
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Last Active</span>
                <span className="font-medium">
                  {new Date(editedAgent.lastActiveAt).toLocaleString()}
                </span>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Timestamps */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base flex items-center gap-2">
            <Clock className="h-4 w-4" />
            Timestamps
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-muted-foreground">Created</span>
            <span>{new Date(editedAgent.createdAt).toLocaleString()}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Updated</span>
            <span>{new Date(editedAgent.updatedAt).toLocaleString()}</span>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};