import { LLMEvent, Message } from '@/lib/types'
import { MessageProcessingContext } from '@/lib/types/streaming'
import { BaseMessageProcessor } from './base-processor'
import { debug } from '@/lib/utils/debug'

// Legacy processors removed - using only new schema

export class ToolCallProcessor extends BaseMessageProcessor {
  constructor() {
    super('ToolCallProcessor')
  }

  canProcess(event: LLMEvent, context: MessageProcessingContext): boolean {
    // Handle tool calls and responses using new schema
    return event.type === 'tool.call' || event.type === 'tool.response'
  }

  process(event: LLMEvent, context: MessageProcessingContext): Message | null {
    const { content, parts } = this.extractContent(event)
    
    const baseMessage = this.createBaseMessage(event, context, 'system')
    
    return {
      ...baseMessage,
      content: content || `Tool ${event.type}`,
      parts,
      usageMetadata: event.usage,
    } as Message
  }
}

export class AgentMessageProcessor extends BaseMessageProcessor {
  constructor() {
    super('AgentMessageProcessor')
  }

  canProcess(event: LLMEvent, context: MessageProcessingContext): boolean {
    // Handle regular agent messages using new schema
    return event.type === 'assistant' || event.type === 'reasoning'
  }

  process(event: LLMEvent, context: MessageProcessingContext): Message | null {
    const { content, parts } = this.extractContent(event)
    
    // Skip events without meaningful content
    if (!content && (!parts || parts.length === 0)) {
      return null
    }

    const messageType = event.type === 'reasoning' ? 'reasoning' : 'agent'
    const baseMessage = this.createBaseMessage(event, context, messageType)
    
    const message: Message = {
      ...baseMessage,
      content,
      parts,
      usageMetadata: event.usage,
    } as Message

    this.logProcessing(event, context, message)
    return message
  }
}

export class SystemMessageProcessor extends BaseMessageProcessor {
  constructor() {
    super('SystemMessageProcessor')
  }

  canProcess(event: LLMEvent, context: MessageProcessingContext): boolean {
    // Fallback processor for any unhandled events
    return true
  }

  process(event: LLMEvent, context: MessageProcessingContext): Message | null {
    // Log unhandled events for debugging
    debug.streaming(
      `${this.processorType}: Unhandled event type`,
      JSON.stringify({
        type: event.type,
        role: event.role,
        agentId: context.agentId
      })
    )
    
    return null
  }
}