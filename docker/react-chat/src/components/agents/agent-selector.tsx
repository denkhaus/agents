"use client";

import { useChatStore } from "@/lib/store";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface AgentSelectorProps {
  value?: string;
  onValueChange?: (value: string) => void;
  placeholder?: string;
}

export function AgentSelector({
  value,
  onValueChange,
  placeholder = "Select an agent",
}: AgentSelectorProps) {
  const { agents } = useChatStore();

  return (
    <Select value={value} onValueChange={onValueChange}>
      <SelectTrigger className="w-full">
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {agents.map((agent) => (
          <SelectItem key={agent.id} value={agent.id}>
            <div className="flex items-center gap-2">
              <span>{agent.name}</span>
              <div
                className={`h-2 w-2 rounded-full ml-auto ${
                  agent.status === "online"
                    ? "bg-green-500"
                    : agent.status === "busy"
                    ? "bg-yellow-500"
                    : "bg-gray-500"
                }`}
              />
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
