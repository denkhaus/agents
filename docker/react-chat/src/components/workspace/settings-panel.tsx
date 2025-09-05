'use client'

import { useTheme } from '@/components/providers/theme-provider'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { ConnectionDebugger } from '@/components/debug/connection-debugger'
import { useState } from 'react'

export function SettingsPanel() {
  const { theme, setTheme } = useTheme()
  const [serverUrl, setServerUrl] = useState(process.env.NEXT_PUBLIC_API_URL || 'http://localhost:6999')
  const [autoConnect, setAutoConnect] = useState(true)
  const [soundEnabled, setSoundEnabled] = useState(true)

  const handleSaveSettings = () => {
    // Here you would typically save settings to localStorage or a backend
    console.log('Settings saved:', {
      serverUrl,
      autoConnect,
      soundEnabled,
      theme
    })
  }

  return (
    <div className="h-full p-6 space-y-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Settings</h2>
        <p className="text-muted-foreground">
          Configure your multi-agent chat preferences
        </p>
      </div>

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Appearance</CardTitle>
            <CardDescription>
              Customize the look and feel of the application
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="dark-mode">Dark mode</Label>
                <div className="text-sm text-muted-foreground">
                  Switch between light and dark themes
                </div>
              </div>
              <Switch
                id="dark-mode"
                checked={theme === 'dark'}
                onCheckedChange={(checked) => setTheme(checked ? 'dark' : 'light')}
              />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Connection</CardTitle>
            <CardDescription>
              Configure server connection settings
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="server-url">Server URL</Label>
              <Input
                id="server-url"
                value={serverUrl}
                onChange={(e) => setServerUrl(e.target.value)}
                placeholder="http://localhost:8080"
              />
              <div className="text-sm text-muted-foreground">
                The URL of your agent backend server
              </div>
            </div>
            
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="auto-connect">Auto-connect</Label>
                <div className="text-sm text-muted-foreground">
                  Automatically connect to agents when available
                </div>
              </div>
              <Switch
                id="auto-connect"
                checked={autoConnect}
                onCheckedChange={setAutoConnect}
              />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Notifications</CardTitle>
            <CardDescription>
              Configure notification and sound preferences
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="sound-enabled">Sound notifications</Label>
                <div className="text-sm text-muted-foreground">
                  Play sounds for new messages and events
                </div>
              </div>
              <Switch
                id="sound-enabled"
                checked={soundEnabled}
                onCheckedChange={setSoundEnabled}
              />
            </div>
          </CardContent>
        </Card>

        <Button onClick={handleSaveSettings} className="w-full">
          Save Settings
        </Button>
        
        {/* Debug Section */}
        <ConnectionDebugger />
      </div>
    </div>
  )
}