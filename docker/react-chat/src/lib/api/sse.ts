import { InterAgentEvent, AgentRunRequest } from '@/lib/types'
import { apiClient } from './client'

export type SSEEventHandler = {
  onMessage?: (event: InterAgentEvent) => void
  onInterAgentEvent?: (event: InterAgentEvent) => void
  onAgentStatusChange?: (agentId: string, status: string) => void
  onConnectionStatusChange?: (connected: boolean) => void
  onError?: (error: Event) => void
}

class SSEService {
  private eventSource: EventSource | null = null
  private handlers: SSEEventHandler = {}
  private isConnected = false
  private abortController: AbortController | null = null;

  connect(agents: string[], sessionId: string, userId: string, handlers: SSEEventHandler) {
    this.disconnect(); // Close any existing connection
    this.handlers = handlers;
    this.abortController = new AbortController();
    const signal = this.abortController.signal;

    const url = apiClient.getSSEUrl(agents, sessionId, userId);

    fetch(url, {
      signal,
      headers: {
        'Accept': 'text/event-stream',
        'Cache-Control': 'no-cache',
      },
    }).then(async (response) => {
      if (!response.body) {
        throw new Error('No response body');
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();

      this.isConnected = true;
      this.handlers.onConnectionStatusChange?.(true);
      console.log('SSE connection established');

      try {
        while (true) {
          const { done, value } = await reader.read();
          if (done || signal.aborted) {
            console.log('SSE stream finished or aborted.');
            break;
          }

          const chunk = decoder.decode(value, { stream: true });
          const lines = chunk.split('\n');

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              try {
                const data = JSON.parse(line.slice(6));
                this.handleEvent(data);
              } catch (error) {
                console.error('Error parsing SSE data:', error, 'Raw data:', line);
              }
            }
          }
        }
      } catch (error) {
        if (error.name !== 'AbortError') {
            console.error('SSE read error:', error);
            if (this.handlers.onError) {
                this.handlers.onError(error as Event);
            }
        }
      } finally {
        this.isConnected = false;
        this.handlers.onConnectionStatusChange?.(false);
        console.log('SSE connection closed.');
      }
    }).catch((error) => {
        if (error.name === 'AbortError') {
            console.log('SSE fetch aborted by client.');
            return;
        }
      console.error('SSE fetch setup error:', error);
      this.isConnected = false;
      this.handlers.onConnectionStatusChange?.(false);
      if (this.handlers.onError) {
        this.handlers.onError(error);
      }
      
      // Optional: Reconnect logic
      setTimeout(() => {
        if (!this.isConnected && !signal.aborted) {
          console.log('Attempting to reconnect SSE...');
          this.connect(agents, sessionId, userId, this.handlers);
        }
      }, 5000);
    });
  }

  connectForAgentRun(request: AgentRunRequest, handlers: SSEEventHandler) {
    this.disconnect()
    this.handlers = handlers

    // For agent run SSE, we need to POST to the run_sse endpoint
    const url = apiClient.getRunSSEUrl()
    
    // Create a fetch request that will establish SSE connection
    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'text/event-stream',
        'Cache-Control': 'no-cache',
      },
      body: JSON.stringify(request),
    }).then(async (response) => {
      if (!response.body) {
        throw new Error('No response body')
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder()

      this.isConnected = true

      try {
        while (true) {
          const { done, value } = await reader.read()
          if (done) break

          const chunk = decoder.decode(value)
          const lines = chunk.split('\n')

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              try {
                const data = JSON.parse(line.slice(6))
                this.handleEvent(data)
              } catch (error) {
                console.error('Error parsing SSE data:', error)
              }
            }
          }
        }
      } catch (error) {
        console.error('SSE read error:', error)
        this.handlers.onError?.(error as Event)
      } finally {
        this.isConnected = false
        this.handlers.onConnectionStatusChange?.(false)
      }
    }).catch((error) => {
      console.error('SSE connection error:', error)
      this.isConnected = false
      this.handlers.onConnectionStatusChange?.(false)
      this.handlers.onError?.(error)
    })
  }

  private handleEvent(event: InterAgentEvent) {
    switch (event.type) {
      case 'agent_list':
        // Handle agent list updates
        this.handlers.onMessage?.(event)
        break
      
      case 'inter_agent':
      case 'communication':
        this.handlers.onInterAgentEvent?.(event)
        break
      
      case 'heartbeat':
        // Handle heartbeat - update connection status and log
        console.log('Received heartbeat:', event.timestamp)
        this.handlers.onConnectionStatusChange?.(true)
        break
      
      default:
        // Handle other event types
        this.handlers.onMessage?.(event)
        break
    }
  }

  disconnect() {
    if (this.abortController) {
        this.abortController.abort();
        this.abortController = null;
    }
    if (this.eventSource) {
      this.eventSource.close()
      this.eventSource = null
    }
    // Note: isConnected and onConnectionStatusChange are handled in the fetch promise finally block
  }

  getConnectionStatus(): boolean {
    return this.isConnected
  }
}

export const sseService = new SSEService()