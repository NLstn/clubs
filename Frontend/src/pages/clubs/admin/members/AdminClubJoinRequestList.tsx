import { useState } from "react";
import { useParams } from "react-router-dom";
import { ODataTable, ODataTableColumn, Button } from '@/components/ui';
import api from "../../../../utils/api";

interface JoinRequest {
    ID: string;
    Email: string;
}

interface AdminClubJoinRequestListProps {
    onRequestsChange?: () => void;
}

const AdminClubJoinRequestList = ({ onRequestsChange }: AdminClubJoinRequestListProps) => {
    const { id } = useParams();
    const [refreshKey, setRefreshKey] = useState(0);

    const handleApprove = async (requestId: string) => {
        try {
            // OData v2: Use Accept action on JoinRequest entity
            await api.post(`/api/v2/JoinRequests('${requestId}')/Accept`, {});
            setRefreshKey(prev => prev + 1); // Refresh the table
            onRequestsChange?.(); // Notify parent component about the change
        } catch (error) {
            console.error("Error approving join request:", error);
        }
    };

    const handleReject = async (requestId: string) => {
        try {
            // OData v2: Use Reject action on JoinRequest entity
            await api.post(`/api/v2/JoinRequests('${requestId}')/Reject`, {});
            setRefreshKey(prev => prev + 1); // Refresh the table
            onRequestsChange?.(); // Notify parent component about the change
        } catch (error) {
            console.error("Error rejecting join request:", error);
        }
    };

    // Define table columns
    const columns: ODataTableColumn<JoinRequest>[] = [
        {
            key: 'Email',
            header: 'Email',
            render: (request) => request.Email,
            sortable: true,
        },
        {
            key: 'actions',
            header: 'Actions',
            render: (request) => (
                <div className="table-actions">
                    <Button
                        size="sm"
                        variant="accept"
                        onClick={() => handleApprove(request.ID)}
                        style={{marginRight: '8px'}}
                    >
                        Approve
                    </Button>
                    <Button
                        size="sm"
                        variant="cancel"
                        onClick={() => handleReject(request.ID)}
                    >
                        Reject
                    </Button>
                </div>
            )
        }
    ];

    return (
        <div>
            <p>People who want to join your club via invitation link:</p>
            <ODataTable
                key={refreshKey}
                endpoint="/api/v2/JoinRequests"
                filter={`ClubID eq '${id}'`}
                expand="User"
                columns={columns}
                keyExtractor={(request) => request.ID}
                pageSize={10}
                emptyMessage="No pending join requests"
                loadingMessage="Loading join requests..."
            />
        </div>
    )
}

export default AdminClubJoinRequestList;
