import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { lazy, Suspense } from 'react';
import ProtectedRoute from './components/auth/ProtectedRoute';
import { AuthProvider } from './context/AuthProvider';
import { ThemeProvider } from './context/ThemeProvider';

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
const ProfilePreferences = lazy(() => import('./pages/profile/ProfilePreferences'));
const ProfileInvites = lazy(() => import('./pages/profile/ProfileInvites'));
const ProfileFines = lazy(() => import('./pages/profile/ProfileFines'));
const ProfileShifts = lazy(() => import('./pages/profile/ProfileShifts'));
const ProfileSessions = lazy(() => import('./pages/profile/ProfileSessions'));
const ProfilePrivacy = lazy(() => import('./pages/profile/ProfilePrivacy'));
const ProfileAPIKeys = lazy(() => import('./pages/profile/ProfileAPIKeys'));
const ProfileNotificationSettings = lazy(() => import('./pages/profile/ProfileNotificationSettings'));
const EventDetails = lazy(() => import('./pages/clubs/events/EventDetails'));
const AdminEventDetails = lazy(() => import('./pages/clubs/admin/events/AdminEventDetails'));
const TeamDetails = lazy(() => import('./pages/teams/TeamDetails'));
const AdminTeamDetails = lazy(() => import('./pages/teams/AdminTeamDetails'));
const ButtonDemo = lazy(() => import('./pages/ButtonDemo'));

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

function App() {
    return (
        <ThemeProvider>
            <AuthProvider>
                <BrowserRouter>
                    <Suspense fallback={<PageLoader />}>
                        <Routes data-testid="routes">
                        <Route path="/" element={
                            <ProtectedRoute>
                                <Dashboard />
                            </ProtectedRoute>
                        } />

                        <Route path="/clubs" element={
                            <ProtectedRoute>
                                <ClubList />
                            </ProtectedRoute>
                        } />

                        <Route
                            path="/clubs/:id"
                            element={
                                <ProtectedRoute>
                                    <ClubDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:id/admin"
                            element={
                                <ProtectedRoute>
                                    <AdminClubDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:id/admin/members"
                            element={
                                <ProtectedRoute>
                                    <AdminClubDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:id/admin/teams"
                            element={
                                <ProtectedRoute>
                                    <AdminClubDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:id/admin/fines"
                            element={
                                <ProtectedRoute>
                                    <AdminClubDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:id/admin/events"
                            element={
                                <ProtectedRoute>
                                    <AdminClubDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:id/admin/news"
                            element={
                                <ProtectedRoute>
                                    <AdminClubDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:id/admin/settings"
                            element={
                                <ProtectedRoute>
                                    <AdminClubDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:clubId/events/:eventId"
                            element={
                                <ProtectedRoute>
                                    <EventDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:clubId/admin/events/:eventId"
                            element={
                                <ProtectedRoute>
                                    <AdminEventDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:clubId/teams/:teamId"
                            element={
                                <ProtectedRoute>
                                    <TeamDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/clubs/:clubId/teams/:teamId/admin"
                            element={
                                <ProtectedRoute>
                                    <AdminTeamDetails />
                                </ProtectedRoute>
                            }
                        />

                        <Route
                            path="/createClub"
                            element={
                                <ProtectedRoute>
                                    <CreateClub />
                                </ProtectedRoute>
                            }
                        />

                        <Route path="/profile" element={
                            <ProtectedRoute>
                                <Profile />
                            </ProtectedRoute>
                        } />

                        <Route path="/profile/preferences" element={
                            <ProtectedRoute>
                                <ProfilePreferences />
                            </ProtectedRoute>
                        } />

                        <Route path="/profile/privacy" element={
                            <ProtectedRoute>
                                <ProfilePrivacy />
                            </ProtectedRoute>
                        } />

                        <Route path="/profile/invites" element={
                            <ProtectedRoute>
                                <ProfileInvites />
                            </ProtectedRoute>
                        } />
                        <Route path="/profile/fines" element={
                            <ProtectedRoute>
                                <ProfileFines />
                            </ProtectedRoute>
                        } />
                        <Route path="/profile/shifts" element={
                            <ProtectedRoute>
                                <ProfileShifts />
                            </ProtectedRoute>
                        } />
                        <Route path="/profile/sessions" element={
                            <ProtectedRoute>
                                <ProfileSessions />
                            </ProtectedRoute>
                        } />
                        <Route path="/profile/api-keys" element={
                            <ProtectedRoute>
                                <ProfileAPIKeys />
                            </ProtectedRoute>
                        } />
                        <Route path="/profile/notifications" element={
                            <ProtectedRoute>
                                <ProfileNotificationSettings />
                            </ProtectedRoute>
                        } />

                        <Route path="/demo/button" element={
                            <ProtectedRoute>
                                <ButtonDemo />
                            </ProtectedRoute>
                        } />

                        <Route path="/login" element={<Login />} />
                        <Route path="/auth/magic" element={<MagicLinkHandler />} />
                        <Route path="/auth/callback" element={<KeycloakCallback />} />
                        <Route path="/signup" element={
                            <ProtectedRoute>
                                <Signup />
                            </ProtectedRoute>
                        } />
                        <Route path="/join/:clubId" element={<JoinClub />} />

                    </Routes>
                </Suspense>
            </BrowserRouter>
        </AuthProvider>
        </ThemeProvider>
    );
}

export default App;