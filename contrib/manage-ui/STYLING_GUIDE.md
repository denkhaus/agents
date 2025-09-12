# 🎨 shadcn/ui Styling Guide

## 1. **CSS-Variablen ändern** (Empfohlen)

### Beispiel: Blaues Theme zu Grünem Theme
```css
/* In src/app/globals.css */
:root {
  /* Primärfarbe von Blau zu Grün ändern */
  --primary: oklch(0.6 0.2 142); /* Grün statt Blau */
  --primary-foreground: oklch(0.985 0 0);
  
  /* Akzentfarbe anpassen */
  --accent: oklch(0.9 0.1 142); /* Helles Grün */
  --accent-foreground: oklch(0.2 0.1 142);
  
  /* Border-Radius für rundere Ecken */
  --radius: 1rem; /* Statt 0.625rem */
}
```

### Beispiel: Lila Corporate Theme
```css
:root {
  --primary: oklch(0.5 0.2 270); /* Lila */
  --primary-foreground: oklch(0.985 0 0);
  --secondary: oklch(0.95 0.05 270); /* Helles Lila */
  --accent: oklch(0.9 0.1 270);
  --border: oklch(0.9 0.02 270);
}
```

## 2. **Komponenten-spezifisches Styling**

### Button-Varianten erweitern
```typescript
// In src/components/ui/button.tsx
const buttonVariants = cva(
  "inline-flex items-center justify-center...",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        // Neue Variante hinzufügen
        gradient: "bg-gradient-to-r from-purple-500 to-pink-500 text-white hover:from-purple-600 hover:to-pink-600",
        neon: "bg-cyan-500 text-black hover:bg-cyan-400 shadow-lg shadow-cyan-500/50",
      },
      size: {
        default: "h-10 px-4 py-2",
        // Neue Größe hinzufügen
        xl: "h-14 px-8 py-4 text-lg",
      }
    }
  }
)
```

### Verwendung:
```tsx
<Button variant="gradient" size="xl">Gradient Button</Button>
<Button variant="neon">Neon Button</Button>
```

## 3. **Tailwind CSS Klassen überschreiben**

### Direkte Klassen-Überschreibung
```tsx
<Button className="bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700">
  Custom Button
</Button>
```

### Mit cn() Utility für bessere Kontrolle
```tsx
import { cn } from "@/lib/utils"

<Button 
  className={cn(
    "bg-emerald-500 hover:bg-emerald-600", 
    "shadow-lg shadow-emerald-500/25",
    "border-2 border-emerald-400"
  )}
>
  Emerald Button
</Button>
```

## 4. **Theme-Presets erstellen**

### Theme-Switcher implementieren
```typescript
// src/lib/themes.ts
export const themes = {
  default: {
    primary: "oklch(0.21 0.006 285.885)",
    secondary: "oklch(0.967 0.001 286.375)",
    accent: "oklch(0.967 0.001 286.375)",
  },
  ocean: {
    primary: "oklch(0.5 0.2 220)", // Blau
    secondary: "oklch(0.9 0.1 220)",
    accent: "oklch(0.7 0.15 200)",
  },
  forest: {
    primary: "oklch(0.4 0.2 142)", // Grün
    secondary: "oklch(0.9 0.1 142)",
    accent: "oklch(0.6 0.15 120)",
  },
  sunset: {
    primary: "oklch(0.6 0.2 30)", // Orange
    secondary: "oklch(0.9 0.1 30)",
    accent: "oklch(0.7 0.15 50)",
  }
}
```

### Theme anwenden
```typescript
function applyTheme(themeName: keyof typeof themes) {
  const theme = themes[themeName]
  const root = document.documentElement
  
  Object.entries(theme).forEach(([key, value]) => {
    root.style.setProperty(`--${key}`, value)
  })
}
```

## 5. **Spezifische Komponenten anpassen**

### Task Cards mit benutzerdefinierten Farben
```tsx
// In TaskCard Komponente
const getTaskCardStyle = (state: TaskState) => {
  switch (state) {
    case 'completed':
      return "bg-gradient-to-br from-green-50 to-emerald-50 border-green-200"
    case 'in-progress':
      return "bg-gradient-to-br from-blue-50 to-cyan-50 border-blue-200"
    case 'blocked':
      return "bg-gradient-to-br from-red-50 to-rose-50 border-red-200"
    default:
      return "bg-white border-gray-200"
  }
}

<div className={cn("rounded-lg p-4", getTaskCardStyle(task.state))}>
  {/* Task content */}
</div>
```

