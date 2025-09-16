/**
 * Task Node Component
 * Custom ReactFlow node for displaying tasks
 */

import React from 'react';
import { Handle, Position, NodeProps } from 'reactflow';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { MarkdownRenderer } from '@/components/ui/markdown-renderer';
import { TaskNodeData } from '@/types/reactflow.types';
import { TaskState } from '@/types/task.types';
import { 
  Clock, 
  User, 
  AlertCircle, 
  CheckCircle2, 
  Play, 
  X
} from 'lucide-react';
import { cn } from '@/lib/utils';

const taskStateConfig = {
  [TaskState.PENDING]: {
    color: 'bg-slate-100 border-slate-300 text-slate-700',
    darkColor: 'dark:bg-slate-800 dark:border-slate-600 dark:text-slate-300',
    icon: Clock,
    badge: 'secondary'
  },
  [TaskState.IN_PROGRESS]: {
    color: 'bg-blue-50 border-blue-300 text-blue-700',
    darkColor: 'dark:bg-blue-900/20 dark:border-blue-600 dark:text-blue-300',
    icon: Play,
    badge: 'default'
  },
  [TaskState.COMPLETED]: {
    color: 'bg-green-50 border-green-300 text-green-700',
    darkColor: 'dark:bg-green-900/20 dark:border-green-600 dark:text-green-300',
    icon: CheckCircle2,
    badge: 'default'
  },
  [TaskState.BLOCKED]: {
    color: 'bg-red-50 border-red-300 text-red-700',
    darkColor: 'dark:bg-red-900/20 dark:border-red-600 dark:text-red-300',
    icon: AlertCircle,
    badge: 'destructive'
  },
  [TaskState.CANCELLED]: {
    color: 'bg-gray-50 border-gray-300 text-gray-700',
    darkColor: 'dark:bg-gray-800 dark:border-gray-600 dark:text-gray-300',
    icon: X,
    badge: 'secondary'
  }
} as const;

export const TaskNode: React.FC<NodeProps<TaskNodeData>> = ({ 
  data, 
  selected 
}) => {
  const { task, isHighlighted, showDetails } = data;
  const config = taskStateConfig[task.state];
  const Icon = config.icon;

  return (
    <div className="task-node">
      {/* Input Handle */}
      <Handle
        type="target"
        position={Position.Left}
        className="w-3 h-3 border-2 border-background"
        style={{ background: '#6b7280' }}
      />

      <Card 
        className={cn(
          "w-80 transition-all duration-200 cursor-pointer",
          config.color,
          config.darkColor,
          selected && "ring-2 ring-primary ring-offset-2",
          isHighlighted && "shadow-lg scale-105"
        )}
      >
        <CardHeader className="pb-2">
          <div className="flex items-start justify-between gap-2">
            <div className="flex-1 min-w-0">
              <h3 className="font-medium text-sm leading-tight truncate">
                {task.title}
              </h3>
              <div className="flex items-center gap-2 mt-1">
                <Badge 
                  variant={config.badge as any}
                  className="text-xs"
                >
                  <Icon className="w-3 h-3 mr-1" />
                  {task.state.replace('-', ' ')}
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

        <CardContent className="pt-0">
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
        style={{ background: '#6b7280' }}
      />
    </div>
  );
};