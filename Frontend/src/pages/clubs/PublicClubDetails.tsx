import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../../components/layout/Layout';
import PageHeader from '../../components/layout/PageHeader';
import ClubNotFound from './ClubNotFound';
import { Button, Card } from '@/components/ui';
import { useT } from '../../hooks/useTranslation';
import './ClubDetails.css';

interface PublicClubDetails {
    ID: string;
    Name: string;
    Description: string | null;
    LogoURL: string | null;
    IsMember: boolean;
}

const PublicClubDetails = () => {
    const { t } = useT();
    const { id } = useParams();
    const navigate = useNavigate();
    const [clubDetails, setClubDetails] = useState<PublicClubDetails | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [clubNotFound, setClubNotFound] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            if (!id) {
                setError('No club ID provided');
                setLoading(false);
                return;
            }

            try {
                // OData v2: Fetch public club details
                const response = await api.get(`/api/v2/Clubs('${id}')/GetPublicDetails()`);
                const data = response.data.value || response.data;
                
                setClubDetails(data);
                
                // If user is actually a member, redirect to the full club details page
                if (data.IsMember) {
                    navigate(`/clubs/${id}`, { replace: true });
                    return;
                }
                
                setLoading(false);
            } catch (err: unknown) {
                console.error('Error fetching public club details:', err);
                
                // Check if it's a 404 or 403 error (club not found or not discoverable)
                if (err && typeof err === 'object' && 'response' in err) {
                    const axiosError = err as { response?: { status?: number } };
                    if (axiosError.response?.status === 404 || axiosError.response?.status === 403) {
                        setClubNotFound(true);
                    } else {
                        setError(t('clubs.errors.loadingClub'));
                    }
                } else {
                    setError(t('clubs.errors.loadingClub'));
                }
                setLoading(false);
            }
        };

        fetchData();
    }, [id, t, navigate]);

    if (loading) return <div>Loading...</div>;
    if (clubNotFound) return <ClubNotFound clubId={id} />;
    if (error) return <div className="error">{error}</div>;
    if (!clubDetails) return <div>Club not found</div>;

    return (
        <Layout title={clubDetails.Name}>
            <div className="club-details-container">
                {/* Club Header */}
                <PageHeader>
                    {/* Club Logo */}
                    <div className="club-logo-section">
                        {clubDetails.LogoURL ? (
                            <img
                                src={clubDetails.LogoURL}
                                alt={`${clubDetails.Name} logo`}
                                className="club-logo"
                            />
                        ) : (
                            <div className="club-logo-placeholder">
                                <span className="logo-placeholder-text">
                                    {clubDetails.Name.charAt(0).toUpperCase()}
                                </span>
                            </div>
                        )}
                    </div>

                    <div className="club-main-info">
                        <h1 className="club-title">{clubDetails.Name}</h1>
                        {clubDetails.Description && (
                            <p className="club-description">{clubDetails.Description}</p>
                        )}
                    </div>
                </PageHeader>

                {/* Content - Limited view for non-members */}
                <div className="club-content">
                    <Card variant="light" padding="lg">
                        <h3>{t('clubs.interestedInJoining')}</h3>
                        <p>{t('clubs.publicPreview', { clubName: clubDetails.Name })}</p>
                        <Button 
                            variant="primary"
                            onClick={() => navigate(`/join/${id}`)}
                            style={{ marginTop: 'var(--space-md)' }}
                        >
                            {t('clubs.requestToJoin')}
                        </Button>
                        <Button 
                            variant="secondary"
                            onClick={() => navigate('/clubs')}
                            style={{ marginTop: 'var(--space-md)', marginLeft: 'var(--space-sm)' }}
                        >
                            {t('clubs.backToMyClubs')}
                        </Button>
                    </Card>
                </div>
            </div>
        </Layout>
    );
};

export default PublicClubDetails;
