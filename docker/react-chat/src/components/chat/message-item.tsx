'use client'

import { Message } from '@/lib/types'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { formatDistanceToNow } from 'date-fns'
import { Bot, User, ArrowRight, Settings, Play, CheckCircle } from 'lucide-react'

interface MessageItemProps {
  message: Message
}

export function MessageItem({ message }: MessageItemProps) {
  const isUser = message.sender === 'user'
  const isInterAgent = message.type === 'inter_agent'
  const isSystem = message.type === 'system'
  const isToolCall = message.parts?.some(part => part.functionCall)
  const isToolResponse = message.parts?.some(part => part.functionResponse)

  const formatMessageContent = () => {
    if (message.parts && message.parts.length > 0) {
      return message.parts.map((part, index) => (
        <div key={index}>
          {part.text && <p className="whitespace-pre-wrap">{part.text}</p>}
          {part.functionCall && (
            <div className="mt-2 p-3 bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <Play className="h-4 w-4 text-blue-600" />
                <strong className="text-blue-800 dark:text-blue-200">Tool Call: {part.functionCall.name}</strong>
                {part.functionCall.id && (
                  <Badge variant="outline" className="text-xs">
                    {part.functionCall.id}
                  </Badge>
                )}
              </div>
              <details className="mt-2">
                <summary className="cursor-pointer text-sm text-blue-600 hover:text-blue-800">
                  View Parameters
                </summary>
                <pre className="mt-2 text-xs overflow-x-auto bg-white dark:bg-gray-900 p-2 rounded border">
                  {JSON.stringify(part.functionCall.args, null, 2)}
                </pre>
              </details>
            </div>
          )}
          {part.functionResponse && (
            <div className="mt-2 p-3 bg-green-50 dark:bg-green-950 border border-green-200 dark:border-green-800 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <CheckCircle className="h-4 w-4 text-green-600" />
                <strong className="text-green-800 dark:text-green-200">Tool Response: {part.functionResponse.name}</strong>
                {part.functionResponse.id && (
                  <Badge variant="outline" className="text-xs">
                    {part.functionResponse.id}
                  </Badge>
                )}
              </div>
              <details className="mt-2">
                <summary className="cursor-pointer text-sm text-green-600 hover:text-green-800">
                  View Response
                </summary>
                <pre className="mt-2 text-xs overflow-x-auto bg-white dark:bg-gray-900 p-2 rounded border max-h-40">
                  {JSON.stringify(part.functionResponse.response, null, 2)}
                </pre>
              </details>
            </div>
          )}
        </div>
      ))
    }
    return <p className="whitespace-pre-wrap">{message.content}</p>
  }

  const getAvatarIcon = () => {
    if (isUser) return <User className="h-4 w-4" />
    if (isToolCall) return <Play className="h-4 w-4 text-blue-600" />
    if (isToolResponse) return <CheckCircle className="h-4 w-4 text-green-600" />
    if (isSystem) return <Settings className="h-4 w-4 text-orange-600" />
    return <Bot className="h-4 w-4" />
  }

  const getCardStyle = () => {
    if (isUser) return 'bg-primary text-primary-foreground'
    if (isToolCall) return 'border-blue-200 dark:border-blue-800'
    if (isToolResponse) return 'border-green-200 dark:border-green-800'
    if (isSystem) return 'border-orange-200 dark:border-orange-800'
    return ''
  }

  const getMessageTypeLabel = () => {
    if (isUser) return 'You'
    if (isToolCall) return `${message.sender} (Tool Call)`
    if (isToolResponse) return `${message.sender} (Tool Response)`
    if (isSystem) return `${message.sender} (System)`
    return message.sender
  }

  return (
    <div className={`flex gap-3 ${isUser ? 'flex-row-reverse' : 'flex-row'}`}>
      <Avatar className="h-8 w-8 flex-shrink-0">
        {getAvatarIcon()}
        <AvatarFallback>
          {isUser ? 'U' : message.sender.slice(0, 2).toUpperCase()}
        </AvatarFallback>
      </Avatar>
      
      <Card className={`max-w-[80%] ${getCardStyle()}`}>
        <CardContent className="p-3">
          <div className="flex items-center gap-2 mb-2">
            <span className="font-medium text-sm">
              {getMessageTypeLabel()}
            </span>
            
            {isInterAgent && message.metadata?.fromAgent && message.metadata?.toAgent && (
              <div className="flex items-center gap-1 text-xs">
                <Badge variant="secondary" className="px-1 py-0">
                  {message.metadata.fromAgent}
                </Badge>
                <ArrowRight className="h-3 w-3" />
                <Badge variant="secondary" className="px-1 py-0">
                  {message.metadata.toAgent}
                </Badge>
              </div>
            )}
            
            <span className="text-xs opacity-70 ml-auto">
              {formatDistanceToNow(message.timestamp, { addSuffix: true })}
            </span>
          </div>
          
          <div className="text-sm">
            {formatMessageContent()}
          </div>
          
          <div className="flex gap-2 mt-2">
            {message.metadata?.partial && (
              <Badge variant="outline" className="text-xs">
                Streaming...
              </Badge>
            )}
            {isToolCall && (
              <Badge variant="secondary" className="text-xs bg-blue-100 text-blue-800">
                Tool Call
              </Badge>
            )}
            {isToolResponse && (
              <Badge variant="secondary" className="text-xs bg-green-100 text-green-800">
                Tool Response
              </Badge>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}