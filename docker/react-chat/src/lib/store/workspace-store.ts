import { create } from "zustand";
import { WorkspaceType } from "@/lib/types";

interface WorkspaceStore {
  // State
  activeWorkspace: WorkspaceType;
  sidebarOpen: boolean;
  sidebarCollapsed: boolean;
  sidebarWidth: number;

  // Actions
  setActiveWorkspace: (workspace: WorkspaceType) => void;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  toggleSidebarCollapsed: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  setSidebarWidth: (width: number) => void;
}

export const useWorkspaceStore = create<WorkspaceStore>((set) => ({
  // Initial state
  activeWorkspace: "chat",
  sidebarOpen: true,
  sidebarCollapsed: false,
  sidebarWidth: 250,

  // Actions
  setActiveWorkspace: (workspace) => set({ activeWorkspace: workspace }),

  toggleSidebar: () =>
    set((state) => ({
      sidebarOpen: !state.sidebarOpen,
    })),

  setSidebarOpen: (open) => set({ sidebarOpen: open }),

  toggleSidebarCollapsed: () =>
    set((state) => ({
      sidebarCollapsed: !state.sidebarCollapsed,
    })),

  setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),

  setSidebarWidth: (width) => set({ sidebarWidth: width }),
}));
