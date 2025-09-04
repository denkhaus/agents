"use client";

import { useWorkspaceStore, useChatStore } from "@/lib/store";
import { WorkspaceSidebar } from "@/components/workspace/workspace-sidebar";
import { TopNavigation } from "@/components/navigation/top-navigation";
import { BottomNavigation } from "@/components/navigation/bottom-navigation";
import { ChatWorkspace } from "@/components/workspace/chat-workspace";
import { SettingsPanel } from "@/components/workspace/settings-panel";
import {
  ResizablePanelGroup,
  ResizablePanel,
  ResizableHandle,
} from "@/components/ui/resizable";
import { useEffect } from "react";

export function MainLayout() {
  const { activeWorkspace, sidebarOpen, sidebarCollapsed } =
    useWorkspaceStore();
  const { agents } = useChatStore();

  const renderWorkspace = () => {
    switch (activeWorkspace) {
      case "chat":
        return <ChatWorkspace />;
      case "settings":
        return <SettingsPanel />;
      default:
        return <ChatWorkspace />;
    }
  };

  const shouldShowBottomNav = activeWorkspace === "chat";

  if (!sidebarOpen) {
    return (
      <div className="h-screen flex flex-col">
        <TopNavigation />
        <main
          className={`flex-1 flex flex-col ${
            shouldShowBottomNav ? "pb-20" : ""
          }`}
        >
          <div className="flex-1 overflow-hidden">{renderWorkspace()}</div>
        </main>
      </div>
    );
  }

  return (
    <div className="h-screen flex flex-col relative">
      <TopNavigation />

      <div className="flex flex-1 overflow-hidden">
        <ResizablePanelGroup
          direction="horizontal"
          className="h-full"
          autoSaveId="resizable-sidebar-layout"
        >
          <ResizablePanel
            defaultSize={20} // Set a reasonable default percentage
            minSize={sidebarCollapsed ? 4 : 15}
            maxSize={sidebarCollapsed ? 8 : 35}
          >
            <WorkspaceSidebar />
          </ResizablePanel>

          {!sidebarCollapsed && <ResizableHandle withHandle />}

          <ResizablePanel defaultSize={80} minSize={50} className="h-full">
            <main
              className={`flex flex-col h-full ${
                shouldShowBottomNav ? "pb-20" : ""
              }`}
            >
              <div className="flex-1 overflow-hidden">{renderWorkspace()}</div>
            </main>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>

      {shouldShowBottomNav && (
        <div className="absolute bottom-0 left-0 right-0 z-50">
          <BottomNavigation />
        </div>
      )}
    </div>
  );
}
