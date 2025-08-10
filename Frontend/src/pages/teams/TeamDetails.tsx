import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../utils/api';
import Layout from '../../components/layout/Layout';
import TeamNews from './TeamNews';
import TeamUpcomingEvents from './TeamUpcomingEvents';
import TeamFines from './TeamFines';
import TeamMembers from './TeamMembers';
import { useT } from '../../hooks/useTranslation';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import './TeamDetails.css';

interface Team {
    id: string;
    name: string;
    description: string;
    createdAt: string;
    clubId: string;
}

const TeamDetails = () => {
    const { t } = useT();
    const { clubId, teamId } = useParams();
    const navigate = useNavigate();
    const { user: currentUser } = useCurrentUser();
    const [team, setTeam] = useState<Team | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [teamNotFound, setTeamNotFound] = useState(false);
    const [isAdmin, setIsAdmin] = useState(false);
    const [userRole, setUserRole] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [teamResponse, adminResponse] = await Promise.all([
                    api.get(`/api/v1/clubs/${clubId}/teams/${teamId}`),
                    api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/isAdmin`)
                ]);
                const teamData = teamResponse.data;
                setTeam(teamData);
                setIsAdmin(adminResponse.data.isAdmin);

                // Get user's role by fetching team members and finding current user
                try {
                    const membersResponse = await api.get(`/api/v1/clubs/${clubId}/teams/${teamId}/members`);
                    const currentUserMember = membersResponse.data.find((member: { userId: string; role: string }) => member.userId === currentUser?.ID);
                    if (currentUserMember) {
                        setUserRole(currentUserMember.role);
                    }
                } catch (memberErr) {
                    // If we can't fetch members, we might not be a member or might not have access
                    console.warn('Could not fetch team member information:', memberErr);
                }
                
                setLoading(false);
            } catch (err: unknown) {
                console.error('Error fetching team details:', err);
                
                // Check if it's a 404 or 403 error (team not found or unauthorized)
                if (err && typeof err === 'object' && 'response' in err) {
                    const axiosError = err as { response?: { status?: number } };
                    if (axiosError.response?.status === 404 || axiosError.response?.status === 403) {
                        setTeamNotFound(true);
                    } else {
                        setError(t('teams.errors.loadingTeam') || 'Error fetching team details');
                    }
                } else {
                    setError(t('teams.errors.loadingTeam') || 'Error fetching team details');
                }
                setLoading(false);
            }
        };

        if (clubId && teamId) {
            fetchData();
        } else {
            setError('No club or team ID provided');
            setLoading(false);
        }
    }, [clubId, teamId, t, currentUser?.ID]);

    if (loading) return <div>Loading...</div>;
    if (teamNotFound) return <div className="error">Team not found or you don't have access to this team</div>;
    if (error) return <div className="error">{error}</div>;
    if (!team) return <div>Team not found</div>;

    return (
        <Layout title={team.name}>
            <div className="club-details-container">
                {/* Team Header */}
                <div className="club-header-section">
                    <div className="club-header-content">
                        {/* Team Logo Placeholder */}
                        <div className="club-logo-section">
                            <div className="club-logo-placeholder">
                                <span className="logo-placeholder-text">
                                    {team.name.charAt(0).toUpperCase()}
                                </span>
                            </div>
                        </div>

                        <div className="club-main-info">
                            <h1 className="club-title">{team.name}</h1>
                            {team.description && (
                                <p className="club-description">{team.description}</p>
                            )}
                            {userRole && (
                                <div className="user-role-container">
                                    <span className="role-label">Your role</span>
                                    <div className={`role-badge role-${userRole}`}>
                                        <span className="role-text">{userRole}</span>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                    
                    {/* Action Buttons */}
                    <div className="club-actions">
                        {isAdmin && (
                            <button 
                                className="button button-primary"
                                onClick={() => navigate(`/clubs/${clubId}/teams/${teamId}/admin`)}
                            >
                                Manage Team
                            </button>
                        )}
                        <button 
                            className="button button-cancel"
                            onClick={() => navigate(`/clubs/${clubId}`)}
                        >
                            Back to Club
                        </button>
                    </div>
                </div>

                {/* Content Sections */}
                <div className="club-content">
                    <TeamNews />
                    <TeamUpcomingEvents />
                    <TeamFines />
                    <TeamMembers />
                </div>
            </div>
        </Layout>
    );
};

export default TeamDetails;