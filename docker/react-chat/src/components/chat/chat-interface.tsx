"use client";

import { useEffect, useState } from "react";
import { ChatHeader } from "./chat-header";
import { MessageList } from "./message-list";
import { MessageInput } from "./message-input";
import { SessionSelector } from "./session-selector";
import { InterAgentEventDisplay } from "@/components/inter-agent/inter-agent-event-display";
import { useChatStore } from "@/lib/store";
import { Button } from "@/components/ui/button";
import { PanelLeftOpen, PanelLeftClose } from "lucide-react";
import {
  ResizablePanelGroup,
  ResizablePanel,
  ResizableHandle,
} from "@/components/ui/resizable";

interface ChatInterfaceProps {
  agentId: string;
}

export function ChatInterface({ agentId }: ChatInterfaceProps) {
  const { createSession, loadSessions } = useChatStore();
  const [showSessionSelector, setShowSessionSelector] = useState(true);

  // Load sessions when agent changes
  useEffect(() => {
    const initializeAgent = async () => {
      await loadSessions(agentId);
      // Create session will be called by setActiveAgent in the parent component
    };

    initializeAgent();
  }, [agentId, loadSessions]);

  return (
    <ResizablePanelGroup direction="horizontal" className="h-full">
      <ResizablePanel
        defaultSize={showSessionSelector ? 25 : 0}
        minSize={showSessionSelector ? 15 : 0}
        maxSize={showSessionSelector ? 40 : 0}
        collapsedSize={0}
        collapsible={true}
        onCollapse={() => setShowSessionSelector(false)}
        onExpand={() => setShowSessionSelector(true)}
        className="mr-2 min-w-[50px] transition-all duration-300 ease-in-out"
      >
        <div className="flex flex-col h-full">
          <div className="flex items-center justify-between p-2 border-b">
            <div className="flex items-center gap-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowSessionSelector(!showSessionSelector)}
              >
                {showSessionSelector ? (
                  <PanelLeftClose className="h-4 w-4" />
                ) : (
                  <PanelLeftOpen className="h-4 w-4" />
                )}
              </Button>
              <ChatHeader agentId={agentId} />
            </div>
          </div>
          <SessionSelector
            agentId={agentId}
            onSessionSelect={(sessionId) => {
              console.log("Selected session:", sessionId);
            }}
          />
        </div>
      </ResizablePanel>

      <ResizableHandle withHandle />

      {/* Main Chat Area */}
      <ResizablePanel>
        <div className="flex-1 flex flex-col h-full">
          <div className="flex-1 flex">
            <div className="flex-1 flex flex-col">
              <MessageList agentId={agentId} />
              <MessageInput agentId={agentId} />
            </div>

            <div className="w-80 border-l">
              <InterAgentEventDisplay />
            </div>
          </div>
        </div>
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}
