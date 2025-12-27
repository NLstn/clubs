import Layout from "../../components/layout/Layout";
import SimpleSettingsLayout from '../../components/layout/SimpleSettingsLayout';
import { ODataTable, ODataTableColumn } from '@/components/ui';
import './Profile.css';

interface Fine {
    ID: string;
    Amount: number;
    Reason: string;
    CreatedAt: string;
    UpdatedAt: string;
    Paid: boolean;
    Club?: {
        Name: string;
    };
}

const ProfileFines = () => {
    const columns: ODataTableColumn<Fine>[] = [
        {
            key: 'clubName',
            header: 'Club Name',
            render: (fine) => fine.Club?.Name || 'Unknown Club',
            sortable: true,
            sortField: 'Club/Name',
        },
        {
            key: 'Amount',
            header: 'Amount',
            render: (fine) => `€${fine.Amount.toFixed(2)}`,
            sortable: true,
        },
        {
            key: 'status',
            header: 'Status',
            render: (fine) => (
                <span className={fine.Paid ? 'status-paid' : 'status-unpaid'}>
                    {fine.Paid ? 'Paid' : 'Unpaid'}
                </span>
            ),
            sortable: true,
            sortField: 'Paid',
        },
        {
            key: 'CreatedAt',
            header: 'Date',
            render: (fine) => new Date(fine.CreatedAt).toLocaleDateString(),
            sortable: true,
        },
        {
            key: 'Reason',
            header: 'Reason',
            render: (fine) => fine.Reason,
            sortable: true,
        }
    ];

    return (
        <Layout title="Fines">
            <SimpleSettingsLayout title="Fines">
                <ODataTable
                    endpoint="/api/v2/Fines"
                    expand="Club"
                    columns={columns}
                    keyExtractor={(fine) => fine.ID}
                    pageSize={10}
                    initialSortField="CreatedAt"
                    initialSortDirection="desc"
                    emptyMessage="No fines found"
                    loadingMessage="Loading fines..."
                />
            </SimpleSettingsLayout>
        </Layout>
    )
}

export default ProfileFines;