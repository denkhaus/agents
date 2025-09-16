/**
 * Theme Toggle Component
 * Switches between light and dark mode with Convex persistence
 */

import React from 'react';
import { Button } from '@/components/ui/button';
import { useUIStore } from '@/stores';
import { useSettingsSync } from '@/hooks/use-settings-sync';
import { Moon, Sun } from 'lucide-react';

export const ThemeToggle: React.FC = () => {
  const { theme, setTheme } = useUIStore();
  const { syncThemeToConvex } = useSettingsSync();

  const toggleDarkMode = React.useCallback(() => {
    const newMode = theme.mode === 'light' ? 'dark' : 'light';
    setTheme({ mode: newMode });
    
    // Sync to Convex
    syncThemeToConvex({ ...theme, mode: newMode });
  }, [theme, setTheme, syncThemeToConvex]);

  React.useEffect(() => {
    // Apply theme to document
    if (theme.mode === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [theme.mode]);

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={toggleDarkMode}
      className="h-8 w-8 p-0"
    >
      {theme.mode === 'dark' ? (
        <Sun className="h-4 w-4" />
      ) : (
        <Moon className="h-4 w-4" />
      )}
      <span className="sr-only">Toggle theme</span>
    </Button>
  );
};