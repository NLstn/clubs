import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import EditEvent from "./EditEvent";
import AddEvent from "./AddEvent";
import EventRSVPList from "./EventRSVPList";
import { Table, TableColumn } from '@/components/ui';
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

const AdminTeamEventList = () => {
    const { clubId, teamId } = useParams();
    const [events, setEvents] = useState<Event[]>([]);
    const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isAddModalOpen, setIsAddModalOpen] = useState(false);
    const [isRSVPModalOpen, setIsRSVPModalOpen] = useState(false);
    const [selectedEventForRSVP, setSelectedEventForRSVP] = useState<Event | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [rsvpCounts, setRsvpCounts] = useState<Record<string, RSVPCounts>>({});

    const fetchEvents = useCallback(async () => {
        if (!clubId || !teamId) return;
        setLoading(true);
        setError(null);
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/events`);
            setEvents(response.data || []);
            const counts: Record<string, RSVPCounts> = {};
            for (const event of response.data || []) {
                try {
                    const rsvpResponse = await api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/events/${event.id}/rsvps`);
                    counts[event.id] = rsvpResponse.data.counts || {};
                } catch (err) {
                    counts[event.id] = {};
                }
            }
            setRsvpCounts(counts);
        } catch (err) {
            console.error('Error fetching events:', err);
            setError('Failed to fetch events');
            setEvents([]);
        } finally {
            setLoading(false);
        }
    }, [clubId, teamId]);

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
        if (!clubId || !teamId) return;
        if (!confirm("Are you sure you want to delete this event? This will also delete all RSVPs.")) {
            return;
        }
        try {
            await api.delete(`/api/v1/clubs/${clubId}/teams/${teamId}/events/${eventId}`);
            fetchEvents();
        } catch (err) {
            console.error('Error deleting event:', err);
            setError('Failed to delete event');
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

    const eventColumns: TableColumn<Event>[] = [
        { key: 'name', header: 'Name', render: (event) => event.name },
        { key: 'start_time', header: 'Start', render: (event) => formatDateTime(event.start_time) },
        { key: 'end_time', header: 'End', render: (event) => formatDateTime(event.end_time) },
        {
            key: 'rsvps',
            header: 'RSVPs',
            render: (event) => {
                const counts = rsvpCounts[event.id] || {};
                const yesCount = counts.yes || 0;
                const noCount = counts.no || 0;
                const maybeCount = counts.maybe || 0;
                return (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                        <div>
                            <span style={{color: 'green'}}>Yes: {yesCount}</span>{' '}
                            <span style={{color: 'red'}}>No: {noCount}</span>{' '}
                            <span style={{color: 'orange'}}>Maybe: {maybeCount}</span>
                        </div>
                        {(yesCount > 0 || noCount > 0 || maybeCount > 0) && (
                            <button onClick={() => handleViewRSVPs(event)} className="button" style={{ fontSize: '0.8em', padding: '4px 10px' }}>View Details</button>
                        )}
                    </div>
                );
            }
        },
        {
            key: 'actions',
            header: 'Actions',
            render: (event) => (
                <div style={{ display: 'flex', gap: '8px' }}>
                    <button onClick={() => handleEditEvent(event)} className="button">Edit</button>
                    <button onClick={() => handleDeleteEvent(event.id)} className="button-cancel">Delete</button>
                </div>
            )
        }
    ];

    if (loading) return <div>Loading events...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div>
            <h2>Team Events</h2>
            <button onClick={() => setIsAddModalOpen(true)} className="button-accept" style={{ marginBottom: '10px' }}>
                Add Event
            </button>
            <Table columns={eventColumns} data={events} keyExtractor={(event) => event.id} emptyMessage="No events found." />

            <AddEvent
                isOpen={isAddModalOpen}
                onClose={() => setIsAddModalOpen(false)}
                clubId={clubId || ''}
                teamId={teamId || ''}
                onSuccess={fetchEvents}
            />
            <EditEvent
                isOpen={isEditModalOpen}
                onClose={handleCloseEditModal}
                clubId={clubId || ''}
                teamId={teamId || ''}
                event={selectedEvent}
                onSuccess={fetchEvents}
            />
            {selectedEventForRSVP && (
                <EventRSVPList
                    isOpen={isRSVPModalOpen}
                    onClose={handleCloseRSVPModal}
                    eventId={selectedEventForRSVP.id}
                    eventName={selectedEventForRSVP.name}
                    clubId={clubId || ''}
                    teamId={teamId || ''}
                />
            )}
        </div>
    );
};

export default AdminTeamEventList;
