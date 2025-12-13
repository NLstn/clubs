import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import { useT } from '../../hooks/useTranslation';
import { useCurrentUser } from '../../hooks/useCurrentUser';

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
    const { user: currentUser } = useCurrentUser();
    const [teams, setTeams] = useState<Team[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchUserTeams = async () => {
            if (!clubId || !currentUser?.ID) {
                setLoading(false);
                return;
            }

            try {
                // OData v2: Query Teams with filter for current user membership
                // Using lambda 'any' operator to filter teams where user is a member
                const response = await api.get(
                    `/api/v2/Teams?$filter=ClubID eq '${clubId}' and TeamMembers/any(tm: tm/UserID eq '${currentUser.ID}')`
                );
                const teamsData = response.data.value || [];
                
                // Map OData response to match expected format
                interface ODataTeam { ID: string; Name: string; Description: string; CreatedAt: string; ClubID: string; }
                const mappedTeams = teamsData.map((team: ODataTeam) => ({
                    id: team.ID,
                    name: team.Name,
                    description: team.Description,
                    createdAt: team.CreatedAt,
                    clubId: team.ClubID
                }));
                setTeams(mappedTeams);
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
    }, [clubId, currentUser?.ID]);

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
