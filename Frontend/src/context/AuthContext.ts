import { createContext } from 'react';
import api from '../utils/api';

export interface AuthContextType {
  isAuthenticated: boolean;
  accessToken: string | null;
  refreshToken: string | null;
  login: () => void;
  logout: () => void;
  api: typeof api;
}

export const AuthContext = createContext<AuthContextType | null>(null);
