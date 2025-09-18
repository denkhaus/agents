/**
 * Settings Canvas Component
 * Main settings interface for application configuration
 */

import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import {
  Settings,
  Palette,
  Bell,
  Keyboard,
  Database,
  Shield,
  Zap,
  Monitor,
} from "lucide-react";

interface SettingsCanvasProps {
  className?: string;
}

export const SettingsCanvas: React.FC<SettingsCanvasProps> = ({ className }) => {
  return (
    <div className={`h-full w-full p-6 overflow-y-auto ${className}`}>
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Header */}
        <div className="space-y-2">
          <h1 className="text-2xl font-bold flex items-center gap-2">
            <Settings className="h-6 w-6" />
            Application Settings
          </h1>
          <p className="text-muted-foreground">
            Configure your project canvas preferences and application behavior.
          </p>
        </div>

        <Separator />

        {/* Settings Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Display Settings */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Palette className="h-5 w-5" />
                Display
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Label htmlFor="dark-mode">Dark Mode</Label>
                <Switch id="dark-mode" />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="minimap">Show MiniMap</Label>
                <Switch id="minimap" defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="background">Canvas Background</Label>
                <Switch id="background" defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="auto-layout">Auto Layout</Label>
                <Switch id="auto-layout" defaultChecked />
              </div>
            </CardContent>
          </Card>

          {/* Notifications */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Bell className="h-5 w-5" />
                Notifications
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Label htmlFor="task-updates">Task Updates</Label>
                <Switch id="task-updates" defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="agent-status">Agent Status Changes</Label>
                <Switch id="agent-status" defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="project-changes">Project Changes</Label>
                <Switch id="project-changes" />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="system-alerts">System Alerts</Label>
                <Switch id="system-alerts" defaultChecked />
              </div>
            </CardContent>
          </Card>

          {/* Keyboard Shortcuts */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Keyboard className="h-5 w-5" />
                Keyboard Shortcuts
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm">New Project</span>
                <Badge variant="outline">Ctrl + N</Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm">Save</span>
                <Badge variant="outline">Ctrl + S</Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm">Zoom Fit</span>
                <Badge variant="outline">Ctrl + 0</Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm">Toggle Sidebar</span>
                <Badge variant="outline">Ctrl + B</Badge>
              </div>
              <Button variant="outline" size="sm" className="w-full mt-2">
                Customize Shortcuts
              </Button>
            </CardContent>
          </Card>

          {/* Performance */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Zap className="h-5 w-5" />
                Performance
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Label htmlFor="hardware-acceleration">Hardware Acceleration</Label>
                <Switch id="hardware-acceleration" defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="smooth-animations">Smooth Animations</Label>
                <Switch id="smooth-animations" defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="auto-save">Auto Save</Label>
                <Switch id="auto-save" defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="lazy-loading">Lazy Loading</Label>
                <Switch id="lazy-loading" />
              </div>
            </CardContent>
          </Card>

          {/* Data & Storage */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Database className="h-5 w-5" />
                Data & Storage
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Local Storage Used</span>
                  <span>2.4 MB / 10 MB</span>
                </div>
                <div className="w-full bg-secondary rounded-full h-2">
                  <div className="bg-primary h-2 rounded-full w-1/4"></div>
                </div>
              </div>
              <Button variant="outline" size="sm" className="w-full">
                Clear Cache
              </Button>
              <Button variant="outline" size="sm" className="w-full">
                Export Data
              </Button>
            </CardContent>
          </Card>

          {/* Privacy & Security */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Shield className="h-5 w-5" />
                Privacy & Security
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Label htmlFor="analytics">Usage Analytics</Label>
                <Switch id="analytics" />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="crash-reports">Crash Reports</Label>
                <Switch id="crash-reports" defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <Label htmlFor="data-encryption">Data Encryption</Label>
                <Switch id="data-encryption" defaultChecked />
              </div>
              <Button variant="outline" size="sm" className="w-full">
                Privacy Policy
              </Button>
            </CardContent>
          </Card>
        </div>

        {/* System Information */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Monitor className="h-5 w-5" />
              System Information
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
              <div>
                <span className="text-muted-foreground">Version</span>
                <p className="font-medium">1.0.0</p>
              </div>
              <div>
                <span className="text-muted-foreground">Build</span>
                <p className="font-medium">2024.01.20</p>
              </div>
              <div>
                <span className="text-muted-foreground">Environment</span>
                <p className="font-medium">Development</p>
              </div>
              <div>
                <span className="text-muted-foreground">Platform</span>
                <p className="font-medium">Web</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Action Buttons */}
        <div className="flex gap-4 pt-4">
          <Button>Save Settings</Button>
          <Button variant="outline">Reset to Defaults</Button>
          <Button variant="outline">Import Settings</Button>
        </div>
      </div>
    </div>
  );
};