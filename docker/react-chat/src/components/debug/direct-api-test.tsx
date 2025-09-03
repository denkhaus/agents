'use client'

import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { apiClient } from '@/lib/api/client'
import { agentApi } from '@/lib/api/agents'

export function DirectApiTest() {
  const [result, setResult] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const testAPI = async () => {
    setLoading(true)
    setError(null)
    setResult(null)

    try {
      console.log('Testing direct API call using apiClient...')
      
      // Test both raw API client and agent API
      console.log('=== Testing apiClient.getAgents() ===')
      const rawResult = await apiClient.getAgents()
      console.log('Raw API result:', rawResult)
      
      console.log('=== Testing agentApi.getAgents() ===')
      const agentResult = await agentApi.getAgents()
      console.log('Agent API result:', agentResult)
      
      setResult({
        raw: rawResult,
        processed: agentResult
      })
    } catch (err: any) {
      console.error('API test error:', err)
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  // Auto-test on mount
  useEffect(() => {
    testAPI()
  }, [])

  return (
    <Card className="m-4">
      <CardHeader>
        <CardTitle>Direct API Test</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <Button onClick={testAPI} disabled={loading}>
            {loading ? 'Testing...' : 'Test API'}
          </Button>
          
          {loading && <div>Loading...</div>}
          
          {error && (
            <div className="text-red-500">
              <strong>Error:</strong> {error}
            </div>
          )}
          
          {result && (
            <div className="text-green-500">
              <strong>Success!</strong> API calls completed:
              <pre className="mt-2 p-2 bg-gray-100 rounded text-black text-xs overflow-auto max-h-64">
                {JSON.stringify(result, null, 2)}
              </pre>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}