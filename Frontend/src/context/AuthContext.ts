import { createContext } from 'react';
import api from '../utils/api';

export interface AuthContextType {
  isAuthenticated: boolean;
  accessToken: string | null;
  refreshToken: string | null;
  login: (accessToken: string, refreshToken: string) => void;
  logout: (logoutFromKeycloak?: boolean) => void;
  api: typeof api;
}

export const AuthContext = createContext<AuthContextType | null>(null);
