import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import api from "../../../../utils/api";
import AddFine from "./AddFine";
import AdminClubFineTemplateList from "./AdminClubFineTemplateList";
import { formatDateTime } from "../../../../utils/dateHelpers";
import { createErrorHandler } from "../../../../utils/errorHandling";

interface Fine {
    id: string;
    userName: string;
    amount: number;
    reason: string;
    createdAt: string;
    updatedAt: string;
    paid: boolean;
}

const AdminClubFineList = () => {

    const { id } = useParams();

    const [fines, setFines] = useState<Fine[]>([]);
    const [showAllFines, setShowAllFines] = useState(false);
    const [showFineTemplates, setShowFineTemplates] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [isModalOpen, setIsModalOpen] = useState(false);

    const handleError = createErrorHandler("AdminClubFineList", setError, "Failed to perform operation");

    const fetchFines = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${id}/fines`);
            // Ensure we always have an array, even if API returns null/undefined
            setFines(Array.isArray(response.data) ? response.data : []);
            setError(null);
        } catch (err) {
            handleError(err);
            // Reset fines to empty array on error to prevent stale data issues
            setFines([]);
        }
    }, [id, handleError]);

    const handleDeleteFine = async (fineId: string) => {
        if (!confirm("Are you sure you want to delete this fine?")) {
            return;
        }

        try {
            await api.delete(`/api/v1/clubs/${id}/fines/${fineId}`);
            fetchFines(); // Refresh the list
        } catch (err) {
            handleError(err);
        }
    };

    useEffect(() => {
        fetchFines();
    }, [fetchFines]);

    const displayedFines = showAllFines ? fines : (fines || []).filter(fine => !fine.paid);

    return (
        <div>
            <div className="fines-header">
                <h3>Fines</h3>
                <div className="fines-controls">
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            checked={showAllFines}
                            onChange={(e) => setShowAllFines(e.target.checked)}
                        />
                        Show all fines
                    </label>
                    <button onClick={() => setShowFineTemplates(true)}>Manage Templates</button>
                </div>
            </div>
            {error && <div className="error">{error}</div>}
            <table>
                <thead>
                    <tr>
                        <th>User</th>
                        <th>Amount</th>
                        <th>Reason</th>
                        <th>Created At</th>
                        <th>Updated At</th>
                        <th>Paid</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {displayedFines && displayedFines.map((fine) => (
                        <tr key={fine.id}>
                            <td>{fine.userName}</td>
                            <td>{fine.amount}</td>
                            <td>{fine.reason}</td>
                            <td>{formatDateTime(fine.createdAt)}</td>
                            <td>{formatDateTime(fine.updatedAt)}</td>
                            <td>{fine.paid ? "Yes" : "No"}</td>
                            <td>
                                <button 
                                    onClick={() => handleDeleteFine(fine.id)}
                                    className="button-cancel"
                                >
                                    Delete
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <button onClick={() => setIsModalOpen(true)} className="button-accept">
                Add Fine
            </button>
            <AddFine 
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                clubId={id || ''}
                onSuccess={fetchFines}
            />
            
            {showFineTemplates && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Manage Fine Templates</h2>
                        <AdminClubFineTemplateList />
                        <div className="modal-actions">
                            <button onClick={() => setShowFineTemplates(false)} className="button-cancel">Close</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

export default AdminClubFineList;