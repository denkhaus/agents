# Project Canvas Setup Guide

## 🚀 Quick Start

### 1. Install Dependencies
```bash
cd project-canvas
npm install
```

### 2. Initialize shadcn/ui Components
Run these commands to add all necessary shadcn/ui components:

```bash
# Core UI Components
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
npx shadcn-ui@latest add input
npx shadcn-ui@latest add label
npx shadcn-ui@latest add textarea
npx shadcn-ui@latest add select
npx shadcn-ui@latest add checkbox
npx shadcn-ui@latest add switch
npx shadcn-ui@latest add slider
npx shadcn-ui@latest add progress

# Navigation & Layout
npx shadcn-ui@latest add navigation-menu
npx shadcn-ui@latest add menubar
npx shadcn-ui@latest add tabs
npx shadcn-ui@latest add separator
npx shadcn-ui@latest add scroll-area
npx shadcn-ui@latest add resizable

# Overlays & Dialogs
npx shadcn-ui@latest add dialog
npx shadcn-ui@latest add sheet
npx shadcn-ui@latest add popover
npx shadcn-ui@latest add dropdown-menu
npx shadcn-ui@latest add hover-card
npx shadcn-ui@latest add tooltip
npx shadcn-ui@latest add alert-dialog

# Data Display
npx shadcn-ui@latest add table
npx shadcn-ui@latest add avatar
npx shadcn-ui@latest add badge
npx shadcn-ui@latest add accordion
npx shadcn-ui@latest add collapsible

# Feedback & Status
npx shadcn-ui@latest add toast
npx shadcn-ui@latest add alert
npx shadcn-ui@latest add skeleton

# Forms & Input
npx shadcn-ui@latest add form
npx shadcn-ui@latest add radio-group
npx shadcn-ui@latest add toggle
npx shadcn-ui@latest add command
npx shadcn-ui@latest add calendar
npx shadcn-ui@latest add date-picker
npx shadcn-ui@latest add input-otp

# Advanced Components
npx shadcn-ui@latest add carousel
npx shadcn-ui@latest add drawer
npx shadcn-ui@latest add sonner
```

### 3. Setup Convex (Optional - for later)
```bash
# Initialize Convex
npx convex dev --once
```

### 4. Start Development Server
```bash
npm run dev
```

## 📁 Project Structure

```
project-canvas/
├── src/
│   ├── components/
│   │   ├── ui/           # shadcn/ui components (auto-generated)
│   │   ├── layout/       # Layout components
│   │   ├── canvas/       # ReactFlow canvas components
│   │   └── workspace/    # Workspace-specific components
│   ├── lib/
│   │   └── utils.ts      # Utility functions
│   ├── stores/           # Zustand stores
│   ├── types/            # TypeScript type definitions
│   ├── data/             # Dummy data and helpers
│   ├── utils/            # Utility functions
│   └── app/
│       └── globals.css   # Global styles
├── components.json       # shadcn/ui configuration
├── tailwind.config.js    # Tailwind CSS configuration
├── vite.config.ts        # Vite configuration
└── package.json          # Dependencies
```

## 🎨 Theme Configuration

The project uses a custom theme with:
- **Primary Color**: Blue (`hsl(221.2 83.2% 53.3%)`)
- **Dark Mode**: Automatic class-based switching
- **Typography**: Tailwind Typography plugin for markdown
- **Animations**: Tailwind Animate for smooth transitions

## 🔧 Key Features Configured

### ✅ Vite Setup
- React + TypeScript
- Path aliases (`@/` → `./src/`)
- Development server on port 3000
- Source maps enabled

### ✅ Tailwind CSS 3
- CSS variables for theming
- Dark mode support
- Custom scrollbar styles
- ReactFlow integration styles
- Typography plugin for markdown

### ✅ shadcn/ui Ready
- All components pre-configured
- Consistent design system
- Accessible components
- Dark mode compatible

### ✅ TypeScript Configuration
- Strict mode enabled
- Path mapping configured
- Modern ES2020 target
- Bundler module resolution

## 🚦 Next Steps

After running the setup commands:

1. **Verify Installation**: `npm run dev` should start without errors
2. **Check Components**: All shadcn/ui components should be in `src/components/ui/`
3. **Test Dark Mode**: Theme switching should work
4. **Validate Types**: TypeScript should compile without errors

## 🐛 Troubleshooting

### Common Issues:

1. **Path Resolution Errors**
   ```bash
   # Make sure @types/node is installed
   npm install -D @types/node
   ```

2. **shadcn/ui Component Errors**
   ```bash
   # Reinstall if components are missing
   npx shadcn-ui@latest add [component-name]
   ```

3. **Tailwind Not Working**
   ```bash
   # Rebuild CSS
   npm run build
   ```

4. **TypeScript Errors**
   ```bash
   # Clear cache and restart
   rm -rf node_modules/.vite
   npm run dev
   ```

## 📚 Documentation Links

- [Vite Documentation](https://vitejs.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [shadcn/ui](https://ui.shadcn.com/)
- [ReactFlow](https://reactflow.dev/)
- [Zustand](https://zustand-demo.pmnd.rs/)
- [Convex](https://docs.convex.dev/)

---

**Ready to build! 🎉**