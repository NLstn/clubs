import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../../components/layout/Layout';
import MyOpenClubFines from './MyOpenClubFines';
import UpcomingEvents from './UpcomingEvents';
import ClubNews from './ClubNews';
import MyTeams from './MyTeams';
import ClubNotFound from './ClubNotFound';
import { useClubSettings } from '../../hooks/useClubSettings';
import { addRecentClub, removeRecentClub } from '../../utils/recentClubs';
import { useT } from '../../hooks/useTranslation';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import './ClubDetails.css';

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
                const [clubResponse, adminResponse] = await Promise.all([
                    api.get(`/api/v1/clubs/${id}`),
                    api.get(`/api/v1/clubs/${id}/isAdmin`)
                ]);
                const clubData = clubResponse.data;
                setClub(clubData);
                setIsAdmin(adminResponse.data.isAdmin);

                // Get user's role by fetching club members and finding current user
                try {
                    const membersResponse = await api.get(`/api/v1/clubs/${id}/members`);
                    const currentUserMember = membersResponse.data.find((member: { userId: string; role: string }) => member.userId === currentUser?.ID);
                    if (currentUserMember) {
                        setUserRole(currentUserMember.role);
                    }
                } catch (memberErr) {
                    // If we can't fetch members, we might not be a member or might not have access
                    console.warn('Could not fetch member information:', memberErr);
                }
                
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
    }, [id, t, currentUser?.ID]);

    const handleLeaveClub = () => {
        setShowLeaveConfirmation(true);
    };

    const confirmLeaveClub = async () => {
        if (!id) return;

        setIsLeavingClub(true);
        try {
            await api.post(`/api/v1/clubs/${id}/leave`);
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
                <div className="club-header-section">
                    <div className="club-main-info">
                        <h1 className="club-title">{club.name}</h1>
                        {club.description && (
                            <p className="club-description">{club.description}</p>
                        )}
                        {userRole && (
                            <div className="user-role-container">
                                <span className="role-label">Your role</span>
                                <div className={`role-badge role-${userRole}`}>
                                    <span className="role-text">{userRole}</span>
                                </div>
                            </div>
                        )}
                    </div>
                    
                    {/* Action Buttons */}
                    <div className="club-actions">
                        {isAdmin && !club.deleted && (
                            <button 
                                className="button button-primary"
                                onClick={() => navigate(`/clubs/${id}/admin`)}
                            >
                                Manage Club
                            </button>
                        )}
                        {userRole && !club.deleted && (
                            <button 
                                className="button button-cancel"
                                onClick={handleLeaveClub}
                            >
                                Leave Club
                            </button>
                        )}
                    </div>
                </div>

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
                    {clubSettings?.teamsEnabled && <MyTeams />}
                    {clubSettings?.finesEnabled && <MyOpenClubFines />}
                </div>
            </div>
            
            {/* Leave Club Confirmation Modal */}
            {showLeaveConfirmation && (
                <div className="modal" onClick={cancelLeaveClub}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <h2>Leave Club</h2>
                        <p>Are you sure you want to leave "{club.name}"?</p>
                        <p>You will no longer have access to club content and will need to request to join again if you want to return.</p>
                        <div className="modal-actions">
                            <button 
                                onClick={confirmLeaveClub} 
                                disabled={isLeavingClub}
                                className="button-cancel"
                                style={{ backgroundColor: '#d32f2f', borderColor: '#d32f2f' }}
                            >
                                {isLeavingClub ? 'Leaving...' : 'Leave Club'}
                            </button>
                            <button 
                                onClick={cancelLeaveClub} 
                                disabled={isLeavingClub}
                                className="button-accept"
                            >
                                Cancel
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </Layout>
    );
};

export default ClubDetails;