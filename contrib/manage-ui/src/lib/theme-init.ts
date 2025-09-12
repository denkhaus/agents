'use client';

// Theme initialization - should be called on app startup
export function initializeTheme() {
  if (typeof window === 'undefined') return 'light';

  // Get stored theme or default to system preference
  const storedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null;
  const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  
  const theme = storedTheme || (systemPrefersDark ? 'dark' : 'light');
  
  // Apply theme to document immediately
  const htmlElement = document.documentElement;
  if (theme === 'dark') {
    htmlElement.classList.add('dark');
  } else {
    htmlElement.classList.remove('dark');
  }
  
  // Store the theme if it wasn't already stored
  if (!storedTheme) {
    localStorage.setItem('theme', theme);
  }
  
  console.log('Theme initialized:', theme, 'Dark class present:', htmlElement.classList.contains('dark'));
  
  return theme;
}

// Listen for system theme changes
export function setupThemeListener() {
  if (typeof window === 'undefined') return;

  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
  
  const handleChange = (e: MediaQueryListEvent) => {
    const storedTheme = localStorage.getItem('theme');
    
    // Only auto-switch if user hasn't manually set a preference
    if (!storedTheme) {
      const newTheme = e.matches ? 'dark' : 'light';
      
      if (newTheme === 'dark') {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    }
  };
  
  mediaQuery.addEventListener('change', handleChange);
  
  return () => mediaQuery.removeEventListener('change', handleChange);
}