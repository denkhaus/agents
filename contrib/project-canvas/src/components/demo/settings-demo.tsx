/**
 * Settings Demo Component
 * Demonstrates the new Settings Store functionality
 */

import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { useSettingsStore } from "@/stores";
import { useEnhancedSettingsSync } from "@/hooks/use-enhanced-settings-sync";

export const SettingsDemo: React.FC = () => {
  const settingsStore = useSettingsStore();
  const {
    settings,
    remoteSettings,
    updateTheme,
    updateSidebarState,
    updateApplicationSettings,
  } = useEnhancedSettingsSync();

  return (
    <div className="space-y-6 p-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Settings Store Demo</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Current State Display */}
          <div>
            <h3 className="font-medium mb-3">Current Settings State</h3>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-muted-foreground">Theme:</span>
                <Badge variant="outline" className="ml-2">
                  {settings.theme}
                </Badge>
              </div>
              <div>
                <span className="text-muted-foreground">Left Sidebar:</span>
                <Badge
                  variant={
                    settings.leftSidebarCollapsed ? "destructive" : "default"
                  }
                  className="ml-2"
                >
                  {settings.leftSidebarCollapsed ? "Collapsed" : "Expanded"}
                </Badge>
              </div>
              <div>
                <span className="text-muted-foreground">Right Sidebar:</span>
                <Badge
                  variant={
                    settings.rightSidebarCollapsed ? "destructive" : "default"
                  }
                  className="ml-2"
                >
                  {settings.rightSidebarCollapsed ? "Collapsed" : "Expanded"}
                </Badge>
              </div>
              <div>
                <span className="text-muted-foreground">Workspace:</span>
                <Badge variant="secondary" className="ml-2">
                  {settings.currentWorkspace}
                </Badge>
              </div>
              <div>
                <span className="text-muted-foreground">Selected Project:</span>
                <Badge variant="outline" className="ml-2">
                  {settings.selectedProjectId || "None"}
                </Badge>
              </div>
              <div>
                <span className="text-muted-foreground">Selected Nodes:</span>
                <Badge variant="outline" className="ml-2">
                  {settings.selectedNodeIds.length} nodes
                </Badge>
              </div>
            </div>
          </div>

          <Separator />

          {/* Theme Controls */}
          <div>
            <h3 className="font-medium mb-3">Theme Settings</h3>
            <div className="flex gap-2">
              <Button
                variant={settings.theme === "light" ? "default" : "outline"}
                size="sm"
                onClick={() => updateTheme("light")}
              >
                Light Mode
              </Button>
              <Button
                variant={settings.theme === "dark" ? "default" : "outline"}
                size="sm"
                onClick={() => updateTheme("dark")}
              >
                Dark Mode
              </Button>
            </div>
          </div>

          <Separator />

          {/* Sidebar Controls */}
          <div>
            <h3 className="font-medium mb-3">Sidebar Settings</h3>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm">Left Sidebar Collapsed</span>
                <Switch
                  checked={settings.leftSidebarCollapsed}
                  onCheckedChange={(checked) =>
                    updateSidebarState(checked, settings.rightSidebarCollapsed)
                  }
                />
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm">Right Sidebar Collapsed</span>
                <Switch
                  checked={settings.rightSidebarCollapsed}
                  onCheckedChange={(checked) =>
                    updateSidebarState(settings.leftSidebarCollapsed, checked)
                  }
                />
              </div>
            </div>
          </div>

          <Separator />

          {/* Application Settings */}
          <div>
            <h3 className="font-medium mb-3">Application Settings</h3>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm">Auto Save</span>
                <Switch
                  checked={settings.autoSave}
                  onCheckedChange={(checked) =>
                    updateApplicationSettings({ autoSave: checked })
                  }
                />
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm">Notifications</span>
                <Switch
                  checked={settings.notifications}
                  onCheckedChange={(checked) =>
                    updateApplicationSettings({ notifications: checked })
                  }
                />
              </div>
            </div>
          </div>

          <Separator />

          {/* Debug Information */}
          <div>
            <h3 className="font-medium mb-3">Debug Information</h3>
            <div className="text-xs space-y-2">
              <div>
                <span className="text-muted-foreground">Last Updated:</span>
                <span className="ml-2">
                  {new Date(settings.lastUpdated).toLocaleString()}
                </span>
              </div>
              <div>
                <span className="text-muted-foreground">
                  Remote Sync Status:
                </span>
                <Badge
                  variant={remoteSettings ? "default" : "destructive"}
                  className="ml-2"
                >
                  {remoteSettings ? "Connected" : "Disconnected"}
                </Badge>
              </div>
              <div>
                <span className="text-muted-foreground">UI Store Sync:</span>
                <Badge variant="default" className="ml-2">
                  Active
                </Badge>
              </div>
            </div>
          </div>

          <Separator />

          {/* Reset Controls */}
          <div>
            <h3 className="font-medium mb-3">Reset Settings</h3>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => {
                settingsStore.resetSettings();
              }}
            >
              Reset All Settings
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
