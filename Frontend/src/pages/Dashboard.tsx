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
    const { news, events, loading: dashboardLoading, error: dashboardError } = useDashboardData();

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

                {/* News Section */}
                {dashboardLoading ? (
                    <div>Loading dashboard...</div>
                ) : (
                    <>
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

                        {/* Events Section */}
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