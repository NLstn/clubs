import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AdminDashboard from './components/admin/AdminDashboard';

const App = () => (
    <Router>
        <Routes>
            <Route path="/admin" element={<AdminDashboard />} />
        </Routes>
    </Router>
);

export default App;