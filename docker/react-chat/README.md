# Multi-Agent Chat Interface

A Google ADK-compatible multi-agent chat interface built with Next.js, shadcn-ui, and Tailwind CSS. This application enables real-time communication with multiple AI agents simultaneously, featuring inter-agent communication visualization and workspace management.

## Features

- **Real-time Multi-Agent Communication**: Chat with multiple agents simultaneously with live streaming support
- **Inter-Agent Communication Visualization**: Monitor and visualize communication between agents
- **Workspace Management**: Switch between Chat and Settings workspaces
- **Dark/Light Theme Support**: Toggle between themes with persistent preference
- **Responsive Design**: Optimized for desktop and mobile viewing
- **ADK-Compatible**: Fully compatible with Google ADK backend protocols
- **SSE Integration**: Server-Sent Events for real-time updates

## Technology Stack

- **Framework**: Next.js 15+ with App Router
- **UI Components**: shadcn-ui
- **Styling**: Tailwind CSS v4
- **State Management**: Zustand
- **Server State**: TanStack Query (React Query)
- **Real-time**: Server-Sent Events (EventSource API)
- **Language**: TypeScript
- **Icons**: Lucide React
- **Animations**: Framer Motion

## Architecture

```
src/
├── app/                    # Next.js App Router pages
├── components/
│   ├── ui/                 # shadcn-ui components
│   ├── chat/               # Chat interface components
│   ├── agents/             # Agent management components
│   ├── workspace/          # Workspace switching
│   ├── navigation/         # Navigation components
│   ├── inter-agent/        # Inter-agent communication
│   ├── layout/             # Layout components
│   └── providers/          # Context providers
├── lib/
│   ├── api/                # API client and services
│   ├── store/              # Zustand stores
│   └── types/              # TypeScript definitions
└── hooks/                  # Custom React hooks
```

## Key Components

### Chat Interface
- **Message List**: Displays conversation history with streaming support
- **Message Input**: Send messages with keyboard shortcuts (Ctrl/Cmd+Enter)
- **Chat Header**: Shows active agent information and status
- **Typing Indicator**: Visual feedback for agent activity

### Agent Management
- **Agent List**: Overview of all available agents
- **Agent Cards**: Individual agent status and capabilities
- **Agent Selector**: Dropdown for selecting target agents
- **Status Indicators**: Real-time agent availability

### Navigation
- **Top Navigation**: Breadcrumb navigation with theme toggle
- **Bottom Navigation**: Agent tabs for quick switching
- **Workspace Sidebar**: Switch between Chat and Settings
- **Theme Toggle**: Dark/light mode switching

### Inter-Agent Features
- **Event Display**: Real-time inter-agent communication events
- **Event Filtering**: Different event types (communication, heartbeat, agent_list)
- **Visual Indicators**: Color-coded event types and status

## Getting Started

### Prerequisites
- Node.js 18+ 
- npm or yarn

### Installation

1. **Navigate to the project directory:**
   ```bash
   cd docker/react-chat
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Set up environment variables:**
   Create a `.env.local` file:
   ```env
   NEXT_PUBLIC_API_URL=http://localhost:6999
   ```

4. **Start the development server:**
   ```bash
   npm run dev
   ```

5. **Open your browser:**
   Navigate to `http://localhost:3000`

### Building for Production

```bash
npm run build
npm start
```

## Backend Integration

This frontend is designed to work with the Go-based agent backend located in the main project. The backend should provide:

- **Agent Discovery**: `GET /list-apps`
- **Session Management**: Session creation and retrieval
- **Message Streaming**: `POST /run_sse` for real-time chat
- **Inter-Agent Events**: SSE stream for agent communication

### API Endpoints Expected

```typescript
// Agent discovery
GET /list-apps
Response: string[] // Array of agent names

// Session management
GET /apps/{appName}/users/{userId}/sessions
POST /apps/{appName}/users/{userId}/sessions
GET /apps/{appName}/users/{userId}/sessions/{sessionId}

// Message sending (streaming)
POST /run_sse
Content-Type: application/json
Accept: text/event-stream

// Multi-agent communication
POST /multi-chat/send
GET /multi-chat/start_sse
```

## Configuration

### Environment Variables

- `NEXT_PUBLIC_API_URL`: Backend server URL (default: http://localhost:6999)

### Theme Configuration

The application supports automatic theme detection and manual toggle. Theme preference is persisted in localStorage.

### Workspace Configuration

Two main workspaces are available:
- **Chat**: Main chat interface with agent communication
- **Settings**: Configuration and preferences

## Development

### Project Structure

- Components are organized by feature (chat, agents, navigation, etc.)
- Each component file is kept under 500 lines for maintainability
- Reusable components follow atomic design principles
- All UI components use shadcn-ui - no custom component creation

### State Management

- **Chat Store**: Manages agents, sessions, and messages
- **Workspace Store**: Handles workspace and sidebar state  
- **Theme Store**: Persists theme preferences
- **Agents Store**: Manages agent-specific state

### API Integration

- Centralized API client with TypeScript interfaces
- SSE service for real-time communication
- Error handling and connection management
- Automatic reconnection for SSE streams

## Contributing

1. Follow the existing code structure and naming conventions
2. Keep component files under 500 lines
3. Use shadcn-ui components exclusively for UI elements
4. Add TypeScript types for all new features
5. Test components with the backend integration

## Troubleshooting

### Common Issues

1. **Connection Issues**: Verify backend server is running on correct port
2. **SSE Not Working**: Check CORS settings on backend
3. **Theme Not Persisting**: Clear localStorage and reload
4. **Build Errors**: Ensure all shadcn-ui components are properly installed

### Backend Compatibility

This frontend expects the backend to implement the ADK-compatible API as described in the main project documentation. Refer to `pkg/multi/plugins/web/server.go` for the expected server implementation.

## License

This project follows the same license as the main agents project.