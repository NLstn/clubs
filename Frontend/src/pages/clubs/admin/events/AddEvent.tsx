import { FC, useState } from "react";
import api from "../../../../utils/api";

interface AddEventProps {
    isOpen: boolean;
    onClose: () => void;
    clubId: string;
    onSuccess: () => void;
}

const AddEvent: FC<AddEventProps> = ({ isOpen, onClose, clubId, onSuccess }) => {
    const [name, setName] = useState<string>('');
    const [startTime, setStartTime] = useState<string>('');
    const [endTime, setEndTime] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    if (!isOpen) return null;

    const handleSubmit = async () => {
        if (!name || !startTime || !endTime) {
            setError("Please fill in all fields");
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
            await api.post(`/api/v1/clubs/${clubId}/events`, { 
                name,
                start_time: startDateTime.toISOString(),
                end_time: endDateTime.toISOString()
            });
            setName('');
            setStartTime('');
            setEndTime('');
            onSuccess();
            onClose();
        } catch (error: unknown) {
            if (error instanceof Error) {
                setError("Failed to add event: " + error.message);
            } else {
                setError("Failed to add event: Unknown error");
            }
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleClose = () => {
        setName('');
        setStartTime('');
        setEndTime('');
        setError(null);
        onClose();
    };

    return (
        <div className="modal">
            <div className="modal-content">
                <h2>Add Event</h2>
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
                    <label htmlFor="startTime">Start Date & Time</label>
                    <input
                        id="startTime"
                        type="datetime-local"
                        value={startTime}
                        onChange={(e) => setStartTime(e.target.value)}
                        disabled={isSubmitting}
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="endTime">End Date & Time</label>
                    <input
                        id="endTime"
                        type="datetime-local"
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
                        {isSubmitting ? 'Adding...' : 'Add Event'}
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

export default AddEvent;