import { useState, useEffect, useCallback, useRef } from 'react';
import { useParams, useNavigate, useSearchParams, useLocation, Link } from 'react-router-dom';
import api, { hardDeleteClub } from '../../../utils/api';
import { parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';
import Layout from '../../../components/layout/Layout';
import PageHeader from '../../../components/layout/PageHeader';
import ClubNotFound from '../ClubNotFound';
import AdminClubMemberList from './members/AdminClubMemberList';
import AdminClubTeamList from './teams/AdminClubTeamList';
import AdminClubFineList from './fines/AdminClubFineList';
import AdminClubEventList from './events/AdminClubEventList';
import AdminClubNewsList from './news/AdminClubNewsList';
import { Input, Modal, Button } from '@/components/ui';
import AdminClubSettings from './settings/AdminClubSettings';
import { useClubSettings } from '../../../hooks/useClubSettings';
import { useT } from '../../../hooks/useTranslation';
import { removeRecentClub } from '../../../utils/recentClubs';
import StatisticsCard from '../../../components/dashboard/StatisticsCard';
import '@/components/ui/Tabs.css';
import './AdminClubDetails.css';

interface Club {
    ID: string;
    Name: string;
    Description?: string;
    LogoURL?: string;
    Deleted?: boolean;
}

interface MemberStats {
    total: number;
    newThisMonth: number;
    pendingInvites: number;
}

const AdminClubDetails = () => {
    const { t } = useT();
    const { id } = useParams();
    const navigate = useNavigate();
    const location = useLocation();
    const [searchParams] = useSearchParams();
    const [club, setClub] = useState<Club | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const { settings: clubSettings, refetch: refetchSettings } = useClubSettings(id);
    
    // Determine current tab from route
    const getCurrentTab = () => {
        const path = location.pathname;
        if (path.includes('/admin/members')) return 'members';
        if (path.includes('/admin/teams')) return 'teams';
        if (path.includes('/admin/fines')) return 'fines';
        if (path.includes('/admin/events')) return 'events';
        if (path.includes('/admin/news')) return 'news';
        if (path.includes('/admin/settings')) return 'settings';
        return 'overview';
    };
    
    const activeTab = getCurrentTab();
    
    
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [clubNotFound, setClubNotFound] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [editForm, setEditForm] = useState({ name: '', description: '' });
    const [isOwner, setIsOwner] = useState(false);
    const [showDeletePopup, setShowDeletePopup] = useState(false);
    const [showHardDeletePopup, setShowHardDeletePopup] = useState(false);
    const [logoUploading, setLogoUploading] = useState(false);
    const [logoError, setLogoError] = useState<string | null>(null);
    const [memberStats, setMemberStats] = useState<MemberStats>({ total: 0, newThisMonth: 0, pendingInvites: 0 });
    const [statsLoading, setStatsLoading] = useState(true);

    const fetchMemberStats = useCallback(async () => {
        if (!id) return;
        
        setStatsLoading(true);
        try {
            // OData v2: Fetch Members and Invites using OData queries
            interface ODataMember { CreatedAt: string; }
            interface ODataInvite { ID: string; }
            const [membersResponse, invitesResponse] = await Promise.all([
                api.get<ODataCollectionResponse<ODataMember>>(`/api/v2/Members?$select=CreatedAt&$filter=ClubID eq '${id}'`),
                api.get<ODataCollectionResponse<ODataInvite>>(`/api/v2/Invites?$select=ID&$filter=ClubID eq '${id}'`),
            ]);

            const members = parseODataCollection(membersResponse.data);
            const invites = parseODataCollection(invitesResponse.data);

            // Calculate members joined in the last 30 days
            const thirtyDaysAgo = new Date();
            thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

            const newMembers = members.filter((member: { CreatedAt: string }) => {
                const joinedDate = new Date(member.CreatedAt);
                return joinedDate >= thirtyDaysAgo;
            });

            setMemberStats({
                total: members.length,
                newThisMonth: newMembers.length,
                pendingInvites: invites.length,
            });
        } catch (err) {
            console.error('Error fetching member stats:', err);
        } finally {
            setStatsLoading(false);
        }
    }, [id]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                // OData v2: Use IsAdmin function and fetch club data
                const [adminResponse, clubResponse] = await Promise.all([
                    api.get(`/api/v2/Clubs('${id}')/IsAdmin()`),
                    api.get(`/api/v2/Clubs('${id}')`),
                    
                ]);

                const adminData = adminResponse.data.value || adminResponse.data;
                if (!adminData.IsAdmin) {
                    navigate(`/clubs/${id}`);
                    return;
                }

                setClub(clubResponse.data);
                setIsOwner(adminResponse.data.IsOwner || false);
                
                // Fetch member statistics for overview tab
                if (activeTab === 'overview') {
                    fetchMemberStats();
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
    }, [id, navigate, t, activeTab, fetchMemberStats]);

    // Redirect to valid tab if current tab becomes unavailable
    useEffect(() => {
        if (clubSettings && id) {
            if (activeTab === 'fines' && !clubSettings.FinesEnabled) {
                navigate(`/clubs/${id}/admin`);
            }
            if (activeTab === 'teams' && !clubSettings.TeamsEnabled) {
                navigate(`/clubs/${id}/admin`);
            }
        }
    }, [clubSettings, activeTab, navigate, id]);

    const updateClub = async () => {
        try {
            // OData v2: Update club using PATCH
            const response = await api.patch(`/api/v2/Clubs('${id}')`, {
                Name: editForm.name,
                Description: editForm.description
            });
            const updatedClub = response.data;
            setClub(updatedClub);
            
            setIsEditing(false);
            setError(null);
        } catch {
            setError('Failed to update club');
        }
    };

    const handleEdit = () => {
        if (club) {
            setEditForm({ name: club.Name, description: club.Description ?? '' });
            setIsEditing(true);
        }
    };

    const handleDeleteClub = async () => {
        if (!club) return;
        
        setShowDeletePopup(true);
    };

    const confirmDeleteClub = async () => {
        if (!club) return;

        try {
            // OData v2: Mark club as deleted using DELETE (or custom HardDelete action if needed)
            await api.delete(`/api/v2/Clubs('${id}')`);
            // Note: We don't remove from recent clubs since owners can still access deleted clubs
            // Navigate to clubs list after deletion
            navigate('/');
        } catch (err: Error | unknown) {
            console.error('Error deleting club:', err instanceof Error ? err.message : 'Unknown error');
            setError('Failed to delete club');
        }
        setShowDeletePopup(false);
    };

    const handleHardDeleteClub = async () => {
        if (!club) return;
        
        setShowHardDeletePopup(true);
    };

    const confirmHardDeleteClub = async () => {
        if (!club) return;

        try {
            await hardDeleteClub(club.ID);
            
            // Remove the club from recent clubs since it no longer exists
            removeRecentClub(club.ID);
            
            // Show success message and navigate back to clubs list
            alert(t('clubs.hardDeleteSuccess'));
            navigate('/clubs');
        } catch (error) {
            console.error('Error permanently deleting club:', error);
            alert(t('clubs.hardDeleteError'));
        } finally {
            setShowHardDeletePopup(false);
        }
    };

    const handleLogoUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;

        // Validate file type
        const allowedTypes = ['image/png', 'image/jpeg', 'image/jpg', 'image/webp', 'image/svg+xml'];
        if (!allowedTypes.includes(file.type)) {
            setLogoError('Please select a PNG, JPEG, WebP, or SVG image file');
            return;
        }

        // Validate file size (5MB max)
        if (file.size > 5 * 1024 * 1024) {
            setLogoError('File size must be less than 5MB');
            return;
        }

        setLogoUploading(true);
        setLogoError(null);

        try {
            const formData = new FormData();
            formData.append('logo', file);

            // OData v2: Use custom handler for logo upload (multipart/form-data)
            const response = await api.post(`/api/v2/Clubs('${id}')/UploadLogo`, formData);

            // Update the club state with the new logo URL
            if (club) {
                setClub({ ...club, LogoURL: response.data.LogoURL });
            }
        } catch (err: unknown) {
            console.error('Error uploading logo:', err);
            setLogoError('Failed to upload logo');
        } finally {
            setLogoUploading(false);
            // Reset the file input
            event.target.value = '';
        }
    };

    const handleLogoDelete = async () => {
        if (!club?.LogoURL) return;

        setLogoUploading(true);
        setLogoError(null);

        try {
            // OData v2: Use DeleteLogo action on Club entity
            await api.post(`/api/v2/Clubs('${id}')/DeleteLogo`, {});
            
            // Update the club state to remove logo URL
            setClub({ ...club, LogoURL: undefined });
        } catch (err: unknown) {
            console.error('Error deleting logo:', err);
            setLogoError('Failed to delete logo');
        } finally {
            setLogoUploading(false);
        }
    };

    if (loading) return <div>Loading...</div>;
    if (clubNotFound) return <ClubNotFound clubId={id} title="Club Administration Not Available" message="The club you are trying to manage does not exist or you do not have admin access to it." />;
    if (error) return <div className="error">{error}</div>;
    if (!club) return <div>Club not found</div>;

    return (
        <Layout title={`${club.Name} - Admin`}>
            <div>
                <div className="tabs-container">
                    <nav className="tabs-nav">
                        <Link 
                            to={`/clubs/${id}/admin`}
                            className={`tab-button ${activeTab === 'overview' ? 'active' : ''}`}
                        >
                            {t('clubs.overview')}
                        </Link>
                        <Link 
                            to={`/clubs/${id}/admin/members`}
                            className={`tab-button ${activeTab === 'members' ? 'active' : ''}`}
                        >
                            {t('clubs.members')}
                        </Link>
                        {clubSettings?.TeamsEnabled && (
                            <Link 
                                to={`/clubs/${id}/admin/teams`}
                                className={`tab-button ${activeTab === 'teams' ? 'active' : ''}`}
                            >
                                {t('clubs.teams')}
                            </Link>
                        )}
                        {clubSettings?.FinesEnabled && (
                            <Link 
                                to={`/clubs/${id}/admin/fines`}
                                className={`tab-button ${activeTab === 'fines' ? 'active' : ''}`}
                            >
                                {t('clubs.fines')}
                            </Link>
                        )}

                        <Link 
                            to={`/clubs/${id}/admin/events`}
                            className={`tab-button ${activeTab === 'events' ? 'active' : ''}`}
                        >
                            {t('clubs.events')}
                        </Link>
                        <Link 
                            to={`/clubs/${id}/admin/news`}
                            className={`tab-button ${activeTab === 'news' ? 'active' : ''}`}
                        >
                            {t('clubs.news')}
                        </Link>
                        <Link 
                            to={`/clubs/${id}/admin/settings`}
                            className={`tab-button ${activeTab === 'settings' ? 'active' : ''}`}
                        >
                            {t('clubs.settings')}
                        </Link>
                    </nav>

                    <div className="tab-content">
                        <div className={`tab-panel ${activeTab === 'overview' ? 'active' : ''}`}>
                            {isEditing ? (
                                <div className="edit-form">
                                    <Input
                                        label={t('clubs.clubName')}
                                        id="clubName"
                                        type="text"
                                        value={editForm.name}
                                        onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                                        placeholder={t('clubs.clubName')}
                                    />
                                    <Input
                                        label={t('clubs.description')}
                                        value={editForm.description}
                                        onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                                        placeholder={t('clubs.clubDescription')}
                                        multiline
                                        rows={4}
                                    />
                                    <div className="form-actions">
                                        <Button variant="accept" onClick={updateClub}>{t('common.save')}</Button>
                                        <Button variant="cancel" onClick={() => setIsEditing(false)}>{t('common.cancel')}</Button>
                                    </div>
                                </div>
                            ) : (
                                <>
                                    <PageHeader
                                        actions={
                                            <>
                                                {!club.Deleted && (
                                                    <>
                                                        <Button variant="accept" onClick={handleEdit}>{t('clubs.editClub')}</Button>
                                                        {isOwner && (
                                                            <Button 
                                                                variant="cancel"
                                                                onClick={handleDeleteClub}
                                                            >
                                                                {t('clubs.deleteClub')}
                                                            </Button>
                                                        )}
                                                    </>
                                                )}
                                                {club.Deleted && isOwner && (
                                                    <Button 
                                                        variant="cancel"
                                                        onClick={handleHardDeleteClub}
                                                        className="hard-delete-button"
                                                    >
                                                        {t('clubs.hardDeleteClub')}
                                                    </Button>
                                                )}
                                            </>
                                        }
                                    >
                                        <div className="club-logo-section">
                                            {/* Single file input for both logo states */}
                                            {!club.Deleted && (
                                                <input
                                                    ref={fileInputRef}
                                                    type="file"
                                                    accept="image/png,image/jpeg,image/jpg,image/webp,image/svg+xml"
                                                    onChange={handleLogoUpload}
                                                    style={{ display: 'none' }}
                                                />
                                            )}
                                            {club.LogoURL ? (
                                                <div className="club-logo-container">
                                                    <img 
                                                        src={club.LogoURL} 
                                                        alt={`${club.Name} logo`}
                                                        className="club-logo"
                                                    />
                                                    {!club.Deleted && (
                                                        <div className="logo-actions">
                                                            <Button
                                                                size="sm"
                                                                variant="secondary"
                                                                onClick={() => fileInputRef.current?.click()}
                                                                disabled={logoUploading}
                                                            >
                                                                {logoUploading ? 'Uploading...' : 'Change'}
                                                            </Button>
                                                            <Button
                                                                size="sm"
                                                                variant="cancel"
                                                                onClick={handleLogoDelete}
                                                                disabled={logoUploading}
                                                            >
                                                                Delete
                                                            </Button>
                                                        </div>
                                                    )}
                                                </div>
                                            ) : (
                                                <div className="club-logo-placeholder">
                                                    <div 
                                                        className="logo-placeholder"
                                                        onClick={!club.Deleted ? () => fileInputRef.current?.click() : undefined}
                                                    >
                                                        {!club.Deleted ? 'Click to upload logo' : 'No logo'}
                                                    </div>
                                                </div>
                                            )}
                                        </div>
                                        <div className="club-details">
                                            <h2>{club.Name}</h2>
                                            <p>{club.Description}</p>
                                            {logoError && (
                                                <div className="logo-error">
                                                    {logoError}
                                                </div>
                                            )}
                                        </div>
                                    </PageHeader>
                                    {club.Deleted && (
                                        <div className="club-deleted-notice" style={{ 
                                            backgroundColor: '#f44336', 
                                            color: 'white', 
                                            padding: '15px', 
                                            marginTop: '15px',
                                            borderRadius: '4px',
                                            fontWeight: 'bold',
                                            border: '1px solid #d32f2f'
                                        }}>
                                            <strong>{t('clubs.clubDeleted')}</strong>
                                        </div>
                                    )}
                                    
                                    {/* Statistics Dashboard */}
                                    <div className="dashboard-grid">
                                        <StatisticsCard
                                            title={t('clubs.members')}
                                            value={memberStats.total}
                                            icon="👥"
                                            subtitle={`${memberStats.pendingInvites} pending invites`}
                                            trend={
                                                memberStats.newThisMonth > 0
                                                    ? { value: memberStats.newThisMonth, isPositive: true }
                                                    : undefined
                                            }
                                            loading={statsLoading}
                                        />
                                    </div>
                                </>
                            )}
                        </div>

                        <div className={`tab-panel ${activeTab === 'members' ? 'active' : ''}`}>
                            <AdminClubMemberList openJoinRequests={searchParams.get('openJoinRequests') === 'true'} />
                        </div>

                        {clubSettings?.TeamsEnabled && (
                            <div className={`tab-panel ${activeTab === 'teams' ? 'active' : ''}`}>
                                <AdminClubTeamList />
                            </div>
                        )}

                        {clubSettings?.FinesEnabled && (
                            <div className={`tab-panel ${activeTab === 'fines' ? 'active' : ''}`}>
                                <AdminClubFineList />
                            </div>
                        )}


                        <div className={`tab-panel ${activeTab === 'events' ? 'active' : ''}`}>
                            <AdminClubEventList />
                        </div>

                        <div className={`tab-panel ${activeTab === 'news' ? 'active' : ''}`}>
                            <AdminClubNewsList />
                        </div>

                        <div className={`tab-panel ${activeTab === 'settings' ? 'active' : ''}`}>
                            <AdminClubSettings onSettingsUpdate={refetchSettings} />
                        </div>
                    </div>
                </div>
            </div>

            {/* Delete Confirmation Popup */}
            <Modal 
                isOpen={showDeletePopup} 
                onClose={() => setShowDeletePopup(false)} 
                title={t('clubs.deleteClub')}
                maxWidth="450px"
            >
                <Modal.Body>
                    <p>{t('clubs.deleteConfirmation', { clubName: club?.Name })}</p>
                </Modal.Body>
                <Modal.Actions>
                    <Button 
                        variant="cancel"
                        onClick={confirmDeleteClub}
                    >
                        {t('common.delete')}
                    </Button>
                    <Button 
                        variant="accept"
                        onClick={() => setShowDeletePopup(false)}
                    >
                        {t('common.cancel')}
                    </Button>
                </Modal.Actions>
            </Modal>

            {/* Hard Delete Confirmation Popup */}
            <Modal 
                isOpen={showHardDeletePopup} 
                onClose={() => setShowHardDeletePopup(false)} 
                title={t('clubs.hardDeleteClub')}
                maxWidth="450px"
            >
                <Modal.Body>
                    <p>{t('clubs.hardDeleteConfirmation', { clubName: club?.Name })}</p>
                </Modal.Body>
                <Modal.Actions>
                    <Button 
                        variant="cancel"
                        onClick={confirmHardDeleteClub}
                        style={{ backgroundColor: '#d32f2f', borderColor: '#d32f2f' }}
                    >
                        {t('clubs.hardDeleteClub')}
                    </Button>
                    <Button 
                        variant="accept"
                        onClick={() => setShowHardDeletePopup(false)}
                    >
                        {t('common.cancel')}
                    </Button>
                </Modal.Actions>
            </Modal>
        </Layout>
    );
};

export default AdminClubDetails;
