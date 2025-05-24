import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../../utils/api';
import Layout from '../../layout/Layout';
import AdminClubMemberList from './AdminClubMemberList';
import AdminClubPendingInviteList from './AdminClubPendingInviteList';
import AdminClubFineList from './AdminClubFineList';
import AdminClubShiftList from './AdminClubShiftList';

interface Club {
    id: string;
    name: string;
    description: string;
}

const AdminClubDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [club, setClub] = useState<Club | null>(null);
    
    
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [editForm, setEditForm] = useState({ name: '', description: '' });

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
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <h2>{club.name}</h2>
                            <button onClick={handleEdit}>Edit Club</button>
                        </div>
                        <p>{club.description}</p>
                    </>
                )}

                <div className="club-info">
                    <AdminClubMemberList />
                    <AdminClubFineList />
                    <AdminClubPendingInviteList />
                    <AdminClubShiftList />
                </div>
            </div>
        </Layout>
    );
};

export default AdminClubDetails;