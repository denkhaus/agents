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

    const response = await fetch(url, {
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error("API Error:", response.status, errorText);
      throw new Error(
        `API request failed: ${response.status} ${response.statusText} - ${errorText}`
      );
    }

    const data = await response.json();
    return data;
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
