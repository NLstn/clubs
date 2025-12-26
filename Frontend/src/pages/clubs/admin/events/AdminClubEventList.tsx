import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import EditEvent from "./EditEvent";
import AddEvent from "./AddEvent";
import EventRSVPList from "./EventRSVPList";
import { Table, TableColumn, Button, Modal, Card } from '@/components/ui';
import api from "../../../../utils/api";
import { parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';
import { RSVPCounts } from "../../../../utils/eventUtils";
import './AdminClubEventList.css';

interface EventRSVP {
    Response: string;
}

interface Shift {
    ID: string;
    StartTime: string;
    EndTime: string;
}

interface Event {
    ID: string;
    Name: string;
    Description: string;
    Location: string;
    StartTime: string;
    EndTime: string;
    EventRSVPs?: EventRSVP[];
    Shifts?: Shift[];
}

const AdminClubEventList = () => {
    const { id } = useParams();
    const [events, setEvents] = useState<Event[]>([]);
    const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isAddModalOpen, setIsAddModalOpen] = useState(false);
    const [isRSVPModalOpen, setIsRSVPModalOpen] = useState(false);
    const [isDetailsModalOpen, setIsDetailsModalOpen] = useState(false);
    const [selectedEventForDetails, setSelectedEventForDetails] = useState<Event | null>(null);
    const [selectedEventForRSVP, setSelectedEventForRSVP] = useState<Event | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [deleteLoading, setDeleteLoading] = useState(false);
    const [rsvpCounts, setRsvpCounts] = useState<Record<string, RSVPCounts>>({});
    const [eventShifts, setEventShifts] = useState<Record<string, Shift[]>>({});

    const fetchEvents = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            // OData v2: Query Events for this club with expanded Shifts
            // We no longer need to expand EventRSVPs since we'll use GetRSVPCounts function
            const response = await api.get<ODataCollectionResponse<Event>>(`/api/v2/Events?$filter=ClubID eq '${id}'&$expand=Shifts`);
            const eventList = parseODataCollection(response.data);
            setEvents(eventList);

            // Fetch RSVP counts using the server-side function (eliminates N+1 query pattern)
            const counts: Record<string, RSVPCounts> = {};
            const shifts: Record<string, Shift[]> = {};
            
            // Fetch RSVP counts for all events in parallel
            const countPromises = eventList.map(event =>
                api.get<{ Yes: number; No: number; Maybe: number }>(
                    `/api/v2/Events('${event.ID}')/GetRSVPCounts()`
                ).then(res => ({ eventId: event.ID, counts: res.data, status: 'fulfilled' as const }))
                 .catch(() => ({ eventId: event.ID, status: 'rejected' as const }))
            );
            
            const countResults = await Promise.all(countPromises);
            
            for (const event of eventList) {
                // Transform PascalCase response to camelCase for frontend
                const result = countResults.find(
                    (r): r is { eventId: string; counts: { Yes: number; No: number; Maybe: number }; status: 'fulfilled' } =>
                        r.status === 'fulfilled' && r.eventId === event.ID
                );
                if (result) {
                    counts[event.ID] = {
                        yes: result.counts.Yes,
                        no: result.counts.No,
                        maybe: result.counts.Maybe
                    };
                } else {
                    counts[event.ID] = { yes: 0, no: 0, maybe: 0 };
                }
                
                // Store shifts from the expanded Shifts
                shifts[event.ID] = event.Shifts || [];
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

    const handleViewDetails = (event: Event) => {
        setSelectedEventForDetails(event);
        setIsDetailsModalOpen(true);
    };

    const handleCloseDetailsModal = () => {
        setSelectedEventForDetails(null);
        setIsDetailsModalOpen(false);
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

        setDeleteLoading(true);
        try {
            // OData v2: Delete event
            await api.delete(`/api/v2/Events('${eventId}')`);
            handleCloseDetailsModal(); // Close the details modal after deletion
            fetchEvents(); // Refresh the list
        } catch (error) {
            console.error("Error deleting event:", error);
            setError(error instanceof Error ? error.message : "Failed to delete event");
        } finally {
            setDeleteLoading(false);
        }
    };

    const handleEditFromDetails = () => {
        if (selectedEventForDetails) {
            handleEditEvent(selectedEventForDetails);
            handleCloseDetailsModal();
        }
    };

    const handleDeleteFromDetails = () => {
        if (selectedEventForDetails) {
            handleDeleteEvent(selectedEventForDetails.ID);
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
            render: (event) => (
                <button
                    className="event-name-link"
                    onClick={() => handleViewDetails(event)}
                >
                    {event.Name}
                </button>
            )
        },
        {
            key: 'start_time',
            header: 'Start',
            render: (event) => formatDateTime(event.StartTime)
        },
        {
            key: 'end_time',
            header: 'End',
            render: (event) => formatDateTime(event.EndTime)
        },
        {
            key: 'rsvps',
            header: 'RSVPs',
            render: (event) => {
                const counts = rsvpCounts[event.ID] || {};
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
        }
    ];

    return (
        <div>
            <h3>Events</h3>
            <Table
                columns={eventColumns}
                data={events}
                keyExtractor={(event) => event.ID}
                loading={loading}
                error={error}
                emptyMessage="No events available"
                loadingMessage="Loading events..."
                errorMessage={error || "Failed to fetch events"}
            />
            <div style={{ marginTop: '20px' }}>
                <Button onClick={() => setIsAddModalOpen(true)} variant="accept">
                    Add Event
                </Button>
            </div>
            <EditEvent
                isOpen={isEditModalOpen}
                onClose={handleCloseEditModal}
                event={selectedEvent ? {
                    id: selectedEvent.ID,
                    name: selectedEvent.Name,
                    description: selectedEvent.Description,
                    location: selectedEvent.Location,
                    start_time: selectedEvent.StartTime,
                    end_time: selectedEvent.EndTime
                } : null}
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
                eventId={selectedEventForRSVP?.ID || ''}
                eventName={selectedEventForRSVP?.Name || ''}
                clubId={id || ''}
            />
            
            {/* Event Details Modal */}
            <Modal
                isOpen={isDetailsModalOpen}
                onClose={handleCloseDetailsModal}
                title="Event Details"
                maxWidth="700px"
            >
                <Modal.Body>
                    {selectedEventForDetails && (
                        <div className="event-details-content">
                            <h2 className="event-details-title">{selectedEventForDetails.Name}</h2>
                            
                            {selectedEventForDetails.Description && (
                                <Card variant="dark" padding="md" className="event-details-section">
                                    <h4>Description</h4>
                                    <p>{selectedEventForDetails.Description}</p>
                                </Card>
                            )}
                            
                            {selectedEventForDetails.Location && (
                                <Card variant="dark" padding="md" className="event-details-section">
                                    <h4>Location</h4>
                                    <p>{selectedEventForDetails.Location}</p>
                                </Card>
                            )}
                            
                            <Card variant="dark" padding="md" className="event-details-section">
                                <h4>Schedule</h4>
                                <div className="event-schedule">
                                    <div><strong>Start:</strong> {formatDateTime(selectedEventForDetails.StartTime)}</div>
                                    <div><strong>End:</strong> {formatDateTime(selectedEventForDetails.EndTime)}</div>
                                </div>
                            </Card>
                            
                            <Card variant="dark" padding="md" className="event-details-section">
                                <h4>RSVPs</h4>
                                <div className="event-rsvp-summary">
                                    <span style={{color: 'var(--color-primary)'}}>
                                        Yes: {rsvpCounts[selectedEventForDetails.ID]?.yes || 0}
                                    </span>
                                    <span style={{color: 'var(--color-cancel)'}}>
                                        No: {rsvpCounts[selectedEventForDetails.ID]?.no || 0}
                                    </span>
                                    <span style={{color: 'orange'}}>
                                        Maybe: {rsvpCounts[selectedEventForDetails.ID]?.maybe || 0}
                                    </span>
                                </div>
                            </Card>

                            {eventShifts[selectedEventForDetails.ID]?.length > 0 && (
                                <Card variant="dark" padding="md" className="event-details-section">
                                    <h4>Shifts</h4>
                                    <div className="event-shifts-list">
                                        {eventShifts[selectedEventForDetails.ID].map((shift) => (
                                            <div key={shift.ID} className="event-shift-item">
                                                {formatDateTime(shift.StartTime)} - {formatDateTime(shift.EndTime)}
                                            </div>
                                        ))}
                                    </div>
                                </Card>
                            )}
                        </div>
                    )}
                </Modal.Body>
                <Modal.Actions>
                    <Button
                        onClick={handleEditFromDetails}
                        variant="accept"
                    >
                        Edit
                    </Button>
                    <Button
                        onClick={handleDeleteFromDetails}
                        variant="cancel"
                        disabled={deleteLoading}
                    >
                        {deleteLoading ? 'Deleting...' : 'Delete'}
                    </Button>
                    <Button onClick={handleCloseDetailsModal} variant="secondary">
                        Close
                    </Button>
                </Modal.Actions>
            </Modal>
        </div>
    );
};

export default AdminClubEventList;