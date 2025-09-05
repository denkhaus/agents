import { AgentEvent, Message } from './'

export interface ConnectionStatus {
  connectionId: string
  isConnected: boolean
  lastConnected?: Date
  errorCount: number
  reconnectAttempts: number
}

export interface MessageCallback {
  (message: Message): void
}

export interface InterAgentCallback {
  (event: AgentEvent): void
}

export interface ConnectionCallback {
  (status: ConnectionStatus): void
}

export interface ErrorCallback {
  (error: Error): void
}

export interface StreamingConnection {
  id: string
  type: 'agent_run' | 'inter_agent' | 'system'
  agentId?: string
  sessionId: string
  eventSource: EventSource | null
  status: ConnectionStatus
  handlers: StreamingHandlers
}

export interface StreamingHandlers {
  onMessage?: (event: AgentEvent) => void
  onInterAgentEvent?: InterAgentCallback
  onConnectionChange?: ConnectionCallback
  onError?: ErrorCallback
}

export interface MessageProcessingContext {
  agentId: string
  sessionId: string
  connectionType: StreamingConnection['type']
  timestamp: Date
}

export interface MessageProcessor {
  canProcess(event: AgentEvent, context: MessageProcessingContext): boolean
  process(event: AgentEvent, context: MessageProcessingContext): Message | null
}

export interface StreamingMessageManagerConfig {
  maxReconnectAttempts: number
  reconnectInterval: number
  backoffMultiplier: number
  connectionTimeout: number
}

export interface SendMessageOptions {
  sessionId?: string
  userId?: string
  onProgress?: (event: AgentEvent) => void
  onError?: (error: Error) => void
}