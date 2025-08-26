import { lazy, Suspense } from 'react';
import type { RouteObject, Router as RouterType } from 'react-router-dom';
import { RouterProvider } from 'react-router-dom';
import ProtectedRoute from './components/auth/ProtectedRoute';
import { AuthProvider } from './context/AuthProvider';

// Lazy load page components for code splitting
const Dashboard = lazy(() => import('./pages/Dashboard'));
const ClubDetails = lazy(() => import('./pages/clubs/ClubDetails'));
const ClubList = lazy(() => import('./pages/clubs/ClubList'));
const AdminClubDetails = lazy(() => import('./pages/clubs/admin/AdminClubDetails'));
const CreateClub = lazy(() => import('./pages/clubs/CreateClub'));
const JoinClub = lazy(() => import('./pages/clubs/JoinClub'));
const Login = lazy(() => import('./pages/auth/Login'));
const MagicLinkHandler = lazy(() => import('./pages/auth/MagicLinkHandler'));
const KeycloakCallback = lazy(() => import('./pages/auth/KeycloakCallback'));
const Signup = lazy(() => import('./pages/auth/Signup'));
const Profile = lazy(() => import('./pages/profile/Profile'));
const ProfileInvites = lazy(() => import('./pages/profile/ProfileInvites'));
const ProfileFines = lazy(() => import('./pages/profile/ProfileFines'));
const ProfileSessions = lazy(() => import('./pages/profile/ProfileSessions'));
const ProfilePrivacy = lazy(() => import('./pages/profile/ProfilePrivacy'));
const ProfileNotificationSettings = lazy(() => import('./pages/profile/ProfileNotificationSettings'));
const EventDetails = lazy(() => import('./pages/clubs/events/EventDetails'));
const AdminEventDetails = lazy(() => import('./pages/clubs/admin/events/AdminEventDetails'));
const TeamDetails = lazy(() => import('./pages/teams/TeamDetails'));
const AdminTeamDetails = lazy(() => import('./pages/teams/AdminTeamDetails'));

// Loading component for suspense fallback
const PageLoader = () => (
    <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '200px',
        color: 'var(--color-text-secondary)'
    }}>
        Loading...
    </div>
);

export const routes: RouteObject[] = [
    { path: '/', element: (
        <ProtectedRoute>
            <Dashboard />
        </ProtectedRoute>
    ) },
    { path: '/clubs', element: (
        <ProtectedRoute>
            <ClubList />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:id', element: (
        <ProtectedRoute>
            <ClubDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:id/admin', element: (
        <ProtectedRoute>
            <AdminClubDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:id/admin/members', element: (
        <ProtectedRoute>
            <AdminClubDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:id/admin/teams', element: (
        <ProtectedRoute>
            <AdminClubDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:id/admin/fines', element: (
        <ProtectedRoute>
            <AdminClubDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:id/admin/events', element: (
        <ProtectedRoute>
            <AdminClubDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:id/admin/news', element: (
        <ProtectedRoute>
            <AdminClubDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:id/admin/settings', element: (
        <ProtectedRoute>
            <AdminClubDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:clubId/events/:eventId', element: (
        <ProtectedRoute>
            <EventDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:clubId/admin/events/:eventId', element: (
        <ProtectedRoute>
            <AdminEventDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:clubId/teams/:teamId', element: (
        <ProtectedRoute>
            <TeamDetails />
        </ProtectedRoute>
    ) },
    { path: '/clubs/:clubId/teams/:teamId/admin', element: (
        <ProtectedRoute>
            <AdminTeamDetails />
        </ProtectedRoute>
    ) },
    { path: '/createClub', element: (
        <ProtectedRoute>
            <CreateClub />
        </ProtectedRoute>
    ) },
    { path: '/profile', element: (
        <ProtectedRoute>
            <Profile />
        </ProtectedRoute>
    ) },
    { path: '/profile/privacy', element: (
        <ProtectedRoute>
            <ProfilePrivacy />
        </ProtectedRoute>
    ) },
    { path: '/profile/invites', element: (
        <ProtectedRoute>
            <ProfileInvites />
        </ProtectedRoute>
    ) },
    { path: '/profile/fines', element: (
        <ProtectedRoute>
            <ProfileFines />
        </ProtectedRoute>
    ) },
    { path: '/profile/sessions', element: (
        <ProtectedRoute>
            <ProfileSessions />
        </ProtectedRoute>
    ) },
    { path: '/profile/notifications', element: (
        <ProtectedRoute>
            <ProfileNotificationSettings />
        </ProtectedRoute>
    ) },
    { path: '/login', element: <Login /> },
    { path: '/auth/magic', element: <MagicLinkHandler /> },
    { path: '/auth/callback', element: <KeycloakCallback /> },
    { path: '/signup', element: (
        <ProtectedRoute>
            <Signup />
        </ProtectedRoute>
    ) },
    { path: '/join/:clubId', element: <JoinClub /> },
];

export function AppRouter({ router }: { router: RouterType }) {
    return (
        <AuthProvider>
            <Suspense fallback={<PageLoader />}>
                <RouterProvider router={router} />
            </Suspense>
        </AuthProvider>
    );
}

export default routes;
