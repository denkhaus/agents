import { 
  AgentRunRequest, 
  MultiChatRequest, 
  MultiChatResponse, 
  ADKSession
} from '@/lib/types'

class ApiClient {
  private baseUrl: string

  constructor() {
    // Default to localhost:6999, can be configured via environment
    this.baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:6999'
  }

  private async request<T>(
    endpoint: string, 
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`
    
    console.log('API Request:', url, options)
    
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    })

    console.log('API Response:', response.status, response.statusText)

    if (!response.ok) {
      const errorText = await response.text()
      console.error('API Error:', response.status, errorText)
      throw new Error(`API request failed: ${response.status} ${response.statusText} - ${errorText}`)
    }

    const data = await response.json()
    console.log('API Data:', data)
    return data
  }

  // Agent endpoints
  async getAgents(): Promise<string[]> {
    return this.request<string[]>('/list-apps')
  }

  // Session endpoints
  async getSessions(appName: string, userId: string): Promise<ADKSession[]> {
    return this.request<ADKSession[]>(`/apps/${appName}/users/${userId}/sessions`)
  }

  async createSession(appName: string, userId: string): Promise<ADKSession> {
    return this.request<ADKSession>(
      `/apps/${appName}/users/${userId}/sessions`,
      { method: 'POST' }
    )
  }

  async getSession(
    appName: string, 
    userId: string, 
    sessionId: string
  ): Promise<ADKSession> {
    return this.request<ADKSession>(
      `/apps/${appName}/users/${userId}/sessions/${sessionId}`
    )
  }

  // Agent run endpoints
  async runAgent(request: AgentRunRequest): Promise<unknown> {
    return this.request('/run', {
      method: 'POST',
      body: JSON.stringify(request),
    })
  }

  // Multi-agent chat endpoints
  async sendMessage(request: MultiChatRequest): Promise<MultiChatResponse> {
    return this.request<MultiChatResponse>('/multi-chat/send', {
      method: 'POST',
      body: JSON.stringify(request),
    })
  }

  // SSE endpoint URL
  getSSEUrl(agents: string[], sessionId: string, userId: string): string {
    const agentList = agents.join(',')
    const params = new URLSearchParams({
      agents: agentList,
      sessionId,
      userId,
    })
    return `${this.baseUrl}/multi-chat/start_sse?${params.toString()}`
  }

  // Agent run SSE endpoint
  getRunSSEUrl(): string {
    return `${this.baseUrl}/run_sse`
  }
}

export const apiClient = new ApiClient()