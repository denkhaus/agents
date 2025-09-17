/**
 * Project Property Panel Component
 * Displays and allows editing of project properties
 */

import React from 'react';
import { BasePropertyPanel } from './base-property-panel';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { 
  Calendar, 
  BarChart3, 
  CheckCircle2, 
  Clock, 
  Folder,
  TrendingUp,
  Users,
  Target
} from 'lucide-react';
import type { Project, PropertyUpdateCallback } from '@/types';

interface ProjectPropertyPanelProps {
  project: Project;
  onUpdate: PropertyUpdateCallback;
  className?: string;
}

export const ProjectPropertyPanel: React.FC<ProjectPropertyPanelProps> = ({
  project,
  onUpdate,
  className
}) => {
  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const getProgressColor = (progress: number) => {
    if (progress >= 80) return 'text-green-600';
    if (progress >= 50) return 'text-blue-600';
    if (progress >= 25) return 'text-yellow-600';
    return 'text-red-600';
  };

  const getProgressBadgeVariant = (progress: number) => {
    if (progress >= 80) return 'default';
    if (progress >= 50) return 'secondary';
    return 'outline';
  };

  return (
    <BasePropertyPanel
      nodeId={project.id}
      nodeType="Project"
      title={project.title}
      description={project.description}
      onUpdate={onUpdate}
      className={className}
    >
      <div className="space-y-4">
        {/* Progress Overview */}
        <div>
          <div className="flex items-center gap-2 text-xs text-muted-foreground mb-3">
            <TrendingUp className="w-3 h-3" />
            <span>Progress Overview</span>
          </div>
          
          <div className="space-y-3">
            {/* Overall Progress */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-medium">Overall Progress</span>
                <Badge 
                  variant={getProgressBadgeVariant(project.progress)} 
                  className="text-xs"
                >
                  {Math.round(project.progress)}%
                </Badge>
              </div>
              <Progress 
                value={project.progress} 
                className="h-2" 
              />
            </div>

            {/* Task Completion Stats */}
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1">
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <CheckCircle2 className="w-3 h-3 text-green-500" />
                  <span>Completed</span>
                </div>
                <p className="text-sm font-semibold text-green-600">
                  {project.completedTasks}
                </p>
              </div>
              
              <div className="space-y-1">
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Target className="w-3 h-3 text-blue-500" />
                  <span>Total Tasks</span>
                </div>
                <p className="text-sm font-semibold text-blue-600">
                  {project.totalTasks}
                </p>
              </div>
            </div>

            {/* Remaining Tasks */}
            <div>
              <div className="flex items-center gap-2 text-xs text-muted-foreground mb-1">
                <Clock className="w-3 h-3" />
                <span>Remaining</span>
              </div>
              <div className="flex items-center gap-2">
                <Badge variant="outline" className="text-xs">
                  {project.totalTasks - project.completedTasks} tasks
                </Badge>
                <span className="text-xs text-muted-foreground">
                  left to complete
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Project Metadata */}
        <div className="pt-2 border-t border-border/50">
          <div className="flex items-center gap-2 text-xs text-muted-foreground mb-3">
            <BarChart3 className="w-3 h-3" />
            <span>Project Details</span>
          </div>

          <div className="space-y-3">
            {/* Creation Date */}
            <div>
              <div className="flex items-center gap-2 text-xs text-muted-foreground mb-1">
                <Calendar className="w-3 h-3" />
                <span>Created</span>
              </div>
              <p className="text-xs font-medium">
                {formatDate(project.createdAt)}
              </p>
            </div>

            {/* Last Updated */}
            <div>
              <div className="flex items-center gap-2 text-xs text-muted-foreground mb-1">
                <Clock className="w-3 h-3" />
                <span>Last Updated</span>
              </div>
              <p className="text-xs font-medium">
                {formatDate(project.updatedAt)}
              </p>
            </div>

            {/* Project Status */}
            <div>
              <div className="flex items-center gap-2 text-xs text-muted-foreground mb-1">
                <Folder className="w-3 h-3" />
                <span>Status</span>
              </div>
              <Badge 
                variant={project.progress === 100 ? 'default' : 'secondary'} 
                className="text-xs"
              >
                {project.progress === 100 ? 'Completed' : 'In Progress'}
              </Badge>
            </div>
          </div>
        </div>

        {/* Quick Actions Info */}
        <div className="pt-2 border-t border-border/50">
          <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
            <Users className="w-3 h-3" />
            <span>Quick Stats</span>
          </div>
          
          <div className="grid grid-cols-2 gap-2 text-xs">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Completion Rate</span>
              <span className={`font-medium ${getProgressColor(project.progress)}`}>
                {project.totalTasks > 0 
                  ? Math.round((project.completedTasks / project.totalTasks) * 100)
                  : 0}%
              </span>
            </div>
            
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Tasks Ratio</span>
              <span className="font-medium">
                {project.completedTasks}/{project.totalTasks}
              </span>
            </div>
          </div>
        </div>
      </div>
    </BasePropertyPanel>
  );
};