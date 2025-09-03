'use client'

import { useChatStore } from '@/lib/store'
import { useAgentConnection } from '@/hooks/use-agent-connection'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export function DebugPanel() {
  const { agents: storeAgents, isConnected } = useChatStore()
  const { agents: hookAgents, isLoading, error, isError, isFetching } = useAgentConnection()

  return (
    <Card className="m-4">
      <CardHeader>
        <CardTitle>Debug Information</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-2 text-sm">
          <div>
            <strong>Hook Agents:</strong> {hookAgents.length} 
            {hookAgents.map(a => ` ${a.name}`).join(', ')}
          </div>
          <div>
            <strong>Store Agents:</strong> {storeAgents.length}
            {storeAgents.map(a => ` ${a.name}`).join(', ')}
          </div>
          <div>
            <strong>Is Loading:</strong> {isLoading ? 'Yes' : 'No'}
          </div>
          <div>
            <strong>Is Fetching:</strong> {isFetching ? 'Yes' : 'No'}
          </div>
          <div>
            <strong>Is Error:</strong> {isError ? 'Yes' : 'No'}
          </div>
          <div>
            <strong>Is Connected:</strong> {isConnected ? 'Yes' : 'No'}
          </div>
          <div>
            <strong>Error:</strong> {error ? error.message : 'None'}
          </div>
          <div>
            <strong>Backend URL:</strong> {process.env.NEXT_PUBLIC_API_URL || 'http://localhost:6999'}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}