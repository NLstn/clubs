import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import Layout from '../components/layout/Layout';

interface Club {
    id: string;
    name: string;
    description: string;
}

interface NewsItem {
    id: string;
    club_id: string;
    title: string;
    content: string;
    created_at: string;
    club: Club;
}

interface EventItem {
    id: string;
    club_id: string;
    name: string;
    start_time: string;
    end_time: string;
    created_at: string;
    club: Club;
}

const Dashboard = () => {
    const navigate = useNavigate();
    const { api } = useAuth();
    const [clubs, setClubs] = useState<Club[]>([]);
    const [news, setNews] = useState<NewsItem[]>([]);
    const [events, setEvents] = useState<EventItem[]>([]);
    const [message, setMessage] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchDashboardData = async () => {
            try {
                setLoading(true);
                
                // Fetch clubs
                const clubsResponse = await api.get('/api/v1/clubs');
                setClubs(clubsResponse.data);

                // Fetch unified news
                const newsResponse = await api.get('/api/v1/dashboard/news');
                setNews(newsResponse.data || []);

                // Fetch unified events
                const eventsResponse = await api.get('/api/v1/dashboard/events');
                setEvents(eventsResponse.data || []);
                
            } catch (error) {
                setMessage('Error fetching dashboard data');
                console.error(error);
            } finally {
                setLoading(false);
            }
        };

        fetchDashboardData();
    }, [api]);

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    return (
        <Layout title="Dashboard">
            <div>
                {message && <p className={`message ${message.includes('Error') ? 'error' : 'success'}`}>
                    {message}
                </p>}

                {loading && <p>Loading dashboard...</p>}

                {!loading && (
                    <>
                        {/* Clubs Overview */}
                        <section className="dashboard-section">
                            <h2>Your Clubs</h2>
                            <div className="clubs-list">
                                {clubs === null || clubs.length === 0 ? (
                                    <p>No clubs available. Create one to get started!</p>
                                ) : (
                                    clubs.map(club => (
                                        <div
                                            key={club.id}
                                            className="card card-clickable"
                                            onClick={() => navigate(`/clubs/${club.id}`)}
                                        >
                                            <h4>{club.name}</h4>
                                            <p>{club.description}</p>
                                        </div>
                                    ))
                                )}
                            </div>
                        </section>

                        {/* Unified News */}
                        <section className="dashboard-section">
                            <h2>Recent News</h2>
                            <div className="news-list">
                                {news.length === 0 ? (
                                    <p>No recent news from your clubs.</p>
                                ) : (
                                    news.slice(0, 5).map(newsItem => (
                                        <div 
                                            key={newsItem.id} 
                                            className="card card-clickable news-item"
                                            onClick={() => navigate(`/clubs/${newsItem.club.id}`)}
                                        >
                                            <div className="news-header">
                                                <h4>{newsItem.title}</h4>
                                                <span className="club-badge">{newsItem.club.name}</span>
                                            </div>
                                            <p className="news-content">{newsItem.content}</p>
                                            <div className="news-meta">
                                                <span className="date">{formatDate(newsItem.created_at)}</span>
                                            </div>
                                        </div>
                                    ))
                                )}
                            </div>
                        </section>

                        {/* Upcoming Events */}
                        <section className="dashboard-section">
                            <h2>Upcoming Events</h2>
                            <div className="events-list">
                                {events.length === 0 ? (
                                    <p>No upcoming events from your clubs.</p>
                                ) : (
                                    events.slice(0, 5).map(event => (
                                        <div 
                                            key={event.id} 
                                            className="card card-clickable event-item"
                                            onClick={() => navigate(`/clubs/${event.club.id}`)}
                                        >
                                            <div className="event-header">
                                                <h4>{event.name}</h4>
                                                <span className="club-badge">{event.club.name}</span>
                                            </div>
                                            <div className="event-meta">
                                                <span className="date">
                                                    {formatDate(event.start_time)} - {formatDate(event.end_time)}
                                                </span>
                                            </div>
                                        </div>
                                    ))
                                )}
                            </div>
                        </section>
                    </>
                )}
            </div>
        </Layout>
    );
};

export default Dashboard;