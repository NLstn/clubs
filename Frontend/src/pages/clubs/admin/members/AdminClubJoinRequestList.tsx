import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Table, TableColumn, Button } from '@/components/ui';
import api from "../../../../utils/api";

interface JoinRequest {
    id: string;
    email: string;
}

interface AdminClubJoinRequestListProps {
    onRequestsChange?: () => void;
}

const AdminClubJoinRequestList = ({ onRequestsChange }: AdminClubJoinRequestListProps) => {
    const { id } = useParams();
    const [joinRequests, setJoinRequests] = useState<JoinRequest[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchJoinRequests = async () => {
            setLoading(true);
            setError(null);
            try {
                // OData v2: Query JoinRequests for this club with expanded User data
                const response = await api.get(`/api/v2/JoinRequests?$filter=ClubID eq '${id}'&$expand=User`);
                setJoinRequests(response.data);
            } catch (error) {
                console.error("Error fetching join requests:", error);
                setError("Failed to load join requests");
            } finally {
                setLoading(false);
            }
        };
        fetchJoinRequests();
    }, [id]);

    // Define table columns
    const columns: TableColumn<JoinRequest>[] = [
        {
            key: 'email',
            header: 'Email',
            render: (request) => request.email
        },
        {
            key: 'actions',
            header: 'Actions',
            render: (request) => (
                <div className="table-actions">
                    <Button
                        size="sm"
                        variant="accept"
                        onClick={() => handleApprove(request.id)}
                        style={{marginRight: '8px'}}
                    >
                        Approve
                    </Button>
                    <Button
                        size="sm"
                        variant="cancel"
                        onClick={() => handleReject(request.id)}
                    >
                        Reject
                    </Button>
                </div>
            )
        }
    ];

    const handleApprove = async (requestId: string) => {
        try {
            // OData v2: Use Accept action on JoinRequest entity
            await api.post(`/api/v2/JoinRequests('${requestId}')/Accept`, {});
            setJoinRequests(joinRequests.filter(request => request.id !== requestId));
            onRequestsChange?.(); // Notify parent component about the change
        } catch (error) {
            console.error("Error approving join request:", error);
        }
    };

    const handleReject = async (requestId: string) => {
        try {
            // OData v2: Use Reject action on JoinRequest entity
            await api.post(`/api/v2/JoinRequests('${requestId}')/Reject`, {});
            setJoinRequests(joinRequests.filter(request => request.id !== requestId));
            onRequestsChange?.(); // Notify parent component about the change
        } catch (error) {
            console.error("Error rejecting join request:", error);
        }
    };

    return (
        <div>
            <p>People who want to join your club via invitation link:</p>
            <Table
                columns={columns}
                data={joinRequests}
                keyExtractor={(request) => request.id}
                loading={loading}
                error={error}
                emptyMessage="No pending join requests"
                loadingMessage="Loading join requests..."
                errorMessage="Failed to load join requests"
            />
        </div>
    )
}

export default AdminClubJoinRequestList;
