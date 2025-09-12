# 🎨 Schnelle Styling-Beispiele

## ✅ **Sofort verfügbar:**

### 1. **Theme-Switcher** (Bereits integriert)
- **Palette-Icon** in der Top-Navigation (neben Dark/Light Toggle)
- **5 vorgefertigte Themes**:
  - 🔵 **Default** (Grau/Blau)
  - 🌊 **Ocean Blue** (Blau-Töne)
  - 🌲 **Forest Green** (Grün-Töne)
  - 🌅 **Sunset Orange** (Orange-Töne)
  - 👑 **Royal Purple** (Lila-Töne)

### 2. **Sofort testen:**
1. Klicken Sie auf das **Palette-Icon** (🎨) oben rechts
2. Wählen Sie ein Theme aus
3. **Sofortige Änderung** aller UI-Elemente

## 🎯 **Weitere Anpassungen:**

### Task-Status Farben ändern
```typescript
// In src/lib/utils.ts - getTaskStateColor Funktion
export function getTaskStateColor(state: TaskState): string {
  switch (state) {
    case 'pending':
      return 'bg-amber-100 text-amber-800 border-amber-200'; // Gelb statt Grau
    case 'in-progress':
      return 'bg-cyan-100 text-cyan-800 border-cyan-200'; // Cyan statt Blau
    case 'completed':
      return 'bg-emerald-100 text-emerald-800 border-emerald-200'; // Smaragd
    case 'blocked':
      return 'bg-rose-100 text-rose-800 border-rose-200'; // Rose statt Rot
    case 'cancelled':
      return 'bg-slate-100 text-slate-600 border-slate-200'; // Slate
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200';
  }
}
```

### Kanban-Spalten mit Gradienten
```typescript
// In KanbanColumn Komponente hinzufügen:
const getColumnGradient = (state?: TaskState) => {
  switch (state) {
    case 'pending':
      return 'bg-gradient-to-b from-yellow-50 to-amber-50';
    case 'in-progress':
      return 'bg-gradient-to-b from-blue-50 to-cyan-50';
    case 'completed':
      return 'bg-gradient-to-b from-green-50 to-emerald-50';
    case 'blocked':
      return 'bg-gradient-to-b from-red-50 to-rose-50';
    case 'cancelled':
      return 'bg-gradient-to-b from-gray-50 to-slate-50';
    default:
      return 'bg-gray-50';
  }
};

// Dann in der className:
<div className={cn(
  "flex flex-col h-full rounded-lg",
  getColumnGradient(state)
)}>
```

### Task-Cards mit Hover-Effekten
```typescript
// In TaskCard Komponente:
<div className={cn(
  "bg-white rounded-lg border border-gray-200 p-4",
  "hover:shadow-lg hover:scale-[1.02] transition-all duration-200",
  "hover:border-primary/50"
)}>
```

### Sidebar mit Custom Gradient
```typescript
// In Sidebar Komponente:
<div className={cn(
  "flex flex-col transition-all duration-300",
  "bg-gradient-to-b from-primary/5 to-primary/10", // Verwendet Theme-Farbe
  "border-r border-primary/20",
  sidebarCollapsed ? "w-16" : "w-64"
)}>
```

## 🚀 **Live-Änderungen testen:**

### Browser DevTools (Sofort sichtbar)
1. **F12** drücken
2. **Elements** Tab
3. **:root** Element finden
4. CSS-Variablen ändern:

```css
/* Beispiele zum Kopieren in DevTools */
--primary: oklch(0.6 0.3 120);     /* Grün */
--primary: oklch(0.5 0.3 270);     /* Lila */
--primary: oklch(0.6 0.3 30);      /* Orange */
--radius: 1.5rem;                  /* Rundere Ecken */
--radius: 0.25rem;                 /* Eckigere Ecken */
```

### Komplettes Custom Theme
```css
/* In DevTools :root einfügen */
--primary: oklch(0.45 0.25 310);        /* Pink-Lila */
--secondary: oklch(0.92 0.05 310);      /* Helles Pink */
--accent: oklch(0.7 0.2 330);           /* Pink-Akzent */
--background: oklch(0.99 0.01 310);     /* Sehr helles Pink */
--radius: 1.2rem;                       /* Runde Ecken */
```

## 🎨 **Erweiterte Anpassungen:**

### 1. Neue Button-Variante hinzufügen
```typescript
// In src/components/ui/button.tsx
const buttonVariants = cva(
  // ... existing code
  {
    variants: {
      variant: {
        // ... existing variants
        neon: "bg-cyan-500 text-black hover:bg-cyan-400 shadow-lg shadow-cyan-500/50 font-bold",
        gradient: "bg-gradient-to-r from-purple-500 to-pink-500 text-white hover:from-purple-600 hover:to-pink-600",
      }
    }
  }
)
```

### 2. Verwendung der neuen Varianten
```tsx
<Button variant="neon">Neon Button</Button>
<Button variant="gradient">Gradient Button</Button>
```

### 3. Task-Cards mit Status-spezifischen Styles
```typescript
// In TaskCard Komponente
const getTaskCardStyle = (state: TaskState) => {
  const baseStyle = "bg-white rounded-lg border p-4 transition-all duration-200";
  
  switch (state) {
    case 'completed':
      return cn(baseStyle, "border-green-200 hover:shadow-green-100 hover:shadow-lg");
    case 'in-progress':
      return cn(baseStyle, "border-blue-200 hover:shadow-blue-100 hover:shadow-lg");
    case 'blocked':
      return cn(baseStyle, "border-red-200 hover:shadow-red-100 hover:shadow-lg");
    default:
      return cn(baseStyle, "border-gray-200 hover:shadow-lg");
  }
};
```

## 🔥 **Sofort-Tipps:**

### 1. **Theme wechseln:** Palette-Icon → Theme auswählen
### 2. **DevTools:** F12 → :root → CSS-Variablen ändern
### 3. **Komponenten:** className mit cn() überschreiben
### 4. **Neue Varianten:** buttonVariants erweitern

---

**🎨 Experimentieren Sie mit den Themes und sehen Sie sofort die Änderungen!**