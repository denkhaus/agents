import {
  AgentRunRequest,
  MultiChatRequest,
  MultiChatResponse,
  ADKSession,
} from "@/lib/types";

class ApiClient {
  private baseUrl: string;

  constructor() {
    // Use different URLs for server-side vs client-side
    const envUrl = process.env.NEXT_PUBLIC_API_URL;
    
    // Check if we're running on the server (Node.js) or client (browser)
    const isServer = typeof window === 'undefined';
    
    if (isServer) {
      // Server-side: use Docker internal network
      this.baseUrl = "http://agents:6999";
    } else {
      // Client-side: use localhost for browser
      this.baseUrl = envUrl || "http://localhost:6999";
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;

    console.log(`[API CLIENT] API Request: ${options.method || 'GET'} ${url}`, {
      headers: options.headers,
      hasBody: !!options.body
    });

    try {
      const response = await fetch(url, {
        headers: {
          "Content-Type": "application/json",
          ...options.headers,
        },
        mode: 'cors',
        credentials: 'omit',
        ...options,
      });

      console.log(`[API CLIENT] API Response: ${response.status} ${response.statusText}`, {
        url,
        status: response.status,
        statusText: response.statusText,
        headers: Object.fromEntries(response.headers.entries())
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error("[API CLIENT] API Error:", response.status, errorText);
        throw new Error(
          `API request failed: ${response.status} ${response.statusText} - ${errorText}`
        );
      }

      // Handle empty responses (like DELETE operations)
      const contentLength = response.headers.get('content-length');
      const contentType = response.headers.get('content-type');
      
      // Check if response has content and is JSON
      if (response.status === 204 || contentLength === '0') {
        // No content responses (like DELETE operations)
        return undefined as T;
      } else if (contentType && contentType.includes('application/json')) {
        const data = await response.json();
        console.log(`[API CLIENT] API Response JSON:`, typeof data === 'object' ? JSON.stringify(data).substring(0, 200) : data);
        return data;
      } else {
        // Try to get text content for non-JSON responses
        const text = await response.text();
        console.log(`[API CLIENT] API Response Text:`, text.substring(0, 200));
        return (text || undefined) as T;
      }
    } catch (error) {
      if (error instanceof TypeError && error.message.includes('fetch')) {
        console.error('[API CLIENT] Network error - API server may be unreachable:', error);
        throw new Error('Network error: Unable to connect to API server. Please check if the server is running.');
      }
      
      // Handle CORS errors specifically
      if (error instanceof TypeError && (
        error.message.includes('CORS') || 
        error.message.includes('cross-origin') ||
        error.message.includes('Failed to fetch')
      )) {
        console.error('[API CLIENT] CORS error detected:', error);
        throw new Error('CORS error: The API server does not allow this request from the browser. This operation will be performed locally only.');
      }
      
      console.error('[API CLIENT] Request error:', error);
      throw error;
    }
  }

  // Agent endpoints
  async getAgents(): Promise<string[]> {
    return this.request<string[]>("/list-apps");
  }

  // Session endpoints
  async getSessions(appName: string, userId: string): Promise<ADKSession[]> {
    return this.request<ADKSession[]>(
      `/apps/${appName}/users/${userId}/sessions`
    );
  }

  async createSession(appName: string, userId: string): Promise<ADKSession> {
    return this.request<ADKSession>(
      `/apps/${appName}/users/${userId}/sessions`,
      { method: "POST" }
    );
  }

  async deleteSession(appName: string, userId: string, sessionId: string): Promise<void> {
    const url = `/apps/${appName}/users/${userId}/sessions/${sessionId}`;
    console.log('Attempting to delete session:', { appName, userId, sessionId, url });
    
    try {
      const result = await this.request<void>(url, { method: "DELETE" });
      console.log('Delete session successful');
      return result;
    } catch (error) {
      console.error('Delete session failed:', error);
      throw error;
    }
  }

  async getSession(
    appName: string,
    userId: string,
    sessionId: string
  ): Promise<ADKSession> {
    return this.request<ADKSession>(
      `/apps/${appName}/users/${userId}/sessions/${sessionId}`
    );
  }

  async addSessionEvent(
    appName: string,
    userId: string,
    sessionId: string,
    eventData: any
  ): Promise<void> {
    return this.request<void>(
      `/apps/${appName}/users/${userId}/sessions/${sessionId}/events`,
      {
        method: "POST",
        body: JSON.stringify(eventData),
      }
    );
  }

  // Agent run endpoints
  async runAgent(request: AgentRunRequest): Promise<unknown> {
    return this.request("/run", {
      method: "POST",
      body: JSON.stringify(request),
    });
  }

  // Multi-agent chat endpoints
  async sendMessage(request: MultiChatRequest): Promise<MultiChatResponse> {
    return this.request<MultiChatResponse>("/multi-chat/send", {
      method: "POST",
      body: JSON.stringify(request),
    });
  }

  // SSE endpoint URL
  getSSEUrl(agents: string[], sessionId: string, userId: string): string {
    const agentList = agents.join(",");
    const params = new URLSearchParams({
      agents: agentList,
      sessionId,
      userId,
    });
    return `${this.baseUrl}/multi-chat/start_sse?${params.toString()}`;
  }

  // Agent run SSE endpoint
  getRunSSEUrl(): string {
    return `${this.baseUrl}/run_sse`;
  }
}

export const apiClient = new ApiClient();
