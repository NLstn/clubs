import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../../utils/api';
import Layout from '../../../components/layout/Layout';
import AdminClubMemberList from './members/AdminClubMemberList';
import AdminClubFineList from './fines/AdminClubFineList';
import AdminClubShiftList from './shifts/AdminClubShiftList';
import AdminClubEventList from './events/AdminClubEventList';
import AdminClubNewsList from './news/AdminClubNewsList';
import AdminClubSettings from './settings/AdminClubSettings';
import { useClubSettings } from '../../../hooks/useClubSettings';

interface Club {
    id: string;
    name: string;
    description: string;
}

const AdminClubDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [club, setClub] = useState<Club | null>(null);
    const { settings: clubSettings, refetch: refetchSettings } = useClubSettings(id);
    
    
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [editForm, setEditForm] = useState({ name: '', description: '' });
    const [activeTab, setActiveTab] = useState('overview');

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
            if (activeTab === 'shifts' && !clubSettings.shiftsEnabled) {
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
                            Overview
                        </button>
                        <button 
                            className={`tab-button ${activeTab === 'members' ? 'active' : ''}`}
                            onClick={() => setActiveTab('members')}
                        >
                            Members
                        </button>
                        {clubSettings?.finesEnabled && (
                            <button 
                                className={`tab-button ${activeTab === 'fines' ? 'active' : ''}`}
                                onClick={() => setActiveTab('fines')}
                            >
                                Fines
                            </button>
                        )}
                        {clubSettings?.shiftsEnabled && (
                            <button 
                                className={`tab-button ${activeTab === 'shifts' ? 'active' : ''}`}
                                onClick={() => setActiveTab('shifts')}
                            >
                                Shifts
                            </button>
                        )}

                        <button 
                            className={`tab-button ${activeTab === 'events' ? 'active' : ''}`}
                            onClick={() => setActiveTab('events')}
                        >
                            Events
                        </button>
                        <button 
                            className={`tab-button ${activeTab === 'news' ? 'active' : ''}`}
                            onClick={() => setActiveTab('news')}
                        >
                            News
                        </button>
                        <button 
                            className={`tab-button ${activeTab === 'settings' ? 'active' : ''}`}
                            onClick={() => setActiveTab('settings')}
                        >
                            Settings
                        </button>
                    </nav>

                    <div className="tab-content">
                        <div className={`tab-panel ${activeTab === 'overview' ? 'active' : ''}`}>
                            {isEditing ? (
                                <div className="edit-form">
                                    <div className="form-group">
                                        <label htmlFor="clubName">Club Name</label>
                                        <input
                                            id="clubName"
                                            type="text"
                                            value={editForm.name}
                                            onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                                            placeholder="Club Name"
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label htmlFor="clubDescription">Description</label>
                                        <textarea
                                            id="clubDescription"
                                            value={editForm.description}
                                            onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                                            placeholder="Club Description"
                                        />
                                    </div>
                                    <div className="form-actions">
                                        <button onClick={updateClub} className="button-accept">Save</button>
                                        <button onClick={() => setIsEditing(false)} className="button-cancel">Cancel</button>
                                    </div>
                                </div>
                            ) : (
                                <>
                                    <div className="club-header">
                                        <h2>{club.name}</h2>
                                        <button onClick={handleEdit}>Edit Club</button>
                                    </div>
                                    <p>{club.description}</p>
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

                        {clubSettings?.shiftsEnabled && (
                            <div className={`tab-panel ${activeTab === 'shifts' ? 'active' : ''}`}>
                                <AdminClubShiftList />
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