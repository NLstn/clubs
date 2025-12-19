import { useState, useMemo } from "react";
import { useParams } from "react-router-dom";
import { ODataTable, ODataTableColumn, Button, ButtonState } from '@/components/ui';
import api from "../../../../utils/api";
import { ODataFilter } from '../../../../utils/odata';

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
    const [processingRequests, setProcessingRequests] = useState<Record<string, ButtonState>>({});

    const handleApprove = async (requestId: string) => {
        setProcessingRequests(prev => ({ ...prev, [`approve-${requestId}`]: 'loading' }));
        
        try {
            // OData v2: Use Accept action on JoinRequest entity
            await api.post(`/api/v2/JoinRequests('${requestId}')/Accept`, {});
            setProcessingRequests(prev => ({ ...prev, [`approve-${requestId}`]: 'success' }));
            
            setTimeout(() => {
                setRefreshKey(prev => prev + 1); // Refresh the table
                onRequestsChange?.(); // Notify parent component about the change
                setProcessingRequests(prev => {
                    const newState = { ...prev };
                    delete newState[`approve-${requestId}`];
                    return newState;
                });
            }, 1000);
        } catch (error) {
            console.error("Error approving join request:", error);
            setProcessingRequests(prev => ({ ...prev, [`approve-${requestId}`]: 'error' }));
            setTimeout(() => {
                setProcessingRequests(prev => {
                    const newState = { ...prev };
                    delete newState[`approve-${requestId}`];
                    return newState;
                });
            }, 3000);
        }
    };

    const handleReject = async (requestId: string) => {
        setProcessingRequests(prev => ({ ...prev, [`reject-${requestId}`]: 'loading' }));
        
        try {
            // OData v2: Use Reject action on JoinRequest entity
            await api.post(`/api/v2/JoinRequests('${requestId}')/Reject`, {});
            setProcessingRequests(prev => ({ ...prev, [`reject-${requestId}`]: 'success' }));
            
            setTimeout(() => {
                setRefreshKey(prev => prev + 1); // Refresh the table
                onRequestsChange?.(); // Notify parent component about the change
                setProcessingRequests(prev => {
                    const newState = { ...prev };
                    delete newState[`reject-${requestId}`];
                    return newState;
                });
            }, 1000);
        } catch (error) {
            console.error("Error rejecting join request:", error);
            setProcessingRequests(prev => ({ ...prev, [`reject-${requestId}`]: 'error' }));
            setTimeout(() => {
                setProcessingRequests(prev => {
                    const newState = { ...prev };
                    delete newState[`reject-${requestId}`];
                    return newState;
                });
            }, 3000);
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
                        state={processingRequests[`approve-${request.ID}`] || 'idle'}
                        successMessage="Approved!"
                        errorMessage="Failed"
                    >
                        Approve
                    </Button>
                    <Button
                        size="sm"
                        variant="cancel"
                        onClick={() => handleReject(request.ID)}
                        state={processingRequests[`reject-${request.ID}`] || 'idle'}
                        successMessage="Rejected"
                        errorMessage="Failed"
                    >
                        Reject
                    </Button>
                </div>
            )
        }
    ];

    // Use ODataFilter helpers to safely escape values and prevent filter injection
    const filter = useMemo(() => ODataFilter.eq('ClubID', id || ''), [id]);

    return (
        <div>
            <p>People who want to join your club via invitation link:</p>
            <ODataTable
                key={refreshKey}
                endpoint="/api/v2/JoinRequests"
                filter={filter}
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
