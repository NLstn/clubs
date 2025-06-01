import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

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
  const refreshToken = localStorage.getItem('refresh_token');
  if (!refreshToken) {
    // No refresh token available - logout and redirect to login
    localStorage.removeItem('auth_token');
    localStorage.removeItem('refresh_token');
    window.location.href = '/login';
    throw new Error('No refresh token available');
  }

  console.log('Refreshing token...');

  try {
    const response = await axios.post(`${API_BASE_URL}/api/v1/auth/refreshToken`, {}, {
      headers: {
        'Authorization': refreshToken
      }
    });

    const { access: newAccessToken } = response.data;
    localStorage.setItem('auth_token', newAccessToken);
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
    let token = localStorage.getItem('auth_token');
    
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

export default api;