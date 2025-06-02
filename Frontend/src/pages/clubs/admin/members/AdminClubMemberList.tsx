import { useEffect, useState, useCallback } from "react";
import InviteMember from "./InviteMember";
import api from "../../../../utils/api";
import { useParams } from "react-router-dom";

interface Member {
    id: string;
    name: string;
    role: string;
    joinedAt: string;
}

interface JoinRequest {
    id: string;
    email: string;
}

const AdminClubMemberList = () => {
    const { id } = useParams();

    const [members, setMembers] = useState<Member[]>([]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [showPendingInvites, setShowPendingInvites] = useState(false);
    const [joinRequests, setJoinRequests] = useState<JoinRequest[]>([]);
    const [error, setError] = useState<string | null>(null);

    const fetchMembers = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${id}/members`);
            const sortedMembers = response.data.sort((a: Member, b: Member) => {
                const roleOrder = { owner: 0, admin: 1, member: 2 };
                return (roleOrder[a.role as keyof typeof roleOrder] || 2) - (roleOrder[b.role as keyof typeof roleOrder] || 2);
            });
            setMembers(sortedMembers);
        } catch {
            setError("Failed to fetch members");
        }
    }, [id]);

    const fetchJoinRequests = async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${id}/joinRequests`);
            setJoinRequests(response.data);
        } catch (error) {
            console.error("Error fetching join requests:", error);
        }
    };

    const handleShowPendingInvites = () => {
        fetchJoinRequests();
        setShowPendingInvites(true);
    };

    useEffect(() => {
        fetchMembers();
    }, [fetchMembers]);

    const handleRoleChange = async (memberId: string, newRole: string) => {
        try {
            await api.patch(`/api/v1/clubs/${id}/members/${memberId}`, { role: newRole });
            const updatedMembers = members.map(member => 
                member.id === memberId ? { ...member, role: newRole } : member
            );
            const sortedMembers = updatedMembers.sort((a: Member, b: Member) => {
                const roleOrder = { owner: 0, admin: 1, member: 2 };
                return (roleOrder[a.role as keyof typeof roleOrder] || 2) - (roleOrder[b.role as keyof typeof roleOrder] || 2);
            });
            setMembers(sortedMembers);
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
                        <th>Joined</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {members && members.map((member) => (
                        <tr key={member.id}>
                            <td>{member.name}</td>
                            <td>{member.role}</td>
                            <td>{member.joinedAt ? new Date(member.joinedAt).toLocaleDateString() : 'N/A'}</td>
                            <td>
                                <div className="member-actions">
                                    {member.role !== 'owner' && (
                                        <button
                                            onClick={() => deleteMember(member.id)}
                                            className="action-button remove"
                                            aria-label="Remove member"
                                        >
                                            Remove
                                        </button>
                                    )}
                                    {member.role === 'member' && (
                                        <button
                                            onClick={() => handleRoleChange(member.id, 'admin')}
                                            className="action-button promote"
                                        >
                                            Promote
                                        </button>
                                    )}
                                    {member.role === 'admin' && (
                                        <>
                                            <button
                                                onClick={() => handleRoleChange(member.id, 'member')}
                                                className="action-button demote"
                                            >
                                                Demote
                                            </button>
                                            <button
                                                onClick={() => handleRoleChange(member.id, 'owner')}
                                                className="action-button promote"
                                            >
                                                Promote
                                            </button>
                                        </>
                                    )}
                                    {member.role === 'owner' && (
                                        <button
                                            onClick={() => handleRoleChange(member.id, 'admin')}
                                            className="action-button demote"
                                        >
                                            Demote
                                        </button>
                                    )}
                                </div>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <div className="member-actions buttons" style={{ marginTop: 'var(--space-md)' }}>
                <button onClick={() => setIsModalOpen(true)} className="button-accept">Invite Member</button>
                <button onClick={handleShowPendingInvites}>View Pending Invites</button>
            </div>
            <InviteMember
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onSubmit={sendInvite}
            />
            
            {showPendingInvites && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Pending Invites</h2>
                        {joinRequests.length === 0 ? (
                            <p>No pending invites</p>
                        ) : (
                            <table>
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
                        )}
                        <div className="modal-actions">
                            <button onClick={() => setShowPendingInvites(false)} className="button-cancel">Close</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminClubMemberList;