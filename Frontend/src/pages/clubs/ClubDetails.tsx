import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../../components/layout/Layout';
import MyOpenClubFines from './MyOpenClubFines';

interface Club {
    id: string;
    name: string;
    description: string;
}

const ClubDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [club, setClub] = useState<Club | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [isAdmin, setIsAdmin] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [clubResponse, adminResponse] = await Promise.all([
                    api.get(`/api/v1/clubs/${id}`),
                    api.get(`/api/v1/clubs/${id}/isAdmin`)
                ]);
                setClub(clubResponse.data);
                setIsAdmin(adminResponse.data.isAdmin);
                setLoading(false);
            } catch {
                setError('Error fetching club details');
                setLoading(false);
            }
        };

        fetchData();
    }, [id]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!club) return <div>Club not found</div>;

    return (
        <Layout title={club.name}>
            <div>
                <h2>{club.name}</h2>
                <div className="club-info">
                    <p>{club.description}</p>
                    <MyOpenClubFines />
                    {isAdmin && (
                        <button 
                            className="button"
                            onClick={() => navigate(`/clubs/${id}/admin`)}
                            style={{ marginTop: '20px' }}
                        >
                            Manage Club
                        </button>
                    )}
                </div>
            </div>
        </Layout>
    );
};

export default ClubDetails;