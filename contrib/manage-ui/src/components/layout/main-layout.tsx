'use client';

import { useEffect } from 'react';
import { Sidebar } from './sidebar';
import { TopNavigation } from './top-navigation';
import { useUIStore } from '@/lib/stores/ui-store';
import { ProjectsWorkspace } from '@/components/projects/projects-workspace';
import { initializeTheme, setupThemeListener } from '@/lib/theme-init';

export function MainLayout() {
  const { currentWorkspace, setThemeConfig } = useUIStore();

  // Initialize theme on mount
  useEffect(() => {
    const theme = initializeTheme();
    setThemeConfig({ mode: theme as 'light' | 'dark' });
    
    // Setup system theme listener
    const cleanup = setupThemeListener();
    return cleanup;
  }, [setThemeConfig]);

  const renderWorkspace = () => {
    switch (currentWorkspace) {
      case 'projects':
        return <ProjectsWorkspace />;
      case 'agents':
        return (
          <div className="flex items-center justify-center h-full">
            <div className="text-center">
              <h2 className="text-2xl font-semibold text-gray-900 mb-2">
                Agents Management
              </h2>
              <p className="text-gray-600">Coming soon...</p>
            </div>
          </div>
        );
      case 'monitoring':
        return (
          <div className="flex items-center justify-center h-full">
            <div className="text-center">
              <h2 className="text-2xl font-semibold text-gray-900 mb-2">
                System Monitoring
              </h2>
              <p className="text-gray-600">Coming soon...</p>
            </div>
          </div>
        );
      case 'settings':
        return (
          <div className="flex items-center justify-center h-full">
            <div className="text-center">
              <h2 className="text-2xl font-semibold text-gray-900 mb-2">
                Settings
              </h2>
              <p className="text-gray-600">Coming soon...</p>
            </div>
          </div>
        );
      default:
        return (
          <div className="flex items-center justify-center h-full">
            <div className="text-center">
              <h2 className="text-2xl font-semibold text-gray-900 mb-2">
                Welcome
              </h2>
              <p className="text-gray-600">Select a workspace from the sidebar</p>
            </div>
          </div>
        );
    }
  };

  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="flex-1 flex flex-col overflow-hidden">
        <TopNavigation />
        <main className="flex-1 overflow-hidden">
          {renderWorkspace()}
        </main>
      </div>
    </div>
  );
}