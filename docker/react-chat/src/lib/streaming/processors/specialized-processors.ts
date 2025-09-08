import { LLMEvent, Message, EventType } from '@/lib/types'
import { MessageProcessingContext } from '@/lib/types/streaming'
import { BaseMessageProcessor } from './base-processor'
import { debug } from '@/lib/utils/debug'

// Legacy processors removed - using only new schema

export class InterAgentProcessor extends BaseMessageProcessor {
  constructor() {
    super('InterAgentProcessor')
  }

  canProcess(event: LLMEvent, context: MessageProcessingContext): boolean {
    return event.type === EventType.INTER_AGENT || !!event.inter_agent
  }

  process(event: LLMEvent, context: MessageProcessingContext): Message | null {
    const { content, parts } = this.extractContent(event)
    
    if (!content && (!parts || parts.length === 0)) {
      return null
    }

    const baseMessage = this.createBaseMessage(event, context, 'inter_agent')
    
    return {
      ...baseMessage,
      content,
      parts,
      usageMetadata: event.usage,
      metadata: {
        ...baseMessage.metadata,
        fromAgent: event.inter_agent?.from_agent,
        toAgent: event.inter_agent?.to_agent,
        interAgentType: event.inter_agent?.type,
      }
    } as Message
  }
}

export class ToolCallProcessor extends BaseMessageProcessor {
  constructor() {
    super('ToolCallProcessor')
  }

  canProcess(event: LLMEvent, context: MessageProcessingContext): boolean {
    // Handle tool calls and responses using new schema
    return event.type === EventType.TOOL_CALL || event.type === EventType.TOOL_RESPONSE
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
    return event.type === EventType.ASSISTANT || event.type === EventType.REASONING
  }

  process(event: LLMEvent, context: MessageProcessingContext): Message | null {
    const { content, parts } = this.extractContent(event)
    
    // Skip events without meaningful content
    if (!content && (!parts || parts.length === 0)) {
      return null
    }

    const messageType = event.type === EventType.REASONING ? 'reasoning' : 'agent'
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