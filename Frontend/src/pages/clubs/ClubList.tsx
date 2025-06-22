import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../../components/layout/Layout';
import api from '../../utils/api';
import './ClubList.css';

interface Club {
    id: string;
    name: string;
    description: string;
    user_role: string;
    created_at: string;
    deleted?: boolean;
}

const ClubList = () => {
    const [clubs, setClubs] = useState<Club[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        fetchClubs();
    }, []);

    const fetchClubs = async () => {
        try {
            const response = await api.get('/api/v1/clubs');
            setClubs(response.data);
        } catch (err: Error | unknown) {
            console.error('Error fetching clubs:', err);
            setError('Failed to fetch clubs');
        } finally {
            setLoading(false);
        }
    };

    const handleClubClick = (clubId: string) => {
        navigate(`/clubs/${clubId}`);
    };

    const adminClubs = clubs.filter(club => club.user_role === 'owner' || club.user_role === 'admin');
    const memberClubs = clubs.filter(club => club.user_role === 'member');

    if (loading) {
        return (
            <Layout title="My Clubs">
                <div>Loading clubs...</div>
            </Layout>
        );
    }

    if (error) {
        return (
            <Layout title="My Clubs">
                <div className="error">{error}</div>
            </Layout>
        );
    }

    return (
        <Layout title="My Clubs">
            <div className="clubs-container">
                {adminClubs.length > 0 && (
                    <div className="clubs-section">
                        <h2>Clubs I Manage</h2>
                        <div className="clubs-grid">
                            {adminClubs.map(club => (
                                <div 
                                    key={club.id} 
                                    className="club-card"
                                    onClick={() => handleClubClick(club.id)}
                                >
                                    <div className="club-header">
                                        <h3>{club.name}</h3>
                                        <span className={`role-badge ${club.user_role}`}>
                                            {club.user_role}
                                        </span>
                                    </div>
                                    <p className="club-description">{club.description}</p>
                                    {club.deleted && (
                                        <div className="club-deleted-badge">
                                            Deleted
                                        </div>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {memberClubs.length > 0 && (
                    <div className="clubs-section">
                        <h2>Clubs I'm a Member Of</h2>
                        <div className="clubs-grid">
                            {memberClubs.map(club => (
                                <div 
                                    key={club.id} 
                                    className="club-card"
                                    onClick={() => handleClubClick(club.id)}
                                >
                                    <div className="club-header">
                                        <h3>{club.name}</h3>
                                        <span className={`role-badge ${club.user_role}`}>
                                            {club.user_role}
                                        </span>
                                    </div>
                                    <p className="club-description">{club.description}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {clubs.length === 0 && (
                    <div className="empty-state">
                        <h2>No Clubs Yet</h2>
                        <p>You're not a member of any clubs yet.</p>
                        <button 
                            onClick={() => navigate('/createClub')}
                            className="button-primary"
                        >
                            Create Your First Club
                        </button>
                    </div>
                )}
            </div>
        </Layout>
    );
};

export default ClubList;