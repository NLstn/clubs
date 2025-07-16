import React, { useState, useEffect } from 'react';
import { AuthContext } from './AuthContext';
import api from '../utils/api';
import keycloakService from '../utils/keycloak';

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
    console.log('AuthProvider login called with:', {
      accessToken: newAccessToken ? 'Present' : 'Missing',
      refreshToken: newRefreshToken ? 'Present' : 'Missing'
    });
    setAccessToken(newAccessToken);
    setRefreshToken(newRefreshToken);
    console.log('AuthProvider tokens set');
  };

  const logout = async (logoutFromKeycloak: boolean = true) => {
    console.log('Logout called with logoutFromKeycloak:', logoutFromKeycloak);
    
    if (refreshToken) {
      try {
        if (logoutFromKeycloak) {
          // Use the dedicated Keycloak logout endpoint with ID token
          const idToken = localStorage.getItem('keycloak_id_token');
          const response = await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/auth/keycloak/logout`, {
            method: 'POST',
            headers: {
              'Authorization': refreshToken,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              post_logout_redirect_uri: window.location.origin + '/login',
              id_token: idToken
            })
          });

          if (response.ok) {
            const data = await response.json();
            console.log('Keycloak logout response:', data);
            if (data.logoutURL) {
              // Clear local state first
              setAccessToken(null);
              setRefreshToken(null);
              // Also clear all Keycloak user data
              keycloakService.clearAllKeycloakData();
              // Remove the stored ID token
              localStorage.removeItem('keycloak_id_token');
              // Set a flag to force fresh login next time
              localStorage.setItem('force_keycloak_login', 'true');
              console.log('Redirecting to Keycloak logout:', data.logoutURL);
              // Then redirect to Keycloak logout
              window.location.href = data.logoutURL;
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
    localStorage.removeItem('keycloak_id_token');
    // Set a flag to force fresh login next time
    localStorage.setItem('force_keycloak_login', 'true');
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, accessToken, refreshToken, login, logout, api }}>
      {children}
    </AuthContext.Provider>
  );
};