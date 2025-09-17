/**
 * Settings Store - Zustand
 * Manages persistent user settings and preferences
 */

import { create } from "zustand";
import { devtools, persist } from "zustand/middleware";
import type { UUID, WorkspaceType, ThemeConfig } from "@/types";

export interface UserSettings {
  // UI Preferences
  theme: ThemeConfig["mode"];
  leftSidebarCollapsed: boolean;
  rightSidebarCollapsed: boolean;
  
  // Workspace State
  currentWorkspace: WorkspaceType;
  selectedProjectId: UUID | null;
  selectedNodeIds: UUID[];
  
  // Application Settings
  autoSave: boolean;
  notifications: boolean;
  language: string;
  
  // Canvas Settings
  showMiniMap: boolean;
  showBackground: boolean;
  autoLayout: boolean;
  canvasZoom: number;
  
  // Timestamps
  lastUpdated: number;
}

interface SettingsStore extends UserSettings {
  // Actions
  updateTheme: (theme: ThemeConfig["mode"]) => void;
  updateSidebarState: (left: boolean, right: boolean) => void;
  updateWorkspace: (workspace: WorkspaceType) => void;
  updateSelectedProject: (projectId: UUID | null) => void;
  updateSelectedNodes: (nodeIds: UUID[]) => void;
  updateApplicationSettings: (settings: Partial<Pick<UserSettings, 'autoSave' | 'notifications' | 'language'>>) => void;
  updateCanvasSettings: (settings: Partial<Pick<UserSettings, 'showMiniMap' | 'showBackground' | 'autoLayout' | 'canvasZoom'>>) => void;
  
  // Bulk operations
  updateSettings: (settings: Partial<UserSettings>) => void;
  resetSettings: () => void;
  
  // Sync operations
  loadFromRemote: (remoteSettings: Partial<UserSettings>) => void;
  getSettingsForSync: () => UserSettings;
}

const defaultSettings: UserSettings = {
  // UI Preferences
  theme: "light",
  leftSidebarCollapsed: false,
  rightSidebarCollapsed: true,
  
  // Workspace State
  currentWorkspace: "projects",
  selectedProjectId: null,
  selectedNodeIds: [],
  
  // Application Settings
  autoSave: true,
  notifications: true,
  language: "en",
  
  // Canvas Settings
  showMiniMap: true,
  showBackground: true,
  autoLayout: true,
  canvasZoom: 100,
  
  // Timestamps
  lastUpdated: Date.now(),
};

export const useSettingsStore = create<SettingsStore>()(
  devtools(
    persist(
      (set, get) => ({
        ...defaultSettings,

        // Theme Actions
        updateTheme: (theme) =>
          set(
            {
              theme,
              lastUpdated: Date.now(),
            },
            false,
            "settings/updateTheme"
          ),

        // Sidebar Actions
        updateSidebarState: (leftSidebarCollapsed, rightSidebarCollapsed) =>
          set(
            {
              leftSidebarCollapsed,
              rightSidebarCollapsed,
              lastUpdated: Date.now(),
            },
            false,
            "settings/updateSidebarState"
          ),

        // Workspace Actions
        updateWorkspace: (currentWorkspace) =>
          set(
            {
              currentWorkspace,
              lastUpdated: Date.now(),
            },
            false,
            "settings/updateWorkspace"
          ),

        updateSelectedProject: (selectedProjectId) =>
          set(
            {
              selectedProjectId,
              lastUpdated: Date.now(),
            },
            false,
            "settings/updateSelectedProject"
          ),

        updateSelectedNodes: (selectedNodeIds) =>
          set(
            {
              selectedNodeIds,
              lastUpdated: Date.now(),
            },
            false,
            "settings/updateSelectedNodes"
          ),

        // Application Settings
        updateApplicationSettings: (settings) =>
          set(
            {
              ...settings,
              lastUpdated: Date.now(),
            },
            false,
            "settings/updateApplicationSettings"
          ),

        // Canvas Settings
        updateCanvasSettings: (settings) =>
          set(
            {
              ...settings,
              lastUpdated: Date.now(),
            },
            false,
            "settings/updateCanvasSettings"
          ),

        // Bulk Operations
        updateSettings: (settings) =>
          set(
            {
              ...settings,
              lastUpdated: Date.now(),
            },
            false,
            "settings/updateSettings"
          ),

        resetSettings: () =>
          set(
            {
              ...defaultSettings,
              lastUpdated: Date.now(),
            },
            false,
            "settings/resetSettings"
          ),

        // Sync Operations
        loadFromRemote: (remoteSettings) =>
          set(
            (state) => {
              // Only update if remote settings are newer
              if (remoteSettings.lastUpdated && remoteSettings.lastUpdated > state.lastUpdated) {
                return {
                  ...state,
                  ...remoteSettings,
                };
              }
              return state;
            },
            false,
            "settings/loadFromRemote"
          ),

        getSettingsForSync: () => {
          const state = get();
          return {
            theme: state.theme,
            leftSidebarCollapsed: state.leftSidebarCollapsed,
            rightSidebarCollapsed: state.rightSidebarCollapsed,
            currentWorkspace: state.currentWorkspace,
            selectedProjectId: state.selectedProjectId,
            selectedNodeIds: state.selectedNodeIds,
            autoSave: state.autoSave,
            notifications: state.notifications,
            language: state.language,
            showMiniMap: state.showMiniMap,
            showBackground: state.showBackground,
            autoLayout: state.autoLayout,
            canvasZoom: state.canvasZoom,
            lastUpdated: state.lastUpdated,
          };
        },
      }),
      {
        name: "user-settings",
        // Only persist certain settings locally
        partialize: (state) => ({
          theme: state.theme,
          leftSidebarCollapsed: state.leftSidebarCollapsed,
          rightSidebarCollapsed: state.rightSidebarCollapsed,
          currentWorkspace: state.currentWorkspace,
          selectedProjectId: state.selectedProjectId,
          autoSave: state.autoSave,
          notifications: state.notifications,
          language: state.language,
          showMiniMap: state.showMiniMap,
          showBackground: state.showBackground,
          autoLayout: state.autoLayout,
          canvasZoom: state.canvasZoom,
          lastUpdated: state.lastUpdated,
        }),
      }
    ),
    { name: "settings-store" }
  )
);