import { AgentId } from "../constants/agents";

export interface Message {
  id: string;
  content: string;
  timestamp: Date;
  sender: AgentId;
  type: "user" | "agent" | "inter_agent" | "system" | "reasoning";
  metadata?: {
    partial?: boolean;
    done?: boolean;
    invocationId?: string;
    eventType?: string;
    model?: string;
    created?: number;
    // Inter-agent specific metadata
    fromAgent?: string;
    toAgent?: string;
    interAgentType?: string;
    // Streaming-specific metadata
    streamingKey?: string;
    chunkIndex?: number;
    isChunk?: boolean;
    isCompletion?: boolean;
    finalizedAt?: string;
    hasStructuredThoughts?: boolean;
  };
  parts?: MessagePart[];
  usageMetadata?: {
    prompt_token_count?: number;
    candidates_token_count?: number;
    total_token_count?: number;
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
