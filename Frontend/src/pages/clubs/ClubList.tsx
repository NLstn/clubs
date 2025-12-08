import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../../components/layout/Layout';
import { Button, Card } from '../../components/ui';
import api from '../../utils/api';
import { addRecentClub } from '../../utils/recentClubs';
import { useT } from '../../hooks/useTranslation';
import './ClubList.css';

interface Club {
    id: string;
    name: string;
    description: string;
    user_role: string;
    created_at: string;
    deleted?: boolean;
    user_teams?: Team[];
}

interface Team {
    id: string;
    name: string;
    description: string;
    createdAt: string;
    clubId: string;
}

const ClubList = () => {
    const { t } = useT();
    const [clubs, setClubs] = useState<Club[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    const translateRole = (role: string): string => {
        return t(`clubs.roles.${role}`);
    };

    useEffect(() => {
        fetchClubs();
    }, []);

    const fetchClubs = async () => {
        try {
            const response = await api.get('/api/v1/clubs');
            setClubs(response.data || []);
        } catch (err: Error | unknown) {
            console.error('Error fetching clubs:', err);
            setError('Failed to fetch clubs');
        } finally {
            setLoading(false);
        }
    };

    const handleClubClick = (clubId: string, clubName: string) => {
        // Add to recent clubs when clicking
        addRecentClub(clubId, clubName);
        navigate(`/clubs/${clubId}`);
    };

    const handleTeamClick = (e: React.MouseEvent, clubId: string, teamId: string) => {
        e.stopPropagation(); // Prevent club card click
        navigate(`/clubs/${clubId}/teams/${teamId}`);
    };

    const adminClubs = clubs.filter(club => club.user_role === 'owner' || club.user_role === 'admin');
    const memberClubs = clubs.filter(club => club.user_role === 'member');

    if (loading) {
        return (
            <Layout title="My Clubs">
                <div>Loading clubs...</div>
            </Layout>
        );
    }

    if (error) {
        return (
            <Layout title="My Clubs">
                <div className="error">{error}</div>
            </Layout>
        );
    }

    return (
        <Layout title="My Clubs" showRecentClubs={false}>
            <div className="clubs-container">
                {adminClubs.length > 0 && (
                    <div className="clubs-section">
                        <h2>Clubs I Manage</h2>
                        <div className="clubs-grid">
                            {adminClubs.map(club => (
                                <Card
                                    key={club.id}
                                    variant="light"
                                    padding="lg"
                                    clickable
                                    hover
                                    onClick={() => handleClubClick(club.id, club.name)}
                                    className="club-card"
                                >
                                    <div className="club-header">
                                        <h3>{club.name}</h3>
                                        <span className={`role-badge ${club.user_role}`}>
                                            {translateRole(club.user_role)}
                                        </span>
                                    </div>
                                    <p className="club-description">{club.description}</p>
                                    {club.deleted && (
                                        <div className="club-deleted-badge">
                                            Deleted
                                        </div>
                                    )}
                                    {club.user_teams && club.user_teams.length > 0 && (
                                        <div className="club-teams-section">
                                            <h4 className="teams-title">My Teams</h4>
                                            <div className="teams-list">
                                                {club.user_teams.map(team => (
                                                    <div 
                                                        key={team.id}
                                                        className="team-badge"
                                                        onClick={(e) => handleTeamClick(e, club.id, team.id)}
                                                        title={team.description}
                                                    >
                                                        {team.name}
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    )}
                                </Card>
                            ))}
                        </div>
                    </div>
                )}

                {memberClubs.length > 0 && (
                    <div className="clubs-section">
                        <h2>Clubs I'm a Member Of</h2>
                        <div className="clubs-grid">
                            {memberClubs.map(club => (
                                <Card
                                    key={club.id}
                                    variant="light"
                                    padding="lg"
                                    clickable
                                    hover
                                    onClick={() => handleClubClick(club.id, club.name)}
                                    className="club-card"
                                >
                                    <div className="club-header">
                                        <h3>{club.name}</h3>
                                        <span className={`role-badge ${club.user_role}`}>
                                            {translateRole(club.user_role)}
                                        </span>
                                    </div>
                                    <p className="club-description">{club.description}</p>
                                    {club.user_teams && club.user_teams.length > 0 && (
                                        <div className="club-teams-section">
                                            <h4 className="teams-title">My Teams</h4>
                                            <div className="teams-list">
                                                {club.user_teams.map(team => (
                                                    <div 
                                                        key={team.id}
                                                        className="team-badge"
                                                        onClick={(e) => handleTeamClick(e, club.id, team.id)}
                                                        title={team.description}
                                                    >
                                                        {team.name}
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    )}
                                </Card>
                            ))}
                        </div>
                    </div>
                )}

                {adminClubs.length === 0 && memberClubs.length === 0 && (
                    <div className="empty-state">
                        <h2>No Clubs Yet</h2>
                        <p>You're not a member of any clubs yet.</p>
                        <Button 
                            onClick={() => navigate('/createClub')}
                            variant="primary"
                        >
                            Create Your First Club
                        </Button>
                    </div>
                )}
            </div>
        </Layout>
    );
};

export default ClubList;