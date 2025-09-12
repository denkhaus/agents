'use client';

import { useState } from 'react';
import { TaskCard } from './task-card';
import { useProjectStore } from '@/lib/stores/project-store';
import { useUIStore } from '@/lib/stores/ui-store';
import { buildTaskHierarchy } from '@/lib/utils';
import type { Task, TaskHierarchy } from '@/lib/types';

interface TaskHierarchyItemProps {
  item: TaskHierarchy;
  onChatTask: (taskId: string) => void;
  expandedTasks: Set<string>;
  onToggleExpand: (taskId: string) => void;
}

function TaskHierarchyItem({ 
  item, 
  onChatTask, 
  expandedTasks, 
  onToggleExpand 
}: TaskHierarchyItemProps) {
  const isExpanded = expandedTasks.has(item.task.id);
  const hasChildren = item.children.length > 0;

  return (
    <div className="space-y-2">
      <TaskCard
        task={item.task}
        onChat={onChatTask}
        onToggleExpand={onToggleExpand}
        isExpanded={isExpanded}
        hasChildren={hasChildren}
        level={item.level}
        showHierarchy={true}
      />
      
      {/* Render children if expanded */}
      {isExpanded && hasChildren && (
        <div className="space-y-2">
          {item.children.map((child) => (
            <TaskHierarchyItem
              key={child.task.id}
              item={child}
              onChatTask={onChatTask}
              expandedTasks={expandedTasks}
              onToggleExpand={onToggleExpand}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export function TaskHierarchyView() {
  const { tasks } = useProjectStore();
  const { openChatPanel, selectionState } = useUIStore();
  const [expandedTasks, setExpandedTasks] = useState<Set<string>>(new Set());

  // Build hierarchy from flat task list
  const hierarchy = buildTaskHierarchy(tasks);

  const handleToggleExpand = (taskId: string) => {
    const newExpanded = new Set(expandedTasks);
    if (newExpanded.has(taskId)) {
      newExpanded.delete(taskId);
    } else {
      newExpanded.add(taskId);
    }
    setExpandedTasks(newExpanded);
  };

  const handleChatTask = (taskId: string) => {
    openChatPanel('task', taskId);
  };

  if (tasks.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            No tasks found
          </h3>
          <p className="text-gray-600">
            This project doesn&apos;t have any tasks yet.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-4">
      <div className="mb-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-2">
          Task Hierarchy
        </h3>
        <p className="text-sm text-gray-600">
          {tasks.length} tasks • Click arrows to expand/collapse subtasks
        </p>
      </div>

      <div className="space-y-4">
        {hierarchy.map((item) => (
          <TaskHierarchyItem
            key={item.task.id}
            item={item}
            onChatTask={handleChatTask}
            expandedTasks={expandedTasks}
            onToggleExpand={handleToggleExpand}
          />
        ))}
      </div>
    </div>
  );
}