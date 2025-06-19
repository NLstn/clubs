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
  // Only use cookies
  return getCookie('access_token');
};

const getRefreshToken = (): string | null => {
  // Only use cookies
  return getCookie('refresh_token');
};

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
  // Initialize with tokens from cookies or localStorage
  const [accessToken, setAccessToken] = useState<string | null>(() => getAccessToken());
  const [refreshToken, setRefreshToken] = useState<string | null>(() => getRefreshToken());
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => !!getAccessToken());

  useEffect(() => {
    // Update authentication state based on cookie presence
    if (accessToken) {
      setIsAuthenticated(true);
    } else {
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
        // Cookies are sent automatically with credentials: 'include'
        await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/logout`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
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