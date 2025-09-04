"use client";

import { useState, useEffect } from "react";
import { useChatStore } from "@/lib/store";
import { AgentCard } from "./agent-card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { InterAgentEventDisplay } from "@/components/inter-agent/inter-agent-event-display";
import { SessionSelector } from "@/components/chat/session-selector";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

interface AgentListProps {
  agentId?: string;
  onSessionSelect?: (sessionId: string) => void;
}

export function AgentList({ agentId, onSessionSelect }: AgentListProps) {
  const { agents, interAgentEvents } = useChatStore();
  const [openItems, setOpenItems] = useState<string[]>([]);

  useEffect(() => {
    if (typeof window !== "undefined") {
      const storedOpenItems = localStorage.getItem("accordionOpenItems");
      if (storedOpenItems) {
        setOpenItems(JSON.parse(storedOpenItems));
      } else {
        // Set default open items if nothing is stored
        setOpenItems(["agents", "chat-sessions"]);
      }
    }
  }, []);

  const handleValueChange = (value: string[]) => {
    setOpenItems(value);
    if (typeof window !== "undefined") {
      localStorage.setItem("accordionOpenItems", JSON.stringify(value));
    }
  };

  return (
    <div className="h-full flex flex-col bg-card text-card-foreground">
      <Accordion
        type="multiple"
        value={openItems}
        onValueChange={handleValueChange}
        className="flex-1"
      >
        {agentId && (
          <AccordionItem value="chat-sessions" className="border-b">
            <AccordionTrigger className="px-4 py-2 hover:no-underline">
              <div className="flex items-center gap-2">
                <h2 className="font-semibold">Chat Sessions</h2>
              </div>
            </AccordionTrigger>
            <AccordionContent className="pb-0">
              <div className="px-4">
                <SessionSelector
                  agentId={agentId}
                  onSessionSelect={onSessionSelect}
                />
              </div>
            </AccordionContent>
          </AccordionItem>
        )}

        <AccordionItem value="agents" className="border-b">
          <AccordionTrigger className="px-4 py-2 hover:no-underline">
            <div className="flex items-center gap-2">
              <h2 className="font-semibold">Available Agents</h2>
              <span className="text-sm text-muted-foreground">
                ({agents.length})
              </span>
            </div>
          </AccordionTrigger>
          <AccordionContent className="pb-0 h-full flex flex-col">
            <ScrollArea className="flex-1">
              <div className="p-4 space-y-3">
                {agents.length === 0 ? (
                  <div className="text-center py-8">
                    <p className="text-muted-foreground">No agents available</p>
                    <p className="text-sm text-muted-foreground mt-1">
                      Agents will appear here when they come online
                    </p>
                  </div>
                ) : (
                  agents.map((agent, index: number) => (
                    <div key={agent.id}>
                      <AgentCard agent={agent} />
                      {index < agents.length - 1 && (
                        <Separator className="mt-3" />
                      )}
                    </div>
                  ))
                )}
              </div>
            </ScrollArea>
          </AccordionContent>
        </AccordionItem>

        <AccordionItem value="inter-agent" className="border-b-0">
          <AccordionTrigger className="px-4 py-2 hover:no-underline">
            <div className="flex items-center gap-2">
              <h2 className="font-semibold">Inter-Agent Communication</h2>
              <span className="text-sm text-muted-foreground">
                ({interAgentEvents.length})
              </span>
            </div>
          </AccordionTrigger>
          <AccordionContent className="pb-0">
            <div className="h-96">
              <InterAgentEventDisplay />
            </div>
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>
  );
}
