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
  // Only use cookies
  return getCookie('access_token');
};

const getRefreshToken = (): string | null => {
  // Only use cookies
  return getCookie('refresh_token');
};

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // Enable sending cookies with requests
});

let isRefreshing = false;
let failedQueue: { resolve: () => void; reject: (error: unknown) => void; }[] = [];

const processQueue = (error: unknown) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve();
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
    // No refresh token available - redirect to login
    window.location.href = '/login';
    throw new Error('No refresh token available');
  }

  console.log('Refreshing token...');

  try {
    // Cookies are sent automatically with withCredentials: true
    await axios.post(`${API_BASE_URL}/api/v1/auth/refreshToken`, {}, {
      withCredentials: true,
    });

    // Response is now 204 No Content, no JSON body
    // The new access token is set as a cookie automatically
    return getCookie('access_token') || '';
  } catch (error) {
    window.location.href = '/login';
    throw error;
  }
};

// Request interceptor to handle token refresh
api.interceptors.request.use(
  async (config) => {
    let token = getAccessToken();
    
    if (token && isTokenExpired(token)) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({
            resolve: () => {
              // Cookies are handled automatically, no need to set headers
              resolve(config);
            },
            reject
          });
        });
      }

      isRefreshing = true;
      try {
        token = await refreshAuthToken();
        processQueue(null);
      } catch (error) {
        processQueue(error);
        throw error;
      } finally {
        isRefreshing = false;
      }
    }

    // Cookies are sent automatically with withCredentials: true
    // No need to set Authorization headers
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