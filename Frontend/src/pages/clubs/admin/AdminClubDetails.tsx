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
import AdminClubSettings from './settings/AdminClubSettings';
import { useClubSettings } from '../../../hooks/useClubSettings';
import { useT } from '../../../hooks/useTranslation';
import { removeRecentClub } from '../../../utils/recentClubs';
import './AdminClubDetails.css';

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
                                    <div className="form-group">
                                        <label htmlFor="clubName">{t('clubs.clubName')}</label>
                                        <input
                                            id="clubName"
                                            type="text"
                                            value={editForm.name}
                                            onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                                            placeholder={t('clubs.clubName')}
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label htmlFor="clubDescription">{t('clubs.description')}</label>
                                        <textarea
                                            id="clubDescription"
                                            value={editForm.description}
                                            onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                                            placeholder={t('clubs.clubDescription')}
                                        />
                                    </div>
                                    <div className="form-actions">
                                        <button onClick={updateClub} className="button-accept">{t('common.save')}</button>
                                        <button onClick={() => setIsEditing(false)} className="button-cancel">{t('common.cancel')}</button>
                                    </div>
                                </div>
                            ) : (
                                <>
                                    <div className="club-header">
                                        <div className="club-info">
                                            <div className="club-logo-section">
                                                {club.logo_url ? (
                                                    <div className="club-logo-container">
                                                        <img 
                                                            src={club.logo_url} 
                                                            alt={`${club.name} logo`}
                                                            className="club-logo"
                                                        />
                                                        {!club.deleted && (
                                                            <div className="logo-actions">
                                                                <input
                                                                    type="file"
                                                                    id="logo-upload"
                                                                    accept="image/png,image/jpeg,image/jpg,image/webp"
                                                                    onChange={handleLogoUpload}
                                                                    style={{ display: 'none' }}
                                                                />
                                                                <button
                                                                    onClick={() => document.getElementById('logo-upload')?.click()}
                                                                    className="logo-change-btn"
                                                                    disabled={logoUploading}
                                                                >
                                                                    {logoUploading ? 'Uploading...' : 'Change'}
                                                                </button>
                                                                <button
                                                                    onClick={handleLogoDelete}
                                                                    className="logo-delete-btn"
                                                                    disabled={logoUploading}
                                                                >
                                                                    Delete
                                                                </button>
                                                            </div>
                                                        )}
                                                    </div>
                                                ) : (
                                                    <div className="club-logo-placeholder">
                                                        <div 
                                                            className="logo-placeholder"
                                                            onClick={!club.deleted ? () => document.getElementById('logo-upload')?.click() : undefined}
                                                        >
                                                            {!club.deleted ? 'Click to upload logo' : 'No logo'}
                                                        </div>
                                                        {!club.deleted && (
                                                            <input
                                                                type="file"
                                                                id="logo-upload"
                                                                accept="image/png,image/jpeg,image/jpg,image/webp"
                                                                onChange={handleLogoUpload}
                                                                style={{ display: 'none' }}
                                                            />
                                                        )}
                                                    </div>
                                                )}
                                            </div>
                                            <div className="club-details">
                                                <h2>{club.name}</h2>
                                                <p>{club.description}</p>
                                                {logoError && (
                                                    <div className="logo-error">
                                                        {logoError}
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                        <div className="club-actions">
                                            {!club.deleted && (
                                                <>
                                                    <button onClick={handleEdit} className="button-accept">{t('clubs.editClub')}</button>
                                                    {isOwner && (
                                                        <button 
                                                            onClick={handleDeleteClub} 
                                                            className="button-cancel"
                                                            style={{ marginLeft: '10px' }}
                                                        >
                                                            {t('clubs.deleteClub')}
                                                        </button>
                                                    )}
                                                </>
                                            )}
                                            {club.deleted && isOwner && (
                                                <button 
                                                    onClick={handleHardDeleteClub} 
                                                    className="button-cancel"
                                                    style={{ backgroundColor: '#d32f2f', borderColor: '#d32f2f' }}
                                                >
                                                    {t('clubs.hardDeleteClub')}
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                    {club.deleted && (
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
                                </>
                            )}
                        </div>

                        <div className={`tab-panel ${activeTab === 'members' ? 'active' : ''}`}>
                            <AdminClubMemberList />
                        </div>

                        {clubSettings?.teamsEnabled && (
                            <div className={`tab-panel ${activeTab === 'teams' ? 'active' : ''}`}>
                                <AdminClubTeamList />
                            </div>
                        )}

                        {clubSettings?.finesEnabled && (
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
            {showDeletePopup && (
                <div className="modal" onClick={() => setShowDeletePopup(false)}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <h2>{t('clubs.deleteClub')}</h2>
                        <p>{t('clubs.deleteConfirmation', { clubName: club?.name })}</p>
                        <div className="modal-actions">
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
                        </div>
                    </div>
                </div>
            )}

            {/* Hard Delete Confirmation Popup */}
            {showHardDeletePopup && (
                <div className="modal" onClick={() => setShowHardDeletePopup(false)}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <h2>{t('clubs.hardDeleteClub')}</h2>
                        <p>{t('clubs.hardDeleteConfirmation', { clubName: club?.name })}</p>
                        <div className="modal-actions">
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
                        </div>
                    </div>
                </div>
            )}
        </Layout>
    );
};

export default AdminClubDetails;