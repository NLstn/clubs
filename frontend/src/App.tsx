import { BrowserRouter, Routes, Route } from 'react-router-dom';
import AdminDashboard from './components/admin/AdminDashboard';
import ClubDetails from './components/admin/ClubDetails';

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/admin" element={<AdminDashboard />} />
                <Route path="/admin/clubs/:id" element={<ClubDetails />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;