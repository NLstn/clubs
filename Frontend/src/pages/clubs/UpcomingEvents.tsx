import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Button } from '../../components/ui';
import api from '../../utils/api';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import { buildODataQuery, odataExpandWithOptions, ODataFilter, parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';
import '../../styles/events.css';

interface ODataEvent { 
    ID: string; 
    Name: string; 
    Description: string; 
    Location: string; 
    StartTime: string; 
    EndTime: string; 
    EventRSVPs?: Array<{ Response: string; }>; 
}

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
    const { user: currentUser } = useCurrentUser();
    const [events, setEvents] = useState<Event[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchUpcomingEvents = async () => {
        if (!id || !currentUser?.ID) return;
        
        try {
            // OData v2: Use navigation property with filter for upcoming events and nested expand for user's RSVPs
            const now = new Date().toISOString();
            const encodedId = encodeURIComponent(id);
            const encodedUserId = encodeURIComponent(currentUser.ID);
            
            const query = buildODataQuery({
                filter: ODataFilter.ge('EndTime', now),
                orderby: 'StartTime',
                expand: odataExpandWithOptions('EventRSVPs', {
                    filter: ODataFilter.eq('UserID', encodedUserId)
                })
            });
            
            const response = await api.get<ODataCollectionResponse<ODataEvent>>(`/api/v2/Clubs('${encodedId}')/Events${query}`);
            const eventsData = parseODataCollection(response.data);
            // Map OData response to match expected format
            const mappedEvents = eventsData.map((event: ODataEvent) => ({
                id: event.ID,
                name: event.Name,
                description: event.Description,
                location: event.Location,
                start_time: event.StartTime,
                end_time: event.EndTime,
                user_rsvp: event.EventRSVPs && event.EventRSVPs.length > 0 ? { response: event.EventRSVPs[0].Response } : undefined
            }));
            setEvents(mappedEvents);
        } catch (error) {
            console.error("Error fetching upcoming events:", error);
            setError("Failed to load upcoming events");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUpcomingEvents();
    }, [id, currentUser?.ID]); // eslint-disable-line react-hooks/exhaustive-deps

    const handleRSVP = async (eventId: string, response: string) => {
        try {
            // OData v2: Use AddRSVP action on Event entity
            await api.post(`/api/v2/Events('${eventId}')/AddRSVP`, { response });
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
                                style={{ textDecoration: 'none' }}
                            >
                                <Button variant="secondary" size="sm">
                                    View Details
                                </Button>
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
                                <Button
                                    size="sm"
                                    variant={event.user_rsvp?.response === 'yes' ? 'accept' : 'primary'}
                                    onClick={() => handleRSVP(event.id, 'yes')}
                                >
                                    Yes
                                </Button>
                                <Button
                                    size="sm"
                                    variant={event.user_rsvp?.response === 'maybe' ? 'maybe' : 'primary'}
                                    onClick={() => handleRSVP(event.id, 'maybe')}
                                >
                                    Maybe
                                </Button>
                                <Button
                                    size="sm"
                                    variant={event.user_rsvp?.response === 'no' ? 'cancel' : 'primary'}
                                    onClick={() => handleRSVP(event.id, 'no')}
                                >
                                    No
                                </Button>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default UpcomingEvents;