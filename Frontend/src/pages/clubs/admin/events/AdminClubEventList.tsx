import { useEffect, useState, useCallback } from "react";
import { useParams, Link } from "react-router-dom";
import EditEvent from "./EditEvent";
import AddEvent from "./AddEvent";
import EventRSVPList from "./EventRSVPList";
import { Table, TableColumn, Button } from '@/components/ui';
import api from "../../../../utils/api";
import { parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';
import { calculateRSVPCounts, RSVPCounts } from "../../../../utils/eventUtils";

interface Event {
    id: string;
    name: string;
    description: string;
    location: string;
    start_time: string;
    end_time: string;
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
            // OData v2: Query Events for this club
            const response = await api.get(`/api/v2/Events?$filter=ClubID eq '${id}'`);
            setEvents(response.data || []);

            // Fetch RSVP counts and shifts for each event
            // TODO: Optimize to avoid N+1 query pattern - see issue #472
            // Consider using $expand=EventRSVPs in initial query or batch endpoint
            const counts: Record<string, RSVPCounts> = {};
            const shifts: Record<string, Shift[]> = {};
            for (const event of response.data || []) {
                try {
                    // Fetch RSVP counts
                    // OData v2: Use EventRSVPs navigation and compute counts client-side
                    interface EventRSVP { Response: string; }
                    const rsvpResponse = await api.get<ODataCollectionResponse<EventRSVP>>(`/api/v2/Events('${event.id}')/EventRSVPs`);
                    const rsvpList = parseODataCollection(rsvpResponse.data);
                    // Compute counts by grouping RSVPs by Response field
                    const computedCounts = calculateRSVPCounts(rsvpList);
                    counts[event.id] = computedCounts;
                    
                    // Fetch event shifts
                    // OData v2: Query Shifts for this event
                    const shiftsResponse = await api.get(`/api/v2/Shifts?$filter=EventID eq '${event.id}'`);
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
            // OData v2: Delete event
            await api.delete(`/api/v2/Events('${eventId}')`);
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

    // Define table columns
    const eventColumns: TableColumn<Event>[] = [
        {
            key: 'name',
            header: 'Name',
            render: (event) => event.name
        },
        {
            key: 'start_time',
            header: 'Start',
            render: (event) => formatDateTime(event.start_time)
        },
        {
            key: 'end_time',
            header: 'End',
            render: (event) => formatDateTime(event.end_time)
        },
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
                            <Button
                                onClick={() => handleViewRSVPs(event)}
                                variant="secondary"
                                size="sm"
                            >
                                View Details
                            </Button>
                        )}
                    </div>
                );
            }
        },
        {
            key: 'shifts',
            header: 'Shifts',
            render: (event) => {
                const shifts = eventShifts[event.id] || [];
                return (
                    <span style={{color: shifts.length > 0 ? 'blue' : 'gray'}}>
                        {shifts.length} shift{shifts.length !== 1 ? 's' : ''}
                    </span>
                );
            }
        },
        {
            key: 'actions',
            header: 'Actions',
            render: (event) => (
                <div style={{ display: 'flex', gap: '5px', flexWrap: 'wrap' }}>
                    <Link 
                        to={`/clubs/${id}/admin/events/${event.id}`}
                        style={{ textDecoration: 'none' }}
                    >
                        <Button variant="secondary" size="sm">
                            View Details
                        </Button>
                    </Link>
                    <Button
                        onClick={() => handleEditEvent(event)}
                        variant="accept"
                        size="sm"
                    >
                        Edit
                    </Button>
                    <Button
                        onClick={() => handleDeleteEvent(event.id)}
                        variant="cancel"
                        size="sm"
                    >
                        Delete
                    </Button>
                </div>
            )
        }
    ];

    return (
        <div>
            <h3>Events</h3>
            <Table
                columns={eventColumns}
                data={events}
                keyExtractor={(event) => event.id}
                loading={loading}
                error={error}
                emptyMessage="No events available"
                loadingMessage="Loading events..."
                errorMessage={error || "Failed to fetch events"}
            />
            <div style={{ marginBottom: '20px' }}>
                <Button onClick={() => setIsAddModalOpen(true)} variant="accept">
                    Add Event
                </Button>
            </div>
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