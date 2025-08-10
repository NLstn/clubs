import { FC, useState, useEffect } from "react";
import api from "../../../../utils/api";
import { Input, Modal } from '@/components/ui';

interface Event {
    id: string;
    name: string;
    description: string;
    location: string;
    start_time: string;
    end_time: string;
}

interface EditEventProps {
    isOpen: boolean;
    onClose: () => void;
    clubId: string;
    teamId: string;
    event: Event | null;
    onSuccess: () => void;
}

const EditEvent: FC<EditEventProps> = ({ isOpen, onClose, clubId, teamId, event, onSuccess }) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [location, setLocation] = useState('');
    const [startTime, setStartTime] = useState('');
    const [endTime, setEndTime] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        if (event) {
            setName(event.name || '');
            setDescription(event.description || '');
            setLocation(event.location || '');
            try {
                setStartTime(new Date(event.start_time).toISOString().slice(0,16));
                setEndTime(new Date(event.end_time).toISOString().slice(0,16));
            } catch {
                setStartTime('');
                setEndTime('');
            }
        }
    }, [event]);

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
            await api.put(`/api/v1/clubs/${clubId}/teams/${teamId}/events/${event.id}`, {
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
        onClose();
    };

    return (
        <Modal isOpen={isOpen} onClose={handleClose} title="Edit Event">
            <Modal.Error error={error} />
            <Modal.Body>
                <div className="modal-form-section">
                    <Input
                        label="Event Name"
                        id="eventName"
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        disabled={isSubmitting}
                    />
                    <Input
                        label="Description"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
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
                <button onClick={handleSubmit} className="button-accept" disabled={isSubmitting}>{isSubmitting ? (<><Modal.LoadingSpinner />Saving...</>) : 'Save Changes'}</button>
                <button onClick={handleClose} className="button-cancel" disabled={isSubmitting}>Cancel</button>
            </Modal.Actions>
        </Modal>
    );
};

export default EditEvent;
