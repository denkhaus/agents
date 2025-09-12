# Multi-Agent System Admin Frontend

A modern, responsive admin interface for managing projects and tasks within the Multi-Agent System. Built with NextJS 15, TypeScript, Tailwind CSS v4, and shadcn/ui components.

## 🚀 Features

### MVP Features (Completed)
- ✅ **Project Management**: Create, view, and manage projects
- ✅ **Hierarchical Task Management**: Kanban-style view with task hierarchy
- ✅ **Real-time Chat Integration**: LLM-powered assistance for projects and tasks
- ✅ **Responsive Design**: Desktop-first with mobile support
- ✅ **Modern UI**: Clean, intuitive interface with dark/light theme support
- ✅ **State Management**: Zustand for efficient state handling
- ✅ **Type Safety**: Full TypeScript implementation

### Planned Features
- 🔄 **Real-time Updates**: SSE integration for live collaboration
- 🔄 **Advanced Filtering**: Complex search and filtering capabilities
- 🔄 **Drag & Drop**: Task state management via drag and drop
- 🔄 **Bulk Operations**: Multi-select and bulk task operations
- 🔄 **Agent Management**: Interface for managing AI agents
- 🔄 **Analytics Dashboard**: Project progress and performance metrics

## 🛠 Tech Stack

- **Framework**: NextJS 15 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **Components**: shadcn/ui (Radix UI primitives)
- **State Management**: Zustand
- **Data Fetching**: TanStack Query (React Query)
- **Icons**: Lucide React
- **Real-time**: Server-Sent Events (SSE)

## 📁 Project Structure

```
src/
├── app/                    # NextJS App Router
│   ├── globals.css        # Global styles and CSS variables
│   ├── layout.tsx         # Root layout
│   └── page.tsx           # Main page
├── components/
│   ├── ui/                # shadcn/ui base components
│   ├── layout/            # Layout components (sidebar, navigation)
│   ├── projects/          # Project-specific components
│   ├── tasks/             # Task-specific components
│   ├── kanban/            # Kanban view components
│   └── chat/              # LLM chat components
├── lib/
│   ├── api/               # API client functions
│   ├── stores/            # Zustand stores
│   ├── types/             # TypeScript type definitions
│   └── utils/             # Utility functions
└── styles/                # Additional styles
```

## 🚦 Getting Started

### Prerequisites

- Node.js 18+ 
- npm, yarn, or pnpm

### Installation

1. **Navigate to the project directory**:
   ```bash
   cd contrib/manage-ui
   ```

2. **Install dependencies**:
   ```bash
   npm install
   # or
   yarn install
   # or
   pnpm install
   ```

3. **Set up environment variables**:
   ```bash
   cp .env.example .env.local
   ```
   
   Configure the following variables:
   ```env
   NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api
   ```

4. **Run the development server**:
   ```bash
   npm run dev
   # or
   yarn dev
   # or
   pnpm dev
   ```

