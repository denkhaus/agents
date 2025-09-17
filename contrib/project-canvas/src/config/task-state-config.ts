/**
 * Task State Configuration
 * Central configuration for task state styling and behavior
 */

import React from "react";
import { Clock, AlertCircle, CheckCircle2, Play, X } from "lucide-react";
import { TaskState } from "@/types/task.types";

export interface TaskStateConfig {
  color: string;
  darkColor: string;
  icon: React.ComponentType<{ className?: string }>;
  badge: "default" | "secondary" | "destructive" | "outline";
}

export const taskStateConfig: Record<TaskState, TaskStateConfig> = {
  [TaskState.PENDING]: {
    color: "bg-slate-100 border-slate-300 text-slate-700",
    darkColor: "dark:bg-slate-800 dark:border-slate-600 dark:text-slate-300",
    icon: Clock,
    badge: "secondary",
  },
  [TaskState.IN_PROGRESS]: {
    color: "bg-blue-50 border-blue-300 text-blue-700",
    darkColor: "dark:bg-blue-900 dark:border-blue-600 dark:text-blue-300",
    icon: Play,
    badge: "default",
  },
  [TaskState.COMPLETED]: {
    color: "bg-green-50 border-green-300 text-green-700",
    darkColor: "dark:bg-green-900 dark:border-green-600 dark:text-green-300",
    icon: CheckCircle2,
    badge: "default",
  },
  [TaskState.BLOCKED]: {
    color: "bg-red-50 border-red-300 text-red-700",
    darkColor: "dark:bg-red-900 dark:border-red-600 dark:text-red-300",
    icon: AlertCircle,
    badge: "destructive",
  },
  [TaskState.CANCELLED]: {
    color: "bg-gray-50 border-gray-300 text-gray-700",
    darkColor: "dark:bg-gray-800 dark:border-gray-600 dark:text-gray-300",
    icon: X,
    badge: "secondary",
  },
};

/**
 * Helper function to get task state configuration
 */
export const getTaskStateConfig = (state: TaskState): TaskStateConfig => {
  return taskStateConfig[state] || taskStateConfig[TaskState.PENDING];
};
