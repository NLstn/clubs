import { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../../../utils/api';
import { useT } from '../../../../hooks/useTranslation';
import { Input, Modal, Button, Card } from '@/components/ui';
import './AdminClubTeamList.css';

interface Team {
    id: string;
    name: string;
    description: string;
    createdAt: string;
    clubId: string;
}

interface TeamMember {
    id: string;
    userId: string;
    name: string;
    role: 'admin' | 'member';
    joinedAt: string;
}

interface ClubMember {
    id: string;
    userId: string;
    name: string;
    role: string;
}

interface TeamMemberResponse {
    ID: string;
    UserID: string;
    Role: string;
    CreatedAt: string;
    User?: {
        FirstName: string;
        LastName: string;
    };
}

const AdminClubTeamList = () => {
    const { t } = useT();
    const { id: clubId } = useParams();
    const [teams, setTeams] = useState<Team[]>([]);
    const [selectedTeam, setSelectedTeam] = useState<Team | null>(null);
    const [teamMembers, setTeamMembers] = useState<TeamMember[]>([]);
    const [clubMembers, setClubMembers] = useState<ClubMember[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showAddMemberModal, setShowAddMemberModal] = useState(false);
    const [newTeamName, setNewTeamName] = useState('');
    const [newTeamDescription, setNewTeamDescription] = useState('');
    const [editTeamName, setEditTeamName] = useState('');
    const [editTeamDescription, setEditTeamDescription] = useState('');

    const fetchTeams = useCallback(async () => {
        try {
            // OData v2: Query Teams for this club
            const response = await api.get(`/api/v2/Teams?$filter=ClubID eq '${clubId}'`);
            setTeams(response.data);
        } catch {
            setError('Failed to fetch teams');
        }
    }, [clubId]);

    const fetchClubMembers = useCallback(async () => {
        try {
            // OData v2: Query Members for this club
            const response = await api.get(`/api/v2/Members?$filter=ClubID eq '${clubId}'&$expand=User`);
            setClubMembers(response.data);
        } catch {
            setError('Failed to fetch club members');
        }
    }, [clubId]);

    const fetchTeamMembers = useCallback(async (teamId: string) => {
        try {
            // OData v2: Use TeamMembers navigation with User expansion
            const response = await api.get(`/api/v2/Teams('${teamId}')/TeamMembers?$expand=User&$orderby=Role,User/FirstName`);
            const teamMembers = response.data.value || response.data;
            
            // Transform TeamMember entities with nested User to flat structure expected by UI
            const transformedMembers = teamMembers.map((tm: TeamMemberResponse) => ({
                id: tm.ID,
                userId: tm.UserID,
                role: tm.Role,
                joinedAt: tm.CreatedAt,
                name: tm.User ? `${tm.User.FirstName} ${tm.User.LastName}` : 'Unknown'
            }));
            
            // Ensure we always set an array, even if API returns null (prevents .map() crashes)
            setTeamMembers(transformedMembers || []);
        } catch {
            setError('Failed to fetch team members');
            setTeamMembers([]); // Reset to empty array on error to prevent crashes
        }
    }, []);

    useEffect(() => {
        const loadData = async () => {
            setLoading(true);
            await Promise.all([fetchTeams(), fetchClubMembers()]);
            setLoading(false);
        };
        loadData();
    }, [fetchTeams, fetchClubMembers]);

    useEffect(() => {
        if (selectedTeam) {
            // Calling fetchTeamMembers here is the correct pattern for data fetching
            // eslint-disable-next-line react-hooks/set-state-in-effect
            fetchTeamMembers(selectedTeam.id);
        }
    }, [selectedTeam, fetchTeamMembers]);

    const handleCreateTeam = async () => {
        if (!newTeamName.trim()) return;

        try {
            // OData v2: Create new team
            await api.post(`/api/v2/Teams`, {
                ClubID: clubId,
                Name: newTeamName,
                Description: newTeamDescription,
            });
            setNewTeamName('');
            setNewTeamDescription('');
            setShowCreateModal(false);
            await fetchTeams();
        } catch {
            setError('Failed to create team');
        }
    };

    const handleUpdateTeam = async () => {
        if (!selectedTeam || !editTeamName.trim()) return;

        try {
            // OData v2: Update team using PATCH
            await api.patch(`/api/v2/Teams('${selectedTeam.id}')`, {
                Name: editTeamName,
                Description: editTeamDescription,
            });
            setShowEditModal(false);
            await fetchTeams();
            // Update selected team
            setSelectedTeam(prev => prev ? { ...prev, name: editTeamName, description: editTeamDescription } : null);
        } catch {
            setError('Failed to update team');
        }
    };

    const handleDeleteTeam = async (teamId: string) => {
        if (!confirm(t('teams.deleteConfirmation'))) return;

        try {
            // OData v2: Delete team
            await api.delete(`/api/v2/Teams('${teamId}')`);
            await fetchTeams();
            if (selectedTeam?.id === teamId) {
                setSelectedTeam(null);
                setTeamMembers([]);
            }
        } catch {
            setError('Failed to delete team');
        }
    };

    const handleAddMember = async (userId: string, role: string) => {
        if (!selectedTeam) return;

        try {
            // OData v2: Create TeamMember entity directly
            await api.post(`/api/v2/TeamMembers`, {
                TeamID: selectedTeam.id,
                UserID: userId,
                Role: role
            });
            setShowAddMemberModal(false);
            await fetchTeamMembers(selectedTeam.id);
        } catch {
            setError('Failed to add team member');
        }
    };

    const handleUpdateMemberRole = async (memberId: string, newRole: string) => {
        if (!selectedTeam) return;

        try {
            // OData v2: Find TeamMember and update role
            const tmResponse = await api.get(`/api/v2/TeamMembers?$filter=TeamID eq '${selectedTeam.id}' and MemberID eq '${memberId}'`);
            const teamMember = tmResponse.data.value[0];
            if (teamMember) {
                await api.patch(`/api/v2/TeamMembers('${teamMember.ID}')`, {
                    Role: newRole
                });
            }
            await fetchTeamMembers(selectedTeam.id);
        } catch {
            setError('Failed to update member role');
        }
    };

    const handleRemoveMember = async (memberId: string) => {
        if (!selectedTeam || !confirm(t('teams.removeMemberConfirmation'))) return;

        try {
            // OData v2: Find TeamMember and delete
            const tmResponse = await api.get(`/api/v2/TeamMembers?$filter=TeamID eq '${selectedTeam.id}' and MemberID eq '${memberId}'`);
            const teamMember = tmResponse.data.value[0];
            if (teamMember) {
                await api.delete(`/api/v2/TeamMembers('${teamMember.ID}')`);
            }
            await fetchTeamMembers(selectedTeam.id);
        } catch {
            setError('Failed to remove team member');
        }
    };

    const openEditModal = (team: Team) => {
        setSelectedTeam(team);
        setEditTeamName(team.name);
        setEditTeamDescription(team.description);
        setShowEditModal(true);
    };

    const getAvailableMembers = () => {
        // Use defensive programming to prevent crashes if teamMembers is null/undefined
        const teamMemberIds = (teamMembers || []).map(tm => tm.userId);
        return clubMembers.filter(cm => !teamMemberIds.includes(cm.userId));
    };

    if (loading) return <div>Loading teams...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="teams-management">
            <div className="teams-layout">
                {/* Teams List */}
                <div className="teams-list-section">
                    <div className="section-header">
                        <h3>{t('teams.title')}</h3>
                        <Button 
                            variant="accept"
                            onClick={() => setShowCreateModal(true)}
                        >
                            {t('teams.createTeam')}
                        </Button>
                    </div>

                    <div className="teams-grid">
                        {teams.map(team => (
                            <Card
                                key={team.id}
                                variant="light"
                                padding="md"
                                clickable
                                hover
                                onClick={() => setSelectedTeam(team)}
                                className={`team-card ${selectedTeam?.id === team.id ? 'selected' : ''}`}
                            >
                                <div className="team-header">
                                    <h4>{team.name}</h4>
                                    <div className="team-actions">
                                        <Button 
                                            size="sm"
                                            variant="secondary"
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                openEditModal(team);
                                            }}
                                        >
                                            {t('common.edit')}
                                        </Button>
                                        <Button 
                                            size="sm"
                                            variant="cancel"
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleDeleteTeam(team.id);
                                            }}
                                        >
                                            {t('common.delete')}
                                        </Button>
                                    </div>
                                </div>
                                {team.description && <p className="team-description">{team.description}</p>}
                            </Card>
                        ))}
                    </div>
                </div>

                {/* Team Members Section */}
                {selectedTeam && (
                    <div className="team-members-section">
                        <div className="section-header">
                            <h3>{t('teams.membersOf', { teamName: selectedTeam.name })}</h3>
                            <Button 
                                variant="accept"
                                onClick={() => setShowAddMemberModal(true)}
                            >
                                {t('teams.addMember')}
                            </Button>
                        </div>

                        <table className="members-table">
                            <thead>
                                <tr>
                                    <th>{t('common.name')}</th>
                                    <th>{t('common.role')}</th>
                                    <th>{t('teams.joinedAt')}</th>
                                    <th>{t('common.actions')}</th>
                                </tr>
                            </thead>
                            <tbody>
                                {teamMembers.length === 0 ? (
                                    <tr>
                                        <td colSpan={4} style={{textAlign: 'center', fontStyle: 'italic', padding: 'var(--space-lg)'}}>
                                            {t('teams.noMembers')}
                                        </td>
                                    </tr>
                                ) : (
                                    teamMembers.map(member => (
                                        <tr key={member.id}>
                                            <td>{member.name}</td>
                                            <td>{t(`teams.roles.${member.role}`)}</td>
                                            <td>{new Date(member.joinedAt).toLocaleDateString()}</td>
                                            <td>
                                                <div className="member-actions">
                                                    {member.role === 'member' && (
                                                        <Button
                                                            onClick={() => handleUpdateMemberRole(member.id, 'admin')}
                                                            variant="primary"
                                                            size="sm"
                                                        >
                                                            {t('teams.promoteToAdmin')}
                                                        </Button>
                                                    )}
                                                    {member.role === 'admin' && (
                                                        <Button
                                                            onClick={() => handleUpdateMemberRole(member.id, 'member')}
                                                            variant="maybe"
                                                            size="sm"
                                                        >
                                                            {t('teams.demoteToMember')}
                                                        </Button>
                                                    )}
                                                    <Button
                                                        onClick={() => handleRemoveMember(member.id)}
                                                        variant="cancel"
                                                        size="sm"
                                                    >
                                                        {t('common.remove')}
                                                    </Button>
                                                </div>
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>

            {/* Create Team Modal */}
            <Modal isOpen={showCreateModal} onClose={() => setShowCreateModal(false)} title={t('teams.createTeam')}>
                <Modal.Body>
                    <div className="modal-form-section">
                        <Input
                            label={t('teams.teamName')}
                            type="text"
                            value={newTeamName}
                            onChange={(e) => setNewTeamName(e.target.value)}
                            placeholder={t('teams.teamNamePlaceholder')}
                        />
                        <Input
                            label={t('teams.description')}
                            value={newTeamDescription}
                            onChange={(e) => setNewTeamDescription(e.target.value)}
                            placeholder={t('teams.descriptionPlaceholder')}
                            multiline
                            rows={3}
                        />
                    </div>
                </Modal.Body>
                <Modal.Actions>
                    <Button variant="accept" onClick={handleCreateTeam}>
                        {t('common.create')}
                    </Button>
                    <Button variant="cancel" onClick={() => setShowCreateModal(false)}>
                        {t('common.cancel')}
                    </Button>
                </Modal.Actions>
            </Modal>

            {/* Edit Team Modal */}
            <Modal isOpen={showEditModal} onClose={() => setShowEditModal(false)} title={t('teams.editTeam')}>
                <Modal.Body>
                    <div className="modal-form-section">
                        <Input
                            label={t('teams.teamName')}
                            type="text"
                            value={editTeamName}
                            onChange={(e) => setEditTeamName(e.target.value)}
                            placeholder={t('teams.teamNamePlaceholder')}
                        />
                        <Input
                            label={t('teams.description')}
                            value={editTeamDescription}
                            onChange={(e) => setEditTeamDescription(e.target.value)}
                            placeholder={t('teams.descriptionPlaceholder')}
                            multiline
                            rows={3}
                        />
                    </div>
                </Modal.Body>
                <Modal.Actions>
                    <Button variant="accept" onClick={handleUpdateTeam}>
                        {t('common.save')}
                    </Button>
                    <Button variant="cancel" onClick={() => setShowEditModal(false)}>
                        {t('common.cancel')}
                    </Button>
                </Modal.Actions>
            </Modal>

            {/* Add Member Modal */}
            <Modal isOpen={showAddMemberModal} onClose={() => setShowAddMemberModal(false)} title={t('teams.addMember')} maxWidth="600px">
                <Modal.Body>
                    <div className="available-members">
                        {getAvailableMembers().map(member => (
                            <div key={member.id} className="member-option">
                                <span>{member.name}</span>
                                <div className="role-actions">
                                    <Button
                                        variant="accept"
                                        size="sm"
                                        onClick={() => handleAddMember(member.userId, 'member')}
                                    >
                                        {t('teams.addAsMember')}
                                    </Button>
                                    <Button
                                        variant="secondary"
                                        size="sm"
                                        onClick={() => handleAddMember(member.userId, 'admin')}
                                    >
                                        {t('teams.addAsAdmin')}
                                    </Button>
                                </div>
                            </div>
                        ))}
                        {getAvailableMembers().length === 0 && (
                            <p>{t('teams.noAvailableMembers')}</p>
                        )}
                    </div>
                </Modal.Body>
                <Modal.Actions>
                    <Button variant="cancel" onClick={() => setShowAddMemberModal(false)}>
                        {t('common.cancel')}
                    </Button>
                </Modal.Actions>
            </Modal>
        </div>
    );
};

export default AdminClubTeamList;
