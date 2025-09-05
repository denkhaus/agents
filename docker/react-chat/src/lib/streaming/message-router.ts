import { AgentEvent, Message } from '@/lib/types'
import { MessageProcessor, MessageProcessingContext } from '@/lib/types/streaming'
import { UserResponseProcessor } from './processors/user-response-processor'
import { AgentListProcessor, InterAgentProcessor, ToolCallProcessor, SystemMessageProcessor } from './processors/specialized-processors'
import { debug } from '@/lib/utils/debug'

export class MessageEventRouter {
  private processors: MessageProcessor[]

  constructor() {
    // Order matters: more specific processors first
    this.processors = [
      new AgentListProcessor(),
      new InterAgentProcessor(),
      new ToolCallProcessor(),
      new UserResponseProcessor(),
      new SystemMessageProcessor(), // Fallback processor
    ]
  }

  processEvent(event: AgentEvent, context: MessageProcessingContext): Message | null {
    // Find the first processor that can handle this event
    for (const processor of this.processors) {
      if (processor.canProcess(event, context)) {
        try {
          const result = processor.process(event, context)
          if (result) {
            debug.streaming(
              'MessageRouter: Event processed successfully',
              JSON.stringify({
                processor: processor.constructor.name,
                eventType: event.type,
                eventObject: event.object,
                messageId: result.id,
                messageType: result.type
              })
            )
          }
          return result
        } catch (error) {
          debug.error(
            `MessageRouter: Error processing event with ${processor.constructor.name}:`,
            error,
            event
          )
          // Continue to next processor on error
        }
      }
      // Only log when processors can't process for debugging
      // debug.streaming(
      //   'MessageRouter: Processor cannot process event',
      //   JSON.stringify({
      //     processor: processor.constructor.name,
      //     eventType: event.type,
      //     eventObject: event.object,
      //     agentId: context.agentId
      //   })
      // )
    }

    // No processor could handle this event
    debug.warn(
      'MessageRouter: No processor found for event',
      JSON.stringify({
        type: event.type,
        object: event.object,
        fromAgent: event.fromAgent,
        toAgent: event.toAgent,
        agentId: context.agentId
      })
    )
    
    return null
  }

  addProcessor(processor: MessageProcessor) {
    this.processors.unshift(processor) // Add to front for priority
  }

  removeProcessor(processorClass: new () => MessageProcessor) {
    this.processors = this.processors.filter(
      p => !(p instanceof processorClass)
    )
  }
}