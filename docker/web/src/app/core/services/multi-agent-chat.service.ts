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

import { Injectable, NgZone } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { URLUtil } from '../../../utils/url-util';
import { 
  AgentInfo, 
  MultiChatRequest, 
  MultiChatResponse, 
  InterAgentEvent 
} from '../models/MultiAgentChat';

@Injectable({
  providedIn: 'root'
})
export class MultiAgentChatService {
  private apiServerDomain = URLUtil.getApiServerBaseUrl();
  private eventsSubject = new BehaviorSubject<InterAgentEvent[]>([]);
  private agentsSubject = new BehaviorSubject<AgentInfo[]>([]);
  private eventSource?: EventSource;
  private isConnected = false;

  public events$ = this.eventsSubject.asObservable();
  public agents$ = this.agentsSubject.asObservable();

  constructor(
    private http: HttpClient,
    private zone: NgZone
  ) {}

  /**
   * Send a message between agents
   */
  sendMessage(request: MultiChatRequest): Observable<MultiChatResponse> {
    const url = `${this.apiServerDomain}/multi-chat/send`;
    return this.http.post<MultiChatResponse>(url, request);
  }

  /**
   * Connect to SSE stream for real-time inter-agent events
   */
  connectToEventStream(agents: string[], sessionId: string, userId: string): void {
    if (this.isConnected) {
      this.disconnect();
    }

    const agentList = agents.join(',');
    const url = `${this.apiServerDomain}/multi-chat/start_sse?agents=${encodeURIComponent(agentList)}&sessionId=${encodeURIComponent(sessionId)}&userId=${encodeURIComponent(userId)}`;
    
    this.eventSource = new EventSource(url);
    this.isConnected = true;

    this.eventSource.onmessage = (event) => {
      this.zone.run(() => {
        try {
          const data: InterAgentEvent = JSON.parse(event.data);
          this.handleIncomingEvent(data);
        } catch (error) {
          console.error('Error parsing SSE event:', error);
        }
      });
    };

    this.eventSource.onerror = (error) => {
      console.error('SSE connection error:', error);
      this.isConnected = false;
    };

    this.eventSource.onopen = () => {
      console.log('Multi-agent chat SSE connection established');
    };
  }

  /**
   * Disconnect from SSE stream
   */
  disconnect(): void {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = undefined;
    }
    this.isConnected = false;
  }

  /**
   * Get current connection status
   */
  getConnectionStatus(): boolean {
    return this.isConnected;
  }

  /**
   * Clear all events
   */
  clearEvents(): void {
    this.eventsSubject.next([]);
  }

  /**
   * Get current events
   */
  getCurrentEvents(): InterAgentEvent[] {
    return this.eventsSubject.value;
  }

  /**
   * Get current agents
   */
  getCurrentAgents(): AgentInfo[] {
    return this.agentsSubject.value;
  }

  private handleIncomingEvent(event: InterAgentEvent): void {
    switch (event.type) {
      case 'agent_list':
        if (event.agents) {
          this.agentsSubject.next(event.agents);
        }
        break;
      
      case 'inter_agent':
      case 'communication':
        this.addEvent(event);
        break;
      
      case 'heartbeat':
        // Handle heartbeat - could update connection status
        break;
      
      default:
        // Handle other event types or add to general events
        this.addEvent(event);
        break;
    }
  }

  private addEvent(event: InterAgentEvent): void {
    const currentEvents = this.eventsSubject.value;
    const updatedEvents = [...currentEvents, event];
    this.eventsSubject.next(updatedEvents);
  }
}