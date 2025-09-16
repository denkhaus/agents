/**
 * Theme Toggle Component
 * Switches between light and dark mode
 */

import React from 'react';
import { Button } from '@/components/ui/button';
import { useUIStore } from '@/stores';
import { Moon, Sun } from 'lucide-react';

export const ThemeToggle: React.FC = () => {
  const { darkMode, toggleDarkMode } = useUIStore();

  React.useEffect(() => {
    // Apply theme to document
    if (darkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [darkMode]);

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={toggleDarkMode}
      className="h-8 w-8 p-0"
    >
      {darkMode ? (
        <Sun className="h-4 w-4" />
      ) : (
        <Moon className="h-4 w-4" />
      )}
      <span className="sr-only">Toggle theme</span>
    </Button>
  );
};