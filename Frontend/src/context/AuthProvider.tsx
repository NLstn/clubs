import React, { useState, useEffect } from 'react';
import { AuthContext } from './AuthContext';
import api from '../utils/api';

// Cookie utility functions
const getCookie = (name: string): string | null => {
  const nameEQ = name + "=";
  const ca = document.cookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
};

const getAccessToken = (): string | null => {
  // Try cookie first, then fall back to localStorage for backwards compatibility
  return getCookie('access_token') || localStorage.getItem('auth_token');
};

const getRefreshToken = (): string | null => {
  // Try cookie first, then fall back to localStorage for backwards compatibility
  return getCookie('refresh_token') || localStorage.getItem('refresh_token');
};

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
  // Initialize with tokens from cookies or localStorage
  const [accessToken, setAccessToken] = useState<string | null>(() => getAccessToken());
  const [refreshToken, setRefreshToken] = useState<string | null>(() => getRefreshToken());
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => !!getAccessToken());

  useEffect(() => {
    // Maintain localStorage for backwards compatibility when not using cookies
    if (accessToken) {
      if (!getCookie('access_token')) {
        localStorage.setItem('auth_token', accessToken);
      }
      if (refreshToken && !getCookie('refresh_token')) {
        localStorage.setItem('refresh_token', refreshToken);
      }
      setIsAuthenticated(true);
    } else {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('refresh_token');
      setIsAuthenticated(false);
    }
  }, [accessToken, refreshToken]);

  const login = (newAccessToken: string, newRefreshToken: string) => {
    setAccessToken(newAccessToken);
    setRefreshToken(newRefreshToken);
  };

  const logout = async () => {
    const currentRefreshToken = getRefreshToken();
    if (currentRefreshToken) {
      try {
        const headers: Record<string, string> = {
          'Content-Type': 'application/json',
        };
        
        // Only set Authorization header if not using cookies
        if (!getCookie('refresh_token')) {
          headers['Authorization'] = currentRefreshToken;
        }

        await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/logout`, {
          method: 'POST',
          headers,
          credentials: 'include', // Include cookies
        });
      } catch (error) {
        console.error('Error during logout:', error);
      }
    }
    setAccessToken(null);
    setRefreshToken(null);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, accessToken, refreshToken, login, logout, api }}>
      {children}
    </AuthContext.Provider>
  );
};