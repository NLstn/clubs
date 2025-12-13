import { FC, useState, useEffect } from "react";
import { Table, TableColumn, Button, Card } from '@/components/ui';
import Modal from '@/components/ui/Modal';
import api from "../../../../utils/api";

interface EventRSVP {
    ID: string;
    EventID: string;
    UserID: string;
    Response: string;
    CreatedAt: string;
    UpdatedAt: string;
    User: {
        ID: string;
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
            // OData v2: Use EventRSVPs navigation property with User expansion
            const response = await api.get(`/api/v2/Events('${eventId}')/EventRSVPs?$expand=User`);
            const rsvpList = response.data.value || [];
            setRsvps(rsvpList);
            
            // Compute counts by grouping RSVPs by Response field
            const computedCounts = rsvpList.reduce((acc: RSVPCounts, rsvp: EventRSVP) => {
                const responseKey = rsvp.Response.toLowerCase() as keyof RSVPCounts;
                acc[responseKey] = (acc[responseKey] || 0) + 1;
                return acc;
            }, {});
            setCounts(computedCounts);
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
            render: (rsvp) => rsvp.User ? `${rsvp.User.FirstName} ${rsvp.User.LastName}`.trim() || 'Unknown' : 'Unknown'
        },
        {
            key: 'response',
            header: 'Response',
            render: (rsvp) => (
                <span style={{ 
                    color: rsvp.Response === 'yes' ? 'var(--color-primary)' : 
                           rsvp.Response === 'no' ? 'var(--color-cancel)' : 'orange',
                    fontWeight: 'bold'
                }}>
                    {rsvp.Response === 'yes' ? 'Yes' : rsvp.Response === 'no' ? 'No' : 'Maybe'}
                </span>
            )
        },
        {
            key: 'date_responded',
            header: 'Date Responded',
            render: (rsvp) => formatDateTime(rsvp.UpdatedAt)
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
                <Card variant="dark" padding="md" style={{ marginBottom: '20px' }}>
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
                </Card>

                {loading && <p>Loading RSVPs...</p>}
                {error && <p style={{ color: 'red' }}>Error: {error}</p>}
                
                {!loading && !error && (
                    <div style={{ maxHeight: '400px', overflowY: 'auto' }}>
                        <Table
                            columns={rsvpColumns}
                            data={rsvps}
                            keyExtractor={(rsvp) => rsvp.ID}
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
