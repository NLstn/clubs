import { FC, useState, useEffect, useCallback } from "react";
import { Table, TableColumn, Input } from '@/components/ui';
import Modal from '@/components/ui/Modal';
import api from "../../../../utils/api";
import { useClubSettings } from "../../../../hooks/useClubSettings";

interface Event {
    id: string;
    name: string;
    description: string;
    location: string;
    start_time: string;
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
    const [description, setDescription] = useState<string>('');
    const [location, setLocation] = useState<string>('');
    const [startTime, setStartTime] = useState<string>('');
    const [endTime, setEndTime] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { settings: clubSettings, refetch: refetchClubSettings } = useClubSettings(clubId);
    
    // Shift management state
    const [shifts, setShifts] = useState<Shift[]>([]);
    const [shiftStartTime, setShiftStartTime] = useState<string>('');
    const [shiftEndTime, setShiftEndTime] = useState<string>('');
    const [isAddingShift, setIsAddingShift] = useState(false);
    const [activeTab, setActiveTab] = useState<'event' | 'shifts'>('event');

    // Define table columns for shifts
    const shiftColumns: TableColumn<Shift>[] = [
        {
            key: 'startTime',
            header: 'Start Time',
            render: (shift) => new Date(shift.startTime).toLocaleString()
        },
        {
            key: 'endTime',
            header: 'End Time',
            render: (shift) => new Date(shift.endTime).toLocaleString()
        }
    ];

    const fetchEventShifts = useCallback(async () => {
        if (!event || !clubId || !clubSettings?.shiftsEnabled) return;
        
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/events/${event.id}/shifts`);
            setShifts(response.data || []);
        } catch (error) {
            console.error("Error fetching event shifts:", error);
        }
    }, [event, clubId, clubSettings?.shiftsEnabled]);

    useEffect(() => {
        if (event) {
            setName(event.name || '');
            setDescription(event.description || '');
            setLocation(event.location || '');
            
            // Convert timestamps to datetime-local format (YYYY-MM-DDTHH:MM)
            try {
                if (event.start_time) {
                    const startDate = new Date(event.start_time);
                    if (!isNaN(startDate.getTime())) {
                        setStartTime(startDate.toISOString().slice(0, 16));
                    } else {
                        setStartTime('');
                    }
                } else {
                    setStartTime('');
                }

                if (event.end_time) {
                    const endDate = new Date(event.end_time);
                    if (!isNaN(endDate.getTime())) {
                        setEndTime(endDate.toISOString().slice(0, 16));
                    } else {
                        setEndTime('');
                    }
                } else {
                    setEndTime('');
                }
                
                setError(null);
            } catch {
                setError("This event has invalid date/time data. Please enter valid dates and times to update.");
                setStartTime('');
                setEndTime('');
            }
            
            fetchEventShifts();
        }
    }, [event, fetchEventShifts]);

    // Refresh settings when modal opens to get latest settings
    useEffect(() => {
        if (isOpen) {
            refetchClubSettings();
        }
    }, [isOpen, refetchClubSettings]);

    // Reset to event tab if shifts become unavailable
    useEffect(() => {
        if (clubSettings && activeTab === 'shifts' && !clubSettings.shiftsEnabled) {
            setActiveTab('event');
        }
    }, [clubSettings, activeTab]);


    const handleAddShift = async () => {
        if (!shiftStartTime || !shiftEndTime || !event || !clubId) {
            setError("Please fill in both start and end times for the shift");
            return;
        }

        const shiftStart = new Date(shiftStartTime);
        const shiftEnd = new Date(shiftEndTime);

        if (shiftStart >= shiftEnd) {
            setError("Shift end time must be after start time");
            return;
        }

        setIsAddingShift(true);
        setError(null);
        
        try {
            await api.post(`/api/v1/clubs/${clubId}/events/${event.id}/shifts`, {
                startTime: shiftStart.toISOString(),
                endTime: shiftEnd.toISOString()
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
        if (!name || !startTime || !endTime) {
            setError("Please fill in all required fields (name, start time, and end time)");
            return;
        }

        const startDateTime = new Date(startTime);
        const endDateTime = new Date(endTime);

        if (startDateTime >= endDateTime) {
            setError("End date/time must be after start date/time");
            return;
        }

        setError(null);
        setIsSubmitting(true);
        
        try {
            await api.put(`/api/v1/clubs/${clubId}/events/${event.id}`, { 
                name,
                description,
                location,
                start_time: startDateTime.toISOString(),
                end_time: endDateTime.toISOString()
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
        setActiveTab('event'); // Reset to event tab
        onClose();
    };

    return (
        <Modal 
            isOpen={isOpen} 
            onClose={handleClose}
            title="Edit Event"
            maxWidth="800px"
        >
            <Modal.Body>
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
                        {clubSettings?.shiftsEnabled && (
                            <button 
                                className={`tab-button ${activeTab === 'shifts' ? 'active' : ''}`}
                                onClick={() => setActiveTab('shifts')}
                            >
                                Shifts
                            </button>
                        )}
                    </nav>

                    <div className="tab-content">
                        {/* Event Details Tab */}
                        <div className={`tab-panel ${activeTab === 'event' ? 'active' : ''}`}>
                            <Input
                                label="Event Name"
                                id="eventName"
                                type="text"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                placeholder="Event Name"
                                disabled={isSubmitting}
                            />

                            <Input
                                label="Description"
                                value={description}
                                onChange={(e) => setDescription(e.target.value)}
                                placeholder="Event description (optional)"
                                disabled={isSubmitting}
                                multiline
                                rows={3}
                            />

                            <Input
                                label="Location"
                                id="eventLocation"
                                type="text"
                                value={location}
                                onChange={(e) => setLocation(e.target.value)}
                                placeholder="Event location (optional)"
                                disabled={isSubmitting}
                            />

                            <Input
                                label="Start Date & Time"
                                id="startTime"
                                type="datetime-local"
                                value={startTime}
                                onChange={(e) => setStartTime(e.target.value)}
                                disabled={isSubmitting}
                            />

                            <Input
                                label="End Date & Time"
                                id="endTime"
                                type="datetime-local"
                                value={endTime}
                                onChange={(e) => setEndTime(e.target.value)}
                                disabled={isSubmitting}
                            />
                        </div>

                        {/* Shifts Tab */}
                        {clubSettings?.shiftsEnabled && (
                            <div className={`tab-panel ${activeTab === 'shifts' ? 'active' : ''}`}>
                                <h3>Event Shifts</h3>
                                
                                {/* Add Shift Form */}
                                <div style={{ marginBottom: '20px', padding: '15px', border: '1px solid #ddd', borderRadius: '5px' }}>
                                    <h4>Add New Shift</h4>
                                    <Input
                                        label="Shift Start Time"
                                        id="shiftStartTime"
                                        type="datetime-local"
                                        value={shiftStartTime}
                                        onChange={(e) => setShiftStartTime(e.target.value)}
                                        disabled={isAddingShift}
                                    />
                                    <Input
                                        label="Shift End Time"
                                        id="shiftEndTime"
                                        type="datetime-local"
                                        value={shiftEndTime}
                                        onChange={(e) => setShiftEndTime(e.target.value)}
                                        disabled={isAddingShift}
                                    />
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
                                    <Table
                                        columns={shiftColumns}
                                        data={shifts}
                                        keyExtractor={(shift) => shift.id}
                                        emptyMessage="No shifts scheduled for this event yet."
                                    />
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </Modal.Body>
            <Modal.Actions>
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
            </Modal.Actions>
        </Modal>
    );
};

export default EditEvent;