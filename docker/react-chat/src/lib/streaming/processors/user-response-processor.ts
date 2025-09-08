import { LLMEvent, Message, EventType } from '@/lib/types'
import { MessageProcessingContext } from "@/lib/types/streaming";
import { BaseMessageProcessor } from "./base-processor";
import { debug } from "@/lib/utils/debug";

export class UserResponseProcessor extends BaseMessageProcessor {
  constructor() {
    super("UserResponseProcessor");
  }

  canProcess(event: LLMEvent, context: MessageProcessingContext): boolean {
    // Handle regular agent responses - check if it's an assistant or reasoning message
    const isValidType = event.type === EventType.ASSISTANT || event.type === EventType.REASONING

    const hasContent = !!(event.parts && event.parts.length > 0);

    const result =
      context.connectionType === "agent_run" &&
      isValidType &&
      hasContent &&
      !this.isInterAgentMessage(event) &&
      !this.isSystemMessage(event);

    // Only log the final result for debugging
    // debug.streaming(
    //   `UserResponseProcessor: canProcess check for event`,
    //   JSON.stringify({
    //     eventType: event.type,
    //     connectionType: context.connectionType,
    //     isValidObject,
    //     hasRelevantContent,
    //     isInterAgent: this.isInterAgentMessage(event),
    //     isSystem: this.isSystemMessage(event),
    //     result
    //   })
    // )

    return result;
  }

  process(event: LLMEvent, context: MessageProcessingContext): Message | null {
    debug.streaming(
      `UserResponseProcessor: Processing event`,
      JSON.stringify({
        eventType: event.type,
        hasContent: !!event.content,
        contentPreview: event.content
          ? String(event.content).substring(0, 100)
          : "null",
        invocationId: event.invocationId,
        done: event.done,
        partial: event.partial,
      })
    );

    const { content, parts } = this.extractContent(event);

    debug.streaming(
      `UserResponseProcessor: Extracted content`,
      JSON.stringify({
        contentLength: content.length,
        hasParts: !!parts,
        contentPreview: content.substring(0, 100),
      })
    );

    // Generate unique message ID for each invocation
    // For streaming messages, use invocationId as base but ensure uniqueness
    const messageId = this.generateUniqueMessageId(event);

    const baseMessage = this.createBaseMessage(event, context, "agent");

    // Use the event's done/partial flags directly
    const isDone = event.done === true;
    const isPartial = event.partial === true;

    const message: Message = {
      ...baseMessage,
      id: messageId, // Override with our unique ID
      content: content || "", // Ensure content is never undefined
      parts,
      metadata: {
        ...baseMessage.metadata,
        streamingKey: event.invocationId || event.id,
        chunkIndex: this.generateChunkIndex(event),
        done: isDone,
        partial: isPartial,
      },
      usageMetadata: event.usage,
    } as Message;

    this.logProcessing(event, context, message);
    debug.streaming(
      `UserResponseProcessor: Returning message`,
      JSON.stringify({
        messageId: message.id,
        contentLength: message.content.length,
        isPartial: message.metadata?.partial,
        isDone: message.metadata?.done,
        streamingKey: message.metadata?.streamingKey,
      })
    );
    return message;
  }

  private generateUniqueMessageId(event: LLMEvent): string {
    // For streaming messages, we need a stable ID that represents the entire message stream
    // Use invocationId as the base for accumulation, but add chunk info for uniqueness during processing
    if (event.invocationId) {
      return event.invocationId;
    }

    // Fallback to event ID or generate new one
    return (
      event.id || `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    );
  }

  private generateChunkIndex(event: LLMEvent): number {
    // Generate a chunk index based on timestamp or sequence
    // This helps with ordering chunks if needed
    return event.timestamp ? Math.floor(event.timestamp * 1000) : Date.now();
  }

  private isInterAgentMessage(event: LLMEvent): boolean {
    return event.type === EventType.INTER_AGENT || !!event.inter_agent
  }

  private isSystemMessage(event: LLMEvent): boolean {
    return event.type === EventType.TOOL_CALL || event.type === EventType.TOOL_RESPONSE
  }
}
