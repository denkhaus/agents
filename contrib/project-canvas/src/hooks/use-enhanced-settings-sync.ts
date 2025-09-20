/**
 * Enhanced Settings Sync Hook
 * Syncs user preferences between Settings Store, UI Store, and Convex database
 */

import React from "react";
import { useQuery, useMutation } from "convex/react";
import { api } from "../../convex/_generated/api";
import { useSettingsStore, type UserSettings } from "@/stores/settingsStore";
import { useUIStore } from "@/stores/uiStore";
import { useProjectStore } from "@/stores/projectStore";
import { UIStore } from "@/stores/uiStore"; // Re-import UIStore interface
import { ProjectStore } from "@/stores/projectStore"; // Re-import ProjectStore interface

// Simple user ID for demo purposes - in a real app this would come from auth
const USER_ID = "default-user";

export const useEnhancedSettingsSync = () => {
  const settingsStore = useSettingsStore();
  const uiStore = useUIStore();
  const projectStore = useProjectStore();

  // Query user settings from Convex
  const remoteSettings = useQuery(api.settings.getSettings, {
    userId: USER_ID,
  });

  // Mutation to update settings in Convex
  const updateRemoteSettings = useMutation(api.settings.updateSettings);

  // Sync remote settings to local stores when they change
  React.useEffect(() => {
    if (remoteSettings) {
      // Load settings into settings store
      settingsStore.loadFromRemote({
        theme: remoteSettings.theme,
        leftSidebarCollapsed: remoteSettings.leftSidebarCollapsed,
        rightSidebarCollapsed: remoteSettings.rightSidebarCollapsed,
        selectedProjectId: remoteSettings.selectedProjectId,
        selectedNodeIds: remoteSettings.selectedNodeIds || [],
        autoSave: remoteSettings.autoSave,
        notifications: remoteSettings.notifications,
        language: remoteSettings.language,
        showMiniMap: remoteSettings.showMiniMap,
        showBackground: remoteSettings.showBackground,
        autoLayout: remoteSettings.autoLayout,
        projectCanvasViewport: remoteSettings.projectCanvasViewport,
        agentsCanvasViewport: remoteSettings.agentsCanvasViewport,
        lastUpdated: remoteSettings.updatedAt || Date.now(),
      });
    }
  }, [remoteSettings, settingsStore]);

  // Sync settings store to UI store
  React.useEffect(() => {
    const settings = settingsStore.getSettingsForSync();

    // Update UI store with settings
    if (settings.theme && settings.theme !== uiStore.theme.mode) {
      uiStore.setTheme({ mode: settings.theme });
    }

    if (settings.leftSidebarCollapsed !== uiStore.sidebarCollapsed) {
      uiStore.setLeftSidebarCollapsed(settings.leftSidebarCollapsed);
    }

    if (settings.rightSidebarCollapsed !== uiStore.rightSidebarCollapsed) {
      uiStore.setRightSidebarCollapsed(settings.rightSidebarCollapsed);
    }

    if (settings.selectedNodeIds) {
      if (
        settings.selectedNodeIds.length !==
          uiStore.selection.selectedNodes.length ||
        !settings.selectedNodeIds.every((id) =>
          uiStore.selection.selectedNodes.includes(id)
        )
      ) {
        uiStore.setSelectedNodes(settings.selectedNodeIds);
      }
    }
  }, [settingsStore, uiStore]);

  // Sync settings store to project store
  React.useEffect(() => {
    const settings = settingsStore.getSettingsForSync();

    if (settings.selectedProjectId) {
      // Find the project in the store
      const project = projectStore.projects.find(
        (p) => p.id === settings.selectedProjectId
      );
      if (project && project.id !== projectStore.currentProject?.id) {
        projectStore.setCurrentProject(project);
      }
    } else if (projectStore.currentProject) {
      // If no selectedProjectId in settings but a project is current, clear it
      projectStore.setCurrentProject(null);
    }
  }, [
    settingsStore.selectedProjectId,
    projectStore.projects, // Depend on projects to ensure it runs when projects are loaded
    projectStore.currentProject?.id,
    projectStore,
  ]);

  // Sync local changes to remote
  const syncToRemote = React.useCallback(
    async (settingsUpdate: Partial<UserSettings>) => {
      try {
        await updateRemoteSettings({
          userId: USER_ID,
          theme: settingsUpdate.theme,
          notifications: settingsUpdate.notifications,
          autoSave: settingsUpdate.autoSave,
          language: settingsUpdate.language,
          leftSidebarCollapsed: settingsUpdate.leftSidebarCollapsed,
          rightSidebarCollapsed: settingsUpdate.rightSidebarCollapsed,
          selectedProjectId: settingsUpdate.selectedProjectId || undefined,
          selectedNodeIds: settingsUpdate.selectedNodeIds,
          showMiniMap: settingsUpdate.showMiniMap,
          showBackground: settingsUpdate.showBackground,
          autoLayout: settingsUpdate.autoLayout,
          projectCanvasViewport: settingsUpdate.projectCanvasViewport,
          agentsCanvasViewport: settingsUpdate.agentsCanvasViewport,
        });
      } catch (error) {
        console.error("Failed to sync settings to remote:", error);
      }
    },
    [updateRemoteSettings]
  );

  // Listen for UI store changes and sync to settings store + remote
  React.useEffect(() => {
    const unsubscribe = useUIStore.subscribe(
      (state: UIStore, prevState: UIStore) => {
        const updates: Partial<UserSettings> = {};
        let hasChanges = false;

        // Check for theme changes
        if (state.theme.mode !== prevState.theme.mode) {
          updates.theme = state.theme.mode;
          settingsStore.updateTheme(state.theme.mode);
          hasChanges = true;
        }

        // Check for sidebar changes
        if (
          state.sidebarCollapsed !== prevState.sidebarCollapsed ||
          state.rightSidebarCollapsed !== prevState.rightSidebarCollapsed
        ) {
          updates.leftSidebarCollapsed = state.sidebarCollapsed;
          updates.rightSidebarCollapsed = state.rightSidebarCollapsed;
          settingsStore.updateSidebarState(
            state.sidebarCollapsed,
            state.rightSidebarCollapsed
          );
          hasChanges = true;
        }

        // Current workspace is now handled by React Router, no need to sync via settings
        // if (state.currentWorkspace !== prevState.currentWorkspace) {
        //   updates.currentWorkspace = state.currentWorkspace;
        //   settingsStore.updateWorkspace(state.currentWorkspace);
        //   hasChanges = true;
        // }

        // Check for selection changes
        if (
          JSON.stringify(state.selection.selectedNodes) !==
          JSON.stringify(prevState.selection.selectedNodes)
        ) {
          updates.selectedNodeIds = state.selection.selectedNodes;
          settingsStore.updateSelectedNodes(state.selection.selectedNodes);
          hasChanges = true;
        }

        // Sync to remote if there are changes
        if (hasChanges) {
          syncToRemote(updates);
        }
      }
    );

    return unsubscribe;
  }, [uiStore, settingsStore, syncToRemote]);

  // Listen for project store changes
  React.useEffect(() => {
    const unsubscribe = useProjectStore.subscribe(
      (state: ProjectStore, prevState: ProjectStore) => {
        if (state.currentProject?.id !== prevState.currentProject?.id) {
          const projectId = state.currentProject?.id || null;
          settingsStore.updateSelectedProject(projectId);
          syncToRemote({ selectedProjectId: projectId });
        }
      }
    );

    return unsubscribe;
  }, [projectStore, settingsStore, syncToRemote]);

  return {
    settings: settingsStore.getSettingsForSync(),
    remoteSettings,
    syncToRemote,
    // Convenience methods
    updateTheme: (theme: UserSettings["theme"]) => {
      settingsStore.updateTheme(theme);
      syncToRemote({ theme });
    },
    updateSidebarState: (left: boolean, right: boolean) => {
      settingsStore.updateSidebarState(left, right);
      syncToRemote({
        leftSidebarCollapsed: left,
        rightSidebarCollapsed: right,
      });
    },
    updateApplicationSettings: (
      settings: Partial<
        Pick<UserSettings, "autoSave" | "notifications" | "language">
      >
    ) => {
      settingsStore.updateApplicationSettings(settings);
      syncToRemote(settings);
    },
  };
};
