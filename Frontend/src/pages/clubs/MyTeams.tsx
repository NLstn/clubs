import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../utils/api';
import { useCurrentUser } from '../../hooks/useCurrentUser';
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
                const response = await api.get(`/api/v1/clubs/${clubId}/teams?user=${currentUser.ID}`);
                setTeams(response.data || []);
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

    if (loading) return <div className="loading-text">Loading teams...</div>;
    if (error) return <div className="error-text">{error}</div>;
    if (teams.length === 0) return null; // Don't show anything if user has no teams

    return (
        <div className="content-section">
            <h3>{t('teams.myTeams') || 'My Teams'}</h3>
            <div className="teams-list">
                {teams.map(team => (
                    <div key={team.id} className="team-item">
                        <h4 className="team-name">{team.name}</h4>
                        {team.description && (
                            <p className="team-description">{team.description}</p>
                        )}
                    </div>
                ))}
            </div>
        </div>
    );
};

export default MyTeams;
