import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../../../../utils/api";

interface Invite {
    id: string;
    email: string;
}

const AdminClubPendingInviteList = () => {

    const { id } = useParams();
    const [invites, setInvites] = useState<Invite[]>([]);

    useEffect(() => {
        const fetchInvites = async () => {
            try {
                const response = await api.get(`/api/v1/clubs/${id}/invites`);
                setInvites(response.data);
            } catch (error) {
                console.error("Error fetching invites:", error);
            }
        };
        fetchInvites();
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
                    {invites.map((invite) => (
                        <tr key={invite.id}>
                            <td>{invite.email}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

export default AdminClubPendingInviteList;