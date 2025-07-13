import { useEffect, useState, useCallback } from "react";
import InviteMember from "./InviteMember";
import AdminClubJoinRequestList from "./AdminClubJoinRequestList";
import AdminClubPendingInviteList from "./AdminClubPendingInviteList";
import api from "../../../../utils/api";
import { useParams } from "react-router-dom";
import { useT } from "../../../../hooks/useTranslation";
import { useCurrentUser } from "../../../../hooks/useCurrentUser";

interface Member {
    id: string;
    name: string;
    role: string;
    joinedAt: string;
    userId?: string; // Add userId to identify the current user
}

const AdminClubMemberList = () => {
    const { t } = useT();
    const { id } = useParams();
    const { user: currentUser } = useCurrentUser();

    const [members, setMembers] = useState<Member[]>([]);
    const [currentUserRole, setCurrentUserRole] = useState<string | null>(null);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [showPendingInvites, setShowPendingInvites] = useState(false);
    const [showJoinRequests, setShowJoinRequests] = useState(false);
    const [showInviteLink, setShowInviteLink] = useState(false);
    const [inviteLink, setInviteLink] = useState('');
    const [error, setError] = useState<string | null>(null);

    const translateRole = (role: string): string => {
        return t(`clubs.roles.${role}`) || role;
    };

    const sortMembersByRole = useCallback((members: Member[]): Member[] => {
        const roleOrder: { [key: string]: number } = { 
            'owner': 0, 
            'admin': 1, 
            'member': 2 
        };
        
        return [...members].sort((a, b) => {
            const aOrder = roleOrder[a.role.toLowerCase()] ?? 999;
            const bOrder = roleOrder[b.role.toLowerCase()] ?? 999;
            return aOrder - bOrder;
        });
    }, []);

    const fetchMembers = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${id}/members`);
            const sortedMembers = sortMembersByRole(response.data);
            setMembers(sortedMembers);

            // Set the current user's role if available
            const currentUserMember = response.data.find((member: Member) => member.userId === currentUser?.ID);
            if (currentUserMember) {
                setCurrentUserRole(currentUserMember.role);
            }
        } catch {
            setError("Failed to fetch members");
        }
    }, [id, sortMembersByRole, currentUser]);

    const handleShowPendingInvites = () => {
        setShowPendingInvites(true);
    };

    const handleShowInviteLink = async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${id}/inviteLink`);
            const fullLink = `${window.location.origin}${response.data.inviteLink}`;
            setInviteLink(fullLink);
            setShowInviteLink(true);
        } catch (error) {
            console.error("Error fetching invite link:", error);
            setError("Failed to generate invite link");
        }
    };

    const copyToClipboard = async () => {
        try {
            await navigator.clipboard.writeText(inviteLink);
            // You could add a toast notification here
            alert('Invite link copied to clipboard!');
        } catch (err) {
            console.error('Failed to copy link:', err);
            alert('Failed to copy link. Please copy manually.');
        }
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
            const sortedMembers = sortMembersByRole(updatedMembers);
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
            await api.post(`/api/v1/clubs/${id}/invites`, { email });
            setIsModalOpen(false);
        } catch {
            setError('Failed to send invite');
        }
    };

    // Permission logic based on backend rules (with desired admin permissions)
    const canChangeRole = (currentUserRole: string | null, targetMemberRole: string, newRole: string): boolean => {
        if (!currentUserRole) return false;
        
        // Owners can change any role to any role
        if (currentUserRole === 'owner') {
            return true;
        }
        
        // Admins can:
        // - Change members to any role (oldRole == "member")
        // - Change any role to admin (newRole == "admin") 
        // - Demote other admins to members (for better UX, even if backend might restrict this)
        // - BUT cannot touch owners (cannot demote or promote owners)
        if (currentUserRole === 'admin') {
            // Admins cannot change owner roles at all
            if (targetMemberRole === 'owner') {
                return false;
            }
            return targetMemberRole === 'member' || newRole === 'admin' || (targetMemberRole === 'admin' && newRole === 'member');
        }
        
        // Members cannot change roles
        return false;
    };

    // Permission logic for member deletion - admins can delete members, owners can delete anyone except other owners
    const canDeleteMember = (currentUserRole: string | null, targetMemberRole: string): boolean => {
        if (!currentUserRole) return false;
        
        // Owners can delete anyone except other owners
        if (currentUserRole === 'owner') {
            return targetMemberRole !== 'owner';
        }
        
        // Admins can delete regular members (but not other admins or owners)
        if (currentUserRole === 'admin') {
            return targetMemberRole === 'member';
        }
        
        // Regular members cannot delete anyone
        return false;
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
                            <td>{translateRole(member.role)}</td>
                            <td>{member.joinedAt ? new Date(member.joinedAt).toLocaleDateString() : 'N/A'}</td>
                            <td>
                                <div className="member-actions">
                                    {canDeleteMember(currentUserRole, member.role) && (
                                        <button
                                            onClick={() => deleteMember(member.id)}
                                            className="action-button remove"
                                            aria-label="Remove member"
                                        >
                                            Remove
                                        </button>
                                    )}
                                    {member.role === 'member' && canChangeRole(currentUserRole, member.role, 'admin') && (
                                        <button
                                            onClick={() => handleRoleChange(member.id, 'admin')}
                                            className="action-button promote"
                                        >
                                            Promote
                                        </button>
                                    )}
                                    {member.role === 'admin' && (
                                        <>
                                            {canChangeRole(currentUserRole, member.role, 'member') && (
                                                <button
                                                    onClick={() => handleRoleChange(member.id, 'member')}
                                                    className="action-button demote"
                                                >
                                                    Demote
                                                </button>
                                            )}
                                            {canChangeRole(currentUserRole, member.role, 'owner') && (
                                                <button
                                                    onClick={() => handleRoleChange(member.id, 'owner')}
                                                    className="action-button promote"
                                                >
                                                    Promote
                                                </button>
                                            )}
                                        </>
                                    )}
                                    {member.role === 'owner' && canChangeRole(currentUserRole, member.role, 'admin') && (
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
                <button onClick={handleShowInviteLink} className="button-accept">Generate Invite Link</button>
                <button onClick={handleShowPendingInvites}>View Pending Invites</button>
                <button onClick={() => setShowJoinRequests(true)}>View Join Requests</button>
            </div>
            <InviteMember
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onSubmit={sendInvite}
            />

            {showInviteLink && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Club Invitation Link</h2>
                        <p>Share this link with people you want to invite to the club:</p>
                        <div className="invite-link-container" style={{ 
                            background: '#f5f5f5', 
                            color: '#333',
                            padding: '10px', 
                            borderRadius: '5px', 
                            marginBottom: '15px',
                            wordBreak: 'break-all',
                            border: '1px solid #ddd'
                        }}>
                            {inviteLink}
                        </div>
                        <div className="modal-actions">
                            <button onClick={copyToClipboard} className="button-accept">
                                Copy Link
                            </button>
                            <button onClick={() => setShowInviteLink(false)} className="button-cancel">
                                Close
                            </button>
                        </div>
                        <div style={{ marginTop: '15px', fontSize: '0.9em', color: '#666' }}>
                            <p><strong>Note:</strong> Anyone with this link can request to join your club. 
                            Join requests still require admin approval.</p>
                        </div>
                    </div>
                </div>
            )}
            
            {showPendingInvites && (
                <div className="modal">
                    <div className="modal-content">
                        <AdminClubPendingInviteList />
                        <div className="modal-actions">
                            <button onClick={() => setShowPendingInvites(false)} className="button-cancel">Close</button>
                        </div>
                    </div>
                </div>
            )}

            {showJoinRequests && (
                <div className="modal">
                    <div className="modal-content">
                        <AdminClubJoinRequestList />
                        <div className="modal-actions">
                            <button onClick={() => setShowJoinRequests(false)} className="button-cancel">Close</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminClubMemberList;