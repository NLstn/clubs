import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../../components/layout/Layout';
import { Input, Modal, Button, Tabs } from '@/components/ui';
import './AdminTeamDetails.css';

interface Team {
    ID: string;
    Name: string;
    Description: string;
    CreatedAt: string;
    ClubID: string;
}

interface TeamStats {
    MemberCount: number;
    AdminCount: number;
    UpcomingEvents: number;
    TotalEvents: number;
    UnpaidFines: number;
    TotalFines: number;
}

interface TeamOverview {
    Team: Team;
    Stats: TeamStats;
    UserRole: string;
    IsAdmin: boolean;
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
                // OData v2: Use GetOverview function on Team entity
                const overviewResponse = await api.get(`/api/v2/Teams('${teamId}')/GetOverview()`);
                setTeamOverview(overviewResponse.data);
                
                if (!overviewResponse.data.IsAdmin) {
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
            // OData v2: Use standard navigation to access Events
            const response = await api.get(`/api/v2/Teams('${teamId}')/Events`);
            setEvents(response.data?.value || []);
        } catch (err) {
            console.error('Error fetching events:', err);
        }
    }, [teamId]);

    const fetchFines = useCallback(async () => {
        try {
            // OData v2: Use GetFines function on Team entity
            const response = await api.get(`/api/v2/Teams('${teamId}')/GetFines()`);
            setFines(response.data || []);
        } catch (err) {
            console.error('Error fetching fines:', err);
        }
    }, [teamId]);

    useEffect(() => {
        if (activeTab === 'events' && teamOverview) {
            fetchEvents();
        }
        if (activeTab === 'fines' && teamOverview) {
            fetchFines();
        }
    }, [activeTab, teamOverview, fetchEvents, fetchFines]);

    const handleEdit = () => {
        if (teamOverview?.Team) {
            setEditForm({ 
                name: teamOverview.Team.Name, 
                description: teamOverview.Team.Description 
            });
            setIsEditing(true);
        }
    };

    const updateTeam = async () => {
        try {
            // OData v2: Update team using PATCH
            await api.patch(`/api/v2/Teams('${teamId}')`, {
                Name: editForm.name,
                Description: editForm.description
            });
            // Refresh team data
            // OData v2: Use GetOverview function
            const response = await api.get(`/api/v2/Teams('${teamId}')/GetOverview()`);
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
            // OData v2: Create event with TeamID
            await api.post(`/api/v2/Events`, {
                ...newEvent,
                TeamID: teamId,
                ClubID: clubId
            });
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
            // OData v2: Delete event
            await api.delete(`/api/v2/Events('${eventId}')`);
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
            // OData v2: Create fine with TeamID
            await api.post(`/api/v2/Fines`, {
                ...newFine,
                TeamID: teamId,
                ClubID: clubId
            });
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
            // OData v2: Delete fine
            await api.delete(`/api/v2/Fines('${fineId}')`);
            fetchFines();
        } catch {
            setError('Failed to delete fine');
        }
    };

    if (loading) return <div>Loading team admin...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!teamOverview) return <div>Team not found</div>;

    const { Team: team, Stats: stats } = teamOverview;

    return (
        <Layout title={`${team.Name} - Admin`}>
            <div className="admin-team-container">
                <Tabs
                    tabs={[
                        { id: 'overview', label: 'Overview' },
                        { id: 'events', label: 'Events' },
                        { id: 'fines', label: 'Fines' }
                    ]}
                    activeTab={activeTab}
                    onTabChange={(tabId) => setActiveTab(tabId as 'overview' | 'events' | 'fines')}
                >
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
                                        <Button onClick={updateTeam} variant="accept">Save</Button>
                                        <Button onClick={() => setIsEditing(false)} variant="cancel">Cancel</Button>
                                    </div>
                                </div>
                            ) : (
                                <>
                                    <div className="team-header">
                                        <div className="team-info">
                                            <h2>{team.Name}</h2>
                                            <p>{team.Description}</p>
                                        </div>
                                        <div className="team-actions">
                                            <Button onClick={handleEdit} variant="accept">Edit Team</Button>
                                            <Button 
                                                onClick={() => navigate(`/clubs/${clubId}/teams/${teamId}`)} 
                                                variant="secondary"
                                            >
                                                View Team
                                            </Button>
                                        </div>
                                    </div>

                                    {/* Team Stats */}
                                    <div className="team-stats-section">
                                        <h3>Team Statistics</h3>
                                        <div className="stats-grid">
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.MemberCount}</div>
                                                <div className="stat-label">Members</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.AdminCount}</div>
                                                <div className="stat-label">Admins</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.UpcomingEvents}</div>
                                                <div className="stat-label">Upcoming Events</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.TotalEvents}</div>
                                                <div className="stat-label">Total Events</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.UnpaidFines}</div>
                                                <div className="stat-label">Unpaid Fines</div>
                                            </div>
                                            <div className="stat-card">
                                                <div className="stat-number">{stats.TotalFines}</div>
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
                                <Button 
                                    onClick={() => setShowCreateEventModal(true)} 
                                    variant="accept"
                                >
                                    Create Event
                                </Button>
                            </div>

                            <div className="events-list">
                                {events.map(event => (
                                    <div key={event.id} className="event-card">
                                        <div className="event-header">
                                            <h4>{event.name}</h4>
                                            <Button 
                                                onClick={() => deleteEvent(event.id)}
                                                variant="cancel"
                                                size="sm"
                                            >
                                                Delete
                                            </Button>
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
                                <Button 
                                    onClick={() => setShowCreateFineModal(true)} 
                                    variant="accept"
                                >
                                    Create Fine
                                </Button>
                            </div>

                            <div className="fines-list">
                                {fines.map(fine => (
                                    <div key={fine.id} className="fine-card">
                                        <div className="fine-header">
                                            <h4>{fine.userName}</h4>
                                            <div className="fine-amount">${fine.amount}</div>
                                            <Button 
                                                onClick={() => deleteFine(fine.id)}
                                                variant="cancel"
                                                size="sm"
                                            >
                                                Delete
                                            </Button>
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
                </Tabs>
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
                    <Button onClick={createEvent} variant="accept">
                        Create Event
                    </Button>
                    <Button onClick={() => setShowCreateEventModal(false)} variant="cancel">
                        Cancel
                    </Button>
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
                    <Button onClick={createFine} variant="accept">
                        Create Fine
                    </Button>
                    <Button onClick={() => setShowCreateFineModal(false)} variant="cancel">
                        Cancel
                    </Button>
                </Modal.Actions>
            </Modal>
        </Layout>
    );
};

export default AdminTeamDetails;