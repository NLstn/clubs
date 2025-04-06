import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_HOST || '';

const api = axios.create({
  baseURL: API_BASE_URL,
});

// Request interceptor to add auth token to all requests
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export default api;