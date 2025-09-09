import { render, screen, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import '@testing-library/jest-dom';
import App from '../App';

// Mock useAuth hook
vi.mock('../hooks/useAuth', () => ({
  useAuth: () => ({
    isAuthenticated: true,
    user: { id: '1', email: 'test@example.com' },
    api: {},
    logout: vi.fn(),
    login: vi.fn(),
    accessToken: 'mock-token',
    refreshToken: 'mock-refresh-token'
  })
}));

// Mock i18n
vi.mock('../i18n/useT', () => ({
  useT: () => ({
    t: (key: string) => key
  })
}));

// Mock react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
    i18n: {
      language: 'en',
      changeLanguage: vi.fn()
    }
  })
}));

// Mock the entire router module
vi.mock('react-router-dom', () => ({
  BrowserRouter: ({ children }: { children: React.ReactNode }) => <div data-testid="browser-router">{children}</div>,
  Routes: ({ children }: { children: React.ReactNode }) => <div data-testid="routes">{children}</div>,
  Route: ({ element }: { element: React.ReactNode }) => <div data-testid="route">{element}</div>,
  useNavigate: () => vi.fn(),
  useParams: () => ({ clubId: '1', eventId: '1' }),
  useLocation: () => ({ pathname: '/', search: '', hash: '', state: null }),
  Navigate: ({ to }: { to: string }) => <div data-testid="navigate">Navigate to: {to}</div>,
  Outlet: () => <div data-testid="outlet">Outlet</div>,
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => <a href={to}>{children}</a>,
}));

// Mock AuthProvider
vi.mock('../context/AuthProvider', () => ({
  AuthProvider: ({ children }: { children: React.ReactNode }) => <div data-testid="auth-provider">{children}</div>,
}));

// Mock all page components
vi.mock('../pages/Dashboard', () => ({
  default: () => <div data-testid="dashboard">Dashboard</div>,
}));

vi.mock('../pages/clubs/ClubDetails', () => ({
  default: () => <div data-testid="club-details">ClubDetails</div>,
}));

vi.mock('../pages/clubs/ClubList', () => ({
  default: () => <div data-testid="club-list">ClubList</div>,
}));

vi.mock('../pages/clubs/admin/AdminClubDetails', () => ({
  default: () => <div data-testid="admin-club-details">AdminClubDetails</div>,
}));

vi.mock('../pages/clubs/CreateClub', () => ({
  default: () => <div data-testid="create-club">CreateClub</div>,
}));

vi.mock('../pages/clubs/JoinClub', () => ({
  default: () => <div data-testid="join-club">JoinClub</div>,
}));

vi.mock('../pages/auth/Login', () => ({
  default: () => <div data-testid="login">Login</div>,
}));

vi.mock('../pages/auth/MagicLinkHandler', () => ({
  default: () => <div data-testid="magic-link-handler">MagicLinkHandler</div>,
}));

vi.mock('../pages/auth/KeycloakCallback', () => ({
  default: () => <div data-testid="keycloak-callback">KeycloakCallback</div>,
}));

vi.mock('../pages/auth/Signup', () => ({
  default: () => <div data-testid="signup">Signup</div>,
}));

vi.mock('../components/auth/ProtectedRoute', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div data-testid="protected-route">{children}</div>,
}));

// Mock AuthProvider
vi.mock('../context/AuthProvider', () => ({
  AuthProvider: ({ children }: { children: React.ReactNode }) => <div data-testid="auth-provider">{children}</div>,
}));

vi.mock('../pages/profile/Profile', () => ({
  default: () => <div data-testid="profile">Profile</div>,
}));

vi.mock('../pages/profile/ProfileInvites', () => ({
  default: () => <div data-testid="profile-invites">ProfileInvites</div>,
}));

vi.mock('../pages/profile/ProfileFines', () => ({
  default: () => <div data-testid="profile-fines">ProfileFines</div>,
}));

vi.mock('../pages/profile/ProfileSessions', () => ({
  default: () => <div data-testid="profile-sessions">ProfileSessions</div>,
}));

