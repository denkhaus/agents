"use client";

import { useChatStore } from "@/lib/store";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { normalizeToAgentId, AgentId } from "@/lib/constants/agents";

export function AgentNavigationBar() {
  const { agents, activeAgentId, setActiveAgent } = useChatStore();

  if (agents.length === 0) {
    return (
      <div className="border-t bg-background p-2">
        <div className="text-center text-sm text-muted-foreground py-2">
          No agents available - waiting for connection...
        </div>
      </div>
    );
  }

  return (
    <div className="border-t bg-background p-4">
      <div className="flex justify-center">
        <Tabs value={activeAgentId || ""} onValueChange={(value: string) => setActiveAgent(normalizeToAgentId(value) || value as AgentId)}>
          <TabsList
            className="grid h-12 w-fit max-w-4xl"
            style={{ gridTemplateColumns: `repeat(${agents.length}, 1fr)` }}
          >
            {agents.map((agent) => (
              <TabsTrigger
                key={agent.id}
                value={agent.id}
                className="flex items-center gap-2 px-3 py-2"
              >
                <div className="flex items-center gap-2 min-w-0">
                  <Avatar className="h-6 w-6">
                    <AvatarImage src={agent.avatar} alt={agent.name} />
                    <AvatarFallback className="text-xs">
                      {agent.name.slice(0, 2).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>

                  <span className="truncate text-sm font-medium">
                    {agent.name}
                  </span>

                  <Badge
                    variant={
                      agent.status === "online" ? "default" : "secondary"
                    }
                    className="h-2 w-2 p-0 rounded-full"
                  />
                </div>
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
      </div>
    </div>
  );
}
