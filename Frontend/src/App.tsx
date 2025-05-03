import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Dashboard from './components/Dashboard';
import ClubDetails from './components/clubs/ClubDetails';
import AdminClubDetails from './components/clubs/admin/AdminClubDetails';
import CreateClub from './components/clubs/CreateClub';
import Login from './components/auth/Login';
import MagicLinkHandler from './components/auth/MagicLinkHandler';
import ProtectedRoute from './components/auth/ProtectedRoute';
import { AuthProvider } from './context/AuthContext';
import Profile from './components/profile/Profile';
import ProfileInvites from './components/profile/ProfileInvites';

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
                    
                    <Route path="/profile/invites"  element={
                        <ProtectedRoute>
                            <ProfileInvites />
                        </ProtectedRoute>
                    } />

                    <Route path="/login" element={<Login />} />
                    <Route path="/auth/magic" element={<MagicLinkHandler />} />

                </Routes>
            </BrowserRouter>
        </AuthProvider>
    );
}

export default App;