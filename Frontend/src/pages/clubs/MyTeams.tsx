import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import { useT } from '../../hooks/useTranslation';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import { parseODataCollection, type ODataCollectionResponse } from '@/utils/odata';

interface Team {
    ID: string;
    Name: string;
    Description: string;
    CreatedAt: string;
    ClubID: string;
}

interface TeamMemberResponse {
    ID: string;
    TeamID: string;
    UserID: string;
    Role: string;
    Team?: Team;
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
                // OData v2: Single query with $expand to get teams for this user and club
                const encodedUserId = encodeURIComponent(currentUser.ID);
                const encodedClubId = encodeURIComponent(clubId);
                
                const response = await api.get<ODataCollectionResponse<TeamMemberResponse>>(
                    `/api/v2/TeamMembers?$filter=UserID eq '${encodedUserId}'&$expand=Team($filter=ClubID eq '${encodedClubId}')`
                );
                const teamMembersData = parseODataCollection(response.data);
                
                // Extract teams from the expanded navigation property, filtering out null/undefined Teams
                const userTeams = teamMembersData
                    .map((tm: TeamMemberResponse) => tm.Team)
                    .filter((team): team is Team => team !== null && team !== undefined);
                
                setTeams(userTeams);
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
                        key={team.ID} 
                        className="team-item clickable"
                        onClick={() => handleTeamClick(team.ID)}
                        style={{ cursor: 'pointer' }}
                    >
                        <h4 className="team-name">{team.Name}</h4>
                        {team.Description && (
                            <p className="team-description">{team.Description}</p>
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
