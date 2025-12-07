import { FC, useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { Button } from "../../../components/ui";
import api from "../../../utils/api";
import "./EventDetails.css";

interface UserRSVP {
    id: string;
    event_id: string;
    user_id: string;
    response: string;
    created_at: string;
    updated_at: string;
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

const EventDetails: FC = () => {
    const { clubId, eventId } = useParams<{ clubId: string; eventId: string }>();
    const navigate = useNavigate();
    const [eventData, setEventData] = useState<EventDetailsData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [rsvpLoading, setRsvpLoading] = useState(false);

    const fetchEventDetails = async (abortSignal?: AbortSignal) => {
        if (!clubId || !eventId) return;
        
        setLoading(true);
        setError(null);
        
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/events/${eventId}`, {
                signal: abortSignal
            });
            if (!abortSignal?.aborted) {
                setEventData(response.data);
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

    const handleRSVP = async (response: string) => {
        if (!clubId || !eventId || rsvpLoading) return;
        
        setRsvpLoading(true);
        
        try {
            await api.post(`/api/v1/clubs/${clubId}/events/${eventId}/rsvp`, { response });
            // Refresh event details to update RSVP status
            await fetchEventDetails();
        } catch (error) {
            console.error("Error updating RSVP:", error);
            alert("Failed to update RSVP. Please try again.");
        } finally {
            setRsvpLoading(false);
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
            <div className="container">
                <div className="loading">Loading event details...</div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="container">
                <div className="error-container">
                    <h2>Error</h2>
                    <p>{error}</p>
                    <div className="button-group">
                        <Button onClick={() => navigate(-1)} variant="secondary">
                            Go Back
                        </Button>
                        <Link to={`/clubs/${clubId}`} style={{ textDecoration: 'none' }}>
                            <Button variant="primary">
                                Back to Club
                            </Button>
                        </Link>
                    </div>
                </div>
            </div>
        );
    }

    if (!eventData) {
        return (
            <div className="container">
                <div>No event data available</div>
            </div>
        );
    }

    const { user_rsvp } = eventData;

    return (
        <div className="container">
            <div className="page-header">
                <div className="breadcrumb">
                    <Link to={`/clubs/${clubId}`}>Club</Link>
                    <span> / </span>
                    <span>Event Details</span>
                </div>
                <Button onClick={() => navigate(-1)} variant="secondary">
                    Go Back
                </Button>
            </div>

            <div className="event-details-card">
                <h1>{eventData.name}</h1>
                
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
                
                <div className="event-info">
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
                        <h3>Your RSVP</h3>
                        <div className="rsvp-section">
                            {user_rsvp ? (
                                <div className="current-rsvp">
                                    <p>
                                        You have responded: 
                                        <span className={`rsvp-status ${user_rsvp.response}`}>
                                            {user_rsvp.response === 'yes' ? ' Yes' : 
                                             user_rsvp.response === 'no' ? ' No' : ' Maybe'}
                                        </span>
                                    </p>
                                    <p className="rsvp-date">
                                        Responded on: {formatDateTime(user_rsvp.updated_at)}
                                    </p>
                                </div>
                            ) : (
                                <p>You haven't responded to this event yet.</p>
                            )}
                            
                            <div className="rsvp-actions">
                                <h4>Update your response:</h4>
                                <div className="button-group">
                                    <Button 
                                        variant="primary"
                                        onClick={() => handleRSVP('yes')}
                                        disabled={rsvpLoading}
                                        className={user_rsvp?.response === 'yes' ? 'active' : ''}
                                    >
                                        {rsvpLoading ? 'Updating...' : 'Yes, I\'ll attend'}
                                    </Button>
                                    <Button 
                                        variant="maybe"
                                        onClick={() => handleRSVP('maybe')}
                                        disabled={rsvpLoading}
                                        className={user_rsvp?.response === 'maybe' ? 'active' : ''}
                                    >
                                        {rsvpLoading ? 'Updating...' : 'Maybe, I\'m not sure'}
                                    </Button>
                                    <Button 
                                        variant="cancel"
                                        onClick={() => handleRSVP('no')}
                                        disabled={rsvpLoading}
                                        className={user_rsvp?.response === 'no' ? 'active' : ''}
                                    >
                                        {rsvpLoading ? 'Updating...' : 'No, I can\'t attend'}
                                    </Button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div className="info-section">
                        <h3>Event Details</h3>
                        <div className="meta-info">
                            <p><strong>Created:</strong> {formatDateTime(eventData.created_at)}</p>
                            {eventData.updated_at !== eventData.created_at && (
                                <p><strong>Last Updated:</strong> {formatDateTime(eventData.updated_at)}</p>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default EventDetails;
