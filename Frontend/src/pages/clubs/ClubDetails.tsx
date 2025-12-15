import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../../components/layout/Layout';
import PageHeader from '../../components/layout/PageHeader';
import MyOpenClubFines from './MyOpenClubFines';
import UpcomingEvents from './UpcomingEvents';
import ClubNews from './ClubNews';
import MyTeams from './MyTeams';
import ReadonlyMemberList from './ReadonlyMemberList';
import ClubNotFound from './ClubNotFound';
import { useClubSettings } from '../../hooks/useClubSettings';
import { Modal, Button } from '@/components/ui';
import { addRecentClub, removeRecentClub } from '../../utils/recentClubs';
import { useT } from '../../hooks/useTranslation';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import './ClubDetails.css';

interface Club {
    id: string;
    name: string;
    description: string;
    logo_url?: string;
    deleted?: boolean;
}

const ClubDetails = () => {
    const { t } = useT();
    const { id } = useParams();
    const navigate = useNavigate();
    const { user: currentUser } = useCurrentUser();
    const [club, setClub] = useState<Club | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [clubNotFound, setClubNotFound] = useState(false);
    const [isAdmin, setIsAdmin] = useState(false);
    const [userRole, setUserRole] = useState<string | null>(null);
    const [showLeaveConfirmation, setShowLeaveConfirmation] = useState(false);
    const [isLeavingClub, setIsLeavingClub] = useState(false);
    const { settings: clubSettings } = useClubSettings(id);

    useEffect(() => {
        const fetchData = async () => {
            try {
                // OData v2: Fetch club and check admin status
                const [clubResponse, adminResponse] = await Promise.all([
                    api.get(`/api/v2/Clubs('${id}')`),
                    api.get(`/api/v2/Clubs('${id}')/IsAdmin()`)
                ]);
                const clubData = clubResponse.data;
                setClub({
                    id: clubData.ID,
                    name: clubData.Name,
                    description: clubData.Description,
                    logo_url: clubData.LogoURL,
                    deleted: clubData.Deleted
                });
                
                // Handle OData function response - might be wrapped in 'value' property
                const adminData = adminResponse.data.value || adminResponse.data;
                setIsAdmin(adminData.IsAdmin);

                // Get user's role by fetching club members and finding current user
                try {
                    interface ODataMember { UserID: string; Role: string; }
                    const membersResponse = await api.get(`/api/v2/Members?$filter=ClubID eq '${id}'&$expand=User`);
                    const members = membersResponse.data.value || [];
                    const currentUserMember = members.find((member: ODataMember) => member.UserID === currentUser?.ID);
                    if (currentUserMember) {
                        setUserRole(currentUserMember.Role);
                    }
                } catch (memberErr) {
                    // If we can't fetch members, we might not be a member or might not have access
                    console.warn('Could not fetch member information:', memberErr);
                }
                
                // Track this club visit
                if (clubData && clubData.ID && clubData.Name) {
                    addRecentClub(clubData.ID, clubData.Name);
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
                        setError(t('clubs.errors.loadingClub'));
                    }
                } else {
                    setError(t('clubs.errors.loadingClub'));
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
    }, [id, t, currentUser?.ID]);

    const handleLeaveClub = () => {
        setShowLeaveConfirmation(true);
    };

    const confirmLeaveClub = async () => {
        if (!id) return;

        setIsLeavingClub(true);
        try {
            // OData v2: Use Leave action on Club entity
            await api.post(`/api/v2/Clubs('${id}')/Leave`);
            // Remove from recent clubs since user is no longer a member
            removeRecentClub(id);
            // Navigate back to clubs list
            navigate('/clubs');
        } catch (error) {
            console.error('Error leaving club:', error);
            if (error && typeof error === 'object' && 'response' in error) {
                const axiosError = error as { response?: { data?: string } };
                const errorMessage = axiosError.response?.data || 'Failed to leave club';
                alert(errorMessage);
            } else {
                alert('Failed to leave club. Please try again.');
            }
        } finally {
            setIsLeavingClub(false);
            setShowLeaveConfirmation(false);
        }
    };

    const cancelLeaveClub = () => {
        setShowLeaveConfirmation(false);
    };

    if (loading) return <div>Loading...</div>;
    if (clubNotFound) return <ClubNotFound clubId={id} />;
    if (error) return <div className="error">{error}</div>;
    if (!club) return <div>Club not found</div>;

    return (
        <Layout title={club.name}>
            <div className="club-details-container">
                {/* Club Header */}
                <PageHeader
                    actions={
                        <>
                            {isAdmin && !club.deleted && (
                                <Button 
                                    variant="primary"
                                    onClick={() => navigate(`/clubs/${id}/admin`)}
                                >
                                    Manage Club
                                </Button>
                            )}
                            {userRole && !club.deleted && (
                                <Button 
                                    variant="cancel"
                                    onClick={handleLeaveClub}
                                >
                                    Leave Club
                                </Button>
                            )}
                        </>
                    }
                >
                    {/* Club Logo */}
                    <div className="club-logo-section">
                        {club.logo_url ? (
                            <img
                                src={club.logo_url}
                                alt={`${club.name} logo`}
                                className="club-logo"
                            />
                        ) : (
                            <div className="club-logo-placeholder">
                                <span className="logo-placeholder-text">
                                    {club.name.charAt(0).toUpperCase()}
                                </span>
                            </div>
                        )}
                    </div>

                    <div className="club-main-info">
                        <h1 className="club-title">{club.name}</h1>
                        {club.description && (
                            <p className="club-description">{club.description}</p>
                        )}
                    </div>
                </PageHeader>

                {/* Deleted Club Notice */}
                {club.deleted && (
                    <div className="club-deleted-notice">
                        <div className="notice-icon">⚠️</div>
                        <div className="notice-content">
                            <h3>Club Deleted</h3>
                            <p>{t('clubs.clubDeleted')}</p>
                        </div>
                    </div>
                )}

                {/* Content Sections */}
                <div className="club-content">
                    <ClubNews />
                    <UpcomingEvents />
                    {clubSettings?.TeamsEnabled && <MyTeams />}
                    {clubSettings?.FinesEnabled && <MyOpenClubFines />}
                    {clubSettings?.MembersListVisible && <ReadonlyMemberList />}
                </div>
            </div>
            
            {/* Leave Club Confirmation Modal */}
            <Modal 
                isOpen={showLeaveConfirmation} 
                onClose={cancelLeaveClub} 
                title="Leave Club"
                maxWidth="450px"
            >
                <Modal.Body>
                    <p>Are you sure you want to leave "{club.name}"?</p>
                    <p>You will no longer have access to club content and will need to request to join again if you want to return.</p>
                </Modal.Body>
                <Modal.Actions>
                    <Button 
                        variant="cancel"
                        onClick={confirmLeaveClub}
                        disabled={isLeavingClub}
                        style={{ backgroundColor: '#d32f2f', borderColor: '#d32f2f' }}
                    >
                        {isLeavingClub ? (
                            <>
                                <Modal.LoadingSpinner />
                                Leaving...
                            </>
                        ) : (
                            'Leave Club'
                        )}
                    </Button>
                    <Button 
                        variant="accept"
                        onClick={cancelLeaveClub}
                        disabled={isLeavingClub}
                    >
                        Cancel
                    </Button>
                </Modal.Actions>
            </Modal>
        </Layout>
    );
};

export default ClubDetails;