5. **Open your browser**:
   Navigate to [http://localhost:3000](http://localhost:3000)

## 🏗 Architecture

### State Management

The application uses Zustand for state management with the following stores:

- **`useProjectStore`**: Manages projects, tasks, and related data
- **`useUIStore`**: Handles UI state, selections, modals, and view configurations

### API Integration

The API client is built with:
- RESTful endpoints for CRUD operations
- Server-Sent Events for real-time updates
- Type-safe request/response handling
- Error handling and retry logic

### Component Architecture

- **Layout Components**: Sidebar navigation and top bar
- **Workspace Components**: Modular workspace areas
- **Feature Components**: Project and task management
- **UI Components**: Reusable shadcn/ui components

## 🎨 UI/UX Design

### Design Principles

- **Hierarchy First**: Clear visual hierarchy for project/task relationships
- **Progressive Disclosure**: Show relevant information at the right time
- **Consistent Interactions**: Standardized patterns across the interface
- **Accessibility**: WCAG 2.1 AA compliance

### Kanban View

The hierarchical Kanban view displays:
- **Projects**: Top-level project cards
- **Tasks by State**: Columns for pending, in-progress, completed, blocked, cancelled
- **Task Cards**: Compact cards with essential information
- **Chat Integration**: Quick access to LLM assistance

### Color System

- **Primary**: Blue (#3b82f6) for actions and highlights
- **States**: Color-coded task states (green=completed, blue=in-progress, etc.)
- **Neutral**: Gray scale for backgrounds and secondary content
- **Semantic**: Red for destructive actions, yellow for warnings

## 🔌 Backend Integration

### Expected API Endpoints

```
# Projects
GET    /api/projects
POST   /api/projects
GET    /api/projects/:id
PUT    /api/projects/:id
DELETE /api/projects/:id

# Tasks
GET    /api/projects/:id/tasks
POST   /api/tasks
GET    /api/tasks/:id
PUT    /api/tasks/:id
DELETE /api/tasks/:id

# Chat
POST   /api/chat/sessions
GET    /api/chat/sessions/:id/messages
POST   /api/chat/sessions/:id/messages

# SSE
GET    /api/sse/projects
GET    /api/sse/tasks
GET    /api/sse/chat/:sessionId
```

### Data Models

The frontend expects data models matching the Go structs:

- **Project**: ID, title, description, progress metrics, timestamps
- **Task**: ID, project_id, parent_id, title, description, state, complexity, depth, etc.
- **TaskState**: pending, in-progress, completed, blocked, cancelled

## 🧪 Development

### Mock Data

The MVP includes mock data for development and testing:
- Sample projects with realistic metadata
- Hierarchical task structures
- Various task states and complexities

### Adding New Features

1. **Define Types**: Add TypeScript interfaces in `src/lib/types/`
2. **Create API Functions**: Add client functions in `src/lib/api/`
3. **Update Stores**: Extend Zustand stores as needed
4. **Build Components**: Create reusable components
5. **Add Routes**: Integrate with the workspace system

### Code Style

- **TypeScript**: Strict mode enabled
- **ESLint**: NextJS recommended configuration
- **Prettier**: Consistent code formatting
- **Naming**: PascalCase for components, camelCase for functions/variables

## 🚀 Deployment

### Build for Production

```bash
npm run build
npm run start
```

### Environment Variables

Production environment variables:
```env
NEXT_PUBLIC_API_BASE_URL=https://your-api-domain.com/api
NEXT_PUBLIC_WS_URL=wss://your-api-domain.com/ws
```

### Docker Support

```dockerfile
# Example Dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```

## 🤝 Contributing

### Development Workflow

1. Create feature branch from `main`
2. Implement changes with tests
3. Update documentation
4. Submit pull request

### Code Quality

- All components must be TypeScript
- Follow established patterns
- Include proper error handling
- Write meaningful commit messages

## 📚 Documentation

### Component Documentation

Each component includes:
- TypeScript interfaces for props
- JSDoc comments for complex logic
- Usage examples in Storybook (planned)

### API Documentation

- OpenAPI/Swagger specs for backend integration
- Type definitions generated from Go structs
- Error response documentation

## 🐛 Troubleshooting

### Common Issues

1. **Build Errors**: Check TypeScript errors and missing dependencies
2. **API Connection**: Verify backend is running and CORS is configured
3. **Styling Issues**: Ensure Tailwind CSS is properly configured

### Debug Mode

Enable debug logging:
```env
NEXT_PUBLIC_DEBUG=true
```

## 📄 License

This project is part of the Multi-Agent System and follows the same licensing terms.

## 🔗 Related Projects

- [Multi-Agent System Backend](../../) - Go-based backend API
- [Agent Framework](../../pkg/) - Core agent functionality
- [Project Tools](../../pkg/tools/project/) - Project management tools

---

**Version**: 1.0.0 MVP  
**Last Updated**: December 2024  
**Maintainer**: Multi-Agent System Team