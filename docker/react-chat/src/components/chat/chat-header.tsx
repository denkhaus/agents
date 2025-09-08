"use client";

import { useChatStore } from "@/lib/store";
import { Badge } from "@/components/ui/badge";
import { AgentInfo } from '@/lib/types'

interface ChatHeaderProps {
  agentId: string;
}

export function ChatHeader({ agentId }: ChatHeaderProps) {
  const { agents } = useChatStore();
  const agent = agents.find((a: Agent) => a.id === agentId);

  if (!agent) {
    return (
      <div className="border-b p-4">
        <div className="flex items-center gap-3">
          <div className="text-lg font-semibold">Agent not found</div>
        </div>
      </div>
    );
  }

  return (
    <div className="border-b p-4">
      <div className="flex items-center justify-center gap-3">
        <h2 className="text-lg font-semibold">{agent.name}</h2>
        {agent.capabilities.length > 0 && (
          <Badge variant="outline" className="text-xs">
            {agent.capabilities.length} capabilities
          </Badge>
        )}
      </div>
    </div>
  );
}
