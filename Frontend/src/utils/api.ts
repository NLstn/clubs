import axios from 'axios';
import { jwtDecode } from 'jwt-decode';
import storage from './isomorphicStorage';

const API_BASE_URL = import.meta.env.VITE_API_HOST || '';

const api = axios.create({
  baseURL: API_BASE_URL,
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
  const refreshToken = storage.getItem('refresh_token');
  if (!refreshToken) {
    // No refresh token available - logout and redirect to login
    storage.removeItem('auth_token');
    storage.removeItem('refresh_token');
    if (typeof window !== 'undefined') {
      window.location.href = '/login';
    }
    throw new Error('No refresh token available');
  }

  try {
    const response = await axios.post(`${API_BASE_URL}/api/v1/auth/refreshToken`, {}, {
      headers: {
        'Authorization': refreshToken
      }
    });

    const { access: newAccessToken, refresh: newRefreshToken } = response.data;
    storage.setItem('auth_token', newAccessToken);
    storage.setItem('refresh_token', newRefreshToken);
    return newAccessToken;
  } catch (error) {
    storage.removeItem('auth_token');
    storage.removeItem('refresh_token');
    if (typeof window !== 'undefined') {
      window.location.href = '/login';
    }
    throw error;
  }
};

// Request interceptor to add auth token and handle token refresh
api.interceptors.request.use(
  async (config) => {
    let token = storage.getItem('auth_token');
    
    if (token && isTokenExpired(token)) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({
            resolve: (newToken) => {
              config.headers.Authorization = `Bearer ${newToken}`;
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

    if (token) {
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

// Club API functions
export const hardDeleteClub = async (clubId: string) => {
  return api.delete(`/api/v1/clubs/${clubId}/hard-delete`);
};

export default api;