import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../utils/api';

interface Event {
    id: string;
    name: string;
    start_date: string;
    start_time: string;
    end_date: string;
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

    const formatDateTime = (date: string, time: string) => {
        try {
            const dateTime = new Date(`${date}T${time}`);
            return dateTime.toLocaleString();
        } catch {
            return `${date} ${time}`;
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
                    <div key={event.id} className="event-card" style={{
                        border: '1px solid #ddd',
                        borderRadius: '8px',
                        padding: '16px',
                        marginBottom: '12px',
                        backgroundColor: '#f9f9f9'
                    }}>
                        <h4>{event.name}</h4>
                        <p>
                            <strong>Start:</strong> {formatDateTime(event.start_date, event.start_time)}
                        </p>
                        <p>
                            <strong>End:</strong> {formatDateTime(event.end_date, event.end_time)}
                        </p>
                        
                        <div className="rsvp-section" style={{ marginTop: '12px' }}>
                            <p>
                                <strong>RSVP:</strong> 
                                {event.user_rsvp ? (
                                    <span style={{ 
                                        marginLeft: '8px',
                                        color: event.user_rsvp.response === 'yes' ? 'green' : 'red',
                                        fontWeight: 'bold'
                                    }}>
                                        {event.user_rsvp.response === 'yes' ? 'Yes' : 'No'}
                                    </span>
                                ) : (
                                    <span style={{ marginLeft: '8px', color: '#666' }}>No response</span>
                                )}
                            </p>
                            <div className="rsvp-buttons" style={{ marginTop: '8px' }}>
                                <button
                                    onClick={() => handleRSVP(event.id, 'yes')}
                                    className={event.user_rsvp?.response === 'yes' ? 'button-accept' : 'button'}
                                    style={{ marginRight: '8px' }}
                                >
                                    Yes
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