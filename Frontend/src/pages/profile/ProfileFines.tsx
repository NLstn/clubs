import { useEffect, useState } from "react";
import api from "../../utils/api";
import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import Table, { TableColumn } from "../../components/ui/Table";

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
            const response = await api.get('/api/v1/me/fines');
            if (response.status === 200) {
                const data = response.data;
                setFines(data);
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

    return (
        <Layout title="Fines">
            <div style={{
                display: 'flex',
                minHeight: 'calc(100vh - 90px)',
                width: '100%',
                position: 'relative'
            }}>
                <ProfileSidebar />
                <div style={{
                    flex: '1 1 auto',
                    padding: '20px',
                    maxWidth: 'calc(100% - 200px)'
                }}>
                    <h1>Fines</h1>
                    <Table
                        columns={columns}
                        data={fines}
                        keyExtractor={(fine) => fine.id}
                        loading={loading}
                        error={error}
                        emptyMessage="No fines found"
                        loadingMessage="Loading fines..."
                        errorMessage="Failed to load fines"
                    />
                </div>
            </div>
        </Layout>
    )
}

export default ProfileFines;