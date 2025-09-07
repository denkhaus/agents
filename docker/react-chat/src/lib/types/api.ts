import { AgentId } from "../constants/agents";
import { AgentInfo } from "./agent";

export interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: string;
}

// Event types matching server schema
export type EventType = 
  | "assistant"
  | "tool.call" 
  | "tool.response"
  | "reasoning"
  | "inter_agent";

// Part types matching server schema
export interface TextPart {
  content: string;
}

export interface FunctionCallPart {
  name: string;
  args: unknown;
  id?: string;
}

export interface FunctionResponsePart {
  name: string;
  args?: unknown;
  id?: string;
  response: unknown;
}

export interface InterAgentPart {
  from_agent_id: string;
  to_agent_id: string;
  message: string;
  direction: "sent" | "received";
}

export type Part = TextPart | FunctionCallPart | FunctionResponsePart | InterAgentPart;

// Usage metadata matching server schema
export interface UsageMetaData {
  prompt_token_count?: number;
  candidates_token_count?: number;
  total_token_count?: number;
}

// Incoming request structures
export interface PartIncoming {
  text?: string;
  inlineData?: {
    data: string;
    mimeType: string;
    displayName?: string;
  };
  functionCall?: {
    name: string;
    args?: Record<string, unknown>;
  };
  functionResponse?: {
    name: string;
    response: unknown;
    id?: string;
  };
}

export interface Content {
  role: string;
  parts: PartIncoming[];
}

export interface AgentRunRequest {
  appName: string;
  fromAgentId: string;
  toAgentId: string;
  sessionId: string;
  content: Content;
  streaming: boolean;
}

// Session structure matching server schema
export interface ADKSession {
  appName: string;
  agentId: string;
  id: string;
  createTime: number;
  lastUpdateTime: number;
  state: Record<string, Uint8Array>;
  events: LLMEvent[];
}

// Inter-agent communication data
export interface InterAgentData {
  from_agent: string;
  to_agent: string;
  type: "communication" | "received";
}

// Main event structure matching server LLMEvent
export interface LLMEvent {
  usage?: UsageMetaData;
  done?: boolean;
  partial?: boolean;
  type?: EventType;
  created?: number;
  model?: string;
  role?: string;
  parts?: Part[];
  timestamp?: number;
  id?: string;
  invocationId?: string;
  author?: string;
  inter_agent?: InterAgentData;
}


export interface SSEEventData {
  type: string;
  data: LLMEvent;
}
