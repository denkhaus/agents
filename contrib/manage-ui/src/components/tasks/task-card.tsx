'use client';

import { MessageSquare, ChevronDown, ChevronRight, Users, Clock, AlertCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatRelativeTime, getTaskStateColor, formatEstimate, cn } from '@/lib/utils';
import type { Task } from '@/lib/types';

interface TaskCardProps {
  task: Task;
  onChat: (taskId: string) => void;
  onToggleExpand?: (taskId: string) => void;
  isExpanded?: boolean;
  hasChildren?: boolean;
  level?: number;
  showHierarchy?: boolean;
}

export function TaskCard({ 
  task, 
  onChat, 
  onToggleExpand, 
  isExpanded = false, 
  hasChildren = false,
  level = 0,
  showHierarchy = false
}: TaskCardProps) {
  const indentLevel = showHierarchy ? level * 16 : 0;

  return (
    <div 
      className={cn(
        "bg-white rounded-lg border border-gray-200 hover:shadow-md transition-all duration-200",
        showHierarchy && "ml-" + indentLevel
      )}
      style={showHierarchy ? { marginLeft: `${indentLevel}px` } : undefined}
    >
      <div className="p-4">
        {/* Header with expand/collapse and chat */}
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-start space-x-2 flex-1">
            {/* Hierarchy toggle */}
            {showHierarchy && hasChildren && onToggleExpand && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onToggleExpand(task.id)}
                className="h-6 w-6 p-0 mt-0.5"
              >
                {isExpanded ? (
                  <ChevronDown className="h-3 w-3" />
                ) : (
                  <ChevronRight className="h-3 w-3" />
                )}
              </Button>
            )}
            
            {/* Task content */}
            <div className="flex-1 min-w-0">
              <h4 className="font-medium text-gray-900 mb-1 line-clamp-2">
                {task.title}
              </h4>
              <p className="text-xs text-gray-500 mb-2">
                ID: {task.id} • Depth: {task.depth}
              </p>
            </div>
          </div>
          
          {/* Chat button */}
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onChat(task.id)}
            className="h-6 w-6 p-0 ml-2 flex-shrink-0"
          >
            <MessageSquare className="h-3 w-3" />
          </Button>
        </div>

        {/* Description */}
        {task.description && (
          <p className="text-sm text-gray-600 mb-3 line-clamp-2">
            {task.description}
          </p>
        )}

        {/* Status and metadata */}
        <div className="space-y-2">
          {/* State and complexity */}
          <div className="flex items-center justify-between">
            <Badge className={cn("text-xs", getTaskStateColor(task.state))}>
              {task.state.replace('-', ' ')}
            </Badge>
            <div className="flex items-center space-x-2 text-xs text-gray-500">
              <span>Complexity: {task.complexity}/10</span>
            </div>
          </div>

          {/* Time and assignment */}
          <div className="flex items-center justify-between text-xs text-gray-500">
            <div className="flex items-center space-x-1">
              <Clock className="h-3 w-3" />
              <span>{formatEstimate(task.estimate)}</span>
            </div>
            {task.assigned_agent && (
              <div className="flex items-center space-x-1">
                <Users className="h-3 w-3" />
                <span className="truncate max-w-20">{task.assigned_agent}</span>
              </div>
            )}
          </div>

          {/* Dependencies and warnings */}
          {(task.dependencies && task.dependencies.length > 0) && (
            <div className="flex items-center space-x-1 text-xs text-amber-600">
              <AlertCircle className="h-3 w-3" />
              <span>Depends on {task.dependencies.length} task(s)</span>
            </div>
          )}

          {/* Updated timestamp */}
          <div className="text-xs text-gray-400 pt-1 border-t border-gray-100">
            Updated {formatRelativeTime(task.updated_at)}
          </div>
        </div>
      </div>
    </div>
  );
}