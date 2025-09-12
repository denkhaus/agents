// API client configuration and base functions

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api';

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public response?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

interface RequestConfig {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  headers?: Record<string, string>;
  body?: unknown;
  signal?: AbortSignal;
}

async function apiRequest<T>(
  endpoint: string,
  config: RequestConfig = {}
): Promise<T> {
  const {
    method = 'GET',
    headers = {},
    body,
    signal
  } = config;

  const url = `${API_BASE_URL}${endpoint}`;
  
  const requestHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
    ...headers,
  };

  // Add authentication header if available
  const token = localStorage.getItem('auth_token');
  if (token) {
    requestHeaders['Authorization'] = `Bearer ${token}`;
  }

  try {
    const response = await fetch(url, {
      method,
      headers: requestHeaders,
      body: body ? JSON.stringify(body) : undefined,
      signal,
    });

    if (!response.ok) {
      let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
      let errorData;

      try {
        errorData = await response.json();
        errorMessage = errorData.message || errorData.error || errorMessage;
      } catch {
        // If response is not JSON, use status text
      }

      throw new ApiError(errorMessage, response.status, errorData);
    }

    // Handle empty responses
    const contentType = response.headers.get('content-type');
    if (!contentType || !contentType.includes('application/json')) {
      return {} as T;
    }

    return await response.json();
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }

    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        throw new ApiError('Request was cancelled', 0);
      }
      throw new ApiError(`Network error: ${error.message}`, 0);
    }

    throw new ApiError('Unknown error occurred', 0);
  }
}

// HTTP method helpers
export const api = {
  get: <T>(endpoint: string, signal?: AbortSignal) =>
    apiRequest<T>(endpoint, { method: 'GET', signal }),

  post: <T>(endpoint: string, body?: unknown, signal?: AbortSignal) =>
    apiRequest<T>(endpoint, { method: 'POST', body, signal }),

  put: <T>(endpoint: string, body?: unknown, signal?: AbortSignal) =>
    apiRequest<T>(endpoint, { method: 'PUT', body, signal }),

  patch: <T>(endpoint: string, body?: unknown, signal?: AbortSignal) =>
    apiRequest<T>(endpoint, { method: 'PATCH', body, signal }),

  delete: <T>(endpoint: string, signal?: AbortSignal) =>
    apiRequest<T>(endpoint, { method: 'DELETE', signal }),
};

// SSE connection helper
export function createSSEConnection(
  endpoint: string,
  onMessage: (event: MessageEvent) => void,
  onError?: (error: Event) => void,
  onOpen?: (event: Event) => void
): EventSource {
  const url = `${API_BASE_URL}${endpoint}`;
  const eventSource = new EventSource(url);

  eventSource.onmessage = onMessage;
  
  if (onError) {
    eventSource.onerror = onError;
  }
  
  if (onOpen) {
    eventSource.onopen = onOpen;
  }

  return eventSource;
}

// Request timeout helper
export function withTimeout<T>(
  promise: Promise<T>,
  timeoutMs: number = 10000
): Promise<T> {
  return Promise.race([
    promise,
    new Promise<never>((_, reject) =>
      setTimeout(() => reject(new ApiError('Request timeout', 408)), timeoutMs)
    ),
  ]);
}