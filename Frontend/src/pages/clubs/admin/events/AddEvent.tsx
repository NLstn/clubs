import { FC, useState } from "react";
import api from "../../../../utils/api";
import { Input, Modal } from '@/components/ui';

interface AddEventProps {
    isOpen: boolean;
    onClose: () => void;
    clubId: string;
    onSuccess: () => void;
}

const AddEvent: FC<AddEventProps> = ({ isOpen, onClose, clubId, onSuccess }) => {
    const [name, setName] = useState<string>('');
    const [description, setDescription] = useState<string>('');
    const [location, setLocation] = useState<string>('');
    const [startTime, setStartTime] = useState<string>('');
    const [endTime, setEndTime] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    if (!isOpen) return null;

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
            await api.post(`/api/v1/clubs/${clubId}/events`, { 
                name,
                description,
                location,
                start_time: startDateTime.toISOString(),
                end_time: endDateTime.toISOString()
            });
            setName('');
            setDescription('');
            setLocation('');
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
        setDescription('');
        setLocation('');
        setStartTime('');
        setEndTime('');
        setError(null);
        onClose();
    };

    return (
        <Modal isOpen={isOpen} onClose={handleClose} title="Add Event">
            <Modal.Error error={error} />
            
            <Modal.Body>
                <div className="modal-form-section">
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
            </Modal.Body>

            <Modal.Actions>
                <button 
                    onClick={handleSubmit} 
                    className="button-accept"
                    disabled={isSubmitting}
                >
                    {isSubmitting ? (
                        <>
                            <Modal.LoadingSpinner />
                            Adding...
                        </>
                    ) : (
                        'Add Event'
                    )}
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

export default AddEvent;