import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import api from "../../../../utils/api";
import AddFine from "./AddFine";
import AdminClubFineTemplateList from "./AdminClubFineTemplateList";

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

    const fetchFines = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${id}/fines`);
            setFines(response.data);
        } catch (err) {
            setError("Failed to fetch fines: " + err);
        }
    }, [id]);

    useEffect(() => {
        fetchFines();
    }, [fetchFines]);

    const displayedFines = showAllFines ? fines : fines.filter(fine => !fine.paid);

    return (
        <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 'var(--space-md)' }}>
                <h3>Fines</h3>
                <div style={{ display: 'flex', gap: 'var(--space-sm)' }}>
                    <label style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-xs)' }}>
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
                    </tr>
                </thead>
                <tbody>
                    {displayedFines && displayedFines.map((fine) => (
                        <tr key={fine.id}>
                            <td>{fine.userName}</td>
                            <td>{fine.amount}</td>
                            <td>{fine.reason}</td>
                            <td>{new Date(fine.createdAt).toLocaleString()}</td>
                            <td>{new Date(fine.updatedAt).toLocaleString()}</td>
                            <td>{fine.paid ? "Yes" : "No"}</td>
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