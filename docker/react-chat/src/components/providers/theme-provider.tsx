'use client'

import * as React from 'react'
import { useThemeStore } from '@/lib/store'

type Theme = 'dark' | 'light'

type ThemeProviderContextType = {
  theme: Theme
  setTheme: (theme: Theme) => void
  toggleTheme: () => void
}

const ThemeProviderContext = React.createContext<ThemeProviderContextType | undefined>(undefined)

export function ThemeProvider({
  children,
  ...props
}: {
  children: React.ReactNode
}) {
  const { isDarkMode, setTheme: setStoreTheme, toggleTheme: toggleStoreTheme } = useThemeStore()
  
  const theme: Theme = isDarkMode ? 'dark' : 'light'
  
  React.useEffect(() => {
    const root = window.document.documentElement
    
    root.classList.remove('light', 'dark')
    root.classList.add(theme)
  }, [theme])

  const setTheme = (newTheme: Theme) => {
    setStoreTheme(newTheme === 'dark')
  }

  const toggleTheme = () => {
    toggleStoreTheme()
  }

  const value = {
    theme,
    setTheme,
    toggleTheme,
  }

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {children}
    </ThemeProviderContext.Provider>
  )
}

export const useTheme = () => {
  const context = React.useContext(ThemeProviderContext)

  if (context === undefined)
    throw new Error('useTheme must be used within a ThemeProvider')

  return context
}