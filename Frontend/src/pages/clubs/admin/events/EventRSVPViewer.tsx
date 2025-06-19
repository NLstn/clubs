import { useState, useEffect } from 'react';
import api from '../../../../utils/api';
import RSVPHistoryModal from './RSVPHistoryModal';

interface User {
    id: string;
    name: string;
    email: string;
}

interface RSVP {
    id: string;
    response: string;
    user: User;
}

interface RSVPCounts {
    yes?: number;
    no?: number;
}

interface EventRSVPViewerProps {
    isOpen: boolean;
    onClose: () => void;
    clubId: string;
    eventId: string;
    eventName: string;
}

const EventRSVPViewer = ({ isOpen, onClose, clubId, eventId, eventName }: EventRSVPViewerProps) => {
    const [rsvps, setRsvps] = useState<RSVP[]>([]);
    const [counts, setCounts] = useState<RSVPCounts>({});
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
    const [isHistoryModalOpen, setIsHistoryModalOpen] = useState(false);

    const fetchRSVPs = async () => {
        if (!isOpen || !eventId) return;
        
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

    const handleViewHistory = (userId: string) => {
        setSelectedUserId(userId);
        setIsHistoryModalOpen(true);
    };

    const handleCloseHistoryModal = () => {
        setSelectedUserId(null);
        setIsHistoryModalOpen(false);
    };

    if (!isOpen) return null;

    const yesRSVPs = rsvps.filter(rsvp => rsvp.response === 'yes');
    const noRSVPs = rsvps.filter(rsvp => rsvp.response === 'no');

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '800px', width: '90%' }}>
                <div className="modal-header">
                    <h2>RSVPs for {eventName}</h2>
                    <button onClick={onClose} className="modal-close-button">&times;</button>
                </div>
                
                <div className="modal-body">
                    {loading && <p>Loading RSVPs...</p>}
                    {error && <p style={{color: 'red'}}>Error: {error}</p>}
                    
                    {!loading && !error && (
                        <>
                            <div style={{ marginBottom: '20px' }}>
                                <h3>Summary</h3>
                                <p>
                                    <span style={{color: 'green', fontWeight: 'bold'}}>Yes: {counts.yes || 0}</span>
                                    {' | '}
                                    <span style={{color: 'red', fontWeight: 'bold'}}>No: {counts.no || 0}</span>
                                    {' | '}
                                    <span style={{color: 'gray'}}>Total: {rsvps.length}</span>
                                </p>
                            </div>

                            <div style={{ display: 'flex', gap: '20px' }}>
                                {/* Yes RSVPs */}
                                <div style={{ flex: 1 }}>
                                    <h4 style={{ color: 'green' }}>Attending ({yesRSVPs.length})</h4>
                                    {yesRSVPs.length === 0 ? (
                                        <p style={{ fontStyle: 'italic', color: 'gray' }}>No one is attending yet</p>
                                    ) : (
                                        <table className="basic-table" style={{ width: '100%' }}>
                                            <thead>
                                                <tr>
                                                    <th>Name</th>
                                                    <th>Email</th>
                                                    <th>Actions</th>
                                                </tr>
                                            </thead>
                                            <tbody>
                                                {yesRSVPs.map(rsvp => (
                                                    <tr key={rsvp.id}>
                                                        <td>{rsvp.user.name}</td>
                                                        <td>{rsvp.user.email}</td>
                                                        <td>
                                                            <button
                                                                onClick={() => handleViewHistory(rsvp.user.id)}
                                                                className="button-accept"
                                                                style={{ fontSize: '12px', padding: '4px 8px' }}
                                                            >
                                                                View History
                                                            </button>
                                                        </td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    )}
                                </div>

                                {/* No RSVPs */}
                                <div style={{ flex: 1 }}>
                                    <h4 style={{ color: 'red' }}>Not Attending ({noRSVPs.length})</h4>
                                    {noRSVPs.length === 0 ? (
                                        <p style={{ fontStyle: 'italic', color: 'gray' }}>No one has declined yet</p>
                                    ) : (
                                        <table className="basic-table" style={{ width: '100%' }}>
                                            <thead>
                                                <tr>
                                                    <th>Name</th>
                                                    <th>Email</th>
                                                    <th>Actions</th>
                                                </tr>
                                            </thead>
                                            <tbody>
                                                {noRSVPs.map(rsvp => (
                                                    <tr key={rsvp.id}>
                                                        <td>{rsvp.user.name}</td>
                                                        <td>{rsvp.user.email}</td>
                                                        <td>
                                                            <button
                                                                onClick={() => handleViewHistory(rsvp.user.id)}
                                                                className="button-accept"
                                                                style={{ fontSize: '12px', padding: '4px 8px' }}
                                                            >
                                                                View History
                                                            </button>
                                                        </td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    )}
                                </div>
                            </div>
                        </>
                    )}
                </div>
                
                <div className="modal-footer">
                    <button onClick={onClose} className="button-cancel">Close</button>
                </div>

                {selectedUserId && (
                    <RSVPHistoryModal
                        isOpen={isHistoryModalOpen}
                        onClose={handleCloseHistoryModal}
                        clubId={clubId}
                        eventId={eventId}
                        userId={selectedUserId}
                    />
                )}
            </div>
        </div>
    );
};

export default EventRSVPViewer;