import { FC, useState, useEffect } from "react";
import { Table, TableColumn, Button } from '@/components/ui';
import Modal from '@/components/ui/Modal';
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
    maybe?: number;
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

    // Define table columns for RSVPs
    const rsvpColumns: TableColumn<EventRSVP>[] = [
        {
            key: 'member',
            header: 'Member',
            render: (rsvp) => rsvp.user ? `${rsvp.user.FirstName} ${rsvp.user.LastName}`.trim() || 'Unknown' : 'Unknown'
        },
        {
            key: 'response',
            header: 'Response',
            render: (rsvp) => (
                <span style={{ 
                    color: rsvp.response === 'yes' ? 'var(--color-primary)' : 
                           rsvp.response === 'no' ? 'var(--color-cancel)' : 'orange',
                    fontWeight: 'bold'
                }}>
                    {rsvp.response === 'yes' ? 'Yes' : rsvp.response === 'no' ? 'No' : 'Maybe'}
                </span>
            )
        },
        {
            key: 'date_responded',
            header: 'Date Responded',
            render: (rsvp) => formatDateTime(rsvp.updated_at)
        }
    ];

    if (!isOpen) return null;

    return (
        <Modal 
            isOpen={isOpen} 
            onClose={onClose}
            title={`RSVPs for ${eventName}`}
            maxWidth="800px"
        >
            <Modal.Body>
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
                        <span style={{ color: 'orange', fontWeight: 'bold' }}>
                            Maybe: {counts.maybe || 0}
                        </span>
                        <span style={{ color: 'var(--color-text-secondary)', fontWeight: 'bold' }}>
                            Total: {(counts.yes || 0) + (counts.no || 0) + (counts.maybe || 0)}
                        </span>
                    </div>
                </div>

                {loading && <p>Loading RSVPs...</p>}
                {error && <p style={{ color: 'red' }}>Error: {error}</p>}
                
                {!loading && !error && (
                    <div style={{ maxHeight: '400px', overflowY: 'auto' }}>
                        <Table
                            columns={rsvpColumns}
                            data={rsvps}
                            keyExtractor={(rsvp) => rsvp.id}
                            emptyMessage="No RSVPs yet for this event."
                        />
                    </div>
                )}
            </Modal.Body>
            <Modal.Actions>
                <Button onClick={onClose} variant="cancel">
                    Close
                </Button>
            </Modal.Actions>
        </Modal>
    );
};

export default EventRSVPList;
