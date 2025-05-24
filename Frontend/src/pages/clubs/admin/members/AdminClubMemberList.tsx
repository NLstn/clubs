import { useEffect, useState } from "react";
import InviteMember from "./InviteMember";
import api from "../../../../utils/api";
import { useParams } from "react-router-dom";

interface Member {
    id: string;
    name: string;
    role: string;
}

const AdminClubMemberList = () => {
    const { id } = useParams();

    const [members, setMembers] = useState<Member[]>([]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchMembers = async () => {
            try {
                const response = await api.get(`/api/v1/clubs/${id}/members`);
                setMembers(response.data);
            } catch {
                setError("Failed to fetch members");
            }
        };

        fetchMembers();
    }
    , [id]);

    const handleRoleChange = async (memberId: string, newRole: string) => {
        try {
            await api.patch(`/api/v1/clubs/${id}/members/${memberId}`, { role: newRole });
            setMembers(members.map(member => 
                member.id === memberId ? { ...member, role: newRole } : member
            ));
        } catch {
            setError('Failed to update member role');
        }
    };

    const deleteMember = async (memberId: string) => {
        try {
            await api.delete(`/api/v1/clubs/${id}/members/${memberId}`);
            setMembers(members.filter(member => member.id !== memberId));
        } catch {
            setError('Failed to delete member');
        }
    };

    const sendInvite = async (email: string) => {
        try {
            await api.post(`/api/v1/clubs/${id}/joinRequests`, { email });
            setIsModalOpen(false);
        } catch {
            setError('Failed to send invite');
        }
    };

    if (error) return <div className="error">{error}</div>;

    return (
        <div>
            <h3>Members</h3>
            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Role</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {members && members.map((member) => (
                        <tr key={member.id}>
                            <td>{member.name}</td>
                            <td>
                                <select
                                    value={member.role}
                                    onChange={(e) => handleRoleChange(member.id, e.target.value)}
                                >
                                    <option value="member">Member</option>
                                    <option value="admin">Admin</option>
                                    <option value="owner">Owner</option>
                                </select>
                            </td>
                            <td className="delete-cell">
                                <button
                                    onClick={() => deleteMember(member.id)}
                                    className="delete-button"
                                    aria-label="Delete member"
                                >
                                    ×
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <button onClick={() => setIsModalOpen(true)} className="button-accept">Invite Member</button>
            <InviteMember
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onSubmit={sendInvite}
            />
        </div>
    );
};

export default AdminClubMemberList;