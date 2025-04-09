import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Dashboard from './components/Dashboard';
import ClubDetails from './components/ClubDetails';
import Login from './components/auth/Login';
import MagicLinkHandler from './components/auth/MagicLinkHandler';
import ProtectedRoute from './components/auth/ProtectedRoute';
import { AuthProvider } from './context/AuthContext';

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

                    <Route path="/login" element={<Login />} />
                    <Route path="/auth/magic" element={<MagicLinkHandler />} />

                </Routes>
            </BrowserRouter>
        </AuthProvider>
    );
}

export default App;