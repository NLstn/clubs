import MockAdapter from 'axios-mock-adapter';
import axios from 'axios';
import api from '../utils/api';

// Create mock adapters
export const mockApi = new MockAdapter(api);
export const mockAxios = new MockAdapter(axios);

// Mock data for notifications
export const mockNotifications = [
  {
    id: '1',
    userId: 'user-1',
    type: 'info',
    title: 'Welcome',
    message: 'Welcome to the club!',
    read: false,
    createdAt: '2024-01-01T10:00:00Z',
    clubId: 'club-1'
  },
  {
    id: '2',
    userId: 'user-1',
    type: 'warning',
    title: 'Event Reminder',
    message: 'Don\'t forget about the upcoming event!',
    read: true,
    createdAt: '2024-01-02T10:00:00Z',
    eventId: 'event-1'
  }
];

export const mockClubs = [
  {
    id: 'club-1',
    name: 'Test Club 1',
    description: 'A test club',
    memberCount: 10,
    createdAt: '2024-01-01T10:00:00Z'
  },
  {
    id: 'club-2',
    name: 'Test Club 2',
    description: 'Another test club',
    memberCount: 5,
    createdAt: '2024-01-02T10:00:00Z'
  }
];

export const mockUsers = [
  {
    id: 'user-1',
    firstName: 'John',
    lastName: 'Doe',
    email: 'john.doe@example.com',
    username: 'johndoe'
  }
];

export const mockActivities = [
  {
    id: 'activity-1',
    type: 'member_added',
    description: 'John Doe joined the club',
    createdAt: '2024-01-01T10:00:00Z',
    clubId: 'club-1'
  }
];

export const mockFines = [
  {
    id: 'fine-1',
    memberId: 'user-1',
    amount: 25.00,
    reason: 'Late to practice',
    dueDate: '2024-02-01',
    paid: false,
    createdAt: '2024-01-01T10:00:00Z'
  }
];

export const mockTeams = [
  {
    id: 'team-1',
    name: 'Team A',
    description: 'First team',
    memberCount: 15,
    clubId: 'club-1'
  }
];

// Setup default mocks for common API endpoints
export const setupDefaultApiMocks = () => {
  // Reset all mocks
  mockApi.reset();
  mockAxios.reset();

  // Notifications endpoints
  mockApi.onGet('/api/v1/notifications').reply(200, mockNotifications);
  mockApi.onGet('/api/v1/notifications/count').reply(200, { count: 1 });
  mockApi.onPut(/\/api\/v1\/notifications\/[\w-]+/).reply(200);
  mockApi.onDelete(/\/api\/v1\/notifications\/[\w-]+/).reply(200);
  mockApi.onPost('/api/v1/notifications/mark-all-read').reply(200);

  // Auth endpoints
  mockApi.onPost('/api/v1/auth/refreshToken').reply(200, {
    access: 'mock-access-token',
    refresh: 'mock-refresh-token'
  });
  mockAxios.onPost('http://localhost:8080/api/v1/auth/refreshToken').reply(200, {
    access: 'mock-access-token',
    refresh: 'mock-refresh-token'
  });
  mockApi.onPost('/api/v1/auth/logout').reply(200);

  // User endpoints  
  mockApi.onGet('/api/v1/users/me').reply(200, mockUsers[0]);
  mockApi.onPut('/api/v1/users/me').reply(200, mockUsers[0]);

  // Club endpoints
  mockApi.onGet('/api/v1/clubs').reply(200, mockClubs);
  mockApi.onGet(/\/api\/v1\/clubs\/[\w-]+/).reply(200, mockClubs[0]);
  mockApi.onPost('/api/v1/clubs').reply(201, mockClubs[0]);
  mockApi.onPut(/\/api\/v1\/clubs\/[\w-]+/).reply(200, mockClubs[0]);
  mockApi.onDelete(/\/api\/v1\/clubs\/[\w-]+/).reply(200);

  // Dashboard endpoints
  mockApi.onGet('/api/v1/dashboard/activity').reply(200, mockActivities);

  // Club-specific endpoints
  mockApi.onGet(/\/api\/v1\/clubs\/[\w-]+\/members/).reply(200, mockUsers);
  mockApi.onGet(/\/api\/v1\/clubs\/[\w-]+\/fines/).reply(200, mockFines);
  mockApi.onGet(/\/api\/v1\/clubs\/[\w-]+\/teams/).reply(200, mockTeams);
  mockApi.onPost(/\/api\/v1\/clubs\/[\w-]+\/teams/).reply(201, mockTeams[0]);
  mockApi.onPut(/\/api\/v1\/clubs\/[\w-]+\/teams\/[\w-]+/).reply(200, mockTeams[0]);
  mockApi.onDelete(/\/api\/v1\/clubs\/[\w-]+\/teams\/[\w-]+/).reply(200);

  // Fine endpoints
  mockApi.onPost(/\/api\/v1\/clubs\/[\w-]+\/fines/).reply(201, mockFines[0]);
  mockApi.onPut(/\/api\/v1\/clubs\/[\w-]+\/fines\/[\w-]+/).reply(200, mockFines[0]);
  mockApi.onDelete(/\/api\/v1\/clubs\/[\w-]+\/fines\/[\w-]+/).reply(200);

  // Events endpoints
  mockApi.onGet(/\/api\/v1\/clubs\/[\w-]+\/events/).reply(200, []);
  
  // Join requests endpoints
  mockApi.onGet(/\/api\/v1\/clubs\/[\w-]+\/join-requests/).reply(200, []);
  
  // Invites endpoints
  mockApi.onGet(/\/api\/v1\/clubs\/[\w-]+\/invites/).reply(200, []);

  // Default fallback for any unmatched requests
  mockApi.onAny().reply(404, { message: 'API endpoint not mocked' });
  mockAxios.onAny().reply(404, { message: 'Axios endpoint not mocked' });
};

// Setup error mocks for testing error states
export const setupErrorApiMocks = () => {
  mockApi.reset();
  mockAxios.reset();
  mockApi.onAny().networkError();
  mockAxios.onAny().networkError();
};

// Clean up mocks
export const cleanupApiMocks = () => {
  mockApi.reset();
  mockAxios.reset();
  mockApi.restore();
  mockAxios.restore();
};
