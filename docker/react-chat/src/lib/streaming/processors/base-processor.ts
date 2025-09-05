import { AgentEvent, Message } from '@/lib/types'
import { MessageProcessor, MessageProcessingContext } from '@/lib/types/streaming'
import { debug } from '@/lib/utils/debug'

export abstract class BaseMessageProcessor implements MessageProcessor {
  protected readonly processorType: string

  constructor(processorType: string) {
    this.processorType = processorType
  }

  abstract canProcess(event: AgentEvent, context: MessageProcessingContext): boolean
  
  abstract process(event: AgentEvent, context: MessageProcessingContext): Message | null

  protected createBaseMessage(
    event: AgentEvent, 
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
        fromAgent: event.fromAgent,
        toAgent: event.toAgent,
        eventType: event.type,
      }
    }
  }

  protected extractContent(event: AgentEvent): { content: string; parts?: unknown[] } {
    let content = ""
    let parts = undefined

    // Handle different content formats
    if (typeof event.content === 'object' && event.content !== null && 'parts' in event.content) {
      const eventContent = event.content as { parts?: Array<{ text?: string }> }
      if (Array.isArray(eventContent.parts)) {
        const textParts = eventContent.parts.filter(part => part.text)
        if (textParts.length > 0) {
          content = textParts.map(part => part.text).join('')
        }
        parts = eventContent.parts
      }
    } else if (typeof event.content === 'string') {
      content = event.content
    } else if (event.content && typeof event.content === 'object') {
      // Handle cases where content is an object with text property
      const contentObj = event.content as { text?: string };
      if (contentObj.text) {
        content = contentObj.text;
      } else {
        // Fallback: stringify the content object
        content = JSON.stringify(contentObj)
      }
    }

    // If no content extracted yet, try other event properties
    if (!content && event.message) {
      content = typeof event.message === 'string' ? event.message : JSON.stringify(event.message)
    }

    return { content, parts }
  }

  protected logProcessing(event: AgentEvent, context: MessageProcessingContext, result: Message | null) {
    debug.streaming(
      `${this.processorType}: Processing event`,
      JSON.stringify({
        type: event.type,
        object: event.object,
        fromAgent: event.fromAgent,
        toAgent: event.toAgent,
        agentId: context.agentId,
        processed: !!result
      })
    )
  }
}