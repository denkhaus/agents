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

import { Component, OnInit, OnDestroy, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subscription } from 'rxjs';
import { MultiAgentChatService } from '../../core/services/multi-agent-chat.service';
import { AgentService } from '../../core/services/agent.service';
import { 
  AgentInfo, 
  MultiChatRequest, 
  InterAgentEvent 
} from '../../core/models/MultiAgentChat';

@Component({
  selector: 'app-multi-agent-chat',
  templateUrl: './multi-agent-chat.component.html',
  styleUrls: ['./multi-agent-chat.component.scss'],
  imports: [CommonModule, FormsModule]
})
export class MultiAgentChatComponent implements OnInit, OnDestroy {
  @Input() sessionId: string = '';
  @Input() userId: string = 'user';
  @Input() availableAgents: string[] = [];

  agents: AgentInfo[] = [];
  events: InterAgentEvent[] = [];
  selectedFromAgent: string = 'user';
  selectedToAgent: string = '';
  messageText: string = '';
  isConnected: boolean = false;
  isLoading: boolean = false;

  private subscriptions: Subscription[] = [];

  constructor(
    private multiAgentChatService: MultiAgentChatService,
    private agentService: AgentService
  ) {}

  ngOnInit(): void {
    this.initializeComponent();
  }

  ngOnDestroy(): void {
    this.subscriptions.forEach(sub => sub.unsubscribe());
    this.multiAgentChatService.disconnect();
  }

  private initializeComponent(): void {
    // Subscribe to agents
    const agentsSubscription = this.multiAgentChatService.agents$.subscribe(
      agents => {
        this.agents = agents;
        if (agents.length > 0 && !this.selectedToAgent) {
          this.selectedToAgent = agents[0].name;
        }
      }
    );
    this.subscriptions.push(agentsSubscription);

    // Subscribe to events
    const eventsSubscription = this.multiAgentChatService.events$.subscribe(
      events => {
        this.events = events;
      }
    );
    this.subscriptions.push(eventsSubscription);

    // Connect to event stream if agents are available
    if (this.availableAgents.length > 0) {
      this.connectToEventStream();
    }
  }

  connectToEventStream(): void {
    if (this.availableAgents.length === 0 || !this.sessionId) {
      console.warn('Cannot connect: missing agents or session ID');
      return;
    }

    this.multiAgentChatService.connectToEventStream(
      this.availableAgents,
      this.sessionId,
      this.userId
    );
    this.isConnected = true;
  }

  disconnect(): void {
    this.multiAgentChatService.disconnect();
    this.isConnected = false;
  }

  sendMessage(): void {
    if (!this.messageText.trim() || !this.selectedToAgent) {
      return;
    }

    this.isLoading = true;

    const request: MultiChatRequest = {
      fromAgent: this.selectedFromAgent,
      toAgent: this.selectedToAgent,
      message: this.messageText.trim(),
      sessionId: this.sessionId,
      userId: this.userId
    };

    const sendSubscription = this.multiAgentChatService.sendMessage(request).subscribe({
      next: (response) => {
        console.log('Message sent successfully:', response);
        this.messageText = '';
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error sending message:', error);
        this.isLoading = false;
      }
    });

    this.subscriptions.push(sendSubscription);
  }

  clearEvents(): void {
    this.multiAgentChatService.clearEvents();
  }

  getEventDisplayText(event: InterAgentEvent): string {
    if (event.content && event.content.parts && event.content.parts.length > 0) {
      return event.content.parts[0].text || '';
    }
    return event.message || '';
  }

  getEventTimestamp(event: InterAgentEvent): string {
    return new Date(event.timestamp * 1000).toLocaleTimeString();
  }

  getFromAgentOptions(): string[] {
    const options = ['user'];
    return options.concat(this.agents.map(agent => agent.name));
  }

  getToAgentOptions(): string[] {
    return this.agents.map(agent => agent.name);
  }

  isInterAgentEvent(event: InterAgentEvent): boolean {
    return event.type === 'inter_agent' || 
           (event.interAgent !== undefined);
  }

  getEventTypeClass(event: InterAgentEvent): string {
    if (this.isInterAgentEvent(event)) {
      return 'inter-agent-event';
    }
    return 'system-event';
  }

  trackByEventId(index: number, event: InterAgentEvent): string {
    return event.id || event.invocationId || `${index}-${event.timestamp}`;
  }

  onEnterKeydown(event: Event): void {
    const keyboardEvent = event as KeyboardEvent;
    if (keyboardEvent.ctrlKey) {
      keyboardEvent.preventDefault();
      this.sendMessage();
    }
  }
}