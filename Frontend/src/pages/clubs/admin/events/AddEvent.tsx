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
    const [isRecurring, setIsRecurring] = useState<boolean>(false);
    const [recurrencePattern, setRecurrencePattern] = useState<string>('weekly');
    const [recurrenceInterval, setRecurrenceInterval] = useState<number>(1);
    const [recurrenceEnd, setRecurrenceEnd] = useState<string>('');
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

        if (isRecurring && !recurrenceEnd) {
            setError("Please specify when recurring events should end");
            return;
        }

        if (isRecurring) {
            const recurrenceEndDateTime = new Date(recurrenceEnd);
            if (recurrenceEndDateTime <= startDateTime) {
                setError("Recurrence end must be after start date/time");
                return;
            }
        }

        setError(null);
        setIsSubmitting(true);
        
        try {
            const endpoint = isRecurring 
                ? `/api/v1/clubs/${clubId}/events/recurring`
                : `/api/v1/clubs/${clubId}/events`;

            const payload = {
                name,
                description,
                location,
                start_time: startDateTime.toISOString(),
                end_time: endDateTime.toISOString(),
                ...(isRecurring && {
                    recurrence_pattern: recurrencePattern,
                    recurrence_interval: recurrenceInterval,
                    recurrence_end: new Date(recurrenceEnd).toISOString()
                })
            };

            await api.post(endpoint, payload);
            
            // Reset form
            setName('');
            setDescription('');
            setLocation('');
            setStartTime('');
            setEndTime('');
            setIsRecurring(false);
            setRecurrencePattern('weekly');
            setRecurrenceInterval(1);
            setRecurrenceEnd('');
            
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
        setIsRecurring(false);
        setRecurrencePattern('weekly');
        setRecurrenceInterval(1);
        setRecurrenceEnd('');
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

                    <div className="modal-form-section">
                        <label className="checkbox-label">
                            <input
                                type="checkbox"
                                checked={isRecurring}
                                onChange={(e) => setIsRecurring(e.target.checked)}
                                disabled={isSubmitting}
                            />
                            Make this a recurring event
                        </label>
                    </div>

                    {isRecurring && (
                        <>
                            <div className="modal-form-section">
                                <label htmlFor="recurrencePattern">Recurrence Pattern</label>
                                <select
                                    id="recurrencePattern"
                                    value={recurrencePattern}
                                    onChange={(e) => setRecurrencePattern(e.target.value)}
                                    disabled={isSubmitting}
                                >
                                    <option value="daily">Daily</option>
                                    <option value="weekly">Weekly</option>
                                    <option value="monthly">Monthly</option>
                                </select>
                            </div>

                            <Input
                                label={`Every ${recurrenceInterval} ${recurrencePattern === 'daily' ? 'day(s)' : recurrencePattern === 'weekly' ? 'week(s)' : 'month(s)'}`}
                                id="recurrenceInterval"
                                type="number"
                                min="1"
                                max="52"
                                value={recurrenceInterval.toString()}
                                onChange={(e) => setRecurrenceInterval(parseInt(e.target.value) || 1)}
                                disabled={isSubmitting}
                            />

                            <Input
                                label="Stop Recurring After"
                                id="recurrenceEnd"
                                type="date"
                                value={recurrenceEnd}
                                onChange={(e) => setRecurrenceEnd(e.target.value)}
                                disabled={isSubmitting}
                            />
                        </>
                    )}
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