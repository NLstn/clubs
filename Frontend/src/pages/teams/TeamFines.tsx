import { useMemo } from "react";
import { useParams } from "react-router-dom";
import { ODataTable, ODataTableColumn } from '@/components/ui';

interface Fine {
    ID: string;
    TeamID: string;
    Amount: number;
    Reason: string;
    CreatedAt: string;
    UpdatedAt: string;
    Paid: boolean;
    User?: {
        FirstName: string;
        LastName: string;
    };
}

const TeamFines = () => {
    const { teamId } = useParams();

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
                {fine.User 
                    ? `${fine.User.FirstName} ${fine.User.LastName}`.trim() 
                    : 'Unknown'}
            </span>,
            className: 'hide-small',
            sortable: true,
            sortField: 'User/FirstName',
        }
    ];

    const filter = useMemo(() => {
        return `TeamID eq '${teamId}' and Paid eq false`;
    }, [teamId]);

    return (
        <div className="content-section">
            <h3>My Team Fines</h3>
            <ODataTable
                endpoint="/api/v2/Fines"
                filter={filter}
                expand="User"
                columns={columns}
                keyExtractor={(fine) => fine.ID}
                pageSize={10}
                initialSortField="CreatedAt"
                initialSortDirection="desc"
                emptyMessage="No open team fines"
                loadingMessage="Loading team fines..."
            />
        </div>
    );
};

export default TeamFines;