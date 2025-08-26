import React, { useState, useEffect } from 'react';
import { AuthContext } from './AuthContext';
import api from '../utils/api';
import keycloakService from '../utils/keycloak';
import storage from '../utils/isomorphicStorage';

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
  const [accessToken, setAccessToken] = useState<string | null>(storage.getItem('auth_token'));
  const [refreshToken, setRefreshToken] = useState<string | null>(storage.getItem('refresh_token'));
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(!!accessToken);

  useEffect(() => {
    if (accessToken) {
      storage.setItem('auth_token', accessToken);
      storage.setItem('refresh_token', refreshToken || '');
      setIsAuthenticated(true);
    } else {
      storage.removeItem('auth_token');
      storage.removeItem('refresh_token');
      setIsAuthenticated(false);
    }
  }, [accessToken, refreshToken]);

  const login = (newAccessToken: string, newRefreshToken: string) => {
    setAccessToken(newAccessToken);
    setRefreshToken(newRefreshToken);
  };

  const logout = async (logoutFromKeycloak: boolean = true) => {
    if (refreshToken) {
      try {
        if (logoutFromKeycloak) {
          // Use the dedicated Keycloak logout endpoint with ID token
          const idToken = storage.getItem('keycloak_id_token');
          const response = await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/keycloak/logout`, {
            method: 'POST',
            headers: {
              'Authorization': refreshToken,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              post_logout_redirect_uri: typeof window !== 'undefined' ? window.location.origin + '/login' : '/login',
              id_token: idToken
            })
          });

          if (response.ok) {
            const data = await response.json();
            if (data.logoutURL) {
              // Clear local state first
              setAccessToken(null);
              setRefreshToken(null);
              // Also clear all Keycloak user data
              keycloakService.clearAllKeycloakData();
              // Remove the stored ID token
              storage.removeItem('keycloak_id_token');
              // Set a flag to force fresh login next time
              storage.setItem('force_keycloak_login', 'true');
              // Then redirect to Keycloak logout
              if (typeof window !== 'undefined') {
                window.location.href = data.logoutURL;
              }
              return; // Don't clear state again below
            }
          }
        } else {
          // Regular logout without Keycloak
          await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/logout`, {
            method: 'POST',
            headers: {
              'Authorization': refreshToken
            }
          });
        }
      } catch (error) {
        console.error('Error during logout:', error);
      }
    }
    
    // Clear local state for non-Keycloak logout or if Keycloak logout failed
    setAccessToken(null);
    setRefreshToken(null);
    // Also clear all Keycloak user data
    keycloakService.clearAllKeycloakData();
    // Remove the stored ID token
    storage.removeItem('keycloak_id_token');
    // Set a flag to force fresh login next time
    storage.setItem('force_keycloak_login', 'true');
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, accessToken, refreshToken, login, logout, api }}>
      {children}
    </AuthContext.Provider>
  );
};