import { useParams } from "react-router-dom";
import { ODataTable, ODataTableColumn } from '@/components/ui';

interface Invite {
    ID: string;
    Email: string;
}

const AdminClubPendingInviteList = () => {

    const { id } = useParams();

    // Define table columns
    const columns: ODataTableColumn<Invite>[] = [
        {
            key: 'Email',
            header: 'Email',
            render: (invite) => invite.Email,
            sortable: true,
        }
    ];

    return (
        <div>
            <ODataTable
                endpoint="/api/v2/Invites"
                filter={`ClubID eq '${id}'`}
                columns={columns}
                keyExtractor={(invite) => invite.ID}
                pageSize={10}
                emptyMessage="No pending invites found"
                loadingMessage="Loading pending invites..."
            />
        </div>
    )
}

export default AdminClubPendingInviteList;