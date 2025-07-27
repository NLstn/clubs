import { useEffect, useState, useCallback } from "react";
import { useParams, Link } from "react-router-dom";
import EditEvent from "./EditEvent";
import AddEvent from "./AddEvent";
import EventRSVPList from "./EventRSVPList";
import api from "../../../../utils/api";

interface Event {
    id: string;
    name: string;
    description: string;
    location: string;
    start_time: string;
    end_time: string;
}

interface RSVPCounts {
    yes?: number;
    no?: number;
    maybe?: number;
}

interface Shift {
    id: string;
    startTime: string;
    endTime: string;
}

const AdminClubEventList = () => {
    const { id } = useParams();
    const [events, setEvents] = useState<Event[]>([]);
    const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isAddModalOpen, setIsAddModalOpen] = useState(false);
    const [isRSVPModalOpen, setIsRSVPModalOpen] = useState(false);
    const [selectedEventForRSVP, setSelectedEventForRSVP] = useState<Event | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [rsvpCounts, setRsvpCounts] = useState<Record<string, RSVPCounts>>({});
    const [eventShifts, setEventShifts] = useState<Record<string, Shift[]>>({});

    const fetchEvents = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await api.get(`/api/v1/clubs/${id}/events`);
            setEvents(response.data || []);

            // Fetch RSVP counts and shifts for each event
            const counts: Record<string, RSVPCounts> = {};
            const shifts: Record<string, Shift[]> = {};
            for (const event of response.data || []) {
                try {
                    // Fetch RSVP counts
                    const rsvpResponse = await api.get(`/api/v1/clubs/${id}/events/${event.id}/rsvps`);
                    counts[event.id] = rsvpResponse.data.counts || {};
                    
                    // Fetch event shifts
                    const shiftsResponse = await api.get(`/api/v1/clubs/${id}/events/${event.id}/shifts`);
                    shifts[event.id] = shiftsResponse.data || [];
                } catch (err) {
                    console.warn(`Failed to fetch data for event ${event.id}:`, err);
                    counts[event.id] = {};
                    shifts[event.id] = [];
                }
            }
            setRsvpCounts(counts);
            setEventShifts(shifts);
        } catch (error) {
            console.error("Error fetching events:", error);
            setError(error instanceof Error ? error.message : "Failed to fetch events");
            setEvents([]);
        } finally {
            setLoading(false);
        }
    }, [id]);

    useEffect(() => {
        fetchEvents();
    }, [fetchEvents]);

    const handleEditEvent = (event: Event) => {
        setSelectedEvent(event);
        setIsEditModalOpen(true);
    };

    const handleCloseEditModal = () => {
        setSelectedEvent(null);
        setIsEditModalOpen(false);
    };

    const handleViewRSVPs = (event: Event) => {
        setSelectedEventForRSVP(event);
        setIsRSVPModalOpen(true);
    };

    const handleCloseRSVPModal = () => {
        setSelectedEventForRSVP(null);
        setIsRSVPModalOpen(false);
    };

    const handleDeleteEvent = async (eventId: string) => {
        if (!confirm("Are you sure you want to delete this event? This will also delete all RSVPs.")) {
            return;
        }

        try {
            await api.delete(`/api/v1/clubs/${id}/events/${eventId}`);
            fetchEvents(); // Refresh the list
        } catch (error) {
            console.error("Error deleting event:", error);
            setError(error instanceof Error ? error.message : "Failed to delete event");
        }
    };

    const formatDateTime = (timestamp: string) => {
        // Handle undefined, null, or empty values
        if (!timestamp || timestamp === 'null') {
            return 'Not set';
        }
        
        try {
            const dateTime = new Date(timestamp);
            // Check if the date is valid
            if (isNaN(dateTime.getTime())) {
                return 'Invalid date/time';
            }
            return dateTime.toLocaleString();
        } catch {
            return 'Parse error';
        }
    };

    return (
        <div>
            <h3>Events</h3>
            {loading && <p>Loading events...</p>}
            {error && <p style={{color: 'red'}}>Error: {error}</p>}
            {!loading && !error && (
                <>
                    <table className="basic-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Start</th>
                                <th>End</th>
                                <th>RSVPs</th>
                                <th>Shifts</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {events.length === 0 ? (
                                <tr>
                                    <td colSpan={6} style={{textAlign: 'center', fontStyle: 'italic'}}>
                                        No events available
                                    </td>
                                </tr>
                            ) : (
                                events.map(event => {
                                    const counts = rsvpCounts[event.id] || {};
                                    const shifts = eventShifts[event.id] || [];
                                    const yesCount = counts.yes || 0;
                                    const noCount = counts.no || 0;
                                    const maybeCount = counts.maybe || 0;
                                    
                                    return (
                                        <tr key={event.id}>
                                            <td>{event.name}</td>
                                            <td>{formatDateTime(event.start_time)}</td>
                                            <td>{formatDateTime(event.end_time)}</td>
                                            <td>
                                                <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                                                    <div>
                                                        <span style={{color: 'green'}}>Yes: {yesCount}</span>{' '}
                                                        <span style={{color: 'red'}}>No: {noCount}</span>{' '}
                                                        <span style={{color: 'orange'}}>Maybe: {maybeCount}</span>
                                                    </div>
                                                    {(yesCount > 0 || noCount > 0 || maybeCount > 0) && (
                                                        <button
                                                            onClick={() => handleViewRSVPs(event)}
                                                            className="button"
                                                            style={{ 
                                                                fontSize: '0.8em', 
                                                                padding: '4px 10px',
                                                                backgroundColor: '#007bff',
                                                                color: 'white',
                                                                border: '1px solid #007bff',
                                                                borderRadius: '4px',
                                                                cursor: 'pointer',
                                                                fontWeight: '500'
                                                            }}
                                                        >
                                                            View Details
                                                        </button>
                                                    )}
                                                </div>
                                            </td>
                                            <td>
                                                <span style={{color: shifts.length > 0 ? 'blue' : 'gray'}}>
                                                    {shifts.length} shift{shifts.length !== 1 ? 's' : ''}
                                                </span>
                                            </td>
                                            <td>
                                                <Link 
                                                    to={`/clubs/${id}/admin/events/${event.id}`}
                                                    className="button-info"
                                                    style={{marginRight: '5px'}}
                                                >
                                                    View Details
                                                </Link>
                                                <button
                                                    onClick={() => handleEditEvent(event)}
                                                    className="button-accept"
                                                    style={{marginRight: '5px'}}
                                                >
                                                    Edit
                                                </button>
                                                <button
                                                    onClick={() => handleDeleteEvent(event.id)}
                                                    className="button-cancel"
                                                >
                                                    Delete
                                                </button>
                                            </td>
                                        </tr>
                                    );
                                })
                            )}
                        </tbody>
                    </table>
                    <div style={{ marginTop: '20px' }}>
                        <button onClick={() => setIsAddModalOpen(true)} className="button-accept">
                            Add Event
                        </button>
                    </div>
                </>
            )}
            <EditEvent
                isOpen={isEditModalOpen}
                onClose={handleCloseEditModal}
                event={selectedEvent}
                clubId={id}
                onSuccess={fetchEvents}
            />
            <AddEvent 
                isOpen={isAddModalOpen}
                onClose={() => setIsAddModalOpen(false)}
                clubId={id || ''}
                onSuccess={fetchEvents}
            />
            <EventRSVPList
                isOpen={isRSVPModalOpen}
                onClose={handleCloseRSVPModal}
                eventId={selectedEventForRSVP?.id || ''}
                eventName={selectedEventForRSVP?.name || ''}
                clubId={id || ''}
            />
        </div>
    );
};

export default AdminClubEventList;