/**
 * Settings Workspace Component
 * Displays application preferences and settings navigation
 */

import React from "react";
import { Button } from "@/components/ui/button";
import { ChevronRight } from "lucide-react";

export const SettingsWorkspace: React.FC = () => (
  <div className="space-y-2">
    <h3 className="text-sm font-medium">Preferences</h3>
    <div className="space-y-1">
      {["Display", "Notifications", "Keyboard", "Advanced"].map((setting) => (
        <Button
          key={setting}
          variant="ghost"
          className="w-full justify-between h-8 px-2"
        >
          <span className="text-xs">{setting}</span>
          <ChevronRight className="h-3 w-3" />
        </Button>
      ))}
    </div>
  </div>
);