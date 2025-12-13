import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../../utils/api";
import { Table, TableColumn } from '@/components/ui';

interface Fine {
    id: string;
    teamId: string;
    amount: number;
    reason: string;
    createdAt: string;
    updatedAt: string;
    paid: boolean;
    createdByName: string;
}

const TeamFines = () => {
    const { teamId } = useParams();
    const [fines, setFines] = useState<Fine[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchFines = useCallback(async () => {
        setLoading(true);
        try {
            // OData v2: Use Fines navigation property with User expansion
            const response = await api.get(`/api/v2/Teams('${teamId}')/Fines?$expand=User`);
            const finesData = response.data.value || response.data || [];
            // Map OData response to expected format
            interface ODataFine { ID: string; TeamID: string; Amount: number; Reason: string; CreatedAt: string; UpdatedAt: string; Paid: boolean; createdByName?: string; CreatedByName?: string; }
            const mappedFines = finesData.map((fine: ODataFine) => ({
                id: fine.ID,
                teamId: fine.TeamID,
                amount: fine.Amount,
                reason: fine.Reason,
                createdAt: fine.CreatedAt,
                updatedAt: fine.UpdatedAt,
                paid: fine.Paid,
                createdByName: fine.createdByName || fine.CreatedByName || 'Unknown'
            }));
            setFines(mappedFines);
            setError(null);
        } catch (err) {
            setError("Failed to fetch team fines: " + err);
            setFines([]);
        } finally {
            setLoading(false);
        }
    }, [teamId]);

    useEffect(() => {
        fetchFines();
    }, [fetchFines]);

    // Define table columns
    const columns: TableColumn<Fine>[] = [
        {
            key: 'reason',
            header: 'Reason',
            render: (fine) => <span>{fine.reason}</span>
        },
        {
            key: 'amount',
            header: 'Amount',
            render: (fine) => <span>${fine.amount.toFixed(2)}</span>
        },
        {
            key: 'createdAt',
            header: 'Created At',
            render: (fine) => <span>{new Date(fine.createdAt).toLocaleString()}</span>,
            className: 'hide-mobile'
        },
        {
            key: 'createdBy',
            header: 'Created By',
            render: (fine) => <span>{fine.createdByName}</span>,
            className: 'hide-small'
        }
    ];

    // Filter open fines and calculate total
    const openFines = fines.filter(fine => !fine.paid);
    const totalAmount = openFines.reduce((sum, fine) => sum + fine.amount, 0);

    return (
        <div className="content-section">
            <h3>My Team Fines</h3>
            <Table
                columns={columns}
                data={openFines}
                keyExtractor={(fine) => fine.id}
                loading={loading}
                error={error}
                emptyMessage="No open team fines"
                loadingMessage="Loading team fines..."
                errorMessage="Failed to load team fines"
                footer={
                    openFines.length > 0 ? (
                        <div>Total: ${totalAmount.toFixed(2)} across {openFines.length} fine{openFines.length !== 1 ? 's' : ''}</div>
                    ) : null
                }
            />
        </div>
    );
};

export default TeamFines;