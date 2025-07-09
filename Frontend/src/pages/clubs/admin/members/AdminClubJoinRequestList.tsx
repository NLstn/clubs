import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../../../../utils/api";

interface JoinRequest {
    id: string;
    email: string;
}

const AdminClubJoinRequestList = () => {
    const { id } = useParams();
    const [joinRequests, setJoinRequests] = useState<JoinRequest[]>([]);

    useEffect(() => {
        const fetchJoinRequests = async () => {
            try {
                const response = await api.get(`/api/v1/clubs/${id}/joinRequests`);
                setJoinRequests(response.data);
            } catch (error) {
                console.error("Error fetching join requests:", error);
            }
        };
        fetchJoinRequests();
    }, [id]);

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
            {joinRequests.length === 0 ? (
                <p>No pending join requests</p>
            ) : (
                <table className="basic-table">
                    <thead>
                        <tr>
                            <th>Email</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {joinRequests.map((request) => (
                            <tr key={request.id}>
                                <td>{request.email}</td>
                                <td>
                                    <button 
                                        onClick={() => handleApprove(request.id)}
                                        className="button-accept"
                                        style={{marginRight: '8px'}}
                                    >
                                        Approve
                                    </button>
                                    <button 
                                        onClick={() => handleReject(request.id)}
                                        className="button-cancel"
                                    >
                                        Reject
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    )
}

export default AdminClubJoinRequestList;