### Sidebar mit Custom Styling
```tsx
// In Sidebar Komponente
<div className={cn(
  "flex flex-col transition-all duration-300",
  "bg-gradient-to-b from-slate-900 to-slate-800", // Custom Gradient
  "border-r border-slate-700",
  sidebarCollapsed ? "w-16" : "w-64"
)}>
```

## 6. **Animationen hinzufügen**

### Hover-Effekte für Cards
```css
/* In globals.css */
.task-card {
  @apply transition-all duration-200 ease-in-out;
}

.task-card:hover {
  @apply transform scale-105 shadow-xl;
}

.task-card:hover .task-title {
  @apply text-primary;
}
```

### Framer Motion für erweiterte Animationen
```tsx
import { motion } from "framer-motion"

<motion.div
  initial={{ opacity: 0, y: 20 }}
  animate={{ opacity: 1, y: 0 }}
  transition={{ duration: 0.3 }}
  className="task-card"
>
  {/* Task content */}
</motion.div>
```

## 7. **Dark Mode anpassen**

### Custom Dark Mode Farben
```css
.dark {
  --background: oklch(0.1 0.02 270); /* Dunkles Lila statt Grau */
  --foreground: oklch(0.95 0.02 270);
  --primary: oklch(0.7 0.2 270); /* Helles Lila */
  --card: oklch(0.15 0.02 270);
}
```

## 8. **Praktische Beispiele für Ihr Projekt**

### 1. Agent-Status Farben
```typescript
export const getAgentStatusColor = (status: string) => {
  const colors = {
    online: "bg-green-100 text-green-800 border-green-200",
    busy: "bg-yellow-100 text-yellow-800 border-yellow-200",
    offline: "bg-gray-100 text-gray-800 border-gray-200",
    error: "bg-red-100 text-red-800 border-red-200"
  }
  return colors[status] || colors.offline
}
```

### 2. Projekt-Progress mit Gradient
```tsx
<div className="w-full bg-gray-200 rounded-full h-3">
  <div 
    className="bg-gradient-to-r from-blue-500 to-purple-600 h-3 rounded-full transition-all duration-500"
    style={{ width: `${progress}%` }}
  />
</div>
```

### 3. Kanban-Spalten mit Theme-Farben
```tsx
const getColumnStyle = (state: TaskState) => {
  const styles = {
    pending: "bg-gradient-to-b from-gray-50 to-gray-100 border-gray-200",
    'in-progress': "bg-gradient-to-b from-blue-50 to-blue-100 border-blue-200",
    completed: "bg-gradient-to-b from-green-50 to-green-100 border-green-200",
    blocked: "bg-gradient-to-b from-red-50 to-red-100 border-red-200",
    cancelled: "bg-gradient-to-b from-gray-50 to-gray-100 border-gray-300"
  }
  return styles[state]
}
```

## 9. **Live Theme Editor** (Erweitert)

```tsx
// Theme Editor Komponente
export function ThemeEditor() {
  const [primaryColor, setPrimaryColor] = useState("#3b82f6")
  
  useEffect(() => {
    document.documentElement.style.setProperty(
      '--primary', 
      `oklch(0.5 0.2 ${hexToHue(primaryColor)})`
    )
  }, [primaryColor])
  
  return (
    <div className="p-4 space-y-4">
      <label>
        Primary Color:
        <input 
          type="color" 
          value={primaryColor}
          onChange={(e) => setPrimaryColor(e.target.value)}
        />
      </label>
      {/* Weitere Farb-Controls */}
    </div>
  )
}
```

## 10. **Schnelle Änderungen testen**

### Browser DevTools verwenden
1. F12 öffnen
2. Elements Tab
3. `:root` Element finden
4. CSS-Variablen live ändern:
```css
--primary: oklch(0.6 0.3 120); /* Grün */
--radius: 1.5rem; /* Rundere Ecken */
```

### Empfohlene Workflow
1. **DevTools** für schnelle Tests
2. **CSS-Variablen** für Theme-Änderungen
3. **Komponenten-Varianten** für neue Styles
4. **Tailwind-Klassen** für spezifische Anpassungen

---

**Tipp**: Beginnen Sie mit CSS-Variablen für globale Änderungen und verwenden Sie Komponenten-Varianten für spezifische UI-Elemente!