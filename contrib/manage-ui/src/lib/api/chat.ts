import { api, createSSEConnection } from './client';
import type {
  ChatSession,
  ChatMessage,
  ChatRequest,
  ApiResponse,
} from '../types';

// Chat API functions
export const chatApi = {
  // Session management
  async createChatSession(entityType: 'project' | 'task', entityId: string): Promise<ChatSession> {
    const response = await api.post<ApiResponse<ChatSession>>('/chat/sessions', {
      entity_type: entityType,
      entity_id: entityId,
    });
    return response.data;
  },

  async getChatSession(sessionId: string): Promise<ChatSession> {
    const response = await api.get<ApiResponse<ChatSession>>(`/chat/sessions/${sessionId}`);
    return response.data;
  },

  async getChatSessions(entityType?: 'project' | 'task', entityId?: string): Promise<ChatSession[]> {
    const params = new URLSearchParams();
    if (entityType) params.append('entity_type', entityType);
    if (entityId) params.append('entity_id', entityId);
    
    const queryString = params.toString();
    const endpoint = queryString ? `/chat/sessions?${queryString}` : '/chat/sessions';
    
    const response = await api.get<ApiResponse<ChatSession[]>>(endpoint);
    return response.data;
  },

  async deleteChatSession(sessionId: string): Promise<void> {
    await api.delete(`/chat/sessions/${sessionId}`);
  },

  // Message management
  async getChatMessages(sessionId: string): Promise<ChatMessage[]> {
    const response = await api.get<ApiResponse<ChatMessage[]>>(`/chat/sessions/${sessionId}/messages`);
    return response.data;
  },

  async sendChatMessage(sessionId: string, message: string): Promise<ChatMessage> {
    const response = await api.post<ApiResponse<ChatMessage>>(`/chat/sessions/${sessionId}/messages`, {
      content: message,
    });
    return response.data;
  },

  // Quick chat (creates session and sends message in one call)
  async quickChat(data: ChatRequest): Promise<{ session: ChatSession; message: ChatMessage }> {
    const response = await api.post<ApiResponse<{ session: ChatSession; message: ChatMessage }>>('/chat/quick', data);
    return response.data;
  },

  // SSE connection for real-time chat
  createChatSSEConnection(
    sessionId: string,
    onMessage: (message: ChatMessage) => void,
    onError?: (error: Event) => void,
    onOpen?: (event: Event) => void
  ): EventSource {
    return createSSEConnection(
      `/chat/sessions/${sessionId}/stream`,
      (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type === 'message') {
            onMessage(data.message);
          }
        } catch (error) {
          console.error('Failed to parse SSE message:', error);
        }
      },
      onError,
      onOpen
    );
  },
};