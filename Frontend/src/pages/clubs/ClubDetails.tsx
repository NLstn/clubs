import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../../components/layout/Layout';
import MyOpenClubFines from './MyOpenClubFines';
import UpcomingEvents from './UpcomingEvents';
import ClubNews from './ClubNews';
import ClubNotFound from './ClubNotFound';
import { useClubSettings } from '../../hooks/useClubSettings';
import { addRecentClub, removeRecentClub } from '../../utils/recentClubs';
import { useT } from '../../hooks/useTranslation';

interface Club {
    id: string;
    name: string;
    description: string;
    deleted?: boolean;
}

const ClubDetails = () => {
    const { t } = useT();
    const { id } = useParams();
    const navigate = useNavigate();
    const [club, setClub] = useState<Club | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [clubNotFound, setClubNotFound] = useState(false);
    const [isAdmin, setIsAdmin] = useState(false);
    const { settings: clubSettings } = useClubSettings(id);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [clubResponse, adminResponse] = await Promise.all([
                    api.get(`/api/v1/clubs/${id}`),
                    api.get(`/api/v1/clubs/${id}/isAdmin`)
                ]);
                const clubData = clubResponse.data;
                setClub(clubData);
                setIsAdmin(adminResponse.data.isAdmin);
                
                // Track this club visit
                if (clubData && clubData.id && clubData.name) {
                    addRecentClub(clubData.id, clubData.name);
                }
                
                setLoading(false);
            } catch (err: unknown) {
                console.error('Error fetching club details:', err);
                
                // Check if it's a 404 or 403 error (club not found or unauthorized)
                if (err && typeof err === 'object' && 'response' in err) {
                    const axiosError = err as { response?: { status?: number } };
                    if (axiosError.response?.status === 404 || axiosError.response?.status === 403) {
                        setClubNotFound(true);
                        // Remove this club from recent clubs since it doesn't exist or user can't access it
                        if (id) {
                            removeRecentClub(id);
                        }
                    } else {
                        setError(t('clubs.errors.loadingClub') || 'Error fetching club details');
                    }
                } else {
                    setError(t('clubs.errors.loadingClub') || 'Error fetching club details');
                }
                setLoading(false);
            }
        };

        if (id) {
            fetchData();
        } else {
            setError('No club ID provided');
            setLoading(false);
        }
    }, [id, t]);

    if (loading) return <div>Loading...</div>;
    if (clubNotFound) return <ClubNotFound clubId={id} />;
    if (error) return <div className="error">{error}</div>;
    if (!club) return <div>Club not found</div>;

    return (
        <Layout title={club.name}>
            <div>
                <h2>{club.name}</h2>
                {club.deleted && (
                    <div className="club-deleted-notice" style={{ 
                        backgroundColor: '#f44336', 
                        color: 'white', 
                        padding: '15px', 
                        marginBottom: '20px',
                        borderRadius: '4px',
                        fontWeight: 'bold'
                    }}>
                        {t('clubs.clubDeleted')}
                    </div>
                )}
                <div className="club-info">
                    <p>{club.description}</p>
                    <ClubNews />
                    <UpcomingEvents />
                    {clubSettings?.finesEnabled && <MyOpenClubFines />}
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