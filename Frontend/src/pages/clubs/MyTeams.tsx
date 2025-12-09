import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import { useT } from '../../hooks/useTranslation';

interface Team {
    id: string;
    name: string;
    description: string;
    createdAt: string;
    clubId: string;
}

const MyTeams = () => {
    const { t } = useT();
    const { id: clubId } = useParams();
    const navigate = useNavigate();
    const [teams, setTeams] = useState<Team[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchUserTeams = async () => {
            if (!clubId) {
                setLoading(false);
                return;
            }

            try {
                // OData v2: Query Teams through Members relationship
                // Get current user's team memberships for this club
                const response = await api.get(
                    `/api/v2/Members?$filter=ClubID eq '${clubId}'&$expand=Teams($expand=Team)&$select=Teams`
                );
                const members = response.data.value || [];
                // Extract unique teams from the member's team relationships
                interface ODataTeam { ID: string; Name: string; Description: string; CreatedAt: string; ClubID: string; }
                interface ODataTeamMember { Team: ODataTeam; }
                interface ODataMember { Teams?: ODataTeamMember[]; }
                const teamsSet = new Set<string>();
                const teamsMap = new Map<string, Team>();
                members.forEach((member: ODataMember) => {
                    member.Teams?.forEach((teamMember: ODataTeamMember) => {
                        const team = teamMember.Team;
                        if (team && !teamsSet.has(team.ID)) {
                            teamsSet.add(team.ID);
                            teamsMap.set(team.ID, {
                                id: team.ID,
                                name: team.Name,
                                description: team.Description,
                                createdAt: team.CreatedAt,
                                clubId: team.ClubID
                            });
                        }
                    });
                });
                setTeams(Array.from(teamsMap.values()));
                setError(null);
            } catch (err) {
                console.error('Error fetching user teams:', err);
                // If it's a 403 error, user might not have access or teams might be disabled
                if (err && typeof err === 'object' && 'response' in err) {
                    const axiosError = err as { response?: { status?: number } };
                    if (axiosError.response?.status === 403) {
                        // User doesn't have access, don't show error
                        setTeams([]);
                        setError(null);
                        return;
                    }
                }
                setError('Failed to fetch teams');
                setTeams([]);
            } finally {
                setLoading(false);
            }
        };

        fetchUserTeams();
    }, [clubId]);

    const handleTeamClick = (teamId: string) => {
        navigate(`/clubs/${clubId}/teams/${teamId}`);
    };

    if (loading) return <div className="loading-text">Loading teams...</div>;
    if (error) return <div className="error-text">{error}</div>;
    if (teams.length === 0) return null; // Don't show anything if user has no teams

    return (
        <div className="content-section">
            <h3>{t('teams.myTeams')}</h3>
            <div className="teams-list">
                {teams.map(team => (
                    <div 
                        key={team.id} 
                        className="team-item clickable"
                        onClick={() => handleTeamClick(team.id)}
                        style={{ cursor: 'pointer' }}
                    >
                        <h4 className="team-name">{team.name}</h4>
                        {team.description && (
                            <p className="team-description">{team.description}</p>
                        )}
                        <div className="team-actions">
                            <span className="team-link">View Team →</span>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default MyTeams;
