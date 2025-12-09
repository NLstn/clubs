import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../../utils/api";
import { Table, TableColumn } from '@/components/ui';

interface Fine {
    id: string;
    clubId: string;
    amount: number;
    reason: string;
    createdAt: string;
    updatedAt: string;
    paid: boolean;
    createdByName: string;
}

const MyOpenClubFines = () => {
    const { id } = useParams();
    const [fines, setFines] = useState<Fine[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchFines = useCallback(async () => {
        setLoading(true);
        try {
            // OData v2: Query Fines filtered by club ID, expand creator for name
            const response = await api.get(
                `/api/v2/Fines?$filter=ClubID eq '${id}'&$expand=CreatedByUser`
            );
            interface ODataFine { ID: string; ClubID: string; Amount: number; Reason: string; CreatedAt: string; UpdatedAt: string; Paid: boolean; CreatedByUser?: { FirstName: string; LastName: string; }; }
            const finesData = response.data.value || [];
            // Map OData response to match expected format
            const mappedFines = finesData.map((fine: ODataFine) => ({
                id: fine.ID,
                clubId: fine.ClubID,
                amount: fine.Amount,
                reason: fine.Reason,
                createdAt: fine.CreatedAt,
                updatedAt: fine.UpdatedAt,
                paid: fine.Paid,
                createdByName: fine.CreatedByUser ? 
                    `${fine.CreatedByUser.FirstName} ${fine.CreatedByUser.LastName}`.trim() : 
                    'Unknown'
            }));
            setFines(mappedFines);
            setError(null);
        } catch (err) {
            setError("Failed to fetch fines: " + err);
            setFines([]);
        } finally {
            setLoading(false);
        }
    }, [id]);

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
            <h3>My Open Fines</h3>
            <Table
                columns={columns}
                data={openFines}
                keyExtractor={(fine) => fine.id}
                loading={loading}
                error={error}
                emptyMessage="No open fines"
                loadingMessage="Loading fines..."
                errorMessage="Failed to load fines"
                footer={
                    openFines.length > 0 ? (
                        <div>Total: ${totalAmount.toFixed(2)} across {openFines.length} fine{openFines.length !== 1 ? 's' : ''}</div>
                    ) : null
                }
            />
        </div>
    );
};

export default MyOpenClubFines;