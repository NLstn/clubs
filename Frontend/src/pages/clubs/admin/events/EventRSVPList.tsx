import { FC, useState, useEffect } from "react";
import api from "../../../../utils/api";

interface EventRSVP {
    id: string;
    event_id: string;
    user_id: string;
    response: string;
    created_at: string;
    updated_at: string;
    user: {
        id: string;
        FirstName: string;
        LastName: string;
        Email: string;
    };
}

interface RSVPCounts {
    yes?: number;
    no?: number;
}

interface EventRSVPListProps {
    isOpen: boolean;
    onClose: () => void;
    eventId: string;
    eventName: string;
    clubId: string;
}

const EventRSVPList: FC<EventRSVPListProps> = ({ isOpen, onClose, eventId, eventName, clubId }) => {
    const [rsvps, setRsvps] = useState<EventRSVP[]>([]);
    const [counts, setCounts] = useState<RSVPCounts>({});
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchRSVPs = async () => {
        if (!isOpen || !eventId || !clubId) return;
        
        setLoading(true);
        setError(null);
        
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/events/${eventId}/rsvps`);
            setRsvps(response.data.rsvps || []);
            setCounts(response.data.counts || {});
        } catch (error) {
            console.error("Error fetching RSVPs:", error);
            setError("Failed to load RSVPs");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchRSVPs();
    }, [isOpen, eventId, clubId]); // eslint-disable-line react-hooks/exhaustive-deps

    const formatDateTime = (timestamp: string) => {
        try {
            const dateTime = new Date(timestamp);
            return dateTime.toLocaleString();
        } catch {
            return timestamp;
        }
    };

    if (!isOpen) return null;

    return (
        <div className="modal">
            <div className="modal-content" style={{ maxWidth: '800px', width: '90%' }}>
                <h2>RSVPs for {eventName}</h2>
                
                {/* Summary */}
                <div style={{ marginBottom: '20px', padding: '15px', backgroundColor: 'var(--color-background)', borderRadius: 'var(--border-radius-md)', border: '1px solid var(--color-border)' }}>
                    <h3 style={{ margin: '0 0 10px 0', color: 'var(--color-text)' }}>Summary</h3>
                    <div style={{ display: 'flex', gap: '20px' }}>
                        <span style={{ color: 'var(--color-primary)', fontWeight: 'bold' }}>
                            Yes: {counts.yes || 0}
                        </span>
                        <span style={{ color: 'var(--color-cancel)', fontWeight: 'bold' }}>
                            No: {counts.no || 0}
                        </span>
                        <span style={{ color: 'var(--color-text-secondary)', fontWeight: 'bold' }}>
                            Total: {(counts.yes || 0) + (counts.no || 0)}
                        </span>
                    </div>
                </div>

                {loading && <p>Loading RSVPs...</p>}
                {error && <p style={{ color: 'red' }}>Error: {error}</p>}
                
                {!loading && !error && (
                    <>
                        {rsvps.length === 0 ? (
                            <p style={{ fontStyle: 'italic', color: '#666' }}>
                                No RSVPs yet for this event.
                            </p>
                        ) : (
                            <div style={{ maxHeight: '400px', overflowY: 'auto' }}>
                                <table className="basic-table">
                                    <thead>
                                        <tr>
                                            <th>Member</th>
                                            <th>Response</th>
                                            <th>Date Responded</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {rsvps.map((rsvp) => (
                                            <tr key={rsvp.id}>
                                                <td>{rsvp.user ? `${rsvp.user.FirstName} ${rsvp.user.LastName}`.trim() || 'Unknown' : 'Unknown'}</td>
                                                <td>
                                                    <span style={{ 
                                                        color: rsvp.response === 'yes' ? 'var(--color-primary)' : 'var(--color-cancel)',
                                                        fontWeight: 'bold'
                                                    }}>
                                                        {rsvp.response === 'yes' ? 'Yes' : 'No'}
                                                    </span>
                                                </td>
                                                <td>{formatDateTime(rsvp.updated_at)}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </>
                )}

                <div className="modal-actions" style={{ marginTop: '20px' }}>
                    <button onClick={onClose} className="button-cancel">
                        Close
                    </button>
                </div>
            </div>
        </div>
    );
};

export default EventRSVPList;
