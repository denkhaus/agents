/**
 * Convex Settings Hook
 * Syncs user preferences between UI store and Convex database
 */

import React from "react";
import { useQuery, useMutation } from "convex/react";
import { api } from "../../convex/_generated/api";
import { useUIStore } from "@/stores";
import { ThemeConfig } from "@/types";

// Simple user ID for demo purposes - in a real app this would come from auth
const USER_ID = "default-user";

export const useSettingsSync = () => {
  const { theme, setTheme } = useUIStore();
  
  // Query user settings from Convex
  const settings = useQuery(api.settings.getSettings, { userId: USER_ID });
  
  // Mutation to update settings in Convex
  const updateThemeSetting = useMutation(api.settings.updateTheme);
  
  // Sync Convex settings to UI store when they change
  React.useEffect(() => {
    if (settings) {
      // Apply theme from Convex settings if it exists and differs from current
      if (settings.theme && settings.theme !== theme.mode) {
        setTheme({ mode: settings.theme });
      }
    }
  }, [settings, theme.mode, setTheme]);
  
  // Sync UI store changes to Convex
  const syncThemeToConvex = React.useCallback((newTheme: ThemeConfig) => {
    // Update Convex settings
    updateThemeSetting({
      userId: USER_ID,
      theme: newTheme.mode,
    }).catch(console.error);
  }, [updateThemeSetting]);
  
  return {
    settings,
    syncThemeToConvex,
  };
};