import { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../../../utils/api';
import { useT } from '../../../../hooks/useTranslation';
import { Input, Modal } from '@/components/ui';

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
            const response = await api.get(`/api/v1/clubs/${clubId}/teams`);
            setTeams(response.data);
        } catch {
            setError('Failed to fetch teams');
        }
    }, [clubId]);

    const fetchClubMembers = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/members`);
            setClubMembers(response.data);
        } catch {
            setError('Failed to fetch club members');
        }
    }, [clubId]);

    const fetchTeamMembers = useCallback(async (teamId: string) => {
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/members`);
            // Ensure we always set an array, even if API returns null (prevents .map() crashes)
            setTeamMembers(response.data || []);
        } catch {
            setError('Failed to fetch team members');
            setTeamMembers([]); // Reset to empty array on error to prevent crashes
        }
    }, [clubId]);

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
            await api.post(`/api/v1/clubs/${clubId}/teams`, {
                name: newTeamName,
                description: newTeamDescription
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
            await api.put(`/api/v1/clubs/${clubId}/teams/${selectedTeam.id}`, {
                name: editTeamName,
                description: editTeamDescription
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
            await api.delete(`/api/v1/clubs/${clubId}/teams/${teamId}`);
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
            await api.post(`/api/v1/clubs/${clubId}/teams/${selectedTeam.id}/members`, {
                userId,
                role
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
            await api.patch(`/api/v1/clubs/${clubId}/teams/${selectedTeam.id}/members/${memberId}`, {
                role: newRole
            });
            await fetchTeamMembers(selectedTeam.id);
        } catch {
            setError('Failed to update member role');
        }
    };

    const handleRemoveMember = async (memberId: string) => {
        if (!selectedTeam || !confirm(t('teams.removeMemberConfirmation'))) return;

        try {
            await api.delete(`/api/v1/clubs/${clubId}/teams/${selectedTeam.id}/members/${memberId}`);
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
                        <button 
                            onClick={() => setShowCreateModal(true)} 
                            className="button-accept"
                        >
                            {t('teams.createTeam')}
                        </button>
                    </div>

                    <div className="teams-grid">
                        {teams.map(team => (
                            <div 
                                key={team.id} 
                                className={`team-card ${selectedTeam?.id === team.id ? 'selected' : ''}`}
                                onClick={() => setSelectedTeam(team)}
                            >
                                <div className="team-header">
                                    <h4>{team.name}</h4>
                                    <div className="team-actions">
                                        <button 
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                openEditModal(team);
                                            }}
                                            className="action-button edit"
                                        >
                                            {t('common.edit')}
                                        </button>
                                        <button 
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleDeleteTeam(team.id);
                                            }}
                                            className="action-button remove"
                                        >
                                            {t('common.delete')}
                                        </button>
                                    </div>
                                </div>
                                {team.description && <p className="team-description">{team.description}</p>}
                            </div>
                        ))}
                    </div>
                </div>

                {/* Team Members Section */}
                {selectedTeam && (
                    <div className="team-members-section">
                        <div className="section-header">
                            <h3>{t('teams.membersOf', { teamName: selectedTeam.name })}</h3>
                            <button 
                                onClick={() => setShowAddMemberModal(true)} 
                                className="button-accept"
                            >
                                {t('teams.addMember')}
                            </button>
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
                                            {t('teams.noMembers') || 'No team members yet.'}
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
                                                        <button
                                                            onClick={() => handleUpdateMemberRole(member.id, 'admin')}
                                                            className="action-button promote"
                                                        >
                                                            {t('teams.promoteToAdmin')}
                                                        </button>
                                                    )}
                                                    {member.role === 'admin' && (
                                                        <button
                                                            onClick={() => handleUpdateMemberRole(member.id, 'member')}
                                                            className="action-button demote"
                                                        >
                                                            {t('teams.demoteToMember')}
                                                        </button>
                                                    )}
                                                    <button
                                                        onClick={() => handleRemoveMember(member.id)}
                                                        className="action-button remove"
                                                    >
                                                        {t('common.remove')}
                                                    </button>
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
                    <button onClick={handleCreateTeam} className="button-accept">
                        {t('common.create')}
                    </button>
                    <button onClick={() => setShowCreateModal(false)} className="button-cancel">
                        {t('common.cancel')}
                    </button>
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
                    <button onClick={handleUpdateTeam} className="button-accept">
                        {t('common.save')}
                    </button>
                    <button onClick={() => setShowEditModal(false)} className="button-cancel">
                        {t('common.cancel')}
                    </button>
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
                                    <button
                                        onClick={() => handleAddMember(member.userId, 'member')}
                                        className="button-accept"
                                    >
                                        {t('teams.addAsMember')}
                                    </button>
                                    <button
                                        onClick={() => handleAddMember(member.userId, 'admin')}
                                        className="button-secondary"
                                    >
                                        {t('teams.addAsAdmin')}
                                    </button>
                                </div>
                            </div>
                        ))}
                        {getAvailableMembers().length === 0 && (
                            <p>{t('teams.noAvailableMembers')}</p>
                        )}
                    </div>
                </Modal.Body>
                <Modal.Actions>
                    <button onClick={() => setShowAddMemberModal(false)} className="button-cancel">
                        {t('common.cancel')}
                    </button>
                </Modal.Actions>
            </Modal>
        </div>
    );
};

export default AdminClubTeamList;
