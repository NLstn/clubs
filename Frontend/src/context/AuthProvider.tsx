import React, { useState, useEffect } from 'react';
import { AuthContext } from './AuthContext';
import api from '../utils/api';

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
  const [accessToken, setAccessToken] = useState<string | null>(localStorage.getItem('auth_token'));
  const [refreshToken, setRefreshToken] = useState<string | null>(localStorage.getItem('refresh_token'));
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(!!accessToken);

  useEffect(() => {
    if (accessToken) {
      localStorage.setItem('auth_token', accessToken);
      localStorage.setItem('refresh_token', refreshToken || '');
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
    if (refreshToken) {
      try {
        await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/logout`, {
          method: 'POST',
          headers: {
            'Authorization': refreshToken
          }
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