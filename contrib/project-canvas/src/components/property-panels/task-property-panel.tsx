/**
 * Task Property Panel Component
 * Displays and allows editing of task properties
 */

import React from 'react';
import { BasePropertyPanel } from './base-property-panel';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { 
  Clock, 
  User, 
  Calendar, 
  BarChart3, 
  GitBranch, 
  Target,
  CheckCircle2,
  Play,
  AlertCircle,
  X as XIcon
} from 'lucide-react';
import { cn } from '@/lib/utils';
import type { Task, TaskState, PropertyUpdateCallback } from '@/types';

interface TaskPropertyPanelProps {
  task: Task;
  onUpdate: PropertyUpdateCallback;
  className?: string;
}

const taskStateConfig = {
  pending: {
    icon: Clock,
    color: 'bg-slate-100 text-slate-700 border-slate-300',
    darkColor: 'dark:bg-slate-800 dark:text-slate-300 dark:border-slate-600'
  },
  'in-progress': {
    icon: Play,
    color: 'bg-blue-100 text-blue-700 border-blue-300',
    darkColor: 'dark:bg-blue-900 dark:text-blue-300 dark:border-blue-600'
  },
  completed: {
    icon: CheckCircle2,
    color: 'bg-green-100 text-green-700 border-green-300',
    darkColor: 'dark:bg-green-900 dark:text-green-300 dark:border-green-600'
  },
  blocked: {
    icon: AlertCircle,
    color: 'bg-red-100 text-red-700 border-red-300',
    darkColor: 'dark:bg-red-900 dark:text-red-300 dark:border-red-600'
  },
  cancelled: {
    icon: XIcon,
    color: 'bg-gray-100 text-gray-700 border-gray-300',
    darkColor: 'dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600'
  }
};

export const TaskPropertyPanel: React.FC<TaskPropertyPanelProps> = ({
  task,
  onUpdate,
  className
}) => {
  const config = taskStateConfig[task.state as TaskState];
  const StateIcon = config?.icon || Clock;

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const formatEstimate = (minutes?: number) => {
    if (!minutes) return 'Not estimated';
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    if (hours === 0) return `${mins}m`;
    if (mins === 0) return `${hours}h`;
    return `${hours}h ${mins}m`;
  };

  return (
    <BasePropertyPanel
      nodeId={task.id}
      nodeType="Task"
      title={task.title}
      description={task.description}
      onUpdate={onUpdate}
      className={className}
    >
      {/* Task Status */}
      <div className="space-y-3">
        <div>
          <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
            <BarChart3 className="w-3 h-3" />
            <span>Status</span>
          </div>
          <Badge 
            variant="outline" 
            className={cn(
              "text-xs border",
              config?.color,
              config?.darkColor
            )}
          >
            <StateIcon className="w-3 h-3 mr-1" />
            {task.state.replace('-', ' ').replace(/\b\w/g, l => l.toUpperCase())}
          </Badge>
        </div>

        {/* Complexity */}
        <div>
          <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
            <Target className="w-3 h-3" />
            <span>Complexity</span>
          </div>
          <div className="flex items-center gap-2">
            <Badge variant="secondary" className="text-xs">
              Level {task.complexity}
            </Badge>
            <Progress 
              value={task.complexity * 10} 
              className="h-1 flex-1" 
            />
          </div>
        </div>

        {/* Time Estimate */}
        {task.estimate && (
          <div>
            <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
              <Clock className="w-3 h-3" />
              <span>Estimate</span>
            </div>
            <p className="text-xs font-medium">
              {formatEstimate(task.estimate)}
            </p>
          </div>
        )}

        {/* Assigned Agent */}
        {task.assignedAgent && (
          <div>
            <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
              <User className="w-3 h-3" />
              <span>Assigned Agent</span>
            </div>
            <div className="flex items-center gap-2">
              <Avatar className="w-5 h-5">
                <AvatarFallback className="text-xs">
                  {task.assignedAgent.slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <span className="text-xs font-medium">Agent</span>
            </div>
          </div>
        )}

        {/* Dependencies */}
        {task.dependencies.length > 0 && (
          <div>
            <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
              <GitBranch className="w-3 h-3" />
              <span>Dependencies</span>
            </div>
            <Badge variant="outline" className="text-xs">
              {task.dependencies.length} task{task.dependencies.length !== 1 ? 's' : ''}
            </Badge>
          </div>
        )}

        {/* Dependents */}
        {task.dependents.length > 0 && (
          <div>
            <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
              <GitBranch className="w-3 h-3 rotate-180" />
              <span>Dependents</span>
            </div>
            <Badge variant="outline" className="text-xs">
              {task.dependents.length} task{task.dependents.length !== 1 ? 's' : ''}
            </Badge>
          </div>
        )}

        {/* Timestamps */}
        <div className="grid grid-cols-1 gap-2 pt-2 border-t border-border/50">
          <div>
            <div className="flex items-center gap-2 text-xs text-muted-foreground mb-1">
              <Calendar className="w-3 h-3" />
              <span>Created</span>
            </div>
            <p className="text-xs font-medium">
              {formatDate(task.createdAt)}
            </p>
          </div>

          {task.completedAt && (
            <div>
              <div className="flex items-center gap-2 text-xs text-muted-foreground mb-1">
                <CheckCircle2 className="w-3 h-3" />
                <span>Completed</span>
              </div>
              <p className="text-xs font-medium">
                {formatDate(task.completedAt)}
              </p>
            </div>
          )}
        </div>

        {/* Hierarchy Info */}
        <div className="pt-2 border-t border-border/50">
          <div className="flex items-center justify-between text-xs">
            <span className="text-muted-foreground">Depth Level</span>
            <Badge variant="outline" className="text-xs">
              D{task.depth}
            </Badge>
          </div>
        </div>
      </div>
    </BasePropertyPanel>
  );
};