import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Layout from './layout/Layout';

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
        <Layout title="Dashboard">
            <div>
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
                                className="card"
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
        </Layout>
    );
};

export default Dashboard;