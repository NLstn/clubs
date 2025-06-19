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
  // Initialize with tokens from cookies
  const [accessToken, setAccessToken] = useState<string | null>(() => getAccessToken());
  const [refreshToken, setRefreshToken] = useState<string | null>(() => getRefreshToken());
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => !!getAccessToken());

  // Check for delayed cookie availability on startup
  useEffect(() => {
    // If we don't have tokens initially, check a few times in case cookies are delayed
    if (!accessToken) {
      let attempts = 0;
      const maxAttempts = 10; // Check for up to 2 seconds (10 * 200ms)
      
      const checkForCookies = () => {
        const currentAccessToken = getAccessToken();
        const currentRefreshToken = getRefreshToken();
        
        if (currentAccessToken) {
          // Found cookies! Update state
          setAccessToken(currentAccessToken);
          setRefreshToken(currentRefreshToken);
        } else if (attempts < maxAttempts) {
          attempts++;
          setTimeout(checkForCookies, 200);
        }
      };
      
      // Start checking after a small delay
      setTimeout(checkForCookies, 100);
    }
  }, []); // Only run once on mount

  useEffect(() => {
    // Update authentication state based on cookie presence
    if (accessToken) {
      setIsAuthenticated(true);
    } else {
      setIsAuthenticated(false);
    }
  }, [accessToken]); // Only depend on accessToken since that's what we check

  const login = () => {
    // With cookie-only auth, tokens are set automatically as cookies
    // Add a small delay to ensure cookies are available after being set by the browser
    const checkCookiesWithRetry = (attempts = 0) => {
      const currentAccessToken = getAccessToken();
      const currentRefreshToken = getRefreshToken();
      
      if (currentAccessToken || attempts >= 5) {
        // Either we found tokens or we've exhausted our attempts
        setAccessToken(currentAccessToken);
        setRefreshToken(currentRefreshToken);
      } else {
        // Retry after a short delay
        setTimeout(() => checkCookiesWithRetry(attempts + 1), 50);
      }
    };
    
    checkCookiesWithRetry();
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