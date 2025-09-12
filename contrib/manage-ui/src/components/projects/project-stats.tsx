'use client';

import { BarChart3, CheckCircle, Clock, AlertCircle, XCircle } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { useProjectStore } from '@/lib/stores/project-store';
import { calculateProjectProgress } from '@/lib/utils';
import type { Project, Task } from '@/lib/types';

interface ProjectStatsProps {
  projects: Project[];
  allTasks: Task[];
}

export function ProjectStats({ projects, allTasks }: ProjectStatsProps) {
  // Calculate overall statistics
  const totalProjects = projects.length;
  const totalTasks = allTasks.length;
  
  const tasksByState = {
    completed: allTasks.filter(t => t.state === 'completed').length,
    inProgress: allTasks.filter(t => t.state === 'in-progress').length,
    pending: allTasks.filter(t => t.state === 'pending').length,
    blocked: allTasks.filter(t => t.state === 'blocked').length,
    cancelled: allTasks.filter(t => t.state === 'cancelled').length,
  };

  const overallProgress = totalTasks > 0 ? (tasksByState.completed / totalTasks) * 100 : 0;

  // Calculate project completion distribution
  const projectsByCompletion = {
    high: projects.filter(p => p.progress >= 75).length,
    medium: projects.filter(p => p.progress >= 25 && p.progress < 75).length,
    low: projects.filter(p => p.progress < 25).length,
  };

  const stats = [
    {
      label: 'Total Projects',
      value: totalProjects,
      icon: BarChart3,
      color: 'text-blue-600',
      bgColor: 'bg-blue-100',
    },
    {
      label: 'Completed Tasks',
      value: tasksByState.completed,
      icon: CheckCircle,
      color: 'text-green-600',
      bgColor: 'bg-green-100',
    },
    {
      label: 'In Progress',
      value: tasksByState.inProgress,
      icon: Clock,
      color: 'text-blue-600',
      bgColor: 'bg-blue-100',
    },
    {
      label: 'Blocked Tasks',
      value: tasksByState.blocked,
      icon: AlertCircle,
      color: 'text-red-600',
      bgColor: 'bg-red-100',
    },
  ];

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 mb-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-gray-900">
          System Overview
        </h3>
        <Badge variant="secondary">
          {Math.round(overallProgress)}% Complete
        </Badge>
      </div>

      {/* Main Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <div key={stat.label} className="flex items-center space-x-3 p-3 rounded-lg bg-gray-50">
              <div className={`p-2 rounded-lg ${stat.bgColor}`}>
                <Icon className={`h-5 w-5 ${stat.color}`} />
              </div>
              <div>
                <p className="text-2xl font-bold text-gray-900">{stat.value}</p>
                <p className="text-sm text-gray-600">{stat.label}</p>
              </div>
            </div>
          );
        })}
      </div>

      {/* Progress Breakdown */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Task Status Distribution */}
        <div>
          <h4 className="font-medium text-gray-900 mb-3">Task Status Distribution</h4>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Completed</span>
              <div className="flex items-center space-x-2">
                <div className="w-24 bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-green-500 h-2 rounded-full" 
                    style={{ width: `${(tasksByState.completed / totalTasks) * 100}%` }}
                  />
                </div>
                <span className="text-sm font-medium">{tasksByState.completed}</span>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">In Progress</span>
              <div className="flex items-center space-x-2">
                <div className="w-24 bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-blue-500 h-2 rounded-full" 
                    style={{ width: `${(tasksByState.inProgress / totalTasks) * 100}%` }}
                  />
                </div>
                <span className="text-sm font-medium">{tasksByState.inProgress}</span>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Pending</span>
              <div className="flex items-center space-x-2">
                <div className="w-24 bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-gray-400 h-2 rounded-full" 
                    style={{ width: `${(tasksByState.pending / totalTasks) * 100}%` }}
                  />
                </div>
                <span className="text-sm font-medium">{tasksByState.pending}</span>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Blocked</span>
              <div className="flex items-center space-x-2">
                <div className="w-24 bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-red-500 h-2 rounded-full" 
                    style={{ width: `${(tasksByState.blocked / totalTasks) * 100}%` }}
                  />
                </div>
                <span className="text-sm font-medium">{tasksByState.blocked}</span>
              </div>
            </div>
          </div>
        </div>

        {/* Project Completion Distribution */}
        <div>
          <h4 className="font-medium text-gray-900 mb-3">Project Completion</h4>
          <div className="space-y-3">
            <div className="flex items-center justify-between p-2 rounded bg-green-50">
              <span className="text-sm text-green-800">High Progress (75%+)</span>
              <Badge className="bg-green-100 text-green-800">
                {projectsByCompletion.high} projects
              </Badge>
            </div>
            <div className="flex items-center justify-between p-2 rounded bg-yellow-50">
              <span className="text-sm text-yellow-800">Medium Progress (25-75%)</span>
              <Badge className="bg-yellow-100 text-yellow-800">
                {projectsByCompletion.medium} projects
              </Badge>
            </div>
            <div className="flex items-center justify-between p-2 rounded bg-red-50">
              <span className="text-sm text-red-800">Low Progress (&lt;25%)</span>
              <Badge className="bg-red-100 text-red-800">
                {projectsByCompletion.low} projects
              </Badge>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}