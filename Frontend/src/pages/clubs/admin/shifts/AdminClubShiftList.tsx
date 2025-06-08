import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import EditShift from "./EditShift";
import api from "../../../../utils/api";

interface Shift {
    id: string;
    startTime: string;
    endTime: string;
    eventId: string;
}

interface Event {
    id: string;
    name: string;
}

const AdminClubShiftList = () => {

    const { id } = useParams();
    const [shifts, setShifts] = useState<Shift[]>([]);
    const [selectedShift, setSelectedShift] = useState<Shift | null>(null);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [events, setEvents] = useState<{[key: string]: string}>({});

    const fetchShifts = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await api.get(`/api/v1/clubs/${id}/shifts`);
            const shiftsData = response.data || [];
            setShifts(shiftsData);
            
            // Fetch event names for all shifts since they all have eventId now
            const eventIds = [...new Set(shiftsData.map((shift: Shift) => shift.eventId))];
            if (eventIds.length > 0) {
                try {
                    const eventsResponse = await api.get(`/api/v1/clubs/${id}/events`);
                    const eventsData = eventsResponse.data || [];
                    const eventMap: {[key: string]: string} = {};
                    eventsData.forEach((event: Event) => {
                        eventMap[event.id] = event.name;
                    });
                    setEvents(eventMap);
                } catch (err) {
                    console.warn("Failed to fetch event names:", err);
                }
            }
        } catch (error) {
            console.error("Error fetching shifts:", error);
            setError(error instanceof Error ? error.message : "Failed to fetch shifts");
            setShifts([]);
        } finally {
            setLoading(false);
        }
    }, [id]);

    useEffect(() => {
        fetchShifts();
    }, [fetchShifts]);

    const handleEditShift = (shift: Shift) => {
        setSelectedShift(shift);
        setIsEditModalOpen(true);
    };

    const handleCloseEditModal = () => {
        setSelectedShift(null);
        setIsEditModalOpen(false);
    };

    return (
        <div>
            <h3>Shifts</h3>
            {loading && <p>Loading shifts...</p>}
            {error && <p style={{color: 'red'}}>Error: {error}</p>}
            {!loading && !error && (
                <>
                    <table className="basic-table">
                        <thead>
                            <tr>
                                <th>Start Time</th>
                                <th>End Time</th>
                                <th>Event</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {shifts.length === 0 ? (
                                <tr>
                                    <td colSpan={4} style={{textAlign: 'center', fontStyle: 'italic'}}>
                                        No shifts available
                                    </td>
                                </tr>
                            ) : (
                                shifts.map(shift => (
                                    <tr key={shift.id}>
                                        <td>{new Date(shift.startTime).toLocaleString()}</td>
                                        <td>{new Date(shift.endTime).toLocaleString()}</td>
                                        <td>
                                            <span style={{color: 'blue'}}>
                                                {events[shift.eventId] || 'Unknown Event'}
                                            </span>
                                        </td>
                                        <td>
                                            <button
                                                onClick={() => handleEditShift(shift)}
                                                className="button-accept"
                                            >
                                                Edit
                                            </button>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                    <p style={{fontStyle: 'italic', color: '#666', marginTop: '10px'}}>
                        Shifts can only be created through events. Go to Events tab to create shifts.
                    </p>
                </>
            )}
            <EditShift
                isOpen={isEditModalOpen}
                onClose={handleCloseEditModal}
                shift={selectedShift}
                clubId={id}
            />
        </div>
    );
};

export default AdminClubShiftList;