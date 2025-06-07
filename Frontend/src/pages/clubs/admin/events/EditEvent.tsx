import { FC, useState, useEffect } from "react";
import api from "../../../../utils/api";

interface Event {
    id: string;
    name: string;
    start_date: string;
    start_time: string;
    end_date: string;
    end_time: string;
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

    useEffect(() => {
        if (event) {
            setName(event.name);
            setStartDate(event.start_date);
            setStartTime(event.start_time);
            setEndDate(event.end_date);
            setEndTime(event.end_time);
        }
    }, [event]);

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
            <div className="modal-content">
                <h2>Edit Event</h2>
                {error && <p style={{ color: 'red' }}>{error}</p>}
                
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