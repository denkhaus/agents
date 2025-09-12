'use client';

import { Plus, Filter, SortAsc, Grid, List, MoreHorizontal } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useUIStore } from '@/lib/stores/ui-store';
import { useProjectStore } from '@/lib/stores/project-store';

export function ProjectsHeader() {
  const { 
    viewConfig, 
    setViewConfig,
    selectionState,
    openCreateProjectModal,
    openCreateTaskModal 
  } = useUIStore();
  
  const { currentProject } = useProjectStore();

  const handleViewTypeChange = () => {
    setViewConfig({
      type: viewConfig.type === 'kanban' ? 'list' : 'kanban'
    });
  };

  const getHeaderTitle = () => {
    if (selectionState?.selectedProjectId && currentProject) {
      return currentProject.title;
    }
    return 'All Projects';
  };

  const getHeaderActions = () => {
    if (selectionState?.selectedProjectId) {
      // Project-specific actions
      return (
        <>
          <Button
            onClick={openCreateTaskModal}
            size="sm"
            className="h-9"
          >
            <Plus className="h-4 w-4 mr-2" />
            New Task
          </Button>
          
          <Button
            variant="outline"
            size="sm"
            onClick={handleViewTypeChange}
            className="h-9"
          >
            {viewConfig.type === 'kanban' ? (
              <List className="h-4 w-4 mr-2" />
            ) : (
              <Grid className="h-4 w-4 mr-2" />
            )}
            {viewConfig.type === 'kanban' ? 'List View' : 'Kanban View'}
          </Button>
        </>
      );
    }

    // Projects list actions
    return (
      <Button
        onClick={openCreateProjectModal}
        size="sm"
        className="h-9"
      >
        <Plus className="h-4 w-4 mr-2" />
        New Project
      </Button>
    );
  };

  return (
    <div className="bg-white border-b border-gray-200 px-6 py-4">
      <div className="flex items-center justify-between">
        {/* Left side - Title and breadcrumb */}
        <div className="flex items-center space-x-4">
          <div>
            <h2 className="text-lg font-semibold text-gray-900">
              {getHeaderTitle()}
            </h2>
            {selectionState?.selectedProjectId && currentProject && (
              <p className="text-sm text-gray-600 mt-1">
                {currentProject.total_tasks} tasks • {Math.round(currentProject.progress)}% complete
              </p>
            )}
          </div>
        </div>

        {/* Right side - Actions */}
        <div className="flex items-center space-x-2">
          {/* Filter button */}
          <Button
            variant="outline"
            size="sm"
            className="h-9"
          >
            <Filter className="h-4 w-4 mr-2" />
            Filter
          </Button>

          {/* Sort button */}
          <Button
            variant="outline"
            size="sm"
            className="h-9"
          >
            <SortAsc className="h-4 w-4 mr-2" />
            Sort
          </Button>

          {/* Dynamic actions */}
          {getHeaderActions()}

          {/* More options */}
          <Button
            variant="outline"
            size="sm"
            className="h-9 w-9 p-0"
          >
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}