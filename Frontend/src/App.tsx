import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import ClubDetails from './pages/clubs/ClubDetails';
import ClubList from './pages/clubs/ClubList';
import AdminClubDetails from './pages/clubs/admin/AdminClubDetails';
import CreateClub from './pages/clubs/CreateClub';
import JoinClub from './pages/clubs/JoinClub';
import Login from './pages/auth/Login';
import MagicLinkHandler from './pages/auth/MagicLinkHandler';
import KeycloakCallback from './pages/auth/KeycloakCallback';
import Signup from './pages/auth/Signup';
import ProtectedRoute from './components/auth/ProtectedRoute';
import { AuthProvider } from './context/AuthProvider';
import Profile from './pages/profile/Profile';
import ProfileInvites from './pages/profile/ProfileInvites';
import ProfileFines from './pages/profile/ProfileFines';
import ProfileSessions from './pages/profile/ProfileSessions';
import ProfilePrivacy from './pages/profile/ProfilePrivacy';
import NotificationSettings from './pages/settings/NotificationSettings';

function App() {
    return (
        <AuthProvider>
            <BrowserRouter>
                <Routes>
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
                        path="/createClub"
                        element={
                            <ProtectedRoute>
                                <CreateClub />
                            </ProtectedRoute>
                        }
                    />

                    <Route path="/profile"  element={
                        <ProtectedRoute>
                            <Profile />
                        </ProtectedRoute>
                    } />
                    
                    <Route path="/profile/privacy"  element={
                        <ProtectedRoute>
                            <ProfilePrivacy />
                        </ProtectedRoute>
                    } />
                    
                    <Route path="/profile/invites"  element={
                        <ProtectedRoute>
                            <ProfileInvites />
                        </ProtectedRoute>
                    } />
                    <Route path="/profile/fines"  element={
                        <ProtectedRoute>
                            <ProfileFines />
                        </ProtectedRoute>
                    } />
                    <Route path="/profile/sessions"  element={
                        <ProtectedRoute>
                            <ProfileSessions />
                        </ProtectedRoute>
                    } />
                    <Route path="/profile/notifications"  element={
                        <ProtectedRoute>
                            <NotificationSettings />
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
            </BrowserRouter>
        </AuthProvider>
    );
}

export default App;