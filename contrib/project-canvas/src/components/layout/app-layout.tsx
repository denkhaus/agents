/**
 * Main Application Layout
 * Provides the overall structure with sidebar and main content area
 */

import React from 'react';
import { Sidebar } from './sidebar';
import { TopNavigation } from './top-navigation';
import { useUIStore } from '@/stores';
import { cn } from '@/lib/utils';

interface AppLayoutProps {
  children: React.ReactNode;
}

export const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const { sidebarCollapsed } = useUIStore();

  return (
    <div className="h-screen flex flex-col bg-background">
      {/* Top Navigation */}
      <TopNavigation />
      
      {/* Main Content Area */}
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar */}
        <Sidebar />
        
        {/* Main Content */}
        <main 
          className={cn(
            "flex-1 transition-all duration-300 ease-in-out",
            "bg-background border-l border-border",
            sidebarCollapsed ? "ml-16" : "ml-64"
          )}
        >
          {children}
        </main>
      </div>
    </div>
  );
};