import { useEffect, useState, useCallback } from "react";
import AdminClubJoinRequestList from "./AdminClubJoinRequestList";
import AdminClubPendingInviteList from "./AdminClubPendingInviteList";
import { Table, TableColumn, Input } from '@/components/ui';
import Modal from '@/components/ui/Modal';
import api from "../../../../utils/api";
import { useParams } from "react-router-dom";
import { useT } from "../../../../hooks/useTranslation";
import { useCurrentUser } from "../../../../hooks/useCurrentUser";
import { useOwnerCount } from "../../../../hooks/useOwnerCount";

interface Member {
    id: string;
    name: string;
    role: string;
    joinedAt: string;
    userId?: string; // Add userId to identify the current user
    birthDate?: string; // Add birth date field
}

const AdminClubMemberList = () => {
    const { t } = useT();
    const { id } = useParams();
    const { user: currentUser } = useCurrentUser();
    const { ownerCount, refetch: refetchOwnerCount } = useOwnerCount(id || '');

    const [members, setMembers] = useState<Member[]>([]);
    const [currentUserRole, setCurrentUserRole] = useState<string | null>(null);
    const [showManageInvites, setShowManageInvites] = useState(false);
    const [inviteEmail, setInviteEmail] = useState('');
    const [showJoinRequests, setShowJoinRequests] = useState(false);
    const [showInviteLink, setShowInviteLink] = useState(false);
    const [inviteLink, setInviteLink] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

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
            setLoading(true);
            setError(null);
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
        } finally {
            setLoading(false);
        }
    }, [id, sortMembersByRole, currentUser]);

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
            // Refetch owner count when roles change
            refetchOwnerCount();
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
            setInviteEmail(''); // Clear the email input after successful invite
            // Keep the modal open but refresh the pending invites list
            // The pending invites list will automatically refresh due to useEffect
        } catch {
            setError('Failed to send invite');
        }
    };

    // Permission logic based on backend rules (with desired admin permissions)
    const canChangeRole = (currentUserRole: string | null, targetMemberRole: string, newRole: string, targetMember: Member): boolean => {
        if (!currentUserRole) return false;
        
        // Check if current user is trying to demote themselves as the last owner
        if (currentUser && targetMember.userId === currentUser.ID && 
            targetMemberRole === 'owner' && newRole !== 'owner' && ownerCount <= 1) {
            return false; // Prevent last owner from demoting themselves
        }
        
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

    // Define table columns
    const columns: TableColumn<Member>[] = [
        {
            key: 'name',
            header: 'Name',
            render: (member) => member.name
        },
        {
            key: 'role',
            header: 'Role',
            render: (member) => translateRole(member.role)
        },
        {
            key: 'joined',
            header: 'Joined',
            render: (member) => member.joinedAt ? new Date(member.joinedAt).toLocaleDateString() : 'N/A'
        },
        {
            key: 'birthDate',
            header: 'Birth Date',
            render: (member) => member.birthDate ? new Date(member.birthDate).toLocaleDateString() : 'Not shared'
        },
        {
            key: 'actions',
            header: 'Actions',
            render: (member) => (
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
                    {member.role === 'member' && canChangeRole(currentUserRole, member.role, 'admin', member) && (
                        <button
                            onClick={() => handleRoleChange(member.id, 'admin')}
                            className="action-button promote"
                        >
                            Promote
                        </button>
                    )}
                    {member.role === 'admin' && (
                        <>
                            {canChangeRole(currentUserRole, member.role, 'member', member) && (
                                <button
                                    onClick={() => handleRoleChange(member.id, 'member')}
                                    className="action-button demote"
                                >
                                    Demote
                                </button>
                            )}
                            {canChangeRole(currentUserRole, member.role, 'owner', member) && (
                                <button
                                    onClick={() => handleRoleChange(member.id, 'owner')}
                                    className="action-button promote"
                                >
                                    Promote
                                </button>
                            )}
                        </>
                    )}
                    {member.role === 'owner' && canChangeRole(currentUserRole, member.role, 'admin', member) && (
                        <button
                            onClick={() => handleRoleChange(member.id, 'admin')}
                            className="action-button demote"
                        >
                            Demote
                        </button>
                    )}
                </div>
            )
        }
    ];

    if (error) return <div className="error">{error}</div>;

    return (
        <div>
            <h3>Members</h3>
            <Table
                columns={columns}
                data={members}
                keyExtractor={(member) => member.id}
                loading={loading}
                error={error}
                emptyMessage="No members found"
                loadingMessage="Loading members..."
                errorMessage="Failed to load members"
                footer={
                    members.length > 0 ? (
                        <div>
                            {t('clubs.totalMembers', { count: members.length }) || `Total: ${members.length} members`}
                        </div>
                    ) : null
                }
            />
            <div className="member-actions buttons" style={{ marginTop: 'var(--space-md)' }}>
                <button onClick={() => setShowManageInvites(true)} className="button-accept">Manage Invites</button>
                <button onClick={handleShowInviteLink} className="button-accept">Generate Invite Link</button>
                <button onClick={() => setShowJoinRequests(true)}>View Join Requests</button>
            </div>
            {showManageInvites && (
                <Modal 
                    isOpen={showManageInvites} 
                    onClose={() => setShowManageInvites(false)}
                    title="Manage Invites"
                >
                    <Modal.Body>
                        {/* Invite Member Section */}
                        <div style={{ marginBottom: '24px', paddingBottom: '16px', borderBottom: '1px solid #e0e0e0' }}>
                            <h4 style={{ marginBottom: '12px', fontSize: '1.1em', fontWeight: '600' }}>Invite New Member</h4>
                            <div className="modal-form-section">
                                <Input
                                    label="Email"
                                    id="inviteEmail"
                                    type="email"
                                    value={inviteEmail}
                                    onChange={(e) => setInviteEmail(e.target.value)}
                                    placeholder="Enter email"
                                />
                                <button 
                                    onClick={() => inviteEmail && sendInvite(inviteEmail)} 
                                    disabled={!inviteEmail} 
                                    className="button-accept"
                                    style={{ marginTop: '12px' }}
                                >
                                    Send Invite
                                </button>
                            </div>
                        </div>
                        
                        {/* Pending Invites Section */}
                        <div>
                            <h4 style={{ marginBottom: '12px', fontSize: '1.1em', fontWeight: '600' }}>Pending Invites</h4>
                            <AdminClubPendingInviteList />
                        </div>
                    </Modal.Body>
                    <Modal.Actions>
                        <button onClick={() => setShowManageInvites(false)} className="button-cancel">Close</button>
                    </Modal.Actions>
                </Modal>
            )}

            {showInviteLink && (
                <Modal 
                    isOpen={showInviteLink} 
                    onClose={() => setShowInviteLink(false)}
                    title="Club Invitation Link"
                >
                    <Modal.Body>
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
                        <div style={{ marginTop: '15px', fontSize: '0.9em', color: '#666' }}>
                            <p><strong>Note:</strong> Anyone with this link can request to join your club. 
                            Join requests still require admin approval.</p>
                        </div>
                    </Modal.Body>
                    <Modal.Actions>
                        <button onClick={copyToClipboard} className="button-accept">
                            Copy Link
                        </button>
                        <button onClick={() => setShowInviteLink(false)} className="button-cancel">
                            Close
                        </button>
                    </Modal.Actions>
                </Modal>
            )}
            
            {showJoinRequests && (
                <Modal 
                    isOpen={showJoinRequests} 
                    onClose={() => setShowJoinRequests(false)}
                    title="Join Requests"
                >
                    <Modal.Body>
                        <AdminClubJoinRequestList />
                    </Modal.Body>
                    <Modal.Actions>
                        <button onClick={() => setShowJoinRequests(false)} className="button-cancel">Close</button>
                    </Modal.Actions>
                </Modal>
            )}
        </div>
    );
};

export default AdminClubMemberList;