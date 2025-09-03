# Multi-Agent Chat Integration

This document describes the Multi-Agent Chat functionality added to the ADK Web Interface.

## Overview

The Multi-Agent Chat feature enables real-time communication between different agents in the ADK system. It provides a visual interface for monitoring and facilitating inter-agent communication.

## Components

### Backend Components

#### 1. Multi-Agent Chat Service (`pkg/multi/plugins/web/server/multi_chat.go`)

**Key Features:**
- SSE (Server-Sent Events) streaming for real-time communication
- RESTful API endpoints for sending messages between agents
- Integration with the existing ADK server architecture

**Endpoints:**
- `POST /multi-chat/send` - Send messages between agents
- `GET /multi-chat/start_sse` - Establish SSE connection for real-time events

**Data Models:**
- `MultiChatRequest` - Request structure for sending messages
- `InterAgentEvent` - Event structure for inter-agent communication
- `AgentInfo` - Agent metadata structure

### Frontend Components

#### 1. Multi-Agent Chat Service (`src/app/core/services/multi-agent-chat.service.ts`)

**Responsibilities:**
- Manage SSE connections for real-time events
- Send HTTP requests for message transmission
- Handle event parsing and distribution

#### 2. Multi-Agent Chat Component (`src/app/components/multi-agent-chat/`)

**Features:**
- Real-time event display
- Agent selection interface
- Message composition and sending
- Connection status monitoring
- Event filtering and categorization

**UI Elements:**
- Agent list with roles and IDs
- Message input with from/to agent selection
- Real-time event stream display
- Connection status indicator

#### 3. Models (`src/app/core/models/MultiAgentChat.ts`)

**Interfaces:**
- `AgentInfo` - Agent metadata
- `MultiChatRequest` - Message request structure
- `InterAgentEvent` - Event data structure
- `MultiChatResponse` - Response structure

## Integration

### Chat Component Integration

The Multi-Agent Chat is integrated as a new tab in the main chat interface:

```html
<mat-tab>
  <ng-template mat-tab-label>
    <span class="tab-label">Multi-Agent</span>
  </ng-template>
  <app-multi-agent-chat
    [sessionId]="sessionId"
    [userId]="userId"
    [availableAgents]="availableAgents"
  ></app-multi-agent-chat>
</mat-tab>
```

### Server Integration

The multi-agent chat functionality is integrated into the existing server through:

```go
// In server.go New() function
s.registerMultiChatRoutes()

// Optional chat processor integration
func WithChatProcessor(processor multi.ChatProcessor) Option {
    return func(s *Server) { 
        s.chatProcessor = processor
        s.setupInterAgentInterceptor()
    }
}
```

## Usage

### 1. Starting a Multi-Agent Chat Session

1. Navigate to the Multi-Agent tab in the ADK Web interface
2. The system automatically connects to the SSE stream
3. Available agents are displayed in the agents panel

### 2. Sending Messages

1. Select the source agent (From field)
2. Select the target agent (To field)
3. Type your message
4. Click Send or press Ctrl+Enter

### 3. Monitoring Communication

- Real-time events appear in the communication events panel
- Inter-agent events are highlighted differently from system events
- Connection status is displayed in the header

## Event Types

### Inter-Agent Events
- **Type:** `inter_agent` or `communication`
- **Display:** Shows from/to agents with arrow indicator
- **Styling:** Blue background with distinct formatting

### System Events
- **Type:** `agent_list`, `heartbeat`, etc.
- **Display:** System-level information
- **Styling:** Gray background

## Configuration

### Backend Configuration

```go
// Enable multi-agent chat in server setup
server := server.New(agents, 
    server.WithChatProcessor(chatProcessor),
    // other options...
)
```

### Frontend Configuration

The component automatically configures itself based on:
- Available agents from the agent service
- Current session and user context
- Server endpoints from URL utilities

## API Reference

### REST Endpoints

#### Send Message
```
POST /multi-chat/send
Content-Type: application/json

{
  "fromAgent": "user",
  "toAgent": "agent1",
  "message": "Hello, agent1!",
  "sessionId": "session-123",
  "userId": "user-456"
}
```

#### SSE Stream
```
GET /multi-chat/start_sse?agents=agent1,agent2&sessionId=session-123&userId=user-456
Accept: text/event-stream
```

### Event Format

```json
{
  "type": "inter_agent",
  "id": "event-id",
  "timestamp": 1640995200,
  "content": {
    "role": "assistant",
    "parts": [{"text": "Message content"}]
  },
  "interAgent": {
    "fromAgent": "agent1",
    "toAgent": "agent2",
    "type": "communication"
  }
}
```

## Styling

The component uses a comprehensive SCSS stylesheet with:
- Responsive grid layouts for agent cards
- Color-coded event types
- Loading states and animations
- Connection status indicators
- Accessible form controls

## Testing

Basic unit tests are provided in `multi-agent-chat.component.spec.ts` covering:
- Component initialization
- Agent selection logic
- Event handling
- Service integration

## Future Enhancements

Potential improvements include:
- Message history persistence
- Agent status indicators (online/offline)
- Message threading and replies
- File attachments in inter-agent messages
- Advanced filtering and search
- Performance optimizations for large agent networks