"use client";

import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export function ApiMonitor() {
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    // Override console.log to capture API logs
    const originalLog = console.log;
    const originalError = console.error;
    
    console.log = function(...args) {
      // Capture API-related logs
      const message = args.map(arg => 
        typeof arg === 'object' ? JSON.stringify(arg) : String(arg)
      ).join(' ');
      
      if (message.includes('API') || message.includes('apiClient')) {
        setLogs(prev => [...prev.slice(-50), `[LOG] ${message}`]);
      }
      
      originalLog.apply(console, args);
    };
    
    console.error = function(...args) {
      // Capture API-related errors
      const message = args.map(arg => 
        typeof arg === 'object' ? JSON.stringify(arg) : String(arg)
      ).join(' ');
      
      if (message.includes('API') || message.includes('apiClient')) {
        setLogs(prev => [...prev.slice(-50), `[ERROR] ${message}`]);
      }
      
      originalError.apply(console, args);
    };
    
    return () => {
      console.log = originalLog;
      console.error = originalError;
    };
  }, []);

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>API Monitor</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="h-64 overflow-y-auto bg-muted p-2 rounded">
          {logs.length > 0 ? (
            <pre className="text-xs font-mono">
              {logs.map((log, index) => (
                <div 
                  key={index} 
                  className={log.includes('ERROR') ? "text-red-500" : "text-muted-foreground"}
                >
                  {log}
                </div>
              ))}
            </pre>
          ) : (
            <p className="text-muted-foreground text-sm">No API logs yet...</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}