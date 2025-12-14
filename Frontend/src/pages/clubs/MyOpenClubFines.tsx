import { useMemo } from "react";
import { useParams } from "react-router-dom";
import { ODataTable, ODataTableColumn } from '@/components/ui';
import { ODataFilter } from '../../utils/odata';

interface Fine {
    ID: string;
    ClubID: string;
    Amount: number;
    Reason: string;
    CreatedAt: string;
    UpdatedAt: string;
    Paid: boolean;
    CreatedByUser?: {
        FirstName: string;
        LastName: string;
    };
}

const MyOpenClubFines = () => {
    const { id } = useParams();

    // Define table columns
    const columns: ODataTableColumn<Fine>[] = [
        {
            key: 'Reason',
            header: 'Reason',
            render: (fine) => <span>{fine.Reason}</span>,
            sortable: true,
        },
        {
            key: 'Amount',
            header: 'Amount',
            render: (fine) => <span>${fine.Amount.toFixed(2)}</span>,
            sortable: true,
        },
        {
            key: 'CreatedAt',
            header: 'Created At',
            render: (fine) => <span>{new Date(fine.CreatedAt).toLocaleString()}</span>,
            className: 'hide-mobile',
            sortable: true,
        },
        {
            key: 'createdBy',
            header: 'Created By',
            render: (fine) => <span>
                {fine.CreatedByUser 
                    ? `${fine.CreatedByUser.FirstName} ${fine.CreatedByUser.LastName}`.trim() 
                    : 'Unknown'}
            </span>,
            className: 'hide-small',
            sortable: true,
            sortField: 'CreatedByUser/FirstName',
        }
    ];

    const filter = useMemo(() => {
        // Use ODataFilter helpers to safely escape values and prevent filter injection
        return ODataFilter.and(
            ODataFilter.eq('ClubID', id || ''),
            ODataFilter.eq('Paid', false)
        );
    }, [id]);

    return (
        <div className="content-section">
            <h3>My Open Fines</h3>
            <ODataTable
                endpoint="/api/v2/Fines"
                filter={filter}
                expand="CreatedByUser"
                columns={columns}
                keyExtractor={(fine) => fine.ID}
                pageSize={10}
                initialSortField="CreatedAt"
                initialSortDirection="desc"
                emptyMessage="No open fines"
                loadingMessage="Loading fines..."
            />
        </div>
    );
};

export default MyOpenClubFines;