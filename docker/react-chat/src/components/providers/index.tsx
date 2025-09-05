'use client'

import { ThemeProvider } from './theme-provider'
import { QueryProvider } from './query-provider'
import { Toaster } from '@/components/ui/sonner'
import { DebugPanel } from '@/components/debug/debug-panel'

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider>
      <QueryProvider>
        {children}
        <Toaster />
        <DebugPanel />
      </QueryProvider>
    </ThemeProvider>
  )
}