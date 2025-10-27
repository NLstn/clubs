import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../../components/layout/Layout';
import { Input, Modal } from '@/components/ui';
import './AdminTeamDetails.css';

interface Team {
    id: string;
    name: string;
    description: string;
    createdAt: string;
    clubId: string;
}

interface TeamStats {
    member_count: number;
    admin_count: number;
    upcoming_events: number;
    total_events: number;
    unpaid_fines: number;
    total_fines: number;
}

interface TeamOverview {
    team: Team;
    stats: TeamStats;
    user_role: string;
    is_admin: boolean;
}

interface Event {
    id: string;
    name: string;
    description: string;
    location: string;
    start_time: string;
    end_time: string;
}

interface Fine {
    id: string;
    userId: string;
    userName: string;
    reason: string;
    amount: number;
    createdAt: string;
    paid: boolean;
}

const AdminTeamDetails = () => {
    const { clubId, teamId } = useParams();
    const navigate = useNavigate();
    const [teamOverview, setTeamOverview] = useState<TeamOverview | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [activeTab, setActiveTab] = useState('overview');
    const [isEditing, setIsEditing] = useState(false);
    const [editForm, setEditForm] = useState({ name: '', description: '' });

    // Events state
    const [events, setEvents] = useState<Event[]>([]);
    const [showCreateEventModal, setShowCreateEventModal] = useState(false);
    const [newEvent, setNewEvent] = useState({
        name: '',
        description: '',
        location: '',
        start_time: '',
        end_time: ''
    });

    // Fines state
    const [fines, setFines] = useState<Fine[]>([]);
    const [showCreateFineModal, setShowCreateFineModal] = useState(false);
    const [newFine, setNewFine] = useState({
        userId: '',
        reason: '',
        amount: 0
    });

    useEffect(() => {
        const fetchTeamData = async () => {
            if (!clubId || !teamId) {
                setError('Missing club or team ID');
                setLoading(false);
                return;
            }

            try {
                const overviewResponse = await api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/overview`);
                setTeamOverview(overviewResponse.data);
                
                if (!overviewResponse.data.is_admin) {
                    navigate(`/clubs/${clubId}/teams/${teamId}`);
                    return;
                }

                setError(null);
            } catch (err: unknown) {
                console.error('Error fetching team details:', err);
                if (err && typeof err === 'object' && 'response' in err) {
                    const axiosError = err as { response?: { status?: number } };
                    if (axiosError.response?.status === 404) {
                        setError('Team not found');
                    } else if (axiosError.response?.status === 403) {
                        setError('You do not have admin access to this team');
                    } else {
                        setError('Failed to load team details');
                    }
                } else {
                    setError('Failed to load team details');
                }
            } finally {
                setLoading(false);
            }
        };

        fetchTeamData();
    }, [clubId, teamId, navigate]);

    const fetchEvents = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/events`);
            setEvents(response.data || []);
        } catch (err) {
            console.error('Error fetching events:', err);
        }
    }, [clubId, teamId]);

    const fetchFines = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/fines`);
            setFines(response.data || []);
        } catch (err) {
            console.error('Error fetching fines:', err);
        }
    }, [clubId, teamId]);

    useEffect(() => {
        if (activeTab === 'events' && teamOverview) {
            fetchEvents();
        }
        if (activeTab === 'fines' && teamOverview) {
            fetchFines();
        }
    }, [activeTab, teamOverview, fetchEvents, fetchFines]);

    const handleEdit = () => {
        if (teamOverview?.team) {
            setEditForm({ 
                name: teamOverview.team.name, 
                description: teamOverview.team.description 
            });
            setIsEditing(true);
        }
    };

    const updateTeam = async () => {
        try {
            await api.put(`/api/v1/clubs/${clubId}/teams/${teamId}`, editForm);
            // Refresh team data
            const response = await api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/overview`);
            setTeamOverview(response.data);
            setIsEditing(false);
            setError(null);
        } catch {
            setError('Failed to update team');
        }
    };

    const createEvent = async () => {
        if (!newEvent.name || !newEvent.start_time || !newEvent.end_time) {
            setError('Please fill in all required fields');
            return;
        }

        try {
            await api.post(`/api/v1/clubs/${clubId}/teams/${teamId}/events`, newEvent);
            setNewEvent({
                name: '',
                description: '',
                location: '',
                start_time: '',
                end_time: ''
            });
            setShowCreateEventModal(false);
            fetchEvents();
        } catch {
            setError('Failed to create event');
        }
    };

    const deleteEvent = async (eventId: string) => {
        if (!confirm('Are you sure you want to delete this event?')) return;

        try {
            await api.delete(`/api/v1/clubs/${clubId}/teams/${teamId}/events/${eventId}`);
            fetchEvents();
        } catch {
            setError('Failed to delete event');
        }
    };

    const createFine = async () => {
        if (!newFine.userId || !newFine.reason || newFine.amount <= 0) {
            setError('Please fill in all required fields');
            return;
        }

        try {
            await api.post(`/api/v1/clubs/${clubId}/teams/${teamId}/fines`, newFine);
            setNewFine({
                userId: '',
                reason: '',
                amount: 0
            });
            setShowCreateFineModal(false);
            fetchFines();
        } catch {
            setError('Failed to create fine');
        }
    };

    const deleteFine = async (fineId: string) => {
        if (!confirm('Are you sure you want to delete this fine?')) return;

        try {
            await api.delete(`/api/v1/clubs/${clubId}/teams/${teamId}/fines/${fineId}`);
            fetchFines();
        } catch {
            setError('Failed to delete fine');
        }
    };

    if (loading) return <div>Loading team admin...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!teamOverview) return <div>Team not found</div>;

    const { team, stats } = teamOverview;

    return (
        <Layout title={`${team.name} - Admin`}>
            <div className="admin-team-container">
                <div className="tabs-container">
                    <nav className="tabs-nav">
                        <button 
                            className={`tab-button ${activeTab === 'overview' ? 'active' : ''}`}
                            onClick={() => setActiveTab('overview')}
                        >
                            Overview
                        </button>
                        <button 
                            className={`tab-button ${activeTab === 'events' ? 'active' : ''}`}
                            onClick={() => setActiveTab('events')}
                        >
                            Events
                        </button>
                        <button 
                            className={`tab-button ${activeTab === 'fines' ? 'active' : ''}`}
                            onClick={() => setActiveTab('fines')}
                        >
                            Fines
                        </button>
                    </nav>

                    <div className="tab-content">
                        <div className={`tab-panel ${activeTab === 'overview' ? 'active' : ''}`}>
                            {isEditing ? (
                                <div className="edit-form">
                                    <Input
                                        label="Team Name"
                                        type="text"
                                        value={editForm.name}
                                        onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                                        placeholder="Team Name"
                                    />
                                    <Input
                                        label="Description"
                                        value={editForm.description}
                                        onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                                        placeholder="Team Description"
                                        multiline
                                        rows={4}
                                    />
                                    <div className="form-actions">
                                        <button onClick={updateTeam} className="button-accept">Save</button>
                                        <button onClick={() => setIsEditing(false)} className="button-cancel">Cancel</button>
                                    </div>
                                </div>
                            ) : (
                                <>
                                    <div className="team-header">
                                        <div className="team-info">
                                            <h2>{team.name}</h2>
                                            <p>{team.description}</p>
                                        </div>
                                        <div className="team-actions">
                                            <button onClick={handleEdit} className="button-accept">Edit Team</button>
                                            <button 
                                                onClick={() => navigate(`/clubs/${clubId}/teams/${teamId}`)} 
                                                className="button-secondary"
                                            >
                                                View Team
                                            </button>
                                        </div>
                                    </div>

                                    {/* Team Stats */}
                                    <div className="team-stats-section">
                                        <h3>Team Statistics</h3>
                                        <div className="stats-grid">
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.member_count}</div>
                                                <div className="stat-label">Members</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.admin_count}</div>
                                                <div className="stat-label">Admins</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.upcoming_events}</div>
                                                <div className="stat-label">Upcoming Events</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.total_events}</div>
                                                <div className="stat-label">Total Events</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.unpaid_fines}</div>
                                                <div className="stat-label">Unpaid Fines</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.total_fines}</div>
                                                <div className="stat-label">Total Fines</div>
                                            </div>
                                        </div>
                                    </div>
                                </>
                            )}
                        </div>

                        <div className={`tab-panel ${activeTab === 'events' ? 'active' : ''}`}>
                            <div className="section-header">
                                <h3>Team Events</h3>
                                <button 
                                    onClick={() => setShowCreateEventModal(true)} 
                                    className="button-accept"
                                >
                                    Create Event
                                </button>
                            </div>

                            <div className="events-list">
                                {events.map(event => (
                                    <div key={event.id} className="event-card">
                                        <div className="event-header">
                                            <h4>{event.name}</h4>
                                            <button 
                                                onClick={() => deleteEvent(event.id)}
                                                className="button-cancel button-sm"
                                            >
                                                Delete
                                            </button>
                                        </div>
                                        <p>{event.description}</p>
                                        <div className="event-details">
                                            <span>📅 {new Date(event.start_time).toLocaleString()}</span>
                                            {event.location && <span>📍 {event.location}</span>}
                                        </div>
                                    </div>
                                ))}
                                {events.length === 0 && (
                                    <div className="no-content">No events created yet.</div>
                                )}
                            </div>
                        </div>

                        <div className={`tab-panel ${activeTab === 'fines' ? 'active' : ''}`}>
                            <div className="section-header">
                                <h3>Team Fines</h3>
                                <button 
                                    onClick={() => setShowCreateFineModal(true)} 
                                    className="button-accept"
                                >
                                    Create Fine
                                </button>
                            </div>

                            <div className="fines-list">
                                {fines.map(fine => (
                                    <div key={fine.id} className="fine-card">
                                        <div className="fine-header">
                                            <h4>{fine.userName}</h4>
                                            <div className="fine-amount">${fine.amount}</div>
                                            <button 
                                                onClick={() => deleteFine(fine.id)}
                                                className="button-cancel button-sm"
                                            >
                                                Delete
                                            </button>
                                        </div>
                                        <p>{fine.reason}</p>
                                        <div className="fine-status">
                                            Status: {fine.paid ? 'Paid' : 'Unpaid'}
                                        </div>
                                    </div>
                                ))}
                                {fines.length === 0 && (
                                    <div className="no-content">No fines created yet.</div>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Create Event Modal */}
            <Modal isOpen={showCreateEventModal} onClose={() => setShowCreateEventModal(false)} title="Create Event">
                <Modal.Body>
                    <div className="modal-form-section">
                        <Input
                            label="Event Name"
                            type="text"
                            value={newEvent.name}
                            onChange={(e) => setNewEvent({ ...newEvent, name: e.target.value })}
                            placeholder="Event name"
                        />
                        <Input
                            label="Description"
                            value={newEvent.description}
                            onChange={(e) => setNewEvent({ ...newEvent, description: e.target.value })}
                            placeholder="Event description"
                            multiline
                            rows={3}
                        />
                        <Input
                            label="Location"
                            type="text"
                            value={newEvent.location}
                            onChange={(e) => setNewEvent({ ...newEvent, location: e.target.value })}
                            placeholder="Event location"
                        />
                        <Input
                            label="Start Time"
                            type="datetime-local"
                            value={newEvent.start_time}
                            onChange={(e) => setNewEvent({ ...newEvent, start_time: e.target.value })}
                        />
                        <Input
                            label="End Time"
                            type="datetime-local"
                            value={newEvent.end_time}
                            onChange={(e) => setNewEvent({ ...newEvent, end_time: e.target.value })}
                        />
                    </div>
                </Modal.Body>
                <Modal.Actions>
                    <button onClick={createEvent} className="button-accept">
                        Create Event
                    </button>
                    <button onClick={() => setShowCreateEventModal(false)} className="button-cancel">
                        Cancel
                    </button>
                </Modal.Actions>
            </Modal>

            {/* Create Fine Modal */}
            <Modal isOpen={showCreateFineModal} onClose={() => setShowCreateFineModal(false)} title="Create Fine">
                <Modal.Body>
                    <div className="modal-form-section">
                        <Input
                            label="User ID"
                            type="text"
                            value={newFine.userId}
                            onChange={(e) => setNewFine({ ...newFine, userId: e.target.value })}
                            placeholder="User ID to fine"
                        />
                        <Input
                            label="Reason"
                            value={newFine.reason}
                            onChange={(e) => setNewFine({ ...newFine, reason: e.target.value })}
                            placeholder="Reason for fine"
                            multiline
                            rows={3}
                        />
                        <Input
                            label="Amount"
                            type="number"
                            value={newFine.amount.toString()}
                            onChange={(e) => setNewFine({ ...newFine, amount: parseFloat(e.target.value) || 0 })}
                            placeholder="Fine amount"
                        />
                    </div>
                </Modal.Body>
                <Modal.Actions>
                    <button onClick={createFine} className="button-accept">
                        Create Fine
                    </button>
                    <button onClick={() => setShowCreateFineModal(false)} className="button-cancel">
                        Cancel
                    </button>
                </Modal.Actions>
            </Modal>
        </Layout>
    );
};

export default AdminTeamDetails;