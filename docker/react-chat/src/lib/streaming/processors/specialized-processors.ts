import { AgentEvent, Message } from '@/lib/types'
import { MessageProcessingContext } from '@/lib/types/streaming'
import { BaseMessageProcessor } from './base-processor'
import { debug } from '@/lib/utils/debug'
import { normalizeToAgentId, normalizeToAgentName } from '@/lib/constants/agents'

export class AgentListProcessor extends BaseMessageProcessor {
  constructor() {
    super('AgentListProcessor')
  }

  canProcess(event: AgentEvent, context: MessageProcessingContext): boolean {
    // Handle agent list events
    return !!(
      event.type === 'agent_list' ||
      event.object === 'agent_list'
    )
  }

  process(event: AgentEvent, context: MessageProcessingContext): Message | null {
    // Agent list events are typically control messages that don't need to be displayed
    // Return null to indicate no message should be created
    // But log for debugging
    debug.streaming(
      `${this.processorType}: Agent list event processed`,
      JSON.stringify({
        type: event.type,
        object: event.object,
        agentId: context.agentId
      })
    )
    
    return null
  }
}

export class InterAgentProcessor extends BaseMessageProcessor {
  constructor() {
    super('InterAgentProcessor')
  }

  canProcess(event: AgentEvent, context: MessageProcessingContext): boolean {
    // Handle inter-agent communication events
    return !!(
      event.type === 'inter_agent' ||
      event.type === 'communication' ||
      event.object === 'inter_agent' ||
      (event.fromAgent && event.toAgent) ||
      (event.object === 'message' && (event.fromAgent || event.toAgent))
    )
  }

  process(event: AgentEvent, context: MessageProcessingContext): Message | null {
    const { content, parts } = this.extractContent(event)
    
    // Skip events without meaningful content
    if (!content && (!parts || parts.length === 0)) {
      return null
    }

    // Determine the correct agent for this message based on event type
    const rawFromAgent = event.fromAgent || event.interAgent?.fromAgent || event.author
    const rawToAgent = event.toAgent || event.interAgent?.toAgent
    const eventType = event.interAgent?.type || event.type || 'inter_agent'
    
    // Normalize agent identifiers to IDs
    const fromAgent = rawFromAgent ? (normalizeToAgentId(rawFromAgent) || rawFromAgent) : undefined
    const toAgent = rawToAgent ? (normalizeToAgentId(rawToAgent) || rawToAgent) : undefined
    const contextAgentId = normalizeToAgentId(context.agentId) || context.agentId
    
    // For "received" events, only show in the target agent's window
    if (eventType === 'received' && toAgent && contextAgentId !== toAgent) {
      return null
    }
    
    // For "communication" events, only show in the sender's window
    if (eventType === 'communication' && fromAgent && contextAgentId !== fromAgent) {
      return null
    }
    
    // For other inter-agent events, check if this agent is involved
    if (fromAgent && toAgent && contextAgentId !== fromAgent && contextAgentId !== toAgent) {
      return null
    }

    const baseMessage = this.createBaseMessage(event, context, 'inter_agent')
    
    const message: Message = {
      ...baseMessage,
      content,
      parts,
      metadata: {
        ...baseMessage.metadata,
        fromAgent,
        toAgent,
        eventType,
      }
    } as Message

    this.logProcessing(event, context, message)
    return message
  }
}

export class ToolCallProcessor extends BaseMessageProcessor {
  constructor() {
    super('ToolCallProcessor')
  }

  canProcess(event: AgentEvent, context: MessageProcessingContext): boolean {
    // Handle tool calls and responses
    return !!(
      event.object === 'tool_call' ||
      event.object === 'tool.call' ||  // Handle dot notation
      event.object === 'tool_response' ||
      event.object === 'tool.response' ||  // Handle dot notation
      event.object === 'tool_code'  // Handle tool code events
    )
  }

  process(event: AgentEvent, context: MessageProcessingContext): Message | null {
    const { content, parts } = this.extractContent(event)
    
    // Even tool calls without visible content should be processed for system tracking
    const baseMessage = this.createBaseMessage(event, context, 'system')
    
    // Check if this is a tool call with function call data
    let messageParts = parts;
    if (event.object === 'tool_code') {
      const code = (event.content as any)?.code || (event as any).code;
      if (code) {
        messageParts = [{
          tool_code: { code }
        }];
      }
    } else if (!messageParts && event.content && typeof event.content === 'object') {
      const contentObj = event.content as any;
      
      // Handle function call data
      if (contentObj.functionCall) {
        messageParts = [{
          functionCall: contentObj.functionCall
        }];
      } else if (contentObj.functionResponse) {
        messageParts = [{
          functionResponse: contentObj.functionResponse
        }];
      }
    }
    
    const message: Message = {
      ...baseMessage,
      content: content || `Tool ${event.object || event.type}`,
      parts: messageParts,
      metadata: {
        ...baseMessage.metadata,
        toolCallType: event.object || event.type,
      }
    } as Message

    this.logProcessing(event, context, message)
    return message
  }
}

export class AgentMessageProcessor extends BaseMessageProcessor {
  constructor() {
    super('AgentMessageProcessor')
  }

  canProcess(event: AgentEvent, context: MessageProcessingContext): boolean {
    // Handle regular agent messages
    return !!(
      event.object === 'message' ||
      (event.type === 'agent' && !event.fromAgent && !event.toAgent)
    )
  }

  process(event: AgentEvent, context: MessageProcessingContext): Message | null {
    const { content, parts } = this.extractContent(event)
    
    // Skip events without meaningful content
    if (!content && (!parts || parts.length === 0)) {
      return null
    }

    // Normalize agent identifiers to IDs
    const eventAgentId = (event as any).agentId || event.author
    const normalizedAgentId = normalizeToAgentId(eventAgentId) || eventAgentId
    
    // Only process if this message is for the current agent context
    const contextAgentId = normalizeToAgentId(context.agentId) || context.agentId
    if (normalizedAgentId !== contextAgentId) {
      return null
    }

    const baseMessage = this.createBaseMessage(event, context, 'agent')
    
    const message: Message = {
      ...baseMessage,
      content,
      parts,
      sender: normalizedAgentId,
      metadata: {
        ...baseMessage.metadata,
        fromAgent: normalizedAgentId,
      }
    } as Message

    this.logProcessing(event, context, message)
    return message
  }
}

export class SystemMessageProcessor extends BaseMessageProcessor {
  constructor() {
    super('SystemMessageProcessor')
  }

  canProcess(event: AgentEvent, context: MessageProcessingContext): boolean {
    // Handle system messages, heartbeats, and other control messages
    return !!(
      event.type === 'system' ||
      event.type === 'heartbeat' ||
      event.object === 'system' ||
      (!event.type && !event.object) // fallback for untyped system events
    )
  }

  process(event: AgentEvent, context: MessageProcessingContext): Message | null {
    // Most system events don't need to be displayed as messages
    // Return null to indicate no message should be created
    // But log for debugging
    debug.streaming(
      `${this.processorType}: System event processed`,
      JSON.stringify({
        type: event.type,
        object: event.object,
        agentId: context.agentId
      })
    )
    
    return null
  }
}