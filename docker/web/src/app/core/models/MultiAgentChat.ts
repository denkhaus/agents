/**
 * @license
 * Copyright 2025 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export interface AgentInfo {
  id: string;
  name: string;
  role: string;
}

export interface MultiChatRequest {
  fromAgent: string;
  toAgent: string;
  message: string;
  sessionId: string;
  userId: string;
}

export interface InterAgentEvent {
  type: string;
  fromAgent?: string;
  toAgent?: string;
  message?: string;
  timestamp: number;
  agents?: AgentInfo[];
  id?: string;
  invocationId?: string;
  author?: string;
  content?: {
    role: string;
    parts: Array<{text: string}>;
  };
  interAgent?: {
    fromAgent: string;
    toAgent: string;
    type: string;
  };
}

export interface MultiChatResponse {
  success: boolean;
  message: string;
  from: string;
  to: string;
}