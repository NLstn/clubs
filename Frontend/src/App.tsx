import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import AdminDashboard from './components/admin/AdminDashboard';
import ClubDetails from './components/admin/ClubDetails';
import Login from './components/auth/Login';
import MagicLinkHandler from './components/auth/MagicLinkHandler';
import ProtectedRoute from './components/auth/ProtectedRoute';
import { AuthProvider } from './context/AuthContext';

function App() {
    return (
        <AuthProvider>
            <BrowserRouter>
                <Routes>
                    {/* Root redirect */}
                    <Route path="/" element={<Navigate to="/admin" replace />} />

                    {/* Auth routes */}
                    <Route path="/login" element={<Login />} />
                    <Route path="/auth/magic" element={<MagicLinkHandler />} />

                    {/* Protected routes */}
                    <Route 
                        path="/admin" 
                        element={
                            <ProtectedRoute>
                                <AdminDashboard />
                            </ProtectedRoute>
                        } 
                    />
                    <Route 
                        path="/admin/clubs/:id" 
                        element={
                            <ProtectedRoute>
                                <ClubDetails />
                            </ProtectedRoute>
                        } 
                    />

                </Routes>
            </BrowserRouter>
        </AuthProvider>
    );
}

export default App;