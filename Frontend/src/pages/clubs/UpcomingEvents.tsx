import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import api from '../../utils/api';

interface Event {
    id: string;
    name: string;
    description: string;
    location: string;
    start_time: string;
    end_time: string;
    user_rsvp?: {
        response: string;
    };
}

const UpcomingEvents = () => {
    const { id } = useParams();
    const [events, setEvents] = useState<Event[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchUpcomingEvents = async () => {
        if (!id) return;
        
        try {
            const response = await api.get(`/api/v1/clubs/${id}/events/upcoming`);
            setEvents(response.data || []);
        } catch (error) {
            console.error("Error fetching upcoming events:", error);
            setError("Failed to load upcoming events");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUpcomingEvents();
    }, [id]); // eslint-disable-line react-hooks/exhaustive-deps

    const handleRSVP = async (eventId: string, response: string) => {
        try {
            await api.post(`/api/v1/clubs/${id}/events/${eventId}/rsvp`, { response });
            // Refresh events to update RSVP status
            fetchUpcomingEvents();
        } catch (error) {
            console.error("Error updating RSVP:", error);
            alert("Failed to update RSVP. Please try again.");
        }
    };

    const formatDateTime = (timestamp: string) => {
        try {
            const dateTime = new Date(timestamp);
            return dateTime.toLocaleString();
        } catch {
            return timestamp;
        }
    };

    if (loading) return <div>Loading upcoming events...</div>;
    if (error) return <div className="error">{error}</div>;
    if (events.length === 0) return <div>No upcoming events.</div>;

    return (
        <div>
            <h3>Upcoming Events</h3>
            <div className="events-list">
                {events.map(event => (
                    <div key={event.id} className="event-card">
                        <div className="event-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
                            <h4 style={{ margin: 0 }}>{event.name}</h4>
                            <Link 
                                to={`/clubs/${id}/events/${event.id}`} 
                                className="button-info"
                            >
                                View Details
                            </Link>
                        </div>
                        <p>
                            <strong>Start:</strong> {formatDateTime(event.start_time)}
                        </p>
                        <p>
                            <strong>End:</strong> {formatDateTime(event.end_time)}
                        </p>
                        {event.location && (
                            <p>
                                <strong>Location:</strong> {event.location}
                            </p>
                        )}
                        
                        <div className="rsvp-section">
                            <p>
                                <strong>RSVP:</strong> 
                                {event.user_rsvp ? (
                                    <span className={`rsvp-status ${event.user_rsvp.response}`}>
                                        {event.user_rsvp.response === 'yes' ? 'Yes' : 
                                         event.user_rsvp.response === 'no' ? 'No' : 'Maybe'}
                                    </span>
                                ) : (
                                    <span className="rsvp-status none">No response</span>
                                )}
                            </p>
                            <div className="rsvp-buttons">
                                <button
                                    onClick={() => handleRSVP(event.id, 'yes')}
                                    className={event.user_rsvp?.response === 'yes' ? 'button-accept' : 'button'}
                                >
                                    Yes
                                </button>
                                <button
                                    onClick={() => handleRSVP(event.id, 'maybe')}
                                    className={event.user_rsvp?.response === 'maybe' ? 'button-maybe' : 'button'}
                                >
                                    Maybe
                                </button>
                                <button
                                    onClick={() => handleRSVP(event.id, 'no')}
                                    className={event.user_rsvp?.response === 'no' ? 'button-cancel' : 'button'}
                                >
                                    No
                                </button>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default UpcomingEvents;