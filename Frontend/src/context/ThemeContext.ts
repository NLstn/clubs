import { createContext } from 'react';

export type ThemeMode = 'light' | 'dark' | 'system';

export interface ThemeContextType {
  theme: ThemeMode;
  setTheme: (theme: ThemeMode) => void;
  effectiveTheme: 'light' | 'dark'; // The actual theme being applied (resolved from system if needed)
}

export const ThemeContext = createContext<ThemeContextType | undefined>(undefined);
