"use client";

import { useState } from 'react';
import { ConnectionDebugger } from './connection-debugger';
import { DirectApiTest } from './direct-api-test';
import { ApiMonitor } from './api-monitor';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

export function DebugPanel() {
  const [isOpen, setIsOpen] = useState(false);

  if (!isOpen) {
    return (
      <Button
        onClick={() => setIsOpen(true)}
        className="fixed bottom-4 right-4 z-50"
        variant="outline"
      >
        Open Debug Panel
      </Button>
    );
  }

  return (
    <div className="fixed bottom-4 right-4 z-50 w-[500px] max-h-[80vh] overflow-y-auto">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-lg">Debug Panel</CardTitle>
          <Button 
            onClick={() => setIsOpen(false)} 
            variant="ghost" 
            size="sm"
          >
            Close
          </Button>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="connection" className="w-full">
            <TabsList className="grid w-full grid-cols-3">
              <TabsTrigger value="connection">Connection</TabsTrigger>
              <TabsTrigger value="api">API Test</TabsTrigger>
              <TabsTrigger value="monitor">Monitor</TabsTrigger>
            </TabsList>
            <TabsContent value="connection">
              <ConnectionDebugger />
            </TabsContent>
            <TabsContent value="api">
              <DirectApiTest />
            </TabsContent>
            <TabsContent value="monitor">
              <ApiMonitor />
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  );
}