vi.mock('../pages/profile/ProfilePrivacy', () => ({
  default: () => <div data-testid="profile-privacy">ProfilePrivacy</div>,
}));

vi.mock('../pages/profile/ProfileNotificationSettings', () => ({
  default: () => <div data-testid="profile-notification-settings">ProfileNotificationSettings</div>,
}));

vi.mock('../pages/profile/ProfileShifts', () => ({
  default: () => <div data-testid="profile-shifts">ProfileShifts</div>,
}));

vi.mock('../pages/clubs/events/EventDetails', () => ({
  default: () => <div data-testid="event-details">EventDetails</div>,
}));

vi.mock('../pages/clubs/admin/events/AdminEventDetails', () => ({
  default: () => <div data-testid="admin-event-details">AdminEventDetails</div>,
}));

vi.mock('../pages/teams/TeamDetails', () => ({
  default: () => <div data-testid="team-details">TeamDetails</div>,
}));

vi.mock('../pages/teams/AdminTeamDetails', () => ({
  default: () => <div data-testid="admin-team-details">AdminTeamDetails</div>,
}));

describe('App', () => {
  it('renders without crashing', async () => {
    await act(async () => {
      render(<App />);
    });
    
    expect(screen.getByTestId('auth-provider')).toBeInTheDocument();
    expect(screen.getByTestId('browser-router')).toBeInTheDocument();
    expect(screen.getByTestId('routes')).toBeInTheDocument();
  });

  it('renders all route components', async () => {
    await act(async () => {
      render(<App />);
    });
    
    // Wait for all components to load and render
    await waitFor(() => {
      expect(screen.getByTestId('dashboard')).toBeInTheDocument();
    });
    
    // Check that all page components are rendered within their routes
    expect(screen.getByTestId('club-list')).toBeInTheDocument();
    expect(screen.getByTestId('club-details')).toBeInTheDocument();
    expect(screen.getAllByTestId('admin-club-details')).toHaveLength(7); // Now we have 7 admin routes
    expect(screen.getByTestId('create-club')).toBeInTheDocument();
    expect(screen.getByTestId('profile')).toBeInTheDocument();
    expect(screen.getByTestId('profile-invites')).toBeInTheDocument();
    expect(screen.getByTestId('profile-fines')).toBeInTheDocument();
    expect(screen.getByTestId('profile-sessions')).toBeInTheDocument();
    expect(screen.getByTestId('profile-privacy')).toBeInTheDocument();
    expect(screen.getByTestId('profile-notification-settings')).toBeInTheDocument();
    expect(screen.getByTestId('profile-shifts')).toBeInTheDocument();
    expect(screen.getByTestId('login')).toBeInTheDocument();
    expect(screen.getByTestId('magic-link-handler')).toBeInTheDocument();
    expect(screen.getByTestId('keycloak-callback')).toBeInTheDocument();
    expect(screen.getByTestId('signup')).toBeInTheDocument();
    expect(screen.getByTestId('join-club')).toBeInTheDocument();
    expect(screen.getByTestId('event-details')).toBeInTheDocument();
    expect(screen.getByTestId('admin-event-details')).toBeInTheDocument();
  });

  it('wraps protected routes with ProtectedRoute component', async () => {
    await act(async () => {
      render(<App />);
    });
    
    // Wait for components to load
    await waitFor(() => {
      expect(screen.getByTestId('dashboard')).toBeInTheDocument();
    });
    
    // Most routes should be protected
    const protectedRoutes = screen.getAllByTestId('protected-route');
    expect(protectedRoutes.length).toBeGreaterThan(0);
  });

  it('wraps the app with AuthProvider', async () => {
    await act(async () => {
      render(<App />);
    });
    expect(screen.getByTestId('auth-provider')).toBeInTheDocument();
  });

  it('uses BrowserRouter for routing', async () => {
    await act(async () => {
      render(<App />);
    });
    expect(screen.getByTestId('browser-router')).toBeInTheDocument();
  });
});
