# Development Guide

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ 
- npm or pnpm
- Convex account (free at https://convex.dev)

### Setup Steps

1. **Install Dependencies**
   ```bash
   cd project-canvas
   npm install
   ```

2. **Initialize Convex**
   ```bash
   npx convex dev
   ```
   - Follow the prompts to create/login to Convex account
   - This will create `.env.local` with your `VITE_CONVEX_URL`

3. **Start Development**
   ```bash
   npm run dev
   ```
   This automatically starts both:
   - Vite dev server (http://localhost:3000)
   - Convex dev server (real-time backend)

## 📁 Project Structure

```
project-canvas/
├── src/
│   ├── components/
│   │   ├── ui/              # shadcn/ui components (auto-generated)
│   │   ├── layout/          # App layout components
│   │   ├── canvas/          # ReactFlow visualization
│   │   └── providers/       # Context providers
│   ├── hooks/               # Custom React hooks
│   ├── stores/              # Zustand state management
│   ├── types/               # TypeScript definitions
│   ├── utils/               # Utility functions
│   └── data/                # Dummy/test data
├── convex/                  # Backend functions & schema
│   ├── schema.ts           # Database schema
│   ├── projects.ts         # Project CRUD operations
│   ├── tasks.ts            # Task CRUD operations
│   ├── agents.ts           # Agent management
│   └── events.ts           # Real-time events
└── public/                  # Static assets
```

## 🔧 Available Scripts

- `npm run dev` - Start both Convex and Vite dev servers
- `npm run dev:vite` - Start only Vite dev server
- `npm run dev:convex` - Start only Convex dev server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## 🗄️ Database Schema

### Projects
- `title`, `description` (editable in frontend)
- `totalTasks`, `completedTasks`, `progress` (calculated)
- Auto-generated: `_id`, `_creationTime`

### Tasks
- `title`, `description` (editable in frontend)
- `state`, `complexity`, `depth` (LLM managed)
- `dependencies`, `dependents` (LLM managed)
- `positionX`, `positionY` (UI managed)
- `assignedAgent` (LLM managed)

### Agents
- `name`, `role`, `description`, `status`
- `capabilities`, `currentTasks`
- System managed only

### Events
- Real-time event log for live updates
- Auto-cleanup of old events

## 🔄 Real-time Features

### Live Updates
- Project statistics update in real-time
- Task state changes broadcast to all clients
- Agent status updates
- Drag & drop position sync

### Optimistic Updates
- UI updates immediately on user action
- Background sync with Convex
- Automatic conflict resolution

## 🎨 UI Components

All UI components are from shadcn/ui:
- Fully accessible (ARIA compliant)
- Dark/light mode support
- Customizable with CSS variables
- TypeScript ready

### Key Components Used
- `Button`, `Card`, `Badge` - Basic UI
- `Dialog`, `Sheet`, `Popover` - Overlays
- `Progress`, `Separator` - Data display
- `Toast`, `Alert` - Notifications
- `Tabs`, `Accordion` - Navigation

## 🔐 Environment Variables

Create `.env.local`:
```bash
VITE_CONVEX_URL=https://your-deployment.convex.cloud
```

## 🚨 Troubleshooting

### Convex Connection Issues
1. Check `.env.local` has correct `VITE_CONVEX_URL`
2. Ensure `npx convex dev` is running
3. Verify Convex dashboard shows active deployment

### Build Errors
1. Run `npm run lint` to check for TypeScript errors
2. Clear node_modules: `rm -rf node_modules && npm install`
3. Check all imports use correct paths with `@/` alias

### Real-time Not Working
1. Check browser console for WebSocket errors
2. Verify Convex functions are deployed
3. Check network tab for failed API calls

## 📚 Key Technologies

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool & dev server
- **Tailwind CSS 3** - Styling
- **shadcn/ui** - Component library
- **ReactFlow** - Flow diagrams
- **Zustand** - State management
- **Convex** - Real-time backend
- **Lucide React** - Icons

## 🎯 Development Workflow

1. **Frontend Changes**: Edit React components, hooks, stores
2. **Backend Changes**: Edit Convex functions in `convex/`
3. **Schema Changes**: Update `convex/schema.ts`
4. **UI Changes**: Use shadcn/ui components or add new ones
5. **Testing**: Use dummy data or connect to real Convex data

## 🔄 Data Flow

```
User Action → Zustand Store → Optimistic Update → Convex Mutation → Real-time Broadcast → All Clients Update
```

---

**Happy coding! 🚀**