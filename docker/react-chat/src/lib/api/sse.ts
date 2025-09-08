import { LLMEvent, AgentRunRequest, EventType } from '@/lib/types'
import { apiClient } from './client'

export type SSEEventHandler = {
  /** Handles all types of agent events: regular responses, tool calls, system messages */
  onMessage?: (event: LLMEvent) => void
  /** Specifically handles inter-agent communication events */
  onInterAgentEvent?: (event: LLMEvent) => void
  /** Handles agent status changes */
  onAgentStatusChange?: (agentId: string, status: string) => void
  /** Handles connection status changes */
  onConnectionStatusChange?: (connected: boolean) => void
  /** Handles connection and parsing errors */
  onError?: (error: Event) => void
}

class SSEService {
  private mainChatAbortController: AbortController | null = null;
  private agentRunAbortController: AbortController | null = null;
  private mainChatHandlers: SSEEventHandler = {};
  private agentRunHandlers: SSEEventHandler = {};
  private isMainChatConnected = false;
  private isAgentRunConnected = false;

  connectMainChat(agents: string[], sessionId: string, userId: string, handlers: SSEEventHandler) {
    this.disconnect('mainChat'); // Ensure only one main chat connection
    this.mainChatHandlers = handlers;
    this.mainChatAbortController = new AbortController();
    const signal = this.mainChatAbortController.signal;

    const url = apiClient.getSSEUrl(agents, sessionId, userId);

    fetch(url, {
      method: 'GET',
      signal,
      headers: {
        'Accept': 'text/event-stream',
        'Cache-Control': 'no-cache',
      },
    }).then(async (response) => {
      if (!response.ok) {
        throw new Error(`SSE request failed: ${response.status} ${response.statusText}`);
      }
      
      if (!response.body) {
        throw new Error('No response body');
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();

      this.isMainChatConnected = true;
      this.mainChatHandlers.onConnectionStatusChange?.(true);

      try {
        while (true) {
          const { done, value } = await reader.read();
          if (done || signal.aborted) {
            break;
          }

          const chunk = decoder.decode(value, { stream: true });
          const lines = chunk.split('\n');

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              try {
                const data = JSON.parse(line.slice(6));
                this.handleEvent(data, 'mainChat');
              } catch (error) {
                console.error('Error parsing main chat SSE data:', error, 'Raw data:', line);
              }
            }
          }
        }
      } catch (error) {
        if ((error as Error).name !== 'AbortError') {
            console.error('Main chat SSE read error:', error);
            if (this.mainChatHandlers.onError) {
                this.mainChatHandlers.onError(error as Event);
            }
        }
      } finally {
        if (!signal.aborted) {
          this.isMainChatConnected = false
          this.mainChatHandlers.onConnectionStatusChange?.(false)
        }
      }
    }).catch((error) => {
        if (error.name === 'AbortError') {
            return;
        }
      console.error('Main chat SSE fetch setup error:', error);
      this.isMainChatConnected = false;
      this.mainChatHandlers.onConnectionStatusChange?.(false);
      if (this.mainChatHandlers.onError) {
        this.mainChatHandlers.onError(error);
      }
      
      // Optional: Reconnect logic
      setTimeout(() => {
        if (!this.isMainChatConnected && !signal.aborted) {
          this.connectMainChat(agents, sessionId, userId, this.mainChatHandlers);
        }
      }, 5000);
    });
  }

  /**
   * Connects to agent run SSE endpoint for user-to-agent communication.
   * Used for sending messages to agents and receiving their responses.
   */
  connectAgentRun(request: AgentRunRequest, handlers: SSEEventHandler) {
    this.disconnect('agentRun'); // Ensure only one agent run connection
    this.agentRunHandlers = handlers;
    this.agentRunAbortController = new AbortController();
    const signal = this.agentRunAbortController.signal;

    const url = apiClient.getRunSSEUrl();
    
    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'text/event-stream',
        'Cache-Control': 'no-cache',
        'Connection': 'keep-alive',
      },
      body: JSON.stringify(request),
      signal,
    }).then(async (response) => {
      if (!response.ok) {
        throw new Error(`Agent run SSE request failed: ${response.status} ${response.statusText}`);
      }
      
      if (!response.body) {
        throw new Error('No response body');
      }

      // Check if response is actually SSE
      const contentType = response.headers.get('content-type');
      if (!contentType?.includes('text/event-stream')) {
        console.warn('Response is not SSE, content-type:', contentType);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();

      this.isAgentRunConnected = true;
      this.agentRunHandlers.onConnectionStatusChange?.(true);

      try {
        while (true) {
          const { done, value } = await reader.read();
          if (done || signal.aborted) {
            break;
          }

          const chunk = decoder.decode(value, { stream: true });
          const lines = chunk.split('\n');

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              try {
                const data = JSON.parse(line.slice(6));
                this.handleEvent(data, 'agentRun');
              } catch (error) {
                console.error('Error parsing agent run SSE data:', error, 'Raw data:', line);
              }
            } else if (line.trim() === '' || line.startsWith(':')) {
              // Skip empty lines and comments
              continue;
            } else if (line.trim() !== '' && !line.startsWith('event:') && !line.startsWith('id:')) {
              console.warn('Unexpected SSE line format:', line);
            }
          }
        }
      } catch (error) {
        if ((error as Error).name !== 'AbortError') {
            console.error('Agent run SSE read error:', error);
            if (this.agentRunHandlers.onError) {
                this.agentRunHandlers.onError(error as Event);
            }
        }
      } finally {
        if (!signal.aborted) {
          this.isAgentRunConnected = false
          this.agentRunHandlers.onConnectionStatusChange?.(false)
        }
      }
    }).catch((error) => {
        if (error.name === 'AbortError') {
            return;
        }
      console.error('Agent run SSE fetch setup error:', error);
      this.isAgentRunConnected = false;
      this.agentRunHandlers.onConnectionStatusChange?.(false);
      if (this.agentRunHandlers.onError) {
        this.agentRunHandlers.onError(error);
      }
      
      // Optional: Reconnect logic
      setTimeout(() => {
        if (!this.isAgentRunConnected && !signal.aborted) {
          this.connectAgentRun(request, this.agentRunHandlers);
        }
      }, 5000);
    });
  }


  private handleEvent(event: LLMEvent, connectionType: 'mainChat' | 'agentRun') {
    const handlers = connectionType === 'mainChat' ? this.mainChatHandlers : this.agentRunHandlers;

    // Handle stream termination event
    if (event.done === true && !event.type && !event.parts) {
      if (connectionType === 'agentRun') {
        this.isAgentRunConnected = false;
        handlers.onConnectionStatusChange?.(false);
      }
      return;
    }

    switch (event.type) {
      case 'agent_list':
        handlers.onMessage?.(event);
        break;
      
      case 'inter_agent':
      case 'communication':
        handlers.onInterAgentEvent?.(event);
        break;
      
      case 'heartbeat':
        // Keep connection alive - don't change status on heartbeat
        // The connection is already established, heartbeat just confirms it's alive
        break;
      
      default:
        // Handle tool calls and responses as regular messages
        if (event.type === EventType.TOOL_CALL || event.type === EventType.TOOL_RESPONSE || event.parts) {
          handlers.onMessage?.(event);
        } else {
          console.log('Unhandled SSE event type:', event.type, event);
        }
        break;
    }
  }

  disconnect(type?: 'mainChat' | 'agentRun') {
    if (type === 'mainChat' || !type) {
      if (this.mainChatAbortController) {
          this.mainChatAbortController.abort();
          this.mainChatAbortController = null;
      }
      this.isMainChatConnected = false;
      this.mainChatHandlers.onConnectionStatusChange?.(false);
    }
    if (type === 'agentRun' || !type) {
      if (this.agentRunAbortController) {
          this.agentRunAbortController.abort();
          this.agentRunAbortController = null;
      }
      this.isAgentRunConnected = false;
      this.agentRunHandlers.onConnectionStatusChange?.(false);
    }
  }

  getConnectionStatus(): boolean {
    return this.isMainChatConnected || this.isAgentRunConnected;
  }
}

export const sseService = new SSEService()