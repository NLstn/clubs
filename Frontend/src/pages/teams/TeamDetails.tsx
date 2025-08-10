import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../../components/layout/Layout';
import { useT } from '../../hooks/useTranslation';
import './TeamDetails.css';

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

const TeamDetails = () => {
    const { t } = useT();
    const { clubId, teamId } = useParams();
    const navigate = useNavigate();
    const [teamOverview, setTeamOverview] = useState<TeamOverview | null>(null);
    const [upcomingEvents, setUpcomingEvents] = useState<Event[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchTeamData = async () => {
            if (!clubId || !teamId) {
                setError('Missing club or team ID');
                setLoading(false);
                return;
            }

            try {
                const [overviewResponse, eventsResponse] = await Promise.all([
                    api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/overview`),
                    api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/events/upcoming`)
                ]);

                setTeamOverview(overviewResponse.data);
                setUpcomingEvents(eventsResponse.data || []);
                setError(null);
            } catch (err: unknown) {
                console.error('Error fetching team details:', err);
                if (err && typeof err === 'object' && 'response' in err) {
                    const axiosError = err as { response?: { status?: number } };
                    if (axiosError.response?.status === 404) {
                        setError('Team not found');
                    } else if (axiosError.response?.status === 403) {
                        setError('You do not have access to this team');
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
    }, [clubId, teamId]);

    if (loading) return <div>Loading team details...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!teamOverview) return <div>Team not found</div>;

    const { team, stats, user_role, is_admin } = teamOverview;

    return (
        <Layout title={team.name}>
            <div className="team-details-container">
                {/* Team Header */}
                <div className="team-header-section">
                    <div className="team-header-content">
                        <div className="team-main-info">
                            <h1 className="team-title">{team.name}</h1>
                            {team.description && (
                                <p className="team-description">{team.description}</p>
                            )}
                            {user_role && (
                                <div className="user-role-container">
                                    <span className="role-label">Your role</span>
                                    <div className={`role-badge role-${user_role}`}>
                                        <span className="role-text">{user_role}</span>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                    
                    {/* Action Buttons */}
                    <div className="team-actions">
                        <button 
                            className="button button-secondary"
                            onClick={() => navigate(`/clubs/${clubId}`)}
                        >
                            Back to Club
                        </button>
                        {is_admin && (
                            <button 
                                className="button button-primary"
                                onClick={() => navigate(`/clubs/${clubId}/teams/${teamId}/admin`)}
                            >
                                Manage Team
                            </button>
                        )}
                    </div>
                </div>

                {/* Team Stats */}
                <div className="team-stats-section">
                    <h3>Team Overview</h3>
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

                {/* Upcoming Events Section */}
                {upcomingEvents.length > 0 && (
                    <div className="team-content-section">
                        <h3>Upcoming Events</h3>
                        <div className="events-list">
                            {upcomingEvents.map(event => (
                                <div key={event.id} className="event-card">
                                    <div className="event-header">
                                        <h4 className="event-name">{event.name}</h4>
                                        <div className="event-time">
                                            {new Date(event.start_time).toLocaleDateString()} at{' '}
                                            {new Date(event.start_time).toLocaleTimeString([], { 
                                                hour: '2-digit', 
                                                minute: '2-digit' 
                                            })}
                                        </div>
                                    </div>
                                    {event.description && (
                                        <p className="event-description">{event.description}</p>
                                    )}
                                    {event.location && (
                                        <p className="event-location">📍 {event.location}</p>
                                    )}
                                    {is_admin && (
                                        <div className="event-actions">
                                            <button 
                                                className="button button-sm button-secondary"
                                                onClick={() => navigate(`/clubs/${clubId}/teams/${teamId}/admin/events`)}
                                            >
                                                Manage Events
                                            </button>
                                        </div>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {/* No Content Messages */}
                {upcomingEvents.length === 0 && (
                    <div className="team-content-section">
                        <div className="no-content-message">
                            <h3>No upcoming events</h3>
                            {is_admin && (
                                <p>
                                    <button 
                                        className="button button-primary"
                                        onClick={() => navigate(`/clubs/${clubId}/teams/${teamId}/admin/events`)}
                                    >
                                        Create First Event
                                    </button>
                                </p>
                            )}
                        </div>
                    </div>
                )}
            </div>
        </Layout>
    );
};

export default TeamDetails;