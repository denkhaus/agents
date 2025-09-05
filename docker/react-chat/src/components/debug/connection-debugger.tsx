"use client";

import { useState, useEffect } from 'react';
import { useStreamingManager } from '@/hooks/use-streaming-manager';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ConnectionStatus } from '@/lib/types/streaming';

export function ConnectionDebugger() {
  const streamingManager = useStreamingManager();
  const [connectionStatus, setConnectionStatus] = useState<Record<string, ConnectionStatus>>({});
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    // Update connection status periodically
    const interval = setInterval(() => {
      const status = streamingManager.getConnectionStatus();
      setConnectionStatus(status);
      setIsConnected(streamingManager.isConnected());
    }, 1000);

    return () => clearInterval(interval);
  }, [streamingManager]);

  const handleTestConnection = () => {
    console.log('[DEBUGGER] Testing connection status');
    const status = streamingManager.getConnectionStatus();
    console.log('[DEBUGGER] Current connection status:', status);
    setConnectionStatus(status);
    setIsConnected(streamingManager.isConnected());
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Connection Debugger</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <span className="font-medium">Overall Connection Status:</span>
            <span className={isConnected ? "text-green-600" : "text-red-600"}>
              {isConnected ? "Connected" : "Disconnected"}
            </span>
          </div>
          
          <div className="space-y-2">
            <h3 className="font-medium">Connection Details:</h3>
            {Object.keys(connectionStatus).length > 0 ? (
              Object.entries(connectionStatus).map(([id, status]) => (
                <div key={id} className="flex items-center justify-between text-sm">
                  <span className="truncate">{id}</span>
                  <span className={status.isConnected ? "text-green-600" : "text-red-600"}>
                    {status.isConnected ? "Connected" : "Disconnected"}
                  </span>
                </div>
              ))
            ) : (
              <p className="text-muted-foreground text-sm">No active connections</p>
            )}
          </div>
          
          <Button onClick={handleTestConnection} variant="outline" size="sm">
            Refresh Connection Status
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}