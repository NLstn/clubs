import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../../../../utils/api";

interface JoinRequest {
    id: string;
    email: string;
}

const AdminClubPendingInviteList = () => {

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

    return (
        <div>
            <h3>Pending Invites</h3>
            <table className="basic-table">
                <thead>
                    <tr>
                        <th>Email</th>
                    </tr>
                </thead>
                <tbody>
                    {joinRequests.map((request) => (
                        <tr key={request.id}>
                            <td>{request.email}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

export default AdminClubPendingInviteList;