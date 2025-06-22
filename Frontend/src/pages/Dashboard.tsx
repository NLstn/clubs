import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { useDashboardData } from '../hooks/useDashboardData';
import Layout from '../components/layout/Layout';

interface Club {
    id: number;
    name: string;
    description: string;
}

const Dashboard = () => {
    const navigate = useNavigate();
    const { api } = useAuth();
    const [clubs, setClubs] = useState<Club[]>([]);
    const [message, setMessage] = useState('');
    const { news, events, activities, loading: dashboardLoading, error: dashboardError } = useDashboardData();

    useEffect(() => {
        const fetchClubs = async () => {
            try {
                const response = await api.get('/api/v1/clubs');
                setClubs(response.data);
            } catch (error) {
                setMessage('Error fetching clubs');
                console.error(error);
            }
        };

        fetchClubs();
    }, [api]);

    const formatDateTime = (timestamp: string) => {
        try {
            const dateTime = new Date(timestamp);
            return dateTime.toLocaleDateString();
        } catch {
            return timestamp;
        }
    };

    const formatEventDateTime = (timestamp: string) => {
        try {
            const dateTime = new Date(timestamp);
            return dateTime.toLocaleString();
        } catch {
            return timestamp;
        }
    };

    return (
        <Layout title="Dashboard">
            <div>
                {message && <p className={`message ${message.includes('Error') ? 'error' : 'success'}`}>
                    {message}
                </p>}
                {dashboardError && <p className="error">{dashboardError}</p>}

                {/* Activity Feed Section */}
                {dashboardLoading ? (
                    <div>Loading dashboard...</div>
                ) : (
                    <>
                        <div className="dashboard-section">
                            <h2>Activity Feed</h2>
                            {activities.length > 0 ? (
                                <div className="activity-feed">
                                    {activities.map(activity => (
                                        <div key={`${activity.type}-${activity.id}`} className="activity-item">
                                            <div className="activity-header">
                                                <div className="activity-type-badge">{activity.type}</div>
                                                <span 
                                                    className="club-badge"
                                                    onClick={() => navigate(`/clubs/${activity.club_id}`)}
                                                >
                                                    {activity.club_name}
                                                </span>
                                            </div>
                                            <h4 className="activity-title">{activity.title}</h4>
                                            {activity.content && (
                                                <p className="activity-content">{activity.content}</p>
                                            )}
                                            {activity.type === 'event' && activity.metadata?.start_time && (
                                                <p className="activity-event-details">
                                                    <strong>Start:</strong> {formatEventDateTime(activity.metadata.start_time)}
                                                    {activity.metadata?.end_time && (
                                                        <>
                                                            <br />
                                                            <strong>End:</strong> {formatEventDateTime(activity.metadata.end_time)}
                                                        </>
                                                    )}
                                                    {activity.metadata?.user_rsvp && (
                                                        <>
                                                            <br />
                                                            <strong>RSVP:</strong> 
                                                            <span className={`rsvp-status ${activity.metadata.user_rsvp.response}`}>
                                                                {activity.metadata.user_rsvp.response === 'yes' ? 'Yes' : 'No'}
                                                            </span>
                                                        </>
                                                    )}
                                                </p>
                                            )}
                                            <small className="activity-meta">
                                                Posted on {formatDateTime(activity.created_at)}
                                            </small>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <div className="empty-state">
                                    <p>No recent activities from your clubs.</p>
                                </div>
                            )}
                        </div>

                        {/* Legacy News Section - keeping for backward compatibility */}
                        <div className="dashboard-section">
                            <h2>Latest News</h2>
                            {news.length > 0 ? (
                                <div className="dashboard-news">
                                    {news.slice(0, 5).map(newsItem => (
                                        <div key={newsItem.id} className="dashboard-news-card">
                                            <div className="news-header">
                                                <h4 className="news-title">{newsItem.title}</h4>
                                                <span 
                                                    className="club-badge"
                                                    onClick={() => navigate(`/clubs/${newsItem.club_id}`)}
                                                >
                                                    {newsItem.club_name}
                                                </span>
                                            </div>
                                            <p className="news-content">{newsItem.content}</p>
                                            <small className="news-meta">
                                                Posted on {formatDateTime(newsItem.created_at)}
                                            </small>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <div className="empty-state">
                                    <p>No recent news from your clubs.</p>
                                </div>
                            )}
                        </div>

                        {/* Legacy Events Section - keeping for backward compatibility */}
                        <div className="dashboard-section">
                            <h2>Upcoming Events</h2>
                            {events.length > 0 ? (
                                <div className="dashboard-events">
                                    {events.slice(0, 5).map(event => (
                                        <div key={event.id} className="dashboard-event-card">
                                            <div className="event-header">
                                                <h4 className="event-title">{event.name}</h4>
                                                <span 
                                                    className="club-badge"
                                                    onClick={() => navigate(`/clubs/${event.club_id}`)}
                                                >
                                                    {event.club_name}
                                                </span>
                                            </div>
                                            <p>
                                                <strong>Start:</strong> {formatEventDateTime(event.start_time)}
                                            </p>
                                            <p>
                                                <strong>End:</strong> {formatEventDateTime(event.end_time)}
                                            </p>
                                            {event.user_rsvp && (
                                                <p>
                                                    <strong>RSVP:</strong> 
                                                    <span className={`rsvp-status ${event.user_rsvp.response}`}>
                                                        {event.user_rsvp.response === 'yes' ? 'Yes' : 'No'}
                                                    </span>
                                                </p>
                                            )}
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <div className="empty-state">
                                    <p>No upcoming events from your clubs.</p>
                                </div>
                            )}
                        </div>
                    </>
                )}

                {/* Clubs Section */}
                <div className="dashboard-section">
                    <h2>Your Clubs</h2>
                    <div className="clubs-list">
                        {clubs === null || clubs.length === 0 ? (
                            <p>No clubs available. Create one to get started!</p>
                        ) : (
                            clubs.map(club => (
                                <div
                                    key={club.id}
                                    className="club-card"
                                    onClick={() => navigate(`/clubs/${club.id}`)}
                                >
                                    <div className="club-card-header">
                                        <h4>{club.name}</h4>
                                    </div>
                                    <p className="club-description">{club.description}</p>
                                    <div className="club-card-footer">
                                        <span className="club-link">View Club →</span>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </div>
            </div>
        </Layout>
    );
};

export default Dashboard;