import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Table, TableColumn } from '@/components/ui';
import api from "../../../../utils/api";

interface JoinRequest {
    id: string;
    email: string;
}

const AdminClubJoinRequestList = () => {
    const { id } = useParams();
    const [joinRequests, setJoinRequests] = useState<JoinRequest[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchJoinRequests = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await api.get(`/api/v1/clubs/${id}/joinRequests`);
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
                    <button 
                        onClick={() => handleApprove(request.id)}
                        className="action-button edit"
                        style={{marginRight: '8px'}}
                    >
                        Approve
                    </button>
                    <button 
                        onClick={() => handleReject(request.id)}
                        className="action-button remove"
                    >
                        Reject
                    </button>
                </div>
            )
        }
    ];

    const handleApprove = async (requestId: string) => {
        try {
            await api.post(`/api/v1/joinRequests/${requestId}/accept`);
            setJoinRequests(joinRequests.filter(request => request.id !== requestId));
        } catch (error) {
            console.error("Error approving join request:", error);
        }
    };

    const handleReject = async (requestId: string) => {
        try {
            await api.post(`/api/v1/joinRequests/${requestId}/reject`);
            setJoinRequests(joinRequests.filter(request => request.id !== requestId));
        } catch (error) {
            console.error("Error rejecting join request:", error);
        }
    };

    return (
        <div>
            <h3>Join Requests</h3>
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
