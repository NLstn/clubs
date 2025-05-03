import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../layout/Layout';
import InviteMember from './InviteMember';

interface Member {
    id: string;
    name: string;
    role: string;
}

interface Club {
    id: string;
    name: string;
    description: string;
}

interface Events {
    id: string;
    name: string;
    date: string;
    description: string;
    begin_time: string;
    end_time: string;
}

interface JoinRequest {
    id: string;
    email: string;
}

const AdminClubDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [club, setClub] = useState<Club | null>(null);
    const [members, setMembers] = useState<Member[]>([]);
    const [events, setEvents] = useState<Events[]>([]);
    const [joinRequests, setJoinRequests] = useState<JoinRequest[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [editForm, setEditForm] = useState({ name: '', description: '' });

    const handleRoleChange = async (memberId: string, newRole: string) => {
        try {
            await api.patch(`/api/v1/clubs/${id}/members/${memberId}`, { role: newRole });
            setMembers(members.map(member => 
                member.id === memberId ? { ...member, role: newRole } : member
            ));
        } catch {
            setError('Failed to update member role');
        }
    };

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [adminResponse, clubResponse, membersResponse, eventsResponse, joinRequestsResponse] = await Promise.all([
                    api.get(`/api/v1/clubs/${id}/isAdmin`),
                    api.get(`/api/v1/clubs/${id}`),
                    api.get(`/api/v1/clubs/${id}/members`),
                    api.get(`/api/v1/clubs/${id}/events`),
                    api.get(`/api/v1/clubs/${id}/joinRequests`)
                ]);

                if (!adminResponse.data.isAdmin) {
                    navigate(`/clubs/${id}`);
                    return;
                }

                setClub(clubResponse.data);
                setMembers(membersResponse.data);
                setEvents(eventsResponse.data);
                setJoinRequests(joinRequestsResponse.data);
                setLoading(false);
            } catch (err: Error | unknown) {
                console.error('Error fetching club details:', err instanceof Error ? err.message : 'Unknown error');
                setError('Error fetching club details');
                setLoading(false);
            }
        };

        fetchData();
    }, [id, navigate]);

    const deleteMember = async (memberId: string) => {
        try {
            await api.delete(`/api/v1/clubs/${id}/members/${memberId}`);
            setMembers(members.filter(member => member.id !== memberId));
        } catch {
            setError('Failed to delete member');
        }
    };

    const sendInvite = async (email: string) => {
        try {
            await api.post(`/api/v1/clubs/${id}/joinRequests`, { email });
            setIsModalOpen(false);
        } catch {
            setError('Failed to send invite');
        }
    };

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
                    <h3>Members</h3>
                    <table className="basic-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Role</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {members && members.map((member) => (
                                <tr key={member.id}>
                                    <td>{member.name}</td>
                                    <td>
                                        <select 
                                            value={member.role} 
                                            onChange={(e) => handleRoleChange(member.id, e.target.value)}
                                        >
                                            <option value="member">Member</option>
                                            <option value="admin">Admin</option>
                                            <option value="owner">Owner</option>
                                        </select>
                                    </td>
                                    <td className="delete-cell">
                                        <button
                                            onClick={() => deleteMember(member.id)}
                                            className="delete-button"
                                            aria-label="Delete member"
                                        >
                                            ×
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                    <button onClick={() => setIsModalOpen(true)}>Invite Member</button>
                    <InviteMember
                        isOpen={isModalOpen}
                        onClose={() => setIsModalOpen(false)}
                        onSubmit={sendInvite}
                    />
                    <h3>Events</h3>
                    <table className="basic-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Date</th>
                                <th>Description</th>
                                <th>Begin Time</th>
                                <th>End Time</th>
                            </tr>
                        </thead>
                        <tbody>
                            {events.map((event) => (
                                <tr key={event.id}>
                                    <td>{event.name}</td>
                                    <td>{event.date}</td>
                                    <td>{event.description}</td>
                                    <td>{event.begin_time}</td>
                                    <td>{event.end_time}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                    <h3>Pending Invites</h3>
                    <table className="basic-table">
                        <thead>
                            <tr>
                                <th>Email</th>
                            </tr>
                        </thead>
                        <tbody>
                            {joinRequests.map((request) => (
                                <tr key={request.id}>
                                    <td>{request.email}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </Layout>
    );
};

export default AdminClubDetails;