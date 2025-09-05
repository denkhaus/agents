export interface Message {
  id: string;
  content: string;
  timestamp: Date;
  sender: "user" | string; // 'user' or agent ID
  type: "user" | "agent" | "inter_agent" | "system";
  metadata?: {
    fromAgent?: string;
    toAgent?: string;
    eventType?: string;
    partial?: boolean;
    done?: boolean;
    invocationId?: string;
    addedToUI?: boolean; // Track if message has been added to UI
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
  agentId: string;
  messages: Message[];
  isActive: boolean;
  lastActivity: Date;
  sessionId?: string;
}

export interface TypingIndicator {
  agentId: string;
  isTyping: boolean;
}
