import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import '@testing-library/jest-dom';
import App from '../App';

// Mock the entire router module
vi.mock('react-router-dom', () => ({
  BrowserRouter: ({ children }: { children: React.ReactNode }) => <div data-testid="browser-router">{children}</div>,
  Routes: ({ children }: { children: React.ReactNode }) => <div data-testid="routes">{children}</div>,
  Route: ({ element }: { element: React.ReactNode }) => <div data-testid="route">{element}</div>,
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

vi.mock('../pages/settings/NotificationSettings', () => ({
  default: () => <div data-testid="notification-settings">NotificationSettings</div>,
}));

describe('App', () => {
  it('renders without crashing', () => {
    render(<App />);
    expect(screen.getByTestId('auth-provider')).toBeInTheDocument();
    expect(screen.getByTestId('browser-router')).toBeInTheDocument();
    expect(screen.getByTestId('routes')).toBeInTheDocument();
  });

  it('renders all route components', () => {
    render(<App />);
    
    // Check that all page components are rendered within their routes
    expect(screen.getByTestId('dashboard')).toBeInTheDocument();
    expect(screen.getByTestId('club-list')).toBeInTheDocument();
    expect(screen.getByTestId('club-details')).toBeInTheDocument();
    expect(screen.getByTestId('admin-club-details')).toBeInTheDocument();
    expect(screen.getByTestId('create-club')).toBeInTheDocument();
    expect(screen.getByTestId('profile')).toBeInTheDocument();
    expect(screen.getByTestId('profile-invites')).toBeInTheDocument();
    expect(screen.getByTestId('profile-fines')).toBeInTheDocument();
    expect(screen.getByTestId('profile-sessions')).toBeInTheDocument();
    expect(screen.getByTestId('notification-settings')).toBeInTheDocument();
    expect(screen.getByTestId('login')).toBeInTheDocument();
    expect(screen.getByTestId('magic-link-handler')).toBeInTheDocument();
    expect(screen.getByTestId('keycloak-callback')).toBeInTheDocument();
    expect(screen.getByTestId('signup')).toBeInTheDocument();
    expect(screen.getByTestId('join-club')).toBeInTheDocument();
  });

  it('wraps protected routes with ProtectedRoute component', () => {
    render(<App />);
    
    // Most routes should be protected
    const protectedRoutes = screen.getAllByTestId('protected-route');
    expect(protectedRoutes.length).toBeGreaterThan(0);
  });

  it('wraps the app with AuthProvider', () => {
    render(<App />);
    expect(screen.getByTestId('auth-provider')).toBeInTheDocument();
  });

  it('uses BrowserRouter for routing', () => {
    render(<App />);
    expect(screen.getByTestId('browser-router')).toBeInTheDocument();
  });
});
