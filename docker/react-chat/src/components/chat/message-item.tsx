'use client'

import { Message } from '@/lib/types'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { formatDistanceToNow } from 'date-fns'
import { Bot, User, ArrowRight } from 'lucide-react'

interface MessageItemProps {
  message: Message
}

export function MessageItem({ message }: MessageItemProps) {
  const isUser = message.sender === 'user'
  const isInterAgent = message.type === 'inter_agent'

  const formatMessageContent = () => {
    if (message.parts && message.parts.length > 0) {
      return message.parts.map((part, index) => (
        <div key={index}>
          {part.text && <p>{part.text}</p>}
          {part.functionCall && (
            <div className="mt-2 p-2 bg-muted rounded text-sm">
              <strong>Function Call:</strong> {part.functionCall.name}
              <pre className="mt-1 text-xs overflow-x-auto">
                {JSON.stringify(part.functionCall.args, null, 2)}
              </pre>
            </div>
          )}
          {part.functionResponse && (
            <div className="mt-2 p-2 bg-muted rounded text-sm">
              <strong>Function Response:</strong> {part.functionResponse.name}
              <pre className="mt-1 text-xs overflow-x-auto">
                {JSON.stringify(part.functionResponse.response, null, 2)}
              </pre>
            </div>
          )}
        </div>
      ))
    }
    return <p>{message.content}</p>
  }

  return (
    <div className={`flex gap-3 ${isUser ? 'flex-row-reverse' : 'flex-row'}`}>
      <Avatar className="h-8 w-8 flex-shrink-0">
        {isUser ? (
          <>
            <User className="h-4 w-4" />
            <AvatarFallback>U</AvatarFallback>
          </>
        ) : (
          <>
            <Bot className="h-4 w-4" />
            <AvatarFallback>
              {message.sender.slice(0, 2).toUpperCase()}
            </AvatarFallback>
          </>
        )}
      </Avatar>
      
      <Card className={`max-w-[80%] ${isUser ? 'bg-primary text-primary-foreground' : ''}`}>
        <CardContent className="p-3">
          <div className="flex items-center gap-2 mb-2">
            <span className="font-medium text-sm">
              {isUser ? 'You' : message.sender}
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
          
          {message.metadata?.partial && (
            <Badge variant="outline" className="mt-2 text-xs">
              Streaming...
            </Badge>
          )}
        </CardContent>
      </Card>
    </div>
  )
}