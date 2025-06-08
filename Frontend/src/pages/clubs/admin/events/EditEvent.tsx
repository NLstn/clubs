import { FC, useState, useEffect, useCallback } from "react";
import api from "../../../../utils/api";

interface Event {
    id: string;
    name: string;
    start_date: string;
    start_time: string;
    end_date: string;
    end_time: string;
}

interface Shift {
    id: string;
    startTime: string;
    endTime: string;
    eventId: string;
}

interface EditEventProps {
    isOpen: boolean;
    onClose: () => void;
    event: Event | null;
    clubId: string | undefined;
    onSuccess: () => void;
}

const EditEvent: FC<EditEventProps> = ({ isOpen, onClose, event, clubId, onSuccess }) => {
    const [name, setName] = useState<string>('');
    const [startDate, setStartDate] = useState<string>('');
    const [startTime, setStartTime] = useState<string>('');
    const [endDate, setEndDate] = useState<string>('');
    const [endTime, setEndTime] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    
    // Shift management state
    const [shifts, setShifts] = useState<Shift[]>([]);
    const [shiftStartTime, setShiftStartTime] = useState<string>('');
    const [shiftEndTime, setShiftEndTime] = useState<string>('');
    const [isAddingShift, setIsAddingShift] = useState(false);
    const [activeTab, setActiveTab] = useState<'event' | 'shifts'>('event');

    const fetchEventShifts = useCallback(async () => {
        if (!event || !clubId) return;
        
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/events/${event.id}/shifts`);
            setShifts(response.data || []);
        } catch (error) {
            console.error("Error fetching event shifts:", error);
        }
    }, [event, clubId]);

    useEffect(() => {
        if (event) {
            setName(event.name);
            setStartDate(event.start_date);
            setStartTime(event.start_time);
            setEndDate(event.end_date);
            setEndTime(event.end_time);
            fetchEventShifts();
        }
    }, [event, fetchEventShifts]);


    const handleAddShift = async () => {
        if (!shiftStartTime || !shiftEndTime || !event || !clubId) {
            setError("Please fill in both start and end times for the shift");
            return;
        }

        if (new Date(shiftStartTime) >= new Date(shiftEndTime)) {
            setError("Shift end time must be after start time");
            return;
        }

        setIsAddingShift(true);
        setError(null);
        
        try {
            await api.post(`/api/v1/clubs/${clubId}/events/${event.id}/shifts`, {
                startTime: shiftStartTime,
                endTime: shiftEndTime
            });
            setShiftStartTime('');
            setShiftEndTime('');
            await fetchEventShifts(); // Refresh the shifts list
        } catch (error: unknown) {
            if (error instanceof Error) {
                setError("Failed to add shift: " + error.message);
            } else {
                setError("Failed to add shift: Unknown error");
            }
        } finally {
            setIsAddingShift(false);
        }
    };

    if (!isOpen || !event) return null;

    const handleSubmit = async () => {
        if (!name || !startDate || !startTime || !endDate || !endTime) {
            setError("Please fill in all fields");
            return;
        }

        const startDateTime = new Date(`${startDate}T${startTime}`);
        const endDateTime = new Date(`${endDate}T${endTime}`);

        if (startDateTime >= endDateTime) {
            setError("End date/time must be after start date/time");
            return;
        }

        setError(null);
        setIsSubmitting(true);
        
        try {
            await api.put(`/api/v1/clubs/${clubId}/events/${event.id}`, { 
                name,
                start_date: startDate,
                start_time: startTime,
                end_date: endDate,
                end_time: endTime
            });
            onSuccess();
            onClose();
        } catch (error: unknown) {
            if (error instanceof Error) {
                setError("Failed to update event: " + error.message);
            } else {
                setError("Failed to update event: Unknown error");
            }
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleClose = () => {
        setError(null);
        onClose();
    };

    return (
        <div className="modal">
            <div className="modal-content" style={{ maxWidth: '800px', width: '90%' }}>
                <h2>Edit Event</h2>
                {error && <p style={{ color: 'red' }}>{error}</p>}
                
                {/* Tab Navigation */}
                <div className="tabs-container">
                    <nav className="tabs-nav">
                        <button 
                            className={`tab-button ${activeTab === 'event' ? 'active' : ''}`}
                            onClick={() => setActiveTab('event')}
                        >
                            Event Details
                        </button>
                        <button 
                            className={`tab-button ${activeTab === 'shifts' ? 'active' : ''}`}
                            onClick={() => setActiveTab('shifts')}
                        >
                            Shifts
                        </button>
                    </nav>

                    <div className="tab-content">
                        {/* Event Details Tab */}
                        <div className={`tab-panel ${activeTab === 'event' ? 'active' : ''}`}>
                            <div className="form-group">
                                <label htmlFor="eventName">Event Name</label>
                                <input
                                    id="eventName"
                                    type="text"
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                    placeholder="Event Name"
                                    disabled={isSubmitting}
                                />
                            </div>

                            <div className="form-group">
                                <label htmlFor="startDate">Start Date</label>
                                <input
                                    id="startDate"
                                    type="date"
                                    value={startDate}
                                    onChange={(e) => setStartDate(e.target.value)}
                                    disabled={isSubmitting}
                                />
                            </div>

                            <div className="form-group">
                                <label htmlFor="startTime">Start Time</label>
                                <input
                                    id="startTime"
                                    type="time"
                                    value={startTime}
                                    onChange={(e) => setStartTime(e.target.value)}
                                    disabled={isSubmitting}
                                />
                            </div>

                            <div className="form-group">
                                <label htmlFor="endDate">End Date</label>
                                <input
                                    id="endDate"
                                    type="date"
                                    value={endDate}
                                    onChange={(e) => setEndDate(e.target.value)}
                                    disabled={isSubmitting}
                                />
                            </div>

                            <div className="form-group">
                                <label htmlFor="endTime">End Time</label>
                                <input
                                    id="endTime"
                                    type="time"
                                    value={endTime}
                                    onChange={(e) => setEndTime(e.target.value)}
                                    disabled={isSubmitting}
                                />
                            </div>
                        </div>

                        {/* Shifts Tab */}
                        <div className={`tab-panel ${activeTab === 'shifts' ? 'active' : ''}`}>
                            <h3>Event Shifts</h3>
                            
                            {/* Add Shift Form */}
                            <div style={{ marginBottom: '20px', padding: '15px', border: '1px solid #ddd', borderRadius: '5px' }}>
                                <h4>Add New Shift</h4>
                                <div className="form-group">
                                    <label htmlFor="shiftStartTime">Shift Start Time</label>
                                    <input
                                        id="shiftStartTime"
                                        type="datetime-local"
                                        value={shiftStartTime}
                                        onChange={(e) => setShiftStartTime(e.target.value)}
                                        disabled={isAddingShift}
                                    />
                                </div>
                                <div className="form-group">
                                    <label htmlFor="shiftEndTime">Shift End Time</label>
                                    <input
                                        id="shiftEndTime"
                                        type="datetime-local"
                                        value={shiftEndTime}
                                        onChange={(e) => setShiftEndTime(e.target.value)}
                                        disabled={isAddingShift}
                                    />
                                </div>
                                <button 
                                    onClick={handleAddShift}
                                    className="button-accept"
                                    disabled={isAddingShift || !shiftStartTime || !shiftEndTime}
                                >
                                    {isAddingShift ? 'Adding...' : 'Add Shift'}
                                </button>
                            </div>

                            {/* Shifts List */}
                            <div>
                                <h4>Current Shifts</h4>
                                {shifts.length === 0 ? (
                                    <p style={{fontStyle: 'italic', color: '#666'}}>No shifts scheduled for this event yet.</p>
                                ) : (
                                    <table className="basic-table">
                                        <thead>
                                            <tr>
                                                <th>Start Time</th>
                                                <th>End Time</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {shifts.map(shift => (
                                                <tr key={shift.id}>
                                                    <td>{new Date(shift.startTime).toLocaleString()}</td>
                                                    <td>{new Date(shift.endTime).toLocaleString()}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                )}
                            </div>
                        </div>
                    </div>
                </div>

                <div className="modal-actions">
                    <button 
                        onClick={handleSubmit} 
                        className="button-accept"
                        disabled={isSubmitting}
                    >
                        {isSubmitting ? 'Updating...' : 'Update Event'}
                    </button>
                    <button 
                        onClick={handleClose} 
                        className="button-cancel"
                        disabled={isSubmitting}
                    >
                        Cancel
                    </button>
                </div>
            </div>
        </div>
    );
};

export default EditEvent;