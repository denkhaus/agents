/**
 * Agents Workspace Component
 * Displays active agents and their status
 */

import React from "react";
import { cn } from "@/lib/utils";

export const AgentsWorkspace: React.FC = () => (
  <div className="space-y-2">
    <h3 className="text-sm font-medium">Active Agents</h3>
    <div className="space-y-1">
      {[
        { name: "Designer", status: "online" },
        { name: "Frontend Dev", status: "busy" },
        { name: "Backend Dev", status: "online" },
        { name: "QA Engineer", status: "idle" },
        { name: "DevOps", status: "online" },
      ].map((agent) => (
        <div
          key={agent.name}
          className="flex items-center gap-2 p-2 rounded-md hover:bg-muted"
        >
          <div
            className={cn(
              "h-2 w-2 rounded-full",
              agent.status === "online" && "bg-green-500",
              agent.status === "busy" && "bg-yellow-500",
              agent.status === "idle" && "bg-gray-400"
            )}
          />
          <span className="text-xs">{agent.name}</span>
        </div>
      ))}
    </div>
  </div>
);