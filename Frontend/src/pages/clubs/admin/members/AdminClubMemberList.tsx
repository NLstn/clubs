import { useEffect, useState, useCallback, useRef } from "react";
import AdminClubJoinRequestList from "./AdminClubJoinRequestList";
import AdminClubPendingInviteList from "./AdminClubPendingInviteList";
import { Table, TableColumn, Input, Button, ButtonState } from '@/components/ui';
import Modal from '@/components/ui/Modal';
import api from "../../../../utils/api";
import { parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';
import { useParams, useSearchParams } from "react-router-dom";
import { useT } from "../../../../hooks/useTranslation";
import { useCurrentUser } from "../../../../hooks/useCurrentUser";
import { useOwnerCount } from "../../../../hooks/useOwnerCount";
import './AdminClubMemberList.css';

interface Member {
    id: string;
    name: string;
    role: string;
    joinedAt: string;
    userId?: string; // Add userId to identify the current user
    birthDate?: string; // Add birth date field
}

interface MemberResponse {
    ID: string;
    ClubID: string;
    UserID: string;
    Role: string;
    CreatedAt: string;
    User?: {
        FirstName: string;
        LastName: string;
        BirthDate?: string;
    };
}

interface AdminClubMemberListProps {
    openJoinRequests?: boolean;
}

const AdminClubMemberList = ({ openJoinRequests = false }: AdminClubMemberListProps) => {
    const { t } = useT();
    const { id } = useParams();
    const [searchParams, setSearchParams] = useSearchParams();
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
    const [inviteLinkError, setInviteLinkError] = useState<string | null>(null);
    const [sendInviteError, setSendInviteError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [joinRequestCount, setJoinRequestCount] = useState<number>(0);
    const [memberActions, setMemberActions] = useState<Record<string, ButtonState>>({});
    const timeoutRefs = useRef<number[]>([]);

    useEffect(() => {
        // Cleanup timeouts on unmount
        const timeouts = timeoutRefs.current;
        return () => {
            timeouts.forEach(clearTimeout);
        };
    }, []);

    const translateRole = (role: string): string => {
        return t(`clubs.roles.${role}`);
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
            // OData v2: Query Members with expanded User data
            const response = await api.get<ODataCollectionResponse<MemberResponse>>(`/api/v2/Members?$filter=ClubID eq '${id}'&$expand=User`);
            const odataMembers = parseODataCollection(response.data);
            
            // Transform OData response (PascalCase) to frontend interface (camelCase)
            const transformedMembers = odataMembers.map((m: MemberResponse) => ({
                id: m.ID,
                userId: m.UserID,
                role: m.Role,
                joinedAt: m.CreatedAt,
                name: m.User ? `${m.User.FirstName} ${m.User.LastName}` : 'Unknown',
                birthDate: m.User?.BirthDate
            }));
            
            const sortedMembers = sortMembersByRole(transformedMembers);
            setMembers(sortedMembers);

            // Set the current user's role if available
            const currentUserMember = transformedMembers.find((member: Member) => member.userId === currentUser?.ID);
            if (currentUserMember) {
                setCurrentUserRole(currentUserMember.role);
            }
        } catch {
            setError("Failed to fetch members");
        } finally {
            setLoading(false);
        }
    }, [id, sortMembersByRole, currentUser]);

    const fetchJoinRequestCount = useCallback(async () => {
        try {
            // OData v2: Use $count to get join requests count
            const response = await api.get(`/api/v2/JoinRequests/$count?$filter=ClubID eq '${id}'`);
            const count = typeof response.data === 'number' ? response.data : parseInt(response.data, 10);

            setJoinRequestCount(Number.isFinite(count) ? count : 0);
        } catch (error) {
            console.error("Error fetching join request count:", error);
            setJoinRequestCount(0);
        }
    }, [id]);

    const handleShowInviteLink = async () => {
        setInviteLinkError(null); // Clear any previous errors
        try {
            // OData v2: Use GetInviteLink function on Club
            const response = await api.get(`/api/v2/Clubs('${id}')/GetInviteLink()`);
            const fullLink = `${window.location.origin}${response.data.InviteLink}`;
            setInviteLink(fullLink);
            setShowInviteLink(true);
        } catch (error) {
            console.error("Error fetching invite link:", error);
            setInviteLinkError('Failed to generate invite link. Please try again.');
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
        fetchJoinRequestCount();
    }, [fetchMembers, fetchJoinRequestCount]);

    // Automatically open join requests modal if specified via URL parameter
    useEffect(() => {
        if (openJoinRequests && !loading) {
            setShowJoinRequests(true);
        }
    }, [openJoinRequests, loading]);

    const handleCloseJoinRequests = () => {
        setShowJoinRequests(false);
        
        // Clear URL parameters when closing the modal
        if (searchParams.get('openJoinRequests') === 'true') {
            const newSearchParams = new URLSearchParams(searchParams);
            newSearchParams.delete('openJoinRequests');
            setSearchParams(newSearchParams);
        }
    };

    const handleRoleChange = async (memberId: string, newRole: string, actionKey: string) => {
        setMemberActions(prev => ({ ...prev, [actionKey]: 'loading' }));
        
        try {
            // OData v2: Use UpdateRole action on Member entity
            await api.post(`/api/v2/Members('${memberId}')/UpdateRole`, { newRole });
            const updatedMembers = members.map(member => 
                member.id === memberId ? { ...member, role: newRole } : member
            );
            const sortedMembers = sortMembersByRole(updatedMembers);
            setMembers(sortedMembers);
            // Refetch owner count when roles change
            refetchOwnerCount();
            
            setMemberActions(prev => ({ ...prev, [actionKey]: 'success' }));
            const timeoutId = window.setTimeout(() => {
                setMemberActions(prev => {
                    const newState = { ...prev };
                    delete newState[actionKey];
                    return newState;
                });
            }, 1000);
            timeoutRefs.current.push(timeoutId);
        } catch {
            // Show error on the button itself, don't replace the whole table
            setMemberActions(prev => ({ ...prev, [actionKey]: 'error' }));
            const timeoutId = window.setTimeout(() => {
                setMemberActions(prev => {
                    const newState = { ...prev };
                    delete newState[actionKey];
                    return newState;
                });
            }, 3000);
            timeoutRefs.current.push(timeoutId);
        }
    };

    const deleteMember = async (memberId: string) => {
        setMemberActions(prev => ({ ...prev, [`delete-${memberId}`]: 'loading' }));
        
        try {
            // OData v2: Delete member using DELETE
            await api.delete(`/api/v2/Members('${memberId}')`);
            setMemberActions(prev => ({ ...prev, [`delete-${memberId}`]: 'success' }));
            
            const timeoutId = window.setTimeout(() => {
                setMembers(members.filter(member => member.id !== memberId));
                setMemberActions(prev => {
                    const newState = { ...prev };
                    delete newState[`delete-${memberId}`];
                    return newState;
                });
            }, 1000);
            timeoutRefs.current.push(timeoutId);
        } catch {
            // Show error on the button itself, don't replace the whole table
            setMemberActions(prev => ({ ...prev, [`delete-${memberId}`]: 'error' }));
            const timeoutId = window.setTimeout(() => {
                setMemberActions(prev => {
                    const newState = { ...prev };
                    delete newState[`delete-${memberId}`];
                    return newState;
                });
            }, 3000);
            timeoutRefs.current.push(timeoutId);
        }
    };

    const sendInvite = async (email: string) => {
        setSendInviteError(null); // Clear any previous errors
        try {
            // OData v2: Use CreateInvite action on Club entity
            await api.post(`/api/v2/Clubs('${id}')/CreateInvite`, { email });
            setInviteEmail(''); // Clear the email input after successful invite
            // Keep the modal open but refresh the pending invites list
            // The pending invites list will automatically refresh due to useEffect
        } catch {
            setSendInviteError('Failed to send invite. Please try again.');
        }
    };

    // Permission logic based on backend rules (with desired admin permissions)
    const canChangeRole = (currentUserRole: string | null, targetMemberRole: string, newRole: string): boolean => {
        if (!currentUserRole) return false;
        
        // Check if trying to demote the last owner (prevents demoting ANY owner when they're the last one)
        if (targetMemberRole === 'owner' && newRole !== 'owner' && ownerCount <= 1) {
            return false; // Prevent last owner from being demoted
        }
        
        // Owners can change any role to any role (except demoting last owner, checked above)
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
                        <Button
                            size="sm"
                            variant="cancel"
                            onClick={() => deleteMember(member.id)}
                            aria-label="Remove member"
                            state={memberActions[`delete-${member.id}`] || 'idle'}
                            successMessage="Removed"
                            errorMessage="Failed"
                        >
                            Remove
                        </Button>
                    )}
                    {member.role === 'member' && canChangeRole(currentUserRole, member.role, 'admin') && (
                        <Button
                            size="sm"
                            variant="secondary"
                            onClick={() => handleRoleChange(member.id, 'admin', `promote-admin-${member.id}`)}
                            state={memberActions[`promote-admin-${member.id}`] || 'idle'}
                            successMessage="Promoted!"
                            errorMessage="Failed"
                        >
                            Promote
                        </Button>
                    )}
                    {member.role === 'admin' && (
                        <>
                            {canChangeRole(currentUserRole, member.role, 'member') && (
                                <Button
                                    size="sm"
                                    variant="secondary"
                                    onClick={() => handleRoleChange(member.id, 'member', `demote-member-${member.id}`)}
                                    state={memberActions[`demote-member-${member.id}`] || 'idle'}
                                    successMessage="Demoted"
                                    errorMessage="Failed"
                                >
                                    Demote
                                </Button>
                            )}
                            {canChangeRole(currentUserRole, member.role, 'owner') && (
                                <Button
                                    size="sm"
                                    variant="secondary"
                                    onClick={() => handleRoleChange(member.id, 'owner', `promote-owner-${member.id}`)}
                                    state={memberActions[`promote-owner-${member.id}`] || 'idle'}
                                    successMessage="Promoted!"
                                    errorMessage="Failed"
                                >
                                    Promote
                                </Button>
                            )}
                        </>
                    )}
                    {member.role === 'owner' && canChangeRole(currentUserRole, member.role, 'admin') && (
                        <Button
                            size="sm"
                            variant="secondary"
                            onClick={() => handleRoleChange(member.id, 'admin', `demote-admin-${member.id}`)}
                            state={memberActions[`demote-admin-${member.id}`] || 'idle'}
                            successMessage="Demoted"
                            errorMessage="Failed"
                        >
                            Demote
                        </Button>
                    )}
                </div>
            )
        }
    ];

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
                            {t('clubs.totalMembers', { count: members.length })}
                        </div>
                    ) : null
                }
            />
            <div className="member-actions buttons" style={{ marginTop: 'var(--space-md)' }}>
                <Button onClick={() => setShowManageInvites(true)} variant="accept">Manage Invites</Button>
                <Button onClick={handleShowInviteLink} variant="accept">Generate Invite Link</Button>
                <Button 
                    onClick={() => setShowJoinRequests(true)}
                    variant="primary"
                    counter={joinRequestCount}
                >
                    View Join Requests
                </Button>
            </div>
            {showManageInvites && (
                <Modal 
                    isOpen={showManageInvites} 
                    onClose={() => setShowManageInvites(false)}
                    title="Manage Invites"
                >
                    <Modal.Error error={sendInviteError} />
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
                                <Button 
                                    variant="accept"
                                    onClick={() => inviteEmail && sendInvite(inviteEmail)}
                                    disabled={!inviteEmail}
                                    style={{ marginTop: '12px' }}
                                >
                                    Send Invite
                                </Button>
                            </div>
                        </div>
                        
                        {/* Pending Invites Section */}
                        <div>
                            <h4 style={{ marginBottom: '12px', fontSize: '1.1em', fontWeight: '600' }}>Pending Invites</h4>
                            <AdminClubPendingInviteList />
                        </div>
                    </Modal.Body>
                    <Modal.Actions>
                        <Button variant="cancel" onClick={() => setShowManageInvites(false)}>Close</Button>
                    </Modal.Actions>
                </Modal>
            )}

            {showInviteLink && (
                <Modal 
                    isOpen={showInviteLink} 
                    onClose={() => setShowInviteLink(false)}
                    title="Club Invitation Link"
                >
                    <Modal.Error error={inviteLinkError} />
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
                        <Button variant="accept" onClick={copyToClipboard}>
                            Copy Link
                        </Button>
                        <Button variant="cancel" onClick={() => setShowInviteLink(false)}>
                            Close
                        </Button>
                    </Modal.Actions>
                </Modal>
            )}
            
            {showJoinRequests && (
                <Modal 
                    isOpen={showJoinRequests} 
                    onClose={handleCloseJoinRequests}
                    title="Join Requests"
                >
                    <Modal.Body>
                        <AdminClubJoinRequestList onRequestsChange={fetchJoinRequestCount} />
                    </Modal.Body>
                    <Modal.Actions>
                        <Button variant="cancel" onClick={handleCloseJoinRequests}>Close</Button>
                    </Modal.Actions>
                </Modal>
            )}
        </div>
    );
};

export default AdminClubMemberList;