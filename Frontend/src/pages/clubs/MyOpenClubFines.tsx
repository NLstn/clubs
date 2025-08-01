import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../../utils/api";
import Table, { TableColumn } from "../../components/ui/Table";

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
            const response = await api.get(`/api/v1/me/fines`);
            if (!response.data) {
                setFines([]);
                return;
            }
            const filteredFines = response.data.filter((fine: Fine) => fine.clubId === id);
            setFines(filteredFines);
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