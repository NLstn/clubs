import { useState, useEffect } from 'react';
import './Dashboard.css';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

interface Club {
    id: number;
    name: string;
    description: string;
}

const Dashboard = () => {
    const navigate = useNavigate();
    const { api } = useAuth();
    const [clubs, setClubs] = useState<Club[]>([]);
    const [message, setMessage] = useState('');

    useEffect(() => {
        fetchClubs();
    }, []);

    const fetchClubs = async () => {
        try {
            const response = await api.get('/api/v1/clubs');
            setClubs(response.data);
        } catch (error) {
            setMessage('Error fetching clubs');
            console.error(error);
        }
    };

    return (
        <div className="dashboard">
            <div className="dashboard-header">
                <h2>Dashboard</h2>
                <button className="create-button" onClick={() => navigate('/createClub')}>
                    Create New Club
                </button>
            </div>
            
            {message && <p className={`message ${message.includes('Error') ? 'error' : 'success'}`}>
                {message}
            </p>}

            <div className="clubs-list">
                {clubs === null || clubs.length === 0 ? (
                    <p>No clubs available. Create one to get started!</p>
                ) : (
                    clubs.map(club => (
                        <div 
                            key={club.id} 
                            className="club-item"
                            onClick={() => navigate(`/clubs/${club.id}`)}
                            style={{ cursor: 'pointer' }}
                        >
                            <h4>{club.name}</h4>
                            <p>{club.description}</p>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
};

export default Dashboard;