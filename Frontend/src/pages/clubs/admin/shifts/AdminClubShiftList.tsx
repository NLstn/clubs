import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import EditShift from "./EditShift";
import AddShift from "./AddShift";
import api from "../../../../utils/api";

interface Shift {
    id: string;
    startTime: string;
    endTime: string;
}

const AdminClubShiftList = () => {

    const { id } = useParams();
    const [shifts, setShifts] = useState<Shift[]>([]);
    const [selectedShift, setSelectedShift] = useState<Shift | null>(null);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isAddModalOpen, setIsAddModalOpen] = useState(false);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchShifts = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await api.get(`/api/v1/clubs/${id}/shifts`);
            setShifts(response.data || []);
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
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {shifts.length === 0 ? (
                                <tr>
                                    <td colSpan={3} style={{textAlign: 'center', fontStyle: 'italic'}}>
                                        No shifts available
                                    </td>
                                </tr>
                            ) : (
                                shifts.map(shift => (
                                    <tr key={shift.id}>
                                        <td>{new Date(shift.startTime).toLocaleString()}</td>
                                        <td>{new Date(shift.endTime).toLocaleString()}</td>
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
                    <button onClick={() => setIsAddModalOpen(true)} className="button-accept">
                        Add Shift
                    </button>
                </>
            )}
            <EditShift
                isOpen={isEditModalOpen}
                onClose={handleCloseEditModal}
                shift={selectedShift}
                clubId={id}
            />
            <AddShift 
                isOpen={isAddModalOpen}
                onClose={() => setIsAddModalOpen(false)}
                clubId={id || ''}
                onSuccess={fetchShifts}
            />
        </div>
    );
};

export default AdminClubShiftList;