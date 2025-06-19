import { useState, useEffect } from 'react';
import api from '../../../../utils/api';

interface User {
    id: string;
    name: string;
    email: string;
}

interface RSVPHistoryEntry {
    id: string;
    response: string;
    created_at: string;
    user: User;
}

interface RSVPHistoryModalProps {
    isOpen: boolean;
    onClose: () => void;
    clubId: string;
    eventId: string;
    userId: string;
}

const RSVPHistoryModal = ({ isOpen, onClose, clubId, eventId, userId }: RSVPHistoryModalProps) => {
    const [history, setHistory] = useState<RSVPHistoryEntry[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [userName, setUserName] = useState<string>('');

    const fetchHistory = async () => {
        if (!isOpen || !eventId || !userId) return;
        
        setLoading(true);
        setError(null);
        
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/rsvp-history?eventid=${eventId}&userid=${userId}`);
            const historyData = response.data.history || [];
            setHistory(historyData);
            
            // Set user name from the first entry if available
            if (historyData.length > 0 && historyData[0].user) {
                setUserName(historyData[0].user.name);
            }
        } catch (error) {
            console.error("Error fetching RSVP history:", error);
            setError("Failed to load RSVP history");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchHistory();
    }, [isOpen, eventId, userId, clubId]); // eslint-disable-line react-hooks/exhaustive-deps

    const formatDateTime = (timestamp: string) => {
        try {
            const dateTime = new Date(timestamp);
            if (isNaN(dateTime.getTime())) {
                return 'Invalid date/time';
            }
            return dateTime.toLocaleString();
        } catch {
            return 'Parse error';
        }
    };

    const getResponseStyle = (response: string) => {
        return {
            color: response === 'yes' ? 'green' : 'red',
            fontWeight: 'bold' as const
        };
    };

    const getResponseText = (response: string) => {
        return response === 'yes' ? 'Will Attend' : 'Will Not Attend';
    };

    if (!isOpen) return null;

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '600px', width: '90%' }}>
                <div className="modal-header">
                    <h2>RSVP History for {userName}</h2>
                    <button onClick={onClose} className="modal-close-button">&times;</button>
                </div>
                
                <div className="modal-body">
                    {loading && <p>Loading RSVP history...</p>}
                    {error && <p style={{color: 'red'}}>Error: {error}</p>}
                    
                    {!loading && !error && (
                        <>
                            {history.length === 0 ? (
                                <p style={{ fontStyle: 'italic', color: 'gray', textAlign: 'center' }}>
                                    No RSVP history found for this user and event.
                                </p>
                            ) : (
                                <>
                                    <p style={{ marginBottom: '15px', color: 'gray' }}>
                                        This shows all RSVP changes in chronological order (oldest first).
                                    </p>
                                    <table className="basic-table" style={{ width: '100%' }}>
                                        <thead>
                                            <tr>
                                                <th>Date & Time</th>
                                                <th>Response</th>
                                                <th>Change #</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {history.map((entry, index) => (
                                                <tr key={entry.id}>
                                                    <td>{formatDateTime(entry.created_at)}</td>
                                                    <td style={getResponseStyle(entry.response)}>
                                                        {getResponseText(entry.response)}
                                                    </td>
                                                    <td style={{ textAlign: 'center' }}>
                                                        {index + 1}
                                                        {index === history.length - 1 && (
                                                            <span style={{ fontSize: '12px', color: 'blue', marginLeft: '5px' }}>
                                                                (Latest)
                                                            </span>
                                                        )}
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                    
                                    {history.length > 1 && (
                                        <div style={{ marginTop: '15px', padding: '10px', backgroundColor: '#f8f9fa', borderRadius: '4px' }}>
                                            <strong>Summary:</strong> This user changed their RSVP {history.length - 1} time{history.length - 1 !== 1 ? 's' : ''} after their initial response.
                                        </div>
                                    )}
                                </>
                            )}
                        </>
                    )}
                </div>
                
                <div className="modal-footer">
                    <button onClick={onClose} className="button-cancel">Close</button>
                </div>
            </div>
        </div>
    );
};

export default RSVPHistoryModal;