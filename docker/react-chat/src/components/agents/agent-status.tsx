'use client'

import { Agent } from '@/lib/types'
import { Badge } from '@/components/ui/badge'

interface AgentStatusProps {
  status: Agent['status']
  size?: 'sm' | 'md' | 'lg'
}

export function AgentStatus({ status, size = 'md' }: AgentStatusProps) {
  const getStatusConfig = (status: Agent['status']) => {
    switch (status) {
      case 'online':
        return {
          color: 'bg-green-500',
          variant: 'default' as const,
          text: 'Online'
        }
      case 'busy':
        return {
          color: 'bg-yellow-500',
          variant: 'secondary' as const,
          text: 'Busy'
        }
      case 'offline':
        return {
          color: 'bg-gray-500',
          variant: 'outline' as const,
          text: 'Offline'
        }
      default:
        return {
          color: 'bg-gray-500',
          variant: 'outline' as const,
          text: 'Unknown'
        }
    }
  }

  const config = getStatusConfig(status)
  const sizeClasses = {
    sm: 'h-1.5 w-1.5',
    md: 'h-2 w-2',
    lg: 'h-3 w-3'
  }

  return (
    <div className="flex items-center gap-2">
      <div className={`rounded-full ${config.color} ${sizeClasses[size]}`} />
      <Badge variant={config.variant} className="text-xs">
        {config.text}
      </Badge>
    </div>
  )
}