import { AgentEvent, Message } from '@/lib/types'
import { MessageProcessingContext } from '@/lib/types/streaming'
import { BaseMessageProcessor } from './base-processor'
import { debug } from '@/lib/utils/debug'

export class UserResponseProcessor extends BaseMessageProcessor {
  constructor() {
    super('UserResponseProcessor')
  }

  canProcess(event: AgentEvent, context: MessageProcessingContext): boolean {
    // Handle regular agent responses to user messages
    // Support both legacy 'message' format and OpenAI completion formats
    // Also handle events without explicit object type but with content
    const isValidObject = (
      event.object === 'message' ||
      event.object === 'chat.completion.chunk' ||
      event.object === 'chat.completion' ||
      !event.object // Handle events without object type but with content
    )
    
    // Accept events with content or those that are explicitly marked as done/partial
    const hasRelevantContent = Boolean(
      event.content || 
      event.partial !== undefined || 
      event.done !== undefined ||
      event.invocationId
    )
    
    const result = (
      context.connectionType === 'agent_run' &&
      isValidObject &&
      hasRelevantContent &&
      !this.isInterAgentMessage(event) &&
      !this.isSystemMessage(event)
    )
    
    // Only log the final result for debugging
    // debug.streaming(
    //   `UserResponseProcessor: canProcess check for event`,
    //   JSON.stringify({
    //     eventType: event.type,
    //     eventObject: event.object,
    //     connectionType: context.connectionType,
    //     isValidObject,
    //     hasRelevantContent,
    //     isInterAgent: this.isInterAgentMessage(event),
    //     isSystem: this.isSystemMessage(event),
    //     result
    //   })
    // )
    
    return result
  }

  process(event: AgentEvent, context: MessageProcessingContext): Message | null {
    debug.streaming(
      `UserResponseProcessor: Processing event`,
      JSON.stringify({
        eventType: event.type,
        eventObject: event.object,
        hasContent: !!event.content,
        contentPreview: event.content ? String(event.content).substring(0, 100) : 'null',
        invocationId: event.invocationId,
        done: event.done,
        partial: event.partial
      })
    )
    
    // Handle final stream termination event for OpenAI completion
    // But still create message if there's content
    if (event.object === 'chat.completion' && event.done === true) {
      const { content, parts } = this.extractContent(event)
      // Only return null if there's no content to show
      if (!content && (!parts || parts.length === 0)) {
        debug.streaming('UserResponseProcessor: Returning null for completion event with no content')
        return null
      }
      // If there's content, process it like a regular message but mark as done
      debug.streaming('UserResponseProcessor: Processing completion event with content')
    }

    const { content, parts } = this.extractContent(event)
    
    debug.streaming(
      `UserResponseProcessor: Extracted content`,
      JSON.stringify({
        contentLength: content.length,
        hasParts: !!parts,
        contentPreview: content.substring(0, 100)
      })
    )
    
    // Generate unique message ID for each invocation
    // For streaming messages, use invocationId as base but ensure uniqueness
    const messageId = this.generateUniqueMessageId(event)
    
    const baseMessage = this.createBaseMessage(event, context, 'agent')
    
    // Determine if message should be treated as partial/streaming
    const isDone = event.done === true || event.object === 'chat.completion'
    const isPartial = event.partial === true || (event.object === 'chat.completion.chunk' && !isDone)
    
    const message: Message = {
      ...baseMessage,
      id: messageId, // Override with our unique ID
      content: content || '', // Ensure content is never undefined
      parts,
      metadata: {
        ...baseMessage.metadata,
        // Add OpenAI-specific metadata for streaming chunks
        isChunk: event.object === 'chat.completion.chunk',
        isCompletion: event.object === 'chat.completion',
        model: event.model,
        created: event.created,
        // Add streaming identifiers for proper accumulation
        streamingKey: event.invocationId || event.id,
        chunkIndex: this.generateChunkIndex(event),
        // Properly handle done state - mark as done if the event indicates completion
        done: isDone,
        partial: isPartial
      }
    } as Message

    this.logProcessing(event, context, message)
    debug.streaming(
      `UserResponseProcessor: Returning message`,
      JSON.stringify({
        messageId: message.id,
        contentLength: message.content.length,
        isPartial: message.metadata?.partial,
        isDone: message.metadata?.done,
        streamingKey: message.metadata?.streamingKey
      })
    )
    return message
  }

  private generateUniqueMessageId(event: AgentEvent): string {
    // For streaming messages, we need a stable ID that represents the entire message stream
    // Use invocationId as the base for accumulation, but add chunk info for uniqueness during processing
    if (event.invocationId) {
      return event.invocationId
    }
    
    // Fallback to event ID or generate new one
    return event.id || `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  }

  private generateChunkIndex(event: AgentEvent): number {
    // Generate a chunk index based on timestamp or sequence
    // This helps with ordering chunks if needed
    return event.timestamp ? Math.floor(event.timestamp * 1000) : Date.now()
  }

  private isInterAgentMessage(event: AgentEvent): boolean {
    return !!(
      event.type === 'inter_agent' ||
      event.type === 'communication' ||
      event.object === 'inter_agent' ||
      event.fromAgent ||
      event.toAgent
    )
  }

  private isSystemMessage(event: AgentEvent): boolean {
    return !!(
      event.object === 'tool_call' ||
      event.object === 'tool_response'
    )
  }

  // Override createBaseMessage to ensure proper streaming message handling
  protected createBaseMessage(
    event: AgentEvent, 
    context: MessageProcessingContext, 
    messageType: Message['type']
  ): Partial<Message> {
    const base = super.createBaseMessage(event, context, messageType)
    
    // Don't override the ID here - let the process method handle it
    return {
      ...base,
      metadata: {
        ...base.metadata,
        // Ensure streaming metadata is properly set
        partial: event.partial ?? (event.object === 'chat.completion.chunk'),
        done: event.done ?? (event.object === 'chat.completion'),
        streamingKey: event.invocationId || event.id
      }
    }
  }
}