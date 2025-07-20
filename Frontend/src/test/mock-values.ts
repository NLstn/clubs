import { vi } from 'vitest';

// Mock auth context values
export const defaultMockAuthValue = {
  isAuthenticated: true,
  isLoading: false,
  logout: vi.fn()
};

export const unauthenticatedMockAuthValue = {
  isAuthenticated: false,
  isLoading: false,
  logout: vi.fn()
};

export const loadingMockAuthValue = {
  isAuthenticated: false,
  isLoading: true,
  logout: vi.fn()
};

// Mock notification hook values
export const defaultMockNotificationValue = {
  notifications: [],
  unreadCount: 0,
  loading: false,
  error: null,
  markAsRead: vi.fn(),
  markAllAsRead: vi.fn(),
  deleteNotification: vi.fn(),
  refreshNotifications: vi.fn(),
  fetchNotifications: vi.fn(),
  fetchUnreadCount: vi.fn()
};

export const mockNotificationsWithData = {
  notifications: [
    {
      id: '1',
      userId: 'user-1',
      type: 'info',
      title: 'Welcome',
      message: 'Welcome to the club!',
      read: false,
      createdAt: '2024-01-01T10:00:00Z',
      clubId: 'club-1'
    }
  ],
  unreadCount: 1,
  loading: false,
  error: null,
  markAsRead: vi.fn(),
  markAllAsRead: vi.fn(),
  deleteNotification: vi.fn(),
  refreshNotifications: vi.fn(),
  fetchNotifications: vi.fn(),
  fetchUnreadCount: vi.fn()
};

// Helper to mock localStorage
export const mockLocalStorage = () => {
  const store: Record<string, string> = {};
  
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value;
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key];
    }),
    clear: vi.fn(() => {
      Object.keys(store).forEach(key => delete store[key]);
    }),
    length: 0,
    key: vi.fn()
  };
};

// Mock window.location
export const mockWindowLocation = () => {
  const mockLocation = {
    href: 'http://localhost:3000',
    origin: 'http://localhost:3000',
    protocol: 'http:',
    host: 'localhost:3000',
    hostname: 'localhost',
    port: '3000',
    pathname: '/',
    search: '',
    hash: '',
    assign: vi.fn(),
    replace: vi.fn(),
    reload: vi.fn()
  };

  Object.defineProperty(window, 'location', {
    value: mockLocation,
    writable: true
  });

  return mockLocation;
};
