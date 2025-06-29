import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../../utils/api';
import Layout from '../../../components/layout/Layout';
import AdminClubMemberList from './members/AdminClubMemberList';
import AdminClubFineList from './fines/AdminClubFineList';
import AdminClubEventList from './events/AdminClubEventList';
import AdminClubNewsList from './news/AdminClubNewsList';
import AdminClubSettings from './settings/AdminClubSettings';
import { useClubSettings } from '../../../hooks/useClubSettings';
import { useT } from '../../../hooks/useTranslation';

interface Club {
    id: string;
    name: string;
    description: string;
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
    const [isEditing, setIsEditing] = useState(false);
    const [editForm, setEditForm] = useState({ name: '', description: '' });
    const [activeTab, setActiveTab] = useState('overview');
    const [isOwner, setIsOwner] = useState(false);

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
            } catch (err: Error | unknown) {
                console.error('Error fetching club details:', err instanceof Error ? err.message : 'Unknown error');
                setError('Error fetching club details');
                setLoading(false);
            }
        };

        fetchData();
    }, [id, navigate]);

    // Reset to valid tab if current tab becomes unavailable
    useEffect(() => {
        if (clubSettings) {
            if (activeTab === 'fines' && !clubSettings.finesEnabled) {
                setActiveTab('overview');
            }
        }
    }, [clubSettings, activeTab]);

    const updateClub = async () => {
        try {
            const response = await api.patch(`/api/v1/clubs/${id}`, editForm);
            setClub(response.data);
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
        
        const confirmText = t('clubs.deleteConfirmation', { clubName: club.name });
        if (!confirm(confirmText)) {
            return;
        }

        try {
            await api.delete(`/api/v1/clubs/${id}`);
            // Navigate to clubs list after deletion
            navigate('/');
        } catch (err: Error | unknown) {
            console.error('Error deleting club:', err instanceof Error ? err.message : 'Unknown error');
            setError('Failed to delete club');
        }
    };

    if (loading) return <div>Loading...</div>;
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
                                        <h2>{club.name}</h2>
                                        <div className="club-actions">
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
                                        </div>
                                    </div>
                                    <p>{club.description}</p>
                                    {club.deleted && (
                                        <div className="club-deleted-notice" style={{ 
                                            backgroundColor: '#ffebee', 
                                            border: '1px solid #f44336', 
                                            padding: '10px', 
                                            marginTop: '10px',
                                            borderRadius: '4px'
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
        </Layout>
    );
};

export default AdminClubDetails;