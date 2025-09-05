"use client";

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { apiClient } from '@/lib/api';

export function DirectApiTest() {
  const [agentId, setAgentId] = useState('');
  const [message, setMessage] = useState('Hello, test message');
  const [response, setResponse] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleTestApi = async () => {
    if (!agentId.trim()) {
      setResponse('Please enter an agent ID');
      return;
    }

    setIsLoading(true);
    setResponse('Testing...');

    try {
      console.log('[DIRECT API TEST] Getting session');
      const session = await apiClient.createSession(agentId, 'testuser');
      console.log('[DIRECT API TEST] Created session:', session);

      console.log('[DIRECT API TEST] Running agent');
      const runResponse = await apiClient.runAgent({
        appName: agentId,
        userID: 'testuser',
        sessionID: session.id,
        streaming: false, // Try non-streaming first
        newMessage: {
          role: 'user',
          parts: [{ text: message }]
        }
      });

      console.log('[DIRECT API TEST] Run response:', runResponse);
      setResponse(JSON.stringify(runResponse, null, 2));
    } catch (error) {
      console.error('[DIRECT API TEST] Error:', error);
      setResponse(`Error: ${error instanceof Error ? error.message : String(error)}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Direct API Test</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <label className="text-sm font-medium">Agent ID</label>
          <Input
            value={agentId}
            onChange={(e) => setAgentId(e.target.value)}
            placeholder="Enter agent ID"
          />
        </div>
        
        <div className="space-y-2">
          <label className="text-sm font-medium">Message</label>
          <Textarea
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Enter test message"
            rows={3}
          />
        </div>
        
        <Button 
          onClick={handleTestApi} 
          disabled={isLoading || !agentId.trim()}
          className="w-full"
        >
          {isLoading ? 'Testing...' : 'Test API'}
        </Button>
        
        {response && (
          <div className="space-y-2">
            <label className="text-sm font-medium">Response</label>
            <Textarea
              value={response}
              readOnly
              placeholder="API response will appear here"
              rows={10}
              className="font-mono text-xs"
            />
          </div>
        )}
      </CardContent>
    </Card>
  );
}