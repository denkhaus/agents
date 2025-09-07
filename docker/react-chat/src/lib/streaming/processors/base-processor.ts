import { LLMEvent, Message, TextPart, FunctionCallPart, FunctionResponsePart } from '@/lib/types'
import { MessageProcessor, MessageProcessingContext } from '@/lib/types/streaming'
import { debug } from '@/lib/utils/debug'

export abstract class BaseMessageProcessor implements MessageProcessor {
  protected readonly processorType: string

  constructor(processorType: string) {
    this.processorType = processorType
  }

  abstract canProcess(event: LLMEvent, context: MessageProcessingContext): boolean
  
  abstract process(event: LLMEvent, context: MessageProcessingContext): Message | null

  protected createBaseMessage(
    event: LLMEvent, 
    context: MessageProcessingContext, 
    messageType: Message['type']
  ): Partial<Message> {
    const messageId = event.id || event.invocationId || `${context.agentId}-${Date.now()}`
    
    return {
      id: messageId,
      timestamp: new Date((event.timestamp || Date.now() / 1000) * 1000),
      sender: event.author || context.agentId,
      type: messageType,
      metadata: {
        invocationId: event.invocationId,
        partial: event.partial || false,
        done: event.done || false,
        eventType: event.type,
        model: event.model,
        created: event.created,
      }
    }
  }

  protected extractContent(event: LLMEvent): { content: string; parts?: unknown[] } {
    let content = ""
    let parts = undefined

    // Extract content from LLMEvent parts
    if (event.parts && Array.isArray(event.parts)) {
      const textParts: string[] = []
      
      for (const part of event.parts) {
        if ('content' in part && typeof part.content === 'string') {
          // TextPart
          textParts.push(part.content)
        }
      }
      
      if (textParts.length > 0) {
        content = textParts.join('')
      }
      parts = event.parts
    }

    return { content, parts }
  }

  protected logProcessing(event: LLMEvent, context: MessageProcessingContext, result: Message | null) {
    debug.streaming(
      `${this.processorType}: Processing event`,
      JSON.stringify({
        type: event.type,
        role: event.role,
        agentId: context.agentId,
        processed: !!result
      })
    )
  }
}