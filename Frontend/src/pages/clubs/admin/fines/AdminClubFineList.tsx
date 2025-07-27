import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import api from "../../../../utils/api";
import AddFine from "./AddFine";
import AdminClubFineTemplateList from "./AdminClubFineTemplateList";
import Table, { TableColumn } from "../../../../components/ui/Table";

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
    const [loading, setLoading] = useState(false);

    const fetchFines = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await api.get(`/api/v1/clubs/${id}/fines`);
            // Ensure we always have an array, even if API returns null/undefined
            setFines(Array.isArray(response.data) ? response.data : []);
        } catch (err) {
            setError("Failed to fetch fines: " + err);
            // Reset fines to empty array on error to prevent stale data issues
            setFines([]);
        } finally {
            setLoading(false);
        }
    }, [id]);

    const handleDeleteFine = async (fineId: string) => {
        if (!confirm("Are you sure you want to delete this fine?")) {
            return;
        }

        try {
            await api.delete(`/api/v1/clubs/${id}/fines/${fineId}`);
            fetchFines(); // Refresh the list
        } catch (err) {
            setError("Failed to delete fine: " + err);
        }
    };

    useEffect(() => {
        fetchFines();
    }, [fetchFines]);

    const displayedFines = showAllFines ? fines : (fines || []).filter(fine => !fine.paid);
    
    // Calculate open fines statistics
    const openFines = fines.filter(fine => !fine.paid);
    const totalOpenFinesCount = openFines.length;
    const totalOpenFinesAmount = openFines.reduce((sum, fine) => sum + fine.amount, 0);

    const columns: TableColumn<Fine>[] = [
        {
            key: 'userName',
            header: 'User',
            render: (fine) => fine.userName
        },
        {
            key: 'amount',
            header: 'Amount',
            render: (fine) => fine.amount
        },
        {
            key: 'reason',
            header: 'Reason',
            render: (fine) => fine.reason
        },
        {
            key: 'createdAt',
            header: 'Created At',
            render: (fine) => new Date(fine.createdAt).toLocaleString()
        },
        {
            key: 'updatedAt',
            header: 'Updated At',
            render: (fine) => new Date(fine.updatedAt).toLocaleString()
        },
        {
            key: 'paid',
            header: 'Paid',
            render: (fine) => fine.paid ? "Yes" : "No"
        },
        {
            key: 'actions',
            header: 'Actions',
            render: (fine) => (
                <button 
                    onClick={() => handleDeleteFine(fine.id)}
                    className="button-cancel"
                >
                    Delete
                </button>
            )
        }
    ];

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
            <Table
                columns={columns}
                data={displayedFines}
                keyExtractor={(fine) => fine.id}
                loading={loading}
                error={error}
                emptyMessage="No fines available"
                loadingMessage="Loading fines..."
                errorMessage="Failed to load fines"
                footer={
                    <div>
                        <span>Open Fines: {totalOpenFinesCount}<br/></span>
                        <span>Total Amount: {totalOpenFinesAmount.toFixed(2)}</span>
                    </div>
                }
            />
            <div style={{ marginTop: '20px' }}>
                <button onClick={() => setIsModalOpen(true)} className="button-accept">
                    Add Fine
                </button>
            </div>
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