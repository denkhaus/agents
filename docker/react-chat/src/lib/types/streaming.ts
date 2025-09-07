import { AgentId } from "../constants/agents";
import { LLMEvent, Message } from "./";

export interface ConnectionStatus {
  connectionId: string;
  isConnected: boolean;
  lastConnected?: Date;
  errorCount: number;
  reconnectAttempts: number;
}

export interface MessageCallback {
  (message: Message): void;
}

export interface InterAgentCallback {
  (event: LLMEvent): void;
}

export interface ConnectionCallback {
  (status: ConnectionStatus): void;
}

export interface ErrorCallback {
  (error: Error): void;
}

export interface StreamingConnection {
  id: string;
  type: "agent_run" | "inter_agent" | "system";
  agentId?: AgentId;
  sessionId: string;
  eventSource: EventSource | null;
  status: ConnectionStatus;
  handlers: StreamingHandlers;
}

export interface StreamingHandlers {
  onMessage?: (event: LLMEvent) => void;
  onInterAgentEvent?: InterAgentCallback;
  onConnectionChange?: ConnectionCallback;
  onError?: ErrorCallback;
}

export interface MessageProcessingContext {
  agentId: AgentId;
  sessionId: string;
  connectionType: StreamingConnection["type"];
  timestamp: Date;
}

export interface MessageProcessor {
  canProcess(event: LLMEvent, context: MessageProcessingContext): boolean;
  process(event: LLMEvent, context: MessageProcessingContext): Message | null;
}

export interface StreamingMessageManagerConfig {
  maxReconnectAttempts: number;
  reconnectInterval: number;
  backoffMultiplier: number;
  connectionTimeout: number;
}

export interface SendMessageOptions {
  onProgress?: (event: LLMEvent) => void;
  onError?: (error: Error) => void;
}
