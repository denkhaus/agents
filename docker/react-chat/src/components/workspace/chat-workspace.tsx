"use client";

import { useChatStore } from "@/lib/store";
import { ChatInterface } from "@/components/chat/chat-interface";
import { AgentList } from "@/components/agents/agent-list";
import { BottomNavigation } from "@/components/navigation/bottom-navigation";
import {
  ResizablePanelGroup,
  ResizablePanel,
  ResizableHandle,
} from "@/components/ui/resizable";

export function ChatWorkspace() {
  const { activeAgentId, agents } = useChatStore();

  return (
    <div className="flex h-full flex-col relative">
      {agents.length === 0 ? (
        <div className="flex h-full items-center justify-center">
          <div className="text-center">
            <h3 className="text-lg font-semibold">No agents available</h3>
            <p className="text-muted-foreground">
              Waiting for agents to come online...
            </p>
          </div>
        </div>
      ) : (
        <>
          <div className="flex-1 overflow-hidden">
            <ResizablePanelGroup direction="horizontal" className="h-full">
              {!activeAgentId ? (
                <ResizablePanel defaultSize={75}>
                  <div className="flex items-center justify-center h-full">
                    <div className="text-center">
                      <h3 className="text-lg font-semibold">Select an agent</h3>
                      <p className="text-muted-foreground">
                        Choose an agent from the bottom navigation to start
                        chatting
                      </p>
                    </div>
                  </div>
                </ResizablePanel>
              ) : (
                <ResizablePanel defaultSize={75} minSize={50}>
                  <ChatInterface agentId={activeAgentId} />
                </ResizablePanel>
              )}

              <ResizableHandle withHandle />

              <ResizablePanel
                defaultSize={25}
                minSize={15}
                maxSize={40}
                className="h-screen"
              >
                <div className="h-screen">
                  <AgentList
                    agentId={activeAgentId}
                    onSessionSelect={(sessionId) => {
                      console.log("Selected session:", sessionId);
                    }}
                  />
                </div>
              </ResizablePanel>
            </ResizablePanelGroup>
          </div>

          {/* Bottom Navigation positioned absolutely */}
          <div className="absolute bottom-0 left-0 right-0 z-50">
            <BottomNavigation />
          </div>
        </>
      )}
    </div>
  );
}
