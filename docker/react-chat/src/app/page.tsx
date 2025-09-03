'use client'

import { MainLayout } from '@/components/layout/main-layout'
import { useAgentConnection } from '@/hooks/use-agent-connection'
import { useEffect } from 'react'

export default function Home() {
  // Initialize agent connection
  const connectionResult = useAgentConnection()
  
  useEffect(() => {
    console.log('Home component mounted')
    console.log('Connection result:', connectionResult)
  }, [connectionResult])

  return <MainLayout />
}
