import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

const API_BASE_URL = import.meta.env.VITE_API_HOST || '';

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

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // Enable sending cookies with requests
});

let isRefreshing = false;
let failedQueue: { resolve: (token: string) => void; reject: (error: unknown) => void; }[] = [];

const processQueue = (error: unknown, token: string | null = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token!);
    }
  });
  failedQueue = [];
};

interface JWTPayload {
  exp: number;
}

const isTokenExpired = (token: string): boolean => {
  try {
    const decoded = jwtDecode<JWTPayload>(token);
    // Check if token will expire in the next 30 seconds
    return (decoded.exp * 1000) < (Date.now() + 30000);
  } catch {
    return true;
  }
};

const refreshAuthToken = async () => {
  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    // No refresh token available - logout and redirect to login
    localStorage.removeItem('auth_token');
    localStorage.removeItem('refresh_token');
    window.location.href = '/login';
    throw new Error('No refresh token available');
  }

  console.log('Refreshing token...');

  try {
    // With cookies enabled, we don't need to send the token in the header
    // The cookie will be sent automatically. But we keep header support for backwards compatibility.
    const headers: Record<string, string> = {};
    if (!getCookie('refresh_token')) {
      // Only set header if not using cookies
      headers['Authorization'] = refreshToken;
    }

    const response = await axios.post(`${API_BASE_URL}/api/v1/auth/refreshToken`, {}, {
      headers,
      withCredentials: true,
    });

    const { access: newAccessToken } = response.data;
    
    // Store in localStorage for backwards compatibility if cookies aren't being used
    if (!getCookie('access_token')) {
      localStorage.setItem('auth_token', newAccessToken);
    }
    
    return newAccessToken;
  } catch (error) {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('refresh_token');
    window.location.href = '/login';
    throw error;
  }
};

// Request interceptor to add auth token and handle token refresh
api.interceptors.request.use(
  async (config) => {
    let token = getAccessToken();
    
    if (token && isTokenExpired(token)) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({
            resolve: (newToken) => {
              // Only set Authorization header if not using cookies
              if (!getCookie('access_token')) {
                config.headers.Authorization = `Bearer ${newToken}`;
              }
              resolve(config);
            },
            reject
          });
        });
      }

      isRefreshing = true;
      try {
        token = await refreshAuthToken();
        processQueue(null, token);
      } catch (error) {
        processQueue(error, null);
        throw error;
      } finally {
        isRefreshing = false;
      }
    }

    // Only set Authorization header if not using cookies
    if (token && !getCookie('access_token')) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for other errors
api.interceptors.response.use(
  (response) => response,
  (error) => Promise.reject(error)
);

export default api;