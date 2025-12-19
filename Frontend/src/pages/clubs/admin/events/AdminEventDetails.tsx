import { FC, useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { Button } from "../../../../components/ui";
import PageHeader from "../../../../components/layout/PageHeader";
import api from "../../../../utils/api";
import EventRSVPList from "./EventRSVPList";
import EditEvent from "./EditEvent";
import "./AdminEventDetails.css";
import '../../../../styles/events.css';

interface EventRSVP {
    Response: string;
}

interface UserRSVP {
    ID: string;
    EventID: string;
    UserID: string;
    Response: string;
    CreatedAt: string;
    UpdatedAt: string;
}

interface RSVPCounts {
    yes?: number;
    no?: number;
    maybe?: number;
}

interface Shift {
    ID: string;
    StartTime: string;
    EndTime: string;
}

interface EventDetailsData {
    ID: string;
    Name: string;
    Description: string;
    Location: string;
    StartTime: string;
    EndTime: string;
    CreatedAt: string;
    CreatedBy: string;
    UpdatedAt: string;
    UpdatedBy: string;
    UserRSVP?: UserRSVP;
    EventRSVPs?: EventRSVP[];
    Shifts?: Shift[];
}

const AdminEventDetails: FC = () => {
    const { clubId, eventId } = useParams<{ clubId: string; eventId: string }>();
    const navigate = useNavigate();
    const [eventData, setEventData] = useState<EventDetailsData | null>(null);
    const [rsvpCounts, setRsvpCounts] = useState<RSVPCounts>({});
    const [eventShifts, setEventShifts] = useState<Shift[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isRSVPModalOpen, setIsRSVPModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [deleteLoading, setDeleteLoading] = useState(false);

    const fetchEventDetails = async (abortSignal?: AbortSignal) => {
        if (!clubId || !eventId) return;
        
        setLoading(true);
        setError(null);
        
        try {
            // OData v2: Fetch Event with expanded Shifts in a single call
            // Use GetRSVPCounts function for efficient server-side aggregation
            const eventResponse = await api.get(
                `/api/v2/Events('${eventId}')?$expand=Shifts`,
                { signal: abortSignal }
            );
            
            if (!abortSignal?.aborted) {
                const event = eventResponse.data;
                setEventData(event);
                
                // Fetch RSVP counts using server-side function
                try {
                    const countsResponse = await api.get<{ Yes: number; No: number; Maybe: number }>(
                        `/api/v2/Events('${eventId}')/GetRSVPCounts()`,
                        { signal: abortSignal }
                    );
                    
                    if (!abortSignal?.aborted) {
                        // Transform PascalCase response to camelCase for frontend
                        setRsvpCounts({
                            yes: countsResponse.data.Yes,
                            no: countsResponse.data.No,
                            maybe: countsResponse.data.Maybe
                        });
                    }
                } catch (countsError) {
                    console.error("Error fetching RSVP counts:", countsError);
                    setRsvpCounts({});
                }
                
                // Set shifts from the expanded Shifts
                if (event.Shifts && Array.isArray(event.Shifts)) {
                    setEventShifts(event.Shifts);
                } else {
                    setEventShifts([]);
                }
            }
        } catch (error: unknown) {
            if (!abortSignal?.aborted) {
                console.error("Error fetching event details:", error);
                if (error && typeof error === 'object' && 'response' in error) {
                    const httpError = error as { response?: { status?: number } };
                    if (httpError.response?.status === 404) {
                        setError("Event not found");
                    } else if (httpError.response?.status === 403) {
                        setError("You don't have permission to view this event");
                    } else {
                        setError("Failed to load event details");
                    }
                } else {
                    setError("Failed to load event details");
                }
            }
        } finally {
            if (!abortSignal?.aborted) {
                setLoading(false);
            }
        }
    };

    useEffect(() => {
        const abortController = new AbortController();
        fetchEventDetails(abortController.signal);
        
        return () => {
            abortController.abort();
        };
    }, [clubId, eventId]); // eslint-disable-line react-hooks/exhaustive-deps

    const handleDeleteEvent = async () => {
        if (!clubId || !eventId) return;
        
        const confirmDelete = window.confirm(
            "Are you sure you want to delete this event? This action cannot be undone."
        );
        
        if (!confirmDelete) return;
        
        setDeleteLoading(true);
        
        try {
            // OData v2: Delete event
            await api.delete(`/api/v2/Events('${eventId}')`);
            navigate(`/clubs/${clubId}/admin`, { 
                state: { message: "Event deleted successfully" } 
            });
        } catch (error) {
            console.error("Error deleting event:", error);
            alert("Failed to delete event. Please try again.");
        } finally {
            setDeleteLoading(false);
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

    const formatDate = (timestamp: string) => {
        try {
            const date = new Date(timestamp);
            return date.toLocaleDateString();
        } catch {
            return timestamp;
        }
    };

    const formatTime = (timestamp: string) => {
        try {
            const time = new Date(timestamp);
            return time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        } catch {
            return timestamp;
        }
    };

    if (loading) {
        return (
            <div className="admin-container">
                <div className="loading">Loading event details...</div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="admin-container">
                <div className="error-container">
                    <h2>Error</h2>
                    <p>{error}</p>
                    <div className="button-group">
                        <Button onClick={() => navigate(-1)} variant="secondary">
                            Go Back
                        </Button>
                        <Link to={`/clubs/${clubId}/admin`} style={{ textDecoration: 'none' }}>
                            <Button variant="primary">
                                Back to Admin
                            </Button>
                        </Link>
                    </div>
                </div>
            </div>
        );
    }

    if (!eventData) {
        return (
            <div className="admin-container">
                <div>No event data available</div>
            </div>
        );
    }

    const { UserRSVP: userRSVP } = eventData;
    const totalRSVPs = (rsvpCounts.yes || 0) + (rsvpCounts.no || 0) + (rsvpCounts.maybe || 0);

    return (
        <div className="admin-container">
            <PageHeader
                variant="simple"
                actions={
                    <>
                        <Button onClick={() => setIsEditModalOpen(true)} variant="primary">
                            Edit Event
                        </Button>
                        <Button onClick={() => navigate(-1)} variant="secondary">
                            Go Back
                        </Button>
                    </>
                }
            >
                <div className="breadcrumb">
                    <Link to={`/clubs/${clubId}/admin`}>Admin Dashboard</Link>
                    <span> / </span>
                    <span>Event Details</span>
                </div>
            </PageHeader>

            <div className="admin-event-details-card">
                <div className="event-header">
                    <h1>{eventData.Name}</h1>
                    <div className="event-actions">
                        <Button 
                            onClick={() => setIsRSVPModalOpen(true)} 
                            variant="secondary"
                            size="sm"
                        >
                            View RSVPs ({totalRSVPs})
                        </Button>
                        <Button 
                            onClick={handleDeleteEvent}
                            disabled={deleteLoading}
                            variant="cancel"
                            size="sm"
                        >
                            {deleteLoading ? 'Deleting...' : 'Delete Event'}
                        </Button>
                    </div>
                </div>
                
                {eventData.Description && (
                    <div className="info-section">
                        <h3>Description</h3>
                        <p className="event-description">{eventData.Description}</p>
                    </div>
                )}
                
                {eventData.Location && (
                    <div className="info-section">
                        <h3>Location</h3>
                        <p className="event-location">{eventData.Location}</p>
                    </div>
                )}
                
                <div className="admin-event-info">
                    <div className="info-section">
                        <h3>Event Schedule</h3>
                        <div className="schedule-details">
                            <div className="schedule-item">
                                <strong>Date:</strong> {formatDate(eventData.StartTime)}
                            </div>
                            <div className="schedule-item">
                                <strong>Start Time:</strong> {formatTime(eventData.StartTime)}
                            </div>
                            <div className="schedule-item">
                                <strong>End Time:</strong> {formatTime(eventData.EndTime)}
                            </div>
                        </div>
                    </div>

                    <div className="info-section">
                        <h3>RSVP Summary</h3>
                        <div className="rsvp-summary">
                            <div className="rsvp-stats">
                                <div className="stat-item yes">
                                    <span className="stat-number">{rsvpCounts.yes || 0}</span>
                                    <span className="stat-label">Yes</span>
                                </div>
                                <div className="stat-item no">
                                    <span className="stat-number">{rsvpCounts.no || 0}</span>
                                    <span className="stat-label">No</span>
                                </div>
                                <div className="stat-item maybe">
                                    <span className="stat-number">{rsvpCounts.maybe || 0}</span>
                                    <span className="stat-label">Maybe</span>
                                </div>
                                <div className="stat-item total">
                                    <span className="stat-number">{totalRSVPs}</span>
                                    <span className="stat-label">Total</span>
                                </div>
                            </div>
                            {userRSVP && (
                                <div className="admin-rsvp">
                                    <p>
                                        Your response: 
                                        <span className={`rsvp-status ${userRSVP.Response}`}>
                                            {userRSVP.Response === 'yes' ? ' Yes' : userRSVP.Response === 'no' ? ' No' : ' Maybe'}
                                        </span>
                                    </p>
                                </div>
                            )}
                        </div>
                    </div>

                    {eventShifts.length > 0 && (
                        <div className="info-section">
                            <h3>Shifts</h3>
                            <div className="shifts-list">
                                {eventShifts.map((shift) => (
                                    <div key={shift.ID} className="shift-item">
                                        <span>{formatTime(shift.StartTime)} - {formatTime(shift.EndTime)}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}

                    <div className="info-section">
                        <h3>Event Management</h3>
                        <div className="meta-info">
                            <p><strong>Created:</strong> {formatDateTime(eventData.CreatedAt)}</p>
                            {eventData.UpdatedAt !== eventData.CreatedAt && (
                                <p><strong>Last Updated:</strong> {formatDateTime(eventData.UpdatedAt)}</p>
                            )}
                            <p><strong>Event ID:</strong> {eventData.ID}</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* RSVP Modal */}
            {isRSVPModalOpen && (
                <EventRSVPList
                    isOpen={isRSVPModalOpen}
                    onClose={() => setIsRSVPModalOpen(false)}
                    eventId={eventId!}
                    eventName={eventData.Name}
                    clubId={clubId!}
                />
            )}

            {/* Edit Modal */}
            {isEditModalOpen && (
                <EditEvent
                    isOpen={isEditModalOpen}
                    onClose={() => setIsEditModalOpen(false)}
                    event={{
                        id: eventData.ID,
                        name: eventData.Name,
                        description: eventData.Description,
                        location: eventData.Location,
                        start_time: eventData.StartTime,
                        end_time: eventData.EndTime
                    }}
                    clubId={clubId}
                    onSuccess={() => {
                        setIsEditModalOpen(false);
                        fetchEventDetails(); // Refresh data after edit
                    }}
                />
            )}
        </div>
    );
};

export default AdminEventDetails;
