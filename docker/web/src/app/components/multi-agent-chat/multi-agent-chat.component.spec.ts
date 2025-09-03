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

import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { FormsModule } from '@angular/forms';
import { MultiAgentChatComponent } from './multi-agent-chat.component';
import { MultiAgentChatService } from '../../core/services/multi-agent-chat.service';
import { AgentService } from '../../core/services/agent.service';

describe('MultiAgentChatComponent', () => {
  let component: MultiAgentChatComponent;
  let fixture: ComponentFixture<MultiAgentChatComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        MultiAgentChatComponent,
        HttpClientTestingModule,
        FormsModule
      ],
      providers: [
        MultiAgentChatService,
        AgentService
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MultiAgentChatComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize with default values', () => {
    expect(component.selectedFromAgent).toBe('user');
    expect(component.selectedToAgent).toBe('');
    expect(component.messageText).toBe('');
    expect(component.isConnected).toBe(false);
    expect(component.isLoading).toBe(false);
  });

  it('should handle agent selection', () => {
    component.availableAgents = ['agent1', 'agent2'];
    component.ngOnInit();
    
    // Simulate agents being loaded
    const mockAgents = [
      { id: '1', name: 'agent1', role: 'assistant' },
      { id: '2', name: 'agent2', role: 'assistant' }
    ];
    
    // This would normally come from the service
    component.agents = mockAgents;
    
    expect(component.getToAgentOptions()).toEqual(['agent1', 'agent2']);
  });
});