import { useEffect, useState } from "react";
import api from "../../utils/api";
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { Table, TableColumn } from '@/components/ui';
import './Profile.css';

interface Fine {
    id: string;
    clubName: string;
    amount: number;
    reason: string;
    createdAt: string;
    updatedAt: string;
    paid: boolean;
}

const ProfileFines = () => {
    const [fines, setFines] = useState<Fine[]>([]);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        fetchFines();
    }, []);

    const fetchFines = async () => {
        setLoading(true);
        setError('');
        try {
            // OData v2: Query Fines with expand Club for club name, filtered to current user
            const response = await api.get('/api/v2/Fines?$expand=Club&$filter=Paid eq false or Paid eq true');
            if (response.status === 200) {
                interface ODataFine { ID: string; Amount: number; Reason: string; CreatedAt: string; UpdatedAt: string; Paid: boolean; Club?: { Name: string; }; }
                const data = response.data.value || [];
                // Map OData response to match expected format
                const mappedFines = data.map((fine: ODataFine) => ({
                    id: fine.ID,
                    clubName: fine.Club?.Name || 'Unknown Club',
                    amount: fine.Amount || 0,
                    reason: fine.Reason || '',
                    createdAt: fine.CreatedAt,
                    updatedAt: fine.UpdatedAt,
                    paid: fine.Paid || false
                }));
                setFines(mappedFines);
            }
        } catch (error) {
            setError('Error fetching fines: ' + error);
        } finally {
            setLoading(false);
        }
    };

    const columns: TableColumn<Fine>[] = [
        {
            key: 'clubName',
            header: 'Club Name',
            render: (fine) => fine.clubName
        },
        {
            key: 'amount',
            header: 'Amount',
            render: (fine) => `€${fine.amount.toFixed(2)}`
        },
        {
            key: 'status',
            header: 'Status',
            render: (fine) => (
                <span className={fine.paid ? 'status-paid' : 'status-unpaid'}>
                    {fine.paid ? 'Paid' : 'Unpaid'}
                </span>
            )
        },
        {
            key: 'date',
            header: 'Date',
            render: (fine) => new Date(fine.createdAt).toLocaleDateString()
        },
        {
            key: 'reason',
            header: 'Reason',
            render: (fine) => fine.reason
        }
    ];

    // Calculate total amount
    const totalAmount = fines.reduce((sum, fine) => sum + fine.amount, 0);

    // Create footer content
    const footerContent = fines.length > 0 ? `Total: €${totalAmount.toFixed(2)}` : null;

    return (
        <Layout title="Fines">
            <ProfileContentLayout title="Fines">
                <Table
                    columns={columns}
                    data={fines}
                    keyExtractor={(fine) => fine.id}
                    loading={loading}
                    error={error}
                    emptyMessage="No fines found"
                    loadingMessage="Loading fines..."
                    errorMessage="Failed to load fines"
                    footer={footerContent}
                />
            </ProfileContentLayout>
        </Layout>
    )
}

export default ProfileFines;