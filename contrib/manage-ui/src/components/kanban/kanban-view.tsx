'use client';

import { useState, useEffect } from 'react';
import { ArrowLeft, MessageSquare, Plus, List, Grid } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useProjectStore } from '@/lib/stores/project-store';
import { useUIStore } from '@/lib/stores/ui-store';
import { TaskCard } from '@/components/tasks/task-card';
import { TaskHierarchyView } from '@/components/tasks/task-hierarchy-view';
import { formatRelativeTime, getTaskStateColor, getTaskStateIcon, formatEstimate, cn } from '@/lib/utils';
import type { Task, TaskState } from '@/lib/types';

interface KanbanTaskCardProps {
  task: Task;
  onChat: (taskId: string) => void;
}

function KanbanTaskCard({ task, onChat }: KanbanTaskCardProps) {
  return (
    <TaskCard
      task={task}
      onChat={onChat}
      showHierarchy={false}
    />
  );
}

interface KanbanColumnProps {
  title: string;
  tasks: Task[];
  state?: TaskState;
  onAddTask?: () => void;
  onChatTask: (taskId: string) => void;
}

function KanbanColumn({ title, tasks, state, onAddTask, onChatTask }: KanbanColumnProps) {
  return (
    <div className="flex flex-col h-full bg-gray-50 rounded-lg">
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold text-gray-900">{title}</h3>
          <div className="flex items-center space-x-2">
            <Badge variant="secondary" className="text-xs">
              {tasks.length}
            </Badge>
            {onAddTask && (
              <Button
                variant="ghost"
                size="sm"
                onClick={onAddTask}
                className="h-6 w-6 p-0"
              >
                <Plus className="h-3 w-3" />
              </Button>
            )}
          </div>
        </div>
      </div>
      
      <div className="flex-1 p-4 space-y-3 overflow-y-auto">
        {tasks.map((task) => (
          <KanbanTaskCard
            key={task.id}
            task={task}
            onChat={onChatTask}
          />
        ))}
        
        {tasks.length === 0 && (
          <div className="text-center text-gray-500 text-sm py-8">
            No tasks in this state
          </div>
        )}
      </div>
    </div>
  );
}

export function KanbanView() {
  const { 
    currentProject, 
    tasks, 
    setTasks, 
    isLoadingTasks
  } = useProjectStore();
  
  const { 
    setSelectedProject, 
    openChatPanel,
    openCreateTaskModal,
    viewConfig,
    setViewConfig,
    selectionState
  } = useUIStore();

  const [isInitialized, setIsInitialized] = useState(false);

  // Load tasks for the selected project
  useEffect(() => {
    const loadTasks = async () => {
      if (!selectionState?.selectedProjectId || isInitialized) return;
      
      try {
        // Import corrected mock data with proper UUIDs and filter by project
        const { mockTasks } = await import('@/lib/data/mock-data-corrected');
        const projectTasks = mockTasks.filter(task => task.project_id === selectionState?.selectedProjectId);
        
        setTasks(projectTasks);
        setIsInitialized(true);
      } catch (error) {
        console.error('Failed to load tasks:', error);
      }
    };

    loadTasks();
  }, [selectionState?.selectedProjectId, setTasks, isInitialized]);

  const handleBackToProjects = () => {
    setSelectedProject(undefined);
    setIsInitialized(false);
  };

  const handleChatTask = (taskId: string) => {
    openChatPanel('task', taskId);
  };

  const handleAddTask = () => {
    openCreateTaskModal();
  };

  const handleViewToggle = () => {
    setViewConfig({
      type: viewConfig.type === 'kanban' ? 'list' : 'kanban'
    });
  };

  // Group tasks by state
  const tasksByState = {
    pending: tasks.filter(t => t.state === 'pending'),
    'in-progress': tasks.filter(t => t.state === 'in-progress'),
    completed: tasks.filter(t => t.state === 'completed'),
    blocked: tasks.filter(t => t.state === 'blocked'),
    cancelled: tasks.filter(t => t.state === 'cancelled'),
  };

  if (isLoadingTasks) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading tasks...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header with back button and view toggle */}
      <div className="p-4 border-b border-gray-200 bg-white">
        <div className="flex items-center justify-between">
          <Button
            variant="ghost"
            onClick={handleBackToProjects}
            className="h-8"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Projects
          </Button>

          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleViewToggle}
              className="h-8"
            >
              {viewConfig.type === 'kanban' ? (
                <>
                  <List className="h-4 w-4 mr-2" />
                  Hierarchy View
                </>
              ) : (
                <>
                  <Grid className="h-4 w-4 mr-2" />
                  Kanban View
                </>
              )}
            </Button>
            
            <Button
              onClick={handleAddTask}
              size="sm"
              className="h-8"
            >
              <Plus className="h-4 w-4 mr-2" />
              New Task
            </Button>
          </div>
        </div>
      </div>

      {/* Content area */}
      <div className="flex-1 overflow-hidden">
        {viewConfig.type === 'kanban' ? (
          /* Kanban board */
          <div className="p-6 h-full overflow-hidden">
            <div className="grid grid-cols-5 gap-6 h-full">
              <KanbanColumn
                title="Pending"
                tasks={tasksByState.pending}
                state="pending"
                onAddTask={handleAddTask}
                onChatTask={handleChatTask}
              />
              <KanbanColumn
                title="In Progress"
                tasks={tasksByState['in-progress']}
                state="in-progress"
                onChatTask={handleChatTask}
              />
              <KanbanColumn
                title="Completed"
                tasks={tasksByState.completed}
                state="completed"
                onChatTask={handleChatTask}
              />
              <KanbanColumn
                title="Blocked"
                tasks={tasksByState.blocked}
                state="blocked"
                onChatTask={handleChatTask}
              />
              <KanbanColumn
                title="Cancelled"
                tasks={tasksByState.cancelled}
                state="cancelled"
                onChatTask={handleChatTask}
              />
            </div>
          </div>
        ) : (
          /* Hierarchy view */
          <TaskHierarchyView />
        )}
      </div>
    </div>
  );
}