import { AgentId } from "../constants/agents";

export interface Message {
  id: string;
  content: string;
  timestamp: Date;
  sender: AgentId;
  type: "user" | "agent" | "inter_agent" | "system" | "reasoning";
  metadata?: {
    fromAgent?: AgentId;
    toAgent?: AgentId;
    eventType?: string;
    partial?: boolean;
    done?: boolean;
    invocationId?: string;
    addedToUI?: boolean; // Track if message has been added to UI
    // Streaming-specific metadata
    streamingKey?: string; // Key for tracking streaming message accumulation
    chunkIndex?: number; // Index of the chunk in streaming sequence
    isChunk?: boolean; // Whether this is a streaming chunk
    isCompletion?: boolean; // Whether this is a completion event
    model?: string; // AI model used
    created?: number; // Creation timestamp
    finalizedAt?: string; // When the message was finalized
    hasStructuredThoughts?: boolean; // Whether the message contains structured thoughts
    object?: string; // Add this line
  };
  parts?: MessagePart[];
  usageMetadata?: {
    promptTokenCount?: number;
    candidatesTokenCount?: number;
    totalTokenCount?: number;
  };
}

export interface MessagePart {
  text?: string;
  functionCall?: {
    name: string;
    args: unknown;
    id?: string;
  };
  functionResponse?: {
    name: string;
    response: unknown;
    id?: string;
  };
}

export interface ChatSession {
  agentId: AgentId;
  messages: Message[];
  isActive: boolean;
  lastActivity: Date;
  sessionId?: string;
}

export interface TypingIndicator {
  agentId: AgentId;
  isTyping: boolean;
}
