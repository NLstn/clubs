import { FC, useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
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

const TeamEventDetails: FC = () => {
    const { clubId, teamId, eventId } = useParams<{ clubId: string; teamId: string; eventId: string }>();
    const navigate = useNavigate();
    const [eventData, setEventData] = useState<EventDetailsData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [rsvpLoading, setRsvpLoading] = useState(false);

    const fetchEventDetails = async (abortSignal?: AbortSignal) => {
        if (!clubId || !teamId || !eventId) return;
        setLoading(true);
        setError(null);
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/events/${eventId}`, { signal: abortSignal });
            if (!abortSignal?.aborted) {
                setEventData(response.data);
            }
        } catch (err: any) {
            if (!abortSignal?.aborted) {
                console.error('Error fetching event details:', err);
                if (err?.response?.status === 404) {
                    setError('Event not found');
                } else if (err?.response?.status === 403) {
                    setError("You don't have permission to view this event");
                } else {
                    setError('Failed to load event details');
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
        return () => abortController.abort();
    }, [clubId, teamId, eventId]); // eslint-disable-line react-hooks/exhaustive-deps

    const handleRSVP = async (response: string) => {
        if (!clubId || !teamId || !eventId || rsvpLoading) return;
        setRsvpLoading(true);
        try {
            await api.post(`/api/v1/clubs/${clubId}/teams/${teamId}/events/${eventId}/rsvp`, { response });
            await fetchEventDetails();
        } catch (err) {
            console.error('Error updating RSVP:', err);
            alert('Failed to update RSVP. Please try again.');
        } finally {
            setRsvpLoading(false);
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
        return <div className="container"><div className="loading">Loading event details...</div></div>;
    }
    if (error) {
        return (
            <div className="container">
                <div className="error-container">
                    <h2>Error</h2>
                    <p>{error}</p>
                    <div className="button-group">
                        <button onClick={() => navigate(-1)} className="button-secondary">Go Back</button>
                        <Link to={`/clubs/${clubId}`} className="button-primary">Back to Club</Link>
                    </div>
                </div>
            </div>
        );
    }
    if (!eventData) {
        return <div className="container"><div>No event data available</div></div>;
    }
    const { user_rsvp } = eventData;
    return (
        <div className="container">
            <div className="page-header">
                <div className="breadcrumb">
                    <Link to={`/clubs/${clubId}`}>Club</Link>
                    <span> / </span>
                    <Link to={`/clubs/${clubId}/teams/${teamId}`}>Team</Link>
                    <span> / </span>
                    <span>Event Details</span>
                </div>
                <button onClick={() => navigate(-1)} className="button-secondary">Go Back</button>
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
                            <div className="schedule-item"><strong>Date:</strong> {formatDate(eventData.start_time)}</div>
                            <div className="schedule-item"><strong>Start:</strong> {formatTime(eventData.start_time)}</div>
                            <div className="schedule-item"><strong>End:</strong> {formatTime(eventData.end_time)}</div>
                        </div>
                    </div>
                </div>
                <div className="rsvp-section">
                    <h3>Your RSVP</h3>
                    <div className="rsvp-buttons">
                        <button onClick={() => handleRSVP('yes')} className={user_rsvp?.response === 'yes' ? 'button-accept' : 'button'} disabled={rsvpLoading}>Yes</button>
                        <button onClick={() => handleRSVP('maybe')} className={user_rsvp?.response === 'maybe' ? 'button-maybe' : 'button'} disabled={rsvpLoading}>Maybe</button>
                        <button onClick={() => handleRSVP('no')} className={user_rsvp?.response === 'no' ? 'button-cancel' : 'button'} disabled={rsvpLoading}>No</button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default TeamEventDetails;
