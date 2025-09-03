import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

interface ThemeStore {
  // State
  isDarkMode: boolean
  
  // Actions
  toggleTheme: () => void
  setTheme: (isDark: boolean) => void
}

export const useThemeStore = create<ThemeStore>()(
  persist(
    (set) => ({
      // Initial state - defaults to dark mode
      isDarkMode: true,
      
      // Actions
      toggleTheme: () => set((state) => ({ 
        isDarkMode: !state.isDarkMode 
      })),
      
      setTheme: (isDark) => set({ isDarkMode: isDark })
    }),
    {
      name: 'theme-storage',
      storage: createJSONStorage(() => localStorage),
    }
  )
)