'use client'

import { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'

interface ApiCall {
  id: string
  method: string
  url: string
  timestamp: Date
  status?: number
  response?: any
  error?: string
  duration?: number
}

export function ApiMonitor() {
  const [apiCalls, setApiCalls] = useState<ApiCall[]>([])
  const [isVisible, setIsVisible] = useState(true)

  useEffect(() => {
    // Override fetch to monitor API calls
    const originalFetch = window.fetch
    
    window.fetch = async (input, init) => {
      const url = typeof input === 'string' ? input : input.url
      const method = init?.method || 'GET'
      const startTime = Date.now()
      
      const callId = `${Date.now()}-${Math.random()}`
      
      // Add the call to our list
      const newCall: ApiCall = {
        id: callId,
        method,
        url,
        timestamp: new Date(),
      }
      
      setApiCalls(prev => [newCall, ...prev.slice(0, 19)]) // Keep last 20 calls
      
      try {
        const response = await originalFetch(input, init)
        const duration = Date.now() - startTime
        const responseClone = response.clone()
        
        try {
          const responseData = await responseClone.json()
          
          setApiCalls(prev => prev.map(call => 
            call.id === callId 
              ? { ...call, status: response.status, response: responseData, duration }
              : call
          ))
        } catch {
          // Response is not JSON
          setApiCalls(prev => prev.map(call => 
            call.id === callId 
              ? { ...call, status: response.status, duration }
              : call
          ))
        }
        
        return response
      } catch (error) {
        const duration = Date.now() - startTime
        setApiCalls(prev => prev.map(call => 
          call.id === callId 
            ? { ...call, error: error instanceof Error ? error.message : 'Unknown error', duration }
            : call
        ))
        throw error
      }
    }

    return () => {
      window.fetch = originalFetch
    }
  }, [])

  const clearCalls = () => setApiCalls([])

  if (!isVisible) {
    return (
      <div className="fixed bottom-4 right-4 z-50">
        <Button onClick={() => setIsVisible(true)} size="sm">
          Show API Monitor
        </Button>
      </div>
    )
  }

  return (
    <Card className="fixed bottom-4 right-4 w-96 max-h-96 z-50 bg-background border shadow-lg">
      <CardHeader className="pb-2">
        <div className="flex justify-between items-center">
          <CardTitle className="text-sm">API Monitor</CardTitle>
          <div className="flex gap-2">
            <Badge variant="outline">{apiCalls.length} calls</Badge>
            <Button onClick={clearCalls} size="sm" variant="outline">Clear</Button>
            <Button onClick={() => setIsVisible(false)} size="sm" variant="outline">Hide</Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="p-2">
        <ScrollArea className="h-64">
          {apiCalls.length === 0 ? (
            <div className="text-center text-muted-foreground py-4">
              No API calls yet
            </div>
          ) : (
            <div className="space-y-2">
              {apiCalls.map((call) => (
                <div key={call.id} className="border rounded p-2 text-xs">
                  <div className="flex justify-between items-start mb-1">
                    <div className="flex items-center gap-2">
                      <Badge variant={call.method === 'GET' ? 'default' : 'secondary'} className="text-xs">
                        {call.method}
                      </Badge>
                      {call.status && (
                        <Badge variant={call.status < 400 ? 'default' : 'destructive'} className="text-xs">
                          {call.status}
                        </Badge>
                      )}
                    </div>
                    <span className="text-muted-foreground">
                      {call.timestamp.toLocaleTimeString()}
                    </span>
                  </div>
                  <div className="font-mono text-xs break-all mb-1">
                    {call.url}
                  </div>
                  {call.error && (
                    <div className="text-red-500 text-xs">
                      Error: {call.error}
                    </div>
                  )}
                  {call.response && (
                    <details className="text-xs">
                      <summary className="cursor-pointer text-muted-foreground">Response</summary>
                      <pre className="mt-1 p-1 bg-muted rounded text-xs overflow-auto max-h-20">
                        {JSON.stringify(call.response, null, 2)}
                      </pre>
                    </details>
                  )}
                  {call.duration && (
                    <div className="text-muted-foreground text-xs">
                      Duration: {call.duration}ms
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </ScrollArea>
      </CardContent>
    </Card>
  )
}