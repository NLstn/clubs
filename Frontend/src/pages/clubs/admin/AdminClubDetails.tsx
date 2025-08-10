import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api, { hardDeleteClub } from '../../../utils/api';
import Layout from '../../../components/layout/Layout';
import ClubNotFound from '../ClubNotFound';
import AdminClubMemberList from './members/AdminClubMemberList';
import AdminClubTeamList from './teams/AdminClubTeamList';
import AdminClubFineList from './fines/AdminClubFineList';
import AdminClubEventList from './events/AdminClubEventList';
import AdminClubNewsList from './news/AdminClubNewsList';
import { Input, Modal } from '@/components/ui';
import AdminClubSettings from './settings/AdminClubSettings';
import { useClubSettings } from '../../../hooks/useClubSettings';
import { useT } from '../../../hooks/useTranslation';
import { removeRecentClub } from '../../../utils/recentClubs';
import AdminClubOverview from './AdminClubOverview';

interface Club {
    id: string;
    name: string;
    description: string;
    logo_url?: string;
    deleted?: boolean;
}

const AdminClubDetails = () => {
    const { t } = useT();
    const { id } = useParams();
    const navigate = useNavigate();
    const [club, setClub] = useState<Club | null>(null);
    const { settings: clubSettings, refetch: refetchSettings } = useClubSettings(id);
    
    
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [clubNotFound, setClubNotFound] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [editForm, setEditForm] = useState({ name: '', description: '' });
    const [activeTab, setActiveTab] = useState('overview');
    const [isOwner, setIsOwner] = useState(false);
    const [showDeletePopup, setShowDeletePopup] = useState(false);
    const [showHardDeletePopup, setShowHardDeletePopup] = useState(false);
    const [logoUploading, setLogoUploading] = useState(false);
    const [logoError, setLogoError] = useState<string | null>(null);
    const [openFinesCount, setOpenFinesCount] = useState(0);
    const [addEventTrigger, setAddEventTrigger] = useState(0);
    const [createTeamTrigger, setCreateTeamTrigger] = useState(0);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [adminResponse, clubResponse] = await Promise.all([
                    api.get(`/api/v1/clubs/${id}/isAdmin`),
                    api.get(`/api/v1/clubs/${id}`),
                    
                ]);

                if (!adminResponse.data.isAdmin) {
                    navigate(`/clubs/${id}`);
                    return;
                }

                setClub(clubResponse.data);
                setIsOwner(adminResponse.data.isOwner || false);
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
    }, [id, navigate, t]);

    // Reset to valid tab if current tab becomes unavailable
    useEffect(() => {
        if (clubSettings) {
            if (activeTab === 'fines' && !clubSettings.finesEnabled) {
                setActiveTab('overview');
            }
            if (activeTab === 'teams' && !clubSettings.teamsEnabled) {
                setActiveTab('overview');
            }
        }
    }, [clubSettings, activeTab]);

    useEffect(() => {
        const fetchOpenFines = async () => {
            if (!id || !clubSettings?.finesEnabled) {
                setOpenFinesCount(0);
                return;
            }
            try {
                const response = await api.get(`/api/v1/clubs/${id}/fines`);
                const fines = Array.isArray(response.data) ? response.data : [];
                setOpenFinesCount(fines.filter((f: { paid: boolean }) => !f.paid).length);
            } catch (err) {
                console.error('Error fetching open fines:', err);
            }
        };
        fetchOpenFines();
    }, [id, clubSettings?.finesEnabled]);

    const updateClub = async () => {
        try {
            const response = await api.patch(`/api/v1/clubs/${id}`, editForm);
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
            setEditForm({ name: club.name, description: club.description });
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
            await api.delete(`/api/v1/clubs/${id}`);
            // Note: We don't remove from recent clubs for soft delete since owners can still access it
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
            await hardDeleteClub(club.id);
            
            // Remove the club from recent clubs since it no longer exists
            removeRecentClub(club.id);
            
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

    const handleQuickAddEvent = () => {
        setActiveTab('events');
        setAddEventTrigger((n) => n + 1);
    };

    const handleQuickCreateTeam = () => {
        setActiveTab('teams');
        setCreateTeamTrigger((n) => n + 1);
    };

    const handleLogoUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;

        // Validate file type
        const allowedTypes = ['image/png', 'image/jpeg', 'image/jpg', 'image/webp'];
        if (!allowedTypes.includes(file.type)) {
            setLogoError('Please select a PNG, JPEG, or WebP image file');
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

            const response = await api.post(`/api/v1/clubs/${id}/logo`, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            // Update the club state with the new logo URL
            if (club) {
                setClub({ ...club, logo_url: response.data.logo_url });
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
        if (!club?.logo_url) return;

        setLogoUploading(true);
        setLogoError(null);

        try {
            await api.delete(`/api/v1/clubs/${id}/logo`);
            
            // Update the club state to remove logo URL
            setClub({ ...club, logo_url: undefined });
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
        <Layout title={`${club.name} - Admin`}>
            <div>
                <div className="tabs-container">
                    <nav className="tabs-nav">
                        <button 
                            className={`tab-button ${activeTab === 'overview' ? 'active' : ''}`}
                            onClick={() => setActiveTab('overview')}
                        >
                            {t('clubs.overview')}
                        </button>
                        <button 
                            className={`tab-button ${activeTab === 'members' ? 'active' : ''}`}
                            onClick={() => setActiveTab('members')}
                        >
                            {t('clubs.members')}
                        </button>
                        {clubSettings?.teamsEnabled && (
                            <button 
                                className={`tab-button ${activeTab === 'teams' ? 'active' : ''}`}
                                onClick={() => setActiveTab('teams')}
                            >
                                {t('clubs.teams')}
                            </button>
                        )}
                        {clubSettings?.finesEnabled && (
                            <button 
                                className={`tab-button ${activeTab === 'fines' ? 'active' : ''}`}
                                onClick={() => setActiveTab('fines')}
                            >
                                {t('clubs.fines')}
                            </button>
                        )}

                        <button 
                            className={`tab-button ${activeTab === 'events' ? 'active' : ''}`}
                            onClick={() => setActiveTab('events')}
                        >
                            {t('clubs.events')}
                        </button>
                        <button 
                            className={`tab-button ${activeTab === 'news' ? 'active' : ''}`}
                            onClick={() => setActiveTab('news')}
                        >
                            {t('clubs.news')}
                        </button>
                        <button 
                            className={`tab-button ${activeTab === 'settings' ? 'active' : ''}`}
                            onClick={() => setActiveTab('settings')}
                        >
                            {t('clubs.settings')}
                        </button>
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
                                        <button onClick={updateClub} className="button-accept">{t('common.save')}</button>
                                        <button onClick={() => setIsEditing(false)} className="button-cancel">{t('common.cancel')}</button>
                                    </div>
                                </div>
                            ) : (
                                <AdminClubOverview
                                    club={club}
                                    isOwner={isOwner}
                                    logoUploading={logoUploading}
                                    logoError={logoError}
                                    onEdit={handleEdit}
                                    onDelete={handleDeleteClub}
                                    onHardDelete={handleHardDeleteClub}
                                    onLogoUpload={handleLogoUpload}
                                    onLogoDelete={handleLogoDelete}
                                    openFinesCount={openFinesCount}
                                    onCreateEvent={handleQuickAddEvent}
                                    onCreateTeam={clubSettings?.teamsEnabled ? handleQuickCreateTeam : undefined}
                                    teamsEnabled={clubSettings?.teamsEnabled}
                                />
                            )}
                        </div>

                        <div className={`tab-panel ${activeTab === 'members' ? 'active' : ''}`}>
                            <AdminClubMemberList />
                        </div>

                        {clubSettings?.teamsEnabled && (
                            <div className={`tab-panel ${activeTab === 'teams' ? 'active' : ''}`}>
                                <AdminClubTeamList createTeamTrigger={createTeamTrigger} />
                            </div>
                        )}

                        {clubSettings?.finesEnabled && (
                            <div className={`tab-panel ${activeTab === 'fines' ? 'active' : ''}`}>
                                <AdminClubFineList />
                            </div>
                        )}


                        <div className={`tab-panel ${activeTab === 'events' ? 'active' : ''}`}>
                            <AdminClubEventList addEventTrigger={addEventTrigger} />
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
                    <p>{t('clubs.deleteConfirmation', { clubName: club?.name })}</p>
                </Modal.Body>
                <Modal.Actions>
                    <button 
                        onClick={confirmDeleteClub} 
                        className="button-cancel"
                    >
                        {t('common.delete')}
                    </button>
                    <button 
                        onClick={() => setShowDeletePopup(false)} 
                        className="button-accept"
                    >
                        {t('common.cancel')}
                    </button>
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
                    <p>{t('clubs.hardDeleteConfirmation', { clubName: club?.name })}</p>
                </Modal.Body>
                <Modal.Actions>
                    <button 
                        onClick={confirmHardDeleteClub} 
                        className="button-cancel"
                        style={{ backgroundColor: '#d32f2f', borderColor: '#d32f2f' }}
                    >
                        {t('clubs.hardDeleteClub')}
                    </button>
                    <button 
                        onClick={() => setShowHardDeletePopup(false)} 
                        className="button-accept"
                    >
                        {t('common.cancel')}
                    </button>
                </Modal.Actions>
            </Modal>
        </Layout>
    );
};

export default AdminClubDetails;