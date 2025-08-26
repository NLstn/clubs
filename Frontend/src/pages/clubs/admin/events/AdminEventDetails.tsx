import { FC, useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import api from "../../../../utils/api";
import EventRSVPList from "./EventRSVPList";
import EditEvent from "./EditEvent";
import "./AdminEventDetails.css";

interface UserRSVP {
    id: string;
    event_id: string;
    user_id: string;
    response: string;
    created_at: string;
    updated_at: string;
}

interface RSVPCounts {
    yes?: number;
    no?: number;
    maybe?: number;
}

interface EventDetailsData {
    id: string;
    name: string;
    description: string;
    location: string;
    start_time: string;
    end_time: string;
    created_at: string;
    created_by: string;
    updated_at: string;
    updated_by: string;
    user_rsvp?: UserRSVP;
}

interface Shift {
    id: string;
    startTime: string;
    endTime: string;
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
            // Fetch event details
            const eventResponse = await api.get(`/api/v1/clubs/${clubId}/events/${eventId}`, {
                signal: abortSignal
            });
            if (!abortSignal?.aborted) {
                setEventData(eventResponse.data);
            }

            // Fetch RSVP counts
            try {
                const rsvpResponse = await api.get(`/api/v1/clubs/${clubId}/events/${eventId}/rsvps`, {
                    signal: abortSignal
                });
                if (!abortSignal?.aborted) {
                    setRsvpCounts(rsvpResponse.data.counts || {});
                }
            } catch (rsvpError) {
                if (!abortSignal?.aborted) {
                    console.error("Error fetching RSVP counts:", rsvpError);
                    setRsvpCounts({});
                }
            }

            // Fetch shifts if available
            try {
                const shiftsResponse = await api.get(`/api/v1/clubs/${clubId}/events/${eventId}/shifts`, {
                    signal: abortSignal
                });
                if (!abortSignal?.aborted) {
                    setEventShifts(shiftsResponse.data || []);
                }
            } catch {
                // Shifts endpoint might not exist, ignore error
                if (!abortSignal?.aborted) {
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
        
        const confirmDelete = typeof window !== 'undefined'
            ? window.confirm("Are you sure you want to delete this event? This action cannot be undone.")
            : false;
        
        if (!confirmDelete) return;
        
        setDeleteLoading(true);
        
        try {
            await api.delete(`/api/v1/clubs/${clubId}/events/${eventId}`);
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
                        <button onClick={() => navigate(-1)} className="button-secondary">
                            Go Back
                        </button>
                        <Link to={`/clubs/${clubId}/admin`} className="button-primary">
                            Back to Admin
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

    const { user_rsvp } = eventData;
    const totalRSVPs = (rsvpCounts.yes || 0) + (rsvpCounts.no || 0) + (rsvpCounts.maybe || 0);

    return (
        <div className="admin-container">
            <div className="page-header">
                <div className="breadcrumb">
                    <Link to={`/clubs/${clubId}/admin`}>Admin Dashboard</Link>
                    <span> / </span>
                    <span>Event Details</span>
                </div>
                <div className="header-actions">
                    <button onClick={() => setIsEditModalOpen(true)} className="button-primary">
                        Edit Event
                    </button>
                    <button onClick={() => navigate(-1)} className="button-secondary">
                        Go Back
                    </button>
                </div>
            </div>

            <div className="admin-event-details-card">
                <div className="event-header">
                    <h1>{eventData.name}</h1>
                    <div className="event-actions">
                        <button 
                            onClick={() => setIsRSVPModalOpen(true)} 
                            className="button-info"
                        >
                            View RSVPs ({totalRSVPs})
                        </button>
                        <button 
                            onClick={handleDeleteEvent}
                            disabled={deleteLoading}
                            className="button-cancel"
                        >
                            {deleteLoading ? 'Deleting...' : 'Delete Event'}
                        </button>
                    </div>
                </div>
                
                {eventData.description && (
                    <div className="info-section">
                        <h3>Description</h3>
                        <p className="event-description">{eventData.description}</p>
                    </div>
                )}
                
                {eventData.location && (
                    <div className="info-section">
                        <h3>Location</h3>
                        <p className="event-location">{eventData.location}</p>
                    </div>
                )}
                
                <div className="admin-event-info">
                    <div className="info-section">
                        <h3>Event Schedule</h3>
                        <div className="schedule-details">
                            <div className="schedule-item">
                                <strong>Date:</strong> {formatDate(eventData.start_time)}
                            </div>
                            <div className="schedule-item">
                                <strong>Start Time:</strong> {formatTime(eventData.start_time)}
                            </div>
                            <div className="schedule-item">
                                <strong>End Time:</strong> {formatTime(eventData.end_time)}
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
                            {user_rsvp && (
                                <div className="admin-rsvp">
                                    <p>
                                        Your response: 
                                        <span className={`rsvp-status ${user_rsvp.response}`}>
                                            {user_rsvp.response === 'yes' ? ' Yes' : user_rsvp.response === 'no' ? ' No' : ' Maybe'}
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
                                    <div key={shift.id} className="shift-item">
                                        <span>{formatTime(shift.startTime)} - {formatTime(shift.endTime)}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}

                    <div className="info-section">
                        <h3>Event Management</h3>
                        <div className="meta-info">
                            <p><strong>Created:</strong> {formatDateTime(eventData.created_at)}</p>
                            {eventData.updated_at !== eventData.created_at && (
                                <p><strong>Last Updated:</strong> {formatDateTime(eventData.updated_at)}</p>
                            )}
                            <p><strong>Event ID:</strong> {eventData.id}</p>
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
                    eventName={eventData.name}
                    clubId={clubId!}
                />
            )}

            {/* Edit Modal */}
            {isEditModalOpen && (
                <EditEvent
                    isOpen={isEditModalOpen}
                    onClose={() => setIsEditModalOpen(false)}
                    event={{
                        id: eventData.id,
                        name: eventData.name,
                        description: eventData.description,
                        location: eventData.location,
                        start_time: eventData.start_time,
                        end_time: eventData.end_time
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
