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

  // Notifications endpoints (v2 OData)
  mockApi.onGet(/\/api\/v2\/Notifications/).reply(200, { value: mockNotifications });
  mockApi.onGet(/\/api\/v2\/Notifications\/\$count/).reply(200, 1);
  mockApi.onPost(/\/api\/v2\/Notifications\('[\w-]+'\)\/MarkAsRead/).reply(200);
  mockApi.onDelete(/\/api\/v2\/Notifications\('[\w-]+'\)/).reply(200);
  mockApi.onPost('/api/v2/MarkAllNotificationsRead').reply(200);

  // Auth endpoints (keep on v1 - these are authentication concerns, not data operations)
  mockApi.onPost('/api/v1/auth/refreshToken').reply(200, {
    access: 'mock-access-token',
    refresh: 'mock-refresh-token'
  });
  mockAxios.onPost('http://localhost:8080/api/v1/auth/refreshToken').reply(200, {
    access: 'mock-access-token',
    refresh: 'mock-refresh-token'
  });
  mockApi.onPost('/api/v1/auth/logout').reply(200);

  // User endpoints (v2 OData)
  mockApi.onGet(/\/api\/v2\/Users/).reply(200, { value: mockUsers });
  mockApi.onPatch(/\/api\/v2\/Users\('[\w-]+'\)/).reply(200, mockUsers[0]);

  // Club endpoints (v2 OData)
  mockApi.onGet(/\/api\/v2\/Clubs$/).reply(200, { value: mockClubs });
  mockApi.onGet(/\/api\/v2\/Clubs\('[\w-]+'\)/).reply(200, mockClubs[0]);
  mockApi.onPost('/api/v2/Clubs').reply(201, mockClubs[0]);
  mockApi.onPatch(/\/api\/v2\/Clubs\('[\w-]+'\)/).reply(200, mockClubs[0]);
  mockApi.onDelete(/\/api\/v2\/Clubs\('[\w-]+'\)/).reply(200);

  // Dashboard endpoints (v2 OData functions)
  mockApi.onGet(/\/api\/v2\/GetDashboardActivities/).reply(200, mockActivities);

  // Members endpoints (v2 OData)
  mockApi.onGet(/\/api\/v2\/Members/).reply(200, { value: mockUsers });

  // Fines endpoints (v2 OData with filters)
  mockApi.onGet(/\/api\/v2\/Fines/).reply(200, { value: mockFines });
  mockApi.onPost('/api/v2/Fines').reply(201, mockFines[0]);
  mockApi.onPatch(/\/api\/v2\/Fines\('[\w-]+'\)/).reply(200, mockFines[0]);
  mockApi.onDelete(/\/api\/v2\/Fines\('[\w-]+'\)/).reply(200);

  // Teams endpoints (v2 OData)
  mockApi.onGet(/\/api\/v2\/Teams/).reply(200, { value: mockTeams });
  mockApi.onPost('/api/v2/Teams').reply(201, mockTeams[0]);
  mockApi.onPatch(/\/api\/v2\/Teams\('[\w-]+'\)/).reply(200, mockTeams[0]);
  mockApi.onDelete(/\/api\/v2\/Teams\('[\w-]+'\)/).reply(200);

  // Events endpoints (v2 OData)
  mockApi.onGet(/\/api\/v2\/Events/).reply(200, { value: [] });
  mockApi.onPost('/api/v2/Events').reply(201, {});
  mockApi.onDelete(/\/api\/v2\/Events\('[\w-]+'\)/).reply(200);
  
  // Join requests endpoints (v2 OData)
  mockApi.onGet(/\/api\/v2\/JoinRequests/).reply(200, { value: [] });
  
  // Invites endpoints (v2 OData)
  mockApi.onGet(/\/api\/v2\/Invites/).reply(200, { value: [] });

  // Search endpoint (v2 OData function)
  mockApi.onGet(/\/api\/v2\/SearchGlobal/).reply(200, { clubs: [], events: [] });

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
