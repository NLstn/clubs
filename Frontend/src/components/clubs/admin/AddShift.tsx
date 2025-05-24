import { FC, useState } from "react";
import api from "../../../utils/api";

interface AddShiftProps {
    isOpen: boolean;
    onClose: () => void;
    clubId: string;
    onSuccess: () => void;
}

const AddShift: FC<AddShiftProps> = ({ isOpen, onClose, clubId, onSuccess }) => {
    const [startTime, setStartTime] = useState<string>('');
    const [endTime, setEndTime] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    if (!isOpen) return null;

    const handleSubmit = async () => {
        if (!startTime || !endTime) {
            setError("Please fill in both start and end times");
            return;
        }

        if (new Date(startTime) >= new Date(endTime)) {
            setError("End time must be after start time");
            return;
        }

        setError(null);
        setIsSubmitting(true);
        
        try {
            await api.post(`/api/v1/clubs/${clubId}/shifts`, { 
                startTime, 
                endTime
            });
            setStartTime('');
            setEndTime('');
            onSuccess();
            onClose();
        } catch (error: unknown) {
            if (error instanceof Error) {
                setError("Failed to add shift: " + error.message);
            } else {
                setError("Failed to add shift: Unknown error");
            }
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="modal">
            <div className="modal-content">
                <h2>Add Shift</h2>
                {error && <div className="error">{error}</div>}
                <div className="form-group">
                    <label htmlFor="startTime">Start Time</label>
                    <input
                        id="startTime"
                        type="datetime-local"
                        value={startTime}
                        onChange={(e) => setStartTime(e.target.value)}
                        placeholder="Select start time"
                    />
                    <label htmlFor="endTime">End Time</label>
                    <input
                        id="endTime"
                        type="datetime-local"
                        value={endTime}
                        onChange={(e) => setEndTime(e.target.value)}
                        placeholder="Select end time"
                    />
                </div>
                <div>
                    <button 
                        onClick={handleSubmit} 
                        disabled={!startTime || !endTime || isSubmitting} 
                        className="button-accept"
                    >
                        {isSubmitting ? "Adding..." : "Add Shift"}
                    </button>
                    <button onClick={onClose} className="button-cancel">Cancel</button>
                </div>
            </div>
        </div>
    );
};

export default AddShift;
